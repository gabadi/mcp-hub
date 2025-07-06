package platform

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewGenericPlatformService(t *testing.T) {
	service := NewGenericPlatformService(nil)
	
	if service == nil {
		t.Fatal("NewGenericPlatformService() returned nil")
	}
	
	if service.logger == nil {
		t.Error("NewGenericPlatformService() should set default logger when nil passed")
	}
}

func TestGenericPlatformService_GetPlatform(t *testing.T) {
	service := NewGenericPlatformService(nil)
	platform := service.GetPlatform()
	
	if platform != PlatformUnknown {
		t.Errorf("GetPlatform() = %v, want %v", platform, PlatformUnknown)
	}
}

func TestGenericPlatformService_GetPlatformName(t *testing.T) {
	service := NewGenericPlatformService(nil)
	name := service.GetPlatformName()
	
	if name != "Generic" {
		t.Errorf("GetPlatformName() = %v, want %v", name, "Generic")
	}
}

func TestGenericPlatformService_GetLogPath(t *testing.T) {
	service := NewGenericPlatformService(nil)
	
	// Test with HOME set
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()
	
	// Test with HOME set
	testHome := "/test/home"
	os.Setenv("HOME", testHome)
	logPath := service.GetLogPath()
	
	expectedPath := filepath.Join(testHome, ".mcp-hub", "logs")
	if logPath != expectedPath {
		t.Errorf("GetLogPath() = %s, want %s", logPath, expectedPath)
	}
	
	// Test fallback to temp directory
	os.Unsetenv("HOME")
	logPath = service.GetLogPath()
	
	// Should contain mcp-hub/logs in temp directory
	if !strings.Contains(logPath, "mcp-hub") || !strings.Contains(logPath, "logs") {
		t.Errorf("GetLogPath() fallback should contain 'mcp-hub' and 'logs', got %s", logPath)
	}
}

func TestGenericPlatformService_GetConfigPath(t *testing.T) {
	service := NewGenericPlatformService(nil)
	
	// Test with HOME set
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()
	
	// Test with HOME set
	testHome := "/test/home"
	os.Setenv("HOME", testHome)
	configPath := service.GetConfigPath()
	
	expectedPath := filepath.Join(testHome, ".mcp-hub")
	if configPath != expectedPath {
		t.Errorf("GetConfigPath() = %s, want %s", configPath, expectedPath)
	}
	
	// Test fallback to temp directory
	os.Unsetenv("HOME")
	configPath = service.GetConfigPath()
	
	// Should contain mcp-hub in temp directory
	if !strings.Contains(configPath, "mcp-hub") {
		t.Errorf("GetConfigPath() fallback should contain 'mcp-hub', got %s", configPath)
	}
}

func TestGenericPlatformService_GetTempPath(t *testing.T) {
	service := NewGenericPlatformService(nil)
	tempPath := service.GetTempPath()
	
	if tempPath == "" {
		t.Error("GetTempPath() should not return empty string")
	}
	
	if !strings.Contains(tempPath, "mcp-hub") {
		t.Errorf("GetTempPath() should contain 'mcp-hub', got %s", tempPath)
	}
}

func TestGenericPlatformService_GetCachePath(t *testing.T) {
	service := NewGenericPlatformService(nil)
	
	// Test with HOME set
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()
	
	// Test with HOME set
	testHome := "/test/home"
	os.Setenv("HOME", testHome)
	cachePath := service.GetCachePath()
	
	expectedPath := filepath.Join(testHome, ".mcp-hub", "cache")
	if cachePath != expectedPath {
		t.Errorf("GetCachePath() = %s, want %s", cachePath, expectedPath)
	}
	
	// Test fallback to temp path
	os.Unsetenv("HOME")
	cachePath = service.GetCachePath()
	
	expectedTempPath := service.GetTempPath()
	if cachePath != expectedTempPath {
		t.Errorf("GetCachePath() fallback = %s, want %s", cachePath, expectedTempPath)
	}
}

func TestGenericPlatformService_GetCommandDetectionMethod(t *testing.T) {
	service := NewGenericPlatformService(nil)
	method := service.GetCommandDetectionMethod()
	
	if method != "which" {
		t.Errorf("GetCommandDetectionMethod() = %s, want 'which'", method)
	}
}

func TestGenericPlatformService_GetCommandDetectionCommand(t *testing.T) {
	service := NewGenericPlatformService(nil)
	cmd := service.GetCommandDetectionCommand()
	
	if cmd != "which" {
		t.Errorf("GetCommandDetectionCommand() = %s, want 'which'", cmd)
	}
}

func TestGenericPlatformService_SupportsClipboard(t *testing.T) {
	service := NewGenericPlatformService(nil)
	supports := service.SupportsClipboard()
	
	if supports {
		t.Error("SupportsClipboard() = true, want false for generic platform")
	}
}

func TestGenericPlatformService_GetClipboardMethod(t *testing.T) {
	service := NewGenericPlatformService(nil)
	method := service.GetClipboardMethod()
	
	if method != ClipboardUnsupported {
		t.Errorf("GetClipboardMethod() = %v, want %v", method, ClipboardUnsupported)
	}
}

func TestGenericPlatformService_GetDefaultPermissions(t *testing.T) {
	service := NewGenericPlatformService(nil)
	
	filePerms := service.GetDefaultFilePermissions()
	if filePerms != 0644 {
		t.Errorf("GetDefaultFilePermissions() = %v, want 0644", filePerms)
	}
	
	dirPerms := service.GetDefaultDirectoryPermissions()
	if dirPerms != 0755 {
		t.Errorf("GetDefaultDirectoryPermissions() = %v, want 0755", dirPerms)
	}
}

func TestGenericPlatformService_GetEnvironmentVariable(t *testing.T) {
	service := NewGenericPlatformService(nil)
	
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

func TestGenericPlatformService_GetHomeDirectory(t *testing.T) {
	service := NewGenericPlatformService(nil)
	
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

func TestGenericPlatformService_GetCurrentUser(t *testing.T) {
	service := NewGenericPlatformService(nil)
	
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