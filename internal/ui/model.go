package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// AppState represents the current application state
type AppState int

const (
	MainNavigation AppState = iota
	SearchMode
	ModalActive
)

// Model represents the main application model
type Model struct {
	// Window dimensions
	width  int
	height int

	// Application state
	state AppState

	// Navigation
	activeColumn int
	selectedItem int

	// Search
	searchQuery   string
	searchActive  bool
	searchResults []string

	// Layout
	columns     []Column
	columnCount int

	// MCP inventory (placeholder for future stories)
	mcpItems []MCPItem
}

// MCPItem represents an MCP in the inventory
type MCPItem struct {
	Name    string
	Type    string
	Active  bool
	Command string
}

// Column represents a UI column
type Column struct {
	Title string
	Items []string
	Width int
}

// NewModel creates a new application model
func NewModel() Model {
	return Model{
		state:        MainNavigation,
		activeColumn: 0,
		selectedItem: 0,
		searchQuery:  "",
		searchActive: false,
		columns:      make([]Column, 3),
		columnCount:  3,
		mcpItems: []MCPItem{
			{Name: "context7", Type: "CMD", Active: true, Command: "npx @context7/mcp-server"},
			{Name: "ht-mcp", Type: "SSE", Active: false, Command: "ht-mcp"},
			{Name: "filesystem", Type: "JSON", Active: true, Command: ""},
		},
	}
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return nil
}

// Getter methods for testing

// GetColumnCount returns the current number of columns
func (m Model) GetColumnCount() int {
	return m.columnCount
}

// GetActiveColumn returns the currently active column index
func (m Model) GetActiveColumn() int {
	return m.activeColumn
}

// GetSelectedItem returns the currently selected item index
func (m Model) GetSelectedItem() int {
	return m.selectedItem
}

// GetState returns the current application state
func (m Model) GetState() AppState {
	return m.state
}

// GetSearchQuery returns the current search query
func (m Model) GetSearchQuery() string {
	return m.searchQuery
}

// GetSearchActive returns whether search is currently active
func (m Model) GetSearchActive() bool {
	return m.searchActive
}
