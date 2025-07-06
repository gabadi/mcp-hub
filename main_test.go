package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"mcp-hub/internal/platform"
	"mcp-hub/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// Test main function components 
func TestMainFunction(t *testing.T) {
	// Test that main function doesn't panic when importing packages
	t.Run("package_imports", func(t *testing.T) {
		// Test that we can create a platform service
		factory := platform.NewPlatformServiceFactoryDefault()
		if factory == nil {
			t.Error("Expected platform service factory to be created")
		}
		
		platformService := factory.CreatePlatformService()
		if platformService == nil {
			t.Error("Expected platform service to be created")
		}
		
		// Test that we can create a UI model
		model := ui.NewModel()
		// Test that model has valid fields
		if model.PlatformService == nil {
			t.Error("Expected UI model to have platform service")
		}
		
		// Test that we can create a tea program
		program := tea.NewProgram(model, tea.WithAltScreen())
		if program == nil {
			t.Error("Expected tea program to be created")
		}
	})
}

// Test runApp function components 
func TestRunApp(t *testing.T) {
	t.Run("log_path_creation", func(t *testing.T) {
		// Test that log path creation works
		platformService := platform.NewPlatformServiceFactoryDefault().CreatePlatformService()
		
		logPath := filepath.Join(platformService.GetLogPath(), "mcp-hub.log")
		logDir := filepath.Dir(logPath)
		
		// Create log directory
		err := os.MkdirAll(logDir, platformService.GetDefaultDirectoryPermissions())
		if err != nil {
			t.Errorf("Failed to create log directory: %v", err)
		}
		
		// Test log file creation
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, platformService.GetDefaultFilePermissions())
		if err != nil {
			t.Errorf("Failed to create log file: %v", err)
		}
		defer func() { _ = file.Close() }()
		
		// Test writing to log file
		originalOutput := log.Writer()
		defer log.SetOutput(originalOutput)
		
		log.SetOutput(file)
		log.Printf("Test log message")
		
		// Verify log file has content
		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Errorf("Failed to read log file: %v", err)
		}
		
		if !strings.Contains(string(content), "Test log message") {
			t.Error("Expected log message to be written to file")
		}
	})
	
	t.Run("platform_service_creation", func(t *testing.T) {
		// Test platform service creation
		factory := platform.NewPlatformServiceFactoryDefault()
		service := factory.CreatePlatformService()
		
		// Test that platform service provides valid paths
		logPath := service.GetLogPath()
		if logPath == "" {
			t.Error("Expected non-empty log path")
		}
		
		configPath := service.GetConfigPath()
		if configPath == "" {
			t.Error("Expected non-empty config path")
		}
		
		permissions := service.GetDefaultFilePermissions()
		if permissions == 0 {
			t.Error("Expected non-zero file permissions")
		}
		
		dirPermissions := service.GetDefaultDirectoryPermissions()
		if dirPermissions == 0 {
			t.Error("Expected non-zero directory permissions")
		}
	})
	
	t.Run("ui_model_creation", func(t *testing.T) {
		// Test UI model creation
		model := ui.NewModel()
		// Test that model has valid fields
		if model.PlatformService == nil {
			t.Error("Expected UI model to have platform service")
		}
		
		// Test that model has initialization command
		cmd := model.Init()
		if cmd == nil {
			t.Error("Expected initialization command from model")
		}
	})
	
	t.Run("tea_program_creation", func(t *testing.T) {
		// Test tea program creation
		model := ui.NewModel()
		program := tea.NewProgram(model, tea.WithAltScreen())
		
		if program == nil {
			t.Error("Expected tea program to be created")
		}
	})
}

// Test error handling scenarios
func TestErrorHandling(t *testing.T) {
	t.Run("invalid_log_directory", func(t *testing.T) {
		// Test behavior when log directory cannot be created
		invalidPath := "/invalid/path/that/does/not/exist"
		
		// This should not panic even if directory creation fails
		err := os.MkdirAll(invalidPath, 0755)
		if err == nil {
			t.Error("Expected error when creating invalid directory")
		}
	})
	
	t.Run("log_file_permissions", func(t *testing.T) {
		// Test log file with specific permissions
		tempDir := t.TempDir()
		logPath := filepath.Join(tempDir, "test.log")
		
		// Create log file with specific permissions
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			t.Errorf("Failed to create log file with permissions: %v", err)
		}
		defer func() { _ = file.Close() }()
		
		// Check file permissions
		info, err := file.Stat()
		if err != nil {
			t.Errorf("Failed to get file info: %v", err)
		}
		
		mode := info.Mode()
		if mode.Perm() != 0600 {
			t.Errorf("Expected file permissions 0600, got %o", mode.Perm())
		}
	})
}

// Test application initialization components
func TestAppInitialization(t *testing.T) {
	t.Run("platform_detection", func(t *testing.T) {
		// Test that platform is detected correctly
		factory := platform.NewPlatformServiceFactoryDefault()
		service := factory.CreatePlatformService()
		
		if service == nil {
			t.Error("Expected platform service to be created")
		}
		
		// Test that platform service provides platform info
		platformName := service.GetPlatformName()
		if platformName == "" {
			t.Error("Expected valid platform name")
		}
	})
	
	t.Run("supported_platforms", func(t *testing.T) {
		// Test that we have supported platforms
		factory := platform.NewPlatformServiceFactoryDefault()
		service := factory.CreatePlatformService()
		
		if service == nil {
			t.Error("Expected platform service to be created")
		}
		
		// Test that platform service provides valid platform info
		platformType := service.GetPlatform()
		if platformType == 0 {
			t.Error("Expected valid platform type")
		}
	})
}

// Test component integration
func TestComponentIntegration(t *testing.T) {
	t.Run("platform_ui_integration", func(t *testing.T) {
		// Test that platform service and UI model work together
		platformService := platform.NewPlatformServiceFactoryDefault().CreatePlatformService()
		model := ui.NewModel()
		
		// Both should be created successfully
		if platformService == nil {
			t.Error("Expected platform service to be created")
		}
		if model.PlatformService == nil {
			t.Error("Expected UI model to have platform service")
		}
		
		// Test that we can get paths from platform service
		logPath := platformService.GetLogPath()
		configPath := platformService.GetConfigPath()
		
		if logPath == "" || configPath == "" {
			t.Error("Expected valid paths from platform service")
		}
	})
}

// Test resource management
func TestResourceManagement(t *testing.T) {
	t.Run("file_cleanup", func(t *testing.T) {
		// Test that files are properly closed
		tempDir := t.TempDir()
		logPath := filepath.Join(tempDir, "test.log")
		
		// Create and close file
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			t.Errorf("Failed to create log file: %v", err)
		}
		
		// Write something to file
		_, err = file.WriteString("test content")
		if err != nil {
			t.Errorf("Failed to write to file: %v", err)
		}
		
		// Close file
		err = file.Close()
		if err != nil {
			t.Errorf("Failed to close file: %v", err)
		}
		
		// Verify file exists and has content
		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Errorf("Failed to read file: %v", err)
		}
		
		if string(content) != "test content" {
			t.Errorf("Expected 'test content', got '%s'", string(content))
		}
	})
}

// Test command line argument handling
func TestCommandLineArgs(t *testing.T) {
	t.Run("no_args", func(t *testing.T) {
		// Test that application doesn't crash with no arguments
		originalArgs := os.Args
		defer func() { os.Args = originalArgs }()
		
		os.Args = []string{"mcp-hub"}
		
		// This should not panic
		// We can't easily test main() directly, but we can verify components work
		factory := platform.NewPlatformServiceFactoryDefault()
		if factory == nil {
			t.Error("Expected platform service factory to work with no args")
		}
	})
	
	t.Run("extra_args", func(t *testing.T) {
		// Test that extra arguments are ignored gracefully
		originalArgs := os.Args
		defer func() { os.Args = originalArgs }()
		
		os.Args = []string{"mcp-hub", "--help", "--version", "extra"}
		
		// This should not panic
		// Current implementation ignores args
		factory := platform.NewPlatformServiceFactoryDefault()
		if factory == nil {
			t.Error("Expected platform service factory to work with extra args")
		}
	})
}

// Test main() function execution with process isolation
func TestMainExecution(t *testing.T) {
	// Skip this test if we can't run the executable
	if testing.Short() {
		t.Skip("Skipping main execution test in short mode")
	}
	
	t.Run("main_process_execution", func(t *testing.T) {
		// Build test binary
		testBinary := filepath.Join(t.TempDir(), "mcp-hub-test")
		cmd := exec.Command("go", "build", "-o", testBinary, ".")
		cmd.Dir = "."
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to build test binary: %v", err)
		}
		
		// Test main function with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		
		cmd = exec.CommandContext(ctx, testBinary)
		// The application should start and immediately exit due to TUI not being able to run in test
		err := cmd.Run()
		
		// We expect the program to exit quickly with an error (can't run TUI in test)
		if err == nil {
			t.Error("Expected TUI to fail in test environment")
		}
	})
}

// Test runApp function error scenarios
func TestRunAppErrorScenarios(t *testing.T) {
	t.Run("invalid_log_directory_graceful_handling", func(t *testing.T) {
		// Create a mock platform service that returns invalid log path
		mockPlatform := platform.GetMockPlatformService()
		
		// Test that runApp doesn't crash even with invalid log directory
		// We can't easily test runApp directly, but we can test log setup
		logPath := filepath.Join("/invalid/nonexistent/path", "mcp-hub.log")
		logDir := filepath.Dir(logPath)
		
		// This should fail gracefully
		err := os.MkdirAll(logDir, mockPlatform.GetDefaultDirectoryPermissions())
		if err == nil {
			t.Error("Expected error when creating invalid directory")
		}
		
		// Application should continue even if log setup fails
		// Test that UI components still work
		model := ui.NewModel()
		if model.PlatformService == nil {
			t.Error("Expected UI model to be created even with log setup failure")
		}
	})
	
	t.Run("readonly_log_directory", func(t *testing.T) {
		// Create a read-only directory
		tempDir := t.TempDir()
		readOnlyDir := filepath.Join(tempDir, "readonly")
		
		// Create directory
		err := os.MkdirAll(readOnlyDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
		
		// Make it read-only
		err = os.Chmod(readOnlyDir, 0555)
		if err != nil {
			t.Fatalf("Failed to make directory read-only: %v", err)
		}
		
		// Try to create log file - should fail
		logPath := filepath.Join(readOnlyDir, "mcp-hub.log")
		_, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err == nil {
			t.Error("Expected error when creating file in read-only directory")
		}
		
		// Application should still work
		factory := platform.NewPlatformServiceFactoryDefault()
		if factory == nil {
			t.Error("Expected platform service factory to work even with log failure")
		}
	})
	
	t.Run("tea_program_error_handling", func(t *testing.T) {
		// Test that Tea program creation is robust
		model := ui.NewModel()
		if model.PlatformService == nil {
			t.Fatal("Failed to create UI model")
		}
		
		// Test program creation with different options
		program := tea.NewProgram(model, tea.WithAltScreen())
		if program == nil {
			t.Error("Expected tea program to be created")
		}
		
		// Test without alt screen
		program2 := tea.NewProgram(model)
		if program2 == nil {
			t.Error("Expected tea program to be created without alt screen")
		}
	})
}

// Test application startup sequence
func TestApplicationStartup(t *testing.T) {
	t.Run("complete_startup_sequence", func(t *testing.T) {
		// Test the complete startup sequence
		// 1. Platform service creation
		platformService := platform.NewPlatformServiceFactoryDefault().CreatePlatformService()
		if platformService == nil {
			t.Fatal("Failed to create platform service")
		}
		
		// 2. Log path setup
		logPath := filepath.Join(platformService.GetLogPath(), "mcp-hub.log")
		logDir := filepath.Dir(logPath)
		
		// 3. UI model creation
		model := ui.NewModel()
		if model.PlatformService == nil {
			t.Fatal("Failed to create UI model")
		}
		
		// 4. Tea program creation
		program := tea.NewProgram(model, tea.WithAltScreen())
		if program == nil {
			t.Fatal("Failed to create tea program")
		}
		
		// Verify all components are properly initialized
		if model.PlatformService == nil {
			t.Error("UI model should have platform service")
		}
		
		// Test that paths are valid
		if logPath == "" {
			t.Error("Log path should not be empty")
		}
		
		if logDir == "" {
			t.Error("Log directory should not be empty")
		}
	})
	
	t.Run("platform_specific_paths", func(t *testing.T) {
		// Test that platform-specific paths are correctly set
		platformService := platform.NewPlatformServiceFactoryDefault().CreatePlatformService()
		
		// Test log path
		logPath := platformService.GetLogPath()
		if logPath == "" {
			t.Error("Log path should not be empty")
		}
		
		// Test config path
		configPath := platformService.GetConfigPath()
		if configPath == "" {
			t.Error("Config path should not be empty")
		}
		
		// Test permissions
		filePerms := platformService.GetDefaultFilePermissions()
		if filePerms == 0 {
			t.Error("File permissions should not be zero")
		}
		
		dirPerms := platformService.GetDefaultDirectoryPermissions()
		if dirPerms == 0 {
			t.Error("Directory permissions should not be zero")
		}
	})
}

// Test log file management
func TestLogFileManagement(t *testing.T) {
	t.Run("log_file_creation_and_cleanup", func(t *testing.T) {
		// Create temporary directory for log testing
		tempDir := t.TempDir()
		logPath := filepath.Join(tempDir, "test.log")
		
		// Test log file creation
		logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			t.Fatalf("Failed to create log file: %v", err)
		}
		
		// Test log output redirection
		originalOutput := log.Writer()
		defer log.SetOutput(originalOutput)
		
		log.SetOutput(logFile)
		log.Printf("Test log message at %s", time.Now().Format(time.RFC3339))
		
		// Test cleanup
		err = logFile.Close()
		if err != nil {
			t.Errorf("Failed to close log file: %v", err)
		}
		
		// Verify log file contains expected content
		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Errorf("Failed to read log file: %v", err)
		}
		
		if !strings.Contains(string(content), "Test log message") {
			t.Error("Log file should contain test message")
		}
	})
	
	t.Run("log_file_append_mode", func(t *testing.T) {
		// Test that log file opens in append mode
		tempDir := t.TempDir()
		logPath := filepath.Join(tempDir, "append-test.log")
		
		// Create initial content
		err := os.WriteFile(logPath, []byte("Initial content\n"), 0600)
		if err != nil {
			t.Fatalf("Failed to create initial log file: %v", err)
		}
		
		// Open in append mode
		logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			t.Fatalf("Failed to open log file in append mode: %v", err)
		}
		defer func() { _ = logFile.Close() }()
		
		// Write additional content
		_, err = logFile.WriteString("Appended content\n")
		if err != nil {
			t.Fatalf("Failed to write to log file: %v", err)
		}
		
		// Verify both contents are present
		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}
		
		contentStr := string(content)
		if !strings.Contains(contentStr, "Initial content") {
			t.Error("Log file should contain initial content")
		}
		if !strings.Contains(contentStr, "Appended content") {
			t.Error("Log file should contain appended content")
		}
	})
}