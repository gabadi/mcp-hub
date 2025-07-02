package components

import (
	"fmt"
	"strings"

	"cc-mcp-manager/internal/ui/services"
	"cc-mcp-manager/internal/ui/types"

	"github.com/charmbracelet/lipgloss"
)

// RenderFourColumnGrid renders the 4-column MCP grid layout
func RenderFourColumnGrid(model types.Model) string {
	filteredMCPs := services.GetFilteredMCPs(model)

	if len(filteredMCPs) == 0 {
		// Show "No results" message when search returns no results
		noResultsStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Align(lipgloss.Center).
			Width(model.Width).
			Height(model.Height - 8)
		return noResultsStyle.Render("No MCPs found matching your search")
	}

	// Calculate grid dimensions - aim for ~10 rows with 4 columns
	gridRows := (len(filteredMCPs) + 3) / 4 // Round up division
	if gridRows < 10 {
		gridRows = 10 // Minimum 10 rows for consistent layout
	}

	// Build the grid as a simple string without column separators
	var gridLines []string

	for row := 0; row < gridRows; row++ {
		var line []string

		for col := 0; col < 4; col++ {
			mcpIndex := row*4 + col

			if mcpIndex < len(filteredMCPs) {
				item := filteredMCPs[mcpIndex]

				// Enhanced status indicator with toggle state
				status := getEnhancedStatusIndicator(model, item)

				// Highlight selected item by comparing index directly
				// Use FilteredSelectedIndex when search is active
				isSelected := false
				if model.SearchQuery != "" {
					isSelected = (mcpIndex == model.FilteredSelectedIndex)
				} else {
					isSelected = (mcpIndex == model.SelectedItem)
				}

				// Create base item text (without styling)
				baseText := fmt.Sprintf("%s %s", status, item.Name)

				// Calculate padding needed BEFORE styling
				currentWidth := lipgloss.Width(baseText)
				paddingNeeded := types.COLUMN_WIDTH - currentWidth
				if paddingNeeded < 0 {
					paddingNeeded = 0
				}

				// Apply padding first
				paddedText := baseText + strings.Repeat(" ", paddingNeeded)

				// Then apply styling based on selection
				if isSelected {
					styledText := lipgloss.NewStyle().
						Background(lipgloss.Color("#7C3AED")).
						Foreground(lipgloss.Color("#FFFFFF")).
						Bold(true).
						Render(paddedText)
					line = append(line, styledText)
				} else {
					line = append(line, paddedText)
				}
			} else {
				// Empty cell with proper spacing
				line = append(line, strings.Repeat(" ", types.COLUMN_WIDTH))
			}
		}

		// Join columns without separators, just spaces
		gridLines = append(gridLines, strings.Join(line, ""))
	}

	// Create container style that fills available space
	gridStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(1).
		Width(model.Width).
		Height(model.Height - 8)

	return gridStyle.Render(fmt.Sprintf("MCP Inventory\n\n%s", strings.Join(gridLines, "\n")))
}

// getEnhancedStatusIndicator returns the appropriate status indicator with toggle operation state
func getEnhancedStatusIndicator(model types.Model, item types.MCPItem) string {
	// Check if this MCP is currently being toggled
	if model.ToggleMCPName == item.Name {
		switch model.ToggleState {
		case types.ToggleLoading:
			return "⏳" // Loading spinner
		case types.ToggleRetrying:
			return "🔄" // Retry indicator
		case types.ToggleSuccess:
			// Show success briefly, then return to normal
			if item.Active {
				return "✅" // Success with active state
			}
			return "✓" // Success with inactive state
		case types.ToggleError:
			return "✗" // Error indicator
		}
	}
	
	// Default status indicators
	if item.Active {
		return "●" // Active
	}
	return "○" // Inactive
}

// RenderMCPList renders a simple list of MCPs for other layouts
func RenderMCPList(model types.Model) string {
	filteredMCPs := services.GetFilteredMCPs(model)

	// Debug: Show count of MCPs
	if len(model.MCPItems) == 0 {
		return fmt.Sprintf("No MCPs loaded from inventory (total: %d)", len(model.MCPItems))
	}

	if len(filteredMCPs) == 0 {
		return fmt.Sprintf("No MCPs matching filter (total: %d, filtered: %d)", len(model.MCPItems), len(filteredMCPs))
	}

	var items []string
	for i, item := range filteredMCPs {
		style := lipgloss.NewStyle().Padding(0, 1)

		// Highlight selected item
		if i == model.SelectedItem {
			style = style.Background(lipgloss.Color("#7C3AED")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)
		}

		// Enhanced status indicator with toggle state
		status := getEnhancedStatusIndicator(model, item)

		itemText := fmt.Sprintf("%s %s", status, item.Name)
		items = append(items, style.Render(itemText))
	}

	return strings.Join(items, "\n")
}
