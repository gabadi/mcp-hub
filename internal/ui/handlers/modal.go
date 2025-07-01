package handlers

import (
	"cc-mcp-manager/internal/ui/services"
	"cc-mcp-manager/internal/ui/types"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
		model.FormData.ActiveField = (model.FormData.ActiveField + 1) % 3 // 3 fields: Name, Command, Args
	case "enter":
		// Submit form if valid
		var valid bool
		model, valid = validateCommandForm(model)
		if valid {
			mcpItem := types.MCPItem{
				Name:    model.FormData.Name,
				Type:    "CMD",
				Active:  false,
				Command: model.FormData.Command,
				Args:    model.FormData.Args,
			}
			var cmd tea.Cmd
			model, cmd = addMCPToInventory(model, mcpItem)
			// Close modal and return to main navigation
			model.State = types.MainNavigation
			model.ActiveModal = types.NoModal
			model.FormData = types.FormData{}
			model.FormErrors = make(map[string]string)
			return model, cmd
		}
	case "backspace":
		// Delete character from active field
		model = deleteCharFromActiveField(model)
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
		model.FormData.ActiveField = (model.FormData.ActiveField + 1) % 2 // 2 fields: Name, URL
	case "enter":
		// Submit form if valid
		var valid bool
		model, valid = validateSSEForm(model)
		if valid {
			mcpItem := types.MCPItem{
				Name:   model.FormData.Name,
				Type:   "SSE",
				Active: false,
				URL:    model.FormData.URL,
			}
			var cmd tea.Cmd
			model, cmd = addMCPToInventory(model, mcpItem)
			// Close modal and return to main navigation
			model.State = types.MainNavigation
			model.ActiveModal = types.NoModal
			model.FormData = types.FormData{}
			model.FormErrors = make(map[string]string)
			return model, cmd
		}
	case "backspace":
		// Delete character from active field
		model = deleteCharFromActiveField(model)
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
		model.FormData.ActiveField = (model.FormData.ActiveField + 1) % 2 // 2 fields: Name, JSONConfig
	case "enter":
		// Submit form if valid (or newline in JSON field)
		if model.FormData.ActiveField == 1 { // JSON field
			model = addCharToActiveField(model, "\n")
		} else {
			var valid bool
			model, valid = validateJSONForm(model)
			if valid {
				mcpItem := types.MCPItem{
					Name:       model.FormData.Name,
					Type:       "JSON",
					Active:     false,
					JSONConfig: model.FormData.JSONConfig,
				}
				var cmd tea.Cmd
				model, cmd = addMCPToInventory(model, mcpItem)
				// Close modal and return to main navigation
				model.State = types.MainNavigation
				model.ActiveModal = types.NoModal
				model.FormData = types.FormData{}
				model.FormErrors = make(map[string]string)
				return model, cmd
			}
		}
	case "backspace":
		// Delete character from active field
		model = deleteCharFromActiveField(model)
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
		}
	case types.AddSSEForm:
		switch model.FormData.ActiveField {
		case 0:
			model.FormData.Name += char
		case 1:
			model.FormData.URL += char
		}
	case types.AddJSONForm:
		switch model.FormData.ActiveField {
		case 0:
			model.FormData.Name += char
		case 1:
			model.FormData.JSONConfig += char
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
	}
	return model
}

func deleteCharFromSSEForm(model types.Model) types.Model {
	switch model.FormData.ActiveField {
	case 0:
		model.FormData.Name = deleteLastChar(model.FormData.Name)
	case 1:
		model.FormData.URL = deleteLastChar(model.FormData.URL)
	}
	return model
}

func deleteCharFromJSONForm(model types.Model) types.Model {
	switch model.FormData.ActiveField {
	case 0:
		model.FormData.Name = deleteLastChar(model.FormData.Name)
	case 1:
		model.FormData.JSONConfig = deleteLastChar(model.FormData.JSONConfig)
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

	// Check for duplicate names
	for _, item := range model.MCPItems {
		if item.Name == model.FormData.Name {
			model.FormErrors["name"] = "Name already exists"
			valid = false
			break
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

	// Check for duplicate names
	for _, item := range model.MCPItems {
		if item.Name == model.FormData.Name {
			model.FormErrors["name"] = "Name already exists"
			valid = false
			break
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
		// Validate JSON syntax
		var js interface{}
		if err := json.Unmarshal([]byte(model.FormData.JSONConfig), &js); err != nil {
			model.FormErrors["json"] = "Invalid JSON: " + err.Error()
			valid = false
		}
	}

	// Check for duplicate names
	for _, item := range model.MCPItems {
		if item.Name == model.FormData.Name {
			model.FormErrors["name"] = "Name already exists"
			valid = false
			break
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
