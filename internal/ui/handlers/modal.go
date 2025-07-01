package handlers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"cc-mcp-manager/internal/ui/services"
	"cc-mcp-manager/internal/ui/types"
)

// SuccessMsg represents a success message
type SuccessMsg struct {
	Message string
}

// hideSuccessMsg returns a command that sends a message to hide the success message after a delay
func hideSuccessMsg() tea.Cmd {
	return tea.Tick(time.Second*3, func(time.Time) tea.Msg {
		return SuccessMsg{Message: ""}
	})
}

// HandleModalKeys handles keyboard input in modal mode
func HandleModalKeys(model types.Model, key string) (types.Model, tea.Cmd) {
	switch model.ActiveModal {
	case types.AddMCPTypeSelection:
		return handleTypeSelectionKeys(model, key)
	case types.AddCommandForm:
		return handleCommandFormKeys(model, key)
	case types.AddSSEForm:
		return handleSSEFormKeys(model, key)
	case types.AddJSONForm:
		return handleJSONFormKeys(model, key)
	case types.DeleteModal:
		return handleDeleteModalKeys(model, key)
	default:
		// Legacy modal handling
		switch key {
		case "enter":
			// Confirm modal action and return to main navigation
			model.State = types.MainNavigation
		}
	}

	return model, nil
}

// handleTypeSelectionKeys handles keyboard input in the type selection modal
func handleTypeSelectionKeys(model types.Model, key string) (types.Model, tea.Cmd) {
	// Initialize selected option if not set
	if model.FormData.ActiveField == 0 {
		model.FormData.ActiveField = 1 // Start with option 1 selected
	}

	switch key {
	case "1":
		// Command/Binary MCP type
		model.ActiveModal = types.AddCommandForm
		model.FormData.ActiveField = 0 // Focus on first field (Name)
	case "2":
		// SSE Server MCP type
		model.ActiveModal = types.AddSSEForm
		model.FormData.ActiveField = 0 // Focus on first field (Name)
	case "3":
		// JSON Configuration MCP type
		model.ActiveModal = types.AddJSONForm
		model.FormData.ActiveField = 0 // Focus on first field (Name)
	case "up", "k":
		// Navigate up in type selection
		if model.FormData.ActiveField > 1 {
			model.FormData.ActiveField--
		}
	case "down", "j":
		// Navigate down in type selection
		if model.FormData.ActiveField < 3 {
			model.FormData.ActiveField++
		}
	case "enter":
		// Select the currently highlighted option
		switch model.FormData.ActiveField {
		case 1:
			model.ActiveModal = types.AddCommandForm
			model.FormData.ActiveField = 0
		case 2:
			model.ActiveModal = types.AddSSEForm
			model.FormData.ActiveField = 0
		case 3:
			model.ActiveModal = types.AddJSONForm
			model.FormData.ActiveField = 0
		}
	}
	return model, nil
}

// handleCommandFormKeys handles keyboard input in the Command/Binary form
func handleCommandFormKeys(model types.Model, key string) (types.Model, tea.Cmd) {
	switch key {
	case "tab":
		// Move to next field
		model.FormData.ActiveField = (model.FormData.ActiveField + 1) % 4 // 4 fields: Name, Command, Args, Environment
	case "enter":
		// Submit form if valid
		var valid bool
		model, valid = validateCommandForm(model)
		if !valid {
			// Focus on first field with error
			model = focusOnFirstErrorField(model)
		} else {
			// Parse args from string to []string
			args := parseArgsString(model.FormData.Args)

			// Parse environment variables from string to map[string]string
			env := parseEnvironmentString(model.FormData.Environment)

			mcpItem := types.MCPItem{
				Name:        model.FormData.Name,
				Type:        "CMD",
				Active:      false,
				Command:     model.FormData.Command,
				Args:        args,
				Environment: env,
			}

			var cmd tea.Cmd
			if model.EditMode {
				// Update existing MCP
				model, cmd = updateMCPInInventory(model, mcpItem)
			} else {
				// Add new MCP
				model, cmd = addMCPToInventory(model, mcpItem)
			}

			// Close modal and return to main navigation
			model.State = types.MainNavigation
			model.ActiveModal = types.NoModal
			model.FormData = types.FormData{}
			model.FormErrors = make(map[string]string)
			model.EditMode = false
			model.EditMCPName = ""
			return model, cmd
		}
	case "backspace":
		// Delete character from active field
		model = deleteCharFromActiveField(model)
	case "ctrl+c", "cmd+c", "⌘c", "command+c":
		// Copy active field content to clipboard
		model = copyActiveFieldToClipboard(model)
	case "ctrl+v", "cmd+v", "⌘v", "command+v":
		// Paste clipboard content to active field
		model = pasteFromClipboardToActiveField(model)
	default:
		// Add character to active field
		if len(key) == 1 {
			model = addCharToActiveField(model, key)
		}
	}

	// Validate form after each change
	model, _ = validateCommandForm(model)

	return model, nil
}

// handleSSEFormKeys handles keyboard input in the SSE Server form
func handleSSEFormKeys(model types.Model, key string) (types.Model, tea.Cmd) {
	switch key {
	case "tab":
		// Move to next field
		model.FormData.ActiveField = (model.FormData.ActiveField + 1) % 3 // 3 fields: Name, URL, Environment
	case "enter":
		// Submit form if valid
		var valid bool
		model, valid = validateSSEForm(model)
		if !valid {
			// Focus on first field with error
			model = focusOnFirstErrorField(model)
		} else {
			// Parse environment variables from string to map[string]string
			env := parseEnvironmentString(model.FormData.Environment)

			mcpItem := types.MCPItem{
				Name:        model.FormData.Name,
				Type:        "SSE",
				Active:      false,
				URL:         model.FormData.URL,
				Environment: env,
			}

			var cmd tea.Cmd
			if model.EditMode {
				// Update existing MCP
				model, cmd = updateMCPInInventory(model, mcpItem)
			} else {
				// Add new MCP
				model, cmd = addMCPToInventory(model, mcpItem)
			}

			// Close modal and return to main navigation
			model.State = types.MainNavigation
			model.ActiveModal = types.NoModal
			model.FormData = types.FormData{}
			model.FormErrors = make(map[string]string)
			model.EditMode = false
			model.EditMCPName = ""
			return model, cmd
		}
	case "backspace":
		// Delete character from active field
		model = deleteCharFromActiveField(model)
	case "ctrl+c", "cmd+c", "⌘c", "command+c":
		// Copy active field content to clipboard
		model = copyActiveFieldToClipboard(model)
	case "ctrl+v", "cmd+v", "⌘v", "command+v":
		// Paste clipboard content to active field
		model = pasteFromClipboardToActiveField(model)
	default:
		// Add character to active field
		if len(key) == 1 {
			model = addCharToActiveField(model, key)
		}
	}

	// Validate form after each change
	model, _ = validateSSEForm(model)

	return model, nil
}

// handleJSONFormKeys handles keyboard input in the JSON Configuration form
func handleJSONFormKeys(model types.Model, key string) (types.Model, tea.Cmd) {
	switch key {
	case "tab":
		// Move to next field
		model.FormData.ActiveField = (model.FormData.ActiveField + 1) % 3 // 3 fields: Name, JSONConfig, Environment
	case "enter":
		// Submit form if valid (or newline in JSON field)
		if model.FormData.ActiveField == 1 { // JSON field
			model = addCharToActiveField(model, "\n")
		} else {
			var valid bool
			model, valid = validateJSONForm(model)
			if !valid {
				// Focus on first field with error
				model = focusOnFirstErrorField(model)
			} else {
				// Parse environment variables from string to map[string]string
				env := parseEnvironmentString(model.FormData.Environment)

				mcpItem := types.MCPItem{
					Name:        model.FormData.Name,
					Type:        "JSON",
					Active:      false,
					JSONConfig:  model.FormData.JSONConfig,
					Environment: env,
				}

				var cmd tea.Cmd
				if model.EditMode {
					// Update existing MCP
					model, cmd = updateMCPInInventory(model, mcpItem)
				} else {
					// Add new MCP
					model, cmd = addMCPToInventory(model, mcpItem)
				}

				// Close modal and return to main navigation
				model.State = types.MainNavigation
				model.ActiveModal = types.NoModal
				model.FormData = types.FormData{}
				model.FormErrors = make(map[string]string)
				model.EditMode = false
				model.EditMCPName = ""
				return model, cmd
			}
		}
	case "backspace":
		// Delete character from active field
		model = deleteCharFromActiveField(model)
	case "ctrl+c", "cmd+c", "⌘c", "command+c":
		// Copy active field content to clipboard
		model = copyActiveFieldToClipboard(model)
	case "ctrl+v", "cmd+v", "⌘v", "command+v":
		// Paste clipboard content to active field
		model = pasteFromClipboardToActiveField(model)
	default:
		// Add character to active field
		if len(key) == 1 {
			model = addCharToActiveField(model, key)
		}
	}

	// Validate form after each change
	model, _ = validateJSONForm(model)

	return model, nil
}

// Helper functions for form field manipulation
func addCharToActiveField(model types.Model, char string) types.Model {
	switch model.ActiveModal {
	case types.AddCommandForm:
		switch model.FormData.ActiveField {
		case 0:
			model.FormData.Name += char
		case 1:
			model.FormData.Command += char
		case 2:
			model.FormData.Args += char
		case 3:
			model.FormData.Environment += char
		}
	case types.AddSSEForm:
		switch model.FormData.ActiveField {
		case 0:
			model.FormData.Name += char
		case 1:
			model.FormData.URL += char
		case 2:
			model.FormData.Environment += char
		}
	case types.AddJSONForm:
		switch model.FormData.ActiveField {
		case 0:
			model.FormData.Name += char
		case 1:
			model.FormData.JSONConfig += char
		case 2:
			model.FormData.Environment += char
		}
	}
	return model
}

func deleteCharFromActiveField(model types.Model) types.Model {
	switch model.ActiveModal {
	case types.AddCommandForm:
		return deleteCharFromCommandForm(model)
	case types.AddSSEForm:
		return deleteCharFromSSEForm(model)
	case types.AddJSONForm:
		return deleteCharFromJSONForm(model)
	}
	return model
}

func deleteCharFromCommandForm(model types.Model) types.Model {
	switch model.FormData.ActiveField {
	case 0:
		model.FormData.Name = deleteLastChar(model.FormData.Name)
	case 1:
		model.FormData.Command = deleteLastChar(model.FormData.Command)
	case 2:
		model.FormData.Args = deleteLastChar(model.FormData.Args)
	case 3:
		model.FormData.Environment = deleteLastChar(model.FormData.Environment)
	}
	return model
}

func deleteCharFromSSEForm(model types.Model) types.Model {
	switch model.FormData.ActiveField {
	case 0:
		model.FormData.Name = deleteLastChar(model.FormData.Name)
	case 1:
		model.FormData.URL = deleteLastChar(model.FormData.URL)
	case 2:
		model.FormData.Environment = deleteLastChar(model.FormData.Environment)
	}
	return model
}

func deleteCharFromJSONForm(model types.Model) types.Model {
	switch model.FormData.ActiveField {
	case 0:
		model.FormData.Name = deleteLastChar(model.FormData.Name)
	case 1:
		model.FormData.JSONConfig = deleteLastChar(model.FormData.JSONConfig)
	case 2:
		model.FormData.Environment = deleteLastChar(model.FormData.Environment)
	}
	return model
}

func deleteLastChar(s string) string {
	if len(s) > 0 {
		return s[:len(s)-1]
	}
	return s
}

// Validation functions that return the updated model and validation status
func validateCommandForm(model types.Model) (types.Model, bool) {
	model.FormErrors = make(map[string]string)
	valid := true

	if model.FormData.Name == "" {
		model.FormErrors["name"] = "Name is required"
		valid = false
	}

	if model.FormData.Command == "" {
		model.FormErrors["command"] = "Command is required"
		valid = false
	}

	// Validate environment variables format if provided
	if model.FormData.Environment != "" {
		if err := validateEnvironmentFormat(model.FormData.Environment); err != nil {
			model.FormErrors["environment"] = err.Error()
			valid = false
		}
	}

	// Check for duplicate names (but allow the current MCP name in edit mode)
	for _, item := range model.MCPItems {
		if item.Name == model.FormData.Name {
			// Allow the current name in edit mode
			if !model.EditMode || item.Name != model.EditMCPName {
				model.FormErrors["name"] = "Name already exists"
				valid = false
				break
			}
		}
	}

	return model, valid
}

func validateSSEForm(model types.Model) (types.Model, bool) {
	model.FormErrors = make(map[string]string)
	valid := true

	if model.FormData.Name == "" {
		model.FormErrors["name"] = "Name is required"
		valid = false
	}

	if model.FormData.URL == "" {
		model.FormErrors["url"] = "URL is required"
		valid = false
	} else {
		// Validate URL format
		if _, err := url.Parse(model.FormData.URL); err != nil {
			model.FormErrors["url"] = "Invalid URL format"
			valid = false
		}
	}

	// Validate environment variables format if provided
	if model.FormData.Environment != "" {
		if err := validateEnvironmentFormat(model.FormData.Environment); err != nil {
			model.FormErrors["environment"] = err.Error()
			valid = false
		}
	}

	// Check for duplicate names (but allow the current MCP name in edit mode)
	for _, item := range model.MCPItems {
		if item.Name == model.FormData.Name {
			// Allow the current name in edit mode
			if !model.EditMode || item.Name != model.EditMCPName {
				model.FormErrors["name"] = "Name already exists"
				valid = false
				break
			}
		}
	}

	return model, valid
}

func validateJSONForm(model types.Model) (types.Model, bool) {
	model.FormErrors = make(map[string]string)
	valid := true

	if model.FormData.Name == "" {
		model.FormErrors["name"] = "Name is required"
		valid = false
	}

	if model.FormData.JSONConfig == "" {
		model.FormErrors["json"] = "JSON configuration is required"
		valid = false
	} else {
		// Validate JSON syntax with enhanced error reporting
		var js interface{}
		if err := json.Unmarshal([]byte(model.FormData.JSONConfig), &js); err != nil {
			// Extract line and column information from JSON error
			enhancedError := enhanceJSONError(err, model.FormData.JSONConfig)
			model.FormErrors["json"] = enhancedError
			valid = false
		}
	}

	// Check for duplicate names (but allow the current MCP name in edit mode)
	for _, item := range model.MCPItems {
		if item.Name == model.FormData.Name {
			// Allow the current name in edit mode
			if !model.EditMode || item.Name != model.EditMCPName {
				model.FormErrors["name"] = "Name already exists"
				valid = false
				break
			}
		}
	}

	return model, valid
}

// addMCPToInventory adds a new MCP to the inventory and saves it
func addMCPToInventory(model types.Model, mcpItem types.MCPItem) (types.Model, tea.Cmd) {
	// Add to model
	model.MCPItems = append(model.MCPItems, mcpItem)

	// Save to storage
	if err := services.SaveInventory(model.MCPItems); err != nil {
		// Show error message instead of success
		model.SuccessMessage = fmt.Sprintf("Failed to save %s: %v", mcpItem.Name, err)
		return model, hideSuccessMsg()
	}

	// Select the newly added item
	model.SelectedItem = len(model.MCPItems) - 1

	// Show success message
	model.SuccessMessage = fmt.Sprintf("Added %s successfully", mcpItem.Name)

	return model, hideSuccessMsg()
}

// updateMCPInInventory updates an existing MCP in the inventory and saves it
func updateMCPInInventory(model types.Model, updatedMCP types.MCPItem) (types.Model, tea.Cmd) {
	// Find and update the MCP in the inventory
	found := false
	for i, mcp := range model.MCPItems {
		if mcp.Name == model.EditMCPName {
			// Preserve the original active status
			updatedMCP.Active = model.MCPItems[i].Active
			model.MCPItems[i] = updatedMCP
			found = true
			break
		}
	}

	if !found {
		// MCP not found, show error
		model.SuccessMessage = fmt.Sprintf("Could not find MCP '%s' to update", model.EditMCPName)
		return model, hideSuccessMsg()
	}

	// Save to storage
	if err := services.SaveInventory(model.MCPItems); err != nil {
		// Show error message instead of success
		model.SuccessMessage = fmt.Sprintf("Failed to update %s: %v", updatedMCP.Name, err)
		return model, hideSuccessMsg()
	}

	// Update selection to the updated item if name changed
	if updatedMCP.Name != model.EditMCPName {
		// Find the updated MCP's new position
		for i, mcp := range model.MCPItems {
			if mcp.Name == updatedMCP.Name {
				model.SelectedItem = i
				break
			}
		}
	}

	// Show success message
	model.SuccessMessage = fmt.Sprintf("Updated %s successfully", updatedMCP.Name)

	return model, hideSuccessMsg()
}

// handleDeleteModalKeys handles keyboard input in the delete confirmation modal
func handleDeleteModalKeys(model types.Model, key string) (types.Model, tea.Cmd) {
	switch key {
	case "enter":
		// Confirm deletion
		var cmd tea.Cmd
		model, cmd = deleteMCPFromInventory(model)
		// Close modal and return to main navigation
		model.State = types.MainNavigation
		model.ActiveModal = types.NoModal
		return model, cmd
	case "esc":
		// Cancel deletion
		model.State = types.MainNavigation
		model.ActiveModal = types.NoModal
		return model, nil
	}
	return model, nil
}

// deleteMCPFromInventory removes the selected MCP from the inventory and saves it
func deleteMCPFromInventory(model types.Model) (types.Model, tea.Cmd) {
	// Get the selected MCP
	filteredMCPs := services.GetFilteredMCPs(model)
	selectedIndex := model.SelectedItem
	if model.SearchQuery != "" {
		selectedIndex = model.FilteredSelectedIndex
	}

	if selectedIndex >= len(filteredMCPs) {
		model.SuccessMessage = "No MCP selected to delete"
		return model, hideSuccessMsg()
	}

	mcpToDelete := filteredMCPs[selectedIndex]

	// Find the MCP in the main inventory and remove it
	for i, mcp := range model.MCPItems {
		if mcp.Name == mcpToDelete.Name {
			// Remove the MCP from the slice
			model.MCPItems = append(model.MCPItems[:i], model.MCPItems[i+1:]...)
			break
		}
	}

	// Adjust selection if needed
	if model.SelectedItem >= len(model.MCPItems) && len(model.MCPItems) > 0 {
		model.SelectedItem = len(model.MCPItems) - 1
	} else if len(model.MCPItems) == 0 {
		model.SelectedItem = 0
	}

	// Adjust filtered selection if needed
	if model.SearchQuery != "" {
		newFilteredMCPs := services.GetFilteredMCPs(model)
		if model.FilteredSelectedIndex >= len(newFilteredMCPs) && len(newFilteredMCPs) > 0 {
			model.FilteredSelectedIndex = len(newFilteredMCPs) - 1
		} else if len(newFilteredMCPs) == 0 {
			model.FilteredSelectedIndex = 0
		}
	}

	// Save to storage
	if err := services.SaveInventory(model.MCPItems); err != nil {
		model.SuccessMessage = fmt.Sprintf("Failed to delete %s: %v", mcpToDelete.Name, err)
		return model, hideSuccessMsg()
	}

	// Show success message
	model.SuccessMessage = fmt.Sprintf("Deleted %s successfully", mcpToDelete.Name)

	return model, hideSuccessMsg()
}

// parseArgsString converts a string of arguments to []string
// Supports both space-separated ("arg1 arg2 arg3") and quoted arguments
func parseArgsString(argsStr string) []string {
	if argsStr == "" {
		return nil
	}

	var args []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(argsStr); i++ {
		char := argsStr[i]

		switch char {
		case '"', '\'':
			if !inQuotes {
				// Start of quoted string
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				// End of quoted string
				inQuotes = false
				quoteChar = 0
			} else {
				// Quote inside different quote type
				current.WriteByte(char)
			}
		case ' ', '\t':
			if inQuotes {
				// Space inside quotes - add to current arg
				current.WriteByte(char)
			} else if current.Len() > 0 {
				// Space outside quotes - end current arg
				args = append(args, current.String())
				current.Reset()
			}
			// Skip spaces outside quotes
		default:
			current.WriteByte(char)
		}
	}

	// Add final argument if exists
	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

// parseEnvironmentString converts a string of environment variables to map[string]string
// Supports formats: "KEY1=value1,KEY2=value2" or "KEY1=value1\nKEY2=value2"
func parseEnvironmentString(envStr string) map[string]string {
	if envStr == "" {
		return nil
	}

	env := make(map[string]string)

	// Support both comma-separated and newline-separated
	lines := strings.FieldsFunc(envStr, func(c rune) bool {
		return c == ',' || c == '\n'
	})

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Split on first = sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" {
				env[key] = value
			}
		}
	}

	if len(env) == 0 {
		return nil
	}

	return env
}

// copyActiveFieldToClipboard copies the content of the active field to clipboard
func copyActiveFieldToClipboard(model types.Model) types.Model {
	content := getActiveFieldContent(model)
	if content == "" {
		return model
	}

	clipboardService := services.NewClipboardService()
	if err := clipboardService.Copy(content); err != nil {
		model.SuccessMessage = "Failed to copy to clipboard: " + err.Error()
		model.SuccessTimer = 180 // Show error message for 3 seconds (60 ticks per second)
	} else {
		model.SuccessMessage = "Copied to clipboard"
		model.SuccessTimer = 120 // Show success message for 2 seconds
	}

	return model
}

// getActiveFieldContent extracts content from the currently active field
func getActiveFieldContent(model types.Model) string {
	switch model.ActiveModal {
	case types.AddCommandForm:
		return getCommandFormFieldContent(model)
	case types.AddSSEForm:
		return getSSEFormFieldContent(model)
	case types.AddJSONForm:
		return getJSONFormFieldContent(model)
	default:
		return ""
	}
}

func getCommandFormFieldContent(model types.Model) string {
	switch model.FormData.ActiveField {
	case 0:
		return model.FormData.Name
	case 1:
		return model.FormData.Command
	case 2:
		return model.FormData.Args
	case 3:
		return model.FormData.Environment
	default:
		return ""
	}
}

func getSSEFormFieldContent(model types.Model) string {
	switch model.FormData.ActiveField {
	case 0:
		return model.FormData.Name
	case 1:
		return model.FormData.URL
	case 2:
		return model.FormData.Environment
	default:
		return ""
	}
}

func getJSONFormFieldContent(model types.Model) string {
	switch model.FormData.ActiveField {
	case 0:
		return model.FormData.Name
	case 1:
		return model.FormData.JSONConfig
	case 2:
		return model.FormData.Environment
	default:
		return ""
	}
}

// pasteFromClipboardToActiveField pastes clipboard content to the active field
func pasteFromClipboardToActiveField(model types.Model) types.Model {
	clipboardService := services.NewClipboardService()

	// Use enhanced paste for better error diagnostics
	content, err := clipboardService.EnhancedPaste()
	if err != nil {
		// Add user feedback for clipboard paste failure with enhanced error information
		model.SuccessMessage = "Failed to paste from clipboard: " + err.Error()
		model.SuccessTimer = 240 // Show error message for 4 seconds to allow reading detailed error
		return model
	}

	switch model.ActiveModal {
	case types.AddCommandForm:
		switch model.FormData.ActiveField {
		case 0:
			model.FormData.Name = content
		case 1:
			model.FormData.Command = content
		case 2:
			model.FormData.Args = content
		case 3:
			model.FormData.Environment = content
		}
	case types.AddSSEForm:
		switch model.FormData.ActiveField {
		case 0:
			model.FormData.Name = content
		case 1:
			model.FormData.URL = content
		case 2:
			model.FormData.Environment = content
		}
	case types.AddJSONForm:
		switch model.FormData.ActiveField {
		case 0:
			model.FormData.Name = content
		case 1:
			model.FormData.JSONConfig = content
		case 2:
			model.FormData.Environment = content
		}
	}

	// Add success feedback for successful paste operation
	model.SuccessMessage = "Pasted from clipboard"
	model.SuccessTimer = 120 // Show success message for 2 seconds

	return model
}

// enhanceJSONError provides detailed JSON error information with line/column details
func enhanceJSONError(err error, jsonContent string) string {
	errStr := err.Error()

	// Check if it's a JSON syntax error with offset information
	if syntaxErr, ok := err.(*json.SyntaxError); ok {
		line, col := getLineColumn(jsonContent, syntaxErr.Offset)
		return fmt.Sprintf("JSON syntax error at line %d, column %d: %s", line, col, errStr)
	}

	// Check if it's a JSON unmarshaling type error
	if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
		line, col := getLineColumn(jsonContent, typeErr.Offset)
		return fmt.Sprintf("JSON type error at line %d, column %d: expected %s but got %s",
			line, col, typeErr.Type.String(), typeErr.Value)
	}

	// For other errors, provide general guidance
	if strings.Contains(errStr, "unexpected end of JSON input") {
		return "Incomplete JSON: missing closing bracket or brace"
	}

	if strings.Contains(errStr, "invalid character") {
		return errStr + " - check for unescaped quotes or special characters"
	}

	return "Invalid JSON: " + errStr
}

// getLineColumn calculates line and column numbers from byte offset
func getLineColumn(content string, offset int64) (line int, col int) {
	line = 1
	col = 1

	for i, r := range content {
		if int64(i) >= offset {
			break
		}
		if r == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}

	return line, col
}

// focusOnFirstErrorField moves focus to the first field that has a validation error
func focusOnFirstErrorField(model types.Model) types.Model {
	switch model.ActiveModal {
	case types.AddCommandForm:
		// Check field order: Name (0), Command (1), Args (2), Environment (3)
		if _, exists := model.FormErrors["name"]; exists {
			model.FormData.ActiveField = 0
		} else if _, exists := model.FormErrors["command"]; exists {
			model.FormData.ActiveField = 1
		}
		// Args and Environment don't have validation errors in current implementation
	case types.AddSSEForm:
		// Check field order: Name (0), URL (1), Environment (2)
		if _, exists := model.FormErrors["name"]; exists {
			model.FormData.ActiveField = 0
		} else if _, exists := model.FormErrors["url"]; exists {
			model.FormData.ActiveField = 1
		}
	case types.AddJSONForm:
		// Check field order: Name (0), JSONConfig (1), Environment (2)
		if _, exists := model.FormErrors["name"]; exists {
			model.FormData.ActiveField = 0
		} else if _, exists := model.FormErrors["json"]; exists {
			model.FormData.ActiveField = 1
		}
	}

	return model
}

// validateEnvironmentFormat validates environment variable format
func validateEnvironmentFormat(envStr string) error {
	if envStr == "" {
		return nil
	}

	lines := parseEnvironmentLines(envStr)
	for _, line := range lines {
		if err := validateEnvironmentLine(line); err != nil {
			return err
		}
	}

	return nil
}

func parseEnvironmentLines(envStr string) []string {
	lines := strings.FieldsFunc(envStr, func(c rune) bool {
		return c == ',' || c == '\n'
	})

	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}
	return cleanLines
}

func validateEnvironmentLine(line string) error {
	if !strings.Contains(line, "=") {
		return fmt.Errorf("invalid format: '%s' - use KEY=value format", line)
	}

	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format: '%s' - use KEY=value format", line)
	}

	key := strings.TrimSpace(parts[0])
	if key == "" {
		return fmt.Errorf("empty key in: '%s'", line)
	}

	return validateEnvironmentKey(key)
}

func validateEnvironmentKey(key string) error {
	for _, r := range key {
		if !isValidKeyChar(r) {
			return fmt.Errorf("invalid key '%s' - use letters, numbers, and underscores only", key)
		}
	}
	return nil
}

func isValidKeyChar(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_'
}
