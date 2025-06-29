// Package config provides configuration management for the MCP Manager CLI.
// It handles loading, validation, and default configuration values.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	// UI Configuration
	UI UIConfig `json:"ui"`
	
	// MCP Configuration  
	MCP MCPConfig `json:"mcp"`
	
	// Debug Configuration
	Debug DebugConfig `json:"debug"`
}

// UIConfig contains user interface settings
type UIConfig struct {
	// Theme settings
	Theme               string `json:"theme"`                // "light", "dark", "auto"
	
	// Layout settings
	MinimumWidth        int    `json:"minimum_width"`        // Minimum terminal width
	HeaderHeight        int    `json:"header_height"`        // Header height in lines
	FooterHeight        int    `json:"footer_height"`        // Footer height in lines
	
	// Navigation settings
	WrapNavigation      bool   `json:"wrap_navigation"`      // Enable wrap-around navigation
	VimKeybindings      bool   `json:"vim_keybindings"`      // Enable vim-style key bindings
	
	// Search settings
	SearchCaseSensitive bool   `json:"search_case_sensitive"` // Case sensitive search
	SearchRegex         bool   `json:"search_regex"`          // Enable regex search
}

// MCPConfig contains MCP-specific settings
type MCPConfig struct {
	// Server discovery settings
	AutoDiscovery       bool     `json:"auto_discovery"`        // Auto-discover MCP servers
	ConfigPaths         []string `json:"config_paths"`          // Paths to search for MCP configs
	
	// Connection settings
	ConnectionTimeout   int      `json:"connection_timeout"`    // Connection timeout in seconds
	MaxRetries          int      `json:"max_retries"`           // Maximum connection retries
	
	// Refresh settings
	AutoRefresh         bool     `json:"auto_refresh"`          // Auto-refresh server status
	RefreshInterval     int      `json:"refresh_interval"`      // Refresh interval in seconds
}

// DebugConfig contains debugging and logging settings
type DebugConfig struct {
	Enabled             bool   `json:"enabled"`               // Enable debug mode
	LogLevel            string `json:"log_level"`             // Log level (debug, info, warn, error)
	LogFile             string `json:"log_file"`              // Log file path
	ShowPerformanceInfo bool   `json:"show_performance_info"` // Show performance metrics
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		UI: UIConfig{
			Theme:               "auto",
			MinimumWidth:        40,
			HeaderHeight:        3,
			FooterHeight:        2,
			WrapNavigation:      true,
			VimKeybindings:      true,
			SearchCaseSensitive: false,
			SearchRegex:         false,
		},
		MCP: MCPConfig{
			AutoDiscovery:     true,
			ConfigPaths:       []string{"~/.config/mcp", "~/.mcp", "./mcp-config"},
			ConnectionTimeout: 5,
			MaxRetries:        3,
			AutoRefresh:       false,
			RefreshInterval:   30,
		},
		Debug: DebugConfig{
			Enabled:             false,
			LogLevel:            "info",
			LogFile:             "",
			ShowPerformanceInfo: false,
		},
	}
}

// ConfigError represents a configuration error
type ConfigError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("config error for field '%s' with value '%v': %s", e.Field, e.Value, e.Message)
}

// Validate validates the configuration and returns any errors
func (c *Config) Validate() error {
	// Validate UI config
	if c.UI.MinimumWidth < 20 {
		return ConfigError{
			Field:   "ui.minimum_width",
			Value:   c.UI.MinimumWidth,
			Message: "must be at least 20 characters",
		}
	}
	
	if c.UI.Theme != "light" && c.UI.Theme != "dark" && c.UI.Theme != "auto" {
		return ConfigError{
			Field:   "ui.theme",
			Value:   c.UI.Theme,
			Message: "must be 'light', 'dark', or 'auto'",
		}
	}
	
	if c.UI.HeaderHeight < 1 || c.UI.HeaderHeight > 10 {
		return ConfigError{
			Field:   "ui.header_height",
			Value:   c.UI.HeaderHeight,
			Message: "must be between 1 and 10",
		}
	}
	
	if c.UI.FooterHeight < 1 || c.UI.FooterHeight > 5 {
		return ConfigError{
			Field:   "ui.footer_height",
			Value:   c.UI.FooterHeight,
			Message: "must be between 1 and 5",
		}
	}
	
	// Validate MCP config
	if c.MCP.ConnectionTimeout < 1 || c.MCP.ConnectionTimeout > 300 {
		return ConfigError{
			Field:   "mcp.connection_timeout",
			Value:   c.MCP.ConnectionTimeout,
			Message: "must be between 1 and 300 seconds",
		}
	}
	
	if c.MCP.MaxRetries < 0 || c.MCP.MaxRetries > 10 {
		return ConfigError{
			Field:   "mcp.max_retries",
			Value:   c.MCP.MaxRetries,
			Message: "must be between 0 and 10",
		}
	}
	
	if c.MCP.RefreshInterval < 5 || c.MCP.RefreshInterval > 3600 {
		return ConfigError{
			Field:   "mcp.refresh_interval",
			Value:   c.MCP.RefreshInterval,
			Message: "must be between 5 and 3600 seconds",
		}
	}
	
	// Validate debug config
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	
	if !validLogLevels[c.Debug.LogLevel] {
		return ConfigError{
			Field:   "debug.log_level",
			Value:   c.Debug.LogLevel,
			Message: "must be one of: debug, info, warn, error",
		}
	}
	
	return nil
}

// getConfigPath returns the default configuration file path
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, ".config", "mcp-manager")
	configFile := filepath.Join(configDir, "config.json")
	
	return configFile, nil
}

// ensureConfigDir creates the configuration directory if it doesn't exist
func ensureConfigDir(configPath string) error {
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return nil
}

// Load loads configuration from the default location or returns default config
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to determine config path: %w", err)
	}
	
	return LoadFromPath(configPath)
}

// LoadFromPath loads configuration from a specific file path
func LoadFromPath(path string) (*Config, error) {
	// If file doesn't exist, return default config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		config := DefaultConfig()
		if validationErr := config.Validate(); validationErr != nil {
			return nil, fmt.Errorf("default config validation failed: %w", validationErr)
		}
		return config, nil
	}
	
	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w", path, err)
	}
	
	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file '%s': %w", path, err)
	}
	
	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return &config, nil
}

// Save saves the configuration to the default location
func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to determine config path: %w", err)
	}
	
	return c.SaveToPath(configPath)
}

// SaveToPath saves the configuration to a specific file path
func (c *Config) SaveToPath(path string) error {
	// Validate before saving
	if err := c.Validate(); err != nil {
		return fmt.Errorf("cannot save invalid config: %w", err)
	}
	
	// Ensure config directory exists
	if err := ensureConfigDir(path); err != nil {
		return err
	}
	
	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file '%s': %w", path, err)
	}
	
	return nil
}

// String returns a string representation of the configuration
func (c *Config) String() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("Config{error: %v}", err)
	}
	return string(data)
}