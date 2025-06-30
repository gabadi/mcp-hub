package internal

import (
	"testing"

	"cc-mcp-manager/internal/ui"
	"cc-mcp-manager/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

// TestResponsiveBreakpoints validates the responsive layout logic
func TestResponsiveBreakpoints(t *testing.T) {
	tests := []struct {
		width    int
		expected string
	}{
		{150, "wide"},   // >=120 chars → 4 columns
		{100, "medium"}, // 80-119 chars → 2 columns
		{60, "narrow"},  // <80 chars → 1 column
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := getLayoutType(tt.width)
			if result != tt.expected {
				t.Errorf("getLayoutType(%d) = %s, want %s", tt.width, result, tt.expected)
			}
		})
	}
}

// getLayoutType determines layout based on terminal width
func getLayoutType(width int) string {
	if width >= 120 {
		return "wide"
	} else if width >= 80 {
		return "medium"
	}
	return "narrow"
}

// TestModelLayoutAdaptation tests the responsive layout adaptation
func TestModelLayoutAdaptation(t *testing.T) {
	model := ui.NewModel()

	tests := []struct {
		name         string
		width        int
		height       int
		expectedCols int
	}{
		{"Wide layout", 150, 50, 4},
		{"Medium layout", 100, 40, 2},
		{"Narrow layout", 60, 30, 1},
		{"Edge case - exactly 120", 120, 40, 4},
		{"Edge case - exactly 80", 80, 30, 2},
		{"Edge case - 79 chars", 79, 30, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate window size change
			newModel, _ := model.Update(tea.WindowSizeMsg{Width: tt.width, Height: tt.height})
			model = newModel.(ui.Model)

			if model.GetColumnCount() != tt.expectedCols {
				t.Errorf("Width %d: expected %d columns, got %d",
					tt.width, tt.expectedCols, model.GetColumnCount())
			}
		})
	}
}

// TestNavigationLogic tests the navigation functionality
func TestNavigationLogic(t *testing.T) {
	model := ui.NewModel()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 150, Height: 50})
	model = newModel.(ui.Model)

	// For 4-column layout (width 150), navigation uses SelectedItem, not ActiveColumn
	// ActiveColumn stays 0 for grid navigation
	initialSelectedItem := model.GetSelectedItem()

	// Test right navigation (should move to next item in row)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = newModel.(ui.Model)
	if model.GetSelectedItem() != initialSelectedItem+1 {
		t.Errorf("Expected selected item %d, got %d", initialSelectedItem+1, model.GetSelectedItem())
	}

	// Test left navigation (should move back)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = newModel.(ui.Model)
	if model.GetSelectedItem() != initialSelectedItem {
		t.Errorf("Expected selected item %d, got %d", initialSelectedItem, model.GetSelectedItem())
	}

	// Test boundary conditions - can't go left from first item in row
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = newModel.(ui.Model)
	if model.GetSelectedItem() != initialSelectedItem {
		t.Errorf("Should stay at item %d, got %d", initialSelectedItem, model.GetSelectedItem())
	}

	// Navigate to position 3 (rightmost in first row for 4-column)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = newModel.(ui.Model)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = newModel.(ui.Model)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = newModel.(ui.Model)
	if model.GetSelectedItem() != 3 {
		t.Errorf("Expected selected item 3, got %d", model.GetSelectedItem())
	}

	// Test boundary conditions - can't go right from last item in row
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = newModel.(ui.Model)
	if model.GetSelectedItem() != 3 {
		t.Errorf("Should stay at item 3, got %d", model.GetSelectedItem())
	}
}

// TestItemNavigation tests navigation within the MCP list
func TestItemNavigation(t *testing.T) {
	model := ui.NewModel()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 150, Height: 50})
	model = newModel.(ui.Model)

	// Ensure we're in the first column (MCP list)
	if model.GetActiveColumn() != 0 {
		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
		model = newModel.(ui.Model)
	}

	initialItem := model.GetSelectedItem()

	// Test down navigation
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(ui.Model)
	if model.GetSelectedItem() <= initialItem {
		t.Errorf("Down navigation should increase selected item")
	}

	// Test up navigation
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(ui.Model)
	if model.GetSelectedItem() != initialItem {
		t.Errorf("Up navigation should return to initial item")
	}

	// Test boundary condition - can't go up from first item
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(ui.Model)
	if model.GetSelectedItem() != 0 {
		t.Errorf("Should stay at item 0, got %d", model.GetSelectedItem())
	}
}

// TestStateTransitions tests the application state transitions
func TestStateTransitions(t *testing.T) {
	model := ui.NewModel()

	// Test transition to search mode
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(ui.Model)
	if model.GetState() != types.SearchActiveNavigation {
		t.Errorf("Tab should activate search mode")
	}

	// Test ESC from search mode
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = newModel.(ui.Model)
	if model.GetState() != types.MainNavigation {
		t.Errorf("ESC should return to main navigation from search mode")
	}

	// Test slash key for search
	newModel, _ = model.Update(tea.KeyMsg{Runes: []rune{'/'}, Type: tea.KeyRunes})
	model = newModel.(ui.Model)
	if model.GetState() != types.SearchActiveNavigation {
		t.Errorf("/ should activate search mode")
	}
}

// TestSearchFunctionality tests the search mode functionality
func TestSearchFunctionality(t *testing.T) {
	model := ui.NewModel()

	// Enter search mode
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(ui.Model)

	// Type search query
	newModel, _ = model.Update(tea.KeyMsg{Runes: []rune{'t'}, Type: tea.KeyRunes})
	model = newModel.(ui.Model)
	newModel, _ = model.Update(tea.KeyMsg{Runes: []rune{'e'}, Type: tea.KeyRunes})
	model = newModel.(ui.Model)
	newModel, _ = model.Update(tea.KeyMsg{Runes: []rune{'s'}, Type: tea.KeyRunes})
	model = newModel.(ui.Model)
	newModel, _ = model.Update(tea.KeyMsg{Runes: []rune{'t'}, Type: tea.KeyRunes})
	model = newModel.(ui.Model)

	if model.GetSearchQuery() != "test" {
		t.Errorf("Expected search query 'test', got '%s'", model.GetSearchQuery())
	}

	// Test backspace
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = newModel.(ui.Model)
	if model.GetSearchQuery() != "tes" {
		t.Errorf("Expected search query 'tes' after backspace, got '%s'", model.GetSearchQuery())
	}

	// Test Enter to apply search
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = newModel.(ui.Model)
	if model.GetState() != types.MainNavigation {
		t.Errorf("Enter should return to main navigation after search")
	}
}
