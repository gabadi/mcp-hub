package components

import (
	"fmt"
	"strings"
	"time"

	"mcp-hub/internal/ui/services"
	"mcp-hub/internal/ui/types"

	"github.com/charmbracelet/lipgloss"
)

// RenderFooter creates the application footer with status information including toggle operations
func RenderFooter(model types.Model) string {
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		Background(lipgloss.Color("#2D2D3D")).
		Padding(0, 2).
		Width(model.Width)

	footerText := getFooterContent(model)
	return footerStyle.Render(footerText)
}

// getFooterContent determines the appropriate footer content based on model state
func getFooterContent(model types.Model) string {
	switch {
	case model.ToggleState != types.ToggleIdle:
		return getToggleFooterContent(model)
	case model.SearchActive:
		return getSearchActiveFooterContent(model)
	case model.SearchQuery != "":
		return getSearchResultsFooterContent(model)
	default:
		return getProjectContextFooterContent(model)
	}
}

// getToggleFooterContent returns footer content for toggle operations
func getToggleFooterContent(model types.Model) string {
	toggleStatus := renderToggleStatus(model)
	if toggleStatus != "" {
		return toggleStatus
	}
	return ""
}

// getSearchActiveFooterContent returns footer content for active search
func getSearchActiveFooterContent(model types.Model) string {
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

	var searchDisplay string
	if model.SearchInputActive {
		searchDisplay = model.SearchQuery + cursor
	} else {
		searchDisplay = model.SearchQuery
	}
	return fmt.Sprintf("Search: %s%s", searchStyle.Render(searchDisplay), modeIndicator)
}

// getSearchResultsFooterContent returns footer content for search results
func getSearchResultsFooterContent(model types.Model) string {
	filteredMCPs := GetFilteredMCPs(model)
	return fmt.Sprintf("Found %d MCPs matching '%s' • ESC to clear • Terminal: %dx%d",
		len(filteredMCPs), model.SearchQuery, model.Width, model.Height)
}

// getProjectContextFooterContent returns footer content for project context
func getProjectContextFooterContent(model types.Model) string {
	var projectContext types.ProjectContext

	// Use model's project context if it has valid display info, otherwise compute it
	if model.ProjectContext.DisplayPath != "" {
		projectContext = model.ProjectContext
	} else {
		projectContext = services.GetProjectContext(model)
	}

	// Format project context display
	contextInfo := fmt.Sprintf("📁 %s • %d/%d MCPs • %s",
		projectContext.DisplayPath,
		projectContext.ActiveMCPs,
		projectContext.TotalMCPs,
		projectContext.SyncStatusText)

	// Add last sync time if available
	if !projectContext.LastSyncTime.IsZero() {
		timeSinceSync := time.Since(projectContext.LastSyncTime)
		if timeSinceSync < time.Hour {
			contextInfo += fmt.Sprintf(" • Last sync: %s ago",
				formatDuration(timeSinceSync))
		} else {
			contextInfo += fmt.Sprintf(" • Last sync: %s",
				projectContext.LastSyncTime.Format("15:04"))
		}
	}

	refreshHint := services.GetRefreshKeyHint(model.ClaudeStatus)
	return fmt.Sprintf("%s • %s", contextInfo, refreshHint)
}

// renderToggleStatus renders the current toggle operation status with visual indicators
func renderToggleStatus(model types.Model) string {
	switch model.ToggleState {
	case types.ToggleIdle:
		return ""
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

// formatDuration formats a duration for display in the footer
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh", int(d.Hours()))
}
