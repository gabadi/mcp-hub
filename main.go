package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type Model struct {
	width         int
	height        int
	activeColumn  int  // 0, 1, or 2 for the three columns
	searchFocused bool // true when search field is focused
	searchQuery   string
	windowWidth   int
	columns       int // number of columns to display (1, 2, or 3)

	// Mock MCP data for initial development
	mcps        []MCP
	selectedMCP int
}

// MCP represents a Model Context Protocol configuration
type MCP struct {
	Name        string
	Type        string // "CMD", "SSE", "JSON", "HTTP"
	Description string
	Active      bool
}

// Initial model state
func initialModel() Model {
	// Mock MCPs for development
	mockMCPs := []MCP{
		{Name: "context7", Type: "SSE", Description: "Library version lookup and best practices", Active: false},
		{Name: "ht-mcp", Type: "CMD", Description: "Manual testing and validation tool", Active: true},
		{Name: "filesystem", Type: "CMD", Description: "File system operations", Active: false},
		{Name: "git", Type: "CMD", Description: "Git repository management", Active: true},
		{Name: "web-search", Type: "HTTP", Description: "Web search capabilities", Active: false},
	}

	return Model{
		mcps:         mockMCPs,
		activeColumn: 0,
		columns:      3, // default to 3 columns
		selectedMCP:  0,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.windowWidth = msg.Width

		// Responsive column layout based on width
		if msg.Width < 60 {
			m.columns = 1
		} else if msg.Width < 120 {
			m.columns = 2
		} else {
			m.columns = 3
		}

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		m.activeColumn = -1 // indicate search is focused
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
		// Toggle MCP active state
		if m.selectedMCP < len(m.mcps) {
			m.mcps[m.selectedMCP].Active = !m.mcps[m.selectedMCP].Active
		}
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "q":
			return m, tea.Quit
		case "k":
			if m.selectedMCP > 0 {
				m.selectedMCP--
			}
		case "j":
			if m.selectedMCP < len(m.mcps)-1 {
				m.selectedMCP++
			}
		case "h":
			if m.activeColumn > 0 {
				m.activeColumn--
			}
		case "l":
			if m.activeColumn < m.columns-1 {
				m.activeColumn++
			}
		case "/":
			m.searchFocused = true
			m.activeColumn = -1
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Header
	header := m.renderHeader()

	// Main content area
	content := m.renderContent()

	// Footer with status
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// renderHeader creates the header with shortcuts and context
func (m Model) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Bold(true)

	shortcuts := "Tab: Search | ↑↓: Navigate | ←→: Columns | Space: Toggle | ESC/q: Exit"
	context := fmt.Sprintf("Columns: %d | MCPs: %d", m.columns, len(m.mcps))

	headerLeft := headerStyle.Render(shortcuts)
	headerRight := headerStyle.Render(context)

	// Calculate spacing
	totalHeaderLength := lipgloss.Width(headerLeft) + lipgloss.Width(headerRight)
	spacingNeeded := m.width - totalHeaderLength
	if spacingNeeded < 0 {
		spacingNeeded = 0
	}

	spacing := lipgloss.NewStyle().
		Background(lipgloss.Color("#7D56F4")).
		Render(fmt.Sprintf("%*s", spacingNeeded, ""))

	return lipgloss.JoinHorizontal(lipgloss.Top, headerLeft, spacing, headerRight)
}

// renderContent creates the main 3-column layout
func (m Model) renderContent() string {
	contentHeight := m.height - 4 // Account for header and footer

	// Search bar
	searchBar := m.renderSearchBar()

	// MCP list (main column)
	mcpList := m.renderMCPList(contentHeight - 2) // Account for search bar

	// Create columns based on responsive layout
	switch m.columns {
	case 1:
		return lipgloss.JoinVertical(lipgloss.Left, searchBar, mcpList)
	case 2:
		leftColumn := lipgloss.JoinVertical(lipgloss.Left, searchBar, mcpList)
		rightColumn := m.renderDetailsColumn(contentHeight)

		leftStyle := lipgloss.NewStyle().Width(m.width / 2)
		rightStyle := lipgloss.NewStyle().Width(m.width / 2)

		return lipgloss.JoinHorizontal(lipgloss.Top,
			leftStyle.Render(leftColumn),
			rightStyle.Render(rightColumn))
	case 3:
		leftColumn := lipgloss.JoinVertical(lipgloss.Left, searchBar, mcpList)
		centerColumn := m.renderDetailsColumn(contentHeight)
		rightColumn := m.renderActionsColumn(contentHeight)

		colWidth := m.width / 3
		colStyle := lipgloss.NewStyle().Width(colWidth)

		return lipgloss.JoinHorizontal(lipgloss.Top,
			colStyle.Render(leftColumn),
			colStyle.Render(centerColumn),
			colStyle.Render(rightColumn))
	default:
		return "Invalid column configuration"
	}
}

// renderSearchBar creates the search input field
func (m Model) renderSearchBar() string {
	searchStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	if m.searchFocused {
		searchStyle = searchStyle.BorderForeground(lipgloss.Color("#7D56F4"))
	} else {
		searchStyle = searchStyle.BorderForeground(lipgloss.Color("#444444"))
	}

	prompt := "Search MCPs: "
	cursor := ""
	if m.searchFocused {
		cursor = "█"
	}

	searchContent := prompt + m.searchQuery + cursor
	return searchStyle.Render(searchContent)
}

// renderMCPList creates the main MCP list
func (m Model) renderMCPList(height int) string {
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1)

	if m.activeColumn == 0 && !m.searchFocused {
		listStyle = listStyle.BorderForeground(lipgloss.Color("#7D56F4"))
	}

	var items []string
	filteredMCPs := m.getFilteredMCPs()

	for i, mcp := range filteredMCPs {
		status := "○"
		if mcp.Active {
			status = "●"
		}

		typeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render(fmt.Sprintf("[%s]", mcp.Type))

		line := fmt.Sprintf("%s %s %s", status, mcp.Name, typeStyle)

		if i == m.selectedMCP && m.activeColumn == 0 && !m.searchFocused {
			line = lipgloss.NewStyle().
				Background(lipgloss.Color("#7D56F4")).
				Foreground(lipgloss.Color("#FAFAFA")).
				Render(line)
		}

		items = append(items, line)
	}

	if len(items) == 0 {
		items = append(items, "No MCPs found")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	return listStyle.Height(height).Render(content)
}

// renderDetailsColumn creates the details/info column
func (m Model) renderDetailsColumn(height int) string {
	detailStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1).
		Height(height)

	if m.activeColumn == 1 && !m.searchFocused {
		detailStyle = detailStyle.BorderForeground(lipgloss.Color("#7D56F4"))
	}

	if m.selectedMCP < len(m.mcps) {
		mcp := m.mcps[m.selectedMCP]
		content := fmt.Sprintf("Name: %s\nType: %s\nStatus: %s\n\nDescription:\n%s",
			mcp.Name,
			mcp.Type,
			map[bool]string{true: "Active", false: "Inactive"}[mcp.Active],
			mcp.Description)
		return detailStyle.Render(content)
	}

	return detailStyle.Render("Select an MCP to view details")
}

// renderActionsColumn creates the actions/help column
func (m Model) renderActionsColumn(height int) string {
	actionStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1).
		Height(height)

	if m.activeColumn == 2 && !m.searchFocused {
		actionStyle = actionStyle.BorderForeground(lipgloss.Color("#7D56F4"))
	}

	actions := `Actions:
	
Space - Toggle MCP
A - Add new MCP
E - Edit MCP
D - Delete MCP
R - Refresh status

Navigation:
↑↓ - Select MCP
←→ - Switch columns
Tab - Focus search
/ - Quick search
ESC - Exit/Clear`

	return actionStyle.Render(actions)
}

// renderFooter creates the status footer
func (m Model) renderFooter() string {
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Background(lipgloss.Color("#2A2A2A")).
		Padding(0, 1)

	activeCount := 0
	for _, mcp := range m.mcps {
		if mcp.Active {
			activeCount++
		}
	}

	status := fmt.Sprintf("Active: %d/%d | Selected: %s | Terminal: %dx%d",
		activeCount, len(m.mcps),
		func() string {
			if m.selectedMCP < len(m.mcps) {
				return m.mcps[m.selectedMCP].Name
			}
			return "None"
		}(),
		m.width, m.height)

	// Fill the entire width
	spacingNeeded := m.width - lipgloss.Width(status)
	if spacingNeeded < 0 {
		spacingNeeded = 0
	}

	statusWithSpacing := status + fmt.Sprintf("%*s", spacingNeeded, "")
	return footerStyle.Render(statusWithSpacing)
}

// getFilteredMCPs returns MCPs filtered by search query
func (m Model) getFilteredMCPs() []MCP {
	if m.searchQuery == "" {
		return m.mcps
	}

	var filtered []MCP
	for _, mcp := range m.mcps {
		if contains(mcp.Name, m.searchQuery) || contains(mcp.Description, m.searchQuery) {
			filtered = append(filtered, mcp)
		}
	}
	return filtered
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(substr) == 0 ||
		len(s) >= len(substr) &&
			findIgnoreCase(s, substr) != -1
}

// findIgnoreCase finds substring in string ignoring case
func findIgnoreCase(s, substr string) int {
	sLower := toLower(s)
	substrLower := toLower(substr)

	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return i
		}
	}
	return -1
}

// toLower converts string to lowercase
func toLower(s string) string {
	result := make([]byte, len(s))
	for i, r := range []byte(s) {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
