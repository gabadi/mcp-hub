package components

import (
	"fmt"
	"strings"

	"cc-mcp-manager/internal/ui/types"
	"github.com/charmbracelet/lipgloss"
)

// RenderFooter creates the application footer with status information
func RenderFooter(model types.Model) string {
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		Background(lipgloss.Color("#2D2D3D")).
		Padding(0, 2).
		Width(model.Width)

	var footerText string

	if model.SearchActive {
		// Show search input with cursor and mode indicator
		searchStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)
		cursor := "_"

		var modeIndicator string
		if model.State == types.SearchActiveNavigation {
			if model.SearchInputActive {
				modeIndicator = " [INPUT MODE]"
			} else {
				modeIndicator = " [NAVIGATION MODE]"
			}
		}

		// Show the search query with cursor in styled box
		searchDisplay := model.SearchQuery
		if model.SearchInputActive {
			searchDisplay = model.SearchQuery + cursor
		} else {
			searchDisplay = model.SearchQuery
		}
		footerText = fmt.Sprintf("Search: %s%s", searchStyle.Render(searchDisplay), modeIndicator)
	} else if model.SearchQuery != "" {
		// Show search results info when not actively searching but have a query
		filteredMCPs := GetFilteredMCPs(model)
		footerText = fmt.Sprintf("Found %d MCPs matching '%s' • ESC to clear • Terminal: %dx%d",
			len(filteredMCPs), model.SearchQuery, model.Width, model.Height)
	} else {
		footerText = fmt.Sprintf("Terminal: %dx%d • Search: '%s' • Use arrow keys to navigate, Tab or / for search",
			model.Width, model.Height, model.SearchQuery)
	}

	return footerStyle.Render(footerText)
}

// GetFilteredMCPs returns MCPs filtered by search query
func GetFilteredMCPs(model types.Model) []types.MCPItem {
	// If no search query, return all MCPs
	if model.SearchQuery == "" {
		return model.MCPItems
	}

	// Filter MCPs by search query directly
	var filtered []types.MCPItem
	query := strings.ToLower(model.SearchQuery)
	for _, item := range model.MCPItems {
		if strings.Contains(strings.ToLower(item.Name), query) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
