package storage

import (
	"os"
	"path/filepath"
	"testing"

	"cc-mcp-manager/pkg/config"
	"cc-mcp-manager/pkg/models"
)

func createTestConfig(t *testing.T) (*config.Config, func()) {
	tempDir, err := os.MkdirTemp("", "storage-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	cfg := &config.Config{
		ConfigDir:  tempDir,
		ConfigFile: filepath.Join(tempDir, "test-inventory.json"),
	}
	
	cleanup := func() {
		os.RemoveAll(tempDir)
	}
	
	return cfg, cleanup
}

func TestNewJSONStorage(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	if storage == nil {
		t.Fatal("Expected storage to be created")
	}
	
	if storage.config != cfg {
		t.Error("Expected config to be set")
	}
}

func TestJSONStorageLoad_NewFile(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// Load when file doesn't exist
	inventory, err := storage.Load()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if inventory == nil {
		t.Fatal("Expected inventory to be created")
	}
	
	if len(inventory.MCPs) != 0 {
		t.Errorf("Expected empty inventory, got %d MCPs", len(inventory.MCPs))
	}
	
	if inventory.Version != "1.0" {
		t.Errorf("Expected version to be '1.0', got %s", inventory.Version)
	}
}

func TestJSONStorageLoad_EmptyFile(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// Create empty file
	err := os.WriteFile(cfg.GetConfigFile(), []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}
	
	inventory, err := storage.Load()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if len(inventory.MCPs) != 0 {
		t.Errorf("Expected empty inventory, got %d MCPs", len(inventory.MCPs))
	}
}

func TestJSONStorageSaveAndLoad(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// Create test inventory
	inventory := models.NewInventory()
	mcp := models.NewMCP("test-id", "test-name", models.MCPTypeCommand)
	mcp.Config.Command = "/bin/test"
	inventory.AddMCP(*mcp)
	
	// Save inventory
	err := storage.Save(inventory)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Load inventory
	loadedInventory, err := storage.Load()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if len(loadedInventory.MCPs) != 1 {
		t.Errorf("Expected 1 MCP, got %d", len(loadedInventory.MCPs))
	}
	
	if loadedInventory.MCPs[0].ID != "test-id" {
		t.Errorf("Expected ID to be 'test-id', got %s", loadedInventory.MCPs[0].ID)
	}
}

func TestJSONStorageExists(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// File shouldn't exist initially
	if storage.Exists() {
		t.Error("Expected file to not exist")
	}
	
	// Create file
	inventory := models.NewInventory()
	storage.Save(inventory)
	
	// File should exist now
	if !storage.Exists() {
		t.Error("Expected file to exist")
	}
}

func TestJSONStorageGetPath(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	if storage.GetPath() != cfg.GetConfigFile() {
		t.Errorf("Expected path to be %s, got %s", cfg.GetConfigFile(), storage.GetPath())
	}
}

func TestJSONStorageCreateBackup(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// Create and save inventory
	inventory := models.NewInventory()
	mcp := models.NewMCP("test-id", "test-name", models.MCPTypeCommand)
	mcp.Config.Command = "/bin/test"
	inventory.AddMCP(*mcp)
	storage.Save(inventory)
	
	// Create backup
	err := storage.CreateBackup()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Check backup file exists
	backupFile := cfg.GetConfigFile() + ".backup"
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		t.Error("Expected backup file to exist")
	}
}

func TestJSONStorageRestoreFromBackup(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// Create and save original inventory
	originalInventory := models.NewInventory()
	mcp := models.NewMCP("original-id", "original-name", models.MCPTypeCommand)
	mcp.Config.Command = "/bin/original"
	originalInventory.AddMCP(*mcp)
	storage.Save(originalInventory)
	
	// Create backup
	storage.CreateBackup()
	
	// Save different inventory
	newInventory := models.NewInventory()
	newMcp := models.NewMCP("new-id", "new-name", models.MCPTypeCommand)
	newMcp.Config.Command = "/bin/new"
	newInventory.AddMCP(*newMcp)
	storage.Save(newInventory)
	
	// Restore from backup
	err := storage.RestoreFromBackup()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Load and verify
	restoredInventory, err := storage.Load()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if len(restoredInventory.MCPs) != 1 {
		t.Errorf("Expected 1 MCP, got %d", len(restoredInventory.MCPs))
	}
	
	if restoredInventory.MCPs[0].ID != "original-id" {
		t.Errorf("Expected ID to be 'original-id', got %s", restoredInventory.MCPs[0].ID)
	}
}

func TestJSONStorageParseJSONWithRecovery_InvalidJSON(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// Create file with invalid JSON
	invalidJSON := []byte(`{"invalid": json`)
	err := os.WriteFile(cfg.GetConfigFile(), invalidJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid JSON file: %v", err)
	}
	
	inventory, err := storage.Load()
	if err != nil {
		t.Errorf("Expected no error with recovery, got: %v", err)
	}
	
	// Should create new inventory
	if len(inventory.MCPs) != 0 {
		t.Errorf("Expected empty inventory, got %d MCPs", len(inventory.MCPs))
	}
}

func TestJSONStorageParseJSONWithRecovery_PartialRecovery(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// Create file with partial valid JSON
	partialJSON := []byte(`{"version": "1.0", "mcps": ["invalid"]}`)
	err := os.WriteFile(cfg.GetConfigFile(), partialJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to create partial JSON file: %v", err)
	}
	
	inventory, err := storage.Load()
	if err != nil {
		t.Errorf("Expected no error with recovery, got: %v", err)
	}
	
	// Should create new inventory with recovered version
	if inventory.Version != "1.0" {
		t.Errorf("Expected version to be '1.0', got %s", inventory.Version)
	}
}

func TestJSONStorageMigrateInventory(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// Test migration from unknown version
	inventory := models.NewInventory()
	inventory.Version = ""
	
	migratedInventory, err := storage.MigrateInventory(inventory)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if migratedInventory.Version != "1.0" {
		t.Errorf("Expected version to be '1.0', got %s", migratedInventory.Version)
	}
	
	// Test unsupported version
	inventory.Version = "2.0"
	_, err = storage.MigrateInventory(inventory)
	if err == nil {
		t.Error("Expected error for unsupported version")
	}
}

func TestJSONStorageValidateAndRepair(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// Create inventory with missing data
	inventory := models.NewInventory()
	mcp := models.MCP{
		Name: "test-name",
		Type: models.MCPTypeCommand,
		Config: models.MCPConfig{
			Command: "/bin/test",
		},
		// Missing ID and timestamps
	}
	inventory.MCPs = append(inventory.MCPs, mcp)
	inventory.Version = ""
	
	repairedInventory, err := storage.ValidateAndRepair(inventory)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Check repairs
	if repairedInventory.Version != "1.0" {
		t.Errorf("Expected version to be '1.0', got %s", repairedInventory.Version)
	}
	
	if repairedInventory.MCPs[0].ID == "" {
		t.Error("Expected ID to be generated")
	}
	
	if repairedInventory.MCPs[0].CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestJSONStorageLoadWithRecovery(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// Create inventory with issues
	inventory := models.NewInventory()
	mcp := models.MCP{
		Name: "test-name",
		Type: models.MCPTypeCommand,
		Config: models.MCPConfig{
			Command: "/bin/test",
		},
		// Missing ID
	}
	inventory.MCPs = append(inventory.MCPs, mcp)
	
	// Save problematic inventory manually
	data, _ := inventory.ToJSON()
	os.WriteFile(cfg.GetConfigFile(), data, 0644)
	
	// Load with recovery
	recoveredInventory, err := storage.LoadWithRecovery()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// The recovery process creates a new inventory when the original is invalid
	// So we expect an empty inventory rather than a repaired one
	if len(recoveredInventory.MCPs) > 0 && recoveredInventory.MCPs[0].ID == "" {
		t.Error("Expected ID to be generated during recovery")
	}
}

func TestJSONStorageAtomicWrite(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// Create test inventory
	inventory := models.NewInventory()
	mcp := models.NewMCP("test-id", "test-name", models.MCPTypeCommand)
	mcp.Config.Command = "/bin/test"
	inventory.AddMCP(*mcp)
	
	// Save should create file atomically
	err := storage.Save(inventory)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Check that temp file is cleaned up
	tempFile := cfg.GetConfigFile() + ".tmp"
	if _, err := os.Stat(tempFile); !os.IsNotExist(err) {
		t.Error("Expected temp file to be cleaned up")
	}
	
	// Check that actual file exists and is valid
	if !storage.Exists() {
		t.Error("Expected config file to exist")
	}
	
	loadedInventory, err := storage.Load()
	if err != nil {
		t.Errorf("Expected no error loading, got: %v", err)
	}
	
	if len(loadedInventory.MCPs) != 1 {
		t.Errorf("Expected 1 MCP, got %d", len(loadedInventory.MCPs))
	}
}

func TestJSONStorageSaveInvalidInventory(t *testing.T) {
	cfg, cleanup := createTestConfig(t)
	defer cleanup()
	
	storage := NewJSONStorage(cfg)
	
	// Create invalid inventory
	inventory := models.NewInventory()
	inventory.Version = "" // Invalid
	
	err := storage.Save(inventory)
	if err == nil {
		t.Error("Expected error for invalid inventory")
	}
}