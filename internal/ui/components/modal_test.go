package components

import (
	"strings"
	"testing"

	"cc-mcp-manager/internal/ui/types"
)

// Test message constants
const (
	NoMCPSelected = "No MCP selected"
)

func TestOverlayModal(t *testing.T) {
	model := types.NewModel()
	model.ActiveModal = types.AddCommandForm
	width, height := 120, 40

	result := OverlayModal(model, width, height, "background content")

	if result == "" {
		t.Error("OverlayModal should return non-empty string")
	}

	// Should contain modal title
	if !strings.Contains(result, "Add New MCP") {
		t.Error("Modal should contain title")
	}
}

func TestOverlayModalDifferentTypes(t *testing.T) {
	tests := []struct {
		name        string
		modalType   types.ModalType
		expectedStr string
	}{
		{
			name:        "Type Selection Modal",
			modalType:   types.AddMCPTypeSelection,
			expectedStr: "Select Type",
		},
		{
			name:        "Command Form Modal",
			modalType:   types.AddCommandForm,
			expectedStr: "Command/Binary",
		},
		{
			name:        "SSE Form Modal",
			modalType:   types.AddSSEForm,
			expectedStr: "SSE Server",
		},
		{
			name:        "JSON Form Modal",
			modalType:   types.AddJSONForm,
			expectedStr: "JSON Configuration",
		},
		{
			name:        "Edit Modal",
			modalType:   types.EditModal,
			expectedStr: "Edit MCP",
		},
		{
			name:        "Delete Modal",
			modalType:   types.DeleteModal,
			expectedStr: "Delete MCP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := types.NewModel()
			model.ActiveModal = tt.modalType

			// For edit/delete modals, ensure we have an MCP to work with
			if tt.modalType == types.EditModal || tt.modalType == types.DeleteModal {
				model.MCPItems = []types.MCPItem{
					{Name: "test-mcp", Type: "CMD", Command: "test"},
				}
				model.SelectedItem = 0
			}

			result := OverlayModal(model, 120, 40, "background")

			if !strings.Contains(result, tt.expectedStr) {
				t.Errorf("Modal should contain %q, got: %s", tt.expectedStr, result)
			}
		})
	}
}

func TestOverlayModalSizeConstraints(t *testing.T) {
	model := types.NewModel()
	model.ActiveModal = types.AddCommandForm

	// Test with very small dimensions
	result := OverlayModal(model, 50, 20, "background")
	if result == "" {
		t.Error("Modal should handle small dimensions gracefully")
	}

	// Test with JSON form (larger modal)
	model.ActiveModal = types.AddJSONForm
	result = OverlayModal(model, 120, 40, "background")
	if result == "" {
		t.Error("JSON modal should render correctly")
	}
}

func TestRenderTypeSelectionContent(t *testing.T) {
	model := types.NewModel()
	model.FormData.ActiveField = 1 // Select first option

	result := renderTypeSelectionContent(model)

	// Should contain all three options
	if !strings.Contains(result, "Command/Binary") {
		t.Error("Should contain Command/Binary option")
	}
	if !strings.Contains(result, "SSE Server") {
		t.Error("Should contain SSE Server option")
	}
	if !strings.Contains(result, "JSON Configuration") {
		t.Error("Should contain JSON Configuration option")
	}

	// Should contain instructions
	if !strings.Contains(result, "Use number keys") {
		t.Error("Should contain usage instructions")
	}
}

func TestRenderTypeSelectionContentHighlight(t *testing.T) {
	tests := []struct {
		activeField    int
		expectedNumber string
	}{
		{0, "1"}, // Default to option 1 when 0
		{1, "1"},
		{2, "2"},
		{3, "3"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			model := types.NewModel()
			model.FormData.ActiveField = tt.activeField

			result := renderTypeSelectionContent(model)

			// The selected option should be highlighted (we can't easily test styling,
			// but we can verify the content is there)
			if !strings.Contains(result, tt.expectedNumber+".") {
				t.Errorf("Should contain option %s", tt.expectedNumber)
			}
		})
	}
}

func TestRenderCommandFormContent(t *testing.T) {
	model := types.NewModel()
	model.FormData.Name = "test-mcp"
	model.FormData.Command = "test-command"
	model.FormData.Args = "arg1 arg2"
	model.FormData.Environment = "VAR=value"
	model.FormData.ActiveField = 0

	result := renderCommandFormContent(model)

	// Should contain all field labels
	if !strings.Contains(result, "Name:") {
		t.Error("Should contain Name field")
	}
	if !strings.Contains(result, "Command:") {
		t.Error("Should contain Command field")
	}
	if !strings.Contains(result, "Args:") {
		t.Error("Should contain Args field")
	}
	if !strings.Contains(result, "Environment:") {
		t.Error("Should contain Environment field")
	}

	// Should contain field values
	if !strings.Contains(result, "test-mcp") {
		t.Error("Should contain name value")
	}
	if !strings.Contains(result, "test-command") {
		t.Error("Should contain command value")
	}

	// Should show cursor for active field
	if !strings.Contains(result, "test-mcp_") {
		t.Error("Should show cursor in active field")
	}
}

func TestRenderCommandFormContentWithErrors(t *testing.T) {
	model := types.NewModel()
	model.FormErrors = map[string]string{
		"name":    "Name is required",
		"command": "Command is required",
	}

	result := renderCommandFormContent(model)

	// Should contain error messages
	if !strings.Contains(result, "Name is required") {
		t.Error("Should display name error")
	}
	if !strings.Contains(result, "Command is required") {
		t.Error("Should display command error")
	}
}

func TestRenderSSEFormContent(t *testing.T) {
	model := types.NewModel()
	model.FormData.Name = "sse-mcp"
	model.FormData.URL = "http://example.com"
	model.FormData.Environment = "API_KEY=secret"
	model.FormData.ActiveField = 1

	result := renderSSEFormContent(model)

	// Should contain all field labels
	if !strings.Contains(result, "Name:") {
		t.Error("Should contain Name field")
	}
	if !strings.Contains(result, "URL:") {
		t.Error("Should contain URL field")
	}
	if !strings.Contains(result, "Environment:") {
		t.Error("Should contain Environment field")
	}

	// Should contain field values
	if !strings.Contains(result, "sse-mcp") {
		t.Error("Should contain name value")
	}
	if !strings.Contains(result, "http://example.com") {
		t.Error("Should contain URL value")
	}

	// Should show cursor for active field (URL)
	if !strings.Contains(result, "http://example.com_") {
		t.Error("Should show cursor in active URL field")
	}

	// Should contain format instructions
	if !strings.Contains(result, "Enter a valid HTTP/HTTPS URL") {
		t.Error("Should contain URL format instructions")
	}
}

func TestRenderSSEFormContentWithErrors(t *testing.T) {
	model := types.NewModel()
	model.FormErrors = map[string]string{
		"url": "Invalid URL format",
	}

	result := renderSSEFormContent(model)

	// Should contain error message
	if !strings.Contains(result, "Invalid URL format") {
		t.Error("Should display URL error")
	}
}

func TestRenderJSONFormContent(t *testing.T) {
	model := types.NewModel()
	model.FormData.Name = "json-mcp"
	model.FormData.JSONConfig = `{"key": "value"}`
	model.FormData.Environment = "ENV=test"
	model.FormData.ActiveField = 0

	result := renderJSONFormContent(model)

	// Should contain all field labels
	if !strings.Contains(result, "Name:") {
		t.Error("Should contain Name field")
	}
	if !strings.Contains(result, "JSON Configuration:") {
		t.Error("Should contain JSON Configuration field")
	}
	if !strings.Contains(result, "Environment:") {
		t.Error("Should contain Environment field")
	}

	// Should contain field values
	if !strings.Contains(result, "json-mcp") {
		t.Error("Should contain name value")
	}
	if !strings.Contains(result, `"key"`) {
		t.Error("Should contain JSON config value")
	}

	// Should show cursor for active field
	if !strings.Contains(result, "json-mcp_") {
		t.Error("Should show cursor in active name field")
	}
}

func TestRenderJSONFormContentWithValidJSON(t *testing.T) {
	model := types.NewModel()
	model.FormData.JSONConfig = `{"valid": true}`

	result := renderJSONFormContent(model)

	// Should show valid JSON indicator
	if !strings.Contains(result, "✓ Valid JSON") {
		t.Error("Should show valid JSON indicator")
	}
}

func TestRenderJSONFormContentWithErrors(t *testing.T) {
	model := types.NewModel()
	model.FormErrors = map[string]string{
		"json": "Invalid JSON syntax",
	}

	result := renderJSONFormContent(model)

	// Should contain error message
	if !strings.Contains(result, "Invalid JSON syntax") {
		t.Error("Should display JSON error")
	}
}

func TestRenderJSONFormContentMultilineJSON(t *testing.T) {
	model := types.NewModel()
	model.FormData.JSONConfig = "{\n  \"key\": \"value\",\n  \"number\": 42\n}"

	result := renderJSONFormContent(model)

	// Should handle multiline JSON
	if !strings.Contains(result, "key") && !strings.Contains(result, "value") {
		t.Error("Should display multiline JSON content")
	}
}

func TestRenderEditModalContent(t *testing.T) {
	model := types.NewModel()
	model.MCPItems = []types.MCPItem{
		{Name: "test-mcp", Type: "CMD", Command: "test-cmd", Active: true},
	}
	model.SelectedItem = 0

	result := renderEditModalContent(model)

	// Should contain MCP details
	if !strings.Contains(result, "test-mcp") {
		t.Error("Should contain MCP name")
	}
	if !strings.Contains(result, "CMD") {
		t.Error("Should contain MCP type")
	}
	if !strings.Contains(result, "test-cmd") {
		t.Error("Should contain MCP command")
	}

	// Should show active status
	if !strings.Contains(result, "[X] Yes") {
		t.Error("Should show active status correctly")
	}
	if !strings.Contains(result, "[ ] No") {
		t.Error("Should show inactive option correctly")
	}

	// Should contain instructions
	if !strings.Contains(result, "Edit the fields") {
		t.Error("Should contain edit instructions")
	}
}

func TestRenderEditModalContentInactiveMCP(t *testing.T) {
	model := types.NewModel()
	model.MCPItems = []types.MCPItem{
		{Name: "inactive-mcp", Type: "SSE", URL: "http://test.com", Active: false},
	}
	model.SelectedItem = 0

	result := renderEditModalContent(model)

	// Should show inactive status
	if !strings.Contains(result, "[ ] Yes") {
		t.Error("Should show inactive Yes option")
	}
	if !strings.Contains(result, "[X] No") {
		t.Error("Should show active No option")
	}
}

func TestRenderEditModalContentNoSelection(t *testing.T) {
	model := types.NewModel()
	model.MCPItems = []types.MCPItem{}
	model.SelectedItem = 0

	result := renderEditModalContent(model)

	if result != NoMCPSelected {
		t.Errorf("Expected '%s', got: %s", NoMCPSelected, result)
	}
}

func TestRenderEditModalContentWithSearch(t *testing.T) {
	model := types.NewModel()
	model.MCPItems = []types.MCPItem{
		{Name: "test-mcp", Type: "CMD", Command: "test-cmd", Active: true},
		{Name: "other-mcp", Type: "SSE", URL: "http://test.com", Active: false},
	}
	model.SearchQuery = "test"
	model.FilteredSelectedIndex = 0

	result := renderEditModalContent(model)

	// Should show the filtered MCP
	if !strings.Contains(result, "test-mcp") {
		t.Error("Should show filtered MCP for editing")
	}
}

func TestRenderDeleteModalContent(t *testing.T) {
	model := types.NewModel()
	model.MCPItems = []types.MCPItem{
		{Name: "delete-me", Type: "CMD", Command: "delete-cmd", Active: false},
	}
	model.SelectedItem = 0

	result := renderDeleteModalContent(model)

	// Should contain MCP details
	if !strings.Contains(result, "delete-me") {
		t.Error("Should contain MCP name")
	}
	if !strings.Contains(result, "CMD") {
		t.Error("Should contain MCP type")
	}
	if !strings.Contains(result, "delete-cmd") {
		t.Error("Should contain MCP command")
	}

	// Should contain warning message
	if !strings.Contains(result, "Are you sure") {
		t.Error("Should contain confirmation message")
	}
	if !strings.Contains(result, "cannot be undone") {
		t.Error("Should contain warning about permanent deletion")
	}

	// Should contain instructions
	if !strings.Contains(result, "Press Enter to confirm") {
		t.Error("Should contain confirmation instructions")
	}
}

func TestRenderDeleteModalContentNoSelection(t *testing.T) {
	model := types.NewModel()
	model.MCPItems = []types.MCPItem{}
	model.SelectedItem = 0

	result := renderDeleteModalContent(model)

	if result != NoMCPSelected {
		t.Errorf("Expected '%s', got: %s", NoMCPSelected, result)
	}
}

func TestRenderDeleteModalContentWithSearch(t *testing.T) {
	model := types.NewModel()
	model.MCPItems = []types.MCPItem{
		{Name: "keep-me", Type: "CMD", Command: "keep-cmd", Active: true},
		{Name: "delete-me", Type: "SSE", URL: "http://test.com", Active: false},
	}
	model.SearchQuery = "delete"
	model.FilteredSelectedIndex = 0

	result := renderDeleteModalContent(model)

	// Should show the filtered MCP for deletion
	if !strings.Contains(result, "delete-me") {
		t.Error("Should show filtered MCP for deletion")
	}
	if strings.Contains(result, "keep-me") {
		t.Error("Should not show non-filtered MCP")
	}
}

func TestModalFieldFocus(t *testing.T) {
	tests := getModalFieldFocusTestCases()

	for _, tt := range tests {
		model := setupModalFieldFocusTest(tt)
		result := renderModalForFocusTest(model, tt.modalType)
		assertModalFieldFocus(t, result, tt)
	}
}

func getModalFieldFocusTestCases() []struct {
	modalType    types.ModalType
	activeField  int
	fieldValue   string
	expectCursor bool
} {
	return []struct {
		modalType    types.ModalType
		activeField  int
		fieldValue   string
		expectCursor bool
	}{
		{types.AddCommandForm, 0, "name-value", true},
		{types.AddCommandForm, 1, "command-value", true},
		{types.AddCommandForm, 2, "args-value", true},
		{types.AddCommandForm, 3, "env-value", true},
		{types.AddSSEForm, 0, "name-value", true},
		{types.AddSSEForm, 1, "url-value", true},
		{types.AddSSEForm, 2, "env-value", true},
		{types.AddJSONForm, 0, "name-value", true},
		{types.AddJSONForm, 1, "json-value", true},
		{types.AddJSONForm, 2, "env-value", true},
	}
}

func setupModalFieldFocusTest(tt struct {
	modalType    types.ModalType
	activeField  int
	fieldValue   string
	expectCursor bool
}) types.Model {
	model := types.NewModel()
	model.ActiveModal = tt.modalType
	model.FormData.ActiveField = tt.activeField

	setModalFieldValue(&model, tt.modalType, tt.activeField, tt.fieldValue)
	return model
}

func setModalFieldValue(model *types.Model, modalType types.ModalType, activeField int, fieldValue string) {
	switch modalType {
	case types.AddCommandForm:
		setCommandFormFieldValue(model, activeField, fieldValue)
	case types.AddSSEForm:
		setSSEFormFieldValue(model, activeField, fieldValue)
	case types.AddJSONForm:
		setJSONFormFieldValue(model, activeField, fieldValue)
	case types.NoModal, types.AddModal, types.AddMCPTypeSelection, types.EditModal, types.DeleteModal:
		// These modal types don't have form fields to set
	}
}

func setCommandFormFieldValue(model *types.Model, activeField int, fieldValue string) {
	switch activeField {
	case 0:
		model.FormData.Name = fieldValue
	case 1:
		model.FormData.Command = fieldValue
	case 2:
		model.FormData.Args = fieldValue
	case 3:
		model.FormData.Environment = fieldValue
	}
}

func setSSEFormFieldValue(model *types.Model, activeField int, fieldValue string) {
	switch activeField {
	case 0:
		model.FormData.Name = fieldValue
	case 1:
		model.FormData.URL = fieldValue
	case 2:
		model.FormData.Environment = fieldValue
	}
}

func setJSONFormFieldValue(model *types.Model, activeField int, fieldValue string) {
	switch activeField {
	case 0:
		model.FormData.Name = fieldValue
	case 1:
		model.FormData.JSONConfig = fieldValue
	case 2:
		model.FormData.Environment = fieldValue
	}
}

func renderModalForFocusTest(model types.Model, modalType types.ModalType) string {
	switch modalType {
	case types.AddCommandForm:
		return renderCommandFormContent(model)
	case types.AddSSEForm:
		return renderSSEFormContent(model)
	case types.AddJSONForm:
		return renderJSONFormContent(model)
	case types.NoModal, types.AddModal, types.AddMCPTypeSelection, types.EditModal, types.DeleteModal:
		return ""
	default:
		return ""
	}
}

func assertModalFieldFocus(t *testing.T, result string, tt struct {
	modalType    types.ModalType
	activeField  int
	fieldValue   string
	expectCursor bool
}) {
	if tt.expectCursor {
		expectedCursor := tt.fieldValue + "_"
		if !strings.Contains(result, expectedCursor) {
			t.Errorf("Modal %v field %d should show cursor, expected '%s' in result",
				tt.modalType, tt.activeField, expectedCursor)
		}
	}
}

func TestOverlayModalUnknownType(t *testing.T) {
	model := types.NewModel()
	model.ActiveModal = types.ModalType(999) // Unknown modal type

	result := OverlayModal(model, 120, 40, "background")

	if !strings.Contains(result, "Unknown Modal") {
		t.Error("Should handle unknown modal type gracefully")
	}
	if !strings.Contains(result, "Unknown modal type") {
		t.Error("Should show unknown modal type message")
	}
}
