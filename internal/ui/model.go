package ui

import (
	"cc-mcp-manager/internal/ui/handlers"
	"cc-mcp-manager/internal/ui/services"
	"cc-mcp-manager/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

// Model is a wrapper around the types.Model to provide UI-specific methods
type Model struct {
	types.Model
}

// NewModel creates a new application model with inventory loaded from storage
func NewModel() Model {
	// Try to load inventory from storage
	mcpItems, err := services.LoadInventory()
	if err != nil {
		// Fall back to default model if loading fails
		return Model{
			Model: types.NewModel(),
		}
	}

	// If no items loaded (empty inventory), start with defaults for first-time users
	if len(mcpItems) == 0 {
		// First-time setup: save defaults to storage
		defaultModel := types.NewModel()
		if saveErr := services.SaveInventory(defaultModel.MCPItems); saveErr != nil {
			// Log error but continue - the app should still work
			// Error is already logged in SaveInventory
		}
		return Model{
			Model: defaultModel,
		}
	}

	// Use loaded inventory
	return Model{
		Model: types.NewModelWithMCPs(mcpItems),
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Model = services.UpdateLayout(m.Model)

	case tea.KeyMsg:
		var cmd tea.Cmd
		m.Model, cmd = handlers.HandleKeyPress(m.Model, msg)
		return m, cmd

	case handlers.SuccessMsg:
		m.Model.SuccessMessage = msg.Message
	}

	return m, nil
}

// All key handling has been moved to handlers package

// Getter methods for testing

// GetColumnCount returns the current number of columns
func (m Model) GetColumnCount() int {
	return m.Model.ColumnCount
}

// GetActiveColumn returns the currently active column index
func (m Model) GetActiveColumn() int {
	return m.Model.ActiveColumn
}

// GetSelectedItem returns the currently selected item index
func (m Model) GetSelectedItem() int {
	return m.Model.SelectedItem
}

// GetState returns the current application state
func (m Model) GetState() types.AppState {
	return m.Model.State
}

// GetSearchQuery returns the current search query
func (m Model) GetSearchQuery() string {
	return m.Model.SearchQuery
}

// GetSearchActive returns whether search is currently active
func (m Model) GetSearchActive() bool {
	return m.Model.SearchActive
}

// GetSearchInputActive returns whether search input is currently active
func (m Model) GetSearchInputActive() bool {
	return m.Model.SearchInputActive
}

// GetFilteredMCPs returns MCPs filtered by search query
func (m Model) GetFilteredMCPs() []types.MCPItem {
	return services.GetFilteredMCPs(m.Model)
}

// All layout and navigation logic has been moved to services and handlers packages
