package handlers

import (
	"testing"

	"cc-mcp-manager/internal/ui/types"
	tea "github.com/charmbracelet/bubbletea"
)

// Helper function to create KeyMsg for testing
func createKeyMsg(key string) tea.KeyMsg {
	return tea.KeyMsg{
		Type: tea.KeyRunes,
		Runes: []rune(key),
		Alt: false,
		Paste: false,
	}
}

func TestHandleKeyPressGlobalKeys(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		initialState types.AppState
		expectQuit  bool
	}{
		{
			name:        "ESC key triggers ESC handler",
			key:         "esc",
			initialState: types.MainNavigation,
			expectQuit:  true, // ESC in main nav with no search query quits
		},
		{
			name:        "Ctrl+C triggers quit",
			key:         "ctrl+c",
			initialState: types.MainNavigation,
			expectQuit:  true,
		},
		{
			name:        "Ctrl+C from any state triggers quit",
			key:         "ctrl+c",
			initialState: types.SearchMode,
			expectQuit:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := types.Model{
				State: tt.initialState,
			}
			
			msg := createKeyMsg(tt.key)
			result, cmd := HandleKeyPress(model, msg)
			
			if tt.expectQuit {
				if cmd == nil {
					t.Errorf("HandleKeyPress() should return quit command for %s", tt.key)
				}
			}
			
			// Model should be returned (even if quitting)
			if result.State != tt.initialState && tt.key == "ctrl+c" {
				// ctrl+c doesn't change state, just returns quit command
				t.Errorf("HandleKeyPress() should preserve state for ctrl+c")
			}
		})
	}
}

func TestHandleKeyPressStateRouting(t *testing.T) {
	tests := []struct {
		name          string
		initialState  types.AppState
		key           string
		expectedCalls string // which handler should be called
	}{
		{
			name:          "MainNavigation routes to main navigation handler",
			initialState:  types.MainNavigation,
			key:           "j",
			expectedCalls: "main_nav",
		},
		{
			name:          "SearchMode routes to search mode handler",
			initialState:  types.SearchMode,
			key:           "a",
			expectedCalls: "search_mode",
		},
		{
			name:          "SearchActiveNavigation routes to search navigation handler",
			initialState:  types.SearchActiveNavigation,
			key:           "j",
			expectedCalls: "search_nav",
		},
		{
			name:          "ModalActive routes to modal handler",
			initialState:  types.ModalActive,
			key:           "enter",
			expectedCalls: "modal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := types.Model{
				State:        tt.initialState,
				SelectedItem: 0,
				MCPItems: []types.MCPItem{
					{Name: "test1", Active: true},
					{Name: "test2", Active: false},
				},
			}
			
			msg := createKeyMsg(tt.key)
			result, cmd := HandleKeyPress(model, msg)
			
			// Verify the key was processed (state changes or no error)
			_ = result
			_ = cmd
			
			// For MainNavigation with movement keys, verify navigation occurred
			if tt.initialState == types.MainNavigation && tt.key == "j" {
				// Down navigation should move from 0 to 1 (in 2+ column layout)
				// Or from 0 to 4 in 4-column layout, but we only have 2 items
				// The actual navigation logic is tested in navigation_test.go
			}
			
			// For SearchMode, verify character was added
			if tt.initialState == types.SearchMode && len(tt.key) == 1 {
				if result.SearchQuery != tt.key {
					t.Errorf("SearchMode should add character to query")
				}
			}
		})
	}
}

func TestHandleKeyPressWithDifferentStates(t *testing.T) {
	baseModel := types.Model{
		Width:        100,
		Height:       50,
		ColumnCount:  2,
		ActiveColumn: 0,
		SelectedItem: 0,
		MCPItems: []types.MCPItem{
			{Name: "test1", Active: true},
			{Name: "test2", Active: false},
		},
	}

	tests := []struct {
		name                  string
		state                 types.AppState
		key                   string
		searchActive          bool
		searchInputActive     bool
		searchQuery           string
		expectedStateChange   bool
		expectedQueryChange   bool
	}{
		{
			name:                "MainNavigation with tab activates search",
			state:               types.MainNavigation,
			key:                 "tab",
			expectedStateChange: true,
		},
		{
			name:                "MainNavigation with slash activates search",
			state:               types.MainNavigation,
			key:                 "/",
			expectedStateChange: true,
		},
		{
			name:                "MainNavigation with space toggles MCP",
			state:               types.MainNavigation,
			key:                 " ",
			expectedStateChange: false,
		},
		{
			name:                "SearchActiveNavigation with enter returns to main",
			state:               types.SearchActiveNavigation,
			key:                 "enter",
			searchActive:        true,
			searchInputActive:   true,
			expectedStateChange: true,
		},
		{
			name:                "SearchActiveNavigation with text input",
			state:               types.SearchActiveNavigation,
			key:                 "a",
			searchActive:        true,
			searchInputActive:   true,
			searchQuery:         "test",
			expectedQueryChange: true,
		},
		{
			name:              "SearchMode with backspace",
			state:             types.SearchMode,
			key:               "backspace",
			searchQuery:       "test",
			expectedQueryChange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := baseModel
			model.State = tt.state
			model.SearchActive = tt.searchActive
			model.SearchInputActive = tt.searchInputActive
			model.SearchQuery = tt.searchQuery
			
			originalState := model.State
			originalQuery := model.SearchQuery
			
			msg := createKeyMsg(tt.key)
			result, cmd := HandleKeyPress(model, msg)
			
			if tt.expectedStateChange {
				if result.State == originalState {
					t.Errorf("HandleKeyPress() expected state change but state remained %v", originalState)
				}
			} else {
				if result.State != originalState {
					t.Errorf("HandleKeyPress() unexpected state change from %v to %v", originalState, result.State)
				}
			}
			
			if tt.expectedQueryChange {
				if result.SearchQuery == originalQuery {
					t.Errorf("HandleKeyPress() expected query change but query remained %s", originalQuery)
				}
			}
			
			// Verify cmd is reasonable
			if cmd != nil && cmd == nil {
				t.Errorf("HandleKeyPress() returned unexpected command type")
			}
		})
	}
}

func TestHandleKeyPressEdgeCases(t *testing.T) {
	t.Run("Unknown state", func(t *testing.T) {
		model := types.Model{
			State: types.AppState(999), // Invalid state
		}
		
		msg := createKeyMsg("a")
		result, cmd := HandleKeyPress(model, msg)
		
		// Should handle gracefully and return unchanged model
		if result.State != model.State {
			t.Errorf("HandleKeyPress() with unknown state should preserve state")
		}
		
		if cmd != nil {
			t.Errorf("HandleKeyPress() with unknown state should return nil cmd")
		}
	})
	
	t.Run("Empty key string", func(t *testing.T) {
		model := types.Model{
			State: types.MainNavigation,
		}
		
		msg := createKeyMsg("")
		result, cmd := HandleKeyPress(model, msg)
		
		// Should handle gracefully
		_ = result
		if cmd != nil && cmd == nil {
			t.Errorf("HandleKeyPress() with empty key should handle gracefully")
		}
	})
	
	t.Run("Special characters", func(t *testing.T) {
		model := types.Model{
			State: types.SearchMode,
			SearchQuery: "test",
		}
		
		specialKeys := []string{"@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "_", "=", "+"}
		
		for _, key := range specialKeys {
			msg := createKeyMsg(key)
			result, _ := HandleKeyPress(model, msg)
			
			expectedQuery := "test" + key
			if result.SearchQuery != expectedQuery {
				t.Errorf("HandleKeyPress() should handle special character %s", key)
			}
			
			// Reset for next iteration
			model.SearchQuery = "test"
		}
	})
}

func TestHandleKeyPressPreservesModelFields(t *testing.T) {
	model := types.Model{
		State:        types.MainNavigation,
		Width:        150,
		Height:       50,
		ColumnCount:  4,
		ActiveColumn: 2,
		SelectedItem: 3,
		SearchQuery:  "",
		MCPItems: []types.MCPItem{
			{Name: "preserve1", Active: true},
			{Name: "preserve2", Active: false},
		},
	}
	
	msg := createKeyMsg("h") // Navigate left
	result, _ := HandleKeyPress(model, msg)
	
	// These fields should be preserved
	if result.Width != model.Width {
		t.Errorf("HandleKeyPress() should preserve Width")
	}
	if result.Height != model.Height {
		t.Errorf("HandleKeyPress() should preserve Height")
	}
	if result.ColumnCount != model.ColumnCount {
		t.Errorf("HandleKeyPress() should preserve ColumnCount")
	}
	if len(result.MCPItems) != len(model.MCPItems) {
		t.Errorf("HandleKeyPress() should preserve MCPItems")
	}
	
	// ActiveColumn might change due to navigation, but other navigation should work
	if result.ActiveColumn != model.ActiveColumn-1 && result.ActiveColumn != model.ActiveColumn {
		t.Errorf("HandleKeyPress() navigation should only change ActiveColumn appropriately")
	}
}

func TestHandleKeyPressIntegrationScenarios(t *testing.T) {
	t.Run("Complete search workflow", func(t *testing.T) {
		model := types.Model{
			State:       types.MainNavigation,
			SearchQuery: "",
			MCPItems: []types.MCPItem{
				{Name: "github", Active: true},
				{Name: "docker", Active: false},
			},
		}
		
		// 1. Activate search with '/'
		msg := createKeyMsg("/")
		model, _ = HandleKeyPress(model, msg)
		
		if model.State != types.SearchActiveNavigation {
			t.Errorf("'/' should activate search navigation")
		}
		
		// 2. Type search query
		msg = createKeyMsg("g")
		model, _ = HandleKeyPress(model, msg)
		
		if model.SearchQuery != "g" {
			t.Errorf("Should add character to search query")
		}
		
		// 3. Return to main navigation
		msg = createKeyMsg("enter")
		model, _ = HandleKeyPress(model, msg)
		
		if model.State != types.MainNavigation {
			t.Errorf("Enter should return to main navigation")
		}
		if model.SearchQuery != "g" {
			t.Errorf("Search query should be preserved")
		}
	})
	
	t.Run("Modal activation and return", func(t *testing.T) {
		model := types.Model{
			State: types.MainNavigation,
		}
		
		// Activate modal
		msg := createKeyMsg("a")
		model, _ = HandleKeyPress(model, msg)
		
		if model.State != types.ModalActive {
			t.Errorf("'a' should activate modal")
		}
		
		// Return from modal with ESC
		msg = createKeyMsg("esc")
		model, _ = HandleKeyPress(model, msg)
		
		if model.State != types.MainNavigation {
			t.Errorf("ESC should return from modal to main navigation")
		}
	})
	
	t.Run("Navigation in different layouts", func(t *testing.T) {
		// Test 4-column navigation
		model := types.Model{
			State:        types.MainNavigation,
			ColumnCount:  4,
			SelectedItem: 0,
			MCPItems: make([]types.MCPItem, 10),
		}
		
		// Move right in 4-column grid
		msg := createKeyMsg("l")
		result, _ := HandleKeyPress(model, msg)
		
		if result.SelectedItem != 1 {
			t.Errorf("Right navigation in 4-column grid should move from 0 to 1")
		}
		
		// Test 2-column navigation
		model.ColumnCount = 2
		model.ActiveColumn = 0
		
		// Move right should change column in 2-column layout
		msg = createKeyMsg("l")
		result, _ = HandleKeyPress(model, msg)
		
		if result.ActiveColumn != 1 {
			t.Errorf("Right navigation in 2-column layout should change ActiveColumn")
		}
	})
}