// Package platform provides mock implementation for testing
package platform

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// MockPlatformService provides a mock implementation of PlatformService for testing
type MockPlatformService struct {
	platform        PlatformType
	platformName    string
	logPath         string
	configPath      string
	tempPath        string
	cachePath       string
	detectionMethod string
	detectionCmd    string
	supportsClip    bool
	clipboardMethod ClipboardMethod
	filePerms       os.FileMode
	dirPerms        os.FileMode
	homeDir         string
	currentUser     string
}

// NewMockPlatformService creates a new mock platform service with default values
func NewMockPlatformService() *MockPlatformService {
	mockUser := "testuser"
	mockHome := "/home/testuser"
	
	// Use current user info if available for more realistic testing
	if currentUser, err := user.Current(); err == nil {
		mockUser = currentUser.Username
		mockHome = currentUser.HomeDir
	}
	
	return &MockPlatformService{
		platform:        PlatformDarwin, // Default to darwin for testing
		platformName:    "darwin",
		logPath:         filepath.Join(mockHome, "Library", "Logs", "mcp-hub"),
		configPath:      filepath.Join(mockHome, "Library", "Application Support", "mcp-hub"),
		tempPath:        "/tmp",
		cachePath:       filepath.Join(mockHome, "Library", "Caches", "mcp-hub"),
		detectionMethod: "which",
		detectionCmd:    "which",
		supportsClip:    true,
		clipboardMethod: ClipboardPbcopy,
		filePerms:       0644,
		dirPerms:        0755,
		homeDir:         mockHome,
		currentUser:     mockUser,
	}
}

// NewMockPlatformServiceForOS creates a mock platform service for a specific OS
func NewMockPlatformServiceForOS(osName string) *MockPlatformService {
	mock := NewMockPlatformService()
	
	switch osName {
	case "darwin":
		mock.platform = PlatformDarwin
		mock.platformName = "darwin"
		mock.detectionMethod = "which"
		mock.detectionCmd = "which"
		mock.supportsClip = true
		mock.clipboardMethod = ClipboardPbcopy
	case "windows":
		mock.platform = PlatformWindows
		mock.platformName = "windows"
		mock.detectionMethod = "where"
		mock.detectionCmd = "where"
		mock.supportsClip = true
		mock.clipboardMethod = ClipboardPowershell
		mock.logPath = filepath.Join(mock.homeDir, "AppData", "Local", "mcp-hub", "logs")
		mock.configPath = filepath.Join(mock.homeDir, "AppData", "Local", "mcp-hub")
		mock.tempPath = filepath.Join(mock.homeDir, "AppData", "Local", "Temp")
		mock.cachePath = filepath.Join(mock.homeDir, "AppData", "Local", "mcp-hub", "cache")
	case "linux":
		mock.platform = PlatformLinux
		mock.platformName = "linux"
		mock.detectionMethod = "which"
		mock.detectionCmd = "which"
		mock.supportsClip = true
		mock.clipboardMethod = ClipboardXclip
		mock.logPath = filepath.Join(mock.homeDir, ".local", "share", "mcp-hub", "logs")
		mock.configPath = filepath.Join(mock.homeDir, ".config", "mcp-hub")
		mock.tempPath = "/tmp"
		mock.cachePath = filepath.Join(mock.homeDir, ".cache", "mcp-hub")
	default:
		mock.platform = PlatformUnknown
		mock.platformName = "unknown"
		mock.detectionMethod = "which"
		mock.detectionCmd = "which"
		mock.supportsClip = false
		mock.clipboardMethod = ClipboardUnsupported
	}
	
	return mock
}

// GetPlatform returns the platform type
func (m *MockPlatformService) GetPlatform() PlatformType {
	return m.platform
}

// GetPlatformName returns the platform name
func (m *MockPlatformService) GetPlatformName() string {
	return m.platformName
}

// GetLogPath returns the log directory path
func (m *MockPlatformService) GetLogPath() string {
	return m.logPath
}

// GetConfigPath returns the config directory path
func (m *MockPlatformService) GetConfigPath() string {
	return m.configPath
}

// GetTempPath returns the temp directory path
func (m *MockPlatformService) GetTempPath() string {
	return m.tempPath
}

// GetCachePath returns the cache directory path
func (m *MockPlatformService) GetCachePath() string {
	return m.cachePath
}

// GetCommandDetectionMethod returns the command detection method
func (m *MockPlatformService) GetCommandDetectionMethod() string {
	return m.detectionMethod
}

// GetCommandDetectionCommand returns the command detection command
func (m *MockPlatformService) GetCommandDetectionCommand() string {
	return m.detectionCmd
}

// SupportsClipboard returns whether clipboard operations are supported
func (m *MockPlatformService) SupportsClipboard() bool {
	return m.supportsClip
}

// GetClipboardMethod returns the clipboard method
func (m *MockPlatformService) GetClipboardMethod() ClipboardMethod {
	return m.clipboardMethod
}

// GetDefaultFilePermissions returns the default file permissions
func (m *MockPlatformService) GetDefaultFilePermissions() os.FileMode {
	return m.filePerms
}

// GetDefaultDirectoryPermissions returns the default directory permissions
func (m *MockPlatformService) GetDefaultDirectoryPermissions() os.FileMode {
	return m.dirPerms
}

// GetEnvironmentVariable returns an environment variable value
func (m *MockPlatformService) GetEnvironmentVariable(key string) string {
	// Return some mock values for common variables
	switch key {
	case "HOME":
		return m.homeDir
	case "USER":
		return m.currentUser
	case "TERM":
		return "xterm-256color"
	case "TERM_PROGRAM":
		return "Mock Terminal"
	case "PATH":
		return "/usr/local/bin:/usr/bin:/bin"
	default:
		return os.Getenv(key) // Fall back to actual env var
	}
}

// GetHomeDirectory returns the home directory
func (m *MockPlatformService) GetHomeDirectory() string {
	return m.homeDir
}

// GetCurrentUser returns the current user
func (m *MockPlatformService) GetCurrentUser() string {
	return m.currentUser
}

// SetPlatform allows setting the platform type for testing
func (m *MockPlatformService) SetPlatform(platform PlatformType) {
	m.platform = platform
	m.platformName = platform.String()
}

// SetSupportsClipboard allows setting clipboard support for testing
func (m *MockPlatformService) SetSupportsClipboard(supports bool) {
	m.supportsClip = supports
}

// SetClipboardMethod allows setting clipboard method for testing
func (m *MockPlatformService) SetClipboardMethod(method ClipboardMethod) {
	m.clipboardMethod = method
}

// SetPaths allows setting custom paths for testing
func (m *MockPlatformService) SetPaths(logPath, configPath, tempPath, cachePath string) {
	m.logPath = logPath
	m.configPath = configPath
	m.tempPath = tempPath
	m.cachePath = cachePath
}

// SetDetectionCommand allows setting custom detection command for testing
func (m *MockPlatformService) SetDetectionCommand(method, cmd string) {
	m.detectionMethod = method
	m.detectionCmd = cmd
}

// SetCurrentUser allows setting the current user for testing
func (m *MockPlatformService) SetCurrentUser(username string) {
	m.currentUser = username
}

// SetHomeDirectory allows setting home directory for testing
func (m *MockPlatformService) SetHomeDirectory(homeDir string) {
	m.homeDir = homeDir
}

// GetMockPlatformService returns a mock platform service for the current OS
func GetMockPlatformService() *MockPlatformService {
	return NewMockPlatformServiceForOS(runtime.GOOS)
}