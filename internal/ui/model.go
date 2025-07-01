package ui

import (
	"fmt"

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
	var model Model

	if err != nil {
		// Fall back to default model if loading fails
		model = Model{
			Model: types.NewModel(),
		}
	} else if len(mcpItems) == 0 {
		// First-time setup: save defaults to storage
		defaultModel := types.NewModel()
		if saveErr := services.SaveInventory(defaultModel.MCPItems); saveErr != nil {
			// Log error but continue - the app should still work
			// Error is already logged in SaveInventory
		}
		model = Model{
			Model: defaultModel,
		}
	} else {
		// Use loaded inventory
		model = Model{
			Model: types.NewModelWithMCPs(mcpItems),
		}
	}

	return model
}

// Init initializes the application and returns initial commands
func (m Model) Init() tea.Cmd {
	// Initialize Claude status on startup
	return handlers.RefreshClaudeStatusCmd()
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

	case handlers.ClaudeStatusMsg:
		// Update model with Claude status
		m.Model = services.UpdateModelWithClaudeStatus(m.Model, msg.Status)
		// Sync MCP status if Claude is available and has active MCPs
		if msg.Status.Available && len(msg.Status.ActiveMCPs) > 0 {
			m.Model = services.SyncMCPStatus(m.Model, msg.Status.ActiveMCPs)
			// Save updated inventory after sync
			if err := services.SaveModelInventory(m.Model); err != nil {
				// Set error message but don't fail
				m.Model.SuccessMessage = fmt.Sprintf("Claude status updated, but failed to save inventory: %v", err)
				m.Model.SuccessTimer = 240 // Show error for 4 seconds
			} else {
				m.Model.SuccessMessage = "Claude status refreshed and MCPs synced"
				m.Model.SuccessTimer = 120 // Show success for 2 seconds
			}
		} else if msg.Status.Available {
			m.Model.SuccessMessage = "Claude status refreshed"
			m.Model.SuccessTimer = 120
		} else {
			m.Model.SuccessMessage = "Claude CLI not available"
			m.Model.SuccessTimer = 180 // Show message for 3 seconds
		}
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
