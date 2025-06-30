package storage

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// FileLock provides cross-platform file locking using channels
type FileLock struct {
	locks   map[string]chan struct{} // Map of file paths to lock channels
	mutex   sync.RWMutex             // Protects the locks map
	timeout time.Duration            // Lock acquisition timeout
}

// NewFileLock creates a new file lock manager
func NewFileLock(timeout time.Duration) *FileLock {
	return &FileLock{
		locks:   make(map[string]chan struct{}),
		timeout: timeout,
	}
}

// Lock acquires a lock for the given file path
func (fl *FileLock) Lock(ctx context.Context, filePath string) (*FileLockHandle, error) {
	fl.mutex.Lock()
	lockChan, exists := fl.locks[filePath]
	if !exists {
		// Create a new lock channel (buffered with capacity 1 for the lock holder)
		lockChan = make(chan struct{}, 1)
		fl.locks[filePath] = lockChan
	}
	fl.mutex.Unlock()

	// Create a context with timeout
	lockCtx, cancel := context.WithTimeout(ctx, fl.timeout)
	defer cancel()

	// Try to acquire the lock
	select {
	case lockChan <- struct{}{}:
		// Lock acquired successfully
		return &FileLockHandle{
			filePath: filePath,
			lockChan: lockChan,
			fileLock: fl,
		}, nil
	case <-lockCtx.Done():
		return nil, fmt.Errorf("failed to acquire lock for %s: %w", filePath, lockCtx.Err())
	}
}

// TryLock attempts to acquire a lock without blocking
func (fl *FileLock) TryLock(filePath string) (*FileLockHandle, error) {
	fl.mutex.Lock()
	lockChan, exists := fl.locks[filePath]
	if !exists {
		// Create a new lock channel
		lockChan = make(chan struct{}, 1)
		fl.locks[filePath] = lockChan
	}
	fl.mutex.Unlock()

	// Try to acquire the lock without blocking
	select {
	case lockChan <- struct{}{}:
		// Lock acquired successfully
		return &FileLockHandle{
			filePath: filePath,
			lockChan: lockChan,
			fileLock: fl,
		}, nil
	default:
		return nil, fmt.Errorf("file %s is currently locked", filePath)
	}
}

// IsLocked checks if a file is currently locked
func (fl *FileLock) IsLocked(filePath string) bool {
	fl.mutex.RLock()
	lockChan, exists := fl.locks[filePath]
	fl.mutex.RUnlock()

	if !exists {
		return false
	}

	// Check if the channel is full (locked)
	select {
	case lockChan <- struct{}{}:
		// Channel wasn't full, remove the item we just added
		<-lockChan
		return false
	default:
		// Channel is full, file is locked
		return true
	}
}

// Cleanup removes unused lock channels to prevent memory leaks
func (fl *FileLock) Cleanup() {
	fl.mutex.Lock()
	defer fl.mutex.Unlock()

	for filePath, lockChan := range fl.locks {
		// Check if the lock is free
		select {
		case lockChan <- struct{}{}:
			// Lock is free, remove the item we just added and delete the channel
			<-lockChan
			close(lockChan)
			delete(fl.locks, filePath)
		default:
			// Lock is in use, keep it
		}
	}
}

// FileLockHandle represents an acquired file lock
type FileLockHandle struct {
	filePath string
	lockChan chan struct{}
	fileLock *FileLock
}

// Unlock releases the file lock
func (flh *FileLockHandle) Unlock() error {
	if flh.lockChan == nil {
		return fmt.Errorf("lock already released")
	}

	// Release the lock by reading from the channel
	select {
	case <-flh.lockChan:
		// Lock released successfully
		flh.lockChan = nil
		return nil
	default:
		return fmt.Errorf("lock was not held")
	}
}

// FilePath returns the file path this lock is for
func (flh *FileLockHandle) FilePath() string {
	return flh.filePath
}

// IsValid returns true if the lock handle is still valid
func (flh *FileLockHandle) IsValid() bool {
	return flh.lockChan != nil
}

// Default file lock manager instance
var defaultFileLock = NewFileLock(30 * time.Second)

// GetDefaultFileLock returns the default file lock manager
func GetDefaultFileLock() *FileLock {
	return defaultFileLock
}

// WithFileLock is a convenience function that acquires a lock, executes a function, and releases the lock
func WithFileLock(ctx context.Context, filePath string, fn func() error) error {
	lock, err := defaultFileLock.Lock(ctx, filePath)
	if err != nil {
		return fmt.Errorf("failed to acquire file lock: %w", err)
	}
	defer lock.Unlock()

	return fn()
}