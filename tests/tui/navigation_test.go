package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Mock model for integration testing
type testModel struct {
	width         int
	height        int
	activeColumn  int
	searchFocused bool
	searchQuery   string
	columns       int
	mcps          []testMCP
	selectedMCP   int
}

type testMCP struct {
	Name        string
	Type        string
	Description string
	Active      bool
}

func newTestModel() testModel {
	return testModel{
		width:        120,
		height:       24,
		activeColumn: 0,
		columns:      3,
		mcps: []testMCP{
			{Name: "context7", Type: "SSE", Description: "Library lookup", Active: false},
			{Name: "ht-mcp", Type: "CMD", Description: "Testing tool", Active: true},
			{Name: "filesystem", Type: "CMD", Description: "File operations", Active: false},
		},
		selectedMCP: 0,
	}
}

func (m testModel) Init() tea.Cmd { return nil }

func (m testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if msg.Width < 60 {
			m.columns = 1
		} else if msg.Width < 120 {
			m.columns = 2
		} else {
			m.columns = 3
		}
	case tea.KeyMsg:
		// Handle search input first
		if m.searchFocused {
			switch msg.Type {
			case tea.KeyEnter:
				m.searchFocused = false
				m.activeColumn = 0
			case tea.KeyBackspace:
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}
			case tea.KeyEscape:
				m.searchFocused = false
				m.searchQuery = ""
				m.activeColumn = 0
			case tea.KeyRunes:
				m.searchQuery += string(msg.Runes)
			}
			return m, nil
		}

		// Global key handlers when not in search mode
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEscape:
			return m, tea.Quit
		case tea.KeyTab:
			m.searchFocused = true
			m.activeColumn = -1
		case tea.KeyUp:
			if m.selectedMCP > 0 {
				m.selectedMCP--
			}
		case tea.KeyDown:
			if m.selectedMCP < len(m.mcps)-1 {
				m.selectedMCP++
			}
		case tea.KeyLeft:
			if m.activeColumn > 0 {
				m.activeColumn--
			}
		case tea.KeyRight:
			if m.activeColumn < m.columns-1 {
				m.activeColumn++
			}
		case tea.KeySpace:
			if m.selectedMCP < len(m.mcps) {
				m.mcps[m.selectedMCP].Active = !m.mcps[m.selectedMCP].Active
			}
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				return m, tea.Quit
			case "/":
				m.searchFocused = true
				m.activeColumn = -1
			}
		}
	}
	return m, nil
}

func (m testModel) View() string {
	return "Test TUI View"
}

// Integration test for complete navigation workflow
func TestNavigationIntegration(t *testing.T) {
	model := newTestModel()

	// Test sequence: navigate, search, toggle, exit
	testSequence := []struct {
		key      string
		expected func(testModel) bool
		desc     string
	}{
		{
			key: "down",
			expected: func(m testModel) bool {
				return m.selectedMCP == 1
			},
			desc: "Navigate down to second MCP",
		},
		{
			key: "right",
			expected: func(m testModel) bool {
				return m.activeColumn == 1
			},
			desc: "Move to second column",
		},
		{
			key: "right",
			expected: func(m testModel) bool {
				return m.activeColumn == 2
			},
			desc: "Move to third column",
		},
		{
			key: "left",
			expected: func(m testModel) bool {
				return m.activeColumn == 1
			},
			desc: "Move back to second column",
		},
		{
			key: "tab",
			expected: func(m testModel) bool {
				return m.searchFocused && m.activeColumn == -1
			},
			desc: "Focus search field",
		},
		{
			key: "h",
			expected: func(m testModel) bool {
				return m.searchQuery == "h"
			},
			desc: "Type in search field",
		},
		{
			key: "t",
			expected: func(m testModel) bool {
				return m.searchQuery == "ht"
			},
			desc: "Continue typing in search field",
		},
		{
			key: "backspace",
			expected: func(m testModel) bool {
				return m.searchQuery == "h"
			},
			desc: "Backspace in search field",
		},
		{
			key: "esc",
			expected: func(m testModel) bool {
				return !m.searchFocused && m.searchQuery == "" && m.activeColumn == 0
			},
			desc: "Exit search mode",
		},
		{
			key: "space",
			expected: func(m testModel) bool {
				// The second MCP (index 1) starts as true, so after toggle it should be false
				return m.mcps[1].Active == false
			},
			desc: "Toggle MCP active state",
		},
	}

	currentModel := model

	for i, test := range testSequence {
		var keyMsg tea.KeyMsg
		switch test.key {
		case "down":
			keyMsg = tea.KeyMsg{Type: tea.KeyDown}
		case "up":
			keyMsg = tea.KeyMsg{Type: tea.KeyUp}
		case "left":
			keyMsg = tea.KeyMsg{Type: tea.KeyLeft}
		case "right":
			keyMsg = tea.KeyMsg{Type: tea.KeyRight}
		case "tab":
			keyMsg = tea.KeyMsg{Type: tea.KeyTab}
		case "esc":
			keyMsg = tea.KeyMsg{Type: tea.KeyEscape}
		case "space":
			keyMsg = tea.KeyMsg{Type: tea.KeySpace}
		case "backspace":
			keyMsg = tea.KeyMsg{Type: tea.KeyBackspace}
		default:
			keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(test.key)}
		}

		updatedModel, _ := currentModel.Update(keyMsg)
		currentModel = updatedModel.(testModel)

		if !test.expected(currentModel) {
			t.Errorf("Step %d failed: %s", i+1, test.desc)
			t.Logf("Current state: selectedMCP=%d, activeColumn=%d, searchFocused=%t, searchQuery='%s'",
				currentModel.selectedMCP, currentModel.activeColumn, currentModel.searchFocused, currentModel.searchQuery)
		}
	}
}

// Test responsive layout behavior
func TestResponsiveLayoutIntegration(t *testing.T) {
	model := newTestModel()

	// Test width changes and navigation constraints
	widthTests := []struct {
		width           int
		expectedColumns int
		maxActiveColumn int
		desc            string
	}{
		{50, 1, 0, "Narrow terminal (1 column)"},
		{80, 2, 1, "Medium terminal (2 columns)"},
		{150, 3, 2, "Wide terminal (3 columns)"},
	}

	for _, test := range widthTests {
		t.Run(test.desc, func(t *testing.T) {
			// Update window size
			resizeMsg := tea.WindowSizeMsg{Width: test.width, Height: 24}
			updatedModel, _ := model.Update(resizeMsg)
			m := updatedModel.(testModel)

			if m.columns != test.expectedColumns {
				t.Errorf("Expected %d columns for width %d, got %d", test.expectedColumns, test.width, m.columns)
			}

			// Try to navigate to the rightmost column
			for i := 0; i < 5; i++ { // Try multiple right arrows
				rightMsg := tea.KeyMsg{Type: tea.KeyRight}
				updatedModel, _ = m.Update(rightMsg)
				m = updatedModel.(testModel)
			}

			if m.activeColumn > test.maxActiveColumn {
				t.Errorf("activeColumn should not exceed %d for %d columns, got %d",
					test.maxActiveColumn, test.expectedColumns, m.activeColumn)
			}
		})
	}
}

// Test search integration with navigation
func TestSearchNavigationIntegration(t *testing.T) {
	model := newTestModel()

	// Enter search mode with '/' key
	slashMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")}
	updatedModel, _ := model.Update(slashMsg)
	m := updatedModel.(testModel)

	if !m.searchFocused {
		t.Error("'/' key should focus search")
	}

	// Navigation keys should not work in search mode
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ = m.Update(downMsg)
	m = updatedModel.(testModel)

	if m.selectedMCP != 0 {
		t.Error("Navigation should not work while search is focused")
	}

	// Type search query
	searchText := "context"
	for _, char := range searchText {
		charMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
		updatedModel, _ = m.Update(charMsg)
		m = updatedModel.(testModel)
	}

	if m.searchQuery != searchText {
		t.Errorf("Expected search query '%s', got '%s'", searchText, m.searchQuery)
	}

	// Enter should exit search mode
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = m.Update(enterMsg)
	m = updatedModel.(testModel)

	if m.searchFocused {
		t.Error("Enter should exit search mode")
	}

	if m.activeColumn != 0 {
		t.Error("Should return to first column after search")
	}
}

// Test edge cases and error conditions
func TestNavigationEdgeCases(t *testing.T) {
	model := newTestModel()

	// Test navigation at boundaries
	// Try to go up from first item
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ := model.Update(upMsg)
	m := updatedModel.(testModel)

	if m.selectedMCP != 0 {
		t.Error("Should stay at first item when trying to go up")
	}

	// Navigate to last item
	for i := 0; i < len(model.mcps); i++ {
		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		updatedModel, _ = m.Update(downMsg)
		m = updatedModel.(testModel)
	}

	// Try to go down from last item
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ = m.Update(downMsg)
	m = updatedModel.(testModel)

	if m.selectedMCP != len(model.mcps)-1 {
		t.Error("Should stay at last item when trying to go down")
	}

	// Test left navigation at first column
	leftMsg := tea.KeyMsg{Type: tea.KeyLeft}
	updatedModel, _ = model.Update(leftMsg)
	m = updatedModel.(testModel)

	if m.activeColumn != 0 {
		t.Error("Should stay at first column when trying to go left")
	}

	// Navigate to last column
	for i := 0; i < model.columns; i++ {
		rightMsg := tea.KeyMsg{Type: tea.KeyRight}
		updatedModel, _ = m.Update(rightMsg)
		m = updatedModel.(testModel)
	}

	// Try to go right from last column
	rightMsg := tea.KeyMsg{Type: tea.KeyRight}
	updatedModel, _ = m.Update(rightMsg)
	m = updatedModel.(testModel)

	if m.activeColumn >= model.columns {
		t.Error("Should not exceed maximum column index")
	}
}

// Test quit functionality
func TestQuitIntegration(t *testing.T) {
	model := newTestModel()

	// Test 'q' key quits when not in search mode
	qMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	_, cmd := model.Update(qMsg)

	if cmd == nil {
		t.Error("'q' key should return quit command when not in search mode")
	}

	// Test 'q' key doesn't quit when in search mode
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(tabMsg)
	m := updatedModel.(testModel)

	_, cmd = m.Update(qMsg)
	if cmd != nil {
		t.Error("'q' key should not quit when in search mode")
	}

	// Test ESC followed by 'q' quits
	escMsg := tea.KeyMsg{Type: tea.KeyEscape}
	updatedModel, _ = m.Update(escMsg)
	m = updatedModel.(testModel)

	_, cmd = m.Update(qMsg)
	if cmd == nil {
		t.Error("'q' key should quit after ESC exits search mode")
	}
}

// Performance test for rapid key sequences
func TestRapidKeySequence(t *testing.T) {
	model := newTestModel()

	start := time.Now()

	// Simulate rapid navigation
	keys := []tea.KeyMsg{
		{Type: tea.KeyDown},
		{Type: tea.KeyDown},
		{Type: tea.KeyUp},
		{Type: tea.KeyRight},
		{Type: tea.KeyLeft},
		{Type: tea.KeyTab},
		{Type: tea.KeyEscape},
		{Type: tea.KeySpace},
	}

	currentModel := model
	for i := 0; i < 100; i++ { // Repeat sequence 100 times
		for _, key := range keys {
			updatedModel, _ := currentModel.Update(key)
			currentModel = updatedModel.(testModel)
		}
	}

	duration := time.Since(start)
	if duration > time.Millisecond*100 {
		t.Errorf("Rapid key sequence took too long: %v", duration)
	}
}
