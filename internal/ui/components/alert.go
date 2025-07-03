package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderAlertOverlay renders a success/error alert as an overlay positioned at the top center
// without affecting the main content layout flow. This follows the same pattern as the modal
// overlay system to ensure consistent behavior.
func RenderAlertOverlay(message string, width, height int, backgroundContent string) string {
	// If no message, return background content unchanged
	if message == "" {
		return backgroundContent
	}

	// Alert styling - maintaining the existing green success style
	// Using exact same colors and styling as the original implementation
	alertStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#51CF66")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 2).
		Align(lipgloss.Center)

	// Calculate alert width to maintain 8px margin (4px on each side)
	// But ensure we don't exceed the terminal width
	alertWidth := width - 8
	if alertWidth < 10 {
		alertWidth = width - 2 // Use almost full width for very small terminals
	}
	if alertWidth < 5 {
		alertWidth = width // Use full width for extremely small terminals
	}
	alertStyle = alertStyle.Width(alertWidth)

	// Create the alert content
	alertContent := alertStyle.Render(message)

	// Position the alert as an overlay at top center with 4px top margin
	// This uses the same lipgloss.Place approach as the modal system
	alertOverlay := lipgloss.Place(
		width,
		height,
		lipgloss.Center, // Center horizontally
		lipgloss.Top,    // Position at top
		alertContent,
		lipgloss.WithWhitespaceChars(" "),
	)

	// The key innovation: we layer the overlay on the background by combining
	// their visual representations. The alert "floats" over the content.
	return combineOverlayWithBackground(backgroundContent, alertOverlay, width, height)
}

// combineOverlayWithBackground creates a true overlay effect by combining
// the background content with the overlay content, ensuring the alert
// appears "floating" over the main interface without disrupting layout
func combineOverlayWithBackground(background, overlay string, width, height int) string {
	// Split content into lines for manipulation
	backgroundLines := strings.Split(background, "\n")
	overlayLines := strings.Split(overlay, "\n")

	// Ensure we have enough lines to work with
	maxLines := height
	if len(backgroundLines) > maxLines {
		maxLines = len(backgroundLines)
	}
	if len(overlayLines) > maxLines {
		maxLines = len(overlayLines)
	}

	// Pad background lines to match dimensions
	for len(backgroundLines) < maxLines {
		backgroundLines = append(backgroundLines, strings.Repeat(" ", width))
	}

	// Pad overlay lines to match dimensions
	for len(overlayLines) < maxLines {
		overlayLines = append(overlayLines, strings.Repeat(" ", width))
	}

	// Combine the layers: where overlay has content, it takes precedence
	// where overlay is transparent/empty, background shows through
	result := make([]string, maxLines)
	for i := 0; i < maxLines; i++ {
		overlayLine := ""
		backgroundLine := ""

		if i < len(overlayLines) {
			overlayLine = overlayLines[i]
		}
		if i < len(backgroundLines) {
			backgroundLine = backgroundLines[i]
		}

		// If overlay line is effectively empty (just spaces), use background
		if strings.TrimSpace(overlayLine) == "" {
			result[i] = backgroundLine
		} else {
			result[i] = overlayLine
		}
	}

	return strings.Join(result, "\n")
}
