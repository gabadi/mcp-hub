package platform

import (
	"strings"
	"testing"
)


func TestNewDarwinPlatformService(t *testing.T) {
	service := NewDarwinPlatformService(nil)

	if service == nil {
		t.Fatal("NewDarwinPlatformService() returned nil")
	}

	if service.logger == nil {
		t.Error("NewDarwinPlatformService() should set default logger when nil passed")
	}
}

func TestDarwinPlatformService_GetPlatform(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	platform := service.GetPlatform()

	if platform != PlatformDarwin {
		t.Errorf("GetPlatform() = %v, want %v", platform, PlatformDarwin)
	}
}

func TestDarwinPlatformService_GetPlatformName(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	name := service.GetPlatformName()

	if name != "macOS" {
		t.Errorf("GetPlatformName() = %v, want %v", name, "macOS")
	}
}

func TestDarwinPlatformService_GetLogPath(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	logPath := service.GetLogPath()

	if !strings.Contains(logPath, "Library/Logs/mcp-hub") {
		t.Errorf("GetLogPath() = %v, should contain 'Library/Logs/mcp-hub'", logPath)
	}
}

func TestDarwinPlatformService_GetConfigPath(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	configPath := service.GetConfigPath()

	if !strings.Contains(configPath, "mcp-hub") {
		t.Errorf("GetConfigPath() = %v, should contain 'mcp-hub'", configPath)
	}
}

func TestDarwinPlatformService_GetCommandDetectionCommand(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	cmd := service.GetCommandDetectionCommand()

	if cmd != whichCmd {
		t.Errorf("GetCommandDetectionCommand() = %v, want %v", cmd, whichCmd)
	}
}

func TestDarwinPlatformService_SupportsClipboard(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	supports := service.SupportsClipboard()

	if !supports {
		t.Error("SupportsClipboard() = false, want true for macOS")
	}
}

func TestDarwinPlatformService_GetClipboardMethod(t *testing.T) {
	service := NewDarwinPlatformService(nil)
	method := service.GetClipboardMethod()

	// Should be either pbcopy or native
	if method != ClipboardPbcopy && method != ClipboardNative {
		t.Errorf("GetClipboardMethod() = %v, want either ClipboardPbcopy or ClipboardNative", method)
	}
}

func TestDarwinPlatformService_GetDefaultPermissions(t *testing.T) {
	service := NewDarwinPlatformService(nil)

	filePerms := service.GetDefaultFilePermissions()
	if filePerms != 0600 {
		t.Errorf("GetDefaultFilePermissions() = %v, want 0600", filePerms)
	}

	dirPerms := service.GetDefaultDirectoryPermissions()
	if dirPerms != 0700 {
		t.Errorf("GetDefaultDirectoryPermissions() = %v, want 0700", dirPerms)
	}
}