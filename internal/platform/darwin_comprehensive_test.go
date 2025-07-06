package platform

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDarwinPlatformService_GetLogPath_ErrorScenarios(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	
	// Test with modified environment to trigger fallback
	originalHome := os.Getenv("HOME")
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Errorf("Failed to set environment variable: %v", err)
		}
	}()
	
	// Test fallback to /tmp when home directory is not available
	if err := os.Unsetenv("HOME"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	logPath := service.GetLogPath()
	
	if logPath != "/tmp/mcp-hub" {
		t.Errorf("Expected fallback to /tmp/mcp-hub, got %s", logPath)
	}
}

func TestDarwinPlatformService_GetConfigPath_ErrorScenarios(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping Darwin-specific test on non-Darwin platform")
	}
	
	service := NewDarwinPlatformService(nil)
	
	// Test with modified environment to trigger fallback
	originalHome := os.Getenv("HOME")
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Errorf("Failed to set environment variable: %v", err)
		}
	}()
	
	// Test fallback behavior when UserConfigDir fails by setting HOME to empty
	if err := os.Unsetenv("HOME"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	configPath := service.GetConfigPath()
	
	// Should fallback to /tmp/mcp-hub when home directory is not available
	if configPath != "/tmp/mcp-hub" {
		t.Errorf("Expected fallback to /tmp/mcp-hub, got %s", configPath)
	}
}

func TestDarwinPlatformService_GetCachePath_ErrorScenarios(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	
	// Test with modified environment to trigger fallback
	originalHome := os.Getenv("HOME")
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Errorf("Failed to set environment variable: %v", err)
		}
	}()
	
	// Test fallback to temp path when home directory is not available
	if err := os.Unsetenv("HOME"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	cachePath := service.GetCachePath()
	
	// Should fallback to temp path
	expectedTempPath := service.GetTempPath()
	if cachePath != expectedTempPath {
		t.Errorf("Expected fallback to %s, got %s", expectedTempPath, cachePath)
	}
}

func TestDarwinPlatformService_GetClipboardMethod_ErrorScenarios(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	
	// Test that it returns ClipboardPbcopy or ClipboardNative
	method := service.GetClipboardMethod()
	
	if method != ClipboardPbcopy && method != ClipboardNative {
		t.Errorf("Expected ClipboardPbcopy or ClipboardNative, got %v", method)
	}
	
	// We can't easily test the fallback scenario without mocking exec.LookPath
	// but we can verify the method is consistent
	method2 := service.GetClipboardMethod()
	if method != method2 {
		t.Error("GetClipboardMethod() should return consistent results")
	}
}

func TestDarwinPlatformService_GetCommandDetectionMethod(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	
	method := service.GetCommandDetectionMethod()
	if method != "which" {
		t.Errorf("Expected 'which', got %s", method)
	}
}

func TestDarwinPlatformService_GetEnvironmentVariable(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	
	// Test with known environment variables
	homeVar := service.GetEnvironmentVariable("HOME")
	if homeVar == "" {
		t.Error("GetEnvironmentVariable(HOME) should not be empty")
	}
	
	userVar := service.GetEnvironmentVariable("USER")
	if userVar == "" {
		t.Error("GetEnvironmentVariable(USER) should not be empty")
	}
	
	// Test with non-existent variable
	nonExistent := service.GetEnvironmentVariable("NON_EXISTENT_VAR_TEST")
	if nonExistent != "" {
		t.Error("GetEnvironmentVariable(NON_EXISTENT_VAR_TEST) should return empty string")
	}
}

func TestDarwinPlatformService_GetHomeDirectory_ErrorScenarios(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	
	// Test with modified environment to trigger fallback
	originalHome := os.Getenv("HOME")
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Errorf("Failed to set environment variable: %v", err)
		}
	}()
	
	// Test fallback to HOME environment variable
	expectedHome := "/test/home"
	if err := os.Setenv("HOME", expectedHome); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	homeDir := service.GetHomeDirectory()
	
	if !strings.Contains(homeDir, expectedHome) && homeDir != expectedHome {
		t.Errorf("Expected home directory to be %s or contain it, got %s", expectedHome, homeDir)
	}
}

func TestDarwinPlatformService_GetCurrentUser_ErrorScenarios(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	
	// Test with modified environment to trigger fallback
	originalUser := os.Getenv("USER")
	defer func() {
		if err := os.Setenv("USER", originalUser); err != nil {
			t.Errorf("Failed to set environment variable: %v", err)
		}
	}()
	
	// Test fallback to USER environment variable
	expectedUser := "testuser"
	if err := os.Setenv("USER", expectedUser); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	currentUser := service.GetCurrentUser()
	
	// Should return either the actual user or fallback to environment variable
	if currentUser == "" {
		t.Error("GetCurrentUser() should not return empty string")
	}
}

func TestDarwinPlatformService_GetTempPath(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	
	tempPath := service.GetTempPath()
	if tempPath == "" {
		t.Error("GetTempPath() should not return empty string")
	}
	
	if !strings.Contains(tempPath, "mcp-hub") {
		t.Errorf("GetTempPath() should contain 'mcp-hub', got %s", tempPath)
	}
	
	if !filepath.IsAbs(tempPath) {
		t.Errorf("GetTempPath() should return absolute path, got %s", tempPath)
	}
}

func TestDarwinPlatformService_AllMethods(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	
	// Test all methods return non-empty/valid values
	if service.GetPlatform() != PlatformDarwin {
		t.Error("GetPlatform() should return PlatformDarwin")
	}
	
	if service.GetPlatformName() != "macOS" {
		t.Error("GetPlatformName() should return 'macOS'")
	}
	
	if service.GetLogPath() == "" {
		t.Error("GetLogPath() should not be empty")
	}
	
	if service.GetConfigPath() == "" {
		t.Error("GetConfigPath() should not be empty")
	}
	
	if service.GetTempPath() == "" {
		t.Error("GetTempPath() should not be empty")
	}
	
	if service.GetCachePath() == "" {
		t.Error("GetCachePath() should not be empty")
	}
	
	if service.GetCommandDetectionMethod() == "" {
		t.Error("GetCommandDetectionMethod() should not be empty")
	}
	
	if service.GetCommandDetectionCommand() == "" {
		t.Error("GetCommandDetectionCommand() should not be empty")
	}
	
	if !service.SupportsClipboard() {
		t.Error("SupportsClipboard() should return true for macOS")
	}
	
	if service.GetClipboardMethod() == ClipboardUnsupported {
		t.Error("GetClipboardMethod() should not return ClipboardUnsupported for macOS")
	}
	
	if service.GetDefaultFilePermissions() == 0 {
		t.Error("GetDefaultFilePermissions() should not be zero")
	}
	
	if service.GetDefaultDirectoryPermissions() == 0 {
		t.Error("GetDefaultDirectoryPermissions() should not be zero")
	}
	
	if service.GetHomeDirectory() == "" {
		t.Error("GetHomeDirectory() should not be empty")
	}
	
	if service.GetCurrentUser() == "" {
		t.Error("GetCurrentUser() should not be empty")
	}
}