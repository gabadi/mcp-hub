package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	// Test that default config is valid
	if err := config.Validate(); err != nil {
		t.Errorf("Default config should be valid: %v", err)
	}
	
	// Test default values
	if config.UI.Theme != "auto" {
		t.Errorf("Expected default theme 'auto', got '%s'", config.UI.Theme)
	}
	
	if config.UI.MinimumWidth != 40 {
		t.Errorf("Expected default minimum width 40, got %d", config.UI.MinimumWidth)
	}
	
	if !config.UI.WrapNavigation {
		t.Error("Expected default wrap navigation to be true")
	}
	
	if !config.UI.VimKeybindings {
		t.Error("Expected default vim keybindings to be true")
	}
	
	if !config.MCP.AutoDiscovery {
		t.Error("Expected default auto discovery to be true")
	}
	
	if config.MCP.ConnectionTimeout != 5 {
		t.Errorf("Expected default connection timeout 5, got %d", config.MCP.ConnectionTimeout)
	}
	
	if config.Debug.Enabled {
		t.Error("Expected default debug to be false")
	}
	
	if config.Debug.LogLevel != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", config.Debug.LogLevel)
	}
}

func TestConfigValidation(t *testing.T) {
	testCases := []struct {
		name        string
		modifyConfig func(*Config)
		expectError bool
		errorField  string
	}{
		{
			name: "valid config",
			modifyConfig: func(c *Config) {
				// No modifications - should be valid
			},
			expectError: false,
		},
		{
			name: "invalid minimum width",
			modifyConfig: func(c *Config) {
				c.UI.MinimumWidth = 10
			},
			expectError: true,
			errorField:  "ui.minimum_width",
		},
		{
			name: "invalid theme",
			modifyConfig: func(c *Config) {
				c.UI.Theme = "invalid"
			},
			expectError: true,
			errorField:  "ui.theme",
		},
		{
			name: "invalid header height",
			modifyConfig: func(c *Config) {
				c.UI.HeaderHeight = 0
			},
			expectError: true,
			errorField:  "ui.header_height",
		},
		{
			name: "invalid footer height",
			modifyConfig: func(c *Config) {
				c.UI.FooterHeight = 10
			},
			expectError: true,
			errorField:  "ui.footer_height",
		},
		{
			name: "invalid connection timeout",
			modifyConfig: func(c *Config) {
				c.MCP.ConnectionTimeout = 0
			},
			expectError: true,
			errorField:  "mcp.connection_timeout",
		},
		{
			name: "invalid max retries",
			modifyConfig: func(c *Config) {
				c.MCP.MaxRetries = -1
			},
			expectError: true,
			errorField:  "mcp.max_retries",
		},
		{
			name: "invalid refresh interval",
			modifyConfig: func(c *Config) {
				c.MCP.RefreshInterval = 1
			},
			expectError: true,
			errorField:  "mcp.refresh_interval",
		},
		{
			name: "invalid log level",
			modifyConfig: func(c *Config) {
				c.Debug.LogLevel = "invalid"
			},
			expectError: true,
			errorField:  "debug.log_level",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := DefaultConfig()
			tc.modifyConfig(config)
			
			err := config.Validate()
			
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected validation error but got none")
					return
				}
				
				if configErr, ok := err.(ConfigError); ok {
					if configErr.Field != tc.errorField {
						t.Errorf("Expected error field '%s', got '%s'", tc.errorField, configErr.Field)
					}
				} else {
					t.Errorf("Expected ConfigError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error, got: %v", err)
				}
			}
		})
	}
}

func TestConfigSaveLoad(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "mcp-manager-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	configPath := filepath.Join(tempDir, "config.json")
	
	// Create a config with custom values
	originalConfig := DefaultConfig()
	originalConfig.UI.Theme = "dark"
	originalConfig.UI.MinimumWidth = 100
	originalConfig.MCP.ConnectionTimeout = 10
	originalConfig.Debug.Enabled = true
	originalConfig.Debug.LogLevel = "debug"
	
	// Save config
	if err := originalConfig.SaveToPath(configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	// Load config
	loadedConfig, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Compare values
	if loadedConfig.UI.Theme != originalConfig.UI.Theme {
		t.Errorf("Theme mismatch: expected '%s', got '%s'", originalConfig.UI.Theme, loadedConfig.UI.Theme)
	}
	
	if loadedConfig.UI.MinimumWidth != originalConfig.UI.MinimumWidth {
		t.Errorf("MinimumWidth mismatch: expected %d, got %d", originalConfig.UI.MinimumWidth, loadedConfig.UI.MinimumWidth)
	}
	
	if loadedConfig.MCP.ConnectionTimeout != originalConfig.MCP.ConnectionTimeout {
		t.Errorf("ConnectionTimeout mismatch: expected %d, got %d", originalConfig.MCP.ConnectionTimeout, loadedConfig.MCP.ConnectionTimeout)
	}
	
	if loadedConfig.Debug.Enabled != originalConfig.Debug.Enabled {
		t.Errorf("Debug.Enabled mismatch: expected %t, got %t", originalConfig.Debug.Enabled, loadedConfig.Debug.Enabled)
	}
	
	if loadedConfig.Debug.LogLevel != originalConfig.Debug.LogLevel {
		t.Errorf("Debug.LogLevel mismatch: expected '%s', got '%s'", originalConfig.Debug.LogLevel, loadedConfig.Debug.LogLevel)
	}
}

func TestLoadNonExistentConfig(t *testing.T) {
	// Try to load from a non-existent path
	nonExistentPath := "/tmp/non-existent-config.json"
	
	config, err := LoadFromPath(nonExistentPath)
	if err != nil {
		t.Fatalf("LoadFromPath should not error on non-existent file: %v", err)
	}
	
	// Should return default config
	defaultConfig := DefaultConfig()
	if config.UI.Theme != defaultConfig.UI.Theme {
		t.Errorf("Expected default theme, got '%s'", config.UI.Theme)
	}
}

func TestLoadInvalidJSONConfig(t *testing.T) {
	// Create temporary file with invalid JSON
	tempDir, err := os.MkdirTemp("", "mcp-manager-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	configPath := filepath.Join(tempDir, "config.json")
	
	// Write invalid JSON
	invalidJSON := `{"ui": {"theme": "dark",}` // Missing closing brace and trailing comma
	if err := os.WriteFile(configPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}
	
	// Try to load
	_, err = LoadFromPath(configPath)
	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
}

func TestSaveInvalidConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mcp-manager-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	configPath := filepath.Join(tempDir, "config.json")
	
	// Create invalid config
	config := DefaultConfig()
	config.UI.MinimumWidth = 10 // Invalid value
	
	// Try to save
	err = config.SaveToPath(configPath)
	if err == nil {
		t.Error("Expected error when saving invalid config")
	}
}

func TestConfigString(t *testing.T) {
	config := DefaultConfig()
	
	str := config.String()
	if len(str) == 0 {
		t.Error("Config string should not be empty")
	}
	
	// Should be valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(str), &parsed); err != nil {
		t.Errorf("Config string should be valid JSON: %v", err)
	}
}

func TestConfigError(t *testing.T) {
	err := ConfigError{
		Field:   "test.field",
		Value:   "invalid_value",
		Message: "test error message",
	}
	
	expected := "config error for field 'test.field' with value 'invalid_value': test error message"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func BenchmarkConfigLoad(b *testing.B) {
	// Create temporary config file
	tempDir, err := os.MkdirTemp("", "mcp-manager-config-bench")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	configPath := filepath.Join(tempDir, "config.json")
	config := DefaultConfig()
	if err := config.SaveToPath(configPath); err != nil {
		b.Fatalf("Failed to save config: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := LoadFromPath(configPath)
		if err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
	}
}

func BenchmarkConfigValidation(b *testing.B) {
	config := DefaultConfig()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := config.Validate()
		if err != nil {
			b.Fatalf("Config validation failed: %v", err)
		}
	}
}