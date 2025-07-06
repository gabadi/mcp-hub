package components

import (
	"fmt"

	"mcp-hub/internal/ui/services"
	"mcp-hub/internal/ui/types"

	"github.com/charmbracelet/lipgloss"
)

// Header string constants
const (
	EnterConfirmEscCancel = "Enter=Confirm • ESC=Cancel"
)

// RenderHeader creates the application header with shortcuts and context
func RenderHeader(model types.Model) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		Background(lipgloss.Color("#1E1E2E")).
		Padding(0, 2).
		Width(model.Width)

	// Build shortcuts display based on current state
	var shortcuts string
	switch model.State {
	case types.MainNavigation:
		shortcuts = "A=Add • D=Delete • E=Edit • /=Search • Tab=Focus Search • R=Refresh Claude • ESC=Exit • ↑↓←→=Navigate"
	case types.SearchMode:
		shortcuts = "Type to search • Enter=Apply • ESC=Cancel"
	case types.SearchActiveNavigation:
		if model.SearchInputActive {
			shortcuts = "Type to search • Tab=Navigate Mode • ↑↓←→=Navigate • Space=Toggle • R=Refresh • Enter=Apply • ESC=Cancel"
		} else {
			shortcuts = "Navigate Mode • Tab=Input Mode • ↑↓←→=Navigate • Space=Toggle • R=Refresh • Enter=Apply • ESC=Cancel"
		}
	case types.ModalActive:
		shortcuts = EnterConfirmEscCancel
	}

	// Context information
	activeCount := 0
	for _, item := range model.MCPItems {
		if item.Active {
			activeCount++
		}
	}

	// Claude status information
	claudeStatusText := services.FormatClaudeStatusForDisplay(model.ClaudeStatus)

	contextInfo := fmt.Sprintf("MCPs: %d/%d Active • Layout: %s • %s",
		activeCount, len(model.MCPItems), GetLayoutName(model), claudeStatusText)

	title := "MCP Manager v1.0"

	// Create header content with proper spacing
	headerContent := fmt.Sprintf("%s\n%s\n%s", title, shortcuts, contextInfo)

	return headerStyle.Render(headerContent)
}

// GetLayoutName returns the current layout name for display
func GetLayoutName(model types.Model) string {
	switch model.ColumnCount {
	case 1:
		return "Narrow"
	case 2:
		return "Medium"
	case 3:
		return "Wide (3-panel)"
	case 4:
		return "Grid (4-column MCP)"
	default:
		return "Unknown"
	}
}
