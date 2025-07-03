package components

import (
	"fmt"
	"strings"
	"time"

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
		// Priority 4: Show project context information
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
		footerText = fmt.Sprintf("%s • %s", contextInfo, refreshHint)
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

// formatDuration formats a duration for display in the footer
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	} else {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
}
