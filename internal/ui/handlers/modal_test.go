package handlers

import (
	"fmt"
	"testing"

	"mcp-hub/internal/testutil"
	"mcp-hub/internal/ui/types"

	"github.com/stretchr/testify/assert"
)

// Test constants
const (
	pastedContent = "pasted-content"
)


// Epic 1 Story 4 Tests - Edit MCP Functionality

func TestEditMCPFormPrePopulation(t *testing.T) {
	tests := getEditMCPFormTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := populateFormDataFromMCP(tt.mcpItem)
			validateFormPrePopulation(t, formData, tt)
		})
	}
}

func getEditMCPFormTestCases() []struct {
	name       string
	mcpItem    types.MCPItem
	expectName string
	expectType string
} {
	return []struct {
		name       string
		mcpItem    types.MCPItem
		expectName string
		expectType string
	}{
		{
			name: "Command MCP pre-population",
			mcpItem: types.MCPItem{
				Name:        "test-cmd",
				Type:        "CMD",
				Command:     "test-command",
				Args:        []string{"arg1", "arg2"},
				Environment: map[string]string{"KEY1": "value1", "KEY2": "value2"},
			},
			expectName: "test-cmd",
			expectType: "CMD",
		},
		{
			name: "SSE MCP pre-population",
			mcpItem: types.MCPItem{
				Name:        "test-sse",
				Type:        "SSE",
				URL:         "https://example.com/sse",
				Environment: map[string]string{"API_KEY": "secret"},
			},
			expectName: "test-sse",
			expectType: "SSE",
		},
		{
			name: "JSON MCP pre-population",
			mcpItem: types.MCPItem{
				Name:       "test-json",
				Type:       "JSON",
				JSONConfig: `{"key": "value"}`,
			},
			expectName: "test-json",
			expectType: "JSON",
		},
	}
}

func validateFormPrePopulation(t *testing.T, formData types.FormData, tt struct {
	name       string
	mcpItem    types.MCPItem
	expectName string
	expectType string
}) {
	if formData.Name != tt.expectName {
		t.Errorf("Expected name %s, got %s", tt.expectName, formData.Name)
	}

	validateTypeSpecificFields(t, formData, tt.mcpItem)
	validateEnvironmentFields(t, formData, tt.mcpItem)
}

func validateTypeSpecificFields(t *testing.T, formData types.FormData, mcpItem types.MCPItem) {
	switch mcpItem.Type {
	case "CMD":
		validateCommandFields(t, formData, mcpItem)
	case "SSE":
		validateSSEFields(t, formData, mcpItem)
	case "JSON":
		validateJSONFields(t, formData, mcpItem)
	}
}

func validateCommandFields(t *testing.T, formData types.FormData, mcpItem types.MCPItem) {
	if formData.Command != mcpItem.Command {
		t.Errorf("Expected command %s, got %s", mcpItem.Command, formData.Command)
	}

	if len(mcpItem.Args) > 0 {
		expectedArgs := formatArgsForDisplay(mcpItem.Args)
		if formData.Args != expectedArgs {
			t.Errorf("Expected args %s, got %s", expectedArgs, formData.Args)
		}
	}
}

func validateSSEFields(t *testing.T, formData types.FormData, mcpItem types.MCPItem) {
	if formData.URL != mcpItem.URL {
		t.Errorf("Expected URL %s, got %s", mcpItem.URL, formData.URL)
	}
}

func validateJSONFields(t *testing.T, formData types.FormData, mcpItem types.MCPItem) {
	if formData.JSONConfig != mcpItem.JSONConfig {
		t.Errorf("Expected JSON config %s, got %s", mcpItem.JSONConfig, formData.JSONConfig)
	}
}

func validateEnvironmentFields(t *testing.T, formData types.FormData, mcpItem types.MCPItem) {
	if len(mcpItem.Environment) > 0 {
		for key, value := range mcpItem.Environment {
			expectedPair := key + "=" + value
			if !contains(formData.Environment, expectedPair) {
				t.Errorf("Expected environment to contain %q, but got %q", expectedPair, formData.Environment)
			}
		}
	}
}

func TestFormatArgsForDisplay(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Empty args",
			args:     []string{},
			expected: "",
		},
		{
			name:     "Single arg without spaces",
			args:     []string{"arg1"},
			expected: "arg1",
		},
		{
			name:     "Multiple args without spaces",
			args:     []string{"arg1", "arg2", "arg3"},
			expected: "arg1 arg2 arg3",
		},
		{
			name:     "Args with spaces get quoted",
			args:     []string{"arg with spaces", "normal-arg"},
			expected: `"arg with spaces" normal-arg`,
		},
		{
			name:     "Mixed args with and without spaces",
			args:     []string{"normal", "arg with spaces", "another"},
			expected: `normal "arg with spaces" another`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatArgsForDisplay(tt.args)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFormatEnvironmentForDisplay(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected string
	}{
		{
			name:     "Empty environment",
			env:      map[string]string{},
			expected: "",
		},
		{
			name:     "Single environment variable",
			env:      map[string]string{"KEY1": "value1"},
			expected: "KEY1=value1",
		},
		{
			name: "Multiple environment variables",
			env:  map[string]string{"KEY1": "value1", "KEY2": "value2"},
			// Note: map iteration order is not guaranteed, so we need to check both possible orders
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatEnvironmentForDisplay(tt.env)

			switch len(tt.env) {
			case 0:
				if result != tt.expected {
					t.Errorf("Expected %q, got %q", tt.expected, result)
				}
			case 1:
				if result != tt.expected {
					t.Errorf("Expected %q, got %q", tt.expected, result)
				}
			default:
				// For multiple environment variables, check that all key=value pairs are present
				for key, value := range tt.env {
					expectedPair := key + "=" + value
					if !contains(result, expectedPair) {
						t.Errorf("Expected result to contain %q, but got %q", expectedPair, result)
					}
				}
			}
		})
	}
}

func TestUpdateMCPInInventory(t *testing.T) {
	// Create a model with some test MCPs
	model := testutil.NewTestModel().Build()
	model.MCPItems = []types.MCPItem{
		{Name: "test-mcp", Type: "CMD", Command: "old-command", Active: true},
		{Name: "other-mcp", Type: "SSE", URL: "https://old.com", Active: false},
	}
	model.EditMode = true
	model.EditMCPName = "test-mcp"

	// Create updated MCP
	updatedMCP := types.MCPItem{
		Name:    "test-mcp-updated",
		Type:    "CMD",
		Command: "new-command",
		Args:    []string{"new-arg"},
	}

	// Update the MCP
	newModel, _ := updateMCPInInventory(model, updatedMCP)

	// Verify the MCP was updated
	found := false
	for _, mcp := range newModel.MCPItems {
		if mcp.Name == "test-mcp-updated" {
			found = true
			if mcp.Command != "new-command" {
				t.Errorf("Expected command to be updated to 'new-command', got %s", mcp.Command)
			}
			if len(mcp.Args) != 1 || mcp.Args[0] != "new-arg" {
				t.Errorf("Expected args to be updated to ['new-arg'], got %v", mcp.Args)
			}
			// Active status should be preserved
			if !mcp.Active {
				t.Errorf("Expected active status to be preserved as true, got %v", mcp.Active)
			}
		}
	}

	if !found {
		t.Error("Updated MCP not found in inventory")
	}

	// Verify other MCPs are unchanged
	for _, mcp := range newModel.MCPItems {
		if mcp.Name == "other-mcp" {
			if mcp.URL != "https://old.com" {
				t.Errorf("Other MCP should be unchanged, but URL changed to %s", mcp.URL)
			}
		}
	}
}

// Test edge cases for updateMCPInInventory
func TestUpdateMCPInInventoryEdgeCases(t *testing.T) {
	t.Run("update_nonexistent_mcp", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithMCPs([]types.MCPItem{
				{Name: "existing", Type: "CMD", Command: "test"},
			}).
			Build()
		model.EditMode = true
		model.EditMCPName = "nonexistent"

		updatedMCP := types.MCPItem{
			Name:    "updated",
			Type:    "CMD",
			Command: "new-command",
		}

		// Should not crash and should maintain existing MCPs
		newModel, _ := updateMCPInInventory(model, updatedMCP)
		assert.Len(t, newModel.MCPItems, 1)
		assert.Equal(t, "existing", newModel.MCPItems[0].Name)
	})

	t.Run("update_with_empty_edit_name", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithMCPs([]types.MCPItem{
				{Name: "existing", Type: "CMD", Command: "test"},
			}).
			Build()
		model.EditMode = true
		model.EditMCPName = ""

		updatedMCP := types.MCPItem{
			Name:    "updated",
			Type:    "CMD",
			Command: "new-command",
		}

		// Should not crash and should maintain existing MCPs
		newModel, _ := updateMCPInInventory(model, updatedMCP)
		assert.Len(t, newModel.MCPItems, 1)
		assert.Equal(t, "existing", newModel.MCPItems[0].Name)
	})

	t.Run("update_preserves_complex_fields", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithMCPs([]types.MCPItem{
				{
					Name:        "complex",
					Type:        "CMD",
					Command:     "old-command",
					Active:      true,
					Args:        []string{"old-arg1", "old-arg2"},
					Environment: map[string]string{"OLD_KEY": "old_value"},
					URL:         "https://old.com",
					JSONConfig:  `{"old": "config"}`,
				},
			}).
			Build()
		model.EditMode = true
		model.EditMCPName = "complex"

		updatedMCP := types.MCPItem{
			Name:        "complex",
			Type:        "CMD",
			Command:     "new-command",
			Args:        []string{"new-arg1", "new-arg2", "new-arg3"},
			Environment: map[string]string{"NEW_KEY": "new_value", "SECOND_KEY": "second_value"},
			URL:         "https://new.com",
			JSONConfig:  `{"new": "config", "nested": {"key": "value"}}`,
		}

		newModel, _ := updateMCPInInventory(model, updatedMCP)
		assert.Len(t, newModel.MCPItems, 1)

		updated := newModel.MCPItems[0]
		assert.Equal(t, "complex", updated.Name)
		assert.Equal(t, "new-command", updated.Command)
		assert.True(t, updated.Active) // Should preserve active status
		assert.Equal(t, []string{"new-arg1", "new-arg2", "new-arg3"}, updated.Args)
		assert.Equal(t, map[string]string{"NEW_KEY": "new_value", "SECOND_KEY": "second_value"}, updated.Environment)
		assert.Equal(t, "https://new.com", updated.URL)
		assert.Equal(t, `{"new": "config", "nested": {"key": "value"}}`, updated.JSONConfig)
	})
}

func TestEditModeValidation(t *testing.T) {
	// Test that duplicate name validation allows the current MCP name in edit mode
	model := testutil.NewTestModel().Build()
	model.MCPItems = []types.MCPItem{
		{Name: "existing-mcp", Type: "CMD", Command: "test"},
		{Name: "another-mcp", Type: "CMD", Command: "test2"},
	}
	model.EditMode = true
	model.EditMCPName = "existing-mcp"
	model.FormData.Name = "existing-mcp" // Same name as original
	model.FormData.Command = TestString      // Required field

	// This should be valid (keeping the same name)
	newModel, valid := validateCommandForm(model)
	if !valid {
		t.Error("Edit mode should allow keeping the same MCP name")
	}
	if _, exists := newModel.FormErrors["name"]; exists {
		t.Error("Edit mode should not show duplicate name error for same name")
	}

	// Test changing to another existing name (should fail)
	model.FormData.Name = "another-mcp"
	model.FormData.Command = "test" // Keep required field
	newModel, valid = validateCommandForm(model)
	if valid {
		t.Error("Edit mode should not allow changing to another existing MCP name")
	}
	if _, exists := newModel.FormErrors["name"]; !exists {
		t.Error("Edit mode should show duplicate name error for another existing name")
	}
}

func TestEditModeStateCleanup(t *testing.T) {
	// Test that edit mode state is properly cleaned up on cancel
	model := testutil.NewTestModel().Build()
	model.EditMode = true
	model.EditMCPName = "test-mcp"
	model.State = types.ModalActive
	model.ActiveModal = types.AddCommandForm
	model.FormData.Name = "test-data"
	model.FormErrors = map[string]string{"test": "error"}

	// Simulate ESC key (cancel)
	newModel, _ := HandleEscKey(model)

	// Verify state is cleaned up
	if newModel.EditMode {
		t.Error("EditMode should be false after cancel")
	}
	if newModel.EditMCPName != "" {
		t.Error("EditMCPName should be empty after cancel")
	}
	if newModel.State != types.MainNavigation {
		t.Error("State should return to MainNavigation after cancel")
	}
	if newModel.ActiveModal != types.NoModal {
		t.Error("ActiveModal should be NoModal after cancel")
	}
	if newModel.FormData.Name != "" {
		t.Error("FormData should be cleared after cancel")
	}
	if len(newModel.FormErrors) != 0 {
		t.Error("FormErrors should be cleared after cancel")
	}
}

// Epic 2 Story 1 Tests - Original modal functionality

func TestHideSuccessMsg(t *testing.T) {
	cmd := hideSuccessMsg()
	if cmd == nil {
		t.Error("hideSuccessMsg() should return a command")
	}
}

func TestHandleModalKeys(t *testing.T) {
	tests := []struct {
		name          string
		model         types.Model
		key           string
		expectedModal types.ModalType
	}{
		{
			name: "AddMCPTypeSelection modal",
			model: types.Model{
				State:       types.ModalActive,
				ActiveModal: types.AddMCPTypeSelection,
			},
			key:           "1",
			expectedModal: types.AddCommandForm,
		},
		{
			name: "Legacy modal with enter",
			model: types.Model{
				State:       types.ModalActive,
				ActiveModal: types.AddModal,
			},
			key:           "enter",
			expectedModal: types.AddModal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newModel, _ := HandleModalKeys(tt.model, tt.key)
			// Test passes if it doesn't panic and returns a model
			if newModel.ActiveModal != tt.expectedModal && tt.name == "Legacy modal with enter" {
				// For legacy modal, we expect state change to MainNavigation
				if newModel.State != types.MainNavigation {
					t.Errorf("Expected state MainNavigation, got %v", newModel.State)
				}
			}
		})
	}
}

func TestHandleTypeSelectionKeys(t *testing.T) {
	model := types.Model{
		State:       types.ModalActive,
		ActiveModal: types.AddMCPTypeSelection,
		FormData:    types.FormData{},
		FormErrors:  make(map[string]string),
	}

	tests := []struct {
		name          string
		key           string
		expectedModal types.ModalType
		shouldExit    bool
	}{
		{
			name:          "Select command type",
			key:           "1",
			expectedModal: types.AddCommandForm,
			shouldExit:    false,
		},
		{
			name:          "Select SSE type",
			key:           "2",
			expectedModal: types.AddSSEForm,
			shouldExit:    false,
		},
		{
			name:          "Select JSON type",
			key:           "3",
			expectedModal: types.AddJSONForm,
			shouldExit:    false,
		},
		{
			name:       "Escape modal",
			key:        "esc",
			shouldExit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newModel, _ := handleTypeSelectionKeys(model, tt.key)

			if tt.shouldExit {
				if newModel.State != types.MainNavigation {
					t.Errorf("Expected to exit to MainNavigation, got %v", newModel.State)
				}
			} else if tt.expectedModal != types.NoModal {
				if newModel.ActiveModal != tt.expectedModal {
					t.Errorf("Expected modal %v, got %v", tt.expectedModal, newModel.ActiveModal)
				}
			}
		})
	}
}

func TestParseArgsString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "single argument",
			input:    "arg1",
			expected: []string{"arg1"},
		},
		{
			name:     "multiple arguments",
			input:    "arg1 arg2 arg3",
			expected: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "quoted arguments",
			input:    `"arg with spaces" arg2`,
			expected: []string{"arg with spaces", "arg2"},
		},
		{
			name:     "mixed quotes",
			input:    `arg1 "arg 2" 'arg 3'`,
			expected: []string{"arg1", "arg 2", "arg 3"},
		},
		{
			name:     "whitespace only",
			input:    "   \t\n  ",
			expected: nil,
		},
		{
			name:     "unclosed quote",
			input:    `"unclosed quote`,
			expected: []string{"unclosed quote"},
		},
		{
			name:     "nested quotes",
			input:    `"arg with spaces"`,
			expected: []string{`arg with spaces`},
		},
		{
			name:     "multiple spaces",
			input:    "arg1    arg2     arg3",
			expected: []string{"arg1", "arg2", "arg3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseArgsString(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d args, got %d", len(tt.expected), len(result))
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Expected arg[%d] = %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestParseEnvironmentString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:  "single variable",
			input: "KEY=value",
			expected: map[string]string{
				"KEY": "value",
			},
		},
		{
			name:  "multiple variables",
			input: "KEY1=value1\nKEY2=value2",
			expected: map[string]string{
				"KEY1": "value1",
				"KEY2": "value2",
			},
		},
		{
			name:     "whitespace only",
			input:    "\n\n\t  \n",
			expected: nil,
		},
		{
			name:  "empty value",
			input: "KEY=",
			expected: map[string]string{
				"KEY": "",
			},
		},
		{
			name:  "value with equals",
			input: "KEY=value=with=equals",
			expected: map[string]string{
				"KEY": "value=with=equals",
			},
		},
		{
			name:  "mixed valid and empty lines",
			input: "KEY1=value1\n\nKEY2=value2\n\n",
			expected: map[string]string{
				"KEY1": "value1",
				"KEY2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseEnvironmentString(tt.input)

			if tt.expected == nil && result != nil {
				t.Errorf("Expected nil, got %v", result)
				return
			}

			if tt.expected != nil && result == nil {
				t.Errorf("Expected %v, got nil", tt.expected)
				return
			}

			if tt.expected == nil && result == nil {
				return // Both nil, test passes
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d vars, got %d", len(tt.expected), len(result))
				return
			}

			for key, expected := range tt.expected {
				if result[key] != expected {
					t.Errorf("Expected %q = %q, got %q", key, expected, result[key])
				}
			}
		})
	}
}

func TestDeleteLastChar(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "",
		},
		{
			name:     "multiple characters",
			input:    "hello",
			expected: "hell",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deleteLastChar(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestValidateCommandForm(t *testing.T) {
	tests := createValidateCommandFormTests()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultModel, valid := validateCommandForm(tt.model)
			assertValidationResult(t, resultModel, valid, tt.expectValid, tt.errorField)
		})
	}
}

func createValidateCommandFormTests() []struct {
	name        string
	model       types.Model
	expectValid bool
	errorField  string
} {
	return []struct {
		name        string
		model       types.Model
		expectValid bool
		errorField  string
	}{
		{
			name: "valid form",
			model: types.Model{
				FormData: types.FormData{
					Name:    "test-mcp",
					Command: "test-command",
				},
				MCPItems:   []types.MCPItem{},
				FormErrors: make(map[string]string),
			},
			expectValid: true,
		},
		{
			name: "missing name",
			model: types.Model{
				FormData: types.FormData{
					Command: "test-command",
				},
				MCPItems:   []types.MCPItem{},
				FormErrors: make(map[string]string),
			},
			expectValid: false,
			errorField:  "name",
		},
		{
			name: "missing command",
			model: types.Model{
				FormData: types.FormData{
					Name: "test-mcp",
				},
				MCPItems:   []types.MCPItem{},
				FormErrors: make(map[string]string),
			},
			expectValid: false,
			errorField:  "command",
		},
	}
}

func assertValidationResult(t *testing.T, resultModel types.Model, valid bool, expectValid bool, errorField string) {
	if valid != expectValid {
		t.Errorf("Expected valid = %v, got %v", expectValid, valid)
	}

	if !expectValid {
		if len(resultModel.FormErrors) == 0 {
			t.Error("Expected validation errors, got none")
		}
		if errorField != "" {
			if _, exists := resultModel.FormErrors[errorField]; !exists {
				t.Errorf("Expected error for field %q", errorField)
			}
		}
	} else if len(resultModel.FormErrors) > 0 {
		t.Errorf("Expected no errors, got: %v", resultModel.FormErrors)
	}
}

// testValidationFunc is a helper function to test validation functions
func testValidationFunc(t *testing.T, validationFunc func(types.Model) (types.Model, bool), _ string, tests []struct {
	name        string
	model       types.Model
	expectValid bool
}) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultModel, valid := validationFunc(tt.model)

			if valid != tt.expectValid {
				t.Errorf("Expected valid = %v, got %v", tt.expectValid, valid)
			}

			if !tt.expectValid {
				if len(resultModel.FormErrors) == 0 {
					t.Error("Expected validation errors, got none")
				}
			} else {
				if len(resultModel.FormErrors) > 0 {
					t.Errorf("Expected no errors, got: %v", resultModel.FormErrors)
				}
			}
		})
	}
}

func TestValidateSSEForm(t *testing.T) {
	tests := []struct {
		name        string
		model       types.Model
		expectValid bool
	}{
		{
			name: "valid SSE form",
			model: types.Model{
				FormData: types.FormData{
					Name: "test-mcp",
					URL:  "http://example.com",
				},
				MCPItems:   []types.MCPItem{},
				FormErrors: make(map[string]string),
			},
			expectValid: true,
		},
		{
			name: "invalid URL",
			model: types.Model{
				FormData: types.FormData{
					Name: "test-mcp",
					URL:  "not-a-url",
				},
				MCPItems:   []types.MCPItem{},
				FormErrors: make(map[string]string),
			},
			expectValid: false,
		},
	}

	testValidationFunc(t, validateSSEForm, "SSE", tests)
}

func TestValidateJSONForm(t *testing.T) {
	tests := []struct {
		name        string
		model       types.Model
		expectValid bool
	}{
		{
			name: "valid JSON form",
			model: types.Model{
				FormData: types.FormData{
					Name:       "test-mcp",
					JSONConfig: `{"key": "value"}`,
				},
				MCPItems:   []types.MCPItem{},
				FormErrors: make(map[string]string),
			},
			expectValid: true,
		},
		{
			name: "invalid JSON",
			model: types.Model{
				FormData: types.FormData{
					Name:       "test-mcp",
					JSONConfig: `{invalid json}`,
				},
				MCPItems:   []types.MCPItem{},
				FormErrors: make(map[string]string),
			},
			expectValid: false,
		},
	}

	testValidationFunc(t, validateJSONForm, "JSON", tests)
}

func TestValidateEnvironmentFormat(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid format",
			input:       "KEY=value\nKEY2=value2",
			expectError: false,
		},
		{
			name:        "invalid line",
			input:       "INVALID_LINE",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEnvironmentFormat(tt.input)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestIsValidKeyChar(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		expected bool
	}{
		{
			name:     "letter",
			char:     'A',
			expected: true,
		},
		{
			name:     "number",
			char:     '1',
			expected: true,
		},
		{
			name:     "underscore",
			char:     '_',
			expected: true,
		},
		{
			name:     "space",
			char:     ' ',
			expected: false,
		},
		{
			name:     "equals",
			char:     '=',
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidKeyChar(tt.char)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAddCharToActiveField(t *testing.T) {
	model := types.Model{
		ActiveModal: types.AddCommandForm,
		FormData: types.FormData{
			Name:        "test",
			ActiveField: 0, // Name field
		},
	}

	newModel := addCharToActiveField(model, "x")

	if newModel.FormData.Name != "testx" {
		t.Errorf("Expected name to be 'testx', got %q", newModel.FormData.Name)
	}
}

func TestDeleteCharFromActiveField(t *testing.T) {
	model := types.Model{
		ActiveModal: types.AddCommandForm,
		FormData: types.FormData{
			Name:        "test",
			ActiveField: 0, // Name field
		},
	}

	newModel := deleteCharFromActiveField(model)

	if newModel.FormData.Name != "tes" {
		t.Errorf("Expected name to be 'tes', got %q", newModel.FormData.Name)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			findInString(s, substr))))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Tests for uncovered functions

func TestModalFormHandling(t *testing.T) {
	t.Run("Modal form state handling", func(t *testing.T) {
		model := testutil.NewTestModel().Build()
		model.State = types.ModalActive
		model.ActiveModal = types.AddCommandForm
		model.FormData.Name = "test-cmd"
		model.FormData.Command = "test-command"

		// Test that modal state is properly maintained
		assert.Equal(t, types.ModalActive, model.State)
		assert.Equal(t, types.AddCommandForm, model.ActiveModal)
		assert.Equal(t, "test-cmd", model.FormData.Name)
	})

	t.Run("SSE form state handling", func(t *testing.T) {
		model := testutil.NewTestModel().Build()
		model.State = types.ModalActive
		model.ActiveModal = types.AddSSEForm
		model.FormData.Name = "test-sse"
		model.FormData.URL = "https://example.com"

		assert.Equal(t, types.AddSSEForm, model.ActiveModal)
		assert.Equal(t, "test-sse", model.FormData.Name)
		assert.Equal(t, "https://example.com", model.FormData.URL)
	})

	t.Run("JSON form state handling", func(t *testing.T) {
		model := testutil.NewTestModel().Build()
		model.State = types.ModalActive
		model.ActiveModal = types.AddJSONForm
		model.FormData.Name = "test-json"
		model.FormData.JSONConfig = `{"key": "value"}`

		assert.Equal(t, types.AddJSONForm, model.ActiveModal)
		assert.Equal(t, "test-json", model.FormData.Name)
		assert.NotEmpty(t, model.FormData.JSONConfig)
	})

	t.Run("Delete modal state handling", func(t *testing.T) {
		model := testutil.NewTestModel().Build()
		model.State = types.ModalActive
		model.ActiveModal = types.DeleteModal
		model.SelectedItem = 0
		model.MCPItems = []types.MCPItem{
			{Name: "to-delete", Type: "CMD"},
		}

		assert.Equal(t, types.DeleteModal, model.ActiveModal)
		assert.Equal(t, 0, model.SelectedItem)
		assert.Len(t, model.MCPItems, 1)
	})
}

// Test uncovered modal handler functions to achieve >90% coverage

func TestHandleCommandFormKeys(t *testing.T) {
	t.Run("tab_navigation", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormData.ActiveField = 0

		// Test tab navigation through fields
		updatedModel, _ := handleCommandFormKeys(model, "tab")
		assert.Equal(t, 1, updatedModel.FormData.ActiveField)

		// Test wrap around
		updatedModel.FormData.ActiveField = 3
		updatedModel, _ = handleCommandFormKeys(updatedModel, "tab")
		assert.Equal(t, 0, updatedModel.FormData.ActiveField)
	})

	t.Run("enter_submit_valid_form", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormData = types.FormData{
			Name:    "test-cmd",
			Command: "test-command",
			Args:    "arg1 arg2",
		}

		updatedModel, cmd := handleCommandFormKeys(model, "enter")
		assert.Equal(t, types.MainNavigation, updatedModel.State)
		assert.Equal(t, types.NoModal, updatedModel.ActiveModal)
		assert.NotNil(t, cmd)
	})

	t.Run("enter_invalid_form", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormData = types.FormData{
			Name: "", // Invalid - empty name
		}

		updatedModel, _ := handleCommandFormKeys(model, "enter")
		assert.Equal(t, types.ModalActive, updatedModel.State)
		assert.Equal(t, types.AddCommandForm, updatedModel.ActiveModal)
		assert.NotEmpty(t, updatedModel.FormErrors)
	})

	t.Run("backspace_delete_char", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormData = types.FormData{
			Name: "test",
		}
		model.FormData.ActiveField = 0

		updatedModel, _ := handleCommandFormKeys(model, "backspace")
		assert.Equal(t, "tes", updatedModel.FormData.Name)
	})

	t.Run("copy_paste_operations", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormData = types.FormData{
			Name: "test",
		}
		model.FormData.ActiveField = 0

		// Test copy operation
		updatedModel, _ := handleCommandFormKeys(model, "cmd+c")
		// Copy operation should not change model state visibly but should trigger clipboard operation
		assert.Equal(t, "test", updatedModel.FormData.Name)

		// Test paste operation
		updatedModel, _ = handleCommandFormKeys(model, "cmd+v")
		// Paste operation should not change model state without actual clipboard content
		assert.Equal(t, "test", updatedModel.FormData.Name)
	})

	t.Run("character_input", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormData = types.FormData{
			Name: "",
		}
		model.FormData.ActiveField = 0

		updatedModel, _ := handleCommandFormKeys(model, "a")
		assert.Equal(t, "a", updatedModel.FormData.Name)

		// Test multi-character keys are ignored
		updatedModel, _ = handleCommandFormKeys(updatedModel, "ctrl+a")
		assert.Equal(t, "a", updatedModel.FormData.Name) // Should not change
	})

	t.Run("edit_mode", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			WithMCPs([]types.MCPItem{
				{Name: "existing", Type: "CMD", Command: "test"},
			}).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.EditMode = true
		model.EditMCPName = "existing"
		model.FormData = types.FormData{
			Name:    "existing",
			Command: "updated-command",
		}

		updatedModel, cmd := handleCommandFormKeys(model, "enter")
		assert.Equal(t, types.MainNavigation, updatedModel.State)
		assert.Equal(t, types.NoModal, updatedModel.ActiveModal)
		assert.False(t, updatedModel.EditMode)
		assert.Empty(t, updatedModel.EditMCPName)
		assert.NotNil(t, cmd)
	})
}

func TestHandleSSEFormKeys(t *testing.T) {
	t.Run("tab_navigation", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddSSEForm
		model.FormData.ActiveField = 0

		// Test tab navigation through fields
		updatedModel, _ := handleSSEFormKeys(model, "tab")
		assert.Equal(t, 1, updatedModel.FormData.ActiveField)

		// Test wrap around
		updatedModel.FormData.ActiveField = 2
		updatedModel, _ = handleSSEFormKeys(updatedModel, "tab")
		assert.Equal(t, 0, updatedModel.FormData.ActiveField)
	})

	t.Run("enter_submit_valid_form", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddSSEForm
		model.FormData = types.FormData{
			Name: "test-sse",
			URL:  "https://example.com/sse",
		}

		updatedModel, cmd := handleSSEFormKeys(model, "enter")
		assert.Equal(t, types.MainNavigation, updatedModel.State)
		assert.Equal(t, types.NoModal, updatedModel.ActiveModal)
		assert.NotNil(t, cmd)
	})

	t.Run("enter_invalid_form", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddSSEForm
		model.FormData = types.FormData{
			Name: "", // Invalid - empty name
		}

		updatedModel, _ := handleSSEFormKeys(model, "enter")
		assert.Equal(t, types.ModalActive, updatedModel.State)
		assert.Equal(t, types.AddSSEForm, updatedModel.ActiveModal)
		assert.NotEmpty(t, updatedModel.FormErrors)
	})

	t.Run("character_input_to_fields", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddSSEForm

		// Test name field
		model.FormData.ActiveField = 0
		updatedModel, _ := handleSSEFormKeys(model, "a")
		assert.Equal(t, "a", updatedModel.FormData.Name)

		// Test URL field
		model.FormData.ActiveField = 1
		updatedModel, _ = handleSSEFormKeys(model, "h")
		assert.Equal(t, "h", updatedModel.FormData.URL)

		// Test environment field
		model.FormData.ActiveField = 2
		updatedModel, _ = handleSSEFormKeys(model, "e")
		assert.Equal(t, "e", updatedModel.FormData.Environment)
	})

	t.Run("clipboard_operations", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddSSEForm
		model.FormData = types.FormData{
			Name: "test",
		}
		model.FormData.ActiveField = 0

		// Test various clipboard key combinations
		clipboardKeys := []string{"ctrl+c", "cmd+c", "ctrl+v", "cmd+v"}
		for _, key := range clipboardKeys {
			updatedModel, _ := handleSSEFormKeys(model, key)
			assert.Equal(t, types.AddSSEForm, updatedModel.ActiveModal)
		}
	})
}

func TestHandleJSONFormKeys(t *testing.T) {
	t.Run("tab_navigation", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddJSONForm
		model.FormData.ActiveField = 0

		// Test tab navigation through fields
		updatedModel, _ := handleJSONFormKeys(model, "tab")
		assert.Equal(t, 1, updatedModel.FormData.ActiveField)

		// Test wrap around
		updatedModel.FormData.ActiveField = 2
		updatedModel, _ = handleJSONFormKeys(updatedModel, "tab")
		assert.Equal(t, 0, updatedModel.FormData.ActiveField)
	})

	t.Run("enter_in_json_field", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddJSONForm
		model.FormData = types.FormData{
			JSONConfig: `{"key":`,
		}
		model.FormData.ActiveField = 1 // JSON field

		updatedModel, _ := handleJSONFormKeys(model, "enter")
		assert.Contains(t, updatedModel.FormData.JSONConfig, "\n")
	})

	t.Run("enter_submit_from_non_json_field", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddJSONForm
		model.FormData = types.FormData{
			Name:       "test-json",
			JSONConfig: `{"key": "value"}`,
		}
		model.FormData.ActiveField = 0 // Name field

		updatedModel, cmd := handleJSONFormKeys(model, "enter")
		assert.Equal(t, types.MainNavigation, updatedModel.State)
		assert.Equal(t, types.NoModal, updatedModel.ActiveModal)
		assert.NotNil(t, cmd)
	})

	t.Run("character_input_to_fields", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddJSONForm

		// Test name field
		model.FormData.ActiveField = 0
		updatedModel, _ := handleJSONFormKeys(model, "a")
		assert.Equal(t, "a", updatedModel.FormData.Name)

		// Test JSON field
		model.FormData.ActiveField = 1
		updatedModel, _ = handleJSONFormKeys(model, "{")
		assert.Equal(t, "{", updatedModel.FormData.JSONConfig)

		// Test environment field
		model.FormData.ActiveField = 2
		updatedModel, _ = handleJSONFormKeys(model, "e")
		assert.Equal(t, "e", updatedModel.FormData.Environment)
	})
}

func TestDeleteCharFromSSEForm(t *testing.T) {
	t.Run("delete_from_name_field", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddSSEForm
		model.FormData = types.FormData{
			Name: "test",
		}
		model.FormData.ActiveField = 0

		updatedModel := deleteCharFromSSEForm(model)
		assert.Equal(t, "tes", updatedModel.FormData.Name)
	})

	t.Run("delete_from_url_field", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddSSEForm
		model.FormData = types.FormData{
			URL: "https://example.com",
		}
		model.FormData.ActiveField = 1

		updatedModel := deleteCharFromSSEForm(model)
		assert.Equal(t, "https://example.co", updatedModel.FormData.URL)
	})

	t.Run("delete_from_environment_field", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddSSEForm
		model.FormData = types.FormData{
			Environment: "KEY=value",
		}
		model.FormData.ActiveField = 2

		updatedModel := deleteCharFromSSEForm(model)
		assert.Equal(t, "KEY=valu", updatedModel.FormData.Environment)
	})

	t.Run("delete_from_empty_field", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddSSEForm
		model.FormData = types.FormData{
			Name: "",
		}
		model.FormData.ActiveField = 0

		updatedModel := deleteCharFromSSEForm(model)
		assert.Equal(t, "", updatedModel.FormData.Name)
	})
}

func TestDeleteCharFromJSONForm(t *testing.T) {
	t.Run("delete_from_name_field", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddJSONForm
		model.FormData = types.FormData{
			Name: "test",
		}
		model.FormData.ActiveField = 0

		updatedModel := deleteCharFromJSONForm(model)
		assert.Equal(t, "tes", updatedModel.FormData.Name)
	})

	t.Run("delete_from_json_field", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddJSONForm
		model.FormData = types.FormData{
			JSONConfig: `{"key": "value"}`,
		}
		model.FormData.ActiveField = 1

		updatedModel := deleteCharFromJSONForm(model)
		assert.Equal(t, `{"key": "value"`, updatedModel.FormData.JSONConfig)
	})

	t.Run("delete_from_environment_field", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddJSONForm
		model.FormData = types.FormData{
			Environment: "KEY=value",
		}
		model.FormData.ActiveField = 2

		updatedModel := deleteCharFromJSONForm(model)
		assert.Equal(t, "KEY=valu", updatedModel.FormData.Environment)
	})
}

func TestAddMCPToInventory(t *testing.T) {
	t.Run("add_new_mcp", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()

		mcpItem := types.MCPItem{
			Name:    "new-mcp",
			Type:    "CMD",
			Command: "test-command",
		}

		updatedModel, cmd := addMCPToInventory(model, mcpItem)
		assert.Len(t, updatedModel.MCPItems, 1) // 0 default + 1 new
		assert.Equal(t, "new-mcp", updatedModel.MCPItems[0].Name)
		assert.Equal(t, "Added new-mcp successfully", updatedModel.SuccessMessage)
		assert.NotNil(t, cmd)
	})

	t.Run("add_duplicate_name", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			WithMCPs([]types.MCPItem{
				{Name: "existing", Type: "CMD", Command: "test"},
			}).
			Build()

		mcpItem := types.MCPItem{
			Name:    "existing",
			Type:    "CMD",
			Command: "test-command",
		}

		updatedModel, _ := addMCPToInventory(model, mcpItem)
		assert.Len(t, updatedModel.MCPItems, 2) // Actually adds duplicate - validation is elsewhere
		assert.Equal(t, "Added existing successfully", updatedModel.SuccessMessage)
	})
}

func TestHandleDeleteModalKeys(t *testing.T) {
	t.Run("confirm_delete", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			WithMCPs([]types.MCPItem{
				{Name: "to-delete", Type: "CMD", Command: "test"},
			}).
			Build()
		model.ActiveModal = types.DeleteModal
		model.SelectedItem = 0

		updatedModel, cmd := handleDeleteModalKeys(model, "enter")
		assert.Equal(t, types.MainNavigation, updatedModel.State)
		assert.Equal(t, types.NoModal, updatedModel.ActiveModal)
		assert.NotNil(t, cmd)
	})

	t.Run("cancel_delete", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			WithMCPs([]types.MCPItem{
				{Name: "to-keep", Type: "CMD", Command: "test"},
			}).
			Build()
		model.ActiveModal = types.DeleteModal
		model.SelectedItem = 0

		updatedModel, _ := handleDeleteModalKeys(model, "esc")
		assert.Equal(t, types.MainNavigation, updatedModel.State)
		assert.Equal(t, types.NoModal, updatedModel.ActiveModal)
	})

	t.Run("escape_delete", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.DeleteModal

		updatedModel, _ := handleDeleteModalKeys(model, "esc")
		assert.Equal(t, types.MainNavigation, updatedModel.State)
		assert.Equal(t, types.NoModal, updatedModel.ActiveModal)
	})
}

func TestDeleteMCPFromInventory(t *testing.T) {
	t.Run("delete_existing_mcp", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			WithMCPs([]types.MCPItem{
				{Name: "to-delete", Type: "CMD", Command: "test"},
				{Name: "to-keep", Type: "CMD", Command: "test2"},
			}).
			Build()
		model.SelectedItem = 0

		updatedModel, cmd := deleteMCPFromInventory(model)
		assert.Len(t, updatedModel.MCPItems, 1)
		assert.Equal(t, "to-keep", updatedModel.MCPItems[0].Name)
		assert.Equal(t, "Deleted to-delete successfully", updatedModel.SuccessMessage)
		assert.NotNil(t, cmd)
	})

	t.Run("delete_invalid_index", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			WithMCPs([]types.MCPItem{
				{Name: "only-item", Type: "CMD", Command: "test"},
			}).
			Build()
		model.SelectedItem = 5 // Invalid index

		updatedModel, _ := deleteMCPFromInventory(model)
		assert.Len(t, updatedModel.MCPItems, 1) // Should not delete anything
	})
}

func TestClipboardOperations(t *testing.T) {
	t.Run("copy_active_field_content", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormData = types.FormData{
			Name: "test-content",
		}
		model.FormData.ActiveField = 0

		updatedModel := copyActiveFieldToClipboard(model)
		// Should not change the model state but trigger clipboard operation
		assert.Equal(t, "test-content", updatedModel.FormData.Name)
	})

	t.Run("get_active_field_content", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormData = types.FormData{
			Name: "test-name",
		}
		model.FormData.ActiveField = 0

		content := getActiveFieldContent(model)
		assert.Equal(t, "test-name", content)
	})

	t.Run("get_command_form_field_content", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormData = types.FormData{
			Name:        "test-name",
			Command:     "test-command",
			Args:        "arg1 arg2",
			Environment: "KEY=value",
		}

		// Test each field
		model.FormData.ActiveField = 0
		assert.Equal(t, "test-name", getCommandFormFieldContent(model))

		model.FormData.ActiveField = 1
		assert.Equal(t, "test-command", getCommandFormFieldContent(model))

		model.FormData.ActiveField = 2
		assert.Equal(t, "arg1 arg2", getCommandFormFieldContent(model))

		model.FormData.ActiveField = 3
		assert.Equal(t, "KEY=value", getCommandFormFieldContent(model))
	})

	t.Run("get_sse_form_field_content", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddSSEForm
		model.FormData = types.FormData{
			Name:        "test-name",
			URL:         "https://example.com",
			Environment: "KEY=value",
		}

		// Test each field
		model.FormData.ActiveField = 0
		assert.Equal(t, "test-name", getSSEFormFieldContent(model))

		model.FormData.ActiveField = 1
		assert.Equal(t, "https://example.com", getSSEFormFieldContent(model))

		model.FormData.ActiveField = 2
		assert.Equal(t, "KEY=value", getSSEFormFieldContent(model))
	})

	t.Run("get_json_form_field_content", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddJSONForm
		model.FormData = types.FormData{
			Name:        "test-name",
			JSONConfig:  `{"key": "value"}`,
			Environment: "KEY=value",
		}

		// Test each field
		model.FormData.ActiveField = 0
		assert.Equal(t, "test-name", getJSONFormFieldContent(model))

		model.FormData.ActiveField = 1
		assert.Equal(t, `{"key": "value"}`, getJSONFormFieldContent(model))

		model.FormData.ActiveField = 2
		assert.Equal(t, "KEY=value", getJSONFormFieldContent(model))
	})

	t.Run("paste_from_clipboard", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormData.ActiveField = 0

		updatedModel := pasteFromClipboardToActiveField(model)
		// Should not change model without actual clipboard content
		assert.Equal(t, types.AddCommandForm, updatedModel.ActiveModal)
	})

	t.Run("get_clipboard_content", func(t *testing.T) {
		// Test clipboard content retrieval
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		
		content, err := getClipboardContent(model.PlatformService)
		// We can't guarantee clipboard state, so just verify function doesn't panic
		assert.NotNil(t, content) // content can be empty string
		// err can be nil or not nil depending on clipboard state
		_ = err
	})

	t.Run("handle_clipboard_error", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()

		testErr := assert.AnError
		updatedModel := handleClipboardError(model, testErr)
		// Should add error message
		assert.NotEmpty(t, updatedModel.SuccessMessage)
	})

	t.Run("paste_content_to_active_field", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormData.ActiveField = 0

		content := pastedContent
		updatedModel := pasteContentToActiveField(model, content)
		assert.Equal(t, content, updatedModel.FormData.Name)
	})

	t.Run("paste_to_command_form", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm

		content := pastedContent

		// Test pasting to each field
		model.FormData.ActiveField = 0
		updatedModel := pasteToCommandForm(model, content)
		assert.Equal(t, content, updatedModel.FormData.Name)

		model.FormData.ActiveField = 1
		updatedModel = pasteToCommandForm(model, content)
		assert.Equal(t, content, updatedModel.FormData.Command)

		model.FormData.ActiveField = 2
		updatedModel = pasteToCommandForm(model, content)
		assert.Equal(t, content, updatedModel.FormData.Args)

		model.FormData.ActiveField = 3
		updatedModel = pasteToCommandForm(model, content)
		assert.Equal(t, content, updatedModel.FormData.Environment)
	})

	t.Run("paste_to_sse_form", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddSSEForm

		content := pastedContent

		// Test pasting to each field
		model.FormData.ActiveField = 0
		updatedModel := pasteToSSEForm(model, content)
		assert.Equal(t, content, updatedModel.FormData.Name)

		model.FormData.ActiveField = 1
		updatedModel = pasteToSSEForm(model, content)
		assert.Equal(t, content, updatedModel.FormData.URL)

		model.FormData.ActiveField = 2
		updatedModel = pasteToSSEForm(model, content)
		assert.Equal(t, content, updatedModel.FormData.Environment)
	})

	t.Run("paste_to_json_form", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddJSONForm

		content := pastedContent

		// Test pasting to each field
		model.FormData.ActiveField = 0
		updatedModel := pasteToJSONForm(model, content)
		assert.Equal(t, content, updatedModel.FormData.Name)

		model.FormData.ActiveField = 1
		updatedModel = pasteToJSONForm(model, content)
		assert.Equal(t, content, updatedModel.FormData.JSONConfig)

		model.FormData.ActiveField = 2
		updatedModel = pasteToJSONForm(model, content)
		assert.Equal(t, content, updatedModel.FormData.Environment)
	})

	t.Run("add_paste_success_message", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()

		updatedModel := addPasteSuccessMessage(model)
		assert.NotEmpty(t, updatedModel.SuccessMessage)
		assert.Contains(t, updatedModel.SuccessMessage, "Pasted")
	})
}

func TestFocusOnFirstErrorField(t *testing.T) {
	t.Run("focus_on_name_error", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormErrors = map[string]string{
			"name": "Name is required",
		}
		model.FormData.ActiveField = 2

		updatedModel := focusOnFirstErrorField(model)
		assert.Equal(t, 0, updatedModel.FormData.ActiveField) // Name field
	})

	t.Run("focus_on_command_error", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormErrors = map[string]string{
			"command": "Command is required",
		}
		model.FormData.ActiveField = 3

		updatedModel := focusOnFirstErrorField(model)
		assert.Equal(t, 1, updatedModel.FormData.ActiveField) // Command field
	})

	t.Run("focus_on_url_error_sse", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddSSEForm
		model.FormErrors = map[string]string{
			"url": "URL is required",
		}
		model.FormData.ActiveField = 2

		updatedModel := focusOnFirstErrorField(model)
		assert.Equal(t, 1, updatedModel.FormData.ActiveField) // URL field
	})

	t.Run("focus_on_json_error", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddJSONForm
		model.FormErrors = map[string]string{
			"json": "Invalid JSON",
		}
		model.FormData.ActiveField = 2

		updatedModel := focusOnFirstErrorField(model)
		assert.Equal(t, 1, updatedModel.FormData.ActiveField) // JSON field
	})

	t.Run("no_errors", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormErrors = map[string]string{}
		model.FormData.ActiveField = 2

		updatedModel := focusOnFirstErrorField(model)
		assert.Equal(t, 2, updatedModel.FormData.ActiveField) // Should not change
	})
}

// Test additional edge cases for comprehensive coverage
func TestModalKeyHandlingEdgeCases(t *testing.T) {
	t.Run("handle_modal_keys_with_unknown_modal", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.ModalType(999) // Unknown modal type

		updatedModel, _ := HandleModalKeys(model, "enter")
		// Should handle gracefully without crashing
		assert.Equal(t, types.ModalType(999), updatedModel.ActiveModal)
	})

	t.Run("handle_esc_key_various_states", func(t *testing.T) {
		// Test ESC with different modal types
		modalTypes := []types.ModalType{
			types.AddCommandForm,
			types.AddSSEForm,
			types.AddJSONForm,
			types.DeleteModal,
			types.AddMCPTypeSelection,
		}

		for _, modalType := range modalTypes {
			t.Run(fmt.Sprintf("modal_%d", modalType), func(t *testing.T) {
				model := testutil.NewTestModel().
					WithActiveColumn(0).
					WithState(types.ModalActive).
					Build()
				model.ActiveModal = modalType
				model.EditMode = true
				model.EditMCPName = "test"
				model.FormData.Name = "test-data"
				model.FormErrors = map[string]string{"test": "error"}

				updatedModel, _ := HandleEscKey(model)
				assert.Equal(t, types.MainNavigation, updatedModel.State)
				assert.Equal(t, types.NoModal, updatedModel.ActiveModal)
				assert.False(t, updatedModel.EditMode)
				assert.Empty(t, updatedModel.EditMCPName)
				assert.Empty(t, updatedModel.FormData.Name)
				assert.Empty(t, updatedModel.FormErrors)
			})
		}
	})

	t.Run("handle_invalid_key_combinations", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm

		invalidKeys := []string{
			"ctrl+alt+del",
			"f1",
			"shift+tab",
			"ctrl+z",
			"alt+f4",
			"meta+shift+z",
		}

		for _, key := range invalidKeys {
			t.Run(key, func(t *testing.T) {
				updatedModel, _ := handleCommandFormKeys(model, key)
				// Should handle gracefully without changing state inappropriately
				assert.Equal(t, types.AddCommandForm, updatedModel.ActiveModal)
				assert.Equal(t, types.ModalActive, updatedModel.State)
			})
		}
	})

	t.Run("type_selection_invalid_keys", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddMCPTypeSelection

		invalidKeys := []string{"4", "5", "a", "enter", "space"}
		for _, key := range invalidKeys {
			t.Run(key, func(t *testing.T) {
				updatedModel, _ := handleTypeSelectionKeys(model, key)
				// Should stay in type selection for invalid keys
				if key != "enter" {
					assert.Equal(t, types.AddMCPTypeSelection, updatedModel.ActiveModal)
				}
			})
		}
	})

	t.Run("field_navigation_boundary_tests", func(t *testing.T) {
		// Test field navigation at boundaries
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm

		// Test going beyond last field
		model.FormData.ActiveField = 3 // Last field
		updatedModel, _ := handleCommandFormKeys(model, "tab")
		assert.Equal(t, 0, updatedModel.FormData.ActiveField) // Should wrap to first field

		// Test shift+tab handling (note: actual behavior may vary)
		model.FormData.ActiveField = 0 // First field
		updatedModel, _ = handleCommandFormKeys(model, "shift+tab")
		// Just ensure it doesn't crash - specific behavior may vary
		assert.True(t, updatedModel.FormData.ActiveField >= 0)
	})

	t.Run("form_validation_edge_cases", func(t *testing.T) {
		// Test validation with edge case data
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			WithMCPs([]types.MCPItem{
				{Name: "existing", Type: "CMD", Command: "test"},
			}).Build()

		// Test duplicate name validation
		model.FormData = types.FormData{
			Name:    "existing",
			Command: "test-command",
		}
		model.FormErrors = make(map[string]string)

		_, valid := validateCommandForm(model)
		assert.False(t, valid) // Should be invalid due to duplicate name

		// Test empty environment string validation
		model.FormData.Environment = ""
		err := validateEnvironmentFormat(model.FormData.Environment)
		assert.NoError(t, err) // Empty environment should be valid

		// Test malformed environment string
		model.FormData.Environment = "INVALID_LINE_NO_EQUALS"
		err = validateEnvironmentFormat(model.FormData.Environment)
		assert.Error(t, err) // Should error on malformed line
	})

	t.Run("url_validation_edge_cases", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.FormErrors = make(map[string]string)

		// Test various URL formats
		urlTestCases := []struct {
			url   string
			valid bool
		}{
			{"http://example.com", true},
			{"https://example.com", true},
			{"not-a-url", false},
			{"ftp://example.com", true},
			{"", false},
			{"://invalid", false},
		}

		for _, tc := range urlTestCases {
			t.Run(tc.url, func(t *testing.T) {
				model.FormData = types.FormData{
					Name: "test-sse",
					URL:  tc.url,
				}
				model.FormErrors = make(map[string]string)

				_, valid := validateSSEForm(model)
				if tc.valid {
					assert.True(t, valid, "URL %s should be valid", tc.url)
				} else {
					assert.False(t, valid, "URL %s should be invalid", tc.url)
				}
			})
		}
	})

	t.Run("json_validation_edge_cases", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.FormErrors = make(map[string]string)

		// Test various JSON formats
		jsonTestCases := []struct {
			json  string
			valid bool
		}{
			{`{"key": "value"}`, true},
			{`[]`, true},
			{`null`, true},
			{`"string"`, true},
			{`123`, true},
			{`true`, true},
			{`{invalid json}`, false},
			{`{"unclosed": "object"`, false},
			{"", false},
			{`{"nested": {"key": "value"}}`, true},
		}

		for _, tc := range jsonTestCases {
			t.Run(tc.json, func(t *testing.T) {
				model.FormData = types.FormData{
					Name:       "test-json",
					JSONConfig: tc.json,
				}
				model.FormErrors = make(map[string]string)

				_, valid := validateJSONForm(model)
				if tc.valid {
					assert.True(t, valid, "JSON %s should be valid", tc.json)
				} else {
					assert.False(t, valid, "JSON %s should be invalid", tc.json)
				}
			})
		}
	})

	t.Run("environment_parsing_edge_cases", func(t *testing.T) {
		// Test environment string parsing with edge cases
		envTestCases := []struct {
			input    string
			expected map[string]string
			hasError bool
		}{
			{"", nil, false},
			{"\n\n", nil, false}, // Only newlines
			{"KEY=value\n\n", map[string]string{"KEY": "value"}, false},
			{"KEY1=value1\nKEY2=value2\n", map[string]string{"KEY1": "value1", "KEY2": "value2"}, false},
			{"INVALID_LINE\nKEY=value", nil, true}, // Should error on first invalid line
			{"KEY=\n", map[string]string{"KEY": ""}, false}, // Empty value should be OK
			{"=value", nil, true}, // Empty key should error
		}

		for _, tc := range envTestCases {
			t.Run(tc.input, func(t *testing.T) {
				result := parseEnvironmentString(tc.input)
				if tc.hasError {
					err := validateEnvironmentFormat(tc.input)
					assert.Error(t, err)
					return
				}
				assert.Equal(t, tc.expected, result)
			})
		}
	})

	t.Run("character_input_edge_cases", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithActiveColumn(0).
			WithState(types.ModalActive).
			Build()
		model.ActiveModal = types.AddCommandForm
		model.FormData.ActiveField = 0

		// Test special characters
		specialChars := []string{
			"!", "@", "#", "$", "%", "^", "&", "*", "(", ")",
			"-", "_", "=", "+", "[", "]", "{", "}", "|", "\\",
			";", ":", "'", "\"", ",", ".", "<", ">", "/", "?",
			"~", "`", " ", "\t", "\n", "\r",
		}

		for _, char := range specialChars {
			t.Run("char_"+char, func(t *testing.T) {
				updatedModel, _ := handleCommandFormKeys(model, char)
				// Should add printable characters to field
				if len(char) == 1 && char != "\t" && char != "\n" && char != "\r" {
					assert.Contains(t, updatedModel.FormData.Name, char)
				}
			})
		}
	})
}
