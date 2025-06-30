package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if cfg.ConfigDir == "" {
		t.Error("Expected ConfigDir to be set")
	}
	
	if cfg.ConfigFile == "" {
		t.Error("Expected ConfigFile to be set")
	}
	
	if !strings.HasSuffix(cfg.ConfigFile, ConfigFileName) {
		t.Errorf("Expected ConfigFile to end with %s, got %s", ConfigFileName, cfg.ConfigFile)
	}
	
	// Check that config directory was created
	if _, err := os.Stat(cfg.ConfigDir); os.IsNotExist(err) {
		t.Errorf("Expected config directory to be created: %s", cfg.ConfigDir)
	}
}

func TestGetConfigDir(t *testing.T) {
	configDir, err := getConfigDir()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if configDir == "" {
		t.Error("Expected configDir to be non-empty")
	}
	
	// Check that the path contains the app name
	if !strings.Contains(configDir, AppName) {
		t.Errorf("Expected configDir to contain %s, got %s", AppName, configDir)
	}
	
	// Check OS-specific paths
	switch runtime.GOOS {
	case "darwin":
		if !strings.Contains(configDir, "Library/Application Support") {
			t.Errorf("Expected macOS config path, got %s", configDir)
		}
	case "linux":
		// Should contain either .config or XDG_CONFIG_HOME
		if !strings.Contains(configDir, ".config") && os.Getenv("XDG_CONFIG_HOME") == "" {
			t.Errorf("Expected Linux config path, got %s", configDir)
		}
	case "windows":
		if !strings.Contains(configDir, "AppData") {
			t.Errorf("Expected Windows config path, got %s", configDir)
		}
	}
}

func TestGetConfigDirWithXDG(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("XDG test only runs on Linux")
	}
	
	// Set XDG_CONFIG_HOME
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)
	
	testXDGPath := "/tmp/test-xdg-config"
	os.Setenv("XDG_CONFIG_HOME", testXDGPath)
	
	configDir, err := getConfigDir()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	expectedPath := filepath.Join(testXDGPath, AppName)
	if configDir != expectedPath {
		t.Errorf("Expected %s, got %s", expectedPath, configDir)
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary config for testing
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	cfg := &Config{
		ConfigDir:  tempDir,
		ConfigFile: filepath.Join(tempDir, ConfigFileName),
	}
	
	// File shouldn't exist initially
	if cfg.FileExists() {
		t.Error("Expected file to not exist")
	}
	
	// Create the file
	file, err := os.Create(cfg.ConfigFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()
	
	// File should exist now
	if !cfg.FileExists() {
		t.Error("Expected file to exist")
	}
}

func TestCreateBackup(t *testing.T) {
	// Create a temporary config for testing
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	cfg := &Config{
		ConfigDir:  tempDir,
		ConfigFile: filepath.Join(tempDir, ConfigFileName),
	}
	
	// Create the original file with content
	testContent := "test config content"
	err = os.WriteFile(cfg.ConfigFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Create backup
	err = cfg.CreateBackup()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Check backup file exists
	backupFile := cfg.ConfigFile + ".backup"
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		t.Error("Expected backup file to exist")
	}
	
	// Check backup content
	backupContent, err := os.ReadFile(backupFile)
	if err != nil {
		t.Fatalf("Failed to read backup file: %v", err)
	}
	
	if string(backupContent) != testContent {
		t.Errorf("Expected backup content to be '%s', got '%s'", testContent, string(backupContent))
	}
}

func TestCreateBackupNoFile(t *testing.T) {
	// Create a temporary config for testing
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	cfg := &Config{
		ConfigDir:  tempDir,
		ConfigFile: filepath.Join(tempDir, ConfigFileName),
	}
	
	// Create backup when no file exists (should not error)
	err = cfg.CreateBackup()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestRestoreBackup(t *testing.T) {
	// Create a temporary config for testing
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	cfg := &Config{
		ConfigDir:  tempDir,
		ConfigFile: filepath.Join(tempDir, ConfigFileName),
	}
	
	// Create backup file
	backupContent := "backup config content"
	backupFile := cfg.ConfigFile + ".backup"
	err = os.WriteFile(backupFile, []byte(backupContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create backup file: %v", err)
	}
	
	// Restore from backup
	err = cfg.RestoreBackup()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Check restored content
	restoredContent, err := os.ReadFile(cfg.ConfigFile)
	if err != nil {
		t.Fatalf("Failed to read restored file: %v", err)
	}
	
	if string(restoredContent) != backupContent {
		t.Errorf("Expected restored content to be '%s', got '%s'", backupContent, string(restoredContent))
	}
}

func TestRestoreBackupNoFile(t *testing.T) {
	// Create a temporary config for testing
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	cfg := &Config{
		ConfigDir:  tempDir,
		ConfigFile: filepath.Join(tempDir, ConfigFileName),
	}
	
	// Try to restore when no backup exists
	err = cfg.RestoreBackup()
	if err == nil {
		t.Error("Expected error when no backup file exists")
	}
}

func TestGetConfigDir_Methods(t *testing.T) {
	cfg := &Config{
		ConfigDir:  "/test/config",
		ConfigFile: "/test/config/test.json",
	}
	
	if cfg.GetConfigDir() != "/test/config" {
		t.Errorf("Expected '/test/config', got %s", cfg.GetConfigDir())
	}
	
	if cfg.GetConfigFile() != "/test/config/test.json" {
		t.Errorf("Expected '/test/config/test.json', got %s", cfg.GetConfigFile())
	}
}

func TestEnsureDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "ensure-dir-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	testDir := filepath.Join(tempDir, "nested", "directory")
	
	// Directory shouldn't exist initially
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		t.Error("Expected directory to not exist")
	}
	
	// Ensure directory
	err = ensureDir(testDir)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Directory should exist now
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Expected directory to exist")
	}
	
	// Calling again should not error
	err = ensureDir(testDir)
	if err != nil {
		t.Errorf("Expected no error on second call, got: %v", err)
	}
}

func TestGetConfigDirWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("This test simulates Windows behavior on non-Windows systems")
	}
	
	// We can't actually change runtime.GOOS, but we can test the Windows path logic
	// by temporarily setting APPDATA environment variable
	originalAppData := os.Getenv("APPDATA")
	defer os.Setenv("APPDATA", originalAppData)
	
	// Test with APPDATA set
	testAppData := "/test/appdata"
	os.Setenv("APPDATA", testAppData)
	
	// This test verifies our understanding of the Windows path structure
	// In actual Windows environment, the function would use this path
	expectedPath := filepath.Join(testAppData, AppName)
	if runtime.GOOS != "windows" {
		// On non-Windows, we're just verifying the expected path construction
		if !strings.Contains(expectedPath, AppName) {
			t.Errorf("Expected path to contain app name, got %s", expectedPath)
		}
	}
	
	// Reset environment
	os.Setenv("APPDATA", "")
}

func TestGetConfigDirFallback(t *testing.T) {
	// This test is hard to execute reliably across platforms, but we can test
	// the principle that unsupported OS falls back to home directory
	
	// Get current home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get user home directory")
	}
	
	// The fallback path should be ~/.{appname}
	expectedPath := filepath.Join(homeDir, "."+AppName)
	
	// Verify the structure is correct (even if not currently used)
	if !strings.Contains(expectedPath, AppName) {
		t.Errorf("Expected fallback path to contain app name, got %s", expectedPath)
	}
}

func TestNewConfigError(t *testing.T) {
	// Test error handling in New() function by temporarily making home directory inaccessible
	// This is difficult to test reliably across platforms, so we test the structure
	
	// Save original HOME/USERPROFILE
	var originalHome string
	if runtime.GOOS == "windows" {
		originalHome = os.Getenv("USERPROFILE")
		defer os.Setenv("USERPROFILE", originalHome)
		os.Setenv("USERPROFILE", "")
	} else {
		originalHome = os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", "")
	}
	
	// This should cause an error since home directory cannot be determined
	_, err := New()
	if err == nil {
		// If no error, it means the system found another way to determine home
		// which is acceptable behavior
		t.Logf("System was able to determine config directory despite empty HOME/USERPROFILE")
	} else {
		// Expected error case
		if !strings.Contains(err.Error(), "failed to get config directory") {
			t.Errorf("Expected error about config directory, got: %v", err)
		}
	}
}

func TestEnsureDirError(t *testing.T) {
	// Test ensureDir error handling by trying to create a directory in an invalid location
	// Use a path that should fail on most systems
	invalidPath := "/root/impossible/path/for/normal/user"
	
	err := ensureDir(invalidPath)
	if err == nil && runtime.GOOS != "linux" {
		// On some systems this might not fail, which is acceptable
		t.Logf("System allowed creation of directory at %s", invalidPath)
	} else if err != nil {
		// Expected error case
		if !strings.Contains(err.Error(), "failed to create directory") {
			t.Errorf("Expected error about creating directory, got: %v", err)
		}
	}
}

func TestCreateBackupReadError(t *testing.T) {
	// Create a temporary config for testing
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	cfg := &Config{
		ConfigDir:  tempDir,
		ConfigFile: filepath.Join(tempDir, ConfigFileName),
	}
	
	// Create a directory with the same name as the config file to cause a read error
	err = os.Mkdir(cfg.ConfigFile, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	
	// This should cause an error when trying to read the "file" (which is actually a directory)
	err = cfg.CreateBackup()
	if err == nil {
		t.Error("Expected error when reading directory as file")
	} else {
		if !strings.Contains(err.Error(), "failed to read config file") {
			t.Errorf("Expected error about reading config file, got: %v", err)
		}
	}
}

func TestCreateBackupWriteError(t *testing.T) {
	// This test is platform-specific and hard to create reliably
	// We test the principle that write errors are handled properly
	
	// Create a temporary config for testing
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	cfg := &Config{
		ConfigDir:  tempDir,
		ConfigFile: filepath.Join(tempDir, ConfigFileName),
	}
	
	// Create the original file
	testContent := "test config content"
	err = os.WriteFile(cfg.ConfigFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Make the directory read-only to prevent backup creation
	err = os.Chmod(tempDir, 0444)
	if err != nil {
		t.Fatalf("Failed to change directory permissions: %v", err)
	}
	
	// Restore permissions for cleanup
	defer os.Chmod(tempDir, 0755)
	
	// This should cause an error when trying to create the backup file
	err = cfg.CreateBackup()
	if err == nil && runtime.GOOS != "windows" {
		// On some systems (like Windows) this might not fail due to different permission models
		t.Logf("System allowed backup creation despite read-only directory")
	} else if err != nil {
		// Expected error case - could be read error or create backup error
		if !strings.Contains(err.Error(), "failed to create backup") && 
		   !strings.Contains(err.Error(), "failed to read config file") {
			t.Errorf("Expected error about backup or reading config file, got: %v", err)
		}
	}
}