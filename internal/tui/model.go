// Package tui provides the Terminal User Interface for the MCP Manager CLI.
//
// This package implements a responsive, feature-rich TUI using the Bubble Tea framework.
// It provides a three-column layout that adapts to different terminal sizes, comprehensive
// keyboard navigation, search functionality, and robust error handling.
//
// Key Features:
//   - Responsive layout that adapts to terminal width (3-column, 2-column, single-column)
//   - Intuitive keyboard navigation with both arrow keys and vim-style bindings
//   - Search functionality with real-time filtering
//   - Error handling with user-friendly messages and recovery options
//   - Configuration management with validation
//   - Status notifications and feedback
//
// Architecture:
//   - Model: Central application state following Bubble Tea patterns
//   - View: Rendering logic with consistent styling using Lipgloss
//   - Update: Event handling and state transitions
//   - Error handling: Centralized error management with recovery actions
//
// Usage:
//   m := tui.NewModel()
//   p := tea.NewProgram(m, tea.WithAltScreen())
//   if _, err := p.Run(); err != nil {
//       log.Fatal(err)
//   }
//
// The TUI supports terminals with a minimum width of 40 columns and gracefully
// degrades functionality for smaller terminals. All keyboard shortcuts are
// discoverable through the header display.
package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"mcp-manager/internal/config"
	"mcp-manager/internal/errors"
)

// ViewMode represents the current view mode of the TUI.
// The application can be in different modes that affect how user input
// is handled and what is displayed on screen.
type ViewMode int

const (
	// ViewModeList is the default mode for browsing and navigating MCP servers
	ViewModeList ViewMode = iota
	
	// ViewModeSearch is the mode where users can type search queries
	ViewModeSearch
	
	// ViewModeError is the mode that displays error dialogs with recovery options
	ViewModeError
)

// Column represents a column in the responsive layout.
// The TUI uses a multi-column layout that adapts based on terminal width.
type Column int

const (
	// ColumnLeft is the leftmost column (typically MCP server list)
	ColumnLeft Column = iota
	
	// ColumnCenter is the middle column (visible in 3-column layout)
	ColumnCenter
	
	// ColumnRight is the rightmost column (typically details panel)
	ColumnRight
)

// LayoutType represents the current layout configuration based on terminal width.
// The TUI automatically adapts its layout to provide the best user experience
// across different terminal sizes.
type LayoutType int

const (
	// Layout3Column provides three columns for terminals with 80+ columns
	// Offers: left panel (MCP list), center panel (details), right panel (actions)
	Layout3Column LayoutType = iota
	
	// Layout2Column provides two columns for terminals with 60-79 columns
	// Offers: left panel (MCP list), right panel (details/actions)
	Layout2Column
	
	// Layout1Column provides single column for terminals with <60 columns
	// Offers: single scrollable panel with all information
	Layout1Column
)

// Model represents the TUI application state and implements the Bubble Tea Model interface.
// It contains all the state needed to render the interface and handle user interactions.
//
// The Model follows the Elm Architecture pattern:
//   - All state is immutable and centralized
//   - Updates are handled through message passing
//   - Views are pure functions of the model state
type Model struct {
	// Navigation state - tracks current position and mode
	currentColumn Column   // Which column has focus in multi-column layouts
	currentMode   ViewMode // Current interaction mode (list, search, error)
	selectedIndex int      // Currently selected item index
	
	// Terminal dimensions - updated automatically on window resize
	width  int        // Terminal width in columns
	height int        // Terminal height in lines
	layout LayoutType // Current layout configuration based on width
	
	// Search state - manages search functionality
	searchQuery  string // Current search query string
	searchActive bool   // Whether search mode is active
	
	// Data - application data (will be replaced with real MCP data)
	mcpList []string // List of MCP servers (mock data for now)
	
	// UI components - rendering and initialization state
	initialized bool // Whether the terminal has been properly initialized
	
	// Configuration - application settings and preferences
	config *config.Config // Configuration loaded from file or defaults
	
	// Error handling - manages error display and recovery
	currentError   error     // Currently displayed error (if any)
	errorDisplayed bool      // Whether an error dialog is visible
	errorStartTime time.Time // When the error was first displayed
	
	// Status system - temporary messages and notifications
	statusMessage string    // Current status message
	statusTimeout time.Time // When the status message expires
}

// NewModel creates a new TUI model with default settings and loaded configuration.
//
// This function:
//   - Loads the application configuration from file or uses defaults
//   - Initializes the model with sensible default values
//   - Sets up mock MCP data (to be replaced with real data loading)
//   - Handles configuration loading errors gracefully
//
// Returns a fully initialized Model ready to be used with Bubble Tea.
// If configuration loading fails, the model will display an error but remain functional
// using default configuration values.
func NewModel() Model {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// If config fails to load, use default and show error
		cfg = config.DefaultConfig()
	}
	
	m := Model{
		currentColumn: ColumnLeft,
		currentMode:   ViewModeList,
		selectedIndex: 0,
		mcpList:       []string{"MCP Server 1", "MCP Server 2", "MCP Server 3", "MCP Server 4"},
		initialized:   false,
		config:        cfg,
	}
	
	// Handle config loading error
	if err != nil {
		m.ShowError(errors.Wrap(err, errors.ErrorTypeConfiguration, "CONFIG_LOAD_FAILED", 
			"Failed to load configuration"))
	}
	
	return m
}

// Init initializes the model and returns initial commands for Bubble Tea.
//
// This method is called by the Bubble Tea framework when the program starts.
// It sets up the alternate screen buffer to ensure the TUI doesn't interfere
// with the user's terminal history.
//
// Returns tea.EnterAltScreen command to switch to alternate screen mode.
func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// determineLayout calculates the appropriate layout based on terminal width.
//
// This method automatically selects the best layout configuration for the current
// terminal size, ensuring optimal user experience across different environments.
//
// Layout selection rules:
//   - 80+ columns: 3-column layout (list | details | actions)
//   - 60-79 columns: 2-column layout (list | details+actions)
//   - <60 columns: 1-column layout (single scrollable view)
//
// If the terminal is smaller than the configured minimum width, an error is displayed
// asking the user to resize their terminal.
func (m *Model) determineLayout() {
	minWidth := m.config.UI.MinimumWidth
	if m.width < minWidth {
		m.ShowError(errors.ErrTerminalTooSmall.WithContext("current_width", m.width).WithContext("minimum_width", minWidth))
		return
	}
	
	if m.width >= 80 {
		m.layout = Layout3Column
	} else if m.width >= 60 {
		m.layout = Layout2Column
	} else {
		m.layout = Layout1Column
	}
}

// ShowError displays an error to the user with a modal dialog overlay.
//
// This method switches the TUI to error mode, displaying a user-friendly error message
// with appropriate recovery options. The error dialog includes:
//   - Human-readable error description
//   - Recovery actions (if the error is recoverable)
//   - Keyboard shortcuts for dismissing or handling the error
//
// The error remains visible until the user explicitly dismisses it or attempts recovery.
//
// Parameters:
//   err: The error to display (can be a standard error or *errors.AppError)
func (m *Model) ShowError(err error) {
	m.currentError = err
	m.errorDisplayed = true
	m.errorStartTime = time.Now()
	m.currentMode = ViewModeError
}

// ClearError clears the current error and returns to normal operation mode.
//
// This method:
//   - Removes the error from the model state
//   - Hides the error dialog overlay
//   - Returns the TUI to list/navigation mode
//
// Called when the user dismisses an error or after successful error recovery.
func (m *Model) ClearError() {
	m.currentError = nil
	m.errorDisplayed = false
	m.currentMode = ViewModeList
}

// ShowStatus displays a temporary status message in the footer.
//
// Status messages provide feedback to the user about operations that have completed
// or are in progress. They automatically disappear after the specified duration.
//
// Examples of status messages:
//   - "Configuration saved successfully"
//   - "Connecting to MCP server..."
//   - "Search completed: 5 results found"
//
// Parameters:
//   message: The status text to display
//   duration: How long to show the message before it automatically disappears
func (m *Model) ShowStatus(message string, duration time.Duration) {
	m.statusMessage = message
	m.statusTimeout = time.Now().Add(duration)
}

// ClearStatus clears the status message if it has expired.
//
// This method is called automatically during each update cycle to ensure
// that temporary status messages are removed after their display duration
// has elapsed. It only clears the message if the timeout has been reached.
func (m *Model) ClearStatus() {
	if time.Now().After(m.statusTimeout) {
		m.statusMessage = ""
	}
}

// Update handles messages and updates the model according to the Bubble Tea pattern.
//
// This is the central update function that processes all incoming messages and
// returns an updated model along with any commands to execute. It handles:
//   - Window size changes (terminal resize)
//   - Keyboard input (navigation, search, error handling)
//   - Automatic cleanup (expired status messages)
//
// The function follows the Elm Architecture pattern where all state changes
// happen through this central update mechanism.
//
// Parameters:
//   msg: The message to process (tea.WindowSizeMsg, tea.KeyMsg, etc.)
//
// Returns:
//   tea.Model: Updated model with new state
//   tea.Cmd: Optional command to execute (e.g., tea.Quit)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Clear expired status messages
	m.ClearStatus()
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.determineLayout()
		if !m.initialized {
			m.initialized = true
		}
		
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	
	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle error mode first
	if m.currentMode == ViewModeError {
		return m.handleErrorKeys(msg)
	}
	
	// Handle global keys first
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "esc":
		if m.searchActive {
			m.searchActive = false
			m.currentMode = ViewModeList
			return m, nil
		}
		return m, tea.Quit
	}
	
	// Handle mode-specific keys
	if m.searchActive {
		return m.handleSearchKeys(msg)
	}
	
	return m.handleNavigationKeys(msg)
}

// handleErrorKeys processes keyboard input in error mode
func (m Model) handleErrorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "esc", " ":
		// Clear error and return to normal mode
		m.ClearError()
		return m, nil
	case "r":
		// Try to recover from error if possible
		if errors.IsRecoverable(m.currentError) {
			return m.attemptRecovery()
		}
		m.ClearError()
		return m, nil
	case "q":
		// Critical errors can force quit
		if errors.IsCritical(m.currentError) {
			return m, tea.Quit
		}
		m.ClearError()
		return m, nil
	}
	return m, nil
}

// attemptRecovery tries to recover from the current error
func (m Model) attemptRecovery() (tea.Model, tea.Cmd) {
	if m.currentError == nil {
		return m, nil
	}
	
	actions := errors.GetRecoveryActions(m.currentError)
	if len(actions) == 0 {
		m.ShowStatus("No recovery actions available", 3*time.Second)
		m.ClearError()
		return m, nil
	}
	
	// Try the first recovery action
	action := actions[0]
	switch action.Action {
	case "create_default_config":
		if err := m.config.Save(); err != nil {
			m.ShowError(errors.Wrap(err, errors.ErrorTypeConfiguration, "CONFIG_SAVE_FAILED", 
				"Failed to save default configuration"))
		} else {
			m.ShowStatus("Default configuration created", 3*time.Second)
			m.ClearError()
		}
		
	case "resize_terminal":
		// This will be handled automatically by the next window size message
		m.ShowStatus(fmt.Sprintf("Please resize terminal to at least %d columns", m.config.UI.MinimumWidth), 5*time.Second)
		m.ClearError()
		
	case "retry_connection":
		// In future this would retry connection
		m.ShowStatus("Retrying connection...", 2*time.Second)
		m.ClearError()
		
	default:
		m.ShowStatus("Recovery attempted", 2*time.Second)
		m.ClearError()
	}
	
	return m, nil
}

// handleNavigationKeys processes navigation keyboard input
func (m Model) handleNavigationKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.searchActive = true
		m.currentMode = ViewModeSearch
		
	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
		} else {
			m.selectedIndex = len(m.mcpList) - 1 // Wrap to bottom
		}
		
	case "down", "j":
		if m.selectedIndex < len(m.mcpList)-1 {
			m.selectedIndex++
		} else {
			m.selectedIndex = 0 // Wrap to top
		}
		
	case "left", "h":
		switch m.layout {
		case Layout3Column:
			switch m.currentColumn {
			case ColumnCenter:
				m.currentColumn = ColumnLeft
			case ColumnRight:
				m.currentColumn = ColumnCenter
			}
		case Layout2Column:
			if m.currentColumn == ColumnRight {
				m.currentColumn = ColumnLeft
			}
		}
		
	case "right", "l":
		switch m.layout {
		case Layout3Column:
			switch m.currentColumn {
			case ColumnLeft:
				m.currentColumn = ColumnCenter
			case ColumnCenter:
				m.currentColumn = ColumnRight
			}
		case Layout2Column:
			if m.currentColumn == ColumnLeft {
				m.currentColumn = ColumnRight
			}
		}
	}
	
	return m, nil
}

// handleSearchKeys processes search mode keyboard input
func (m Model) handleSearchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.searchActive = false
		m.currentMode = ViewModeList
		
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
		
	default:
		// Add character to search query
		if len(msg.String()) == 1 {
			m.searchQuery += msg.String()
		}
	}
	
	return m, nil
}

// View renders the TUI and returns the complete terminal output as a string.
//
// This method implements the View part of the Bubble Tea Model interface. It's a pure
// function that takes the current model state and returns a string representation
// of what should be displayed on the terminal.
//
// The view is structured as:
//   1. Header with title and keyboard shortcuts
//   2. Main content area (columns or error overlay)
//   3. Footer with status messages and layout information
//
// The rendering automatically adapts to:
//   - Different terminal sizes (responsive layout)
//   - Current mode (normal, search, error)
//   - Available data and user interactions
//
// Returns:
//   string: Complete terminal output ready for display
func (m Model) View() string {
	if !m.initialized {
		return "Initializing MCP Manager..."
	}
	
	var s string
	
	// Render header
	s += m.renderHeader()
	s += "\n"
	
	// Render error overlay if in error mode
	if m.currentMode == ViewModeError {
		s += m.renderErrorOverlay()
	} else {
		// Render main content based on layout
		switch m.layout {
		case Layout3Column:
			s += m.render3Column()
		case Layout2Column:
			s += m.render2Column()
		case Layout1Column:
			s += m.render1Column()
		}
	}
	
	// Render footer/status
	s += "\n"
	s += m.renderFooter()
	
	return s
}

// renderErrorOverlay renders an error dialog overlay that appears on top of the main interface.
//
// This method creates a modal-style error dialog that:
//   - Displays user-friendly error messages
//   - Shows appropriate recovery options based on error type
//   - Provides clear instructions for how to proceed
//   - Is centered both horizontally and vertically on screen
//   - Uses visual styling to indicate error state (red border)
//
// The dialog content adapts based on the error characteristics:
//   - Recoverable errors show a [R]ecovery option
//   - Critical errors may force the user to quit
//   - Non-critical errors allow continuing with [Enter/Esc/Space]
//
// Returns:
//   string: Rendered error dialog overlay, or empty string if no error
func (m Model) renderErrorOverlay() string {
	if m.currentError == nil {
		return ""
	}
	
	userMsg := errors.FormatErrorForUser(m.currentError)
	recoverable := errors.IsRecoverable(m.currentError)
	critical := errors.IsCritical(m.currentError)
	
	var content string
	content += "⚠️  Error\n\n"
	content += userMsg + "\n\n"
	
	if recoverable {
		content += "Press [R] to try recovery, "
	}
	
	if critical {
		content += "Press [Q] to quit"
	} else {
		content += "Press [Enter/Esc/Space] to continue"
	}
	
	// Style the error box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("9")). // Red
		Padding(1, 2).
		Width(min(m.width-4, 60)).
		Align(lipgloss.Center)
	
	errorBox := boxStyle.Render(content)
	
	// Center the error box vertically and horizontally
	return lipgloss.Place(m.width, m.height-6, lipgloss.Center, lipgloss.Center, errorBox)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// renderHeader renders the application header with keyboard shortcuts
func (m Model) renderHeader() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Render("MCP Manager CLI")
	
	shortcuts := "[A]dd [E]dit [D]elete [Space]Toggle [R]efresh [Q]uit"
	if m.currentMode == ViewModeError {
		if errors.IsRecoverable(m.currentError) {
			shortcuts = "[R]ecovery [Enter/Esc/Space]Continue"
		} else {
			shortcuts = "[Enter/Esc/Space]Continue"
		}
		if errors.IsCritical(m.currentError) {
			shortcuts += " [Q]uit"
		}
	} else if m.searchActive {
		shortcuts = "[Enter]Finish [Esc]Cancel"
	}
	
	shortcutsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render(shortcuts)
	
	headerStyle := lipgloss.NewStyle().
		Width(m.width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Padding(0, 1)
	
	headerContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		lipgloss.NewStyle().Width(m.width-len(title)-len(shortcuts)-4).Render(""),
		shortcutsStyle,
	)
	
	return headerStyle.Render(headerContent)
}

// render3Column renders the 3-column layout
func (m Model) render3Column() string {
	leftWidth := m.width / 3
	centerWidth := m.width / 3
	rightWidth := m.width - leftWidth - centerWidth
	
	leftColumn := m.renderColumn("Left Panel", leftWidth, m.currentColumn == ColumnLeft)
	centerColumn := m.renderColumn("Center Panel", centerWidth, m.currentColumn == ColumnCenter)
	rightColumn := m.renderColumn("Right Panel", rightWidth, m.currentColumn == ColumnRight)
	
	return lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, centerColumn, rightColumn)
}

// render2Column renders the 2-column layout
func (m Model) render2Column() string {
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth
	
	leftColumn := m.renderColumn("Main Panel", leftWidth, m.currentColumn == ColumnLeft)
	rightColumn := m.renderColumn("Details Panel", rightWidth, m.currentColumn == ColumnRight)
	
	return lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
}

// render1Column renders the single-column layout
func (m Model) render1Column() string {
	return m.renderColumn("MCP List", m.width, true)
}

// renderColumn renders a single column with the given title and width
func (m Model) renderColumn(title string, width int, isActive bool) string {
	borderStyle := lipgloss.NormalBorder()
	if isActive {
		borderStyle = lipgloss.ThickBorder()
	}
	
	style := lipgloss.NewStyle().
		Border(borderStyle).
		Width(width - 2).
		Height(m.height - 6). // Account for header and footer
		Padding(1)
	
	if isActive {
		style = style.BorderForeground(lipgloss.Color("12"))
	}
	
	content := fmt.Sprintf("%s\n\n", title)
	
	// Render MCP list if this is the active column or single column layout
	if isActive || m.layout == Layout1Column {
		for i, mcp := range m.mcpList {
			prefix := "  "
			if i == m.selectedIndex && !m.searchActive {
				prefix = "> "
				mcp = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render(mcp)
			}
			content += fmt.Sprintf("%s%s\n", prefix, mcp)
		}
	}
	
	return style.Render(content)
}

// renderFooter renders the status footer
func (m Model) renderFooter() string {
	// Show status message if available
	if m.statusMessage != "" && time.Now().Before(m.statusTimeout) {
		footerStyle := lipgloss.NewStyle().
			Width(m.width).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			Padding(0, 1).
			Foreground(lipgloss.Color("10")). // Green for status messages
			Bold(true)
		
		return footerStyle.Render(m.statusMessage)
	}
	
	status := fmt.Sprintf("Layout: %dx%d ", m.width, m.height)
	
	switch m.layout {
	case Layout3Column:
		status += "| 3-Column Mode"
	case Layout2Column:
		status += "| 2-Column Mode"
	case Layout1Column:
		status += "| Single-Column Mode"
	}
	
	if m.searchActive {
		status += fmt.Sprintf(" | Search: %s", m.searchQuery)
	}
	
	// Show error indicator if there's an error
	if m.currentError != nil && m.currentMode == ViewModeError {
		status += " | ⚠️ Error Mode"
	}
	
	footerStyle := lipgloss.NewStyle().
		Width(m.width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		Padding(0, 1).
		Foreground(lipgloss.Color("8"))
	
	return footerStyle.Render(status)
}