package services

import (
	"encoding/json"
	"fmt"

	// "log"
	"os"
	"path/filepath"
	"time"

	"cc-mcp-manager/internal/ui/types"
)

// InventoryData wraps MCPItems with metadata for JSON serialization
type InventoryData struct {
	Version   string          `json:"version"`
	Timestamp string          `json:"timestamp"`
	Inventory []types.MCPItem `json:"inventory"`
}

const (
	configFileName = "inventory.json"
	appName        = "cc-mcp-manager"
	configVersion  = "1.0"
)

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	return getConfigPathWithBase("")
}

// getConfigPathWithBase allows overriding the base directory for testing
func getConfigPathWithBase(baseDir string) (string, error) {
	var appConfigDir string

	if baseDir != "" {
		// Use provided base directory (for testing)
		appConfigDir = filepath.Join(baseDir, appName)
	} else {
		// Use system config directory
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user config directory: %w", err)
		}
		appConfigDir = filepath.Join(userConfigDir, appName)
	}

	configPath := filepath.Join(appConfigDir, configFileName)
	// log.Printf("Config file path: %s", configPath)
	return configPath, nil
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	return ensureConfigDirWithBase("")
}

// ensureConfigDirWithBase allows overriding the base directory for testing
func ensureConfigDirWithBase(baseDir string) error {
	var appConfigDir string

	if baseDir != "" {
		// Use provided base directory (for testing)
		appConfigDir = filepath.Join(baseDir, appName)
	} else {
		// Use system config directory
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get user config directory: %w", err)
		}
		appConfigDir = filepath.Join(userConfigDir, appName)
	}

	// Create directory with user read/write/execute permissions only
	err := os.MkdirAll(appConfigDir, 0700)
	if err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", appConfigDir, err)
	}

	// log.Printf("Config directory ensured: %s", appConfigDir)
	return nil
}

// SaveInventory saves the inventory to the JSON config file
func SaveInventory(mcpItems []types.MCPItem) error {
	return saveInventoryWithBase(mcpItems, "")
}

// saveInventoryWithBase allows overriding the base directory for testing
func saveInventoryWithBase(mcpItems []types.MCPItem, baseDir string) error {
	configPath, err := getConfigPathWithBase(baseDir)
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Ensure config directory exists
	if err := ensureConfigDirWithBase(baseDir); err != nil {
		return fmt.Errorf("failed to ensure config directory: %w", err)
	}

	// Create inventory data with metadata
	inventoryData := InventoryData{
		Version:   configVersion,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Inventory: mcpItems,
	}

	// Marshal to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(inventoryData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal inventory data: %w", err)
	}

	// Write to temporary file first for atomic operation
	tempPath := configPath + ".tmp"
	err = os.WriteFile(tempPath, jsonData, 0600) // User read/write only
	if err != nil {
		return fmt.Errorf("failed to write temporary config file %s: %w", tempPath, err)
	}

	// Atomic rename
	err = os.Rename(tempPath, configPath)
	if err != nil {
		// Clean up temporary file on failure
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename temporary config file: %w", err)
	}

	// log.Printf("Inventory saved successfully to: %s", configPath)
	return nil
}

// LoadInventory loads the inventory from the JSON config file
func LoadInventory() ([]types.MCPItem, error) {
	return loadInventoryWithBase("")
}

// loadInventoryWithBase allows overriding the base directory for testing
func loadInventoryWithBase(baseDir string) ([]types.MCPItem, error) {
	configPath, err := getConfigPathWithBase(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// log.Printf("Config file does not exist: %s, will start with empty inventory", configPath)
		return []types.MCPItem{}, nil
	}

	// Read config file
	jsonData, err := os.ReadFile(configPath)
	if err != nil {
		// log.Printf("Failed to read config file %s: %v, falling back to empty inventory", configPath, err)
		return []types.MCPItem{}, nil
	}

	// Try to unmarshal inventory data
	var inventoryData InventoryData
	err = json.Unmarshal(jsonData, &inventoryData)
	if err != nil {
		// Handle corrupted file - backup and start fresh
		backupPath := configPath + ".corrupted." + time.Now().Format("20060102-150405")
		if backupErr := os.Rename(configPath, backupPath); backupErr != nil {
			// log.Printf("Failed to backup corrupted config file: %v", backupErr)
		} else {
			// log.Printf("Corrupted config file backed up to: %s", backupPath)
		}

		// log.Printf("Failed to parse config file %s: %v, falling back to empty inventory", configPath, err)
		return []types.MCPItem{}, nil
	}

	// log.Printf("Inventory loaded successfully from: %s (version: %s, %d items)",
	//	configPath, inventoryData.Version, len(inventoryData.Inventory))

	return inventoryData.Inventory, nil
}

// SaveModelInventory saves the inventory from a model
func SaveModelInventory(model types.Model) error {
	return SaveInventory(model.MCPItems)
}
