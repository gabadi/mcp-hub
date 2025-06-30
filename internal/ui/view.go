package ui

import (
	"fmt"
	"strings"

	"cc-mcp-manager/internal/ui/components"
	"cc-mcp-manager/internal/ui/types"
	"github.com/charmbracelet/lipgloss"
)

// View renders the application interface
func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Loading..."
	}

	// Build the complete interface using components
	header := components.RenderHeader(m.Model)
	body := m.renderBody()
	footer := components.RenderFooter(m.Model)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

// renderHeader creates the application header with shortcuts and context
func (m Model) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		Background(lipgloss.Color("#1E1E2E")).
		Padding(0, 2).
		Width(m.Width)

	// Build shortcuts display based on current state
	var shortcuts string
	switch m.State {
	case types.MainNavigation:
		shortcuts = "A=Add • D=Delete • E=Edit • /=Search • Tab=Focus Search • ESC=Exit • ↑↓←→=Navigate"
	case types.SearchMode:
		shortcuts = "Type to search • Enter=Apply • ESC=Cancel"
	case types.SearchActiveNavigation:
		if m.SearchInputActive {
			shortcuts = "Type to search • Tab=Navigate Mode • ↑↓←→=Navigate • Space=Toggle • Enter=Apply • ESC=Cancel"
		} else {
			shortcuts = "Navigate Mode • Tab=Input Mode • ↑↓←→=Navigate • Space=Toggle • Enter=Apply • ESC=Cancel"
		}
	case types.ModalActive:
		shortcuts = "Enter=Confirm • ESC=Cancel"
	}

	// Context information
	activeCount := 0
	for _, item := range m.MCPItems {
		if item.Active {
			activeCount++
		}
	}

	contextInfo := fmt.Sprintf("MCPs: %d/%d Active • Layout: %s",
		activeCount, len(m.MCPItems), m.getLayoutName())

	title := "MCP Manager v1.0"

	// Create header content with proper spacing
	headerContent := fmt.Sprintf("%s\n%s\n%s", title, shortcuts, contextInfo)

	return headerStyle.Render(headerContent)
}

// renderBody creates the main application body with columns
func (m Model) renderBody() string {
	if m.ColumnCount == 1 {
		return m.renderSingleColumn()
	} else if m.ColumnCount == 2 {
		return m.renderTwoColumns()
	} else if m.ColumnCount == 4 {
		return components.RenderFourColumnGrid(m.Model)
	}
	return m.renderThreeColumns()
}

// renderSingleColumn renders the narrow layout with 1 column
func (m Model) renderSingleColumn() string {
	columnStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1).
		Width(m.Width - 4).
		Height(m.Height - 8) // Account for header and footer

	if m.ActiveColumn == 0 {
		columnStyle = columnStyle.BorderForeground(lipgloss.Color("#7C3AED"))
	}

	content := components.RenderMCPList(m.Model)

	return columnStyle.Render(fmt.Sprintf("MCP Manager\n\n%s", content))
}

// renderTwoColumns renders the medium layout with 2 columns
func (m Model) renderTwoColumns() string {
	columnWidth := (m.Width - 6) / 2
	columnHeight := m.Height - 8

	// Column styles
	leftStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1).
		Width(columnWidth).
		Height(columnHeight)

	rightStyle := leftStyle.Copy()

	// Highlight active column
	if m.ActiveColumn == 0 {
		leftStyle = leftStyle.BorderForeground(lipgloss.Color("#7C3AED"))
	} else {
		rightStyle = rightStyle.BorderForeground(lipgloss.Color("#7C3AED"))
	}

	// Column content
	leftContent := fmt.Sprintf("MCPs\n\n%s", components.RenderMCPList(m.Model))
	rightContent := fmt.Sprintf("Status & Details\n\n%s", m.renderStatusAndDetails())

	leftColumn := leftStyle.Render(leftContent)
	rightColumn := rightStyle.Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
}

// renderThreeColumns renders the wide layout with 3 columns
func (m Model) renderThreeColumns() string {
	columnWidth := (m.Width - 8) / 3
	columnHeight := m.Height - 8

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
		if i == m.ActiveColumn {
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
			Width(m.Width).
			Height(m.Height - 8)
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
				isSelected := (mcpIndex == m.SelectedItem)

				// Create base item text (without styling)
				baseText := fmt.Sprintf("%s %s", status, item.Name)

				// Calculate padding needed BEFORE styling
				currentWidth := lipgloss.Width(baseText)
				paddingNeeded := types.COLUMN_WIDTH - currentWidth
				if paddingNeeded < 0 {
					paddingNeeded = 0
				}

				// Apply padding first
				paddedText := baseText + strings.Repeat(" ", paddingNeeded)

				// Then apply styling to padded text
				if isSelected {
					itemStyle := lipgloss.NewStyle().
						Bold(true).
						Background(lipgloss.Color("#7C3AED")).
						Foreground(lipgloss.Color("#FFFFFF"))
					paddedText = itemStyle.Render(paddedText)
				}

				line = append(line, paddedText)
			} else {
				// Empty cell with proper spacing
				line = append(line, strings.Repeat(" ", types.COLUMN_WIDTH))
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
		Width(m.Width).
		Height(m.Height - 8)

	return gridStyle.Render(gridContent)
}

// renderMCPList renders the list of MCPs with selection highlighting
func (m Model) renderMCPList() string {
	if len(m.MCPItems) == 0 {
		return "No MCPs configured"
	}

	var items []string
	for i, item := range m.MCPItems {
		style := lipgloss.NewStyle().Padding(0, 1)

		// Highlight selected item
		if i == m.SelectedItem {
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
func (m Model) renderMCPColumnList(columnMCPs []types.MCPItem, startIdx int) string {
	if len(columnMCPs) == 0 {
		return ""
	}

	var items []string
	for i, item := range columnMCPs {
		actualIdx := startIdx + i
		style := lipgloss.NewStyle().Padding(0, 1)

		// Highlight selected item
		if actualIdx == m.SelectedItem {
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
	if m.SelectedItem >= len(m.MCPItems) {
		return "No MCP selected"
	}

	item := m.MCPItems[m.SelectedItem]

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
	if m.SelectedItem >= len(m.MCPItems) {
		return "No MCP selected"
	}

	item := m.MCPItems[m.SelectedItem]

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
		Width(m.Width)

	var footerText string
	if m.SearchActive {
		searchStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("#7C3AED")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

		cursor := "_"
		modeIndicator := ""

		// Show dual-mode indicator for SearchActiveNavigation
		if m.State == types.SearchActiveNavigation {
			if m.SearchInputActive {
				modeIndicator = " [INPUT MODE]"
			} else {
				modeIndicator = " [NAVIGATION MODE]"
			}
		}

		footerText = fmt.Sprintf("Search: %s%s", searchStyle.Render(m.SearchQuery+cursor), modeIndicator)
	} else if m.SearchQuery != "" {
		// Show search results info when not actively searching but have a query
		filteredMCPs := m.GetFilteredMCPs()
		footerText = fmt.Sprintf("Found %d MCPs matching '%s' • ESC to clear • Terminal: %dx%d",
			len(filteredMCPs), m.SearchQuery, m.Width, m.Height)
	} else {
		footerText = fmt.Sprintf("Terminal: %dx%d • Search: '%s' • Use arrow keys to navigate, Tab or / for search",
			m.Width, m.Height, m.SearchQuery)
	}

	return footerStyle.Render(footerText)
}

// getLayoutName returns the current layout name for display
func (m Model) getLayoutName() string {
	switch m.ColumnCount {
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
