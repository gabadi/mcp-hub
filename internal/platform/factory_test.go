package platform

import (
	"log"
	"testing"
)

func TestNewPlatformServiceFactory(t *testing.T) {
	logger := log.Default()
	factory := NewPlatformServiceFactory(logger)

	if factory == nil {
		t.Fatal("NewPlatformServiceFactory() returned nil")
	}

	if factory.logger != logger {
		t.Error("NewPlatformServiceFactory() did not set logger correctly")
	}
}

func TestNewPlatformServiceFactoryDefault(t *testing.T) {
	factory := NewPlatformServiceFactoryDefault()

	if factory == nil {
		t.Fatal("NewPlatformServiceFactoryDefault() returned nil")
	}

	if factory.logger == nil {
		t.Error("NewPlatformServiceFactoryDefault() did not set default logger")
	}
}

func TestCreatePlatformServiceForOS(t *testing.T) {
	factory := NewPlatformServiceFactoryDefault()

	tests := []struct {
		osName   string
		expected PlatformType
	}{
		{"darwin", PlatformDarwin},
		{"windows", PlatformWindows},
		{"linux", PlatformLinux},
		{"unknown", PlatformUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.osName, func(t *testing.T) {
			service := factory.CreatePlatformServiceForOS(tt.osName)
			if service == nil {
				t.Fatal("CreatePlatformServiceForOS() returned nil")
			}

			platform := service.GetPlatform()
			if platform != tt.expected {
				t.Errorf("CreatePlatformServiceForOS() platform = %v, want %v", platform, tt.expected)
			}
		})
	}
}

func TestGetPlatformTypeFromOS(t *testing.T) {
	tests := []struct {
		osName   string
		expected PlatformType
	}{
		{"darwin", PlatformDarwin},
		{"windows", PlatformWindows},
		{"linux", PlatformLinux},
		{"freebsd", PlatformUnknown},
		{"", PlatformUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.osName, func(t *testing.T) {
			result := GetPlatformTypeFromOS(tt.osName)
			if result != tt.expected {
				t.Errorf("GetPlatformTypeFromOS() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsPlatformSupported(t *testing.T) {
	tests := []struct {
		platform PlatformType
		expected bool
	}{
		{PlatformDarwin, true},
		{PlatformWindows, true},
		{PlatformLinux, true},
		{PlatformUnknown, false},
	}

	for _, tt := range tests {
		t.Run(tt.platform.String(), func(t *testing.T) {
			result := IsPlatformSupported(tt.platform)
			if result != tt.expected {
				t.Errorf("IsPlatformSupported() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetSupportedPlatforms(t *testing.T) {
	platforms := GetSupportedPlatforms()

	expected := []PlatformType{PlatformDarwin, PlatformWindows, PlatformLinux}
	if len(platforms) != len(expected) {
		t.Errorf("GetSupportedPlatforms() length = %v, want %v", len(platforms), len(expected))
	}

	for i, platform := range platforms {
		if platform != expected[i] {
			t.Errorf("GetSupportedPlatforms()[%d] = %v, want %v", i, platform, expected[i])
		}
	}
}

func TestGetCurrentPlatformType(t *testing.T) {
	platformType := GetCurrentPlatformType()
	
	// Should not be PlatformUnknown on supported platforms
	if platformType == PlatformUnknown {
		t.Error("GetCurrentPlatformType() should not return PlatformUnknown on supported platforms")
	}
	
	// Should be one of the supported platforms
	supportedPlatforms := GetSupportedPlatforms()
	isSupported := false
	for _, supported := range supportedPlatforms {
		if platformType == supported {
			isSupported = true
			break
		}
	}
	
	if !isSupported {
		t.Errorf("GetCurrentPlatformType() = %v, should be one of %v", platformType, supportedPlatforms)
	}
}