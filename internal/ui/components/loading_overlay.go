package components

import (
	"strings"

	"cc-mcp-manager/internal/ui/types"

	"github.com/charmbracelet/lipgloss"
)

// RenderLoadingOverlay renders the loading overlay on top of the main content
func RenderLoadingOverlay(model types.Model, width, height int, baseContent string) string {
	if model.LoadingOverlay == nil || !model.LoadingOverlay.Active {
		return baseContent
	}

	// Create the loading dialog box - minimal, non-invasive design
	dialogContent := renderLoadingDialog(model.LoadingOverlay)

	// Create a compact dialog style with subtle appearance
	dialogWidth := 32 // Even smaller for minimal footprint
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#4B5563")). // More subtle gray
		Background(lipgloss.Color("#111827")).       // Darker, less prominent
		Foreground(lipgloss.Color("#D1D5DB")).       // Softer white
		Padding(1).                                  // Minimal padding
		Width(dialogWidth).
		Align(lipgloss.Center)

	// Apply the dialog style to the content
	styledDialog := dialogStyle.Render(dialogContent)

	// Position the dialog in the top-right corner to be less invasive
	positionedDialog := lipgloss.Place(
		width, height,
		lipgloss.Right, lipgloss.Top,
		styledDialog,
	)

	return positionedDialog
}

// renderLoadingDialog creates the content for the loading dialog
func renderLoadingDialog(overlay *types.LoadingOverlay) string {
	// Get spinner character
	spinnerChar := overlay.Spinner.GetSpinnerChar()

	// Create the dialog content - simplified and less prominent
	message := spinnerChar + " " + overlay.Message
	instruction := "ESC to cancel"

	// Join content with minimal spacing
	content := strings.Join([]string{
		message,
		"",
		instruction,
	}, "\n")

	return content
}

// GetLoadingMessages returns the appropriate loading messages for the given loading type
func GetLoadingMessages(loadingType types.LoadingType) []string {
	switch loadingType {
	case types.LoadingStartup:
		return []string{
			"Initializing MCP Manager...",
			"Loading MCP inventory...",
			"Detecting Claude CLI...",
			"Ready!",
		}
	case types.LoadingRefresh:
		return []string{
			"Refreshing MCP status...",
			"Syncing with Claude CLI...",
			"Updating display...",
			"Complete!",
		}
	default:
		return []string{"Loading..."}
	}
}