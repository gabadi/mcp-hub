package ui

import (
	"testing"

	"cc-mcp-manager/internal/testutil"
	"cc-mcp-manager/internal/ui/types"
	tea "github.com/charmbracelet/bubbletea"
)

func TestModel_GetterMethods(t *testing.T) {
	t.Run("GetColumnCount", func(t *testing.T) {
		tests := []struct {
			name        string
			columnCount int
		}{
			{"1 column", 1},
			{"2 columns", 2},
			{"3 columns", 3},
			{"4 columns", 4},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				model := NewModel()
				model.ColumnCount = tt.columnCount

				result := model.GetColumnCount()

				if result != tt.columnCount {
					t.Errorf("GetColumnCount() = %d, expected %d", result, tt.columnCount)
				}
			})
		}
	})

	t.Run("GetActiveColumn", func(t *testing.T) {
		tests := []struct {
			name         string
			activeColumn int
		}{
			{"First column", 0},
			{"Second column", 1},
			{"Third column", 2},
			{"Fourth column", 3},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				model := NewModel()
				model.ActiveColumn = tt.activeColumn

				result := model.GetActiveColumn()

				if result != tt.activeColumn {
					t.Errorf("GetActiveColumn() = %d, expected %d", result, tt.activeColumn)
				}
			})
		}
	})

	t.Run("GetSelectedItem", func(t *testing.T) {
		tests := []struct {
			name         string
			selectedItem int
		}{
			{"First item", 0},
			{"Middle item", 5},
			{"Last item", 10},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				model := NewModel()
				model.SelectedItem = tt.selectedItem

				result := model.GetSelectedItem()

				if result != tt.selectedItem {
					t.Errorf("GetSelectedItem() = %d, expected %d", result, tt.selectedItem)
				}
			})
		}
	})

	t.Run("GetState", func(t *testing.T) {
		tests := []struct {
			name  string
			state types.AppState
		}{
			{"MainNavigation", types.MainNavigation},
			{"SearchMode", types.SearchMode},
			{"SearchActiveNavigation", types.SearchActiveNavigation},
			{"ModalActive", types.ModalActive},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				model := NewModel()
				model.State = tt.state

				result := model.GetState()

				if result != tt.state {
					t.Errorf("GetState() = %v, expected %v", result, tt.state)
				}
			})
		}
	})

	t.Run("GetSearchQuery", func(t *testing.T) {
		tests := []struct {
			name        string
			searchQuery string
		}{
			{"Empty query", ""},
			{"Simple query", "test"},
			{"Complex query", "github-mcp"},
			{"Special characters", "test@#$%"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				model := NewModel()
				model.SearchQuery = tt.searchQuery

				result := model.GetSearchQuery()

				if result != tt.searchQuery {
					t.Errorf("GetSearchQuery() = %q, expected %q", result, tt.searchQuery)
				}
			})
		}
	})

	t.Run("GetSearchActive", func(t *testing.T) {
		tests := []struct {
			name         string
			searchActive bool
		}{
			{"Search active", true},
			{"Search inactive", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				model := NewModel()
				model.SearchActive = tt.searchActive

				result := model.GetSearchActive()

				if result != tt.searchActive {
					t.Errorf("GetSearchActive() = %v, expected %v", result, tt.searchActive)
				}
			})
		}
	})

	t.Run("GetSearchInputActive", func(t *testing.T) {
		tests := []struct {
			name              string
			searchInputActive bool
		}{
			{"Search input active", true},
			{"Search input inactive", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				model := NewModel()
				model.SearchInputActive = tt.searchInputActive

				result := model.GetSearchInputActive()

				if result != tt.searchInputActive {
					t.Errorf("GetSearchInputActive() = %v, expected %v", result, tt.searchInputActive)
				}
			})
		}
	})
}

func TestModel_GetFilteredMCPs(t *testing.T) {
	tests := []struct {
		name        string
		mcpItems    []types.MCPItem
		searchQuery string
		expected    int
		expectedNames []string
	}{
		{
			name: "No filter returns all",
			mcpItems: []types.MCPItem{
				{Name: "github", Active: true},
				{Name: "docker", Active: false},
			},
			searchQuery:   "",
			expected:      2,
			expectedNames: []string{"github", "docker"},
		},
		{
			name: "Filter by name",
			mcpItems: []types.MCPItem{
				{Name: "github-mcp", Active: true},
				{Name: "docker-mcp", Active: false},
				{Name: "context7", Active: true},
			},
			searchQuery:   "mcp",
			expected:      2,
			expectedNames: []string{"github-mcp", "docker-mcp"},
		},
		{
			name: "Case insensitive filter",
			mcpItems: []types.MCPItem{
				{Name: "GitHub", Active: true},
				{Name: "Docker", Active: false},
			},
			searchQuery:   "git",
			expected:      1,
			expectedNames: []string{"GitHub"},
		},
		{
			name: "No matches",
			mcpItems: []types.MCPItem{
				{Name: "github", Active: true},
				{Name: "docker", Active: false},
			},
			searchQuery:   "nonexistent",
			expected:      0,
			expectedNames: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.MCPItems = tt.mcpItems
			model.SearchQuery = tt.searchQuery

			result := model.GetFilteredMCPs()

			if len(result) != tt.expected {
				t.Errorf("GetFilteredMCPs() returned %d items, expected %d", len(result), tt.expected)
			}

			resultNames := make(map[string]bool)
			for _, item := range result {
				resultNames[item.Name] = true
			}

			for _, expectedName := range tt.expectedNames {
				if !resultNames[expectedName] {
					t.Errorf("GetFilteredMCPs() should include %s", expectedName)
				}
			}
		})
	}
}

func TestModel_NewModel(t *testing.T) {
	model := NewModel()

	// Test initial state
	if model.GetState() != types.MainNavigation {
		t.Errorf("NewModel() initial state should be MainNavigation")
	}

	if model.GetActiveColumn() != 0 {
		t.Errorf("NewModel() initial active column should be 0")
	}

	if model.GetSelectedItem() != 0 {
		t.Errorf("NewModel() initial selected item should be 0")
	}

	if model.GetSearchQuery() != "" {
		t.Errorf("NewModel() initial search query should be empty")
	}

	if model.GetSearchActive() {
		t.Errorf("NewModel() search should be initially inactive")
	}

	if model.GetSearchInputActive() {
		t.Errorf("NewModel() search input should be initially inactive")
	}

	if model.GetColumnCount() != 4 {
		t.Errorf("NewModel() initial column count should be 4")
	}

	// Test that MCPItems are populated
	if len(model.MCPItems) == 0 {
		t.Errorf("NewModel() should have populated MCPItems")
	}
}

func TestModel_Update(t *testing.T) {
	t.Run("WindowSizeMsg updates dimensions", func(t *testing.T) {
		model := NewModel()
		
		msg := tea.WindowSizeMsg{
			Width:  120,
			Height: 40,
		}

		updatedModel, cmd := model.Update(msg)
		m := updatedModel.(Model)

		if m.Width != 120 {
			t.Errorf("Update() should set width to 120, got %d", m.Width)
		}

		if m.Height != 40 {
			t.Errorf("Update() should set height to 40, got %d", m.Height)
		}

		if cmd != nil {
			t.Errorf("Update() with WindowSizeMsg should return nil cmd")
		}
	})

	t.Run("KeyMsg delegates to handlers", func(t *testing.T) {
		model := NewModel()
		model.Width = 120
		model.Height = 40
		
		// Test a key that should change state
		msg := tea.KeyMsg{
			Type: tea.KeyRunes,
			Runes: []rune("/"),
		}

		updatedModel, cmd := model.Update(msg)
		m := updatedModel.(Model)

		// The key handler should process the message
		// (specific behavior tested in handlers package)
		_ = m
		_ = cmd
	})

	t.Run("Unknown message type", func(t *testing.T) {
		model := NewModel()
		
		// Send an unknown message type
		msg := struct{}{}

		updatedModel, cmd := model.Update(msg)
		m := updatedModel.(Model)

		// Should return unchanged model
		if m.GetState() != model.GetState() {
			t.Errorf("Update() with unknown message should preserve state")
		}

		if cmd != nil {
			t.Errorf("Update() with unknown message should return nil cmd")
		}
	})
}

func TestModel_StateConsistency(t *testing.T) {
	t.Run("Model maintains state consistency", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithState(types.SearchActiveNavigation).
			WithSearchActive(true).
			WithSearchInputActive(true).
			WithSearchQuery("test").
			WithSelectedItem(2).
			WithActiveColumn(1).
			Build()

		// Convert to UI model
		uiModel := Model{Model: model}

		// All getters should return consistent values
		if uiModel.GetState() != types.SearchActiveNavigation {
			t.Errorf("State consistency check failed")
		}

		if !uiModel.GetSearchActive() {
			t.Errorf("SearchActive consistency check failed")
		}

		if !uiModel.GetSearchInputActive() {
			t.Errorf("SearchInputActive consistency check failed")
		}

		if uiModel.GetSearchQuery() != "test" {
			t.Errorf("SearchQuery consistency check failed")
		}

		if uiModel.GetSelectedItem() != 2 {
			t.Errorf("SelectedItem consistency check failed")
		}

		if uiModel.GetActiveColumn() != 1 {
			t.Errorf("ActiveColumn consistency check failed")
		}
	})
}

func TestModel_Integration(t *testing.T) {
	t.Run("Model works with testutil builders", func(t *testing.T) {
		// Test that our UI model integrates properly with testutil
		baseModel := testutil.NewTestModel().
			WithWindowSize(150, 50).
			WithState(types.MainNavigation).
			WithSelectedItem(3).
			Build()

		uiModel := Model{Model: baseModel}

		// Test that all properties are accessible
		if uiModel.Width != 150 || uiModel.Height != 50 {
			t.Errorf("Model should preserve window dimensions")
		}

		if uiModel.GetSelectedItem() != 3 {
			t.Errorf("Model should preserve selected item")
		}

		if uiModel.GetState() != types.MainNavigation {
			t.Errorf("Model should preserve state")
		}
	})

	t.Run("Model with realistic MCP data", func(t *testing.T) {
		model := NewModel()
		
		// Model should start with real MCP data
		if len(model.MCPItems) < 10 {
			t.Errorf("NewModel() should have substantial MCP data")
		}

		// Should have some active MCPs by default
		activeCount := 0
		for _, item := range model.MCPItems {
			if item.Active {
				activeCount++
			}
		}

		if activeCount == 0 {
			t.Errorf("NewModel() should have some active MCPs")
		}

		// Test filtering with realistic data
		model.SearchQuery = "github"
		filtered := model.GetFilteredMCPs()

		if len(filtered) == 0 {
			t.Errorf("Should find MCPs matching 'github' in default data")
		}
	})
}