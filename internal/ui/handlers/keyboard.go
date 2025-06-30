package handlers

import (
	"cc-mcp-manager/internal/ui/types"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleKeyPress processes keyboard input based on current state
func HandleKeyPress(model types.Model, msg tea.KeyMsg) (types.Model, tea.Cmd) {
	key := msg.String()

	// Global keys that work in any state
	switch key {
	case "esc":
		return HandleEscKey(model)
	case "ctrl+c":
		return model, tea.Quit
	}

	// State-specific key handling
	switch model.State {
	case types.MainNavigation:
		return HandleMainNavigationKeys(model, key)
	case types.SearchMode:
		return HandleSearchModeKeys(model, key)
	case types.SearchActiveNavigation:
		return HandleSearchNavigationKeys(model, key)
	case types.ModalActive:
		return HandleModalKeys(model, key)
	}

	return model, nil
}
