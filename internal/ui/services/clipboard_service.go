package services

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	atottoclipboard "github.com/atotto/clipboard"
	designclipboard "golang.design/x/clipboard"
)

// ClipboardService provides clipboard operations with enhanced macOS support
type ClipboardService struct {
	isAvailable        *bool
	lastCheck          time.Time
	checkDuration      time.Duration
	mutex              sync.RWMutex
	useDesignClipboard bool
	initialized        bool
}

// NewClipboardService creates a new clipboard service with enhanced macOS support
func NewClipboardService() *ClipboardService {
	cs := &ClipboardService{
		checkDuration: 30 * time.Second, // Cache availability for 30 seconds
	}

	// Try to initialize golang.design/x/clipboard for better macOS support
	if runtime.GOOS == PlatformDarwin {
		if err := designclipboard.Init(); err == nil {
			cs.useDesignClipboard = true
			cs.initialized = true
		}
	}

	return cs
}

// Copy copies text to the system clipboard with macOS fallbacks
func (cs *ClipboardService) Copy(text string) error {
	// Try golang.design/x/clipboard first on macOS
	if cs.useDesignClipboard && runtime.GOOS == PlatformDarwin {
		designclipboard.Write(designclipboard.FmtText, []byte(text))
		return nil
	}

	// Try atotto/clipboard
	err := atottoclipboard.WriteAll(text)
	if err == nil {
		return nil
	}

	// macOS fallback using pbcopy
	if runtime.GOOS == PlatformDarwin {
		return cs.copyWithPbcopy(text)
	}

	return err
}

// Paste gets text from the system clipboard with macOS fallbacks
func (cs *ClipboardService) Paste() (string, error) {
	// Try golang.design/x/clipboard first on macOS
	if cs.useDesignClipboard && runtime.GOOS == PlatformDarwin {
		data := designclipboard.Read(designclipboard.FmtText)
		if len(data) > 0 {
			return string(data), nil
		}
	}

	// Try atotto/clipboard
	content, err := atottoclipboard.ReadAll()
	if err == nil && content != "" {
		return content, nil
	}

	// macOS fallback using pbpaste
	if runtime.GOOS == PlatformDarwin {
		return cs.pasteWithPbpaste()
	}

	return content, err
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

	// Perform the actual availability check using multiple methods
	available := false

	// Test golang.design/x/clipboard if available
	if cs.useDesignClipboard && runtime.GOOS == PlatformDarwin {
		// Try reading to ensure clipboard is functional, but even empty clipboard is considered available
		_ = designclipboard.Read(designclipboard.FmtText)
		available = true
	}

	// Test atotto/clipboard
	if !available {
		_, err := atottoclipboard.ReadAll()
		available = err == nil
	}

	// Test macOS pbpaste fallback
	if !available && runtime.GOOS == PlatformDarwin {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		cmd2 := exec.CommandContext(ctx, "pbpaste")
		_, err := cmd2.Output()
		available = err == nil
	}

	cs.isAvailable = &available
	cs.lastCheck = time.Now()

	return available
}

// copyWithPbcopy uses macOS pbcopy command as fallback
func (cs *ClipboardService) copyWithPbcopy(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// pasteWithPbpaste uses macOS pbpaste command as fallback
func (cs *ClipboardService) pasteWithPbpaste() (string, error) {
	cmd := exec.Command("pbpaste")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(output), "\n\r"), nil
}

// GetDiagnosticInfo returns diagnostic information about clipboard service
func (cs *ClipboardService) GetDiagnosticInfo() map[string]interface{} {
	info := make(map[string]interface{})
	info["os"] = runtime.GOOS
	info["useDesignClipboard"] = cs.useDesignClipboard
	info["initialized"] = cs.initialized
	info["isAvailable"] = cs.IsAvailable()

	// Test clipboard access methods
	if runtime.GOOS == PlatformDarwin {
		// Test pbpaste availability
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		cmd2 := exec.CommandContext(ctx, "pbpaste")
		_, err := cmd2.Output()
		info["pbpasteAvailable"] = err == nil

		// Test pbcopy availability
		cmd3 := exec.CommandContext(ctx, "pbcopy")
		cmd3.Stdin = strings.NewReader("test")
		err2 := cmd3.Run()
		info["pbcopyAvailable"] = err2 == nil
	}

	return info
}

// EnhancedPaste provides detailed error information for debugging
func (cs *ClipboardService) EnhancedPaste() (string, error) {
	var lastErr error

	// Method 1: Try golang.design/x/clipboard on macOS
	var designErr error
	if cs.useDesignClipboard && runtime.GOOS == PlatformDarwin {
		data := designclipboard.Read(designclipboard.FmtText)
		if len(data) > 0 {
			return string(data), nil
		}
		// Store error for fallback reporting
		designErr = fmt.Errorf("golang.design/x/clipboard returned empty data")
	}

	// Method 2: Try atotto/clipboard
	content, err := atottoclipboard.ReadAll()
	if err == nil && content != "" {
		return content, nil
	}
	if err != nil {
		lastErr = fmt.Errorf("atotto/clipboard error: %w", err)
	} else {
		lastErr = fmt.Errorf("atotto/clipboard returned empty content")
	}

	// Include design error if it occurred
	if designErr != nil {
		lastErr = fmt.Errorf("%v, design/clipboard: %v", lastErr, designErr)
	}

	// Method 3: macOS fallback using pbpaste
	if runtime.GOOS == PlatformDarwin {
		content, err := cs.pasteWithPbpaste()
		if err == nil && content != "" {
			return content, nil
		}
		if err != nil {
			return "", fmt.Errorf("all clipboard methods failed - pbpaste error: %w, previous: %v", err, lastErr)
		}
		return "", fmt.Errorf("all clipboard methods failed - pbpaste returned empty, previous: %v", lastErr)
	}

	return "", fmt.Errorf("clipboard access failed: %v", lastErr)
}
