// Package platform provides a cross-platform abstraction layer for system-specific operations.
package platform

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// GenericPlatformService provides generic implementations of platform operations
// This serves as a fallback for unsupported platforms
type GenericPlatformService struct {
	logger *log.Logger
}

// NewGenericPlatformService creates a new generic platform service
func NewGenericPlatformService(logger *log.Logger) *GenericPlatformService {
	if logger == nil {
		logger = log.Default()
	}
	return &GenericPlatformService{
		logger: logger,
	}
}

// GetPlatform returns the platform type
func (g *GenericPlatformService) GetPlatform() PlatformType {
	return PlatformUnknown
}

// GetPlatformName returns the human-readable platform name
func (g *GenericPlatformService) GetPlatformName() string {
	return "Generic"
}

// GetLogPath returns a generic log directory path
func (g *GenericPlatformService) GetLogPath() string {
	homeDir := g.GetHomeDirectory()
	if homeDir == "" {
		// Fallback to temp directory if home directory is not available
		return filepath.Join(os.TempDir(), "mcp-hub", "logs")
	}
	return filepath.Join(homeDir, ".mcp-hub", "logs")
}

// GetConfigPath returns a generic config directory path
func (g *GenericPlatformService) GetConfigPath() string {
	homeDir := g.GetHomeDirectory()
	if homeDir == "" {
		// Fallback to temp directory if home directory is not available
		return filepath.Join(os.TempDir(), "mcp-hub")
	}
	return filepath.Join(homeDir, ".mcp-hub")
}

// GetTempPath returns the generic temporary directory path
func (g *GenericPlatformService) GetTempPath() string {
	tempDir := os.TempDir()
	return filepath.Join(tempDir, "mcp-hub")
}

// GetCachePath returns a generic cache directory path
func (g *GenericPlatformService) GetCachePath() string {
	homeDir := g.GetHomeDirectory()
	if homeDir == "" {
		// Fallback to temp directory if home directory is not available
		return g.GetTempPath()
	}
	return filepath.Join(homeDir, ".mcp-hub", "cache")
}

// GetCommandDetectionMethod returns the command detection method for generic platform
func (g *GenericPlatformService) GetCommandDetectionMethod() string {
	return whichCmd
}

// GetCommandDetectionCommand returns the command detection command for generic platform
func (g *GenericPlatformService) GetCommandDetectionCommand() string {
	return whichCmd
}

// SupportsClipboard returns false as generic platform doesn't support clipboard operations
func (g *GenericPlatformService) SupportsClipboard() bool {
	return false
}

// GetClipboardMethod returns the clipboard method for generic platform
func (g *GenericPlatformService) GetClipboardMethod() ClipboardMethod {
	return ClipboardUnsupported
}

// GetDefaultFilePermissions returns the default file permissions for generic platform
func (g *GenericPlatformService) GetDefaultFilePermissions() os.FileMode {
	return 0644 // User read/write, group and others read
}

// GetDefaultDirectoryPermissions returns the default directory permissions for generic platform
func (g *GenericPlatformService) GetDefaultDirectoryPermissions() os.FileMode {
	return 0755 // User read/write/execute, group and others read/execute
}

// GetEnvironmentVariable returns the value of an environment variable
func (g *GenericPlatformService) GetEnvironmentVariable(key string) string {
	return os.Getenv(key)
}

// GetHomeDirectory returns the user's home directory
func (g *GenericPlatformService) GetHomeDirectory() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir
	}
	// Fallback to HOME environment variable
	return os.Getenv("HOME")
}

// GetCurrentUser returns the current user's username
func (g *GenericPlatformService) GetCurrentUser() string {
	if currentUser, err := user.Current(); err == nil {
		return currentUser.Username
	}
	// Fallback to USER environment variable
	return os.Getenv("USER")
}