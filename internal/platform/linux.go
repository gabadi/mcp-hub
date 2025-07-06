// Package platform provides a cross-platform abstraction layer for system-specific operations.
package platform

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

// LinuxPlatformService provides Linux-specific implementations of platform operations
type LinuxPlatformService struct {
	logger *log.Logger
}

// NewLinuxPlatformService creates a new Linux platform service
func NewLinuxPlatformService(logger *log.Logger) *LinuxPlatformService {
	if logger == nil {
		logger = log.Default()
	}
	return &LinuxPlatformService{
		logger: logger,
	}
}

// GetPlatform returns the platform type
func (l *LinuxPlatformService) GetPlatform() PlatformType {
	return PlatformLinux
}

// GetPlatformName returns the human-readable platform name
func (l *LinuxPlatformService) GetPlatformName() string {
	return "Linux"
}

// GetLogPath returns the Linux-specific log directory path following XDG specifications
func (l *LinuxPlatformService) GetLogPath() string {
	// Follow XDG Base Directory Specification
	xdgDataHome := l.GetEnvironmentVariable("XDG_DATA_HOME")
	if xdgDataHome != "" {
		return filepath.Join(xdgDataHome, "mcp-hub", "logs")
	}
	
	homeDir := l.GetHomeDirectory()
	if homeDir == "" {
		// Fallback to /tmp if home directory is not available
		return "/tmp/mcp-hub/logs"
	}
	return filepath.Join(homeDir, ".local", "share", "mcp-hub", "logs")
}

// GetConfigPath returns the Linux-specific config directory path following XDG specifications
func (l *LinuxPlatformService) GetConfigPath() string {
	// Follow XDG Base Directory Specification
	xdgConfigHome := l.GetEnvironmentVariable("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "mcp-hub")
	}
	
	homeDir := l.GetHomeDirectory()
	if homeDir == "" {
		// Fallback to /tmp if home directory is not available
		return "/tmp/mcp-hub"
	}
	return filepath.Join(homeDir, ".config", "mcp-hub")
}

// GetTempPath returns the Linux-specific temporary directory path
func (l *LinuxPlatformService) GetTempPath() string {
	tempDir := os.TempDir()
	return filepath.Join(tempDir, "mcp-hub")
}

// GetCachePath returns the Linux-specific cache directory path following XDG specifications
func (l *LinuxPlatformService) GetCachePath() string {
	// Follow XDG Base Directory Specification
	xdgCacheHome := l.GetEnvironmentVariable("XDG_CACHE_HOME")
	if xdgCacheHome != "" {
		return filepath.Join(xdgCacheHome, "mcp-hub")
	}
	
	homeDir := l.GetHomeDirectory()
	if homeDir == "" {
		// Fallback to temp directory if home directory is not available
		return l.GetTempPath()
	}
	return filepath.Join(homeDir, ".cache", "mcp-hub")
}

// GetCommandDetectionMethod returns the command detection method for Linux
func (l *LinuxPlatformService) GetCommandDetectionMethod() string {
	return whichCmd
}

// GetCommandDetectionCommand returns the command detection command for Linux
func (l *LinuxPlatformService) GetCommandDetectionCommand() string {
	return whichCmd
}

// SupportsClipboard returns true if clipboard operations are supported on Linux
func (l *LinuxPlatformService) SupportsClipboard() bool {
	// Check for xclip availability
	if _, err := exec.LookPath("xclip"); err == nil {
		return true
	}
	// Check for xsel availability as alternative
	if _, err := exec.LookPath("xsel"); err == nil {
		return true
	}
	// Check for wl-clipboard availability for Wayland
	if _, err := exec.LookPath("wl-copy"); err == nil {
		return true
	}
	return false
}

// GetClipboardMethod returns the clipboard method for Linux
func (l *LinuxPlatformService) GetClipboardMethod() ClipboardMethod {
	// Check for xclip first (most common)
	if _, err := exec.LookPath("xclip"); err == nil {
		return ClipboardXclip
	}
	// Check for xsel as alternative
	if _, err := exec.LookPath("xsel"); err == nil {
		return ClipboardXclip // Use same method type for X11 clipboard tools
	}
	// Check for wl-clipboard for Wayland
	if _, err := exec.LookPath("wl-copy"); err == nil {
		return ClipboardNative // Use native for Wayland
	}
	return ClipboardUnsupported
}

// GetDefaultFilePermissions returns the default file permissions for Linux
func (l *LinuxPlatformService) GetDefaultFilePermissions() os.FileMode {
	return 0600 // User read/write only
}

// GetDefaultDirectoryPermissions returns the default directory permissions for Linux
func (l *LinuxPlatformService) GetDefaultDirectoryPermissions() os.FileMode {
	return 0700 // User read/write/execute only
}

// GetEnvironmentVariable returns the value of an environment variable
func (l *LinuxPlatformService) GetEnvironmentVariable(key string) string {
	return os.Getenv(key)
}

// GetHomeDirectory returns the user's home directory
func (l *LinuxPlatformService) GetHomeDirectory() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir
	}
	// Fallback to HOME environment variable
	return os.Getenv("HOME")
}

// GetCurrentUser returns the current user's username
func (l *LinuxPlatformService) GetCurrentUser() string {
	if currentUser, err := user.Current(); err == nil {
		return currentUser.Username
	}
	// Fallback to USER environment variable
	return os.Getenv("USER")
}