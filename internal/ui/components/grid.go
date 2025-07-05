package components

import (
	"fmt"
	"strings"

	"mcp-hub/internal/ui/services"
	"mcp-hub/internal/ui/types"

	"github.com/charmbracelet/lipgloss"
)

// RenderFourColumnGrid renders the 4-column MCP grid layout
func RenderFourColumnGrid(model types.Model) string {
	filteredMCPs := services.GetFilteredMCPs(model)

	if len(filteredMCPs) == 0 {
		return renderNoResultsMessage(model)
	}

	gridLines := buildGridLines(model, filteredMCPs)
	gridStyle := createGridStyle(model)

	return gridStyle.Render(fmt.Sprintf("MCP Inventory\n\n%s", strings.Join(gridLines, "\n")))
}

// renderNoResultsMessage creates the no results display
func renderNoResultsMessage(model types.Model) string {
	noResultsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Align(lipgloss.Center).
		Width(model.Width).
		Height(model.Height - 8)
	return noResultsStyle.Render("No MCPs found matching your search")
}

// buildGridLines creates the grid content lines
func buildGridLines(model types.Model, filteredMCPs []types.MCPItem) []string {
	// Calculate grid dimensions - aim for ~10 rows with 4 columns
	gridRows := (len(filteredMCPs) + 3) / 4 // Round up division
	if gridRows < 10 {
		gridRows = 10 // Minimum 10 rows for consistent layout
	}

	var gridLines []string
	for row := 0; row < gridRows; row++ {
		line := buildGridRow(model, filteredMCPs, row)
		gridLines = append(gridLines, strings.Join(line, ""))
	}
	return gridLines
}

// buildGridRow creates a single row of the grid
func buildGridRow(model types.Model, filteredMCPs []types.MCPItem, row int) []string {
	var line []string
	for col := 0; col < 4; col++ {
		mcpIndex := row*4 + col
		if mcpIndex < len(filteredMCPs) {
			cellContent := renderGridCell(model, filteredMCPs[mcpIndex], mcpIndex)
			line = append(line, cellContent)
		} else {
			// Empty cell with proper spacing
			line = append(line, strings.Repeat(" ", types.ColumnWidth))
		}
	}
	return line
}

// renderGridCell creates the content for a single grid cell
func renderGridCell(model types.Model, item types.MCPItem, mcpIndex int) string {
	// Enhanced status indicator with toggle state
	status := getEnhancedStatusIndicator(model, item)

	// Determine if this item is selected
	isSelected := isItemSelected(model, mcpIndex)

	// Create base item text (without styling)
	baseText := fmt.Sprintf("%s %s", status, item.Name)

	// Calculate padding needed BEFORE styling
	currentWidth := lipgloss.Width(baseText)
	paddingNeeded := types.ColumnWidth - currentWidth
	if paddingNeeded < 0 {
		paddingNeeded = 0
	}

	// Apply padding first
	paddedText := baseText + strings.Repeat(" ", paddingNeeded)

	// Apply styling based on selection
	if isSelected {
		return lipgloss.NewStyle().
			Background(lipgloss.Color("#7C3AED")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Render(paddedText)
	}
	return paddedText
}

// isItemSelected determines if an item is currently selected
func isItemSelected(model types.Model, mcpIndex int) bool {
	if model.SearchQuery != "" {
		return mcpIndex == model.FilteredSelectedIndex
	}
	return mcpIndex == model.SelectedItem
}

// createGridStyle creates the styling for the grid container
func createGridStyle(model types.Model) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(1).
		Width(model.Width).
		Height(model.Height - 8)
}

// getEnhancedStatusIndicator returns the appropriate status indicator with toggle operation state
func getEnhancedStatusIndicator(model types.Model, item types.MCPItem) string {
	// Check if this MCP is currently being toggled
	if model.ToggleMCPName == item.Name {
		switch model.ToggleState {
		case types.ToggleIdle:
			// Fall through to default status indicators
		case types.ToggleLoading:
			return "⏳" // Loading spinner
		case types.ToggleRetrying:
			return "🔄" // Retry indicator
		case types.ToggleSuccess:
			// Show success briefly, then return to normal
			if item.Active {
				return "✅" // Success - MCP activated
			}
			return "◦" // Success - MCP deactivated (removed)
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
