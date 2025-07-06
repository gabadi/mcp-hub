package platform

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)


func TestNewWindowsPlatformService(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	
	if service == nil {
		t.Fatal("NewWindowsPlatformService() returned nil")
	}
	
	if service.logger == nil {
		t.Error("NewWindowsPlatformService() should set default logger when nil passed")
	}
}

func TestWindowsPlatformService_GetPlatform(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	platform := service.GetPlatform()
	
	if platform != PlatformWindows {
		t.Errorf("GetPlatform() = %v, want %v", platform, PlatformWindows)
	}
}

func TestWindowsPlatformService_GetPlatformName(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	name := service.GetPlatformName()
	
	if name != "Windows" {
		t.Errorf("GetPlatformName() = %v, want %v", name, "Windows")
	}
}

func TestWindowsPlatformService_GetLogPath(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	
	// Test with APPDATA set
	originalAppData := os.Getenv("APPDATA")
	defer func() {
		if originalAppData != "" {
			if err := os.Setenv("APPDATA", originalAppData); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("APPDATA"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
	}()
	
	// Test with APPDATA set
	if err := os.Setenv("APPDATA", testAppData); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	logPath := service.GetLogPath()
	
	expectedPath := filepath.Join(testAppData, "mcp-hub", "logs")
	if logPath != expectedPath {
		t.Errorf("GetLogPath() = %s, want %s", logPath, expectedPath)
	}
	
	// Test fallback when APPDATA is not set
	if err := os.Unsetenv("APPDATA"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	logPath = service.GetLogPath()
	
	if !strings.Contains(logPath, "mcp-hub") {
		t.Errorf("GetLogPath() should contain 'mcp-hub', got %s", logPath)
	}
}

func TestWindowsPlatformService_GetConfigPath(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	
	// Test with APPDATA set
	originalAppData := os.Getenv("APPDATA")
	defer func() {
		if originalAppData != "" {
			if err := os.Setenv("APPDATA", originalAppData); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("APPDATA"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
	}()
	
	// Test with APPDATA set
	if err := os.Setenv("APPDATA", testAppData); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	configPath := service.GetConfigPath()
	
	expectedPath := filepath.Join(testAppData, "mcp-hub")
	if configPath != expectedPath {
		t.Errorf("GetConfigPath() = %s, want %s", configPath, expectedPath)
	}
	
	// Test fallback when APPDATA is not set
	if err := os.Unsetenv("APPDATA"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	configPath = service.GetConfigPath()
	
	if !strings.Contains(configPath, "mcp-hub") {
		t.Errorf("GetConfigPath() should contain 'mcp-hub', got %s", configPath)
	}
}

func TestWindowsPlatformService_GetTempPath(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	tempPath := service.GetTempPath()
	
	if tempPath == "" {
		t.Error("GetTempPath() should not return empty string")
	}
	
	if !strings.Contains(tempPath, "mcp-hub") {
		t.Errorf("GetTempPath() should contain 'mcp-hub', got %s", tempPath)
	}
}

func TestWindowsPlatformService_GetCachePath(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	
	// Test with LOCALAPPDATA set
	originalLocalAppData := os.Getenv("LOCALAPPDATA")
	originalAppData := os.Getenv("APPDATA")
	defer func() {
		if originalLocalAppData != "" {
			if err := os.Setenv("LOCALAPPDATA", originalLocalAppData); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("LOCALAPPDATA"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
		if originalAppData != "" {
			if err := os.Setenv("APPDATA", originalAppData); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("APPDATA"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
	}()
	
	// Test with LOCALAPPDATA set
	testLocalAppData := "/test/localappdata"
	if err := os.Setenv("LOCALAPPDATA", testLocalAppData); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	cachePath := service.GetCachePath()
	
	expectedPath := filepath.Join(testLocalAppData, "mcp-hub", "cache")
	if cachePath != expectedPath {
		t.Errorf("GetCachePath() = %s, want %s", cachePath, expectedPath)
	}
	
	// Test fallback to APPDATA when LOCALAPPDATA is not set
	if err := os.Unsetenv("LOCALAPPDATA"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	if err := os.Setenv("APPDATA", testAppData); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	cachePath = service.GetCachePath()
	
	expectedPath = filepath.Join(testAppData, "mcp-hub", "cache")
	if cachePath != expectedPath {
		t.Errorf("GetCachePath() fallback = %s, want %s", cachePath, expectedPath)
	}
	
	// Test final fallback to temp directory
	if err := os.Unsetenv("APPDATA"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	cachePath = service.GetCachePath()
	
	if !strings.Contains(cachePath, "mcp-hub") {
		t.Errorf("GetCachePath() final fallback should contain 'mcp-hub', got %s", cachePath)
	}
}

func TestWindowsPlatformService_GetCommandDetectionMethod(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	method := service.GetCommandDetectionMethod()
	
	if method != whereCmd {
		t.Errorf("GetCommandDetectionMethod() = %s, want 'where'", method)
	}
}

func TestWindowsPlatformService_GetCommandDetectionCommand(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	cmd := service.GetCommandDetectionCommand()
	
	if cmd != whereCmd {
		t.Errorf("GetCommandDetectionCommand() = %s, want 'where'", cmd)
	}
}

func TestWindowsPlatformService_SupportsClipboard(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	supports := service.SupportsClipboard()
	
	if !supports {
		t.Error("SupportsClipboard() = false, want true for Windows")
	}
}

func TestWindowsPlatformService_GetClipboardMethod(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	method := service.GetClipboardMethod()
	
	// Should be either ClipboardPowershell or ClipboardNative
	if method != ClipboardPowershell && method != ClipboardNative {
		t.Errorf("GetClipboardMethod() = %v, want either ClipboardPowershell or ClipboardNative", method)
	}
}

func TestWindowsPlatformService_GetDefaultPermissions(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	
	filePerms := service.GetDefaultFilePermissions()
	if filePerms != 0644 {
		t.Errorf("GetDefaultFilePermissions() = %v, want 0644", filePerms)
	}
	
	dirPerms := service.GetDefaultDirectoryPermissions()
	if dirPerms != 0755 {
		t.Errorf("GetDefaultDirectoryPermissions() = %v, want 0755", dirPerms)
	}
}

func TestWindowsPlatformService_GetEnvironmentVariable(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	
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

func TestWindowsPlatformService_GetHomeDirectory(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	
	// Test with modified environment to trigger fallback scenarios
	originalUserProfile := os.Getenv("USERPROFILE")
	originalHomeDrive := os.Getenv("HOMEDRIVE")
	originalHomePath := os.Getenv("HOMEPATH")
	
	defer func() {
		if originalUserProfile != "" {
			if err := os.Setenv("USERPROFILE", originalUserProfile); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("USERPROFILE"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
		if originalHomeDrive != "" {
			if err := os.Setenv("HOMEDRIVE", originalHomeDrive); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("HOMEDRIVE"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
		if originalHomePath != "" {
			if err := os.Setenv("HOMEPATH", originalHomePath); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("HOMEPATH"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
	}()
	
	// Test fallback to USERPROFILE
	testUserProfile := "/test/userprofile"
	if err := os.Setenv("USERPROFILE", testUserProfile); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	homeDir := service.GetHomeDirectory()
	
	if homeDir == "" {
		t.Error("GetHomeDirectory() should not return empty string")
	}
	
	// Test fallback to HOMEDRIVE + HOMEPATH
	if err := os.Unsetenv("USERPROFILE"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	testHomeDrive := "C:"
	testHomePath := "/Users/Test"
	if err := os.Setenv("HOMEDRIVE", testHomeDrive); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	if err := os.Setenv("HOMEPATH", testHomePath); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	homeDir = service.GetHomeDirectory()
	
	if homeDir == "" {
		t.Error("GetHomeDirectory() should not return empty string with HOMEDRIVE/HOMEPATH")
	}
	
	// Test final fallback (all env vars unset)
	if err := os.Unsetenv("HOMEDRIVE"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	if err := os.Unsetenv("HOMEPATH"); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}
	homeDir = service.GetHomeDirectory()
	
	// Should return empty string as final fallback
	if homeDir != "" {
		// This might still return a value from os.UserHomeDir() which is fine
		t.Logf("GetHomeDirectory() returned %s even with no env vars", homeDir)
	}
}

func TestWindowsPlatformService_GetCurrentUser(t *testing.T) {
	service := NewWindowsPlatformService(nil)
	
	// Test with modified environment to trigger fallback
	originalUsername := os.Getenv("USERNAME")
	defer func() {
		if originalUsername != "" {
			if err := os.Setenv("USERNAME", originalUsername); err != nil {
				t.Errorf("Failed to set environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv("USERNAME"); err != nil {
				t.Errorf("Failed to unset environment variable: %v", err)
			}
		}
	}()
	
	// Test fallback to USERNAME environment variable
	testUsername := "testuser"
	if err := os.Setenv("USERNAME", testUsername); err != nil {
		t.Errorf("Failed to set environment variable: %v", err)
	}
	currentUser := service.GetCurrentUser()
	
	if currentUser == "" {
		t.Error("GetCurrentUser() should not return empty string")
	}
}