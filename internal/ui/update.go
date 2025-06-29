package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()

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
		// Exit application
		return m, tea.Quit
	}
	return m, nil
}

// handleMainNavigationKeys handles keyboard input in main navigation mode
func (m Model) handleMainNavigationKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		m.navigateUp()
	case "down", "j":
		m.navigateDown()
	case "left", "h":
		m.navigateLeft()
	case "right", "l":
		m.navigateRight()
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
		// Toggle MCP active status
		if m.selectedItem < len(m.mcpItems) {
			m.mcpItems[m.selectedItem].Active = !m.mcpItems[m.selectedItem].Active
		}
	}

	return m, nil
}

// handleSearchModeKeys handles keyboard input in search mode
func (m Model) handleSearchModeKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "enter":
		// Process search and return to main navigation
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
func (m *Model) updateLayout() {
	// Responsive breakpoints as specified in acceptance criteria
	if m.width > 120 {
		// Wide: 3 columns
		m.columnCount = 3
		columnWidth := (m.width - 4) / 3 // Account for spacing
		m.columns = []Column{
			{Title: "MCPs", Width: columnWidth},
			{Title: "Status", Width: columnWidth},
			{Title: "Details", Width: columnWidth},
		}
	} else if m.width >= 80 {
		// Medium: 2 columns
		m.columnCount = 2
		columnWidth := (m.width - 2) / 2
		m.columns = []Column{
			{Title: "MCPs", Width: columnWidth},
			{Title: "Status & Details", Width: columnWidth},
		}
	} else {
		// Narrow: 1 column
		m.columnCount = 1
		m.columns = []Column{
			{Title: "MCP Manager", Width: m.width - 2},
		}
	}

	// Reset active column if it's out of bounds
	if m.activeColumn >= m.columnCount {
		m.activeColumn = m.columnCount - 1
	}
}

// Navigation helper methods

// navigateUp moves selection up within the current column
func (m *Model) navigateUp() {
	if m.activeColumn == 0 {
		// In MCP list column
		if m.selectedItem > 0 {
			m.selectedItem--
		}
	}
	// Other columns don't have navigable items yet
}

// navigateDown moves selection down within the current column
func (m *Model) navigateDown() {
	if m.activeColumn == 0 {
		// In MCP list column
		maxItems := len(m.mcpItems) - 1
		if m.selectedItem < maxItems {
			m.selectedItem++
		}
	}
	// Other columns don't have navigable items yet
}

// navigateLeft moves to the previous column
func (m *Model) navigateLeft() {
	if m.activeColumn > 0 {
		m.activeColumn--
	}
}

// navigateRight moves to the next column
func (m *Model) navigateRight() {
	if m.activeColumn < m.columnCount-1 {
		m.activeColumn++
	}
}