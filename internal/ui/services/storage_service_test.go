package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"mcp-hub/internal/platform"
	"mcp-hub/internal/ui/types"
)

func TestGetConfigPath(t *testing.T) {
	mockPlatform := platform.GetMockPlatformService()
	configPath, err := GetConfigPath(mockPlatform)
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
	mockPlatform := platform.GetMockPlatformService()

	// Test creating config directory
	err := ensureConfigDirWithBase(tempDir, mockPlatform)
	if err != nil {
		t.Fatalf("EnsureConfigDir failed: %v", err)
	}

	// Verify directory was created
	expectedDir := filepath.Join(tempDir, appName)
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Config directory was not created: %s", expectedDir)
	}

	// Check permissions (should match platform service expected permissions)
	info, err := os.Stat(expectedDir)
	if err != nil {
		t.Fatalf("Failed to stat config directory: %v", err)
	}
	expectedPerms := mockPlatform.GetDefaultDirectoryPermissions()
	if info.Mode().Perm() != expectedPerms {
		t.Errorf("Config directory permissions should be %o, got: %o", expectedPerms, info.Mode().Perm())
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

	mockPlatform := platform.GetMockPlatformService()
	
	// Test saving
	err := saveInventoryWithBase(testMCPs, tempDir, mockPlatform)
	if err != nil {
		t.Fatalf("SaveInventory failed: %v", err)
	}

	// Verify file was created
	configPath, _ := getConfigPathWithBase(tempDir, mockPlatform)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created: %s", configPath)
	}

	// Test loading
	loadedMCPs, err := loadInventoryWithBase(tempDir, mockPlatform)
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
	mockPlatform := platform.GetMockPlatformService()

	// Test loading when no file exists
	loadedMCPs, err := loadInventoryWithBase(tempDir, mockPlatform)
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
	mockPlatform := platform.GetMockPlatformService()

	// Create a corrupted config file
	configPath, err := getConfigPathWithBase(tempDir, mockPlatform)
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	// Ensure config directory exists
	if err := ensureConfigDirWithBase(tempDir, mockPlatform); err != nil {
		t.Fatalf("EnsureConfigDir failed: %v", err)
	}

	// Write invalid JSON
	corruptedData := `{"invalid": json syntax`
	err = os.WriteFile(configPath, []byte(corruptedData), 0600)
	if err != nil {
		t.Fatalf("Failed to write corrupted config file: %v", err)
	}

	// Test loading corrupted file
	loadedMCPs, err := loadInventoryWithBase(tempDir, mockPlatform)
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
	mockPlatform := platform.GetMockPlatformService()

	testMCPs := []types.MCPItem{
		{Name: "test", Type: "CMD", Active: true, Command: "test-cmd"},
	}

	// Save inventory
	err := saveInventoryWithBase(testMCPs, tempDir, mockPlatform)
	if err != nil {
		t.Fatalf("SaveInventory failed: %v", err)
	}

	// Verify no temporary files are left behind
	configPath, _ := getConfigPathWithBase(tempDir, mockPlatform)
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
	mockPlatform := platform.GetMockPlatformService()
	model := types.NewModelWithMCPs(testMCPs, mockPlatform)

	// Test SaveModelInventory by calling saveInventoryWithBase directly
	err := saveInventoryWithBase(model.MCPItems, tempDir, mockPlatform)
	if err != nil {
		t.Fatalf("SaveModelInventory failed: %v", err)
	}

	// Verify data was saved correctly
	loadedMCPs, err := loadInventoryWithBase(tempDir, mockPlatform)
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
	mockPlatform := platform.GetMockPlatformService()

	// Test all MCP types
	testMCPs := []types.MCPItem{
		{Name: "cmd-mcp", Type: "CMD", Active: true, Command: "cmd-binary"},
		{Name: "sse-mcp", Type: "SSE", Active: false, Command: "sse-server"},
		{Name: "json-mcp", Type: "JSON", Active: true, Command: "json-config"},
		{Name: "http-mcp", Type: "HTTP", Active: false, Command: "http-endpoint"},
	}

	// Save and load
	err := saveInventoryWithBase(testMCPs, tempDir, mockPlatform)
	if err != nil {
		t.Fatalf("SaveInventory failed: %v", err)
	}

	loadedMCPs, err := loadInventoryWithBase(tempDir, mockPlatform)
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

// Test LoadInventory comprehensive scenarios
func TestLoadInventoryComprehensive(t *testing.T) {
	t.Run("load_empty_inventory", func(t *testing.T) {
		// Test loading when no file exists
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		inventory, err := loadInventoryWithBase(tempDir, mockPlatform)
		if err != nil {
			t.Errorf("LoadInventory should not fail when no file exists: %v", err)
		}
		
		if len(inventory) != 0 {
			t.Errorf("Expected empty inventory, got %d items", len(inventory))
		}
	})
	
	t.Run("load_with_invalid_json", func(t *testing.T) {
		// Test loading with invalid JSON
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		// Create config directory and file with invalid JSON
		appDir := filepath.Join(tempDir, appName)
		err := os.MkdirAll(appDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create app directory: %v", err)
		}
		
		configPath := filepath.Join(appDir, configFileName)
		err = os.WriteFile(configPath, []byte(`{invalid json`), 0600)
		if err != nil {
			t.Fatalf("Failed to write invalid JSON: %v", err)
		}
		
		// Should return empty inventory and create backup
		inventory, err := loadInventoryWithBase(tempDir, mockPlatform)
		if err != nil {
			t.Errorf("LoadInventory should not fail with invalid JSON: %v", err)
		}
		
		if len(inventory) != 0 {
			t.Errorf("Expected empty inventory with invalid JSON, got %d items", len(inventory))
		}
		
		// Verify backup file was created
		backupFiles, err := filepath.Glob(configPath + ".corrupted.*")
		if err != nil {
			t.Errorf("Failed to check for backup files: %v", err)
		}
		if len(backupFiles) == 0 {
			t.Error("Expected backup file to be created")
		}
	})
	
	t.Run("load_with_missing_fields", func(t *testing.T) {
		// Test loading with missing required fields
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		appDir := filepath.Join(tempDir, appName)
		err := os.MkdirAll(appDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create app directory: %v", err)
		}
		
		// Create JSON with missing fields
		incompleteJSON := `{"version": "1.0", "inventory": [{"name": "test"}]}`
		configPath := filepath.Join(appDir, configFileName)
		err = os.WriteFile(configPath, []byte(incompleteJSON), 0600)
		if err != nil {
			t.Fatalf("Failed to write incomplete JSON: %v", err)
		}
		
		// Should still load successfully with default values
		inventory, err := loadInventoryWithBase(tempDir, mockPlatform)
		if err != nil {
			t.Errorf("LoadInventory should handle missing fields: %v", err)
		}
		
		if len(inventory) != 1 {
			t.Errorf("Expected 1 inventory item, got %d", len(inventory))
		}
		
		if len(inventory) > 0 && inventory[0].Name != "test" {
			t.Errorf("Expected item name 'test', got '%s'", inventory[0].Name)
		}
	})
	
	t.Run("load_with_large_inventory", func(t *testing.T) {
		// Test loading with large inventory
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		// Create large inventory
		largeInventory := make([]types.MCPItem, 1000)
		for i := 0; i < 1000; i++ {
			largeInventory[i] = types.MCPItem{
				Name:    fmt.Sprintf("mcp-%d", i),
				Type:    "CMD",
				Active:  i%2 == 0,
				Command: fmt.Sprintf("command-%d", i),
			}
		}
		
		// Save large inventory
		err := saveInventoryWithBase(largeInventory, tempDir, mockPlatform)
		if err != nil {
			t.Fatalf("Failed to save large inventory: %v", err)
		}
		
		// Load large inventory
		loadedInventory, err := loadInventoryWithBase(tempDir, mockPlatform)
		if err != nil {
			t.Fatalf("Failed to load large inventory: %v", err)
		}
		
		if len(loadedInventory) != 1000 {
			t.Errorf("Expected 1000 inventory items, got %d", len(loadedInventory))
		}
	})
	
	t.Run("load_with_version_mismatch", func(t *testing.T) {
		// Test loading with different version
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		appDir := filepath.Join(tempDir, appName)
		err := os.MkdirAll(appDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create app directory: %v", err)
		}
		
		// Create JSON with different version
		oldVersionJSON := `{"version": "0.9", "timestamp": "2023-01-01T00:00:00Z", "inventory": [{"name": "test", "type": "CMD", "active": true, "command": "test-cmd"}]}`
		configPath := filepath.Join(appDir, configFileName)
		err = os.WriteFile(configPath, []byte(oldVersionJSON), 0600)
		if err != nil {
			t.Fatalf("Failed to write old version JSON: %v", err)
		}
		
		// Should still load successfully
		inventory, err := loadInventoryWithBase(tempDir, mockPlatform)
		if err != nil {
			t.Errorf("LoadInventory should handle version mismatch: %v", err)
		}
		
		if len(inventory) != 1 {
			t.Errorf("Expected 1 inventory item, got %d", len(inventory))
		}
	})
}

// Test EnsureConfigDir comprehensive scenarios
func TestEnsureConfigDirComprehensive(t *testing.T) {
	t.Run("ensure_config_dir_success", func(t *testing.T) {
		// Test successful directory creation
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		err := ensureConfigDirWithBase(tempDir, mockPlatform)
		if err != nil {
			t.Errorf("EnsureConfigDir should succeed: %v", err)
		}
		
		// Verify directory exists
		appDir := filepath.Join(tempDir, appName)
		info, err := os.Stat(appDir)
		if err != nil {
			t.Errorf("Config directory should exist: %v", err)
		}
		
		if !info.IsDir() {
			t.Error("Config path should be a directory")
		}
	})
	
	t.Run("ensure_config_dir_already_exists", func(t *testing.T) {
		// Test when directory already exists
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		// Create directory first
		appDir := filepath.Join(tempDir, appName)
		err := os.MkdirAll(appDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create initial directory: %v", err)
		}
		
		// Should not fail when directory exists
		err = ensureConfigDirWithBase(tempDir, mockPlatform)
		if err != nil {
			t.Errorf("EnsureConfigDir should not fail when directory exists: %v", err)
		}
	})
	
	t.Run("ensure_config_dir_nested_creation", func(t *testing.T) {
		// Test nested directory creation
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		// Use a deeper nested path
		nestedBase := filepath.Join(tempDir, "deep", "nested", "path")
		err := ensureConfigDirWithBase(nestedBase, mockPlatform)
		if err != nil {
			t.Errorf("EnsureConfigDir should handle nested paths: %v", err)
		}
		
		// Verify nested directory exists
		appDir := filepath.Join(nestedBase, appName)
		info, err := os.Stat(appDir)
		if err != nil {
			t.Errorf("Nested config directory should exist: %v", err)
		}
		
		if !info.IsDir() {
			t.Error("Nested config path should be a directory")
		}
	})
	
	t.Run("ensure_config_dir_permissions", func(t *testing.T) {
		// Test directory permissions
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		err := ensureConfigDirWithBase(tempDir, mockPlatform)
		if err != nil {
			t.Errorf("EnsureConfigDir should succeed: %v", err)
		}
		
		// Check directory permissions
		appDir := filepath.Join(tempDir, appName)
		info, err := os.Stat(appDir)
		if err != nil {
			t.Errorf("Config directory should exist: %v", err)
		}
		
		expectedPerms := mockPlatform.GetDefaultDirectoryPermissions()
		if info.Mode().Perm() != expectedPerms {
			t.Errorf("Expected directory permissions %o, got %o", expectedPerms, info.Mode().Perm())
		}
	})
	
	t.Run("ensure_config_dir_with_production_platform", func(t *testing.T) {
		// Test with real platform service
		platformService := platform.NewPlatformServiceFactoryDefault().CreatePlatformService()
		
		// Test getting config path (should not be empty)
		configPath := platformService.GetConfigPath()
		if configPath == "" {
			t.Error("Production platform should provide config path")
		}
		
		// Test getting permissions (should not be zero)
		dirPerms := platformService.GetDefaultDirectoryPermissions()
		if dirPerms == 0 {
			t.Error("Production platform should provide directory permissions")
		}
	})
}

// Test SaveModelInventory comprehensive scenarios
func TestSaveModelInventoryComprehensive(t *testing.T) {
	t.Run("save_model_inventory_success", func(t *testing.T) {
		// Test successful model inventory save
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		testMCPs := []types.MCPItem{
			{Name: "model-test", Type: "CMD", Active: true, Command: "model-cmd"},
			{Name: "model-test-2", Type: "SSE", Active: false, Command: "model-sse"},
		}
		
		model := types.NewModelWithMCPs(testMCPs, mockPlatform)
		
		// Save using the model save function
		err := saveInventoryWithBase(model.MCPItems, tempDir, mockPlatform)
		if err != nil {
			t.Errorf("SaveModelInventory should succeed: %v", err)
		}
		
		// Verify data was saved
		loadedMCPs, err := loadInventoryWithBase(tempDir, mockPlatform)
		if err != nil {
			t.Errorf("LoadInventory should succeed after save: %v", err)
		}
		
		if len(loadedMCPs) != 2 {
			t.Errorf("Expected 2 MCPs, got %d", len(loadedMCPs))
		}
	})
	
	t.Run("save_model_inventory_empty", func(t *testing.T) {
		// Test saving empty model inventory
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		model := types.NewModelWithMCPs([]types.MCPItem{}, mockPlatform)
		
		err := saveInventoryWithBase(model.MCPItems, tempDir, mockPlatform)
		if err != nil {
			t.Errorf("SaveModelInventory should handle empty inventory: %v", err)
		}
		
		// Verify empty inventory was saved
		loadedMCPs, err := loadInventoryWithBase(tempDir, mockPlatform)
		if err != nil {
			t.Errorf("LoadInventory should succeed with empty inventory: %v", err)
		}
		
		if len(loadedMCPs) != 0 {
			t.Errorf("Expected empty inventory, got %d items", len(loadedMCPs))
		}
	})
	
	t.Run("save_model_inventory_overwrite", func(t *testing.T) {
		// Test overwriting existing inventory
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		// Save initial inventory
		initialMCPs := []types.MCPItem{
			{Name: "initial", Type: "CMD", Active: true, Command: "initial-cmd"},
		}
		err := saveInventoryWithBase(initialMCPs, tempDir, mockPlatform)
		if err != nil {
			t.Fatalf("Failed to save initial inventory: %v", err)
		}
		
		// Save new inventory via model
		newMCPs := []types.MCPItem{
			{Name: "new-1", Type: "SSE", Active: false, Command: "new-cmd-1"},
			{Name: "new-2", Type: "JSON", Active: true, Command: "new-cmd-2"},
		}
		model := types.NewModelWithMCPs(newMCPs, mockPlatform)
		
		err = saveInventoryWithBase(model.MCPItems, tempDir, mockPlatform)
		if err != nil {
			t.Errorf("SaveModelInventory should overwrite existing: %v", err)
		}
		
		// Verify new inventory replaced old
		loadedMCPs, err := loadInventoryWithBase(tempDir, mockPlatform)
		if err != nil {
			t.Errorf("LoadInventory should succeed after overwrite: %v", err)
		}
		
		if len(loadedMCPs) != 2 {
			t.Errorf("Expected 2 MCPs after overwrite, got %d", len(loadedMCPs))
		}
		
		// Check that old inventory is gone
		for _, mcp := range loadedMCPs {
			if mcp.Name == "initial" {
				t.Error("Old inventory should be overwritten")
			}
		}
	})
}

// Test security and validation scenarios
func TestSecurityAndValidation(t *testing.T) {
	t.Run("validate_config_path_security", func(t *testing.T) {
		// Test config path validation
		mockPlatform := platform.GetMockPlatformService()
		
		// Test with valid config path
		validPath := filepath.Join(mockPlatform.GetConfigPath(), "inventory.json")
		err := validateConfigPath(validPath, mockPlatform)
		if err != nil {
			t.Errorf("ValidateConfigPath should accept valid path: %v", err)
		}
		
		// Test with invalid config path (outside expected directory)
		invalidPath := "/tmp/malicious/inventory.json"
		err = validateConfigPath(invalidPath, mockPlatform)
		if err == nil {
			t.Error("ValidateConfigPath should reject invalid path")
		}
	})
	
	t.Run("safe_file_reading", func(t *testing.T) {
		// Test safe file reading with size limits
		tempDir := t.TempDir()
		_ = tempDir // Use tempDir to avoid unused variable warning
		
		// Create a normal sized file
		normalFile := filepath.Join(tempDir, "normal-inventory.json")
		normalContent := `{"version": "1.0", "timestamp": "2023-01-01T00:00:00Z", "inventory": []}`
		err := os.WriteFile(normalFile, []byte(normalContent), 0600)
		if err != nil {
			t.Fatalf("Failed to create normal file: %v", err)
		}
		
		// Should read successfully
		data, err := readSecureFile(normalFile)
		if err != nil {
			t.Errorf("ReadSecureFile should handle normal file: %v", err)
		}
		
		if string(data) != normalContent {
			t.Error("ReadSecureFile should return correct content")
		}
	})
	
	t.Run("file_type_validation", func(t *testing.T) {
		// Test file type validation
		tempDir := t.TempDir()
		_ = tempDir // Use tempDir to avoid unused variable warning
		
		// Create a directory instead of file
		dirPath := filepath.Join(tempDir, "not-a-file")
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		
		// Should fail when trying to read directory
		_, err = readSecureFile(dirPath)
		if err == nil {
			t.Error("ReadSecureFile should reject directory")
		}
	})
	
	t.Run("file_path_cleaning", func(t *testing.T) {
		// Test file path cleaning to prevent path traversal
		tempDir := t.TempDir()
		
		// Create a normal file
		normalFile := filepath.Join(tempDir, "test-inventory.json")
		content := `{"version": "1.0", "inventory": []}`
		err := os.WriteFile(normalFile, []byte(content), 0600)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		
		// Test reading with path traversal attempt
		traversalPath := filepath.Join(tempDir, "../", filepath.Base(tempDir), "test-inventory.json")
		
		// Should still read the file safely (after path cleaning)
		data, err := readSecureFile(traversalPath)
		if err != nil {
			t.Errorf("ReadSecureFile should handle path traversal safely: %v", err)
		}
		
		if string(data) != content {
			t.Error("ReadSecureFile should return correct content after path cleaning")
		}
	})
	
	t.Run("atomic_file_operations", func(t *testing.T) {
		// Test atomic file operations
		tempDir := t.TempDir()
		mockPlatform := platform.GetMockPlatformService()
		
		testMCPs := []types.MCPItem{
			{Name: "atomic-test", Type: "CMD", Active: true, Command: "atomic-cmd"},
		}
		
		// Save should be atomic
		err := saveInventoryWithBase(testMCPs, tempDir, mockPlatform)
		if err != nil {
			t.Errorf("Atomic save should succeed: %v", err)
		}
		
		// Verify no temporary files remain
		appDir := filepath.Join(tempDir, appName)
		files, err := os.ReadDir(appDir)
		if err != nil {
			t.Errorf("Failed to read app directory: %v", err)
		}
		
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".tmp") {
				t.Errorf("Temporary file should not remain: %s", file.Name())
			}
		}
		
		// Verify final file exists
		configPath := filepath.Join(appDir, configFileName)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("Final config file should exist after atomic save")
		}
	})
}
