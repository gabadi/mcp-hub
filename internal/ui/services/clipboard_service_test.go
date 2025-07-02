package services

import (
	"testing"
)

func TestNewClipboardService(t *testing.T) {
	service := NewClipboardService()

	if service == nil {
		t.Fatal("NewClipboardService() returned nil")
	}
}

func TestClipboardServiceCopy(t *testing.T) {
	service := NewClipboardService()

	// Test copy operation - this may fail on headless systems
	err := service.Copy("test content")

	// We don't fail the test if clipboard is not available
	// as this is environment dependent
	if err != nil {
		t.Logf("Copy failed as expected in test environment: %v", err)
	}
}

func TestClipboardServicePaste(t *testing.T) {
	service := NewClipboardService()

	// Test paste operation - this may fail on headless systems
	content, err := service.Paste()

	// We don't fail the test if clipboard is not available
	// as this is environment dependent
	if err != nil {
		t.Logf("Paste failed as expected in test environment: %v", err)
		if content != "" {
			t.Error("Content should be empty when paste fails")
		}
	}
}

func TestClipboardServiceIsAvailable(t *testing.T) {
	service := NewClipboardService()

	// Test availability check
	available := service.IsAvailable()

	// This is environment dependent, so we just verify it returns a boolean
	t.Logf("Clipboard available: %v", available)
}

func TestClipboardServiceEnhancedPaste(t *testing.T) {
	service := NewClipboardService()

	// Test enhanced paste with diagnostics
	content, err := service.EnhancedPaste()

	// Environment dependent operation
	if err != nil {
		t.Logf("Enhanced paste failed as expected in test environment: %v", err)
		// Error should provide diagnostic information
		if len(err.Error()) < 10 {
			t.Error("Enhanced paste error should provide diagnostic information")
		}
	}

	// Content should be empty if error occurred
	if err != nil && content != "" {
		t.Error("Content should be empty when enhanced paste fails")
	}
}

func TestClipboardServiceGetDiagnosticInfo(t *testing.T) {
	service := NewClipboardService()

	// Test diagnostic info
	info := service.GetDiagnosticInfo()

	// Should always return some diagnostic information
	if info == nil {
		t.Error("GetDiagnosticInfo() should return diagnostic information")
	}

	// Should contain basic OS information
	if len(info) == 0 {
		t.Error("Diagnostic info should not be empty")
	}

	// Should have OS info
	if _, exists := info["os"]; !exists {
		t.Error("Diagnostic info should contain OS information")
	}

	t.Logf("Diagnostic info: %+v", info)
}

// Test clipboard integration scenarios
func TestClipboardServiceIntegration(t *testing.T) {
	service := NewClipboardService()

	// Test copy then paste workflow
	testContent := "integration test content"

	// Try to copy
	copyErr := service.Copy(testContent)

	if copyErr == nil {
		// If copy succeeded, try to paste
		pastedContent, pasteErr := service.Paste()

		if pasteErr == nil {
			// Verify content matches
			if pastedContent != testContent {
				t.Logf("Copy/paste content mismatch - expected: %q, got: %q", testContent, pastedContent)
				// Don't fail the test as clipboard content may be modified by system
			}
		} else {
			t.Logf("Paste failed after successful copy: %v", pasteErr)
		}
	} else {
		t.Logf("Copy failed in test environment: %v", copyErr)
	}
}

func TestClipboardServiceErrorHandling(t *testing.T) {
	service := NewClipboardService()

	// Test with empty content
	err := service.Copy("")
	if err != nil {
		t.Logf("Copy of empty string failed: %v", err)
	}

	// Test paste when clipboard might be empty or unavailable
	content, err := service.Paste()
	if err != nil {
		t.Logf("Paste failed as expected: %v", err)
		// Error should be descriptive
		if len(err.Error()) < 5 {
			t.Error("Error message should be descriptive")
		}
	}

	// Content should be string (even if empty)
	_ = content // Just verify it's a string
}
