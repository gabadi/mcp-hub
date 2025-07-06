package handlers

import (
	"mcp-hub/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

// HandleEscKey handles ESC key behavior based on current state
func HandleEscKey(model types.Model) (types.Model, tea.Cmd) {
	// Priority 1: Check for active loading overlay
	if model.IsLoadingOverlayActive() {
		return handleLoadingCancellation(model)
	}

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
		// Clear edit mode state if canceling edit
		model.EditMode = false
		model.EditMCPName = ""
		// Clear form data and errors
		model.FormData = types.FormData{}
		model.FormErrors = make(map[string]string)
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

// handleLoadingCancellation handles ESC key during loading operations
func handleLoadingCancellation(model types.Model) (types.Model, tea.Cmd) {
	if model.LoadingOverlay == nil || !model.LoadingOverlay.Active {
		return model, nil
	}

	// For now, we'll implement simple cancellation without confirmation prompt
	// This can be enhanced later to show "Cancel operation? [Y/N]" prompt

	loadingType := model.LoadingOverlay.Type
	model.StopLoadingOverlay()

	switch loadingType {
	case types.LoadingStartup:
		// For startup cancellation, exit the application
		return model, tea.Quit
	case types.LoadingRefresh:
		// For refresh cancellation, return to current state
		model.SuccessMessage = "Refresh operation canceled"
		model.SuccessTimer = 120
		return model, TimerCmd("success_timer")
	case types.LoadingClaude:
		// For Claude sync cancellation, return to current state
		model.SuccessMessage = "Claude sync canceled"
		model.SuccessTimer = 120
		return model, TimerCmd("success_timer")
	default:
		// Unknown loading type, just stop loading
		return model, nil
	}
}

// HandleSearchModeKeys handles keyboard input in search mode
func HandleSearchModeKeys(model types.Model, key string) (types.Model, tea.Cmd) {
	switch key {
	case KeyEnter:
		// Return to main navigation with search query preserved
		model.State = types.MainNavigation
		model.SearchActive = false
	case "backspace":
		if len(model.SearchQuery) > 0 {
			model.SearchQuery = model.SearchQuery[:len(model.SearchQuery)-1]
		}
	case KeyCtrlV, "cmd+v", KeyCmdSymbolV, "command+v":
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
