package storage

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestFileLock_BasicLocking(t *testing.T) {
	fl := NewFileLock(1 * time.Second)
	filePath := "/test/file.json"
	
	ctx := context.Background()
	
	// Acquire lock
	lock1, err := fl.Lock(ctx, filePath)
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}
	
	if !lock1.IsValid() {
		t.Error("Lock should be valid after acquisition")
	}
	
	if lock1.FilePath() != filePath {
		t.Errorf("Expected file path %s, got %s", filePath, lock1.FilePath())
	}
	
	// Release lock
	err = lock1.Unlock()
	if err != nil {
		t.Errorf("Failed to release lock: %v", err)
	}
	
	if lock1.IsValid() {
		t.Error("Lock should be invalid after release")
	}
}

func TestFileLock_TryLock(t *testing.T) {
	fl := NewFileLock(1 * time.Second)
	filePath := "/test/file.json"
	
	// Try lock on free file
	lock1, err := fl.TryLock(filePath)
	if err != nil {
		t.Fatalf("Failed to try lock on free file: %v", err)
	}
	defer lock1.Unlock()
	
	// Try lock on locked file
	lock2, err := fl.TryLock(filePath)
	if err == nil {
		lock2.Unlock()
		t.Error("Should not be able to try lock on already locked file")
	}
}

func TestFileLock_IsLocked(t *testing.T) {
	fl := NewFileLock(1 * time.Second)
	filePath := "/test/file.json"
	
	// File should not be locked initially
	if fl.IsLocked(filePath) {
		t.Error("File should not be locked initially")
	}
	
	// Acquire lock
	lock, err := fl.TryLock(filePath)
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}
	
	// File should be locked now
	if !fl.IsLocked(filePath) {
		t.Error("File should be locked after acquisition")
	}
	
	// Release lock
	lock.Unlock()
	
	// File should not be locked anymore
	if fl.IsLocked(filePath) {
		t.Error("File should not be locked after release")
	}
}

func TestFileLock_ConcurrentAccess(t *testing.T) {
	fl := NewFileLock(2 * time.Second)
	filePath := "/test/concurrent.json"
	
	var wg sync.WaitGroup
	var lock1Acquired, lock2Acquired bool
	var lock1Time, lock2Time time.Time
	
	// First goroutine acquires lock, holds it for 500ms
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := context.Background()
		
		lock, err := fl.Lock(ctx, filePath)
		if err != nil {
			t.Errorf("Failed to acquire lock in goroutine 1: %v", err)
			return
		}
		defer lock.Unlock()
		
		lock1Acquired = true
		lock1Time = time.Now()
		
		// Hold lock for 500ms
		time.Sleep(500 * time.Millisecond)
	}()
	
	// Second goroutine tries to acquire lock after 100ms
	wg.Add(1)
	go func() {
		defer wg.Done()
		
		// Wait a bit to ensure first goroutine gets the lock first
		time.Sleep(100 * time.Millisecond)
		
		ctx := context.Background()
		lock, err := fl.Lock(ctx, filePath)
		if err != nil {
			t.Errorf("Failed to acquire lock in goroutine 2: %v", err)
			return
		}
		defer lock.Unlock()
		
		lock2Acquired = true
		lock2Time = time.Now()
	}()
	
	wg.Wait()
	
	if !lock1Acquired {
		t.Error("Lock 1 should have been acquired")
	}
	
	if !lock2Acquired {
		t.Error("Lock 2 should have been acquired")
	}
	
	// Lock 2 should have been acquired after lock 1 was released
	if lock2Time.Before(lock1Time.Add(400 * time.Millisecond)) {
		t.Error("Lock 2 should have been acquired after lock 1 was held for at least 400ms")
	}
}

func TestFileLock_Timeout(t *testing.T) {
	fl := NewFileLock(200 * time.Millisecond)
	filePath := "/test/timeout.json"
	
	ctx := context.Background()
	
	// Acquire lock
	lock1, err := fl.Lock(ctx, filePath)
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}
	defer lock1.Unlock()
	
	// Try to acquire lock again with timeout
	start := time.Now()
	lock2, err := fl.Lock(ctx, filePath)
	duration := time.Since(start)
	
	if err == nil {
		lock2.Unlock()
		t.Error("Should have timed out waiting for lock")
	}
	
	// Should have waited approximately the timeout duration
	expectedDuration := 200 * time.Millisecond
	if duration < expectedDuration || duration > expectedDuration+100*time.Millisecond {
		t.Errorf("Expected timeout around %v, got %v", expectedDuration, duration)
	}
}

func TestFileLock_MultipleFiles(t *testing.T) {
	fl := NewFileLock(1 * time.Second)
	
	file1 := "/test/file1.json"
	file2 := "/test/file2.json"
	
	// Should be able to lock different files simultaneously
	lock1, err := fl.TryLock(file1)
	if err != nil {
		t.Fatalf("Failed to lock file1: %v", err)
	}
	defer lock1.Unlock()
	
	lock2, err := fl.TryLock(file2)
	if err != nil {
		t.Fatalf("Failed to lock file2: %v", err)
	}
	defer lock2.Unlock()
	
	// Both locks should be valid
	if !lock1.IsValid() || !lock2.IsValid() {
		t.Error("Both locks should be valid")
	}
}

func TestFileLock_Cleanup(t *testing.T) {
	fl := NewFileLock(1 * time.Second)
	filePath := "/test/cleanup.json"
	
	// Acquire and release lock
	lock, err := fl.TryLock(filePath)
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}
	lock.Unlock()
	
	// Check that we have a lock channel for this file
	fl.mutex.RLock()
	_, exists := fl.locks[filePath]
	fl.mutex.RUnlock()
	
	if !exists {
		t.Error("Lock channel should exist before cleanup")
	}
	
	// Run cleanup
	fl.Cleanup()
	
	// Check that the lock channel was removed
	fl.mutex.RLock()
	_, exists = fl.locks[filePath]
	fl.mutex.RUnlock()
	
	if exists {
		t.Error("Lock channel should be removed after cleanup")
	}
}

func TestFileLockHandle_DoubleUnlock(t *testing.T) {
	fl := NewFileLock(1 * time.Second)
	filePath := "/test/double_unlock.json"
	
	lock, err := fl.TryLock(filePath)
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}
	
	// First unlock should succeed
	err = lock.Unlock()
	if err != nil {
		t.Errorf("First unlock failed: %v", err)
	}
	
	// Second unlock should fail
	err = lock.Unlock()
	if err == nil {
		t.Error("Second unlock should have failed")
	}
}

func TestWithFileLockOperation(t *testing.T) {
	filePath := "/test/with_lock.json"
	ctx := context.Background()
	
	executed := false
	
	result, err := WithFileLockOperation(ctx, filePath, func() (string, error) {
		executed = true
		return "success", nil
	})
	
	if err != nil {
		t.Errorf("WithFileLockOperation failed: %v", err)
	}
	
	if !executed {
		t.Error("Operation should have been executed")
	}
	
	if result != "success" {
		t.Errorf("Expected result 'success', got '%s'", result)
	}
}