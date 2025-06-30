package storage

import "cc-mcp-manager/pkg/models"

// InventoryStorage defines the interface for MCP inventory persistence
type InventoryStorage interface {
	// Load loads the inventory from storage
	Load() (*models.MCPInventory, error)
	
	// Save saves the inventory to storage
	Save(inventory *models.MCPInventory) error
	
	// Exists checks if the storage file exists
	Exists() bool
	
	// GetPath returns the storage file path
	GetPath() string
	
	// CreateBackup creates a backup of the current storage
	CreateBackup() error
	
	// RestoreFromBackup restores storage from backup
	RestoreFromBackup() error
}