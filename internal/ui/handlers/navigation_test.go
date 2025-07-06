package handlers

import (
	"testing"

	"mcp-hub/internal/testutil"
	"mcp-hub/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestHandleMainNavigationKeys(t *testing.T) {
	t.Run("Arrow key navigation", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithState(types.MainNavigation).
			WithActiveColumn(0).
			WithSelectedItem(0).
			Build()

		tests := []struct {
			name string
			key  tea.KeyMsg
		}{
			{"Up arrow", tea.KeyMsg{Type: tea.KeyUp}},
			{"Down arrow", tea.KeyMsg{Type: tea.KeyDown}},
			{"Left arrow", tea.KeyMsg{Type: tea.KeyLeft}},
			{"Right arrow", tea.KeyMsg{Type: tea.KeyRight}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				newModel, cmd := HandleMainNavigationKeys(model, tt.key.String())
				assert.NotNil(t, newModel)
				_ = cmd // May or may not be nil
			})
		}
	})

	t.Run("Action keys", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithState(types.MainNavigation).
			WithMCPs([]types.MCPItem{
				{Name: "test-mcp", Type: "CMD", Active: false},
			}).
			Build()

		tests := []struct {
			name string
			key  tea.KeyMsg
		}{
			{"Space toggle", tea.KeyMsg{Type: tea.KeySpace}},
			{"Enter select", tea.KeyMsg{Type: tea.KeyEnter}},
			{"Search activate", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")}},
			{"Add MCP", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}},
			{"Edit MCP", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")}},
			{"Delete MCP", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")}},
			{"Refresh", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				newModel, cmd := HandleMainNavigationKeys(model, tt.key.String())
				assert.NotNil(t, newModel)
				_ = cmd
			})
		}
	})
}

func TestHandleSearchNavigationKeys(t *testing.T) {
	t.Run("Search mode navigation", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithState(types.SearchActiveNavigation).
			WithSearchActive(true).
			WithSearchInputActive(true).
			WithSearchQuery("test").
			Build()

		tests := []struct {
			name string
			key  tea.KeyMsg
		}{
			{"Tab toggle", tea.KeyMsg{Type: tea.KeyTab}},
			{"Escape cancel", tea.KeyMsg{Type: tea.KeyEsc}},
			{"Enter apply", tea.KeyMsg{Type: tea.KeyEnter}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				newModel, cmd := HandleSearchNavigationKeys(model, tt.key.String())
				assert.NotNil(t, newModel)
				_ = cmd
			})
		}
	})

	t.Run("Search input mode", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithState(types.SearchActiveNavigation).
			WithSearchActive(true).
			WithSearchInputActive(true).
			WithSearchQuery("").
			Build()

		// Test character input
		key := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}
		newModel, cmd := HandleSearchNavigationKeys(model, key.String())

		assert.NotNil(t, newModel)
		assert.Equal(t, "a", newModel.SearchQuery)
		_ = cmd
	})

	t.Run("Search backspace", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithState(types.SearchActiveNavigation).
			WithSearchActive(true).
			WithSearchInputActive(true).
			WithSearchQuery("test").
			Build()

		key := tea.KeyMsg{Type: tea.KeyBackspace}
		newModel, cmd := HandleSearchNavigationKeys(model, key.String())

		assert.NotNil(t, newModel)
		assert.Equal(t, "tes", newModel.SearchQuery)
		_ = cmd
	})
}

func TestNavigationMethods(t *testing.T) {
	t.Run("NavigateUp", func(t *testing.T) {
		// Test with 4-column layout - NavigateUp should work within grid
		model := testutil.NewTestModel().
			WithWindowSize(120, 40). // 4 columns
			WithMCPs([]types.MCPItem{
				{Name: "mcp-1", Type: "CMD"},
				{Name: "mcp-2", Type: "CMD"},
				{Name: "mcp-3", Type: "CMD"},
				{Name: "mcp-4", Type: "CMD"},
				{Name: "mcp-5", Type: "CMD"},
			}).
			WithSelectedItem(4). // Position 4 (second row, first column)
			Build()

		newModel := NavigateUp(model)
		assert.Equal(t, 0, newModel.SelectedItem) // Should move to position 0 (first row, first column)
	})

	t.Run("NavigateDown", func(t *testing.T) {
		// Test with 4-column layout - NavigateDown should work within grid
		model := testutil.NewTestModel().
			WithWindowSize(120, 40). // 4 columns
			WithMCPs([]types.MCPItem{
				{Name: "mcp-1", Type: "CMD"},
				{Name: "mcp-2", Type: "CMD"},
				{Name: "mcp-3", Type: "CMD"},
				{Name: "mcp-4", Type: "CMD"},
				{Name: "mcp-5", Type: "CMD"},
			}).
			WithSelectedItem(0). // Position 0 (first row, first column)
			Build()

		newModel := NavigateDown(model)
		assert.Equal(t, 4, newModel.SelectedItem) // Should move to position 4 (second row, first column)
	})

	t.Run("NavigateLeft", func(t *testing.T) {
		// Test with 4-column layout - NavigateLeft should work within grid
		model := testutil.NewTestModel().
			WithWindowSize(120, 40). // 4 columns
			WithMCPs([]types.MCPItem{
				{Name: "mcp-1", Type: "CMD"},
				{Name: "mcp-2", Type: "CMD"},
				{Name: "mcp-3", Type: "CMD"},
			}).
			WithSelectedItem(2). // Position 2 (first row, third column)
			Build()

		newModel := NavigateLeft(model)
		assert.Equal(t, 1, newModel.SelectedItem) // Should move to position 1 (first row, second column)
	})

	t.Run("NavigateRight", func(t *testing.T) {
		// Test with 4-column layout - NavigateRight should work within grid
		model := testutil.NewTestModel().
			WithWindowSize(120, 40). // 4 columns
			WithMCPs([]types.MCPItem{
				{Name: "mcp-1", Type: "CMD"},
				{Name: "mcp-2", Type: "CMD"},
				{Name: "mcp-3", Type: "CMD"},
			}).
			WithSelectedItem(0). // Position 0 (first row, first column)
			Build()

		newModel := NavigateRight(model)
		assert.Equal(t, 1, newModel.SelectedItem) // Should move to position 1 (first row, second column)
	})
}

func TestCommandGenerators(t *testing.T) {
	t.Run("RefreshClaudeStatusCmd", func(t *testing.T) {
		cmd := RefreshClaudeStatusCmd()
		assert.NotNil(t, cmd)
	})

	t.Run("EnhancedToggleMCPCmd", func(t *testing.T) {
		mcp := types.MCPItem{Name: "test-mcp", Type: "CMD"}
		cmd := EnhancedToggleMCPCmd(mcp.Name, false, &mcp)
		assert.NotNil(t, cmd)
	})

	t.Run("TimerCmd", func(t *testing.T) {
		cmd := TimerCmd("1")
		assert.NotNil(t, cmd)
	})

	t.Run("StartupLoadingCmd", func(t *testing.T) {
		cmd := StartupLoadingCmd()
		assert.NotNil(t, cmd)
	})

	t.Run("StartupLoadingProgressCmd", func(t *testing.T) {
		cmd := StartupLoadingProgressCmd(1)
		assert.NotNil(t, cmd)
	})

	t.Run("StartupLoadingTimerCmd", func(t *testing.T) {
		cmd := StartupLoadingTimerCmd(1)
		assert.NotNil(t, cmd)
	})

	t.Run("RefreshLoadingCmd", func(t *testing.T) {
		cmd := RefreshLoadingCmd()
		assert.NotNil(t, cmd)
	})

	t.Run("RefreshLoadingProgressCmd", func(t *testing.T) {
		cmd := RefreshLoadingProgressCmd(1)
		assert.NotNil(t, cmd)
	})

	t.Run("RefreshLoadingTimerCmd", func(t *testing.T) {
		cmd := RefreshLoadingTimerCmd(1)
		assert.NotNil(t, cmd)
	})

	t.Run("LoadingSpinnerCmd", func(t *testing.T) {
		cmd := LoadingSpinnerCmd(types.LoadingStartup)
		assert.NotNil(t, cmd)
	})
}

func TestUtilityFunctions(t *testing.T) {
	t.Run("pasteToSearchQuery", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithSearchQuery("").
			Build()

		newModel := pasteToSearchQuery(model)
		var cmd tea.Cmd
		assert.NotNil(t, newModel)
		_ = cmd
	})

	t.Run("populateFormDataFromMCP", func(t *testing.T) {
		mcp := types.MCPItem{
			Name:        "test-mcp",
			Type:        "CMD",
			Command:     "test-command",
			Args:        []string{"arg1", "arg2"},
			Environment: map[string]string{"KEY": "value"},
		}

		formData := populateFormDataFromMCP(mcp)

		assert.Equal(t, mcp.Name, formData.Name)
		assert.Equal(t, mcp.Command, formData.Command)
		assert.NotEmpty(t, formData.Args)
		assert.NotEmpty(t, formData.Environment)
	})

	t.Run("formatArgsForDisplay", func(t *testing.T) {
		args := []string{"arg1", "arg2", "arg3"}
		result := formatArgsForDisplay(args)
		assert.Contains(t, result, "arg1")
		assert.Contains(t, result, "arg2")
		assert.Contains(t, result, "arg3")
	})

	t.Run("formatEnvironmentForDisplay", func(t *testing.T) {
		env := map[string]string{
			"KEY1": "value1",
			"KEY2": "value2",
		}
		result := formatEnvironmentForDisplay(env)
		assert.Contains(t, result, "KEY1=value1")
		assert.Contains(t, result, "KEY2=value2")
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("Navigation with empty MCP list", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithMCPs([]types.MCPItem{}).
			WithSelectedItem(0).
			Build()

		// Should handle empty list gracefully
		newModel := NavigateDown(model)
		assert.Equal(t, 0, newModel.SelectedItem)

		newModel = NavigateUp(model)
		assert.Equal(t, 0, newModel.SelectedItem)
	})

	t.Run("Navigation at boundaries", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithMCPs([]types.MCPItem{
				{Name: "mcp-1", Type: "CMD"},
			}).
			WithSelectedItem(0).
			Build()

		// At top boundary
		newModel := NavigateUp(model)
		assert.Equal(t, 0, newModel.SelectedItem)

		// At bottom boundary
		newModel = NavigateDown(model)
		assert.Equal(t, 0, newModel.SelectedItem)
	})

	t.Run("Column navigation at boundaries", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithActiveColumn(0).
			Build()
		model.ColumnCount = 4

		// At left boundary
		newModel := NavigateLeft(model)
		assert.Equal(t, 0, newModel.ActiveColumn)

		// Move to right boundary
		model.ActiveColumn = 3
		newModel = NavigateRight(model)
		assert.Equal(t, 3, newModel.ActiveColumn)
	})
}

func TestSearchHandling(t *testing.T) {
	t.Run("Search with special characters", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithState(types.SearchActiveNavigation).
			WithSearchActive(true).
			WithSearchInputActive(true).
			WithSearchQuery("").
			Build()

		specialChars := []string{"@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "_", "+", "="}

		for _, char := range specialChars {
			key := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(char)}
			newModel, _ := HandleSearchNavigationKeys(model, key.String())
			assert.Contains(t, newModel.SearchQuery, char)
			model = newModel
		}
	})

	t.Run("Search query clearing", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithState(types.SearchActiveNavigation).
			WithSearchActive(true).
			WithSearchInputActive(true).
			WithSearchQuery("test query").
			Build()

		// Clear with multiple backspaces
		for i := 0; i < 20; i++ {
			key := tea.KeyMsg{Type: tea.KeyBackspace}
			newModel, _ := HandleSearchNavigationKeys(model, key.String())
			model = newModel
		}

		assert.Equal(t, "", model.SearchQuery)
	})
}
