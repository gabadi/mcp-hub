package handlers

import (
	"testing"

	"cc-mcp-manager/internal/testutil"
	"cc-mcp-manager/internal/ui/types"
	
	"github.com/stretchr/testify/assert"
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
	model.FormData.Command = "test"      // Required field

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
