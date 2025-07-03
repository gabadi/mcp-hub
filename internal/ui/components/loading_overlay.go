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

	// Create the loading dialog box
	dialogContent := renderLoadingDialog(model.LoadingOverlay)

	// Create dialog box style - less invasive design
	dialogWidth := 40 // Width of the dialog box (reduced from 60)
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#6B7280")).
		Background(lipgloss.Color("#1F2937")).
		Foreground(lipgloss.Color("#F3F4F6")).
		Padding(1, 2).
		Width(dialogWidth).
		Align(lipgloss.Center)

	// Apply the dialog style to the content
	styledDialog := dialogStyle.Render(dialogContent)

	// Position the dialog in the center of the screen with semi-transparent background
	positionedDialog := lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		styledDialog,
		lipgloss.WithWhitespaceBackground(lipgloss.Color("#1E1E2E")),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#6B7280")),
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