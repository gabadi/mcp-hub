package services

import (
	"testing"

	"mcp-hub/internal/platform"
	"mcp-hub/internal/ui/types"
)

// Test helpers for MCP service tests
func createTestModel() types.Model {
	return types.Model{
		SearchQuery:     "",
		SelectedItem:    0,
		PlatformService: platform.GetMockPlatformService(),
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
	// Note: Enhanced ToggleMCPStatus only sets loading state, doesn't actually toggle immediately
	tests := []struct {
		name            string
		selectedItem    int
		searchQuery     string
		claudeAvailable bool
		expectedState   types.ToggleOperationState
		expectedMCPName string
	}{
		{
			name:            "Toggle with Claude available - should set loading state",
			selectedItem:    0, // context7 (active)
			searchQuery:     "",
			claudeAvailable: true,
			expectedState:   types.ToggleLoading,
			expectedMCPName: "context7",
		},
		{
			name:            "Toggle with Claude unavailable - should error",
			selectedItem:    2, // ht-mcp (inactive)
			searchQuery:     "",
			claudeAvailable: false,
			expectedState:   types.ToggleError,
			expectedMCPName: "ht-mcp",
		},
		{
			name:            "Toggle with filtered results",
			selectedItem:    0, // first in filtered results
			searchQuery:     "github",
			claudeAvailable: true,
			expectedState:   types.ToggleLoading,
			expectedMCPName: "github-mcp",
		},
		{
			name:            "Toggle with out of bounds selection",
			selectedItem:    10, // beyond array bounds
			searchQuery:     "",
			claudeAvailable: true,
			expectedState:   types.ToggleIdle, // should remain idle
			expectedMCPName: "",               // no MCP selected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestModel()
			model.SelectedItem = tt.selectedItem
			model.SearchQuery = tt.searchQuery
			model.ClaudeAvailable = tt.claudeAvailable
			model.ToggleState = types.ToggleIdle // Start with idle state

			updatedModel := ToggleMCPStatus(model, platform.GetMockPlatformService())

			if updatedModel.ToggleState != tt.expectedState {
				t.Errorf("ToggleMCPStatus() state = %v, expected %v", updatedModel.ToggleState, tt.expectedState)
			}

			if updatedModel.ToggleMCPName != tt.expectedMCPName {
				t.Errorf("ToggleMCPStatus() MCPName = %q, expected %q", updatedModel.ToggleMCPName, tt.expectedMCPName)
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
			model := types.Model{
				MCPItems:        tt.mcpItems,
				PlatformService: platform.GetMockPlatformService(),
			}
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
			MCPItems:        []types.MCPItem{},
			SelectedItem:    0,
			PlatformService: platform.GetMockPlatformService(),
		}

		result := ToggleMCPStatus(model, platform.GetMockPlatformService())

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
			MCPItems:        []types.MCPItem{},
			SelectedItem:    0,
			PlatformService: platform.GetMockPlatformService(),
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
		model.ClaudeAvailable = true
		model.ToggleState = types.ToggleIdle

		result := ToggleMCPStatus(model, platform.GetMockPlatformService())

		// Enhanced toggle should set loading state, not immediately toggle
		if result.ToggleState != types.ToggleLoading {
			t.Errorf("ToggleMCPStatus() should set loading state at boundary")
		}

		expectedMCPName := model.MCPItems[model.SelectedItem].Name
		if result.ToggleMCPName != expectedMCPName {
			t.Errorf("ToggleMCPStatus() should set correct MCP name at boundary")
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

// Epic 2 Story 2 Tests - Enhanced Toggle Functionality

func createTestModelWithClaudeStatus(claudeAvailable bool) types.Model {
	model := createTestModel()
	model.ClaudeAvailable = claudeAvailable
	model.ToggleState = types.ToggleIdle
	return model
}

func TestToggleMCPStatusEnhanced(t *testing.T) {
	tests := getToggleMCPStatusTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validateToggleMCPStatusTestCase(t, tt)
		})
	}
}

func getToggleMCPStatusTestCases() []struct {
	name            string
	claudeAvailable bool
	selectedItem    int
	expectedState   types.ToggleOperationState
	expectError     bool
} {
	return []struct {
		name            string
		claudeAvailable bool
		selectedItem    int
		expectedState   types.ToggleOperationState
		expectError     bool
	}{
		{
			name:            "Claude available - should start loading",
			claudeAvailable: true,
			selectedItem:    0,
			expectedState:   types.ToggleLoading,
			expectError:     false,
		},
		{
			name:            "Claude unavailable - should error immediately",
			claudeAvailable: false,
			selectedItem:    0,
			expectedState:   types.ToggleError,
			expectError:     true,
		},
		{
			name:            "No MCP selected - no change",
			claudeAvailable: true,
			selectedItem:    -1,
			expectedState:   types.ToggleIdle,
			expectError:     false,
		},
		{
			name:            "Out of bounds selection - no change",
			claudeAvailable: true,
			selectedItem:    100,
			expectedState:   types.ToggleIdle,
			expectError:     false,
		},
	}
}

func validateToggleMCPStatusTestCase(t *testing.T, tt struct {
	name            string
	claudeAvailable bool
	selectedItem    int
	expectedState   types.ToggleOperationState
	expectError     bool
}) {
	model := createTestModelWithClaudeStatus(tt.claudeAvailable)
	model.SelectedItem = tt.selectedItem

	result := ToggleMCPStatus(model, platform.GetMockPlatformService())

	if result.ToggleState != tt.expectedState {
		t.Errorf("ToggleMCPStatus() state = %v, expected %v", result.ToggleState, tt.expectedState)
	}

	if tt.expectError && result.ToggleError == "" {
		t.Error("ToggleMCPStatus() should set error message when Claude unavailable")
	}

	if !tt.expectError && result.ToggleError != "" && result.ToggleState != types.ToggleIdle {
		t.Errorf("ToggleMCPStatus() should not set error when Claude available, got: %s", result.ToggleError)
	}

	// Check MCP name is set when appropriate
	if tt.expectedState == types.ToggleLoading || tt.expectedState == types.ToggleError {
		if tt.selectedItem >= 0 && tt.selectedItem < len(model.MCPItems) {
			expectedMCPName := model.MCPItems[tt.selectedItem].Name
			if result.ToggleMCPName != expectedMCPName {
				t.Errorf("ToggleMCPStatus() MCPName = %s, expected %s", result.ToggleMCPName, expectedMCPName)
			}
		}
	}
}

func TestEnhancedToggleMCPStatus(t *testing.T) {
	tests := []struct {
		name     string
		mcpName  string
		activate bool
	}{
		{
			name:     "Activate MCP",
			mcpName:  "test-mcp",
			activate: true,
		},
		{
			name:     "Deactivate MCP",
			mcpName:  "test-mcp",
			activate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestModelWithClaudeStatus(true)
			// Add the test MCP to the model
			model.MCPItems = append(model.MCPItems, types.MCPItem{
				Name:   tt.mcpName,
				Active: !tt.activate, // Set opposite of what we want to test toggle
			})

			result := EnhancedToggleMCPStatus(model, tt.mcpName, tt.activate)

			// The result depends on whether Claude CLI is actually available
			// Since we can't easily mock exec.Command, we check the structure
			if result.ToggleMCPName != tt.mcpName {
				t.Errorf("EnhancedToggleMCPStatus() MCPName = %s, expected %s", result.ToggleMCPName, tt.mcpName)
			}

			// Toggle state should be success, error, or retrying
			validStates := []types.ToggleOperationState{
				types.ToggleSuccess,
				types.ToggleError,
				types.ToggleRetrying,
			}
			stateValid := false
			for _, state := range validStates {
				if result.ToggleState == state {
					stateValid = true
					break
				}
			}
			if !stateValid {
				t.Errorf("EnhancedToggleMCPStatus() state = %v, expected one of %v", result.ToggleState, validStates)
			}
		})
	}
}

func TestLegacyToggleMCPStatus(t *testing.T) {
	// Test that legacy function still works for backward compatibility
	model := createTestModel()
	originalActive := model.MCPItems[0].Active

	result := LegacyToggleMCPStatus(model, platform.GetMockPlatformService())

	// Should toggle the first MCP
	if result.MCPItems[0].Active == originalActive {
		t.Error("LegacyToggleMCPStatus() should toggle MCP status")
	}
}

// Benchmark tests for enhanced functionality
func BenchmarkToggleMCPStatusEnhanced(b *testing.B) {
	model := createTestModelWithClaudeStatus(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ToggleMCPStatus(model, platform.GetMockPlatformService())
	}
}

func BenchmarkEnhancedToggleMCPStatus(b *testing.B) {
	model := createTestModelWithClaudeStatus(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EnhancedToggleMCPStatus(model, "test-mcp", true)
	}
}
