package platform

import (
	"testing"
)

func TestPlatformType_String(t *testing.T) {
	tests := []struct {
		platform PlatformType
		expected string
	}{
		{PlatformDarwin, "darwin"},
		{PlatformWindows, "windows"},
		{PlatformLinux, "linux"},
		{PlatformUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.platform.String()
			if result != tt.expected {
				t.Errorf("PlatformType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestClipboardMethod_String(t *testing.T) {
	tests := []struct {
		method   ClipboardMethod
		expected string
	}{
		{ClipboardNative, "native"},
		{ClipboardPbcopy, "pbcopy"},
		{ClipboardXclip, "xclip"},
		{ClipboardPowershell, "powershell"},
		{ClipboardUnsupported, "unsupported"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.method.String()
			if result != tt.expected {
				t.Errorf("ClipboardMethod.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}