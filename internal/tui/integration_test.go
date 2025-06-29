package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TestIntegrationFullFlow tests a complete user interaction flow
func TestIntegrationFullFlow(t *testing.T) {
	m := NewModel()
	
	// Simulate window size initialization
	model, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = model.(Model)
	
	if !m.initialized {
		t.Error("Model should be initialized after window size event")
	}
	
	// Test navigation sequence
	steps := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'j'}}, // down
		{Type: tea.KeyRunes, Runes: []rune{'j'}}, // down
		{Type: tea.KeyRunes, Runes: []rune{'l'}}, // right column
		{Type: tea.KeyRunes, Runes: []rune{'h'}}, // left column
		{Type: tea.KeyTab},                       // enter search
	}
	
	for i, step := range steps {
		model, _ = m.Update(step)
		m = model.(Model)
		
		// Validate state after each step
		switch i {
		case 0:
			if m.selectedIndex != 1 {
				t.Errorf("Step %d: expected selectedIndex 1, got %d", i, m.selectedIndex)
			}
		case 1:
			if m.selectedIndex != 2 {
				t.Errorf("Step %d: expected selectedIndex 2, got %d", i, m.selectedIndex)
			}
		case 2:
			if m.currentColumn != ColumnCenter {
				t.Errorf("Step %d: expected ColumnCenter, got %v", i, m.currentColumn)
			}
		case 3:
			if m.currentColumn != ColumnLeft {
				t.Errorf("Step %d: expected ColumnLeft, got %v", i, m.currentColumn)
			}
		case 4:
			if !m.searchActive {
				t.Errorf("Step %d: expected search to be active", i)
			}
		}
	}
	
	// Test search functionality
	searchSteps := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'s'}},
		{Type: tea.KeyRunes, Runes: []rune{'e'}},
		{Type: tea.KeyRunes, Runes: []rune{'r'}},
		{Type: tea.KeyRunes, Runes: []rune{'v'}},
		{Type: tea.KeyBackspace},
		{Type: tea.KeyEnter},
	}
	
	for _, step := range searchSteps {
		model, _ = m.Update(step)
		m = model.(Model)
	}
	
	if m.searchQuery != "ser" {
		t.Errorf("Expected search query 'ser', got '%s'", m.searchQuery)
	}
	
	if m.searchActive {
		t.Error("Search should be inactive after Enter")
	}
}

// TestPerformanceStartup tests startup performance requirements
func TestPerformanceStartup(t *testing.T) {
	start := time.Now()
	
	m := NewModel()
	model, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = model.(Model)
	
	// Render initial view
	_ = m.View()
	
	elapsed := time.Since(start)
	
	// Should be well under 500ms requirement
	if elapsed > 100*time.Millisecond {
		t.Errorf("Startup took %v, expected under 100ms", elapsed)
	}
}

// TestResponsiveLayoutTransitions tests layout changes
func TestResponsiveLayoutTransitions(t *testing.T) {
	m := NewModel()
	
	// Test all layout transitions
	testCases := []struct {
		width    int
		expected LayoutType
		desc     string
	}{
		{120, Layout3Column, "wide screen"},
		{80, Layout3Column, "3-column minimum"},
		{79, Layout2Column, "2-column threshold"},
		{60, Layout2Column, "2-column minimum"},
		{59, Layout1Column, "single column threshold"},
		{40, Layout1Column, "narrow screen"},
	}
	
	for _, tc := range testCases {
		model, _ := m.Update(tea.WindowSizeMsg{Width: tc.width, Height: 30})
		m = model.(Model)
		
		if m.layout != tc.expected {
			t.Errorf("%s (width %d): expected layout %v, got %v", 
				tc.desc, tc.width, tc.expected, m.layout)
		}
		
		// Test that view renders without panic
		view := m.View()
		if len(view) == 0 {
			t.Errorf("%s: view should not be empty", tc.desc)
		}
	}
}

// TestKeyboardShortcuts tests all documented keyboard shortcuts
func TestKeyboardShortcuts(t *testing.T) {
	m := NewModel()
	model, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = model.(Model)
	
	// Test navigation keys
	navKeys := []struct {
		key   string
		runes []rune
		desc  string
	}{
		{"j", []rune{'j'}, "down"},
		{"k", []rune{'k'}, "up"},
		{"h", []rune{'h'}, "left"},
		{"l", []rune{'l'}, "right"},
		{"down", nil, "arrow down"},
		{"up", nil, "arrow up"},
		{"left", nil, "arrow left"},
		{"right", nil, "arrow right"},
	}
	
	for _, nk := range navKeys {
		var keyMsg tea.KeyMsg
		if nk.runes != nil {
			keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: nk.runes}
		} else {
			switch nk.key {
			case "down":
				keyMsg = tea.KeyMsg{Type: tea.KeyDown}
			case "up":
				keyMsg = tea.KeyMsg{Type: tea.KeyUp}
			case "left":
				keyMsg = tea.KeyMsg{Type: tea.KeyLeft}
			case "right":
				keyMsg = tea.KeyMsg{Type: tea.KeyRight}
			}
		}
		
		// Should not panic
		_, _ = m.Update(keyMsg)
	}
	
	// Test special keys
	specialKeys := []tea.KeyMsg{
		{Type: tea.KeyTab},
		{Type: tea.KeyEnter},
		{Type: tea.KeyBackspace},
		{Type: tea.KeyEsc},
		{Type: tea.KeyRunes, Runes: []rune{'q'}},
	}
	
	for _, key := range specialKeys {
		// Should not panic
		_, _ = m.Update(key)
	}
}

// TestTerminalCompatibility tests basic terminal feature requirements
func TestTerminalCompatibility(t *testing.T) {
	m := NewModel()
	
	// Test various terminal sizes that might be encountered
	terminalSizes := []tea.WindowSizeMsg{
		{Width: 200, Height: 50}, // Large terminal
		{Width: 80, Height: 24},  // Standard terminal
		{Width: 40, Height: 12},  // Small terminal
		{Width: 20, Height: 5},   // Very small terminal
	}
	
	for _, size := range terminalSizes {
		model, _ := m.Update(size)
		m = model.(Model)
		
		// Should not panic and should render
		view := m.View()
		if len(view) == 0 {
			t.Errorf("Terminal size %dx%d: view should not be empty", size.Width, size.Height)
		}
		
		// Should have reasonable layout
		if m.width != size.Width || m.height != size.Height {
			t.Errorf("Terminal size not properly set: expected %dx%d, got %dx%d",
				size.Width, size.Height, m.width, m.height)
		}
	}
}