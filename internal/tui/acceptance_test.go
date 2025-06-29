package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestAC1_ApplicationLaunchWithBubbleTeaTUI tests AC1: Application Launch with Bubble Tea TUI
func TestAC1_ApplicationLaunchWithBubbleTeaTUI(t *testing.T) {
	// Given the MCP Manager CLI is executed
	m := NewModel()
	
	// When the application starts
	model, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = model.(Model)
	
	// Then it displays a Bubble Tea TUI interface with a 3-column layout
	if m.layout != Layout3Column {
		t.Error("Expected 3-column layout for 100-width terminal")
	}
	
	// And the interface renders correctly in terminals 80+ columns wide
	model, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = model.(Model)
	if m.layout != Layout3Column {
		t.Error("Expected 3-column layout for 80-width terminal")
	}
	
	// And the header displays application title and keyboard shortcuts
	view := m.View()
	if !strings.Contains(view, "MCP Manager CLI") {
		t.Error("Header should contain application title 'MCP Manager CLI'")
	}
	
	expectedShortcuts := []string{"[A]dd", "[E]dit", "[D]elete", "[Space]Toggle", "[R]efresh", "[Q]uit"}
	for _, shortcut := range expectedShortcuts {
		if !strings.Contains(view, shortcut) {
			t.Errorf("Header should contain keyboard shortcut: %s", shortcut)
		}
	}
}

// TestAC2_ArrowKeyNavigation tests AC2: Arrow Key Navigation
func TestAC2_ArrowKeyNavigation(t *testing.T) {
	// Given the TUI is displayed with MCP inventory
	m := NewModel()
	model, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = model.(Model)
	
	initialIndex := m.selectedIndex
	
	// When the user presses arrow keys (↑↓←→)
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = model.(Model)
	
	// Then the selection cursor moves within and between columns
	if m.selectedIndex == initialIndex {
		t.Error("Down arrow should move selection")
	}
	
	// Test up arrow
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = model.(Model)
	if m.selectedIndex != initialIndex {
		t.Error("Up arrow should return to initial position")
	}
	
	// Test right arrow (column navigation)
	initialColumn := m.currentColumn
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = model.(Model)
	
	if m.currentColumn == initialColumn {
		t.Error("Right arrow should change column in 3-column layout")
	}
	
	// And the currently selected MCP is visually highlighted
	view := m.View()
	if !strings.Contains(view, ">") {
		t.Error("View should contain selection indicator '>'")
	}
	
	// And navigation wraps appropriately at column boundaries
	// Test wrapping by going to last item and pressing down
	for i := 0; i < len(m.mcpList); i++ {
		model, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = model.(Model)
	}
	if m.selectedIndex != 0 {
		t.Error("Navigation should wrap to beginning when going past end")
	}
}

// TestAC3_SearchFieldNavigation tests AC3: Search Field Navigation
func TestAC3_SearchFieldNavigation(t *testing.T) {
	// Given the TUI is displayed
	m := NewModel()
	model, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = model.(Model)
	
	// When the user presses the Tab key
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = model.(Model)
	
	// Then focus jumps to the search field
	if !m.searchActive {
		t.Error("Tab key should activate search mode")
	}
	
	// And the search field is visually highlighted
	view := m.View()
	if !strings.Contains(view, "[Enter]Finish [Esc]Cancel") {
		t.Error("Search mode should show different keyboard shortcuts")
	}
	
	// And the user can type search terms
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	m = model.(Model)
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	m = model.(Model)
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	m = model.(Model)
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	m = model.(Model)
	
	if m.searchQuery != "test" {
		t.Errorf("Expected search query 'test', got '%s'", m.searchQuery)
	}
	
	// Search query should be visible in footer
	view = m.View()
	if !strings.Contains(view, "Search: test") {
		t.Error("Search query should be visible in status footer")
	}
}

// TestAC4_ApplicationExit tests AC4: Application Exit
func TestAC4_ApplicationExit(t *testing.T) {
	// Given the TUI is running
	m := NewModel()
	model, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = model.(Model)
	
	// When the user presses the ESC key from the main interface
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	
	// Then the application exits cleanly
	if cmd == nil {
		t.Error("ESC key should generate quit command")
	}
	
	// Or When the user presses 'Q' key
	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	
	// Then the application exits with confirmation
	if cmd == nil {
		t.Error("Q key should generate quit command")
	}
}

// TestAC5_ResponsiveLayoutAdaptation tests AC5: Responsive Layout Adaptation
func TestAC5_ResponsiveLayoutAdaptation(t *testing.T) {
	m := NewModel()
	
	// Given the application is running in different terminal widths
	// When the terminal is 80-119 columns wide
	model, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = model.(Model)
	
	// Then the interface displays in 3-column layout
	if m.layout != Layout3Column {
		t.Error("Expected 3-column layout for width 100")
	}
	
	// Test 80 columns (minimum for 3-column)
	model, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 30})
	m = model.(Model)
	if m.layout != Layout3Column {
		t.Error("Expected 3-column layout for width 80")
	}
	
	// When the terminal is 60-79 columns wide
	model, _ = m.Update(tea.WindowSizeMsg{Width: 70, Height: 30})
	m = model.(Model)
	
	// Then the interface adapts to 2-column layout
	if m.layout != Layout2Column {
		t.Error("Expected 2-column layout for width 70")
	}
	
	// Test 60 columns (minimum for 2-column)
	model, _ = m.Update(tea.WindowSizeMsg{Width: 60, Height: 30})
	m = model.(Model)
	if m.layout != Layout2Column {
		t.Error("Expected 2-column layout for width 60")
	}
	
	// When the terminal is less than 60 columns wide
	model, _ = m.Update(tea.WindowSizeMsg{Width: 50, Height: 30})
	m = model.(Model)
	
	// Then the interface displays in single-column layout
	if m.layout != Layout1Column {
		t.Error("Expected single-column layout for width 50")
	}
	
	// Test that views render correctly in all layouts
	layouts := []struct {
		width  int
		layout LayoutType
		name   string
	}{
		{100, Layout3Column, "3-column"},
		{70, Layout2Column, "2-column"},
		{50, Layout1Column, "1-column"},
	}
	
	for _, l := range layouts {
		model, _ = m.Update(tea.WindowSizeMsg{Width: l.width, Height: 30})
		m = model.(Model)
		
		view := m.View()
		if len(view) == 0 {
			t.Errorf("%s layout should render non-empty view", l.name)
		}
		
		if m.layout != l.layout {
			t.Errorf("Width %d should result in %v layout, got %v", l.width, l.layout, m.layout)
		}
	}
}

// TestAC6_KeyboardShortcutsDisplay tests AC6: Keyboard Shortcuts Display
func TestAC6_KeyboardShortcutsDisplay(t *testing.T) {
	// Given the TUI is displayed
	m := NewModel()
	model, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = model.(Model)
	
	// When the interface renders
	view := m.View()
	
	// Then the header shows current keyboard shortcuts
	expectedShortcuts := []string{
		"[A]dd", "[E]dit", "[D]elete", "[Space]Toggle", "[R]efresh", "[Q]uit",
	}
	
	for _, shortcut := range expectedShortcuts {
		if !strings.Contains(view, shortcut) {
			t.Errorf("Header should display keyboard shortcut: %s", shortcut)
		}
	}
	
	// And shortcuts are context-appropriate for current mode
	// Test search mode shortcuts
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = model.(Model)
	view = m.View()
	
	searchShortcuts := []string{"[Enter]Finish", "[Esc]Cancel"}
	for _, shortcut := range searchShortcuts {
		if !strings.Contains(view, shortcut) {
			t.Errorf("Search mode should display shortcut: %s", shortcut)
		}
	}
	
	// And current context is displayed in the status bar
	if !strings.Contains(view, "Search:") {
		t.Error("Status bar should show search context when in search mode")
	}
	
	// Test layout information in status bar
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Exit search
	m = model.(Model)
	view = m.View()
	
	if !strings.Contains(view, "3-Column Mode") {
		t.Error("Status bar should show current layout mode")
	}
	
	if !strings.Contains(view, "100x30") {
		t.Error("Status bar should show terminal dimensions")
	}
}

// TestAllAcceptanceCriteria runs a comprehensive test of all acceptance criteria
func TestAllAcceptanceCriteria(t *testing.T) {
	t.Run("AC1: Application Launch", TestAC1_ApplicationLaunchWithBubbleTeaTUI)
	t.Run("AC2: Arrow Key Navigation", TestAC2_ArrowKeyNavigation)
	t.Run("AC3: Search Field Navigation", TestAC3_SearchFieldNavigation)
	t.Run("AC4: Application Exit", TestAC4_ApplicationExit)
	t.Run("AC5: Responsive Layout", TestAC5_ResponsiveLayoutAdaptation)
	t.Run("AC6: Keyboard Shortcuts", TestAC6_KeyboardShortcutsDisplay)
	
	t.Log("✅ All acceptance criteria validated successfully")
}