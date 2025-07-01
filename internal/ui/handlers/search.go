package handlers

import (
	"strings"

	"cc-mcp-manager/internal/ui/services"
	"cc-mcp-manager/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

// HandleEscKey handles ESC key behavior based on current state
func HandleEscKey(model types.Model) (types.Model, tea.Cmd) {
	switch model.State {
	case types.SearchMode:
		// Clear search and return to main navigation
		model.SearchActive = false
		model.SearchQuery = ""
		model.FilteredSelectedIndex = 0
		model.State = types.MainNavigation
		return model, nil
	case types.SearchActiveNavigation:
		// Clear search and return to main navigation
		model.SearchActive = false
		model.SearchInputActive = false
		model.SearchQuery = ""
		model.FilteredSelectedIndex = 0
		model.State = types.MainNavigation
		return model, nil
	case types.ModalActive:
		// Close modal and return to main navigation
		model.State = types.MainNavigation
		model.ActiveModal = types.NoModal
		return model, nil
	case types.MainNavigation:
		// Clear search if active, otherwise exit application
		if model.SearchQuery != "" {
			model.SearchQuery = ""
			model.SearchResults = nil
			model.SelectedItem = 0 // Reset selection
			model.FilteredSelectedIndex = 0
			return model, nil
		}
		// Exit application
		return model, tea.Quit
	}
	return model, nil
}

// HandleSearchModeKeys handles keyboard input in search mode
func HandleSearchModeKeys(model types.Model, key string) (types.Model, tea.Cmd) {
	switch key {
	case "enter":
		// Return to main navigation with search query preserved
		model.State = types.MainNavigation
		model.SearchActive = false
	case "backspace":
		if len(model.SearchQuery) > 0 {
			model.SearchQuery = model.SearchQuery[:len(model.SearchQuery)-1]
		}
	case "ctrl+v", "cmd+v", "⌘v", "command+v":
		// Paste clipboard content to search query
		model = pasteToSearchQuery(model)
	default:
		// Add character to search query
		if len(key) == 1 {
			model.SearchQuery += key
		}
	}

	return model, nil
}

// pasteToSearchQuery pastes clipboard content to the search query
func pasteToSearchQuery(model types.Model) types.Model {
	clipboardService := services.NewClipboardService()

	// Use enhanced paste for better error diagnostics
	content, err := clipboardService.EnhancedPaste()
	if err != nil {
		// Add user feedback for clipboard paste failure with enhanced error information
		model.SuccessMessage = "Failed to paste from clipboard: " + err.Error()
		model.SuccessTimer = 240 // Show error message for 4 seconds to allow reading detailed error
		return model
	}

	// Clean the pasted content - remove newlines and control characters
	cleanContent := strings.ReplaceAll(content, "\n", " ")
	cleanContent = strings.ReplaceAll(cleanContent, "\r", " ")
	cleanContent = strings.ReplaceAll(cleanContent, "\t", " ")
	// Remove multiple spaces
	cleanContent = strings.Join(strings.Fields(cleanContent), " ")

	// Append to existing search query
	if model.SearchQuery != "" && !strings.HasSuffix(model.SearchQuery, " ") {
		model.SearchQuery += " "
	}
	model.SearchQuery += cleanContent

	// Add success feedback for successful paste operation
	model.SuccessMessage = "Pasted from clipboard"
	model.SuccessTimer = 120 // Show success message for 2 seconds

	return model
}
