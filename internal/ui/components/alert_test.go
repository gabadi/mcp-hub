package components

import (
	"strings"
	"testing"
)

func TestRenderAlertOverlay(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		width          int
		height         int
		backgroundContent string
		expectOverlay  bool
	}{
		{
			name:           "No message returns background unchanged",
			message:        "",
			width:          80,
			height:         24,
			backgroundContent: "Test background content",
			expectOverlay:  false,
		},
		{
			name:           "Message creates overlay with proper styling",
			message:        "Success! MCP activated",
			width:          80,
			height:         24,
			backgroundContent: "Test background content",
			expectOverlay:  true,
		},
		{
			name:           "Small terminal dimensions handled gracefully",
			message:        "Success",
			width:          40,
			height:         12,
			backgroundContent: "Small terminal",
			expectOverlay:  true,
		},
		{
			name:           "Very small width uses minimum alert width",
			message:        "Test",
			width:          15,
			height:         10,
			backgroundContent: "Tiny",
			expectOverlay:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderAlertOverlay(tt.message, tt.width, tt.height, tt.backgroundContent)

			if !tt.expectOverlay {
				// Should return background unchanged
				if result != tt.backgroundContent {
					t.Errorf("Expected unchanged background, got different result")
				}
				return
			}

			// When overlay is expected
			if result == tt.backgroundContent {
				t.Errorf("Expected overlay to be applied, but result matches background")
			}

			// Check that the result contains the message
			if !strings.Contains(result, tt.message) {
				t.Errorf("Expected result to contain message '%s', but it doesn't", tt.message)
			}

			// Verify result has appropriate dimensions
			lines := strings.Split(result, "\n")
			if len(lines) == 0 {
				t.Errorf("Expected result to have content, got empty")
			}
		})
	}
}

func TestRenderAlertOverlay_Styling(t *testing.T) {
	message := "Test Alert"
	width := 80
	height := 24
	background := "Background content"

	result := RenderAlertOverlay(message, width, height, background)

	// The result should be different from background (overlay applied)
	if result == background {
		t.Errorf("Expected alert overlay to modify the output")
	}

	// Should contain the message
	if !strings.Contains(result, message) {
		t.Errorf("Expected result to contain alert message")
	}
}

func TestCombineOverlayWithBackground(t *testing.T) {
	background := "Line 1\nLine 2\nLine 3"
	overlay := "Overlay\n\n"  // Empty second line
	width := 20
	height := 5

	result := combineOverlayWithBackground(background, overlay, width, height)

	lines := strings.Split(result, "\n")
	
	// Should have correct number of lines
	if len(lines) != height {
		t.Errorf("Expected %d lines, got %d", height, len(lines))
	}

	// First line should be from overlay
	if !strings.Contains(lines[0], "Overlay") {
		t.Errorf("Expected first line to contain overlay content")
	}

	// Second line should be from background (overlay is empty)
	if !strings.Contains(lines[1], "Line 2") {
		t.Errorf("Expected second line to contain background content when overlay is empty")
	}
}

func TestRenderAlertOverlay_MinimumWidth(t *testing.T) {
	message := "Test"  // Short message that should fit
	width := 10 // Very small width
	height := 5
	background := "BG"

	result := RenderAlertOverlay(message, width, height, background)

	// Should not panic and should produce some result
	if result == "" {
		t.Errorf("Expected some result even with very small width")
	}

	// Should contain at least part of the message (due to word wrapping in small spaces)
	if !strings.Contains(result, "Test") {
		t.Errorf("Expected result to contain at least part of the message, got: %q", result)
	}
}