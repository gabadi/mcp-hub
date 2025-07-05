package ui

import (
	"strings"
	"testing"

	"cc-mcp-manager/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModelCreation(t *testing.T) {
	model := NewModel()

	if model.State != types.MainNavigation {
		t.Error("Expected initial state to be MainNavigation")
	}

	if len(model.MCPItems) == 0 {
		t.Error("Expected MCPItems to be populated")
	}

	if model.FormErrors == nil {
		t.Error("Expected FormErrors to be initialized")
	}
}

func TestModelInit(t *testing.T) {
	model := NewModel()
	cmd := model.Init()

	// Should return a command to refresh Claude status
	if cmd == nil {
		t.Error("Expected Init to return a command")
	}
}

func TestModelUpdateWindowSize(t *testing.T) {
	model := NewModel()

	// Test window size update
	windowMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	newModel, cmd := model.Update(windowMsg)

	updatedModel := newModel.(Model)
	if updatedModel.Width != 120 {
		t.Errorf("Expected width 120, got %d", updatedModel.Width)
	}
	if updatedModel.Height != 40 {
		t.Errorf("Expected height 40, got %d", updatedModel.Height)
	}

	// Command should be nil for window resize
	if cmd != nil {
		t.Error("Expected no command for window resize")
	}
}

func TestModelUpdateKeyPress(t *testing.T) {
	model := NewModel()

	// Test key press
	keyMsg := tea.KeyMsg{Type: tea.KeyRight}
	newModel, cmd := model.Update(keyMsg)

	// Should return the updated model
	if newModel == nil {
		t.Error("Expected model to be returned")
	}

	// Command may or may not be nil depending on key handling
	_ = cmd // Just verify it compiles
}

func TestModelUpdateSuccessMessage(t *testing.T) {
	model := NewModel()

	// Create a mock success message
	successMsg := struct {
		Message string
	}{
		Message: "Test success message",
	}

	// Since the handlers package may not be available in testing,
	// we'll test with a similar message structure
	_, cmd := model.Update(successMsg)

	// Should handle unknown message types gracefully
	if cmd != nil {
		t.Error("Expected nil command for unknown message type")
	}
}

func TestGetterMethods(t *testing.T) {
	model := NewModel()

	// Set some test values
	model.ColumnCount = 4
	model.ActiveColumn = 2
	model.SelectedItem = 5
	model.State = types.SearchActiveNavigation
	model.SearchQuery = "test query"
	model.SearchActive = true
	model.SearchInputActive = true

	// Test getter methods
	if model.GetColumnCount() != 4 {
		t.Errorf("Expected ColumnCount 4, got %d", model.GetColumnCount())
	}

	if model.GetActiveColumn() != 2 {
		t.Errorf("Expected ActiveColumn 2, got %d", model.GetActiveColumn())
	}

	if model.GetSelectedItem() != 5 {
		t.Errorf("Expected SelectedItem 5, got %d", model.GetSelectedItem())
	}

	if model.GetState() != types.SearchActiveNavigation {
		t.Errorf("Expected State SearchActiveNavigation, got %v", model.GetState())
	}

	if model.GetSearchQuery() != "test query" {
		t.Errorf("Expected SearchQuery 'test query', got %q", model.GetSearchQuery())
	}

	if !model.GetSearchActive() {
		t.Error("Expected SearchActive to be true")
	}

	if !model.GetSearchInputActive() {
		t.Error("Expected SearchInputActive to be true")
	}
}

func TestGetFilteredMCPs(t *testing.T) {
	model := NewModel()

	// Add some test MCPs
	model.MCPItems = []types.MCPItem{
		{Name: "github-mcp", Type: "CMD", Active: true},
		{Name: "docker-mcp", Type: "CMD", Active: false},
		{Name: "context7", Type: "SSE", Active: true},
	}

	// Test without search query
	filtered := model.GetFilteredMCPs()
	if len(filtered) != 3 {
		t.Errorf("Expected 3 filtered MCPs, got %d", len(filtered))
	}

	// Test with search query
	model.SearchQuery = "github"
	filtered = model.GetFilteredMCPs()
	if len(filtered) != 1 {
		t.Errorf("Expected 1 filtered MCP for 'github', got %d", len(filtered))
	}
	if filtered[0].Name != "github-mcp" {
		t.Errorf("Expected filtered MCP to be 'github-mcp', got %q", filtered[0].Name)
	}
}

func TestViewRendering(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40

	result := model.View()

	// Should return non-empty string
	if result == "" {
		t.Error("View should return non-empty string")
	}

	// Should contain some basic elements
	if !strings.Contains(result, "MCP Manager") {
		t.Error("View should contain 'MCP Manager'")
	}
}

func TestViewZeroDimensions(t *testing.T) {
	model := NewModel()
	model.Width = 0
	model.Height = 0

	result := model.View()

	if result != "Loading..." {
		t.Errorf("Expected 'Loading...' for zero dimensions, got: %s", result)
	}
}

func TestViewWithModal(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40
	model.State = types.ModalActive
	model.ActiveModal = types.AddCommandForm

	result := model.View()

	// Should contain modal content
	if !strings.Contains(result, "Add New MCP") || !strings.Contains(result, "Command") {
		t.Error("View should contain modal content when modal is active")
	}
}

func TestRenderHeaderDifferentStates(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40

	tests := []struct {
		state    types.AppState
		expected string
	}{
		{types.MainNavigation, "A=Add"},
		{types.SearchMode, "Type to search"},
		{types.SearchActiveNavigation, "Navigate Mode"},
		{types.ModalActive, "ESC=Cancel"},
	}

	for _, tt := range tests {
		model.State = tt.state
		if tt.state == types.SearchActiveNavigation {
			model.SearchActive = true
		}

		result := model.View()
		if !strings.Contains(result, tt.expected) {
			t.Errorf("Expected header to contain %q for state %v", tt.expected, tt.state)
		}
	}
}

func TestRenderBodyDifferentLayouts(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40

	// Test different column counts
	testCases := []struct {
		columns     int
		description string
	}{
		{1, "narrow layout"},
		{2, "medium layout"},
		{4, "wide layout"},
	}

	for _, tc := range testCases {
		model.ColumnCount = tc.columns
		result := model.View()

		if result == "" {
			t.Errorf("View should render content for %s", tc.description)
		}
	}
}

func TestRenderStatusColumn(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40
	model.ColumnCount = 2 // Ensure we're in a layout that shows status
	model.ClaudeStatus = types.ClaudeStatus{
		Available:  true,
		Version:    "1.0.0",
		ActiveMCPs: []string{"github-mcp", "context7"},
	}

	// Access the internal render method through View
	result := model.View()

	// Should contain status information - check for a more specific element
	if !strings.Contains(result, "Status") && !strings.Contains(result, "Claude") {
		t.Log("View output:", result)
		t.Error("Status column should be present in multi-column layout")
	}
}

func TestRenderDetailsColumn(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40
	model.ColumnCount = 2 // Multi-column layout
	model.MCPItems = []types.MCPItem{
		{Name: "test-mcp", Type: "CMD", Command: "test-cmd", Active: true},
	}
	model.SelectedItem = 0

	result := model.View()

	// Should render view successfully with MCP data
	if result == "" {
		t.Error("View should render content")
	}

	// In multi-column layout, MCPs should be displayed
	if !strings.Contains(result, "MCP") {
		t.Error("View should contain MCP-related content")
	}
}

func TestRenderFooterDifferentStates(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40

	// Test search active state
	model.State = types.SearchActiveNavigation
	model.SearchActive = true
	model.SearchQuery = "test"

	result := model.View()

	if !strings.Contains(result, "Search:") {
		t.Error("Footer should show search information when search is active")
	}
}

func TestGetLayoutName(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40

	// Test different column counts
	testCases := []struct {
		columns     int
		description string
	}{
		{1, "single column layout"},
		{2, "two column layout"},
		{4, "four column layout"},
		{3, "three column layout"},
	}

	for _, tc := range testCases {
		model.ColumnCount = tc.columns
		// Access layout through the view rendering
		result := model.View()

		// Should render successfully for different layouts
		if result == "" {
			t.Errorf("Should render %s", tc.description)
		}
	}
}

func TestModelWithDifferentMCPTypes(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40
	model.MCPItems = []types.MCPItem{
		{Name: "cmd-mcp", Type: "CMD", Command: "cmd", Active: true},
		{Name: "sse-mcp", Type: "SSE", URL: "http://example.com", Active: false},
		{Name: "json-mcp", Type: "JSON", JSONConfig: `{"key":"value"}`, Active: false},
	}

	result := model.View()

	// Should handle different MCP types successfully
	if result == "" {
		t.Error("Should render view with different MCP types")
	}

	// Should contain general MCP-related content
	if !strings.Contains(result, "MCP") {
		t.Error("View should contain MCP-related content")
	}
}

func TestModelSearchFunctionality(t *testing.T) {
	model := NewModel()
	model.MCPItems = []types.MCPItem{
		{Name: "github-mcp", Type: "CMD", Active: true},
		{Name: "docker-mcp", Type: "CMD", Active: false},
		{Name: "context7", Type: "SSE", Active: true},
	}

	// Test search filtering
	model.SearchQuery = "github"
	filtered := model.GetFilteredMCPs()

	if len(filtered) != 1 {
		t.Errorf("Expected 1 filtered result, got %d", len(filtered))
	}

	if filtered[0].Name != "github-mcp" {
		t.Errorf("Expected 'github-mcp', got %q", filtered[0].Name)
	}

	// Test case-insensitive search
	model.SearchQuery = "GITHUB"
	filtered = model.GetFilteredMCPs()

	if len(filtered) != 1 {
		t.Error("Search should be case-insensitive")
	}

	// Test partial match
	model.SearchQuery = "git"
	filtered = model.GetFilteredMCPs()

	if len(filtered) != 1 {
		t.Error("Search should support partial matches")
	}
}

func TestModelWithEmptyMCPList(t *testing.T) {
	model := NewModel()
	model.MCPItems = []types.MCPItem{}

	result := model.View()

	// Should handle empty MCP list gracefully
	if result == "" {
		t.Error("View should handle empty MCP list")
	}
}

func TestModelSuccessMessage(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40
	model.SuccessMessage = "Test success message"
	model.SuccessTimer = 120

	result := model.View()

	// Should render view with success message state
	if result == "" {
		t.Error("View should render when success message is set")
	}
}

func TestModelUpdateUnknownMessage(t *testing.T) {
	model := NewModel()

	// Test with unknown message type
	unknownMsg := struct {
		UnknownField string
	}{
		UnknownField: "unknown",
	}

	newModel, cmd := model.Update(unknownMsg)

	// Should return the same model and nil command
	if newModel == nil {
		t.Error("Should return model for unknown message")
	}
	if cmd != nil {
		t.Error("Should return nil command for unknown message")
	}
}

func TestModelLargeWindow(t *testing.T) {
	model := NewModel()
	model.Width = 200
	model.Height = 60

	result := model.View()

	// Should handle large windows
	if result == "" {
		t.Error("Should handle large window dimensions")
	}
}

func TestModelSmallWindow(t *testing.T) {
	model := NewModel()
	model.Width = 40
	model.Height = 10

	result := model.View()

	// Should handle small windows
	if result == "" {
		t.Error("Should handle small window dimensions")
	}
}
