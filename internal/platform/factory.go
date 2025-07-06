// Package platform provides a cross-platform abstraction layer for system-specific operations.
package platform

import (
	"log"
	"runtime"
)

// PlatformServiceFactory creates platform-specific service instances
type PlatformServiceFactory struct {
	logger *log.Logger
}

// NewPlatformServiceFactory creates a new platform service factory
func NewPlatformServiceFactory(logger *log.Logger) *PlatformServiceFactory {
	return &PlatformServiceFactory{
		logger: logger,
	}
}

// NewPlatformServiceFactoryDefault creates a platform service factory with default logger
func NewPlatformServiceFactoryDefault() *PlatformServiceFactory {
	return &PlatformServiceFactory{
		logger: log.Default(),
	}
}

// CreatePlatformService creates a platform-specific service instance based on the current runtime
func (f *PlatformServiceFactory) CreatePlatformService() PlatformService {
	return f.CreatePlatformServiceForOS(runtime.GOOS)
}

// CreatePlatformServiceForOS creates a platform-specific service instance for the given OS
func (f *PlatformServiceFactory) CreatePlatformServiceForOS(osName string) PlatformService {
	switch osName {
	case "darwin":
		return NewDarwinPlatformService(f.logger)
	case "windows":
		return NewWindowsPlatformService(f.logger)
	case "linux":
		return NewLinuxPlatformService(f.logger)
	default:
		return NewGenericPlatformService(f.logger)
	}
}

// GetPlatformTypeFromOS converts an OS name to a PlatformType
func GetPlatformTypeFromOS(osName string) PlatformType {
	switch osName {
	case "darwin":
		return PlatformDarwin
	case "windows":
		return PlatformWindows
	case "linux":
		return PlatformLinux
	default:
		return PlatformUnknown
	}
}

// GetCurrentPlatformType returns the current platform type
func GetCurrentPlatformType() PlatformType {
	return GetPlatformTypeFromOS(runtime.GOOS)
}

// IsPlatformSupported checks if a platform is supported
func IsPlatformSupported(platform PlatformType) bool {
	return platform != PlatformUnknown
}

// GetSupportedPlatforms returns a list of supported platforms
func GetSupportedPlatforms() []PlatformType {
	return []PlatformType{
		PlatformDarwin,
		PlatformWindows,
		PlatformLinux,
	}
}