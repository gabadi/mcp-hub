package services

import (
	"cc-mcp-manager/internal/ui/types"
)

// UpdateLayout updates the column layout based on terminal width
func UpdateLayout(model types.Model) types.Model {
	// Responsive breakpoints as specified in acceptance criteria
	// Prioritize 4-column MCP grid layout as per wireframe
	switch {
	case model.Width >= types.WideLayoutMin:
		// Wide: 4-column MCP grid for maximum information density
		model.ColumnCount = types.WideColumns
		columnWidth := (model.Width - 10) / 4 // Account for spacing between 4 columns
		model.Columns = []types.Column{
			{Title: "MCPs Column 1", Width: columnWidth},
			{Title: "MCPs Column 2", Width: columnWidth},
			{Title: "MCPs Column 3", Width: columnWidth},
			{Title: "MCPs Column 4", Width: columnWidth},
		}
	case model.Width >= types.MediumLayoutMin:
		// Medium: 2 columns (MCPs + Status/Details)
		model.ColumnCount = types.MediumColumns
		columnWidth := (model.Width - 6) / 2
		model.Columns = []types.Column{
			{Title: "MCPs", Width: columnWidth},
			{Title: "Status & Details", Width: columnWidth},
		}
	default:
		// Narrow: 1 column (all in one)
		model.ColumnCount = types.NarrowColumns
		model.Columns = []types.Column{
			{Title: "MCP Manager", Width: model.Width - 4},
		}
	}

	// Reset active column if it's out of bounds
	if model.ActiveColumn >= model.ColumnCount {
		model.ActiveColumn = model.ColumnCount - 1
	}
	return model
}
