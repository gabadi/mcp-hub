package platform

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewLinuxPlatformService(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	
	if service == nil {
		t.Fatal("NewLinuxPlatformService() returned nil")
	}
	
	if service.logger == nil {
		t.Error("NewLinuxPlatformService() should set default logger when nil passed")
	}
}

func TestLinuxPlatformService_GetPlatform(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	platform := service.GetPlatform()
	
	if platform != PlatformLinux {
		t.Errorf("GetPlatform() = %v, want %v", platform, PlatformLinux)
	}
}

func TestLinuxPlatformService_GetPlatformName(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	name := service.GetPlatformName()
	
	if name != "Linux" {
		t.Errorf("GetPlatformName() = %v, want %v", name, "Linux")
	}
}

func TestLinuxPlatformService_GetLogPath(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	
	// Test with XDG_DATA_HOME set
	originalXdgDataHome := os.Getenv("XDG_DATA_HOME")
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalXdgDataHome != "" {
			if err := os.Setenv("XDG_DATA_HOME", originalXdgDataHome); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("XDG_DATA_HOME"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
		if originalHome != "" {
			if err := os.Setenv("HOME", originalHome); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("HOME"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
	}()
	
	// Test with XDG_DATA_HOME set
	testXdgDataHome := "/test/xdg/data"
	if err := os.Setenv("XDG_DATA_HOME", testXdgDataHome); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	logPath := service.GetLogPath()
	
	expectedPath := filepath.Join(testXdgDataHome, "mcp-hub", "logs")
	if logPath != expectedPath {
		t.Errorf("GetLogPath() = %s, want %s", logPath, expectedPath)
	}
	
	// Test fallback to HOME/.local/share
	if err := os.Unsetenv("XDG_DATA_HOME"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	testHome := "/test/home"
	if err := os.Setenv("HOME", testHome); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	logPath = service.GetLogPath()
	
	expectedPath = filepath.Join(testHome, ".local", "share", "mcp-hub", "logs")
	if logPath != expectedPath {
		t.Errorf("GetLogPath() fallback = %s, want %s", logPath, expectedPath)
	}
	
	// Test final fallback to /tmp
	if err := os.Unsetenv("HOME"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	logPath = service.GetLogPath()
	
	expectedPath = "/tmp/mcp-hub/logs"
	if logPath != expectedPath {
		t.Errorf("GetLogPath() final fallback = %s, want %s", logPath, expectedPath)
	}
}

func TestLinuxPlatformService_GetConfigPath(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	
	// Test with XDG_CONFIG_HOME set
	originalXdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalXdgConfigHome != "" {
			if err := os.Setenv("XDG_CONFIG_HOME", originalXdgConfigHome); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("XDG_CONFIG_HOME"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
		if originalHome != "" {
			if err := os.Setenv("HOME", originalHome); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("HOME"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
	}()
	
	// Test with XDG_CONFIG_HOME set
	testXdgConfigHome := "/test/xdg/config"
	if err := os.Setenv("XDG_CONFIG_HOME", testXdgConfigHome); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	configPath := service.GetConfigPath()
	
	expectedPath := filepath.Join(testXdgConfigHome, "mcp-hub")
	if configPath != expectedPath {
		t.Errorf("GetConfigPath() = %s, want %s", configPath, expectedPath)
	}
	
	// Test fallback to HOME/.config
	if err := os.Unsetenv("XDG_CONFIG_HOME"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	testHome := "/test/home"
	if err := os.Setenv("HOME", testHome); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	configPath = service.GetConfigPath()
	
	expectedPath = filepath.Join(testHome, ".config", "mcp-hub")
	if configPath != expectedPath {
		t.Errorf("GetConfigPath() fallback = %s, want %s", configPath, expectedPath)
	}
	
	// Test final fallback to /tmp
	if err := os.Unsetenv("HOME"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	configPath = service.GetConfigPath()
	
	expectedPath = "/tmp/mcp-hub"
	if configPath != expectedPath {
		t.Errorf("GetConfigPath() final fallback = %s, want %s", configPath, expectedPath)
	}
}

func TestLinuxPlatformService_GetTempPath(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	tempPath := service.GetTempPath()
	
	if tempPath == "" {
		t.Error("GetTempPath() should not return empty string")
	}
	
	if !strings.Contains(tempPath, "mcp-hub") {
		t.Errorf("GetTempPath() should contain 'mcp-hub', got %s", tempPath)
	}
}

func TestLinuxPlatformService_GetCachePath(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	
	// Test with XDG_CACHE_HOME set
	originalXdgCacheHome := os.Getenv("XDG_CACHE_HOME")
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalXdgCacheHome != "" {
			if err := os.Setenv("XDG_CACHE_HOME", originalXdgCacheHome); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("XDG_CACHE_HOME"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
		if originalHome != "" {
			if err := os.Setenv("HOME", originalHome); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("HOME"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
	}()
	
	// Test with XDG_CACHE_HOME set
	testXdgCacheHome := "/test/xdg/cache"
	if err := os.Setenv("XDG_CACHE_HOME", testXdgCacheHome); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	cachePath := service.GetCachePath()
	
	expectedPath := filepath.Join(testXdgCacheHome, "mcp-hub")
	if cachePath != expectedPath {
		t.Errorf("GetCachePath() = %s, want %s", cachePath, expectedPath)
	}
	
	// Test fallback to HOME/.cache
	if err := os.Unsetenv("XDG_CACHE_HOME"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	testHome := "/test/home"
	if err := os.Setenv("HOME", testHome); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	cachePath = service.GetCachePath()
	
	expectedPath = filepath.Join(testHome, ".cache", "mcp-hub")
	if cachePath != expectedPath {
		t.Errorf("GetCachePath() fallback = %s, want %s", cachePath, expectedPath)
	}
	
	// Test final fallback to temp path
	if err := os.Unsetenv("HOME"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	cachePath = service.GetCachePath()
	
	expectedTempPath := service.GetTempPath()
	if cachePath != expectedTempPath {
		t.Errorf("GetCachePath() final fallback = %s, want %s", cachePath, expectedTempPath)
	}
}

func TestLinuxPlatformService_GetCommandDetectionMethod(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	method := service.GetCommandDetectionMethod()
	
	if method != "which" {
		t.Errorf("GetCommandDetectionMethod() = %s, want 'which'", method)
	}
}

func TestLinuxPlatformService_GetCommandDetectionCommand(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	cmd := service.GetCommandDetectionCommand()
	
	if cmd != "which" {
		t.Errorf("GetCommandDetectionCommand() = %s, want 'which'", cmd)
	}
}

func TestLinuxPlatformService_SupportsClipboard(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	supports := service.SupportsClipboard()
	
	// This depends on whether xclip, xsel, or wl-copy is available
	// We can't guarantee the result, but we can test that it returns a boolean
	if supports {
		t.Logf("SupportsClipboard() = true (clipboard tool is available)")
	} else {
		t.Logf("SupportsClipboard() = false (no clipboard tool available)")
	}
}

func TestLinuxPlatformService_GetClipboardMethod(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	method := service.GetClipboardMethod()
	
	// Should be one of the supported clipboard methods
	validMethods := []ClipboardMethod{
		ClipboardXclip,
		ClipboardNative,
		ClipboardUnsupported,
	}
	
	isValid := false
	for _, validMethod := range validMethods {
		if method == validMethod {
			isValid = true
			break
		}
	}
	
	if !isValid {
		t.Errorf("GetClipboardMethod() = %v, want one of %v", method, validMethods)
	}
}

func TestLinuxPlatformService_GetDefaultPermissions(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	
	filePerms := service.GetDefaultFilePermissions()
	if filePerms != 0600 {
		t.Errorf("GetDefaultFilePermissions() = %v, want 0600", filePerms)
	}
	
	dirPerms := service.GetDefaultDirectoryPermissions()
	if dirPerms != 0700 {
		t.Errorf("GetDefaultDirectoryPermissions() = %v, want 0700", dirPerms)
	}
}

func TestLinuxPlatformService_GetEnvironmentVariable(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	
	// Test with a known environment variable
	pathVar := service.GetEnvironmentVariable("PATH")
	if pathVar == "" {
		t.Error("GetEnvironmentVariable(PATH) should not be empty")
	}
	
	// Test with non-existent variable
	nonExistent := service.GetEnvironmentVariable("NON_EXISTENT_VAR_TEST")
	if nonExistent != "" {
		t.Error("GetEnvironmentVariable(NON_EXISTENT_VAR_TEST) should return empty string")
	}
}

func TestLinuxPlatformService_GetHomeDirectory(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	
	// Test with modified environment to trigger fallback
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			if err := os.Setenv("HOME", originalHome); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("HOME"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
	}()
	
	// Test fallback to HOME environment variable
	testHome := "/test/home"
	if err := os.Setenv("HOME", testHome); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	homeDir := service.GetHomeDirectory()
	
	if homeDir == "" {
		t.Error("GetHomeDirectory() should not return empty string")
	}
	
	// Test with HOME unset
	if err := os.Unsetenv("HOME"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	homeDir = service.GetHomeDirectory()
	
	// Should return empty string as fallback
	if homeDir != "" {
		// This might still return a value from os.UserHomeDir() which is fine
		t.Logf("GetHomeDirectory() returned %s even with no HOME env var", homeDir)
	}
}

func TestLinuxPlatformService_GetCurrentUser(t *testing.T) {
	service := NewLinuxPlatformService(nil)
	
	// Test with modified environment to trigger fallback
	originalUser := os.Getenv("USER")
	defer func() {
		if originalUser != "" {
			if err := os.Setenv("USER", originalUser); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("USER"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
	}()
	
	// Test fallback to USER environment variable
	testUser := "testuser"
	if err := os.Setenv("USER", testUser); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	currentUser := service.GetCurrentUser()
	
	if currentUser == "" {
		t.Error("GetCurrentUser() should not return empty string")
	}
}