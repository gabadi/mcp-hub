package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInitialModel(t *testing.T) {
	model := initialModel()

	// Test initial state
	if model.activeColumn != 0 {
		t.Errorf("Expected activeColumn to be 0, got %d", model.activeColumn)
	}

	if model.searchFocused {
		t.Error("Expected searchFocused to be false initially")
	}

	if model.columns != 3 {
		t.Errorf("Expected 3 columns initially, got %d", model.columns)
	}

	if model.selectedMCP != 0 {
		t.Errorf("Expected selectedMCP to be 0, got %d", model.selectedMCP)
	}

	if len(model.mcps) == 0 {
		t.Error("Expected initial MCPs to be populated")
	}
}

func TestResponsiveLayout(t *testing.T) {
	model := initialModel()

	// Test narrow width (1 column)
	narrowMsg := tea.WindowSizeMsg{Width: 50, Height: 24}
	updatedModel, _ := model.Update(narrowMsg)
	m := updatedModel.(Model)

	if m.columns != 1 {
		t.Errorf("Expected 1 column for width 50, got %d", m.columns)
	}

	// Test medium width (2 columns)
	mediumMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updatedModel, _ = model.Update(mediumMsg)
	m = updatedModel.(Model)

	if m.columns != 2 {
		t.Errorf("Expected 2 columns for width 80, got %d", m.columns)
	}

	// Test wide width (3 columns)
	wideMsg := tea.WindowSizeMsg{Width: 150, Height: 24}
	updatedModel, _ = model.Update(wideMsg)
	m = updatedModel.(Model)

	if m.columns != 3 {
		t.Errorf("Expected 3 columns for width 150, got %d", m.columns)
	}
}

func TestKeyboardNavigation(t *testing.T) {
	model := initialModel()
	model.width = 120 // Ensure we have 3 columns
	model.height = 24

	tests := []struct {
		name     string
		key      string
		expected func(Model) bool
	}{
		{
			name: "down arrow increases selected MCP",
			key:  "down",
			expected: func(m Model) bool {
				return m.selectedMCP == 1
			},
		},
		{
			name: "up arrow at position 0 stays at 0",
			key:  "up",
			expected: func(m Model) bool {
				return m.selectedMCP == 0
			},
		},
		{
			name: "right arrow increases active column",
			key:  "right",
			expected: func(m Model) bool {
				return m.activeColumn == 1
			},
		},
		{
			name: "left arrow at column 0 stays at 0",
			key:  "left",
			expected: func(m Model) bool {
				return m.activeColumn == 0
			},
		},
		{
			name: "tab key focuses search",
			key:  "tab",
			expected: func(m Model) bool {
				return m.searchFocused
			},
		},
		{
			name: "slash key focuses search",
			key:  "/",
			expected: func(m Model) bool {
				return m.searchFocused
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "down" {
				keyMsg = tea.KeyMsg{Type: tea.KeyDown}
			} else if tt.key == "up" {
				keyMsg = tea.KeyMsg{Type: tea.KeyUp}
			} else if tt.key == "right" {
				keyMsg = tea.KeyMsg{Type: tea.KeyRight}
			} else if tt.key == "left" {
				keyMsg = tea.KeyMsg{Type: tea.KeyLeft}
			} else if tt.key == "tab" {
				keyMsg = tea.KeyMsg{Type: tea.KeyTab}
			}

			updatedModel, _ := model.Update(keyMsg)
			m := updatedModel.(Model)

			if !tt.expected(m) {
				t.Errorf("Test %s failed", tt.name)
			}
		})
	}
}

func TestSearchFunctionality(t *testing.T) {
	model := initialModel()

	// Focus search
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(tabMsg)
	m := updatedModel.(Model)

	if !m.searchFocused {
		t.Error("Expected search to be focused after tab")
	}

	// Type in search
	searchChar := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("context")}
	updatedModel, _ = m.Update(searchChar)
	m = updatedModel.(Model)

	if m.searchQuery != "context" {
		t.Errorf("Expected search query 'context', got '%s'", m.searchQuery)
	}

	// Test backspace
	backspaceMsg := tea.KeyMsg{Type: tea.KeyBackspace}
	updatedModel, _ = m.Update(backspaceMsg)
	m = updatedModel.(Model)

	if m.searchQuery != "contex" {
		t.Errorf("Expected search query 'contex' after backspace, got '%s'", m.searchQuery)
	}

	// Test ESC clears search
	escMsg := tea.KeyMsg{Type: tea.KeyEscape}
	updatedModel, _ = m.Update(escMsg)
	m = updatedModel.(Model)

	if m.searchFocused {
		t.Error("Expected search focus to be cleared after ESC")
	}

	if m.searchQuery != "" {
		t.Errorf("Expected empty search query after ESC, got '%s'", m.searchQuery)
	}
}

func TestMCPToggle(t *testing.T) {
	model := initialModel()
	originalStatus := model.mcps[0].Active

	// Press space to toggle
	spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
	updatedModel, _ := model.Update(spaceMsg)
	m := updatedModel.(Model)

	if m.mcps[0].Active == originalStatus {
		t.Error("Expected MCP active status to toggle")
	}
}

func TestStringUtilities(t *testing.T) {
	// Test toLower function
	tests := []struct {
		input    string
		expected string
	}{
		{"HELLO", "hello"},
		{"Hello", "hello"},
		{"hello", "hello"},
		{"MixedCase", "mixedcase"},
		{"", ""},
	}

	for _, tt := range tests {
		result := toLower(tt.input)
		if result != tt.expected {
			t.Errorf("toLower(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}

	// Test contains function
	containsTests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"hello world", "hello", true},
		{"hello world", "HELLO", true},
		{"hello world", "world", true},
		{"hello world", "foo", false},
		{"", "", true},
		{"hello", "", true},
	}

	for _, tt := range containsTests {
		result := contains(tt.s, tt.substr)
		if result != tt.expected {
			t.Errorf("contains(%s, %s) = %t, expected %t", tt.s, tt.substr, result, tt.expected)
		}
	}
}

func TestFilteredMCPs(t *testing.T) {
	model := initialModel()

	// Test with empty search
	filtered := model.getFilteredMCPs()
	if len(filtered) != len(model.mcps) {
		t.Errorf("Expected %d MCPs with empty search, got %d", len(model.mcps), len(filtered))
	}

	// Test with search query that should match
	model.searchQuery = "context"
	filtered = model.getFilteredMCPs()

	foundMatch := false
	for _, mcp := range filtered {
		if contains(mcp.Name, "context") || contains(mcp.Description, "context") {
			foundMatch = true
			break
		}
	}

	if !foundMatch {
		t.Error("Expected to find MCP matching 'context' search")
	}

	// Test with search query that should not match anything
	model.searchQuery = "nonexistent"
	filtered = model.getFilteredMCPs()

	if len(filtered) != 0 {
		t.Errorf("Expected 0 MCPs for 'nonexistent' search, got %d", len(filtered))
	}
}

func TestViewRendering(t *testing.T) {
	model := initialModel()
	model.width = 120
	model.height = 24

	// Test that view doesn't panic and returns content
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Test loading state
	emptyModel := Model{}
	loadingView := emptyModel.View()
	if loadingView != "Loading..." {
		t.Errorf("Expected 'Loading...' for empty model, got '%s'", loadingView)
	}
}

func TestInit(t *testing.T) {
	model := initialModel()
	cmd := model.Init()

	// Init should return nil command
	if cmd != nil {
		t.Error("Expected Init() to return nil command")
	}
}

func TestRenderComponents(t *testing.T) {
	model := initialModel()
	model.width = 120
	model.height = 24

	// Test header rendering
	header := model.renderHeader()
	if header == "" {
		t.Error("Expected non-empty header")
	}

	// Test search bar rendering
	searchBar := model.renderSearchBar()
	if searchBar == "" {
		t.Error("Expected non-empty search bar")
	}

	// Test MCP list rendering
	mcpList := model.renderMCPList(20)
	if mcpList == "" {
		t.Error("Expected non-empty MCP list")
	}

	// Test details column rendering
	details := model.renderDetailsColumn(20)
	if details == "" {
		t.Error("Expected non-empty details column")
	}

	// Test actions column rendering
	actions := model.renderActionsColumn(20)
	if actions == "" {
		t.Error("Expected non-empty actions column")
	}

	// Test footer rendering
	footer := model.renderFooter()
	if footer == "" {
		t.Error("Expected non-empty footer")
	}
}

func TestRenderContentDifferentColumns(t *testing.T) {
	model := initialModel()
	model.width = 120
	model.height = 24

	// Test 1 column layout
	model.columns = 1
	content1 := model.renderContent()
	if content1 == "" {
		t.Error("Expected non-empty content for 1 column")
	}

	// Test 2 column layout
	model.columns = 2
	content2 := model.renderContent()
	if content2 == "" {
		t.Error("Expected non-empty content for 2 columns")
	}

	// Test 3 column layout
	model.columns = 3
	content3 := model.renderContent()
	if content3 == "" {
		t.Error("Expected non-empty content for 3 columns")
	}

	// Test invalid column configuration
	model.columns = 0
	contentInvalid := model.renderContent()
	if contentInvalid != "Invalid column configuration" {
		t.Error("Expected 'Invalid column configuration' for invalid columns")
	}
}

func TestActiveColumnHighlight(t *testing.T) {
	model := initialModel()
	model.width = 120
	model.height = 24

	// Test different active columns affect rendering
	for i := 0; i < 3; i++ {
		model.activeColumn = i
		content := model.renderContent()
		if content == "" {
			t.Errorf("Expected non-empty content for active column %d", i)
		}
	}
}

func TestSearchHighlight(t *testing.T) {
	model := initialModel()
	model.width = 120
	model.height = 24

	// Test search focused state
	model.searchFocused = true
	searchBar := model.renderSearchBar()
	if searchBar == "" {
		t.Error("Expected non-empty search bar when focused")
	}

	// Test search with query
	model.searchQuery = "test"
	searchBarWithQuery := model.renderSearchBar()
	if searchBarWithQuery == "" {
		t.Error("Expected non-empty search bar with query")
	}
}

func TestEmptyMCPList(t *testing.T) {
	model := initialModel()
	model.mcps = []MCP{} // Empty list
	model.width = 120
	model.height = 24

	mcpList := model.renderMCPList(20)
	if mcpList == "" {
		t.Error("Expected non-empty MCP list even when empty")
	}
}

func TestDetailsPanelWithoutSelection(t *testing.T) {
	model := initialModel()
	model.selectedMCP = 999 // Invalid selection
	model.width = 120
	model.height = 24

	details := model.renderDetailsColumn(20)
	if details == "" {
		t.Error("Expected non-empty details column even with invalid selection")
	}
}
