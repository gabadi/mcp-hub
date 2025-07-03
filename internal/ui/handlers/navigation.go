package handlers

import (
	"context"
	"strings"
	"time"

	"cc-mcp-manager/internal/ui/services"
	"cc-mcp-manager/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

// HandleMainNavigationKeys handles keyboard input in main navigation mode
func HandleMainNavigationKeys(model types.Model, key string) (types.Model, tea.Cmd) {
	// Handle navigation keys
	if updatedModel, handled := handleNavigationKeys(model, key); handled {
		return updatedModel, nil
	}

	// Handle search keys
	if updatedModel, handled := handleSearchKeys(model, key); handled {
		return updatedModel, nil
	}

	// Handle action keys
	if updatedModel, cmd, handled := handleActionKeys(model, key); handled {
		return updatedModel, cmd
	}

	return model, nil
}

// handleNavigationKeys handles directional navigation keys
func handleNavigationKeys(model types.Model, key string) (types.Model, bool) {
	switch key {
	case "up", "k":
		return NavigateUp(model), true
	case "down", "j":
		return NavigateDown(model), true
	case "left", "h":
		return NavigateLeft(model), true
	case "right", "l":
		return NavigateRight(model), true
	}
	return model, false
}

// handleSearchKeys handles search activation keys
func handleSearchKeys(model types.Model, key string) (types.Model, bool) {
	switch key {
	case "tab", "/":
		// Activate search mode with navigation enabled
		model.State = types.SearchActiveNavigation
		model.SearchActive = true
		model.SearchInputActive = true
		model.SelectedItem = 0 // Reset selection to first item
		model.FilteredSelectedIndex = 0
		return model, true
	}
	return model, false
}

// handleActionKeys handles action keys (add, edit, delete, toggle, refresh)
func handleActionKeys(model types.Model, key string) (types.Model, tea.Cmd, bool) {
	switch key {
	case "a":
		return handleAddMCP(model), nil, true
	case "e":
		return handleEditMCP(model)
	case "d":
		return handleDeleteMCP(model), nil, true
	case " ", "space":
		return handleEnhancedToggleMCP(model)
	case "r", "R":
		return model, RefreshClaudeStatusCmd(), true
	}
	return model, nil, false
}

// handleEnhancedToggleMCP handles the enhanced MCP toggle operation (Epic 2 Story 2)
func handleEnhancedToggleMCP(model types.Model) (types.Model, tea.Cmd, bool) {
	selectedMCP := services.GetSelectedMCP(model)
	if selectedMCP == nil {
		return model, nil, true
	}

	// Set immediate loading state for visual feedback
	updatedModel := services.ToggleMCPStatus(model)

	// If already in error state due to Claude unavailability, don't proceed
	if updatedModel.ToggleState == types.ToggleError {
		return updatedModel, nil, true
	}

	// Create command to perform the actual toggle operation
	activate := !selectedMCP.Active
	cmd := EnhancedToggleMCPCmd(selectedMCP.Name, activate, selectedMCP)

	return updatedModel, cmd, true
}

// handleAddMCP handles the add MCP action
func handleAddMCP(model types.Model) types.Model {
	model.State = types.ModalActive
	model.ActiveModal = types.AddMCPTypeSelection
	model.FormData = types.FormData{}
	model.FormErrors = make(map[string]string)
	return model
}

// handleEditMCP handles the edit MCP action
func handleEditMCP(model types.Model) (types.Model, tea.Cmd, bool) {
	selectedMCP := services.GetSelectedMCP(model)
	if selectedMCP == nil {
		return model, nil, true
	}

	model.State = types.ModalActive
	model.ActiveModal = getEditModalType(selectedMCP.Type)
	model.FormData = populateFormDataFromMCP(*selectedMCP)
	model.FormErrors = make(map[string]string)
	model.EditMode = true
	model.EditMCPName = selectedMCP.Name

	return model, nil, true
}

// handleDeleteMCP handles the delete MCP action
func handleDeleteMCP(model types.Model) types.Model {
	model.State = types.ModalActive
	model.ActiveModal = types.DeleteModal
	return model
}

// getEditModalType returns the appropriate modal type for editing based on MCP type
func getEditModalType(mcpType string) types.ModalType {
	switch mcpType {
	case "CMD":
		return types.AddCommandForm
	case "SSE":
		return types.AddSSEForm
	case "JSON":
		return types.AddJSONForm
	default:
		return types.AddCommandForm // Default fallback
	}
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
		// Enhanced toggle MCP active status
		updatedModel, cmd, _ := handleEnhancedToggleMCP(model)
		return updatedModel, cmd
	case "r", "R":
		// Refresh Claude status
		return model, RefreshClaudeStatusCmd()
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
			// Only move if there's actually an item 4 positions down
			if newIndex < len(filteredMCPs) {
				model.FilteredSelectedIndex = newIndex
			}
		} else {
			newIndex := model.SelectedItem + 4
			// Only move if there's actually an item 4 positions down
			if newIndex < len(filteredMCPs) {
				model.SelectedItem = newIndex
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

// populateFormDataFromMCP converts an MCPItem to FormData for editing (Epic 1 Story 4)
func populateFormDataFromMCP(mcp types.MCPItem) types.FormData {
	formData := types.FormData{
		Name:        mcp.Name,
		Command:     mcp.Command,
		URL:         mcp.URL,
		JSONConfig:  mcp.JSONConfig,
		ActiveField: 0, // Start with first field focused
	}

	// Convert Args slice to string
	if len(mcp.Args) > 0 {
		formData.Args = formatArgsForDisplay(mcp.Args)
	}

	// Convert Environment map to string
	if len(mcp.Environment) > 0 {
		formData.Environment = formatEnvironmentForDisplay(mcp.Environment)
	}

	return formData
}

// formatArgsForDisplay converts []string to display format (Epic 1 Story 4)
func formatArgsForDisplay(args []string) string {
	if len(args) == 0 {
		return ""
	}

	// Join with spaces, quoting arguments that contain spaces
	var formattedArgs []string
	for _, arg := range args {
		if strings.Contains(arg, " ") {
			formattedArgs = append(formattedArgs, `"`+arg+`"`)
		} else {
			formattedArgs = append(formattedArgs, arg)
		}
	}

	return strings.Join(formattedArgs, " ")
}

// formatEnvironmentForDisplay converts map[string]string to display format (Epic 1 Story 4)
func formatEnvironmentForDisplay(env map[string]string) string {
	if len(env) == 0 {
		return ""
	}

	var pairs []string
	for key, value := range env {
		pairs = append(pairs, key+"="+value)
	}

	return strings.Join(pairs, ",")
}

// ClaudeStatusMsg represents a Claude status update message (Epic 2 Story 1)
type ClaudeStatusMsg struct {
	Status types.ClaudeStatus
}

// ToggleResultMsg represents a toggle operation result message (Epic 2 Story 2)
type ToggleResultMsg struct {
	MCPName  string
	Activate bool
	Success  bool
	Error    string
	Retrying bool
}

// RefreshClaudeStatusCmd creates a command to refresh Claude status (Epic 2 Story 1)
func RefreshClaudeStatusCmd() tea.Cmd {
	return func() tea.Msg {
		claudeService := services.NewClaudeService()
		ctx := context.Background()
		status := claudeService.RefreshClaudeStatus(ctx)
		return ClaudeStatusMsg{Status: status}
	}
}

// EnhancedToggleMCPCmd creates a command to perform enhanced MCP toggle (Epic 2 Story 2)
func EnhancedToggleMCPCmd(mcpName string, activate bool, mcpConfig *types.MCPItem) tea.Cmd {
	return func() tea.Msg {
		claudeService := services.NewClaudeService()
		ctx := context.Background()

		// Pass the MCP configuration for add operations
		result, err := claudeService.ToggleMCPStatus(ctx, mcpName, activate, mcpConfig)

		if err != nil {
			return ToggleResultMsg{
				MCPName:  mcpName,
				Activate: activate,
				Success:  false,
				Error:    "Internal error during MCP toggle operation",
				Retrying: false,
			}
		}

		return ToggleResultMsg{
			MCPName:  mcpName,
			Activate: activate,
			Success:  result.Success,
			Error:    result.ErrorMsg,
			Retrying: result.Retryable && !result.Success,
		}
	}
}

// TimerCmd creates a command that sends timer tick messages every ~50ms
func TimerCmd(timerID string) tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(time.Time) tea.Msg {
		return types.TimerTickMsg{ID: timerID}
	})
}
