package handlers

import (
	"testing"

	"cc-mcp-manager/internal/ui/services"
	"cc-mcp-manager/internal/ui/types"
)

// Test helper functions
func createTestNavigationModel() types.Model {
	return types.Model{
		State:        types.MainNavigation,
		ColumnCount:  4,
		ActiveColumn: 0,
		SelectedItem: 0,
		SearchQuery:  "",
		MCPItems: []types.MCPItem{
			{Name: "item1", Active: true},
			{Name: "item2", Active: false},
			{Name: "item3", Active: true},
			{Name: "item4", Active: false},
			{Name: "item5", Active: true},
			{Name: "item6", Active: false},
			{Name: "item7", Active: true},
			{Name: "item8", Active: false},
			{Name: "item9", Active: true},
		},
	}
}

func TestHandleMainNavigationKeys(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		initialModel   types.Model
		expectedState  types.AppState
		expectedActive bool
		expectedInput  bool
	}{
		{
			name:           "Tab activates search with navigation",
			key:            "tab",
			initialModel:   createTestNavigationModel(),
			expectedState:  types.SearchActiveNavigation,
			expectedActive: true,
			expectedInput:  true,
		},
		{
			name:           "Slash activates search with navigation",
			key:            "/",
			initialModel:   createTestNavigationModel(),
			expectedState:  types.SearchActiveNavigation,
			expectedActive: true,
			expectedInput:  true,
		},
		{
			name:          "a key activates modal",
			key:           "a",
			initialModel:  createTestNavigationModel(),
			expectedState: types.ModalActive,
		},
		{
			name:          "e key activates modal",
			key:           "e",
			initialModel:  createTestNavigationModel(),
			expectedState: types.ModalActive,
		},
		{
			name:          "d key activates modal",
			key:           "d",
			initialModel:  createTestNavigationModel(),
			expectedState: types.ModalActive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, cmd := HandleMainNavigationKeys(tt.initialModel, tt.key)

			if result.State != tt.expectedState {
				t.Errorf("HandleMainNavigationKeys() State = %v, expected %v",
					result.State, tt.expectedState)
			}

			if tt.key == "tab" || tt.key == "/" {
				if result.SearchActive != tt.expectedActive {
					t.Errorf("HandleMainNavigationKeys() SearchActive = %v, expected %v",
						result.SearchActive, tt.expectedActive)
				}
				if result.SearchInputActive != tt.expectedInput {
					t.Errorf("HandleMainNavigationKeys() SearchInputActive = %v, expected %v",
						result.SearchInputActive, tt.expectedInput)
				}
				if result.SelectedItem != 0 {
					t.Errorf("HandleMainNavigationKeys() should reset SelectedItem to 0")
				}
			}

			if cmd != nil {
				t.Errorf("HandleMainNavigationKeys() should return nil cmd")
			}
		})
	}
}

func TestHandleMainNavigationMovement(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		initialItem  int
		expectedItem int
	}{
		{
			name:         "up key calls NavigateUp",
			key:          "up",
			initialItem:  4, // Should move to 0 in 4-column grid
			expectedItem: 0,
		},
		{
			name:         "k key calls NavigateUp",
			key:          "k",
			initialItem:  4,
			expectedItem: 0,
		},
		{
			name:         "down key calls NavigateDown",
			key:          "down",
			initialItem:  0, // Should move to 4 in 4-column grid
			expectedItem: 4,
		},
		{
			name:         "j key calls NavigateDown",
			key:          "j",
			initialItem:  0,
			expectedItem: 4,
		},
		{
			name:         "left key calls NavigateLeft",
			key:          "left",
			initialItem:  1, // Should move to 0
			expectedItem: 0,
		},
		{
			name:         "h key calls NavigateLeft",
			key:          "h",
			initialItem:  1,
			expectedItem: 0,
		},
		{
			name:         "right key calls NavigateRight",
			key:          "right",
			initialItem:  0, // Should move to 1
			expectedItem: 1,
		},
		{
			name:         "l key calls NavigateRight",
			key:          "l",
			initialItem:  0,
			expectedItem: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestNavigationModel()
			model.SelectedItem = tt.initialItem

			result, _ := HandleMainNavigationKeys(model, tt.key)

			if result.SelectedItem != tt.expectedItem {
				t.Errorf("HandleMainNavigationKeys() SelectedItem = %d, expected %d",
					result.SelectedItem, tt.expectedItem)
			}
		})
	}
}

func TestHandleSearchNavigationKeys(t *testing.T) {
	tests := []struct {
		name                 string
		key                  string
		initialInputActive   bool
		expectedState        types.AppState
		expectedInputActive  bool
		expectedSearchActive bool
	}{
		{
			name:                 "enter returns to main navigation",
			key:                  "enter",
			initialInputActive:   true,
			expectedState:        types.MainNavigation,
			expectedInputActive:  false,
			expectedSearchActive: false,
		},
		{
			name:                 "tab toggles input active true to false",
			key:                  "tab",
			initialInputActive:   true,
			expectedState:        types.SearchActiveNavigation,
			expectedInputActive:  false,
			expectedSearchActive: true,
		},
		{
			name:                 "tab toggles input active false to true",
			key:                  "tab",
			initialInputActive:   false,
			expectedState:        types.SearchActiveNavigation,
			expectedInputActive:  true,
			expectedSearchActive: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestNavigationModel()
			model.State = types.SearchActiveNavigation
			model.SearchActive = true
			model.SearchInputActive = tt.initialInputActive

			result, cmd := HandleSearchNavigationKeys(model, tt.key)

			if result.State != tt.expectedState {
				t.Errorf("HandleSearchNavigationKeys() State = %v, expected %v",
					result.State, tt.expectedState)
			}

			if result.SearchInputActive != tt.expectedInputActive {
				t.Errorf("HandleSearchNavigationKeys() SearchInputActive = %v, expected %v",
					result.SearchInputActive, tt.expectedInputActive)
			}

			if result.SearchActive != tt.expectedSearchActive {
				t.Errorf("HandleSearchNavigationKeys() SearchActive = %v, expected %v",
					result.SearchActive, tt.expectedSearchActive)
			}

			if cmd != nil {
				t.Errorf("HandleSearchNavigationKeys() should return nil cmd")
			}
		})
	}
}

func TestHandleSearchNavigationTextInput(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		initialQuery  string
		inputActive   bool
		expectedQuery string
	}{
		{
			name:          "add character when input active",
			key:           "a",
			initialQuery:  "test",
			inputActive:   true,
			expectedQuery: "testa",
		},
		{
			name:          "ignore character when input inactive",
			key:           "a",
			initialQuery:  "test",
			inputActive:   false,
			expectedQuery: "test",
		},
		{
			name:          "backspace removes character",
			key:           "backspace",
			initialQuery:  "test",
			inputActive:   true,
			expectedQuery: "tes",
		},
		{
			name:          "backspace on empty query",
			key:           "backspace",
			initialQuery:  "",
			inputActive:   true,
			expectedQuery: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestNavigationModel()
			model.State = types.SearchActiveNavigation
			model.SearchInputActive = tt.inputActive
			model.SearchQuery = tt.initialQuery

			result, _ := HandleSearchNavigationKeys(model, tt.key)

			if result.SearchQuery != tt.expectedQuery {
				t.Errorf("HandleSearchNavigationKeys() SearchQuery = %s, expected %s",
					result.SearchQuery, tt.expectedQuery)
			}
		})
	}
}

func TestNavigateUp(t *testing.T) {
	tests := []struct {
		name         string
		columnCount  int
		activeColumn int
		initialItem  int
		expectedItem int
		mcpCount     int
	}{
		{
			name:         "4-column grid - move up one row",
			columnCount:  4,
			activeColumn: 0,
			initialItem:  4, // second row, first column
			expectedItem: 0, // first row, first column
			mcpCount:     9,
		},
		{
			name:         "4-column grid - already at top row",
			columnCount:  4,
			activeColumn: 0,
			initialItem:  2, // first row
			expectedItem: 2, // should stay the same since it's already in top row
			mcpCount:     9,
		},
		{
			name:         "4-column grid - out of bounds selection",
			columnCount:  4,
			activeColumn: 0,
			initialItem:  15, // way out of bounds
			expectedItem: 8,  // should clamp to last item
			mcpCount:     9,
		},
		{
			name:         "2-column layout - MCP column navigation",
			columnCount:  2,
			activeColumn: 0,
			initialItem:  3,
			expectedItem: 2,
			mcpCount:     9,
		},
		{
			name:         "2-column layout - non-MCP column",
			columnCount:  2,
			activeColumn: 1,
			initialItem:  3,
			expectedItem: 3, // no change
			mcpCount:     9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestNavigationModel()
			model.ColumnCount = tt.columnCount
			model.ActiveColumn = tt.activeColumn
			model.SelectedItem = tt.initialItem
			// Adjust MCP count if needed
			if tt.mcpCount < len(model.MCPItems) {
				model.MCPItems = model.MCPItems[:tt.mcpCount]
			}

			result := NavigateUp(model)

			if result.SelectedItem != tt.expectedItem {
				t.Errorf("NavigateUp() SelectedItem = %d, expected %d",
					result.SelectedItem, tt.expectedItem)
			}
		})
	}
}

func TestNavigateDown(t *testing.T) {
	tests := []struct {
		name         string
		columnCount  int
		activeColumn int
		initialItem  int
		expectedItem int
		mcpCount     int
	}{
		{
			name:         "4-column grid - move down one row",
			columnCount:  4,
			activeColumn: 0,
			initialItem:  0, // first row, first column
			expectedItem: 4, // second row, first column
			mcpCount:     9,
		},
		{
			name:         "4-column grid - already at bottom",
			columnCount:  4,
			activeColumn: 0,
			initialItem:  8, // last item
			expectedItem: 8, // should stay at last
			mcpCount:     9,
		},
		{
			name:         "4-column grid - would exceed bounds",
			columnCount:  4,
			activeColumn: 0,
			initialItem:  5, // would go to 9, but only 9 items (0-8)
			expectedItem: 5, // should stay at current position since moving down would exceed bounds
			mcpCount:     9,
		},
		{
			name:         "2-column layout - MCP column navigation",
			columnCount:  2,
			activeColumn: 0,
			initialItem:  3,
			expectedItem: 4,
			mcpCount:     9,
		},
		{
			name:         "2-column layout - at last item",
			columnCount:  2,
			activeColumn: 0,
			initialItem:  8, // last item
			expectedItem: 8, // should stay
			mcpCount:     9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestNavigationModel()
			model.ColumnCount = tt.columnCount
			model.ActiveColumn = tt.activeColumn
			model.SelectedItem = tt.initialItem
			// Adjust MCP count if needed
			if tt.mcpCount < len(model.MCPItems) {
				model.MCPItems = model.MCPItems[:tt.mcpCount]
			}

			result := NavigateDown(model)

			if result.SelectedItem != tt.expectedItem {
				t.Errorf("NavigateDown() SelectedItem = %d, expected %d",
					result.SelectedItem, tt.expectedItem)
			}
		})
	}
}

func TestNavigateLeft(t *testing.T) {
	tests := []struct {
		name           string
		columnCount    int
		activeColumn   int
		initialItem    int
		expectedItem   int
		expectedColumn int
	}{
		{
			name:           "4-column grid - move left within row",
			columnCount:    4,
			activeColumn:   0,
			initialItem:    1, // first row, second column
			expectedItem:   0, // first row, first column
			expectedColumn: 0,
		},
		{
			name:           "4-column grid - already at leftmost",
			columnCount:    4,
			activeColumn:   0,
			initialItem:    0, // already at leftmost
			expectedItem:   0, // should stay
			expectedColumn: 0,
		},
		{
			name:           "2-column layout - move between columns",
			columnCount:    2,
			activeColumn:   1,
			initialItem:    5,
			expectedItem:   5, // item stays same
			expectedColumn: 0, // column changes
		},
		{
			name:           "2-column layout - already at leftmost column",
			columnCount:    2,
			activeColumn:   0,
			initialItem:    5,
			expectedItem:   5, // no change
			expectedColumn: 0, // no change
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestNavigationModel()
			model.ColumnCount = tt.columnCount
			model.ActiveColumn = tt.activeColumn
			model.SelectedItem = tt.initialItem

			result := NavigateLeft(model)

			if result.SelectedItem != tt.expectedItem {
				t.Errorf("NavigateLeft() SelectedItem = %d, expected %d",
					result.SelectedItem, tt.expectedItem)
			}

			if result.ActiveColumn != tt.expectedColumn {
				t.Errorf("NavigateLeft() ActiveColumn = %d, expected %d",
					result.ActiveColumn, tt.expectedColumn)
			}
		})
	}
}

func TestNavigateRight(t *testing.T) {
	tests := []struct {
		name           string
		columnCount    int
		activeColumn   int
		initialItem    int
		expectedItem   int
		expectedColumn int
		mcpCount       int
	}{
		{
			name:           "4-column grid - move right within row",
			columnCount:    4,
			activeColumn:   0,
			initialItem:    0, // first row, first column
			expectedItem:   1, // first row, second column
			expectedColumn: 0,
			mcpCount:       9,
		},
		{
			name:           "4-column grid - at rightmost in row",
			columnCount:    4,
			activeColumn:   0,
			initialItem:    3, // first row, last column
			expectedItem:   3, // should stay
			expectedColumn: 0,
			mcpCount:       9,
		},
		{
			name:           "4-column grid - would exceed total items",
			columnCount:    4,
			activeColumn:   0,
			initialItem:    8, // last item (index 8 of 9 items)
			expectedItem:   8, // should stay
			expectedColumn: 0,
			mcpCount:       9,
		},
		{
			name:           "2-column layout - move between columns",
			columnCount:    2,
			activeColumn:   0,
			initialItem:    5,
			expectedItem:   5, // item stays same
			expectedColumn: 1, // column changes
			mcpCount:       9,
		},
		{
			name:           "2-column layout - already at rightmost column",
			columnCount:    2,
			activeColumn:   1,
			initialItem:    5,
			expectedItem:   5, // no change
			expectedColumn: 1, // no change
			mcpCount:       9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestNavigationModel()
			model.ColumnCount = tt.columnCount
			model.ActiveColumn = tt.activeColumn
			model.SelectedItem = tt.initialItem
			// Adjust MCP count if needed
			if tt.mcpCount < len(model.MCPItems) {
				model.MCPItems = model.MCPItems[:tt.mcpCount]
			}

			result := NavigateRight(model)

			if result.SelectedItem != tt.expectedItem {
				t.Errorf("NavigateRight() SelectedItem = %d, expected %d",
					result.SelectedItem, tt.expectedItem)
			}

			if result.ActiveColumn != tt.expectedColumn {
				t.Errorf("NavigateRight() ActiveColumn = %d, expected %d",
					result.ActiveColumn, tt.expectedColumn)
			}
		})
	}
}

func TestNavigationWithFilteredResults(t *testing.T) {
	model := createTestNavigationModel()
	model.SearchQuery = "item1" // This should filter to just "item1"
	model.ColumnCount = 4
	model.SelectedItem = 0

	// Test that navigation respects filtered results
	filtered := services.GetFilteredMCPs(model)
	if len(filtered) != 1 {
		t.Fatalf("Expected 1 filtered result, got %d", len(filtered))
	}

	// Test NavigateDown with filtered results
	result := NavigateDown(model)
	// Should stay at 0 since there's only 1 filtered item
	if result.SelectedItem != 0 {
		t.Errorf("NavigateDown() with filtered results: SelectedItem = %d, expected 0",
			result.SelectedItem)
	}

	// Test NavigateRight with filtered results
	result = NavigateRight(model)
	// Should stay at 0 since there's only 1 filtered item
	if result.SelectedItem != 0 {
		t.Errorf("NavigateRight() with filtered results: SelectedItem = %d, expected 0",
			result.SelectedItem)
	}
}

func TestNavigationBoundaryConditions(t *testing.T) {
	t.Run("Empty MCP list", func(t *testing.T) {
		model := types.Model{
			ColumnCount:  4,
			SelectedItem: 0,
			MCPItems:     []types.MCPItem{},
		}

		// Navigation should handle empty list gracefully
		result := NavigateDown(model)
		if result.SelectedItem != 0 {
			t.Errorf("NavigateDown() with empty list should keep SelectedItem at 0")
		}

		result = NavigateUp(model)
		if result.SelectedItem != 0 {
			t.Errorf("NavigateUp() with empty list should keep SelectedItem at 0")
		}
	})

	t.Run("Single item navigation", func(t *testing.T) {
		model := types.Model{
			ColumnCount:  4,
			SelectedItem: 0,
			MCPItems:     []types.MCPItem{{Name: "single", Active: true}},
		}

		// All navigation should stay at item 0
		directions := []func(types.Model) types.Model{
			NavigateUp, NavigateDown, NavigateLeft, NavigateRight,
		}

		for i, navFunc := range directions {
			result := navFunc(model)
			if result.SelectedItem != 0 {
				t.Errorf("Navigation function %d with single item should stay at 0, got %d",
					i, result.SelectedItem)
			}
		}
	})
}

func TestPasteToSearchQuery(t *testing.T) {
	model := types.Model{
		SearchQuery: "existing",
	}

	// Test paste functionality - this will depend on clipboard availability
	result := pasteToSearchQuery(model)

	// The function should handle errors gracefully and return a model
	// We can't test actual clipboard content, but we can verify it doesn't crash
	if result.SearchQuery == "" {
		// If clipboard failed, original query should be preserved with error message
		if result.SuccessMessage == "" {
			t.Error("Expected error message when clipboard paste fails")
		}
		if result.SuccessTimer == 0 {
			t.Error("Expected success timer to be set for error message")
		}
	}
}

func TestRefreshClaudeStatusCmdNavigation(t *testing.T) {
	cmd := RefreshClaudeStatusCmd()
	if cmd == nil {
		t.Error("RefreshClaudeStatusCmd() should return a command")
	}

	// Execute the command and verify it returns a ClaudeStatusMsg
	msg := cmd()
	if _, ok := msg.(ClaudeStatusMsg); !ok {
		t.Errorf("RefreshClaudeStatusCmd() should return ClaudeStatusMsg, got %T", msg)
	}
}
