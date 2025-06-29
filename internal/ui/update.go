package ui

import (
	"strings"
	
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.updateLayout()

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress processes keyboard input based on current state
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Global keys that work in any state
	switch key {
	case "esc":
		return m.handleEscKey()
	case "ctrl+c":
		return m, tea.Quit
	}

	// State-specific key handling
	switch m.state {
	case MainNavigation:
		return m.handleMainNavigationKeys(key)
	case SearchMode:
		return m.handleSearchModeKeys(key)
	case ModalActive:
		return m.handleModalKeys(key)
	}

	return m, nil
}

// handleEscKey handles ESC key behavior based on current state
func (m Model) handleEscKey() (tea.Model, tea.Cmd) {
	switch m.state {
	case SearchMode:
		// Clear search and return to main navigation
		m.searchActive = false
		m.searchQuery = ""
		m.state = MainNavigation
		return m, nil
	case ModalActive:
		// Close modal and return to main navigation
		m.state = MainNavigation
		return m, nil
	case MainNavigation:
		// Clear search if active, otherwise exit application
		if m.searchQuery != "" {
			m.searchQuery = ""
			m.searchResults = nil
			m.selectedItem = 0 // Reset selection
			return m, nil
		}
		// Exit application
		return m, tea.Quit
	}
	return m, nil
}

// handleMainNavigationKeys handles keyboard input in main navigation mode
func (m Model) handleMainNavigationKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		m = m.navigateUp()
	case "down", "j":
		m = m.navigateDown()
	case "left", "h":
		m = m.navigateLeft()
	case "right", "l":
		m = m.navigateRight()
	case "tab":
		// Jump to search field
		m.state = SearchMode
		m.searchActive = true
	case "/":
		// Activate search mode
		m.state = SearchMode
		m.searchActive = true
	case "a":
		// Add MCP (future functionality)
		m.state = ModalActive
	case "e":
		// Edit MCP (future functionality)
		m.state = ModalActive
	case "d":
		// Delete MCP (future functionality)
		m.state = ModalActive
	case " ", "space":
		// Toggle MCP active status - work with filtered results
		filteredMCPs := m.GetFilteredMCPs()
		if m.selectedItem < len(filteredMCPs) {
			// Find the original item and toggle it
			selectedItem := filteredMCPs[m.selectedItem]
			for i := range m.mcpItems {
				if m.mcpItems[i].Name == selectedItem.Name {
					m.mcpItems[i].Active = !m.mcpItems[i].Active
					break
				}
			}
		}
	}

	return m, nil
}

// handleSearchModeKeys handles keyboard input in search mode
func (m Model) handleSearchModeKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "enter":
		// Return to main navigation with search query preserved
		m.state = MainNavigation
		m.searchActive = false
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
	default:
		// Add character to search query
		if len(key) == 1 {
			m.searchQuery += key
		}
	}

	return m, nil
}

// handleModalKeys handles keyboard input in modal mode
func (m Model) handleModalKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "enter":
		// Confirm modal action and return to main navigation
		m.state = MainNavigation
	}

	return m, nil
}

// updateLayout updates the column layout based on terminal width
func (m Model) updateLayout() Model {
	// Responsive breakpoints as specified in acceptance criteria
	// Prioritize 4-column MCP grid layout as per wireframe
	if m.width >= 120 {
		// Wide: 4-column MCP grid for maximum information density
		m.columnCount = 4
		columnWidth := (m.width - 10) / 4 // Account for spacing between 4 columns
		m.columns = []Column{
			{Title: "MCPs Column 1", Width: columnWidth},
			{Title: "MCPs Column 2", Width: columnWidth},
			{Title: "MCPs Column 3", Width: columnWidth},
			{Title: "MCPs Column 4", Width: columnWidth},
		}
	} else if m.width >= 80 {
		// Medium: 2 columns (MCPs + Status/Details)
		m.columnCount = 2
		columnWidth := (m.width - 6) / 2
		m.columns = []Column{
			{Title: "MCPs", Width: columnWidth},
			{Title: "Status & Details", Width: columnWidth},
		}
	} else {
		// Narrow: 1 column (all in one)
		m.columnCount = 1
		m.columns = []Column{
			{Title: "MCP Manager", Width: m.width - 4},
		}
	}

	// Reset active column if it's out of bounds
	if m.activeColumn >= m.columnCount {
		m.activeColumn = m.columnCount - 1
	}
	return m
}

// Navigation helper methods

// navigateUp moves selection up within the grid
func (m Model) navigateUp() Model {
	if m.columnCount == 4 {
		// In 4-column grid - move up one row (subtract 4 from index)
		filteredMCPs := m.GetFilteredMCPs()
		if m.selectedItem >= 4 {
			m.selectedItem -= 4
		}
		// Ensure we don't go below 0 and stay within filtered results
		if m.selectedItem < 0 {
			m.selectedItem = 0
		}
		if m.selectedItem >= len(filteredMCPs) && len(filteredMCPs) > 0 {
			m.selectedItem = len(filteredMCPs) - 1
		}
	} else if m.activeColumn == 0 {
		// In MCP list column for other layouts
		if m.selectedItem > 0 {
			m.selectedItem--
		}
	}
	// Other columns don't have navigable items yet
	return m
}

// navigateDown moves selection down within the grid
func (m Model) navigateDown() Model {
	if m.columnCount == 4 {
		// In 4-column grid - move down one row (add 4 to index)
		filteredMCPs := m.GetFilteredMCPs()
		newIndex := m.selectedItem + 4
		if newIndex < len(filteredMCPs) {
			m.selectedItem = newIndex
		}
		// Stay within filtered results bounds
		if m.selectedItem >= len(filteredMCPs) && len(filteredMCPs) > 0 {
			m.selectedItem = len(filteredMCPs) - 1
		}
	} else if m.activeColumn == 0 {
		// In MCP list column for other layouts
		maxItems := len(m.mcpItems) - 1
		if m.selectedItem < maxItems {
			m.selectedItem++
		}
	}
	// Other columns don't have navigable items yet
	return m
}

// navigateLeft moves to the left within the grid
func (m Model) navigateLeft() Model {
	if m.columnCount == 4 {
		// In 4-column grid - move left within current row
		if m.selectedItem % 4 > 0 {
			m.selectedItem--
		}
	} else {
		// For other layouts, move between columns
		if m.activeColumn > 0 {
			m.activeColumn--
		}
	}
	return m
}

// navigateRight moves to the right within the grid
func (m Model) navigateRight() Model {
	if m.columnCount == 4 {
		// In 4-column grid - move right within current row
		filteredMCPs := m.GetFilteredMCPs()
		if m.selectedItem % 4 < 3 && m.selectedItem + 1 < len(filteredMCPs) {
			m.selectedItem++
		}
	} else {
		// For other layouts, move between columns
		if m.activeColumn < m.columnCount-1 {
			m.activeColumn++
		}
	}
	return m
}

// updateSelectedItemForColumn updates the selected item when switching columns in 4-column layout
func (m Model) updateSelectedItemForColumn() Model {
	if m.columnCount != 4 {
		return m
	}
	
	mcpsPerColumn := (len(m.mcpItems) + 3) / 4
	startIdx := m.activeColumn * mcpsPerColumn
	endIdx := startIdx + mcpsPerColumn
	if endIdx > len(m.mcpItems) {
		endIdx = len(m.mcpItems)
	}
	
	// Set selected item to first item in new column
	if startIdx < len(m.mcpItems) {
		m.selectedItem = startIdx
	}
	return m
}

// applySearch filters MCPs based on search query
func (m Model) applySearch() Model {
	if m.searchQuery == "" {
		// Reset to show all MCPs when search is empty
		m.searchResults = nil
		return m
	}
	
	// Clear previous results and initialize if needed
	m.searchResults = make([]string, 0)
	
	// Search through MCP names (case-insensitive)
	query := strings.ToLower(m.searchQuery)
	for _, item := range m.mcpItems {
		if strings.Contains(strings.ToLower(item.Name), query) {
			m.searchResults = append(m.searchResults, item.Name)
		}
	}
	return m
}
