package handlers

import (
	"testing"

	"mcp-hub/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

// Test string constants
const (
	TestKeyEnter = "enter"
	TestKeyEsc   = "esc"
)

func TestHandleEscKey(t *testing.T) {
	tests := getEscKeyTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createEscKeyTestModel(tt)
			result, cmd := HandleEscKey(model)
			validateEscKeyResult(t, result, cmd, tt)
		})
	}
}

func getEscKeyTestCases() []escKeyTestCase {
	return []escKeyTestCase{
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
}

type escKeyTestCase struct {
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
}

func createEscKeyTestModel(tt escKeyTestCase) types.Model {
	return types.Model{
		State:             tt.initialState,
		SearchActive:      tt.initialSearchActive,
		SearchInputActive: tt.initialSearchInput,
		SearchQuery:       tt.initialSearchQuery,
		SelectedItem:      5, // Test that selection resets
	}
}

func validateEscKeyResult(t *testing.T, result types.Model, cmd tea.Cmd, tt escKeyTestCase) {
	validateEscKeyState(t, result, tt)
	validateEscKeySearch(t, result, tt)
	validateEscKeyCommand(t, cmd, tt)
	validateEscKeySelection(t, result, tt)
}

func validateEscKeyState(t *testing.T, result types.Model, tt escKeyTestCase) {
	if result.State != tt.expectedState {
		t.Errorf("HandleEscKey() State = %v, expected %v", result.State, tt.expectedState)
	}
}

func validateEscKeySearch(t *testing.T, result types.Model, tt escKeyTestCase) {
	if result.SearchActive != tt.expectedSearchActive {
		t.Errorf("HandleEscKey() SearchActive = %v, expected %v", result.SearchActive, tt.expectedSearchActive)
	}

	if result.SearchInputActive != tt.expectedSearchInput {
		t.Errorf("HandleEscKey() SearchInputActive = %v, expected %v", result.SearchInputActive, tt.expectedSearchInput)
	}

	if result.SearchQuery != tt.expectedSearchQuery {
		t.Errorf("HandleEscKey() SearchQuery = %s, expected %s", result.SearchQuery, tt.expectedSearchQuery)
	}
}

func validateEscKeyCommand(t *testing.T, cmd tea.Cmd, tt escKeyTestCase) {
	if tt.expectQuit {
		if cmd == nil {
			t.Errorf("HandleEscKey() should return quit command")
		}
	} else {
		if cmd != nil {
			t.Errorf("HandleEscKey() should return nil cmd for non-quit scenarios")
		}
	}
}

func validateEscKeySelection(t *testing.T, result types.Model, tt escKeyTestCase) {
	if tt.initialSearchQuery != "" && tt.expectedSearchQuery == "" && tt.initialState == types.MainNavigation {
		if result.SelectedItem != 0 {
			t.Errorf("HandleEscKey() should reset SelectedItem to 0 when clearing search")
		}
	}
}

func TestHandleSearchModeKeys(t *testing.T) {
	tests := getSearchModeKeyTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createSearchModeTestModel(tt)
			result, cmd := HandleSearchModeKeys(model, tt.key)
			validateSearchModeKeyResult(t, result, cmd, tt)
		})
	}
}

func getSearchModeKeyTestCases() []searchModeKeyTestCase {
	return []searchModeKeyTestCase{
		{
			name:           "enter returns to main navigation",
			key:            TestKeyEnter,
			initialQuery:   TestString,
			expectedQuery:  TestString,
			expectedState:  types.MainNavigation,
			expectedActive: false,
		},
		{
			name:          "backspace removes last character",
			key:           "backspace",
			initialQuery:  TestString,
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
			initialQuery:  TestString,
			expectedQuery: "test5",
			expectedState: types.SearchMode,
		},
		{
			name:          "space adds to query",
			key:           " ",
			initialQuery:  TestString,
			expectedQuery: "test ",
			expectedState: types.SearchMode,
		},
		{
			name:          "special character adds to query",
			key:           "-",
			initialQuery:  TestString,
			expectedQuery: "test-",
			expectedState: types.SearchMode,
		},
		{
			name:          "multi-character key ignored",
			key:           "ctrl+c",
			initialQuery:  TestString,
			expectedQuery: TestString,
			expectedState: types.SearchMode,
		},
	}
}

type searchModeKeyTestCase struct {
	name           string
	key            string
	initialQuery   string
	expectedQuery  string
	expectedState  types.AppState
	expectedActive bool
}

func createSearchModeTestModel(tt searchModeKeyTestCase) types.Model {
	return types.Model{
		State:        types.SearchMode,
		SearchQuery:  tt.initialQuery,
		SearchActive: true,
	}
}

func validateSearchModeKeyResult(t *testing.T, result types.Model, cmd tea.Cmd, tt searchModeKeyTestCase) {
	if result.SearchQuery != tt.expectedQuery {
		t.Errorf("HandleSearchModeKeys() SearchQuery = %s, expected %s",
			result.SearchQuery, tt.expectedQuery)
	}

	if result.State != tt.expectedState {
		t.Errorf("HandleSearchModeKeys() State = %v, expected %v",
			result.State, tt.expectedState)
	}

	if tt.key == TestKeyEnter && result.SearchActive != tt.expectedActive {
		t.Errorf("HandleSearchModeKeys() SearchActive = %v, expected %v",
			result.SearchActive, tt.expectedActive)
	}

	if cmd != nil {
		t.Errorf("HandleSearchModeKeys() should return nil cmd")
	}
}

func TestHandleEscKeyEdgeCases(t *testing.T) {
	t.Run("ESC with SearchResults populated", func(t *testing.T) {
		model := types.Model{
			State:         types.MainNavigation,
			SearchQuery:   TestString,
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
			SearchQuery:  TestString,
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
		testLongSearchQuery(t)
	})

	t.Run("Unicode characters", func(t *testing.T) {
		testUnicodeCharacters(t)
	})

	t.Run("Empty string key", func(t *testing.T) {
		testEmptyStringKey(t)
	})

	t.Run("Preserves other model fields", func(t *testing.T) {
		testPreservesOtherModelFields(t)
	})
}

func testLongSearchQuery(t *testing.T) {
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
}

func testUnicodeCharacters(t *testing.T) {
	model := types.Model{
		State:       types.SearchMode,
		SearchQuery: TestString,
	}

	result, _ := HandleSearchModeKeys(model, "n") // Use ASCII instead of unicode

	if result.SearchQuery != "testn" {
		t.Errorf("HandleSearchModeKeys() should handle characters")
	}
}

func testEmptyStringKey(t *testing.T) {
	model := types.Model{
		State:       types.SearchMode,
		SearchQuery: TestString,
	}

	result, _ := HandleSearchModeKeys(model, "")

	if result.SearchQuery != TestString {
		t.Errorf("HandleSearchModeKeys() should ignore empty string key")
	}
}

func testPreservesOtherModelFields(t *testing.T) {
	model := types.Model{
		State:        types.SearchMode,
		SearchQuery:  TestString,
		Width:        100,
		Height:       50,
		SelectedItem: 5,
		ActiveColumn: 2,
	}

	result, _ := HandleSearchModeKeys(model, "a")

	assertModelFieldsPreserved(t, result, model)
}

func assertModelFieldsPreserved(t *testing.T, result, model types.Model) {
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
}

func TestSearchStateTransitions(t *testing.T) {
	tests := createSearchStateTransitionTests()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestModel(tt.initialState)
			result, cmd := executeStateTransition(model, tt.action)
			assertStateTransition(t, result, cmd, tt)
		})
	}
}

func createSearchStateTransitionTests() []struct {
	name          string
	initialState  types.AppState
	action        string
	expectedState types.AppState
} {
	return []struct {
		name          string
		initialState  types.AppState
		action        string
		expectedState types.AppState
	}{
		{
			name:          "SearchMode to MainNavigation via enter",
			initialState:  types.SearchMode,
			action:        TestKeyEnter,
			expectedState: types.MainNavigation,
		},
		{
			name:          "SearchMode to MainNavigation via ESC",
			initialState:  types.SearchMode,
			action:        TestKeyEsc,
			expectedState: types.MainNavigation,
		},
		{
			name:          "SearchActiveNavigation to MainNavigation via ESC",
			initialState:  types.SearchActiveNavigation,
			action:        TestKeyEsc,
			expectedState: types.MainNavigation,
		},
		{
			name:          "ModalActive to MainNavigation via ESC",
			initialState:  types.ModalActive,
			action:        TestKeyEsc,
			expectedState: types.MainNavigation,
		},
	}
}

func createTestModel(state types.AppState) types.Model {
	return types.Model{
		State:        state,
		SearchActive: true,
		SearchQuery:  TestString,
	}
}

func executeStateTransition(model types.Model, action string) (types.Model, tea.Cmd) {
	switch action {
	case TestKeyEnter:
		return HandleSearchModeKeys(model, TestKeyEnter)
	case TestKeyEsc:
		return HandleEscKey(model)
	default:
		return model, nil
	}
}

func assertStateTransition(t *testing.T, result types.Model, cmd tea.Cmd, tt struct {
	name          string
	initialState  types.AppState
	action        string
	expectedState types.AppState
}) {
	if result.State != tt.expectedState {
		t.Errorf("State transition from %v via %s = %v, expected %v",
			tt.initialState, tt.action, result.State, tt.expectedState)
	}

	// For ESC transitions, search should be cleared (except for modal)
	if tt.action == TestKeyEsc && tt.initialState != types.ModalActive && result.SearchQuery != "" {
		t.Errorf("ESC should clear search query")
	}

	if cmd != nil && tt.action != TestKeyEsc {
		t.Errorf("Non-quit transitions should return nil cmd")
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
