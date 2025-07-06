package platform

import (
	"testing"
)

// TestPlatformIntegration tests the complete platform abstraction workflow
func TestPlatformIntegration(t *testing.T) {
	// Create a platform service factory
	factory := NewPlatformServiceFactoryDefault()
	if factory == nil {
		t.Fatal("Failed to create platform factory")
	}

	// Create a platform service for the current OS
	service := factory.CreatePlatformService()
	if service == nil {
		t.Fatal("Failed to create platform service")
	}

	// Test basic platform information
	platform := service.GetPlatform()
	if platform == PlatformUnknown {
		t.Error("Platform should not be unknown for current OS")
	}

	platformName := service.GetPlatformName()
	if platformName == "" {
		t.Error("Platform name should not be empty")
	}

	// Test path resolution
	logPath := service.GetLogPath()
	if logPath == "" {
		t.Error("Log path should not be empty")
	}

	configPath := service.GetConfigPath()
	if configPath == "" {
		t.Error("Config path should not be empty")
	}

	tempPath := service.GetTempPath()
	if tempPath == "" {
		t.Error("Temp path should not be empty")
	}

	cachePath := service.GetCachePath()
	if cachePath == "" {
		t.Error("Cache path should not be empty")
	}

	// Test command detection
	detectionCmd := service.GetCommandDetectionCommand()
	if detectionCmd == "" {
		t.Error("Command detection command should not be empty")
	}

	// Test permissions
	filePerms := service.GetDefaultFilePermissions()
	if filePerms == 0 {
		t.Error("File permissions should not be zero")
	}

	dirPerms := service.GetDefaultDirectoryPermissions()
	if dirPerms == 0 {
		t.Error("Directory permissions should not be zero")
	}

	// Test environment utilities
	homeDir := service.GetHomeDirectory()
	if homeDir == "" {
		t.Error("Home directory should not be empty")
	}

	currentUser := service.GetCurrentUser()
	if currentUser == "" {
		t.Error("Current user should not be empty")
	}

	t.Logf("Platform integration test passed for %s (%s)", platformName, platform.String())
	t.Logf("Log path: %s", logPath)
	t.Logf("Config path: %s", configPath)
	t.Logf("Detection command: %s", detectionCmd)
}

// TestAllPlatformServices tests all platform service implementations
func TestAllPlatformServices(t *testing.T) {
	factory := NewPlatformServiceFactoryDefault()

	platforms := []string{"darwin", "windows", "linux", "unknown"}

	for _, osName := range platforms {
		t.Run(osName, func(t *testing.T) {
			service := factory.CreatePlatformServiceForOS(osName)
			if service == nil {
				t.Fatal("Failed to create platform service for", osName)
			}

			// Test that all required methods work
			_ = service.GetPlatform()
			_ = service.GetPlatformName()
			_ = service.GetLogPath()
			_ = service.GetConfigPath()
			_ = service.GetTempPath()
			_ = service.GetCachePath()
			_ = service.GetCommandDetectionCommand()
			_ = service.SupportsClipboard()
			_ = service.GetClipboardMethod()
			_ = service.GetDefaultFilePermissions()
			_ = service.GetDefaultDirectoryPermissions()
			_ = service.GetHomeDirectory()
			_ = service.GetCurrentUser()
		})
	}
}