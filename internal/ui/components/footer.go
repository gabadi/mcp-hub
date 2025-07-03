package components

import (
	"fmt"
	"strings"

	"cc-mcp-manager/internal/ui/services"
	"cc-mcp-manager/internal/ui/types"

	"github.com/charmbracelet/lipgloss"
)

// RenderFooter creates the application footer with status information including toggle operations
func RenderFooter(model types.Model) string {
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		Background(lipgloss.Color("#2D2D3D")).
		Padding(0, 2).
		Width(model.Width)

	var footerText string

	// Priority 1: Show toggle operation status if active
	if model.ToggleState != types.ToggleIdle {
		toggleStatus := renderToggleStatus(model)
		if toggleStatus != "" {
			footerText = toggleStatus
		}
	} else if model.SearchActive {
		// Priority 2: Show search input with cursor and mode indicator
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
		// Priority 3: Show search results info when not actively searching but have a query
		filteredMCPs := GetFilteredMCPs(model)
		footerText = fmt.Sprintf("Found %d MCPs matching '%s' • ESC to clear • Terminal: %dx%d",
			len(filteredMCPs), model.SearchQuery, model.Width, model.Height)
	} else {
		// Priority 4: Show Claude status and default info
		claudeStatusText := services.FormatClaudeStatusForDisplay(model.ClaudeStatus)
		refreshHint := services.GetRefreshKeyHint(model.ClaudeStatus)
		footerText = fmt.Sprintf("%s • %s • Terminal: %dx%d",
			claudeStatusText, refreshHint, model.Width, model.Height)
	}

	return footerStyle.Render(footerText)
}

// renderToggleStatus renders the current toggle operation status with visual indicators
func renderToggleStatus(model types.Model) string {
	switch model.ToggleState {
	case types.ToggleLoading:
		loadingStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)
		return fmt.Sprintf("%s MCP '%s'... ⏳",
			loadingStyle.Render("Toggling"), model.ToggleMCPName)

	case types.ToggleRetrying:
		retryStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF8C00")).
			Bold(true)
		return fmt.Sprintf("%s MCP '%s'... 🔄",
			retryStyle.Render("Retrying"), model.ToggleMCPName)

	case types.ToggleSuccess:
		successStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#51CF66")).
			Bold(true)
		return fmt.Sprintf("%s MCP '%s' ✓",
			successStyle.Render("Success:"), model.ToggleMCPName)

	case types.ToggleError:
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)
		return fmt.Sprintf("%s %s ✗",
			errorStyle.Render("Error:"), model.ToggleError)
	}
	return ""
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
