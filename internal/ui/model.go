package ui

import (
	"strings"
	
	tea "github.com/charmbracelet/bubbletea"
)

// AppState represents the current application state
type AppState int

const (
	MainNavigation AppState = iota
	SearchMode
	SearchActiveNavigation  // Combined search + navigation
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
	searchQuery       string
	searchActive      bool
	searchInputActive bool  // Toggle text input vs navigation
	searchResults     []string

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
		searchQuery:       "",
		searchActive:      false,
		searchInputActive: false,
		columns:      make([]Column, 4),
		columnCount:  4,
mcpItems: []MCPItem{
{Name: "context7", Type: "SSE", Active: true, Command: "npx @context7/mcp-server"},
{Name: "github-mcp", Type: "CMD", Active: true, Command: "github-mcp"},
{Name: "ht-mcp", Type: "CMD", Active: true, Command: "ht-mcp"},
{Name: "filesystem", Type: "CMD", Active: false, Command: "filesystem-mcp"},
{Name: "docker-mcp", Type: "CMD", Active: false, Command: "docker-mcp"},
{Name: "redis-mcp", Type: "CMD", Active: false, Command: "redis-mcp"},
{Name: "jira-mcp", Type: "CMD", Active: false, Command: "jira-mcp"},
{Name: "aws-mcp", Type: "JSON", Active: false, Command: "aws-mcp"},
{Name: "k8s-mcp", Type: "CMD", Active: false, Command: "k8s-mcp"},
{Name: "confluence", Type: "SSE", Active: false, Command: "confluence-mcp"},
{Name: "mongodb", Type: "CMD", Active: false, Command: "mongodb-mcp"},
{Name: "terraform", Type: "CMD", Active: false, Command: "terraform-mcp"},
{Name: "gitlab-mcp", Type: "CMD", Active: false, Command: "gitlab-mcp"},
{Name: "linear-mcp", Type: "CMD", Active: false, Command: "linear-mcp"},
{Name: "postgres", Type: "CMD", Active: false, Command: "postgres-mcp"},
{Name: "elastic", Type: "JSON", Active: false, Command: "elastic-mcp"},
{Name: "bitbucket", Type: "CMD", Active: false, Command: "bitbucket-mcp"},
{Name: "asana-mcp", Type: "CMD", Active: false, Command: "asana-mcp"},
{Name: "notion-mcp", Type: "CMD", Active: false, Command: "notion-mcp"},
{Name: "anthropic", Type: "HTTP", Active: false, Command: "anthropic-mcp"},
{Name: "sourcegraph", Type: "CMD", Active: false, Command: "sourcegraph-mcp"},
{Name: "todoist", Type: "CMD", Active: false, Command: "todoist-mcp"},
{Name: "slack-mcp", Type: "CMD", Active: false, Command: "slack-mcp"},
{Name: "openai-mcp", Type: "HTTP", Active: false, Command: "openai-mcp"},
{Name: "codeberg", Type: "CMD", Active: false, Command: "codeberg-mcp"},
{Name: "calendar", Type: "CMD", Active: false, Command: "calendar-mcp"},
{Name: "discord", Type: "CMD", Active: false, Command: "discord-mcp"},
{Name: "gemini-mcp", Type: "HTTP", Active: false, Command: "gemini-mcp"},
{Name: "gitness", Type: "CMD", Active: false, Command: "gitness-mcp"},
{Name: "email-mcp", Type: "CMD", Active: false, Command: "email-mcp"},
{Name: "teams-mcp", Type: "CMD", Active: false, Command: "teams-mcp"},
{Name: "claude-mcp", Type: "HTTP", Active: false, Command: "claude-mcp"},
{Name: "fossil-mcp", Type: "CMD", Active: false, Command: "fossil-mcp"},
{Name: "browser", Type: "CMD", Active: false, Command: "browser-mcp"},
{Name: "zoom-mcp", Type: "CMD", Active: false, Command: "zoom-mcp"},
},	}
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

// GetSearchInputActive returns whether search input is currently active
func (m Model) GetSearchInputActive() bool {
	return m.searchInputActive
}

// GetFilteredMCPs returns MCPs filtered by search query
func (m Model) GetFilteredMCPs() []MCPItem {
	// If no search query, return all MCPs
	if m.searchQuery == "" {
		return m.mcpItems
	}
	
	// Filter MCPs by search query directly
	var filtered []MCPItem
	query := strings.ToLower(m.searchQuery)
	for _, item := range m.mcpItems {
		if strings.Contains(strings.ToLower(item.Name), query) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
