package services

import (
	"testing"

	"cc-mcp-manager/internal/ui/types"
)

// Test helpers for MCP service tests
func createTestModel() types.Model {
	return types.Model{
		SearchQuery:  "",
		SelectedItem: 0,
		MCPItems: []types.MCPItem{
			{Name: "context7", Type: "SSE", Active: true, Command: "npx @context7/mcp-server"},
			{Name: "github-mcp", Type: "CMD", Active: true, Command: "github-mcp"},
			{Name: "ht-mcp", Type: "CMD", Active: false, Command: "ht-mcp"},
			{Name: "filesystem", Type: "CMD", Active: false, Command: "filesystem-mcp"},
			{Name: "docker-mcp", Type: "CMD", Active: false, Command: "docker-mcp"},
		},
	}
}

func TestGetFilteredMCPs(t *testing.T) {
	tests := []struct {
		name          string
		searchQuery   string
		expectedCount int
		expectedFirst string
	}{
		{
			name:          "Empty query returns all MCPs",
			searchQuery:   "",
			expectedCount: 5,
			expectedFirst: "context7",
		},
		{
			name:          "Case insensitive filtering",
			searchQuery:   "GITHUB",
			expectedCount: 1,
			expectedFirst: "github-mcp",
		},
		{
			name:          "Partial name matching",
			searchQuery:   "mcp",
			expectedCount: 3, // github-mcp, ht-mcp, docker-mcp
			expectedFirst: "github-mcp",
		},
		{
			name:          "No matches found",
			searchQuery:   "nonexistent",
			expectedCount: 0,
			expectedFirst: "",
		},
		{
			name:          "Single character search",
			searchQuery:   "h",
			expectedCount: 2, // github-mcp, ht-mcp
			expectedFirst: "github-mcp",
		},
		{
			name:          "Exact name match",
			searchQuery:   "context7",
			expectedCount: 1,
			expectedFirst: "context7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestModel()
			model.SearchQuery = tt.searchQuery

			filtered := GetFilteredMCPs(model)

			if len(filtered) != tt.expectedCount {
				t.Errorf("GetFilteredMCPs() returned %d items, expected %d", len(filtered), tt.expectedCount)
			}

			if tt.expectedCount > 0 && filtered[0].Name != tt.expectedFirst {
				t.Errorf("GetFilteredMCPs() first item = %s, expected %s", filtered[0].Name, tt.expectedFirst)
			}
		})
	}
}

func TestToggleMCPStatus(t *testing.T) {
	tests := []struct {
		name            string
		selectedItem    int
		searchQuery     string
		expectedActive  bool
		expectedChanged string
	}{
		{
			name:            "Toggle active MCP to inactive",
			selectedItem:    0, // context7 (active)
			searchQuery:     "",
			expectedActive:  false,
			expectedChanged: "context7",
		},
		{
			name:            "Toggle inactive MCP to active",
			selectedItem:    2, // ht-mcp (inactive)
			searchQuery:     "",
			expectedActive:  true,
			expectedChanged: "ht-mcp",
		},
		{
			name:            "Toggle with filtered results",
			selectedItem:    0, // first in filtered results
			searchQuery:     "github",
			expectedActive:  false, // github-mcp starts active
			expectedChanged: "github-mcp",
		},
		{
			name:            "Toggle with out of bounds selection",
			selectedItem:    10, // beyond array bounds
			searchQuery:     "",
			expectedActive:  true, // context7 should remain unchanged
			expectedChanged: "context7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestModel()
			model.SelectedItem = tt.selectedItem
			model.SearchQuery = tt.searchQuery

			// Store original state for comparison
			originalStates := make(map[string]bool)
			for _, item := range model.MCPItems {
				originalStates[item.Name] = item.Active
			}

			updatedModel := ToggleMCPStatus(model)

			// Find the expected changed item
			var changedItem *types.MCPItem
			for i := range updatedModel.MCPItems {
				if updatedModel.MCPItems[i].Name == tt.expectedChanged {
					changedItem = &updatedModel.MCPItems[i]
					break
				}
			}

			if changedItem == nil {
				t.Fatalf("Could not find expected changed item: %s", tt.expectedChanged)
			}

			// For out of bounds test, verify no change occurred
			if tt.selectedItem >= len(GetFilteredMCPs(model)) {
				if changedItem.Active != originalStates[tt.expectedChanged] {
					t.Errorf("ToggleMCPStatus() with out of bounds selection should not change status")
				}
			} else {
				if changedItem.Active != tt.expectedActive {
					t.Errorf("ToggleMCPStatus() changed %s to %v, expected %v",
						tt.expectedChanged, changedItem.Active, tt.expectedActive)
				}
			}
		})
	}
}

func TestGetActiveMCPCount(t *testing.T) {
	tests := []struct {
		name          string
		mcpItems      []types.MCPItem
		expectedCount int
	}{
		{
			name: "Count with mixed active states",
			mcpItems: []types.MCPItem{
				{Name: "active1", Active: true},
				{Name: "inactive1", Active: false},
				{Name: "active2", Active: true},
				{Name: "inactive2", Active: false},
			},
			expectedCount: 2,
		},
		{
			name: "All active MCPs",
			mcpItems: []types.MCPItem{
				{Name: "active1", Active: true},
				{Name: "active2", Active: true},
			},
			expectedCount: 2,
		},
		{
			name: "No active MCPs",
			mcpItems: []types.MCPItem{
				{Name: "inactive1", Active: false},
				{Name: "inactive2", Active: false},
			},
			expectedCount: 0,
		},
		{
			name:          "Empty MCP list",
			mcpItems:      []types.MCPItem{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := types.Model{MCPItems: tt.mcpItems}
			count := GetActiveMCPCount(model)

			if count != tt.expectedCount {
				t.Errorf("GetActiveMCPCount() = %d, expected %d", count, tt.expectedCount)
			}
		})
	}
}

func TestGetSelectedMCP(t *testing.T) {
	tests := []struct {
		name         string
		selectedItem int
		expectNil    bool
		expectedName string
	}{
		{
			name:         "Valid selection returns correct MCP",
			selectedItem: 0,
			expectNil:    false,
			expectedName: "context7",
		},
		{
			name:         "Last valid index",
			selectedItem: 4, // last item in test data
			expectNil:    false,
			expectedName: "docker-mcp",
		},
		{
			name:         "Out of bounds selection returns nil",
			selectedItem: 10,
			expectNil:    true,
		},
		{
			name:         "Negative selection returns nil",
			selectedItem: -1,
			expectNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestModel()
			model.SelectedItem = tt.selectedItem

			selected := GetSelectedMCP(model)

			if tt.expectNil {
				if selected != nil {
					t.Errorf("GetSelectedMCP() expected nil, got %v", selected)
				}
			} else {
				if selected == nil {
					t.Errorf("GetSelectedMCP() returned nil, expected non-nil")
				} else if selected.Name != tt.expectedName {
					t.Errorf("GetSelectedMCP() returned %s, expected %s", selected.Name, tt.expectedName)
				}
			}
		})
	}
}

// Edge case and error condition tests
func TestMCPServiceEdgeCases(t *testing.T) {
	t.Run("ToggleMCPStatus with empty MCPs list", func(t *testing.T) {
		model := types.Model{
			MCPItems:     []types.MCPItem{},
			SelectedItem: 0,
		}

		result := ToggleMCPStatus(model)

		// Should not panic and return unchanged model
		if len(result.MCPItems) != 0 {
			t.Errorf("ToggleMCPStatus() with empty list should return empty list")
		}
	})

	t.Run("GetFilteredMCPs with special characters in search", func(t *testing.T) {
		model := createTestModel()
		model.SearchQuery = "!@#$%"

		filtered := GetFilteredMCPs(model)

		if len(filtered) != 0 {
			t.Errorf("GetFilteredMCPs() with special characters should return no results")
		}
	})

	t.Run("GetSelectedMCP with empty MCPs list", func(t *testing.T) {
		model := types.Model{
			MCPItems:     []types.MCPItem{},
			SelectedItem: 0,
		}

		selected := GetSelectedMCP(model)

		if selected != nil {
			t.Errorf("GetSelectedMCP() with empty list should return nil")
		}
	})
}

// Boundary condition tests
func TestMCPServiceBoundaryConditions(t *testing.T) {
	t.Run("ToggleMCPStatus at exact boundary", func(t *testing.T) {
		model := createTestModel()
		model.SelectedItem = len(model.MCPItems) - 1 // exactly at last valid index

		originalActive := model.MCPItems[model.SelectedItem].Active
		result := ToggleMCPStatus(model)

		if result.MCPItems[model.SelectedItem].Active == originalActive {
			t.Errorf("ToggleMCPStatus() should toggle status at boundary")
		}
	})

	t.Run("GetSelectedMCP at exact boundary", func(t *testing.T) {
		model := createTestModel()
		model.SelectedItem = len(model.MCPItems) - 1

		selected := GetSelectedMCP(model)

		if selected == nil {
			t.Errorf("GetSelectedMCP() should return valid item at boundary")
		}
	})
}
