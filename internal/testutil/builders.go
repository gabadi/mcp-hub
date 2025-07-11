// Package testutil provides utilities for building test models and data.
package testutil

import (
	"mcp-hub/internal/platform"
	"mcp-hub/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

// TestModelBuilder provides a fluent interface for building test models
type TestModelBuilder struct {
	model types.Model
}

// NewTestModel creates a new TestModelBuilder with default values
func NewTestModel() *TestModelBuilder {
	return &TestModelBuilder{
		model: types.NewModel(platform.GetMockPlatformService()),
	}
}

// WithWindowSize sets the window dimensions and updates layout
func (b *TestModelBuilder) WithWindowSize(width, height int) *TestModelBuilder {
	b.model.Width = width
	b.model.Height = height
	// Update layout based on width
	b.updateLayout()
	return b
}

// WithState sets the application state
func (b *TestModelBuilder) WithState(state types.AppState) *TestModelBuilder {
	b.model.State = state
	return b
}

// WithActiveColumn sets the active column
func (b *TestModelBuilder) WithActiveColumn(column int) *TestModelBuilder {
	b.model.ActiveColumn = column
	return b
}

// WithSelectedItem sets the selected item
func (b *TestModelBuilder) WithSelectedItem(item int) *TestModelBuilder {
	b.model.SelectedItem = item
	return b
}

// WithSearchQuery sets the search query
func (b *TestModelBuilder) WithSearchQuery(query string) *TestModelBuilder {
	b.model.SearchQuery = query
	return b
}

// WithSearchActive sets the search active state
func (b *TestModelBuilder) WithSearchActive(active bool) *TestModelBuilder {
	b.model.SearchActive = active
	return b
}

// WithSearchInputActive sets the search input active state
func (b *TestModelBuilder) WithSearchInputActive(active bool) *TestModelBuilder {
	b.model.SearchInputActive = active
	return b
}

// WithMCPs sets the MCP items for testing
func (b *TestModelBuilder) WithMCPs(mcps []types.MCPItem) *TestModelBuilder {
	b.model.MCPItems = mcps
	return b
}

// WithTempStorage is a placeholder for storage configuration (for testing with temp directories)
func (b *TestModelBuilder) WithTempStorage(_ string) *TestModelBuilder {
	// This is used for testing but doesn't modify the model directly
	// The tempDir is handled by test code when calling storage functions
	return b
}

// Build returns the constructed model
func (b *TestModelBuilder) Build() types.Model {
	return b.model
}

// BuildAndUpdate returns the model after applying a bubbletea message
func (b *TestModelBuilder) BuildAndUpdate(_ tea.Msg) types.Model {
	// This will be implemented when we refactor the update logic
	return b.model
}

// updateLayout updates the column layout based on terminal width
func (b *TestModelBuilder) updateLayout() {
	// Responsive breakpoints as specified in acceptance criteria
	switch {
	case b.model.Width >= types.WideLayoutMin:
		// Wide: 4-column MCP grid for maximum information density
		b.model.ColumnCount = types.WideColumns
		columnWidth := (b.model.Width - 10) / 4 // Account for spacing between 4 columns
		b.model.Columns = []types.Column{
			{Title: "MCPs Column 1", Width: columnWidth},
			{Title: "MCPs Column 2", Width: columnWidth},
			{Title: "MCPs Column 3", Width: columnWidth},
			{Title: "MCPs Column 4", Width: columnWidth},
		}
	case b.model.Width >= types.MediumLayoutMin:
		// Medium: 3 columns (MCPs + Status/Details)
		b.model.ColumnCount = 3
		columnWidth := (b.model.Width - 8) / 3
		b.model.Columns = []types.Column{
			{Title: "MCPs", Width: columnWidth},
			{Title: "Status", Width: columnWidth},
			{Title: "Details", Width: columnWidth},
		}
	default:
		// Narrow: 2 columns (all in one)
		b.model.ColumnCount = 2
		columnWidth := (b.model.Width - 6) / 2
		b.model.Columns = []types.Column{
			{Title: "MCPs", Width: columnWidth},
			{Title: "Status", Width: columnWidth},
		}
	}

	// Reset active column if it's out of bounds
	if b.model.ActiveColumn >= b.model.ColumnCount {
		b.model.ActiveColumn = b.model.ColumnCount - 1
	}
}
