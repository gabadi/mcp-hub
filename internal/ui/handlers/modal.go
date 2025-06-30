package handlers

import (
	"cc-mcp-manager/internal/ui/types"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleModalKeys handles keyboard input in modal mode
func HandleModalKeys(model types.Model, key string) (types.Model, tea.Cmd) {
	switch key {
	case "enter":
		// Confirm modal action and return to main navigation
		model.State = types.MainNavigation
	}

	return model, nil
}
