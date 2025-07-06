// Package platform provides a cross-platform abstraction layer for system-specific operations.
package platform

import (
	"os"
)

// PlatformType represents the supported operating system platforms
type PlatformType int

const (
	// PlatformUnknown represents an unknown or unsupported platform
	PlatformUnknown PlatformType = iota
	// PlatformDarwin represents the macOS platform
	PlatformDarwin
	// PlatformWindows represents the Windows platform
	PlatformWindows
	// PlatformLinux represents the Linux platform
	PlatformLinux
)

// ClipboardMethod represents the clipboard access method for the platform
type ClipboardMethod int

const (
	// ClipboardUnsupported indicates clipboard operations are not supported
	ClipboardUnsupported ClipboardMethod = iota
	// ClipboardNative indicates native clipboard support
	ClipboardNative
	// ClipboardPbcopy indicates macOS pbcopy/pbpaste clipboard support
	ClipboardPbcopy
	// ClipboardXclip indicates Linux xclip clipboard support
	ClipboardXclip
	// ClipboardPowershell indicates Windows PowerShell clipboard support
	ClipboardPowershell
)

// PlatformService defines the interface for platform-specific operations
type PlatformService interface {
	// Platform identification
	GetPlatform() PlatformType
	GetPlatformName() string
	
	// Path resolution
	GetLogPath() string
	GetConfigPath() string
	GetTempPath() string
	GetCachePath() string
	
	// Command utilities
	GetCommandDetectionMethod() string
	GetCommandDetectionCommand() string
	
	// Clipboard operations
	SupportsClipboard() bool
	GetClipboardMethod() ClipboardMethod
	
	// File operations
	GetDefaultFilePermissions() os.FileMode
	GetDefaultDirectoryPermissions() os.FileMode
	
	// Environment utilities
	GetEnvironmentVariable(key string) string
	GetHomeDirectory() string
	GetCurrentUser() string
}

// String returns the string representation of a PlatformType
func (pt PlatformType) String() string {
	switch pt {
	case PlatformDarwin:
		return osDarwin
	case PlatformWindows:
		return osWindows
	case PlatformLinux:
		return osLinux
	case PlatformUnknown:
		return unknownPlatform
	default:
		return unknownPlatform
	}
}

// String returns the string representation of a ClipboardMethod
func (cm ClipboardMethod) String() string {
	switch cm {
	case ClipboardNative:
		return "native"
	case ClipboardPbcopy:
		return "pbcopy"
	case ClipboardXclip:
		return "xclip"
	case ClipboardPowershell:
		return "powershell"
	case ClipboardUnsupported:
		return "unsupported"
	default:
		return "unsupported"
	}
}