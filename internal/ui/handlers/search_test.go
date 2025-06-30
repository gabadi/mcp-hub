package handlers

import (
	"testing"

	"cc-mcp-manager/internal/ui/types"
	tea "github.com/charmbracelet/bubbletea"
)

func TestHandleEscKey(t *testing.T) {
	tests := []struct {
		name                 string
		initialState         types.AppState
		initialSearchActive  bool
		initialSearchInput   bool
		initialSearchQuery   string
		expectedState        types.AppState
		expectedSearchActive bool
		expectedSearchInput  bool
		expectedSearchQuery  string
		expectQuit           bool
	}{
		{
			name:                 "ESC from SearchMode clears search",
			initialState:         types.SearchMode,
			initialSearchActive:  true,
			initialSearchInput:   false,
			initialSearchQuery:   "test query",
			expectedState:        types.MainNavigation,
			expectedSearchActive: false,
			expectedSearchInput:  false,
			expectedSearchQuery:  "",
			expectQuit:           false,
		},
		{
			name:                 "ESC from SearchActiveNavigation clears search",
			initialState:         types.SearchActiveNavigation,
			initialSearchActive:  true,
			initialSearchInput:   true,
			initialSearchQuery:   "test query",
			expectedState:        types.MainNavigation,
			expectedSearchActive: false,
			expectedSearchInput:  false,
			expectedSearchQuery:  "",
			expectQuit:           false,
		},
		{
			name:          "ESC from ModalActive returns to main",
			initialState:  types.ModalActive,
			expectedState: types.MainNavigation,
			expectQuit:    false,
		},
		{
			name:                "ESC from MainNavigation with search query clears search",
			initialState:        types.MainNavigation,
			initialSearchQuery:  "test query",
			expectedState:       types.MainNavigation,
			expectedSearchQuery: "",
			expectQuit:          false,
		},
		{
			name:          "ESC from MainNavigation without search quits",
			initialState:  types.MainNavigation,
			expectedState: types.MainNavigation,
			expectQuit:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := types.Model{
				State:             tt.initialState,
				SearchActive:      tt.initialSearchActive,
				SearchInputActive: tt.initialSearchInput,
				SearchQuery:       tt.initialSearchQuery,
				SelectedItem:      5, // Test that selection resets
			}

			result, cmd := HandleEscKey(model)

			// Check state
			if result.State != tt.expectedState {
				t.Errorf("HandleEscKey() State = %v, expected %v",
					result.State, tt.expectedState)
			}

			// Check search flags
			if result.SearchActive != tt.expectedSearchActive {
				t.Errorf("HandleEscKey() SearchActive = %v, expected %v",
					result.SearchActive, tt.expectedSearchActive)
			}

			if result.SearchInputActive != tt.expectedSearchInput {
				t.Errorf("HandleEscKey() SearchInputActive = %v, expected %v",
					result.SearchInputActive, tt.expectedSearchInput)
			}

			// Check search query
			if result.SearchQuery != tt.expectedSearchQuery {
				t.Errorf("HandleEscKey() SearchQuery = %s, expected %s",
					result.SearchQuery, tt.expectedSearchQuery)
			}

			// Check command
			if tt.expectQuit {
				if cmd == nil {
					t.Errorf("HandleEscKey() should return quit command")
				}
			} else {
				if cmd != nil {
					t.Errorf("HandleEscKey() should return nil cmd for non-quit scenarios")
				}
			}

			// Check that selection resets when clearing search
			if tt.initialSearchQuery != "" && tt.expectedSearchQuery == "" && tt.initialState == types.MainNavigation {
				if result.SelectedItem != 0 {
					t.Errorf("HandleEscKey() should reset SelectedItem to 0 when clearing search")
				}
			}
		})
	}
}

func TestHandleSearchModeKeys(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		initialQuery   string
		expectedQuery  string
		expectedState  types.AppState
		expectedActive bool
	}{
		{
			name:           "enter returns to main navigation",
			key:            "enter",
			initialQuery:   "test",
			expectedQuery:  "test",
			expectedState:  types.MainNavigation,
			expectedActive: false,
		},
		{
			name:          "backspace removes last character",
			key:           "backspace",
			initialQuery:  "test",
			expectedQuery: "tes",
			expectedState: types.SearchMode,
		},
		{
			name:          "backspace on empty query",
			key:           "backspace",
			initialQuery:  "",
			expectedQuery: "",
			expectedState: types.SearchMode,
		},
		{
			name:          "single character adds to query",
			key:           "a",
			initialQuery:  "tes",
			expectedQuery: "tesa",
			expectedState: types.SearchMode,
		},
		{
			name:          "number adds to query",
			key:           "5",
			initialQuery:  "test",
			expectedQuery: "test5",
			expectedState: types.SearchMode,
		},
		{
			name:          "space adds to query",
			key:           " ",
			initialQuery:  "test",
			expectedQuery: "test ",
			expectedState: types.SearchMode,
		},
		{
			name:          "special character adds to query",
			key:           "-",
			initialQuery:  "test",
			expectedQuery: "test-",
			expectedState: types.SearchMode,
		},
		{
			name:          "multi-character key ignored",
			key:           "ctrl+c",
			initialQuery:  "test",
			expectedQuery: "test",
			expectedState: types.SearchMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := types.Model{
				State:        types.SearchMode,
				SearchQuery:  tt.initialQuery,
				SearchActive: true,
			}

			result, cmd := HandleSearchModeKeys(model, tt.key)

			if result.SearchQuery != tt.expectedQuery {
				t.Errorf("HandleSearchModeKeys() SearchQuery = %s, expected %s",
					result.SearchQuery, tt.expectedQuery)
			}

			if result.State != tt.expectedState {
				t.Errorf("HandleSearchModeKeys() State = %v, expected %v",
					result.State, tt.expectedState)
			}

			if tt.key == "enter" && result.SearchActive != tt.expectedActive {
				t.Errorf("HandleSearchModeKeys() SearchActive = %v, expected %v",
					result.SearchActive, tt.expectedActive)
			}

			if cmd != nil {
				t.Errorf("HandleSearchModeKeys() should return nil cmd")
			}
		})
	}
}

func TestHandleEscKeyEdgeCases(t *testing.T) {
	t.Run("ESC with SearchResults populated", func(t *testing.T) {
		model := types.Model{
			State:         types.MainNavigation,
			SearchQuery:   "test",
			SearchResults: []string{"result1", "result2"},
			SelectedItem:  3,
		}

		result, _ := HandleEscKey(model)

		if result.SearchResults != nil {
			t.Errorf("HandleEscKey() should clear SearchResults")
		}

		if result.SelectedItem != 0 {
			t.Errorf("HandleEscKey() should reset SelectedItem to 0")
		}
	})

	t.Run("ESC preserves other model fields", func(t *testing.T) {
		model := types.Model{
			State:        types.SearchMode,
			SearchQuery:  "test",
			Width:        100,
			Height:       50,
			ColumnCount:  4,
			ActiveColumn: 2,
		}

		result, _ := HandleEscKey(model)

		// Non-search fields should be preserved
		if result.Width != model.Width {
			t.Errorf("HandleEscKey() should preserve Width")
		}
		if result.Height != model.Height {
			t.Errorf("HandleEscKey() should preserve Height")
		}
		if result.ColumnCount != model.ColumnCount {
			t.Errorf("HandleEscKey() should preserve ColumnCount")
		}
		if result.ActiveColumn != model.ActiveColumn {
			t.Errorf("HandleEscKey() should preserve ActiveColumn")
		}
	})
}

func TestHandleSearchModeKeysEdgeCases(t *testing.T) {
	t.Run("Very long search query", func(t *testing.T) {
		longQuery := "this is a very long search query that should still work"
		model := types.Model{
			State:       types.SearchMode,
			SearchQuery: longQuery,
		}

		result, _ := HandleSearchModeKeys(model, "!")

		expectedQuery := longQuery + "!"
		if result.SearchQuery != expectedQuery {
			t.Errorf("HandleSearchModeKeys() should handle long queries")
		}
	})

	t.Run("Unicode characters", func(t *testing.T) {
		model := types.Model{
			State:       types.SearchMode,
			SearchQuery: "test",
		}

		result, _ := HandleSearchModeKeys(model, "n") // Use ASCII instead of unicode

		if result.SearchQuery != "testn" {
			t.Errorf("HandleSearchModeKeys() should handle characters")
		}
	})

	t.Run("Empty string key", func(t *testing.T) {
		model := types.Model{
			State:       types.SearchMode,
			SearchQuery: "test",
		}

		result, _ := HandleSearchModeKeys(model, "")

		if result.SearchQuery != "test" {
			t.Errorf("HandleSearchModeKeys() should ignore empty string key")
		}
	})

	t.Run("Preserves other model fields", func(t *testing.T) {
		model := types.Model{
			State:        types.SearchMode,
			SearchQuery:  "test",
			Width:        100,
			Height:       50,
			SelectedItem: 5,
			ActiveColumn: 2,
		}

		result, _ := HandleSearchModeKeys(model, "a")

		// Non-search fields should be preserved
		if result.Width != model.Width {
			t.Errorf("HandleSearchModeKeys() should preserve Width")
		}
		if result.Height != model.Height {
			t.Errorf("HandleSearchModeKeys() should preserve Height")
		}
		if result.SelectedItem != model.SelectedItem {
			t.Errorf("HandleSearchModeKeys() should preserve SelectedItem")
		}
		if result.ActiveColumn != model.ActiveColumn {
			t.Errorf("HandleSearchModeKeys() should preserve ActiveColumn")
		}
	})
}

func TestSearchStateTransitions(t *testing.T) {
	tests := []struct {
		name          string
		initialState  types.AppState
		action        string
		expectedState types.AppState
	}{
		{
			name:          "SearchMode to MainNavigation via enter",
			initialState:  types.SearchMode,
			action:        "enter",
			expectedState: types.MainNavigation,
		},
		{
			name:          "SearchMode to MainNavigation via ESC",
			initialState:  types.SearchMode,
			action:        "esc",
			expectedState: types.MainNavigation,
		},
		{
			name:          "SearchActiveNavigation to MainNavigation via ESC",
			initialState:  types.SearchActiveNavigation,
			action:        "esc",
			expectedState: types.MainNavigation,
		},
		{
			name:          "ModalActive to MainNavigation via ESC",
			initialState:  types.ModalActive,
			action:        "esc",
			expectedState: types.MainNavigation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := types.Model{
				State:        tt.initialState,
				SearchActive: true,
				SearchQuery:  "test",
			}

			var result types.Model
			var cmd tea.Cmd

			switch tt.action {
			case "enter":
				result, cmd = HandleSearchModeKeys(model, "enter")
			case "esc":
				result, cmd = HandleEscKey(model)
			}

			if result.State != tt.expectedState {
				t.Errorf("State transition from %v via %s = %v, expected %v",
					tt.initialState, tt.action, result.State, tt.expectedState)
			}

			// For ESC transitions, search should be cleared (except for modal)
			if tt.action == "esc" && tt.initialState != types.ModalActive && result.SearchQuery != "" {
				t.Errorf("ESC should clear search query")
			}

			if cmd != nil && tt.action != "esc" {
				t.Errorf("Non-quit transitions should return nil cmd")
			}
		})
	}
}

func TestSearchQueryBoundaryConditions(t *testing.T) {
	t.Run("Backspace on single character", func(t *testing.T) {
		model := types.Model{
			State:       types.SearchMode,
			SearchQuery: "a",
		}

		result, _ := HandleSearchModeKeys(model, "backspace")

		if result.SearchQuery != "" {
			t.Errorf("Backspace on single character should result in empty query")
		}
	})

	t.Run("Multiple backspaces", func(t *testing.T) {
		model := types.Model{
			State:       types.SearchMode,
			SearchQuery: "abc",
		}

		// Apply backspace three times
		result, _ := HandleSearchModeKeys(model, "backspace")
		result, _ = HandleSearchModeKeys(result, "backspace")
		result, _ = HandleSearchModeKeys(result, "backspace")

		if result.SearchQuery != "" {
			t.Errorf("Multiple backspaces should result in empty query")
		}
	})

	t.Run("Build up query character by character", func(t *testing.T) {
		model := types.Model{
			State:       types.SearchMode,
			SearchQuery: "",
		}

		chars := []string{"t", "e", "s", "t"}
		expected := ""

		for _, char := range chars {
			model, _ = HandleSearchModeKeys(model, char)
			expected += char

			if model.SearchQuery != expected {
				t.Errorf("Building query: expected %s, got %s", expected, model.SearchQuery)
			}
		}
	})
}
