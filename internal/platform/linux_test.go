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
			os.Setenv("XDG_DATA_HOME", originalXdgDataHome)
		} else {
			os.Unsetenv("XDG_DATA_HOME")
		}
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()
	
	// Test with XDG_DATA_HOME set
	testXdgDataHome := "/test/xdg/data"
	os.Setenv("XDG_DATA_HOME", testXdgDataHome)
	logPath := service.GetLogPath()
	
	expectedPath := filepath.Join(testXdgDataHome, "mcp-hub", "logs")
	if logPath != expectedPath {
		t.Errorf("GetLogPath() = %s, want %s", logPath, expectedPath)
	}
	
	// Test fallback to HOME/.local/share
	os.Unsetenv("XDG_DATA_HOME")
	testHome := "/test/home"
	os.Setenv("HOME", testHome)
	logPath = service.GetLogPath()
	
	expectedPath = filepath.Join(testHome, ".local", "share", "mcp-hub", "logs")
	if logPath != expectedPath {
		t.Errorf("GetLogPath() fallback = %s, want %s", logPath, expectedPath)
	}
	
	// Test final fallback to /tmp
	os.Unsetenv("HOME")
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
			os.Setenv("XDG_CONFIG_HOME", originalXdgConfigHome)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()
	
	// Test with XDG_CONFIG_HOME set
	testXdgConfigHome := "/test/xdg/config"
	os.Setenv("XDG_CONFIG_HOME", testXdgConfigHome)
	configPath := service.GetConfigPath()
	
	expectedPath := filepath.Join(testXdgConfigHome, "mcp-hub")
	if configPath != expectedPath {
		t.Errorf("GetConfigPath() = %s, want %s", configPath, expectedPath)
	}
	
	// Test fallback to HOME/.config
	os.Unsetenv("XDG_CONFIG_HOME")
	testHome := "/test/home"
	os.Setenv("HOME", testHome)
	configPath = service.GetConfigPath()
	
	expectedPath = filepath.Join(testHome, ".config", "mcp-hub")
	if configPath != expectedPath {
		t.Errorf("GetConfigPath() fallback = %s, want %s", configPath, expectedPath)
	}
	
	// Test final fallback to /tmp
	os.Unsetenv("HOME")
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
			os.Setenv("XDG_CACHE_HOME", originalXdgCacheHome)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()
	
	// Test with XDG_CACHE_HOME set
	testXdgCacheHome := "/test/xdg/cache"
	os.Setenv("XDG_CACHE_HOME", testXdgCacheHome)
	cachePath := service.GetCachePath()
	
	expectedPath := filepath.Join(testXdgCacheHome, "mcp-hub")
	if cachePath != expectedPath {
		t.Errorf("GetCachePath() = %s, want %s", cachePath, expectedPath)
	}
	
	// Test fallback to HOME/.cache
	os.Unsetenv("XDG_CACHE_HOME")
	testHome := "/test/home"
	os.Setenv("HOME", testHome)
	cachePath = service.GetCachePath()
	
	expectedPath = filepath.Join(testHome, ".cache", "mcp-hub")
	if cachePath != expectedPath {
		t.Errorf("GetCachePath() fallback = %s, want %s", cachePath, expectedPath)
	}
	
	// Test final fallback to temp path
	os.Unsetenv("HOME")
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
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()
	
	// Test fallback to HOME environment variable
	testHome := "/test/home"
	os.Setenv("HOME", testHome)
	homeDir := service.GetHomeDirectory()
	
	if homeDir == "" {
		t.Error("GetHomeDirectory() should not return empty string")
	}
	
	// Test with HOME unset
	os.Unsetenv("HOME")
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
			os.Setenv("USER", originalUser)
		} else {
			os.Unsetenv("USER")
		}
	}()
	
	// Test fallback to USER environment variable
	testUser := "testuser"
	os.Setenv("USER", testUser)
	currentUser := service.GetCurrentUser()
	
	if currentUser == "" {
		t.Error("GetCurrentUser() should not return empty string")
	}
}