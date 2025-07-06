package platform

import (
	"os"
	"path/filepath"
	"testing"
)


func TestNewMockPlatformService(t *testing.T) {
	mock := NewMockPlatformService()
	
	if mock == nil {
		t.Fatal("NewMockPlatformService() returned nil")
	}
	
	// Test default values
	if mock.GetPlatform() != PlatformDarwin {
		t.Errorf("Expected default platform to be PlatformDarwin, got %v", mock.GetPlatform())
	}
	
	if mock.GetPlatformName() != darwinOS {
		t.Errorf("Expected platform name to be 'darwin', got %s", mock.GetPlatformName())
	}
	
	if !mock.SupportsClipboard() {
		t.Error("Expected clipboard support to be true by default")
	}
	
	if mock.GetClipboardMethod() != ClipboardPbcopy {
		t.Errorf("Expected clipboard method to be ClipboardPbcopy, got %v", mock.GetClipboardMethod())
	}
}

func TestNewMockPlatformServiceForOS(t *testing.T) {
	testCases := []struct {
		osName              string
		expectedPlatform    PlatformType
		expectedName        string
		expectedClipboard   ClipboardMethod
		expectedDetection   string
		expectedSupportsClip bool
	}{
		{
			osName:              "darwin",
			expectedPlatform:    PlatformDarwin,
			expectedName:        "darwin",
			expectedClipboard:   ClipboardPbcopy,
			expectedDetection:   "which",
			expectedSupportsClip: true,
		},
		{
			osName:              "windows",
			expectedPlatform:    PlatformWindows,
			expectedName:        "windows",
			expectedClipboard:   ClipboardPowershell,
			expectedDetection:   "where",
			expectedSupportsClip: true,
		},
		{
			osName:              "linux",
			expectedPlatform:    PlatformLinux,
			expectedName:        "linux",
			expectedClipboard:   ClipboardXclip,
			expectedDetection:   "which",
			expectedSupportsClip: true,
		},
		{
			osName:              "unknown",
			expectedPlatform:    PlatformUnknown,
			expectedName:        "unknown",
			expectedClipboard:   ClipboardUnsupported,
			expectedDetection:   "which",
			expectedSupportsClip: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.osName, func(t *testing.T) {
			mock := NewMockPlatformServiceForOS(tc.osName)
			
			if mock == nil {
				t.Fatal("NewMockPlatformServiceForOS() returned nil")
			}
			
			if mock.GetPlatform() != tc.expectedPlatform {
				t.Errorf("Expected platform %v, got %v", tc.expectedPlatform, mock.GetPlatform())
			}
			
			if mock.GetPlatformName() != tc.expectedName {
				t.Errorf("Expected platform name %s, got %s", tc.expectedName, mock.GetPlatformName())
			}
			
			if mock.GetClipboardMethod() != tc.expectedClipboard {
				t.Errorf("Expected clipboard method %v, got %v", tc.expectedClipboard, mock.GetClipboardMethod())
			}
			
			if mock.GetCommandDetectionCommand() != tc.expectedDetection {
				t.Errorf("Expected detection command %s, got %s", tc.expectedDetection, mock.GetCommandDetectionCommand())
			}
			
			if mock.SupportsClipboard() != tc.expectedSupportsClip {
				t.Errorf("Expected clipboard support %v, got %v", tc.expectedSupportsClip, mock.SupportsClipboard())
			}
		})
	}
}

func TestMockPlatformService_GetPaths(t *testing.T) {
	mock := NewMockPlatformService()
	
	// Test that paths are not empty
	if mock.GetLogPath() == "" {
		t.Error("GetLogPath() should not return empty string")
	}
	
	if mock.GetConfigPath() == "" {
		t.Error("GetConfigPath() should not return empty string")
	}
	
	if mock.GetTempPath() == "" {
		t.Error("GetTempPath() should not return empty string")
	}
	
	if mock.GetCachePath() == "" {
		t.Error("GetCachePath() should not return empty string")
	}
	
	// Test that paths contain expected components
	logPath := mock.GetLogPath()
	if !filepath.IsAbs(logPath) {
		t.Errorf("GetLogPath() should return absolute path, got %s", logPath)
	}
	
	configPath := mock.GetConfigPath()
	if !filepath.IsAbs(configPath) {
		t.Errorf("GetConfigPath() should return absolute path, got %s", configPath)
	}
}

func TestMockPlatformService_GetDetectionMethods(t *testing.T) {
	mock := NewMockPlatformService()
	
	detectionMethod := mock.GetCommandDetectionMethod()
	if detectionMethod == "" {
		t.Error("GetCommandDetectionMethod() should not return empty string")
	}
	
	detectionCmd := mock.GetCommandDetectionCommand()
	if detectionCmd == "" {
		t.Error("GetCommandDetectionCommand() should not return empty string")
	}
	
	// Test that method and command are consistent
	if detectionMethod != detectionCmd {
		t.Errorf("Detection method %s should match command %s", detectionMethod, detectionCmd)
	}
}

func TestMockPlatformService_GetPermissions(t *testing.T) {
	mock := NewMockPlatformService()
	
	filePerms := mock.GetDefaultFilePermissions()
	if filePerms == 0 {
		t.Error("GetDefaultFilePermissions() should not return zero permissions")
	}
	
	dirPerms := mock.GetDefaultDirectoryPermissions()
	if dirPerms == 0 {
		t.Error("GetDefaultDirectoryPermissions() should not return zero permissions")
	}
}

func TestMockPlatformService_GetEnvironmentVariable(t *testing.T) {
	mock := NewMockPlatformService()
	
	testCases := []struct {
		key      string
		expected string
	}{
		{"HOME", mock.GetHomeDirectory()},
		{"USER", mock.GetCurrentUser()},
		{"TERM", "xterm-256color"},
		{"TERM_PROGRAM", "Mock Terminal"},
		{"PATH", "/usr/local/bin:/usr/bin:/bin"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			result := mock.GetEnvironmentVariable(tc.key)
			if result != tc.expected {
				t.Errorf("GetEnvironmentVariable(%s) = %s, expected %s", tc.key, result, tc.expected)
			}
		})
	}
	
	// Test unknown environment variable falls back to actual env
	unknownVar := mock.GetEnvironmentVariable("UNKNOWN_VAR_TEST")
	expectedUnknown := os.Getenv("UNKNOWN_VAR_TEST")
	if unknownVar != expectedUnknown {
		t.Errorf("GetEnvironmentVariable(UNKNOWN_VAR_TEST) should fallback to os.Getenv, got %s", unknownVar)
	}
}

func TestMockPlatformService_GetUserInfo(t *testing.T) {
	mock := NewMockPlatformService()
	
	homeDir := mock.GetHomeDirectory()
	if homeDir == "" {
		t.Error("GetHomeDirectory() should not return empty string")
	}
	
	currentUser := mock.GetCurrentUser()
	if currentUser == "" {
		t.Error("GetCurrentUser() should not return empty string")
	}
}

func TestMockPlatformService_SetPlatform(t *testing.T) {
	mock := NewMockPlatformService()
	
	// Test setting different platforms
	platforms := []PlatformType{PlatformWindows, PlatformLinux, PlatformUnknown}
	
	for _, platform := range platforms {
		mock.SetPlatform(platform)
		
		if mock.GetPlatform() != platform {
			t.Errorf("SetPlatform(%v) failed, got %v", platform, mock.GetPlatform())
		}
		
		if mock.GetPlatformName() != platform.String() {
			t.Errorf("SetPlatform(%v) should update platform name to %s, got %s", 
				platform, platform.String(), mock.GetPlatformName())
		}
	}
}

func TestMockPlatformService_SetSupportsClipboard(t *testing.T) {
	mock := NewMockPlatformService()
	
	// Test setting clipboard support
	mock.SetSupportsClipboard(false)
	if mock.SupportsClipboard() != false {
		t.Error("SetSupportsClipboard(false) failed")
	}
	
	mock.SetSupportsClipboard(true)
	if mock.SupportsClipboard() != true {
		t.Error("SetSupportsClipboard(true) failed")
	}
}

func TestMockPlatformService_SetClipboardMethod(t *testing.T) {
	mock := NewMockPlatformService()
	
	// Test setting different clipboard methods
	methods := []ClipboardMethod{
		ClipboardNative,
		ClipboardPbcopy,
		ClipboardXclip,
		ClipboardPowershell,
		ClipboardUnsupported,
	}
	
	for _, method := range methods {
		mock.SetClipboardMethod(method)
		
		if mock.GetClipboardMethod() != method {
			t.Errorf("SetClipboardMethod(%v) failed, got %v", method, mock.GetClipboardMethod())
		}
	}
}

func TestMockPlatformService_SetPaths(t *testing.T) {
	mock := NewMockPlatformService()
	
	// Test setting custom paths
	logPath := "/custom/log/path"
	configPath := "/custom/config/path"
	tempPath := "/custom/temp/path"
	cachePath := "/custom/cache/path"
	
	mock.SetPaths(logPath, configPath, tempPath, cachePath)
	
	if mock.GetLogPath() != logPath {
		t.Errorf("SetPaths() failed to set log path, got %s", mock.GetLogPath())
	}
	
	if mock.GetConfigPath() != configPath {
		t.Errorf("SetPaths() failed to set config path, got %s", mock.GetConfigPath())
	}
	
	if mock.GetTempPath() != tempPath {
		t.Errorf("SetPaths() failed to set temp path, got %s", mock.GetTempPath())
	}
	
	if mock.GetCachePath() != cachePath {
		t.Errorf("SetPaths() failed to set cache path, got %s", mock.GetCachePath())
	}
}

func TestMockPlatformService_SetDetectionCommand(t *testing.T) {
	mock := NewMockPlatformService()
	
	// Test setting custom detection command
	method := "custom-method"
	cmd := "custom-cmd"
	
	mock.SetDetectionCommand(method, cmd)
	
	if mock.GetCommandDetectionMethod() != method {
		t.Errorf("SetDetectionCommand() failed to set method, got %s", mock.GetCommandDetectionMethod())
	}
	
	if mock.GetCommandDetectionCommand() != cmd {
		t.Errorf("SetDetectionCommand() failed to set command, got %s", mock.GetCommandDetectionCommand())
	}
}

func TestMockPlatformService_SetUserInfo(t *testing.T) {
	mock := NewMockPlatformService()
	
	// Test setting custom user info
	username := "testuser123"
	homeDir := "/home/testuser123"
	
	mock.SetCurrentUser(username)
	mock.SetHomeDirectory(homeDir)
	
	if mock.GetCurrentUser() != username {
		t.Errorf("SetCurrentUser() failed, got %s", mock.GetCurrentUser())
	}
	
	if mock.GetHomeDirectory() != homeDir {
		t.Errorf("SetHomeDirectory() failed, got %s", mock.GetHomeDirectory())
	}
}

func TestGetMockPlatformService(t *testing.T) {
	mock := GetMockPlatformService()
	
	if mock == nil {
		t.Fatal("GetMockPlatformService() returned nil")
	}
	
	// Should create a mock service for the current OS
	platform := mock.GetPlatform()
	if platform == PlatformUnknown {
		t.Error("GetMockPlatformService() should not return unknown platform for current OS")
	}
}