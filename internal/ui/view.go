package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the application interface
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Build the complete interface
	header := m.renderHeader()
	body := m.renderBody()
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

// renderHeader creates the application header with shortcuts and context
func (m Model) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		Background(lipgloss.Color("#1E1E2E")).
		Padding(0, 2).
		Width(m.width)

	// Build shortcuts display based on current state
	var shortcuts string
	switch m.state {
	case MainNavigation:
		shortcuts = "A=Add • D=Delete • E=Edit • /=Search • Tab=Focus Search • ESC=Exit • ↑↓←→=Navigate"
	case SearchMode:
		shortcuts = "Type to search • Enter=Apply • ESC=Cancel"
	case ModalActive:
		shortcuts = "Enter=Confirm • ESC=Cancel"
	}

	// Context information
	activeCount := 0
	for _, item := range m.mcpItems {
		if item.Active {
			activeCount++
		}
	}

	contextInfo := fmt.Sprintf("MCPs: %d/%d Active • Layout: %s",
		activeCount, len(m.mcpItems), m.getLayoutName())

	title := "MCP Manager v1.0"

	// Create header content with proper spacing
	headerContent := fmt.Sprintf("%s\n%s\n%s", title, shortcuts, contextInfo)

	return headerStyle.Render(headerContent)
}

// renderBody creates the main application body with columns
func (m Model) renderBody() string {
	if m.columnCount == 1 {
		return m.renderSingleColumn()
	} else if m.columnCount == 2 {
		return m.renderTwoColumns()
	} else if m.columnCount == 4 {
		return m.renderFourColumns()
	}
	return m.renderThreeColumns()
}

// renderSingleColumn renders the narrow layout with 1 column
func (m Model) renderSingleColumn() string {
	columnStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1).
		Width(m.width - 4).
		Height(m.height - 8) // Account for header and footer

	if m.activeColumn == 0 {
		columnStyle = columnStyle.BorderForeground(lipgloss.Color("#7C3AED"))
	}

	content := m.renderMCPList()

	return columnStyle.Render(fmt.Sprintf("MCP Manager\n\n%s", content))
}

// renderTwoColumns renders the medium layout with 2 columns
func (m Model) renderTwoColumns() string {
	columnWidth := (m.width - 6) / 2
	columnHeight := m.height - 8

	// Column styles
	leftStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1).
		Width(columnWidth).
		Height(columnHeight)

	rightStyle := leftStyle.Copy()

	// Highlight active column
	if m.activeColumn == 0 {
		leftStyle = leftStyle.BorderForeground(lipgloss.Color("#7C3AED"))
	} else {
		rightStyle = rightStyle.BorderForeground(lipgloss.Color("#7C3AED"))
	}

	// Column content
	leftContent := fmt.Sprintf("MCPs\n\n%s", m.renderMCPList())
	rightContent := fmt.Sprintf("Status & Details\n\n%s", m.renderStatusAndDetails())

	leftColumn := leftStyle.Render(leftContent)
	rightColumn := rightStyle.Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
}

// renderThreeColumns renders the wide layout with 3 columns
func (m Model) renderThreeColumns() string {
	columnWidth := (m.width - 8) / 3
	columnHeight := m.height - 8

	// Base column style
	columnStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1).
		Width(columnWidth).
		Height(columnHeight)

	// Create three columns with conditional highlighting
	columns := make([]string, 3)

	for i := 0; i < 3; i++ {
		style := columnStyle.Copy()
		if i == m.activeColumn {
			style = style.BorderForeground(lipgloss.Color("#7C3AED"))
		}

		var content string
		switch i {
		case 0:
			content = fmt.Sprintf("MCPs\n\n%s", m.renderMCPList())
		case 1:
			content = fmt.Sprintf("Status\n\n%s", m.renderStatusColumn())
		case 2:
			content = fmt.Sprintf("Details\n\n%s", m.renderDetailsColumn())
		}

		columns[i] = style.Render(content)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, columns[0], columns[1], columns[2])
}

// renderFourColumns renders the clean 4x10 MCP grid layout without column separators
func (m Model) renderFourColumns() string {
	// Get filtered MCPs for search functionality
	filteredMCPs := m.GetFilteredMCPs()
	
	if len(filteredMCPs) == 0 {
		// Show "No results" message when search returns no results
		noResultsStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Align(lipgloss.Center).
			Width(m.width).
			Height(m.height - 8)
		return noResultsStyle.Render("No MCPs found matching your search")
	}

	// Calculate grid dimensions - aim for ~10 rows with 4 columns
	gridRows := (len(filteredMCPs) + 3) / 4 // Round up division
	if gridRows < 10 {
		gridRows = 10 // Minimum 10 rows for consistent layout
	}

	// Build the grid as a simple string without column separators
	var gridLines []string
	
	for row := 0; row < gridRows; row++ {
		var line []string
		
		for col := 0; col < 4; col++ {
			mcpIndex := row*4 + col
			
			if mcpIndex < len(filteredMCPs) {
				item := filteredMCPs[mcpIndex]
				
				// Status indicator
				status := "○"
				if item.Active {
					status = "●"
				}
				
				// Highlight selected item by comparing index directly
				isSelected := (mcpIndex == m.selectedItem)
				
				// Create base item text
				itemText := fmt.Sprintf("%s %s", status, item.Name)
				
				if isSelected {
					// Simple, clean selection highlighting - just bold text
					itemStyle := lipgloss.NewStyle().Bold(true)
					itemText = itemStyle.Render(itemText)
				}
				
				// Fixed width for consistent spacing (about 28 chars per column)
				line = append(line, fmt.Sprintf("%-28s", itemText))
			} else {
				// Empty cell with proper spacing
				line = append(line, fmt.Sprintf("%-28s", ""))
			}
		}
		
		// Join columns without separators, just spaces
		gridLines = append(gridLines, strings.Join(line, ""))
	}
	
	// Join all rows with newlines
	gridContent := strings.Join(gridLines, "\n")
	
	// Apply overall styling to the grid
	gridStyle := lipgloss.NewStyle().
		Padding(2).
		Width(m.width).
		Height(m.height - 8)
	
	return gridStyle.Render(gridContent)
}

// renderMCPList renders the list of MCPs with selection highlighting
func (m Model) renderMCPList() string {
	if len(m.mcpItems) == 0 {
		return "No MCPs configured"
	}

	var items []string
	for i, item := range m.mcpItems {
		style := lipgloss.NewStyle().Padding(0, 1)

		// Highlight selected item
		if i == m.selectedItem {
			style = style.Background(lipgloss.Color("#7C3AED")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)
		}

		// Status indicator
		status := "○"
		if item.Active {
			status = "●"
		}

		itemText := fmt.Sprintf("%s %s", status, item.Name)
		items = append(items, style.Render(itemText))
	}

	return strings.Join(items, "\n")
}

// renderMCPColumnList renders a list of MCPs for a specific column with selection highlighting
func (m Model) renderMCPColumnList(columnMCPs []MCPItem, startIdx int) string {
	if len(columnMCPs) == 0 {
		return ""
	}

	var items []string
	for i, item := range columnMCPs {
		actualIdx := startIdx + i
		style := lipgloss.NewStyle().Padding(0, 1)

		// Highlight selected item
		if actualIdx == m.selectedItem {
			style = style.Background(lipgloss.Color("#7C3AED")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)
		}

		// Status indicator
		status := "○"
		if item.Active {
			status = "●"
		}

		itemText := fmt.Sprintf("%s %s", status, item.Name)
		items = append(items, style.Render(itemText))
	}

	return strings.Join(items, "\n")
}

// renderStatusColumn renders the status information for the selected MCP
func (m Model) renderStatusColumn() string {
	if m.selectedItem >= len(m.mcpItems) {
		return "No MCP selected"
	}

	item := m.mcpItems[m.selectedItem]

	status := "Inactive"
	statusColor := "#FF6B6B"
	if item.Active {
		status = "Active"
		statusColor = "#51CF66"
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(statusColor)).
		Bold(true)

	return fmt.Sprintf("Name: %s\nType: %s\nStatus: %s\n\nCommand:\n%s",
		item.Name,
		item.Type,
		statusStyle.Render(status),
		item.Command)
}

// renderDetailsColumn renders detailed information for the selected MCP
func (m Model) renderDetailsColumn() string {
	if m.selectedItem >= len(m.mcpItems) {
		return "No MCP selected"
	}

	item := m.mcpItems[m.selectedItem]

	// Placeholder details - will be expanded in future stories
	details := []string{
		fmt.Sprintf("MCP: %s", item.Name),
		"Configuration:",
		"  • Auto-start: Yes",
		"  • Timeout: 30s",
		"  • Retry: 3x",
		"",
		"Capabilities:",
		"  • File operations",
		"  • Context search",
		"  • Documentation",
		"",
		"Last started:",
		"  2025-06-29 14:30:22",
	}

	return strings.Join(details, "\n")
}

// renderStatusAndDetails renders combined status and details for 2-column layout
func (m Model) renderStatusAndDetails() string {
	status := m.renderStatusColumn()
	details := m.renderDetailsColumn()

	return fmt.Sprintf("%s\n\n%s", status, details)
}

// renderFooter creates the application footer
func (m Model) renderFooter() string {
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 2).
		Width(m.width)

	var footerText string
	if m.searchActive {
		searchStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("#7C3AED")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

		cursor := "_"
		footerText = fmt.Sprintf("Search: %s", searchStyle.Render(m.searchQuery+cursor))
	} else if m.searchQuery != "" {
		// Show search results info when not actively searching but have a query
		filteredMCPs := m.GetFilteredMCPs()
		footerText = fmt.Sprintf("Found %d MCPs matching '%s' • ESC to clear • Terminal: %dx%d",
			len(filteredMCPs), m.searchQuery, m.width, m.height)
	} else {
		footerText = fmt.Sprintf("Terminal: %dx%d • Search: '%s' • Use arrow keys to navigate, Tab or / for search",
			m.width, m.height, m.searchQuery)
	}

	return footerStyle.Render(footerText)
}

// getLayoutName returns the current layout name for display
func (m Model) getLayoutName() string {
	switch m.columnCount {
	case 1:
		return "Narrow"
	case 2:
		return "Medium"
	case 3:
		return "Wide (3-panel)"
	case 4:
		return "Grid (4-column MCP)"
	default:
		return "Unknown"
	}
}
