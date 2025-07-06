package components

import (
	"strings"
	"testing"

	"mcp-hub/internal/ui/types"
)

func TestRenderLoadingOverlay(t *testing.T) {
	// Test case 1: No loading overlay active
	model := types.Model{
		LoadingOverlay: nil,
	}
	baseContent := "Base Content"
	result := RenderLoadingOverlay(model, 80, 24, baseContent)

	if result != baseContent {
		t.Errorf("Expected base content when no loading overlay, got: %s", result)
	}

	// Test case 2: Loading overlay inactive
	model.LoadingOverlay = &types.LoadingOverlay{
		Active: false,
	}
	result = RenderLoadingOverlay(model, 80, 24, baseContent)

	if result != baseContent {
		t.Errorf("Expected base content when loading overlay inactive, got: %s", result)
	}

	// Test case 3: Loading overlay active
	model.LoadingOverlay = &types.LoadingOverlay{
		Active:      true,
		Message:     "Test loading message",
		Spinner:     types.SpinnerFrame1,
		Cancellable: true,
		Type:        types.LoadingStartup,
	}
	result = RenderLoadingOverlay(model, 80, 24, baseContent)

	if result == baseContent {
		t.Error("Expected modified content when loading overlay active")
	}

	// Check that the result contains the loading message
	if !strings.Contains(result, "Test loading message") {
		t.Error("Expected loading message to be present in result")
	}
}

func TestRenderLoadingDialog(t *testing.T) {
	overlay := &types.LoadingOverlay{
		Active:      true,
		Message:     "Test message",
		Spinner:     types.SpinnerFrame1,
		Cancellable: true,
		Type:        types.LoadingStartup,
	}

	result := renderLoadingDialog(overlay)

	// Check that the dialog contains expected elements (updated for simplified UX)
	if !strings.Contains(result, "Test message") {
		t.Error("Expected dialog to contain message")
	}

	if !strings.Contains(result, "ESC to cancel") {
		t.Error("Expected dialog to contain cancel instruction")
	}

	// Check that spinner character is included
	spinnerChar := overlay.Spinner.GetSpinnerChar()
	if !strings.Contains(result, spinnerChar) {
		t.Error("Expected dialog to contain spinner character")
	}
}

func TestGetLoadingMessages(t *testing.T) {
	// Test startup messages
	startupMessages := GetLoadingMessages(types.LoadingStartup)
	expectedStartup := []string{
		"Initializing MCP Manager...",
		"Loading MCP inventory...",
		"Detecting Claude CLI...",
		"Ready!",
	}

	if len(startupMessages) != len(expectedStartup) {
		t.Errorf("Expected %d startup messages, got %d", len(expectedStartup), len(startupMessages))
	}

	for i, expected := range expectedStartup {
		if startupMessages[i] != expected {
			t.Errorf("Expected startup message %d to be '%s', got '%s'", i, expected, startupMessages[i])
		}
	}

	// Test refresh messages
	refreshMessages := GetLoadingMessages(types.LoadingRefresh)
	expectedRefresh := []string{
		"Refreshing MCP status...",
		"Syncing with Claude CLI...",
		"Updating display...",
		"Complete!",
	}

	if len(refreshMessages) != len(expectedRefresh) {
		t.Errorf("Expected %d refresh messages, got %d", len(expectedRefresh), len(refreshMessages))
	}

	for i, expected := range expectedRefresh {
		if refreshMessages[i] != expected {
			t.Errorf("Expected refresh message %d to be '%s', got '%s'", i, expected, refreshMessages[i])
		}
	}

	// Test default case
	defaultMessages := GetLoadingMessages(types.LoadingType(999))
	if len(defaultMessages) != 1 || defaultMessages[0] != "Loading..." {
		t.Error("Expected default loading message")
	}
}

func TestSpinnerStateGetSpinnerChar(t *testing.T) {
	tests := []struct {
		state    types.SpinnerState
		expected string
	}{
		{types.SpinnerFrame1, "◐"},
		{types.SpinnerFrame2, "◓"},
		{types.SpinnerFrame3, "◑"},
		{types.SpinnerFrame4, "◒"},
	}

	for _, test := range tests {
		result := test.state.GetSpinnerChar()
		if result != test.expected {
			t.Errorf("Expected spinner state %d to return '%s', got '%s'", test.state, test.expected, result)
		}
	}

	// Test default case
	defaultState := types.SpinnerState(999)
	if defaultState.GetSpinnerChar() != "◐" {
		t.Error("Expected default spinner state to return '◐'")
	}
}
