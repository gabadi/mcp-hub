package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"cc-mcp-manager/pkg/config"
	"cc-mcp-manager/pkg/metrics"
	"cc-mcp-manager/pkg/models"
)

// JSONStorage implements InventoryStorage using JSON files
type JSONStorage struct {
	config   *config.Config
	fileLock *FileLock
}

// NewJSONStorage creates a new JSON storage instance
func NewJSONStorage(cfg *config.Config) *JSONStorage {
	return &JSONStorage{
		config:   cfg,
		fileLock: NewFileLock(30 * time.Second),
	}
}

// Load loads the inventory from the JSON file
func (s *JSONStorage) Load() (*models.MCPInventory, error) {
	timer := metrics.NewTimer()
	log.Printf("Loading MCP inventory from: %s", s.config.GetConfigFile())
	
	// Acquire file lock
	ctx := context.Background()
	lockHandle, err := s.fileLock.Lock(ctx, s.config.GetConfigFile())
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock for loading: %w", err)
	}
	defer lockHandle.Unlock()
	
	// Check if file exists (AC4: graceful handling of missing files)
	if !s.Exists() {
		log.Printf("Config file does not exist, creating new inventory")
		return models.NewInventory(), nil
	}
	
	// Read the file
	data, err := os.ReadFile(s.config.GetConfigFile())
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Handle empty file
	if len(data) == 0 {
		log.Printf("Config file is empty, creating new inventory")
		return models.NewInventory(), nil
	}
	
	// Try to parse JSON (AC4: graceful handling of corrupted files)
	inventory, err := s.parseJSONWithRecovery(data)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory: %w", err)
	}
	
	log.Printf("Successfully loaded inventory with %d MCPs", len(inventory.MCPs))
	
	// Record performance metrics
	metrics.RecordLoadOperation(timer.Stop(), len(inventory.MCPs))
	
	return inventory, nil
}

// parseJSONWithRecovery attempts to parse JSON with error recovery
func (s *JSONStorage) parseJSONWithRecovery(data []byte) (*models.MCPInventory, error) {
	// First attempt: try to parse normally
	inventory, err := models.FromJSON(data)
	if err == nil {
		return inventory, nil
	}
	
	log.Printf("Failed to parse config file, attempting recovery: %v", err)
	
	// Create backup before attempting recovery
	if backupErr := s.CreateBackup(); backupErr != nil {
		log.Printf("Warning: Failed to create backup before recovery: %v", backupErr)
	} else {
		metrics.RecordBackupOperation()
	}
	
	// Record recovery operation
	metrics.RecordRecoveryOperation()
	
	// Second attempt: try to parse as generic JSON to check structure
	var genericData interface{}
	if jsonErr := json.Unmarshal(data, &genericData); jsonErr != nil {
		// File is not valid JSON at all
		log.Printf("File is not valid JSON, creating new inventory")
		return models.NewInventory(), nil
	}
	
	// Third attempt: try to extract what we can
	var partialInventory struct {
		Version string        `json:"version"`
		MCPs    []interface{} `json:"mcps"`
	}
	
	if err := json.Unmarshal(data, &partialInventory); err == nil {
		log.Printf("Recovered partial inventory structure, creating new inventory with version info")
		inventory := models.NewInventory()
		if partialInventory.Version != "" {
			inventory.Version = partialInventory.Version
		}
		return inventory, nil
	}
	
	// If all recovery attempts fail, create a new inventory
	log.Printf("All recovery attempts failed, creating new inventory")
	return models.NewInventory(), nil
}

// Save saves the inventory to the JSON file
func (s *JSONStorage) Save(inventory *models.MCPInventory) error {
	timer := metrics.NewTimer()
	log.Printf("Saving MCP inventory to: %s", s.config.GetConfigFile())
	
	// Acquire file lock
	ctx := context.Background()
	lockHandle, err := s.fileLock.Lock(ctx, s.config.GetConfigFile())
	if err != nil {
		return fmt.Errorf("failed to acquire lock for saving: %w", err)
	}
	defer lockHandle.Unlock()
	
	// Validate inventory before saving
	validationTimer := metrics.NewTimer()
	if err := inventory.Validate(); err != nil {
		return fmt.Errorf("invalid inventory: %w", err)
	}
	metrics.RecordValidation(validationTimer.Stop())
	
	// Update inventory metadata
	inventory.UpdatedAt = time.Now()
	inventory.Metadata.FileCount = len(inventory.MCPs)
	inventory.Metadata.LastSync = time.Now()
	
	// Convert to JSON
	data, err := inventory.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal inventory: %w", err)
	}
	
	// Atomic write: write to temporary file first, then rename
	tempFile := s.config.GetConfigFile() + ".tmp"
	
	if err := s.writeFile(tempFile, data); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}
	
	// Atomic rename
	if err := os.Rename(tempFile, s.config.GetConfigFile()); err != nil {
		// Clean up temp file on failure
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}
	
	log.Printf("Successfully saved inventory with %d MCPs", len(inventory.MCPs))
	
	// Record performance metrics
	metrics.RecordSaveOperation(timer.Stop(), len(inventory.MCPs))
	
	return nil
}

// writeFile writes data to a file with proper permissions
func (s *JSONStorage) writeFile(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	
	// Ensure data is written to disk
	return file.Sync()
}

// Exists checks if the storage file exists
func (s *JSONStorage) Exists() bool {
	return s.config.FileExists()
}

// GetPath returns the storage file path
func (s *JSONStorage) GetPath() string {
	return s.config.GetConfigFile()
}

// CreateBackup creates a backup of the current storage file
func (s *JSONStorage) CreateBackup() error {
	if !s.Exists() {
		return nil // No file to backup
	}
	
	if err := s.config.CreateBackup(); err != nil {
		return err
	}
	
	// Record backup operation
	metrics.RecordBackupOperation()
	return nil
}

// RestoreFromBackup restores storage from backup
func (s *JSONStorage) RestoreFromBackup() error {
	return s.config.RestoreBackup()
}

// MigrateInventory performs schema migration if needed
func (s *JSONStorage) MigrateInventory(inventory *models.MCPInventory) (*models.MCPInventory, error) {
	// For now, we only support version 1.0
	// Future versions can implement migration logic here
	
	if inventory.Version == "" {
		log.Printf("Migrating inventory from unknown version to 1.0")
		inventory.Version = "1.0"
		return inventory, nil
	}
	
	if inventory.Version != "1.0" {
		return nil, fmt.Errorf("unsupported inventory version: %s", inventory.Version)
	}
	
	return inventory, nil
}

// ValidateAndRepair validates the inventory and attempts to repair common issues
func (s *JSONStorage) ValidateAndRepair(inventory *models.MCPInventory) (*models.MCPInventory, error) {
	log.Printf("Validating and repairing inventory")
	
	// Migrate if needed
	migratedInventory, err := s.MigrateInventory(inventory)
	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	
	// Fix missing timestamps
	now := time.Now()
	for i := range migratedInventory.MCPs {
		mcp := &migratedInventory.MCPs[i]
		
		if mcp.CreatedAt.IsZero() {
			mcp.CreatedAt = now
		}
		if mcp.UpdatedAt.IsZero() {
			mcp.UpdatedAt = now
		}
		
		// Generate ID if missing
		if mcp.ID == "" {
			mcp.ID = fmt.Sprintf("mcp_%d", i+1)
			log.Printf("Generated missing ID for MCP: %s", mcp.ID)
		}
	}
	
	// Fix inventory metadata
	if migratedInventory.Metadata.Created.IsZero() {
		migratedInventory.Metadata.Created = now
	}
	migratedInventory.Metadata.FileCount = len(migratedInventory.MCPs)
	
	if migratedInventory.UpdatedAt.IsZero() {
		migratedInventory.UpdatedAt = now
	}
	
	// Validate the repaired inventory
	if err := migratedInventory.Validate(); err != nil {
		return nil, fmt.Errorf("inventory validation failed after repair: %w", err)
	}
	
	log.Printf("Successfully validated and repaired inventory")
	return migratedInventory, nil
}

// LoadWithRecovery loads the inventory with enhanced error recovery
func (s *JSONStorage) LoadWithRecovery() (*models.MCPInventory, error) {
	inventory, err := s.Load()
	if err != nil {
		return nil, err
	}
	
	// Validate and repair the loaded inventory
	return s.ValidateAndRepair(inventory)
}

// GetPerformanceMetrics returns current performance metrics
func (s *JSONStorage) GetPerformanceMetrics() map[string]int64 {
	return metrics.GetSummary()
}

// LoadWithLock loads the inventory with explicit file locking
func (s *JSONStorage) LoadWithLock(ctx context.Context) (*models.MCPInventory, error) {
	return WithFileLockOperation(ctx, s.config.GetConfigFile(), func() (*models.MCPInventory, error) {
		return s.loadUnsafe()
	})
}

// SaveWithLock saves the inventory with explicit file locking
func (s *JSONStorage) SaveWithLock(ctx context.Context, inventory *models.MCPInventory) error {
	_, err := WithFileLockOperation(ctx, s.config.GetConfigFile(), func() (*models.MCPInventory, error) {
		err := s.saveUnsafe(inventory)
		return nil, err
	})
	return err
}

// loadUnsafe performs the actual load operation without locking (for internal use)
func (s *JSONStorage) loadUnsafe() (*models.MCPInventory, error) {
	timer := metrics.NewTimer()
	
	// Check if file exists (AC4: graceful handling of missing files)
	if !s.Exists() {
		log.Printf("Config file does not exist, creating new inventory")
		inventory := models.NewInventory()
		metrics.RecordLoadOperation(timer.Stop(), len(inventory.MCPs))
		return inventory, nil
	}
	
	// Read the file
	data, err := os.ReadFile(s.config.GetConfigFile())
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Handle empty file
	if len(data) == 0 {
		log.Printf("Config file is empty, creating new inventory")
		inventory := models.NewInventory()
		metrics.RecordLoadOperation(timer.Stop(), len(inventory.MCPs))
		return inventory, nil
	}
	
	// Try to parse JSON (AC4: graceful handling of corrupted files)
	inventory, err := s.parseJSONWithRecovery(data)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory: %w", err)
	}
	
	log.Printf("Successfully loaded inventory with %d MCPs", len(inventory.MCPs))
	metrics.RecordLoadOperation(timer.Stop(), len(inventory.MCPs))
	return inventory, nil
}

// saveUnsafe performs the actual save operation without locking (for internal use)
func (s *JSONStorage) saveUnsafe(inventory *models.MCPInventory) error {
	timer := metrics.NewTimer()
	
	// Validate inventory before saving
	validationTimer := metrics.NewTimer()
	if err := inventory.Validate(); err != nil {
		return fmt.Errorf("invalid inventory: %w", err)
	}
	metrics.RecordValidation(validationTimer.Stop())
	
	// Update inventory metadata
	inventory.UpdatedAt = time.Now()
	inventory.Metadata.FileCount = len(inventory.MCPs)
	inventory.Metadata.LastSync = time.Now()
	
	// Convert to JSON
	data, err := inventory.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal inventory: %w", err)
	}
	
	// Atomic write: write to temporary file first, then rename
	tempFile := s.config.GetConfigFile() + ".tmp"
	
	if err := s.writeFile(tempFile, data); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}
	
	// Atomic rename
	if err := os.Rename(tempFile, s.config.GetConfigFile()); err != nil {
		// Clean up temp file on failure
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}
	
	log.Printf("Successfully saved inventory with %d MCPs", len(inventory.MCPs))
	metrics.RecordSaveOperation(timer.Stop(), len(inventory.MCPs))
	return nil
}

// IsLocked returns true if the storage file is currently locked
func (s *JSONStorage) IsLocked() bool {
	return s.fileLock.IsLocked(s.config.GetConfigFile())
}

// WithFileLockOperation is a helper function for operations that need file locking
func WithFileLockOperation[T any](ctx context.Context, filePath string, operation func() (T, error)) (T, error) {
	var result T
	lock, err := GetDefaultFileLock().Lock(ctx, filePath)
	if err != nil {
		return result, fmt.Errorf("failed to acquire file lock: %w", err)
	}
	defer lock.Unlock()
	
	return operation()
}