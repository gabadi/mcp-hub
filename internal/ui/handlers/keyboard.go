package handlers

import (
	"cc-mcp-manager/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

// Keyboard key constants
const (
	KeyTab        = "tab"
	KeyEsc        = "esc"
	KeyEnter      = "enter"
	KeyCtrlC      = "ctrl+c"
	KeyCtrlV      = "ctrl+v"
	KeyDown       = "down"
	KeyUp         = "up"
	KeyBackspace  = "backspace"
	KeyCmdC       = "cmd+c"
	KeyCmdV       = "cmd+v"
	KeyCmdSymbolC = "⌘c"
	KeyCmdSymbolV = "⌘v"
)

// HandleKeyPress processes keyboard input based on current state
func HandleKeyPress(model types.Model, msg tea.KeyMsg) (types.Model, tea.Cmd) {
	key := msg.String()

	// Handle special key types for search states
	if model.State == types.SearchActiveNavigation && model.SearchInputActive && msg.Type == tea.KeyRunes && len(msg.Runes) > 0 && key != KeyEnter && key != KeyTab && key != KeyEsc {
		// Add runes to search query (only if not a special command key)
		for _, r := range msg.Runes {
			model.SearchQuery += string(r)
		}
		return model, nil
	}

	// Global keys that work in any state
	switch key {
	case KeyEsc:
		return HandleEscKey(model)
	case KeyCtrlC:
		return model, tea.Quit
	case "ctrl+l":
		// Clear screen and redraw - just return nil cmd as the screen will auto-refresh
		return model, nil
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
