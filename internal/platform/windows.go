// Package platform provides a cross-platform abstraction layer for system-specific operations.
package platform

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

// WindowsPlatformService provides Windows-specific implementations of platform operations
type WindowsPlatformService struct {
	logger *log.Logger
}

// NewWindowsPlatformService creates a new Windows platform service
func NewWindowsPlatformService(logger *log.Logger) *WindowsPlatformService {
	if logger == nil {
		logger = log.Default()
	}
	return &WindowsPlatformService{
		logger: logger,
	}
}

// GetPlatform returns the platform type
func (w *WindowsPlatformService) GetPlatform() PlatformType {
	return PlatformWindows
}

// GetPlatformName returns the human-readable platform name
func (w *WindowsPlatformService) GetPlatformName() string {
	return "Windows"
}

// GetLogPath returns the Windows-specific log directory path
func (w *WindowsPlatformService) GetLogPath() string {
	appData := w.GetEnvironmentVariable("APPDATA")
	if appData == "" {
		// Fallback to temp directory if APPDATA is not available
		return filepath.Join(os.TempDir(), "mcp-hub", "logs")
	}
	return filepath.Join(appData, "mcp-hub", "logs")
}

// GetConfigPath returns the Windows-specific config directory path
func (w *WindowsPlatformService) GetConfigPath() string {
	appData := w.GetEnvironmentVariable("APPDATA")
	if appData == "" {
		// Fallback to temp directory if APPDATA is not available
		return filepath.Join(os.TempDir(), "mcp-hub")
	}
	return filepath.Join(appData, "mcp-hub")
}

// GetTempPath returns the Windows-specific temporary directory path
func (w *WindowsPlatformService) GetTempPath() string {
	tempDir := os.TempDir()
	return filepath.Join(tempDir, "mcp-hub")
}

// GetCachePath returns the Windows-specific cache directory path
func (w *WindowsPlatformService) GetCachePath() string {
	localAppData := w.GetEnvironmentVariable("LOCALAPPDATA")
	if localAppData == "" {
		// Fallback to APPDATA if LOCALAPPDATA is not available
		appData := w.GetEnvironmentVariable("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "mcp-hub", "cache")
		}
		// Final fallback to temp directory
		return filepath.Join(os.TempDir(), "mcp-hub", "cache")
	}
	return filepath.Join(localAppData, "mcp-hub", "cache")
}

// GetCommandDetectionMethod returns the command detection method for Windows
func (w *WindowsPlatformService) GetCommandDetectionMethod() string {
	return whereCmd
}

// GetCommandDetectionCommand returns the command detection command for Windows
func (w *WindowsPlatformService) GetCommandDetectionCommand() string {
	return whereCmd
}

// SupportsClipboard returns true as Windows supports clipboard operations
func (w *WindowsPlatformService) SupportsClipboard() bool {
	return true
}

// GetClipboardMethod returns the clipboard method for Windows
func (w *WindowsPlatformService) GetClipboardMethod() ClipboardMethod {
	// Check if PowerShell is available for clipboard operations
	if _, err := exec.LookPath("powershell"); err == nil {
		return ClipboardPowershell
	}
	// Fallback to native if PowerShell is not available
	return ClipboardNative
}

// GetDefaultFilePermissions returns the default file permissions for Windows
func (w *WindowsPlatformService) GetDefaultFilePermissions() os.FileMode {
	return 0644 // User read/write, group and others read
}

// GetDefaultDirectoryPermissions returns the default directory permissions for Windows
func (w *WindowsPlatformService) GetDefaultDirectoryPermissions() os.FileMode {
	return 0755 // User read/write/execute, group and others read/execute
}

// GetEnvironmentVariable returns the value of an environment variable
func (w *WindowsPlatformService) GetEnvironmentVariable(key string) string {
	return os.Getenv(key)
}

// GetHomeDirectory returns the user's home directory
func (w *WindowsPlatformService) GetHomeDirectory() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir
	}
	// Fallback to USERPROFILE environment variable
	userProfile := os.Getenv("USERPROFILE")
	if userProfile != "" {
		return userProfile
	}
	// Final fallback to HOMEDRIVE + HOMEPATH
	homeDrive := os.Getenv("HOMEDRIVE")
	homePath := os.Getenv("HOMEPATH")
	if homeDrive != "" && homePath != "" {
		return filepath.Join(homeDrive, homePath)
	}
	return ""
}

// GetCurrentUser returns the current user's username
func (w *WindowsPlatformService) GetCurrentUser() string {
	if currentUser, err := user.Current(); err == nil {
		return currentUser.Username
	}
	// Fallback to USERNAME environment variable
	return os.Getenv("USERNAME")
}