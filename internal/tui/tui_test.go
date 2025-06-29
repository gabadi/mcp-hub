package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNewModel tests the creation of a new TUI model
func TestNewModel(t *testing.T) {
	m := NewModel()
	
	if m.currentColumn != ColumnLeft {
		t.Errorf("Expected currentColumn to be ColumnLeft, got %v", m.currentColumn)
	}
	
	if m.currentMode != ViewModeList {
		t.Errorf("Expected currentMode to be ViewModeList, got %v", m.currentMode)
	}
	
	if m.selectedIndex != 0 {
		t.Errorf("Expected selectedIndex to be 0, got %d", m.selectedIndex)
	}
	
	if len(m.mcpList) == 0 {
		t.Error("Expected mcpList to have mock data")
	}
}

// TestLayoutDetermination tests the layout calculation logic
func TestLayoutDetermination(t *testing.T) {
	testCases := []struct {
		width    int
		expected LayoutType
	}{
		{120, Layout3Column},
		{80, Layout3Column},
		{79, Layout2Column},
		{60, Layout2Column},
		{59, Layout1Column},
		{40, Layout1Column},
	}
	
	for _, tc := range testCases {
		m := NewModel()
		m.width = tc.width
		m.determineLayout()
		
		if m.layout != tc.expected {
			t.Errorf("Width %d: expected layout %v, got %v", tc.width, tc.expected, m.layout)
		}
	}
}

// TestNavigationKeys tests keyboard navigation
func TestNavigationKeys(t *testing.T) {
	m := NewModel()
	m.width = 100
	m.height = 30
	m.determineLayout()
	m.initialized = true
	
	// Test down navigation
	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = model.(Model)
	if m.selectedIndex != 1 {
		t.Errorf("Expected selectedIndex to be 1 after 'j', got %d", m.selectedIndex)
	}
	
	// Test up navigation
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = model.(Model)
	if m.selectedIndex != 0 {
		t.Errorf("Expected selectedIndex to be 0 after 'k', got %d", m.selectedIndex)
	}
	
	// Test wrap around down
	for i := 0; i < len(m.mcpList); i++ {
		model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		m = model.(Model)
	}
	if m.selectedIndex != 0 {
		t.Errorf("Expected selectedIndex to wrap to 0, got %d", m.selectedIndex)
	}
	
	// Test wrap around up
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = model.(Model)
	expected := len(m.mcpList) - 1
	if m.selectedIndex != expected {
		t.Errorf("Expected selectedIndex to wrap to %d, got %d", expected, m.selectedIndex)
	}
}

// TestColumnNavigation tests navigation between columns
func TestColumnNavigation(t *testing.T) {
	m := NewModel()
	m.width = 100 // 3-column layout
	m.height = 30
	m.determineLayout()
	m.initialized = true
	
	// Test right navigation in 3-column layout
	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = model.(Model)
	if m.currentColumn != ColumnCenter {
		t.Errorf("Expected currentColumn to be ColumnCenter, got %v", m.currentColumn)
	}
	
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = model.(Model)
	if m.currentColumn != ColumnRight {
		t.Errorf("Expected currentColumn to be ColumnRight, got %v", m.currentColumn)
	}
	
	// Test left navigation
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = model.(Model)
	if m.currentColumn != ColumnCenter {
		t.Errorf("Expected currentColumn to be ColumnCenter, got %v", m.currentColumn)
	}
	
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = model.(Model)
	if m.currentColumn != ColumnLeft {
		t.Errorf("Expected currentColumn to be ColumnLeft, got %v", m.currentColumn)
	}
}

// TestSearchMode tests search functionality
func TestSearchMode(t *testing.T) {
	m := NewModel()
	m.width = 100
	m.height = 30
	m.determineLayout()
	m.initialized = true
	
	// Enter search mode
	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = model.(Model)
	if !m.searchActive {
		t.Error("Expected searchActive to be true after Tab")
	}
	
	if m.currentMode != ViewModeSearch {
		t.Errorf("Expected currentMode to be ViewModeSearch, got %v", m.currentMode)
	}
	
	// Type in search mode
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	m = model.(Model)
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	m = model.(Model)
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	m = model.(Model)
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	m = model.(Model)
	
	if m.searchQuery != "test" {
		t.Errorf("Expected searchQuery to be 'test', got '%s'", m.searchQuery)
	}
	
	// Test backspace
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	m = model.(Model)
	if m.searchQuery != "tes" {
		t.Errorf("Expected searchQuery to be 'tes' after backspace, got '%s'", m.searchQuery)
	}
	
	// Exit search mode with Enter
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = model.(Model)
	if m.searchActive {
		t.Error("Expected searchActive to be false after Enter")
	}
	
	if m.currentMode != ViewModeList {
		t.Errorf("Expected currentMode to be ViewModeList, got %v", m.currentMode)
	}
}

// TestExitKeys tests application exit functionality
func TestExitKeys(t *testing.T) {
	m := NewModel()
	m.width = 100
	m.height = 30
	m.determineLayout()
	m.initialized = true
	
	// Test Q key
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Error("Expected quit command after 'q' key")
	}
	
	// Test ESC key from main interface
	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Error("Expected quit command after ESC key")
	}
}

// TestEscapeFromSearch tests ESC behavior in search mode
func TestEscapeFromSearch(t *testing.T) {
	m := NewModel()
	m.width = 100
	m.height = 30
	m.determineLayout()
	m.initialized = true
	
	// Enter search mode
	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = model.(Model)
	
	// ESC should exit search mode, not application
	model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = model.(Model)
	if m.searchActive {
		t.Error("Expected searchActive to be false after ESC in search mode")
	}
	
	if cmd != nil {
		t.Error("Expected ESC in search mode to NOT quit application")
	}
}

// TestWindowSizeHandling tests terminal size change handling
func TestWindowSizeHandling(t *testing.T) {
	m := NewModel()
	
	// Simulate window size change
	model, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = model.(Model)
	
	if m.width != 120 {
		t.Errorf("Expected width to be 120, got %d", m.width)
	}
	
	if m.height != 40 {
		t.Errorf("Expected height to be 40, got %d", m.height)
	}
	
	if m.layout != Layout3Column {
		t.Errorf("Expected layout to be Layout3Column for width 120, got %v", m.layout)
	}
	
	if !m.initialized {
		t.Error("Expected model to be initialized after window size message")
	}
}

// TestView tests basic view rendering
func TestView(t *testing.T) {
	m := NewModel()
	
	// Test uninitialized view
	view := m.View()
	if view != "Initializing MCP Manager..." {
		t.Errorf("Expected initialization message, got: %s", view)
	}
	
	// Initialize and test
	m.width = 100
	m.height = 30
	m.determineLayout()
	m.initialized = true
	
	view = m.View()
	if len(view) == 0 {
		t.Error("Expected non-empty view after initialization")
	}
	
	// View should contain the title
	if !contains(view, "MCP Manager CLI") {
		t.Error("Expected view to contain 'MCP Manager CLI'")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		   len(s) > len(substr) && contains(s[1:], substr)
}

// Benchmark navigation performance
func BenchmarkNavigation(b *testing.B) {
	m := NewModel()
	m.width = 100
	m.height = 30
	m.determineLayout()
	m.initialized = true
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		m = model.(Model)
	}
}