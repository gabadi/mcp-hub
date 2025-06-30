package types

import (
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
	Width  int
	Height int

	// Application state
	State AppState

	// Navigation
	ActiveColumn int
	SelectedItem int

	// Search
	SearchQuery       string
	SearchActive      bool
	SearchInputActive bool  // Toggle text input vs navigation
	SearchResults     []string

	// Layout
	Columns     []Column
	ColumnCount int

	// MCP inventory (placeholder for future stories)
	MCPItems []MCPItem
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
		State:        MainNavigation,
		ActiveColumn: 0,
		SelectedItem: 0,
		SearchQuery:       "",
		SearchActive:      false,
		SearchInputActive: false,
		Columns:      make([]Column, 4),
		ColumnCount:  4,
		MCPItems: []MCPItem{
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
		},
	}
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return nil
}