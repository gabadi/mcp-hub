package handlers

import (
	"strings"

	"cc-mcp-manager/internal/ui/services"
	"cc-mcp-manager/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

// HandleMainNavigationKeys handles keyboard input in main navigation mode
func HandleMainNavigationKeys(model types.Model, key string) (types.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		model = NavigateUp(model)
	case "down", "j":
		model = NavigateDown(model)
	case "left", "h":
		model = NavigateLeft(model)
	case "right", "l":
		model = NavigateRight(model)
	case "tab":
		// Jump to search field with navigation enabled
		model.State = types.SearchActiveNavigation
		model.SearchActive = true
		model.SearchInputActive = true
		model.SelectedItem = 0 // Reset selection to first item
		model.FilteredSelectedIndex = 0
	case "/":
		// Activate search mode with navigation enabled
		model.State = types.SearchActiveNavigation
		model.SearchActive = true
		model.SearchInputActive = true
		model.SelectedItem = 0 // Reset selection to first item
		model.FilteredSelectedIndex = 0
		// Don't add the "/" character to the search query
	case "a":
		// Add MCP - Start with type selection
		model.State = types.ModalActive
		model.ActiveModal = types.AddMCPTypeSelection
		// Reset form data
		model.FormData = types.FormData{}
		model.FormErrors = make(map[string]string)
	case "e":
		// Edit MCP (future functionality)
		model.State = types.ModalActive
		model.ActiveModal = types.EditModal
	case "d":
		// Delete MCP (future functionality)
		model.State = types.ModalActive
		model.ActiveModal = types.DeleteModal
	case " ", "space":
		// Toggle MCP active status
		model = services.ToggleMCPStatus(model)
	}

	return model, nil
}

// HandleSearchNavigationKeys handles keyboard input in search + navigation mode
func HandleSearchNavigationKeys(model types.Model, key string) (types.Model, tea.Cmd) {
	// Priority 1: Navigation keys (always work)
	switch key {
	case "up", "k":
		model = NavigateUp(model)
		return model, nil
	case "down", "j":
		model = NavigateDown(model)
		return model, nil
	case "left", "h":
		model = NavigateLeft(model)
		return model, nil
	case "right", "l":
		model = NavigateRight(model)
		return model, nil
	case " ", "space":
		// Toggle MCP active status
		model = services.ToggleMCPStatus(model)
		return model, nil
	}

	// Priority 2: Search control keys (mode switching)
	switch key {
	case "enter":
		// Return to main navigation with search query preserved
		model.State = types.MainNavigation
		model.SearchActive = false
		model.SearchInputActive = false
		return model, nil
	case "tab":
		// Toggle between input and navigation modes
		model.SearchInputActive = !model.SearchInputActive
		return model, nil
	}

	// Priority 3: Text input (only when searchInputActive = true)
	if model.SearchInputActive {
		switch key {
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
	}

	return model, nil
}

// Navigation helper methods

// NavigateUp moves selection up within the grid
func NavigateUp(model types.Model) types.Model {
	if model.ColumnCount == 4 {
		// In 4-column grid - move up one row (subtract 4 from index)
		filteredMCPs := services.GetFilteredMCPs(model)

		// Use appropriate index based on search state
		if model.SearchQuery != "" {
			if model.FilteredSelectedIndex >= 4 {
				model.FilteredSelectedIndex -= 4
			}
			// Ensure we don't go below 0 and stay within filtered results
			if model.FilteredSelectedIndex < 0 {
				model.FilteredSelectedIndex = 0
			}
			if model.FilteredSelectedIndex >= len(filteredMCPs) && len(filteredMCPs) > 0 {
				model.FilteredSelectedIndex = len(filteredMCPs) - 1
			}
		} else {
			if model.SelectedItem >= 4 {
				model.SelectedItem -= 4
			}
			// Ensure we don't go below 0
			if model.SelectedItem < 0 {
				model.SelectedItem = 0
			}
			if model.SelectedItem >= len(filteredMCPs) && len(filteredMCPs) > 0 {
				model.SelectedItem = len(filteredMCPs) - 1
			}
		}
	} else if model.ActiveColumn == 0 {
		// In MCP list column for other layouts
		if model.SelectedItem > 0 {
			model.SelectedItem--
		}
	}
	// Other columns don't have navigable items yet
	return model
}

// NavigateDown moves selection down within the grid
func NavigateDown(model types.Model) types.Model {
	if model.ColumnCount == 4 {
		// In 4-column grid - move down one row (add 4 to index)
		filteredMCPs := services.GetFilteredMCPs(model)

		// Use appropriate index based on search state
		if model.SearchQuery != "" {
			newIndex := model.FilteredSelectedIndex + 4
			if newIndex < len(filteredMCPs) {
				model.FilteredSelectedIndex = newIndex
			}
			// Stay within filtered results bounds
			if model.FilteredSelectedIndex >= len(filteredMCPs) && len(filteredMCPs) > 0 {
				model.FilteredSelectedIndex = len(filteredMCPs) - 1
			}
		} else {
			newIndex := model.SelectedItem + 4
			if newIndex < len(filteredMCPs) {
				model.SelectedItem = newIndex
			}
			// Stay within filtered results bounds
			if model.SelectedItem >= len(filteredMCPs) && len(filteredMCPs) > 0 {
				model.SelectedItem = len(filteredMCPs) - 1
			}
		}
	} else if model.ActiveColumn == 0 {
		// In MCP list column for other layouts
		maxItems := len(model.MCPItems) - 1
		if model.SelectedItem < maxItems {
			model.SelectedItem++
		}
	}
	// Other columns don't have navigable items yet
	return model
}

// NavigateLeft moves to the left within the grid
func NavigateLeft(model types.Model) types.Model {
	if model.ColumnCount == 4 {
		// In 4-column grid - move left within current row
		if model.SearchQuery != "" {
			if model.FilteredSelectedIndex%4 > 0 {
				model.FilteredSelectedIndex--
			}
		} else {
			if model.SelectedItem%4 > 0 {
				model.SelectedItem--
			}
		}
	} else {
		// For other layouts, move between columns
		if model.ActiveColumn > 0 {
			model.ActiveColumn--
		}
	}
	return model
}

// NavigateRight moves to the right within the grid
func NavigateRight(model types.Model) types.Model {
	if model.ColumnCount == 4 {
		// In 4-column grid - move right within current row
		filteredMCPs := services.GetFilteredMCPs(model)

		if model.SearchQuery != "" {
			if model.FilteredSelectedIndex%4 < 3 && model.FilteredSelectedIndex+1 < len(filteredMCPs) {
				model.FilteredSelectedIndex++
			}
		} else {
			if model.SelectedItem%4 < 3 && model.SelectedItem+1 < len(filteredMCPs) {
				model.SelectedItem++
			}
		}
	} else {
		// For other layouts, move between columns
		if model.ActiveColumn < model.ColumnCount-1 {
			model.ActiveColumn++
		}
	}
	return model
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
