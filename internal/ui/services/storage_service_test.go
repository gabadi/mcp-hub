package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"mcp-hub/internal/ui/types"
)

func TestGetConfigPath(t *testing.T) {
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	// Should contain the app name
	if !strings.Contains(configPath, appName) {
		t.Errorf("Config path should contain app name %s, got: %s", appName, configPath)
	}

	// Should end with config file name
	if !strings.HasSuffix(configPath, configFileName) {
		t.Errorf("Config path should end with %s, got: %s", configFileName, configPath)
	}
}

func TestEnsureConfigDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test creating config directory
	err := ensureConfigDirWithBase(tempDir)
	if err != nil {
		t.Fatalf("EnsureConfigDir failed: %v", err)
	}

	// Verify directory was created
	expectedDir := filepath.Join(tempDir, appName)
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Config directory was not created: %s", expectedDir)
	}

	// Check permissions (should be 0700)
	info, err := os.Stat(expectedDir)
	if err != nil {
		t.Fatalf("Failed to stat config directory: %v", err)
	}
	if info.Mode().Perm() != 0700 {
		t.Errorf("Config directory permissions should be 0700, got: %o", info.Mode().Perm())
	}
}

func TestSaveAndLoadInventory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test data
	testMCPs := []types.MCPItem{
		{Name: "test1", Type: "CMD", Active: true, Command: "test-command1"},
		{Name: "test2", Type: "SSE", Active: false, Command: "test-command2"},
		{Name: "test3", Type: "JSON", Active: true, Command: "test-command3"},
	}

	// Test saving
	err := saveInventoryWithBase(testMCPs, tempDir)
	if err != nil {
		t.Fatalf("SaveInventory failed: %v", err)
	}

	// Verify file was created
	configPath, _ := getConfigPathWithBase(tempDir)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created: %s", configPath)
	}

	// Test loading
	loadedMCPs, err := loadInventoryWithBase(tempDir)
	if err != nil {
		t.Fatalf("LoadInventory failed: %v", err)
	}

	// Verify loaded data matches saved data
	if len(loadedMCPs) != len(testMCPs) {
		t.Errorf("Expected %d MCPs, got %d", len(testMCPs), len(loadedMCPs))
	}

	for i, expected := range testMCPs {
		if i >= len(loadedMCPs) {
			t.Errorf("Missing MCP at index %d", i)
			continue
		}
		actual := loadedMCPs[i]
		if actual.Name != expected.Name {
			t.Errorf("MCP %d name mismatch: expected %s, got %s", i, expected.Name, actual.Name)
		}
		if actual.Type != expected.Type {
			t.Errorf("MCP %d type mismatch: expected %s, got %s", i, expected.Type, actual.Type)
		}
		if actual.Active != expected.Active {
			t.Errorf("MCP %d active mismatch: expected %t, got %t", i, expected.Active, actual.Active)
		}
		if actual.Command != expected.Command {
			t.Errorf("MCP %d command mismatch: expected %s, got %s", i, expected.Command, actual.Command)
		}
	}
}

func TestLoadInventoryNoFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test loading when no file exists
	loadedMCPs, err := loadInventoryWithBase(tempDir)
	if err != nil {
		t.Fatalf("LoadInventory should not fail when no file exists: %v", err)
	}

	// Should return empty slice
	if len(loadedMCPs) != 0 {
		t.Errorf("Expected empty inventory, got %d items", len(loadedMCPs))
	}
}

func TestLoadInventoryCorruptedFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a corrupted config file
	configPath, err := getConfigPathWithBase(tempDir)
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	// Ensure config directory exists
	if err := ensureConfigDirWithBase(tempDir); err != nil {
		t.Fatalf("EnsureConfigDir failed: %v", err)
	}

	// Write invalid JSON
	corruptedData := `{"invalid": json syntax`
	err = os.WriteFile(configPath, []byte(corruptedData), 0600)
	if err != nil {
		t.Fatalf("Failed to write corrupted config file: %v", err)
	}

	// Test loading corrupted file
	loadedMCPs, err := loadInventoryWithBase(tempDir)
	if err != nil {
		t.Fatalf("LoadInventory should not fail with corrupted file: %v", err)
	}

	// Should return empty slice and backup corrupted file
	if len(loadedMCPs) != 0 {
		t.Errorf("Expected empty inventory for corrupted file, got %d items", len(loadedMCPs))
	}

	// Verify backup file was created
	backupFiles, err := filepath.Glob(configPath + ".corrupted.*")
	if err != nil {
		t.Fatalf("Failed to check for backup files: %v", err)
	}
	if len(backupFiles) == 0 {
		t.Errorf("Expected backup file to be created for corrupted config")
	}
}

func TestJSONSerialization(t *testing.T) {
	// Test MCPItem JSON serialization
	testItem := types.MCPItem{
		Name:    "test-mcp",
		Type:    "CMD",
		Active:  true,
		Command: "test-command",
	}

	// Marshal
	jsonData, err := json.Marshal(testItem)
	if err != nil {
		t.Fatalf("Failed to marshal MCPItem: %v", err)
	}

	// Unmarshal
	var unmarshaled types.MCPItem
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal MCPItem: %v", err)
	}

	// Verify data integrity
	if unmarshaled.Name != testItem.Name {
		t.Errorf("Name mismatch: expected %s, got %s", testItem.Name, unmarshaled.Name)
	}
	if unmarshaled.Type != testItem.Type {
		t.Errorf("Type mismatch: expected %s, got %s", testItem.Type, unmarshaled.Type)
	}
	if unmarshaled.Active != testItem.Active {
		t.Errorf("Active mismatch: expected %t, got %t", testItem.Active, unmarshaled.Active)
	}
	if unmarshaled.Command != testItem.Command {
		t.Errorf("Command mismatch: expected %s, got %s", testItem.Command, unmarshaled.Command)
	}
}

func TestInventoryDataSerialization(t *testing.T) {
	// Test InventoryData structure
	testMCPs := []types.MCPItem{
		{Name: "test1", Type: "CMD", Active: true, Command: "cmd1"},
		{Name: "test2", Type: "SSE", Active: false, Command: "cmd2"},
	}

	inventoryData := InventoryData{
		Version:   configVersion,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Inventory: testMCPs,
	}

	// Marshal
	jsonData, err := json.MarshalIndent(inventoryData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal InventoryData: %v", err)
	}

	// Unmarshal
	var unmarshaled InventoryData
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal InventoryData: %v", err)
	}

	// Verify structure
	if unmarshaled.Version != inventoryData.Version {
		t.Errorf("Version mismatch: expected %s, got %s", inventoryData.Version, unmarshaled.Version)
	}
	if len(unmarshaled.Inventory) != len(inventoryData.Inventory) {
		t.Errorf("Inventory length mismatch: expected %d, got %d", len(inventoryData.Inventory), len(unmarshaled.Inventory))
	}
}

func TestAtomicFileOperations(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	testMCPs := []types.MCPItem{
		{Name: "test", Type: "CMD", Active: true, Command: "test-cmd"},
	}

	// Save inventory
	err := saveInventoryWithBase(testMCPs, tempDir)
	if err != nil {
		t.Fatalf("SaveInventory failed: %v", err)
	}

	// Verify no temporary files are left behind
	configPath, _ := getConfigPathWithBase(tempDir)
	configDir := filepath.Dir(configPath)

	files, err := os.ReadDir(configDir)
	if err != nil {
		t.Fatalf("Failed to read config directory: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tmp") {
			t.Errorf("Temporary file left behind: %s", file.Name())
		}
	}
}

func TestSaveModelInventory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test model
	testMCPs := []types.MCPItem{
		{Name: "model-test", Type: "CMD", Active: true, Command: "model-cmd"},
	}
	model := types.NewModelWithMCPs(testMCPs)

	// Test SaveModelInventory by calling saveInventoryWithBase directly
	err := saveInventoryWithBase(model.MCPItems, tempDir)
	if err != nil {
		t.Fatalf("SaveModelInventory failed: %v", err)
	}

	// Verify data was saved correctly
	loadedMCPs, err := loadInventoryWithBase(tempDir)
	if err != nil {
		t.Fatalf("LoadInventory failed: %v", err)
	}

	if len(loadedMCPs) != 1 {
		t.Errorf("Expected 1 MCP, got %d", len(loadedMCPs))
	}

	if len(loadedMCPs) > 0 && loadedMCPs[0].Name != "model-test" {
		t.Errorf("Expected MCP name 'model-test', got '%s'", loadedMCPs[0].Name)
	}
}

func TestMultipleMCPTypes(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test all MCP types
	testMCPs := []types.MCPItem{
		{Name: "cmd-mcp", Type: "CMD", Active: true, Command: "cmd-binary"},
		{Name: "sse-mcp", Type: "SSE", Active: false, Command: "sse-server"},
		{Name: "json-mcp", Type: "JSON", Active: true, Command: "json-config"},
		{Name: "http-mcp", Type: "HTTP", Active: false, Command: "http-endpoint"},
	}

	// Save and load
	err := saveInventoryWithBase(testMCPs, tempDir)
	if err != nil {
		t.Fatalf("SaveInventory failed: %v", err)
	}

	loadedMCPs, err := loadInventoryWithBase(tempDir)
	if err != nil {
		t.Fatalf("LoadInventory failed: %v", err)
	}

	// Verify all types are preserved
	typeMap := make(map[string]bool)
	for _, mcp := range loadedMCPs {
		typeMap[mcp.Type] = true
	}

	expectedTypes := []string{"CMD", "SSE", "JSON", "HTTP"}
	for _, expectedType := range expectedTypes {
		if !typeMap[expectedType] {
			t.Errorf("MCP type %s was not preserved", expectedType)
		}
	}
}
