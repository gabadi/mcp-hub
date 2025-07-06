// Package platform provides a cross-platform abstraction layer for system-specific operations.
package platform

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

// DarwinPlatformService provides macOS-specific implementations of platform operations
type DarwinPlatformService struct {
	logger *log.Logger
}

// NewDarwinPlatformService creates a new Darwin platform service
func NewDarwinPlatformService(logger *log.Logger) *DarwinPlatformService {
	if logger == nil {
		logger = log.Default()
	}
	return &DarwinPlatformService{
		logger: logger,
	}
}

// GetPlatform returns the platform type
func (d *DarwinPlatformService) GetPlatform() PlatformType {
	return PlatformDarwin
}

// GetPlatformName returns the human-readable platform name
func (d *DarwinPlatformService) GetPlatformName() string {
	return "macOS"
}

// GetLogPath returns the macOS-specific log directory path
func (d *DarwinPlatformService) GetLogPath() string {
	homeDir := d.GetHomeDirectory()
	if homeDir == "" {
		// Fallback to /tmp if home directory is not available
		return tmpMcpHub
	}
	return filepath.Join(homeDir, "Library", "Logs", "mcp-hub")
}

// GetConfigPath returns the macOS-specific config directory path
func (d *DarwinPlatformService) GetConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home directory if UserConfigDir fails
		homeDir := d.GetHomeDirectory()
		if homeDir != "" {
			return filepath.Join(homeDir, ".config", "mcp-hub")
		}
		return tmpMcpHub
	}
	return filepath.Join(configDir, "mcp-hub")
}

// GetTempPath returns the macOS-specific temporary directory path
func (d *DarwinPlatformService) GetTempPath() string {
	tempDir := os.TempDir()
	return filepath.Join(tempDir, "mcp-hub")
}

// GetCachePath returns the macOS-specific cache directory path
func (d *DarwinPlatformService) GetCachePath() string {
	homeDir := d.GetHomeDirectory()
	if homeDir == "" {
		return d.GetTempPath()
	}
	return filepath.Join(homeDir, "Library", "Caches", "mcp-hub")
}

// GetCommandDetectionMethod returns the command detection method for macOS
func (d *DarwinPlatformService) GetCommandDetectionMethod() string {
	return whichCmd
}

// GetCommandDetectionCommand returns the command detection command for macOS
func (d *DarwinPlatformService) GetCommandDetectionCommand() string {
	return whichCmd
}

// SupportsClipboard returns true as macOS supports clipboard operations
func (d *DarwinPlatformService) SupportsClipboard() bool {
	return true
}

// GetClipboardMethod returns the clipboard method for macOS
func (d *DarwinPlatformService) GetClipboardMethod() ClipboardMethod {
	// Check if pbcopy is available
	if _, err := exec.LookPath("pbcopy"); err == nil {
		return ClipboardPbcopy
	}
	// Fallback to native if pbcopy is not available
	return ClipboardNative
}

// GetDefaultFilePermissions returns the default file permissions for macOS
func (d *DarwinPlatformService) GetDefaultFilePermissions() os.FileMode {
	return 0600 // User read/write only
}

// GetDefaultDirectoryPermissions returns the default directory permissions for macOS
func (d *DarwinPlatformService) GetDefaultDirectoryPermissions() os.FileMode {
	return 0700 // User read/write/execute only
}

// GetEnvironmentVariable returns the value of an environment variable
func (d *DarwinPlatformService) GetEnvironmentVariable(key string) string {
	return os.Getenv(key)
}

// GetHomeDirectory returns the user's home directory
func (d *DarwinPlatformService) GetHomeDirectory() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir
	}
	// Fallback to HOME environment variable
	return os.Getenv("HOME")
}

// GetCurrentUser returns the current user's username
func (d *DarwinPlatformService) GetCurrentUser() string {
	if currentUser, err := user.Current(); err == nil {
		return currentUser.Username
	}
	// Fallback to USER environment variable
	return os.Getenv("USER")
}