package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	// ConfigFileName is the name of the config file
	ConfigFileName = "mcp-inventory.json"
	// AppName is the application name used for config directory
	AppName = "mcp-manager"
)

// Config holds configuration for the MCP Manager
type Config struct {
	ConfigDir  string
	ConfigFile string
}

// New creates a new configuration with appropriate directory paths
func New() (*Config, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}
	
	configFile := filepath.Join(configDir, ConfigFileName)
	
	// Ensure config directory exists
	if err := ensureDir(configDir); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	
	config := &Config{
		ConfigDir:  configDir,
		ConfigFile: configFile,
	}
	
	// Log the config file location for user reference (AC5)
	log.Printf("MCP Manager config file location: %s", configFile)
	
	return config, nil
}

// getConfigDir returns the appropriate config directory based on OS
func getConfigDir() (string, error) {
	var configDir string
	
	switch runtime.GOOS {
	case "darwin": // macOS
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, "Library", "Application Support", AppName)
		
	case "linux":
		// Follow XDG Base Directory specification
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			configDir = filepath.Join(xdgConfigHome, AppName)
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configDir = filepath.Join(homeDir, ".config", AppName)
		}
		
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			configDir = filepath.Join(appData, AppName)
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configDir = filepath.Join(homeDir, "AppData", "Roaming", AppName)
		}
		
	default:
		// Fallback for unsupported OS
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, "."+AppName)
	}
	
	return configDir, nil
}

// ensureDir creates a directory if it doesn't exist
func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		log.Printf("Created config directory: %s", dir)
	}
	return nil
}

// FileExists checks if the config file exists
func (c *Config) FileExists() bool {
	_, err := os.Stat(c.ConfigFile)
	return !os.IsNotExist(err)
}

// GetConfigDir returns the config directory path
func (c *Config) GetConfigDir() string {
	return c.ConfigDir
}

// GetConfigFile returns the full path to the config file
func (c *Config) GetConfigFile() string {
	return c.ConfigFile
}

// CreateBackup creates a backup of the config file
func (c *Config) CreateBackup() error {
	if !c.FileExists() {
		return nil // No file to backup
	}
	
	backupFile := c.ConfigFile + ".backup"
	
	input, err := os.ReadFile(c.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	if err := os.WriteFile(backupFile, input, 0644); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	
	log.Printf("Created backup of config file: %s", backupFile)
	return nil
}

// RestoreBackup restores the config file from backup
func (c *Config) RestoreBackup() error {
	backupFile := c.ConfigFile + ".backup"
	
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", backupFile)
	}
	
	input, err := os.ReadFile(backupFile)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}
	
	if err := os.WriteFile(c.ConfigFile, input, 0644); err != nil {
		return fmt.Errorf("failed to restore from backup: %w", err)
	}
	
	log.Printf("Restored config file from backup: %s", backupFile)
	return nil
}