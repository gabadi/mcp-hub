package services

import (
	"sync"
	"time"

	"github.com/atotto/clipboard"
)

// ClipboardService provides clipboard operations
type ClipboardService struct {
	isAvailable   *bool
	lastCheck     time.Time
	checkDuration time.Duration
	mutex         sync.RWMutex
}

// NewClipboardService creates a new clipboard service
func NewClipboardService() *ClipboardService {
	return &ClipboardService{
		checkDuration: 30 * time.Second, // Cache availability for 30 seconds
	}
}

// Copy copies text to the system clipboard
func (cs *ClipboardService) Copy(text string) error {
	return clipboard.WriteAll(text)
}

// Paste gets text from the system clipboard
func (cs *ClipboardService) Paste() (string, error) {
	return clipboard.ReadAll()
}

// IsAvailable checks if clipboard operations are available with caching
func (cs *ClipboardService) IsAvailable() bool {
	cs.mutex.RLock()

	// Check if we have a cached result that's still valid
	if cs.isAvailable != nil && time.Since(cs.lastCheck) < cs.checkDuration {
		result := *cs.isAvailable
		cs.mutex.RUnlock()
		return result
	}

	cs.mutex.RUnlock()

	// Need to check availability
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	// Double-check pattern in case another goroutine updated while we were waiting
	if cs.isAvailable != nil && time.Since(cs.lastCheck) < cs.checkDuration {
		return *cs.isAvailable
	}

	// Perform the actual availability check
	_, err := clipboard.ReadAll()
	available := err == nil

	cs.isAvailable = &available
	cs.lastCheck = time.Now()

	return available
}
