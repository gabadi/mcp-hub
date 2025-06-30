package models

import (
	"testing"
	"time"
)

func TestNewMCP(t *testing.T) {
	mcp := NewMCP("test-id", "test-name", MCPTypeCommand)
	
	if mcp.ID != "test-id" {
		t.Errorf("Expected ID to be 'test-id', got %s", mcp.ID)
	}
	
	if mcp.Name != "test-name" {
		t.Errorf("Expected Name to be 'test-name', got %s", mcp.Name)
	}
	
	if mcp.Type != MCPTypeCommand {
		t.Errorf("Expected Type to be MCPTypeCommand, got %s", mcp.Type)
	}
	
	if !mcp.Enabled {
		t.Error("Expected MCP to be enabled by default")
	}
	
	if mcp.CreatedAt.IsZero() || mcp.UpdatedAt.IsZero() {
		t.Error("Expected timestamps to be set")
	}
}

func TestMCPValidate(t *testing.T) {
	tests := []struct {
		name        string
		mcp         MCP
		expectError bool
	}{
		{
			name: "valid command MCP",
			mcp: MCP{
				ID:        "test",
				Name:      "Test MCP",
				Type:      MCPTypeCommand,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Config: MCPConfig{
					Command: "/bin/test",
				},
			},
			expectError: false,
		},
		{
			name: "valid SSE MCP",
			mcp: MCP{
				ID:        "test",
				Name:      "Test SSE",
				Type:      MCPTypeSSE,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Config: MCPConfig{
					ServerURL: "http://localhost:8080",
				},
			},
			expectError: false,
		},
		{
			name: "valid JSON MCP",
			mcp: MCP{
				ID:        "test",
				Name:      "Test JSON",
				Type:      MCPTypeJSON,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Config: MCPConfig{
					JSONConfig: map[string]interface{}{
						"key": "value",
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid HTTP MCP",
			mcp: MCP{
				ID:        "test",
				Name:      "Test HTTP",
				Type:      MCPTypeHTTP,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Config: MCPConfig{
					Endpoint: "http://api.example.com",
				},
			},
			expectError: false,
		},
		{
			name: "missing name",
			mcp: MCP{
				ID:   "test",
				Type: MCPTypeCommand,
			},
			expectError: true,
		},
		{
			name: "missing ID",
			mcp: MCP{
				Name: "Test",
				Type: MCPTypeCommand,
			},
			expectError: true,
		},
		{
			name: "command MCP missing command",
			mcp: MCP{
				ID:   "test",
				Name: "Test",
				Type: MCPTypeCommand,
			},
			expectError: true,
		},
		{
			name: "SSE MCP missing server URL",
			mcp: MCP{
				ID:   "test",
				Name: "Test",
				Type: MCPTypeSSE,
			},
			expectError: true,
		},
		{
			name: "JSON MCP missing config",
			mcp: MCP{
				ID:   "test",
				Name: "Test",
				Type: MCPTypeJSON,
			},
			expectError: true,
		},
		{
			name: "HTTP MCP missing endpoint",
			mcp: MCP{
				ID:   "test",
				Name: "Test",
				Type: MCPTypeHTTP,
			},
			expectError: true,
		},
		{
			name: "unsupported MCP type",
			mcp: MCP{
				ID:   "test",
				Name: "Test",
				Type: "unsupported",
			},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mcp.Validate()
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestNewInventory(t *testing.T) {
	inv := NewInventory()
	
	if inv.Version != "1.0" {
		t.Errorf("Expected version to be '1.0', got %s", inv.Version)
	}
	
	if len(inv.MCPs) != 0 {
		t.Errorf("Expected empty MCPs slice, got %d items", len(inv.MCPs))
	}
	
	if inv.Metadata.FileCount != 0 {
		t.Errorf("Expected FileCount to be 0, got %d", inv.Metadata.FileCount)
	}
	
	if inv.Metadata.Created.IsZero() || inv.UpdatedAt.IsZero() {
		t.Error("Expected timestamps to be set")
	}
}

func TestInventoryAddMCP(t *testing.T) {
	inv := NewInventory()
	mcp := NewMCP("test-id", "test-name", MCPTypeCommand)
	mcp.Config.Command = "/bin/test"
	
	err := inv.AddMCP(*mcp)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if len(inv.MCPs) != 1 {
		t.Errorf("Expected 1 MCP, got %d", len(inv.MCPs))
	}
	
	if inv.Metadata.FileCount != 1 {
		t.Errorf("Expected FileCount to be 1, got %d", inv.Metadata.FileCount)
	}
	
	// Test duplicate ID
	err = inv.AddMCP(*mcp)
	if err == nil {
		t.Error("Expected error for duplicate ID")
	}
}

func TestInventoryRemoveMCP(t *testing.T) {
	inv := NewInventory()
	mcp := NewMCP("test-id", "test-name", MCPTypeCommand)
	mcp.Config.Command = "/bin/test"
	
	inv.AddMCP(*mcp)
	
	err := inv.RemoveMCP("test-id")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if len(inv.MCPs) != 0 {
		t.Errorf("Expected 0 MCPs, got %d", len(inv.MCPs))
	}
	
	// Test non-existent ID
	err = inv.RemoveMCP("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent ID")
	}
}

func TestInventoryGetMCP(t *testing.T) {
	inv := NewInventory()
	mcp := NewMCP("test-id", "test-name", MCPTypeCommand)
	mcp.Config.Command = "/bin/test"
	
	inv.AddMCP(*mcp)
	
	foundMCP, err := inv.GetMCP("test-id")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if foundMCP.ID != "test-id" {
		t.Errorf("Expected ID to be 'test-id', got %s", foundMCP.ID)
	}
	
	// Test non-existent ID
	_, err = inv.GetMCP("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent ID")
	}
}

func TestInventoryUpdateMCP(t *testing.T) {
	inv := NewInventory()
	mcp := NewMCP("test-id", "test-name", MCPTypeCommand)
	mcp.Config.Command = "/bin/test"
	
	inv.AddMCP(*mcp)
	
	// Update the MCP
	updatedMCP := *mcp
	updatedMCP.Name = "Updated Name"
	
	err := inv.UpdateMCP(updatedMCP)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	foundMCP, _ := inv.GetMCP("test-id")
	if foundMCP.Name != "Updated Name" {
		t.Errorf("Expected name to be 'Updated Name', got %s", foundMCP.Name)
	}
	
	// Test non-existent ID
	nonExistentMCP := NewMCP("non-existent", "test", MCPTypeCommand)
	nonExistentMCP.Config.Command = "/bin/test"
	err = inv.UpdateMCP(*nonExistentMCP)
	if err == nil {
		t.Error("Expected error for non-existent ID")
	}
}

func TestInventoryValidate(t *testing.T) {
	// Valid inventory
	inv := NewInventory()
	mcp := NewMCP("test-id", "test-name", MCPTypeCommand)
	mcp.Config.Command = "/bin/test"
	inv.AddMCP(*mcp)
	
	err := inv.Validate()
	if err != nil {
		t.Errorf("Expected no error for valid inventory, got: %v", err)
	}
	
	// Invalid inventory - missing version
	inv.Version = ""
	err = inv.Validate()
	if err == nil {
		t.Error("Expected error for missing version")
	}
	
	// Invalid inventory - duplicate IDs
	inv.Version = "1.0"
	duplicateMCP := *mcp
	inv.MCPs = append(inv.MCPs, duplicateMCP)
	err = inv.Validate()
	if err == nil {
		t.Error("Expected error for duplicate IDs")
	}
}

func TestInventoryJSON(t *testing.T) {
	inv := NewInventory()
	mcp := NewMCP("test-id", "test-name", MCPTypeCommand)
	mcp.Config.Command = "/bin/test"
	inv.AddMCP(*mcp)
	
	// Test ToJSON
	data, err := inv.ToJSON()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if len(data) == 0 {
		t.Error("Expected JSON data to be non-empty")
	}
	
	// Test FromJSON
	parsedInv, err := FromJSON(data)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if len(parsedInv.MCPs) != 1 {
		t.Errorf("Expected 1 MCP, got %d", len(parsedInv.MCPs))
	}
	
	if parsedInv.MCPs[0].ID != "test-id" {
		t.Errorf("Expected ID to be 'test-id', got %s", parsedInv.MCPs[0].ID)
	}
}

func TestFromJSONInvalid(t *testing.T) {
	// Invalid JSON
	_, err := FromJSON([]byte("invalid json"))
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	
	// Valid JSON but invalid inventory
	invalidData := `{"version": "", "mcps": []}`
	_, err = FromJSON([]byte(invalidData))
	if err == nil {
		t.Error("Expected error for invalid inventory")
	}
}

func TestMCPConfigDefaults(t *testing.T) {
	mcp := MCP{
		ID:        "test",
		Name:      "Test HTTP",
		Type:      MCPTypeHTTP,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Config: MCPConfig{
			Endpoint: "http://api.example.com",
		},
	}
	
	// Validate should set default method
	err := mcp.Validate()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if mcp.Config.Method != "GET" {
		t.Errorf("Expected default method to be 'GET', got %s", mcp.Config.Method)
	}
}