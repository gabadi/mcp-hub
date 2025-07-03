package ui

import (
	"strings"
	"testing"

	"cc-mcp-manager/internal/ui/types"
)

func TestView_MainComposition(t *testing.T) {
	tests := []struct {
		name               string
		width              int
		height             int
		state              types.AppState
		expectedComponents []string
	}{
		{
			name:   "Basic view composition includes all components",
			width:  120,
			height: 40,
			state:  types.MainNavigation,
			expectedComponents: []string{
				"MCP Manager v1.0", // Header
				"MCP Manager",      // Body
				"📁",                // Footer project context
			},
		},
		{
			name:   "Loading state shows loading message",
			width:  0,
			height: 0,
			state:  types.MainNavigation,
			expectedComponents: []string{
				"Loading...",
			},
		},
		{
			name:   "Search active state shows search components",
			width:  120,
			height: 40,
			state:  types.SearchActiveNavigation,
			expectedComponents: []string{
				"Type to search", // Header shortcuts
				"Search:",        // Footer
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.Model.Width = tt.width
			model.Model.Height = tt.height
			model.Model.State = tt.state
			if tt.state == types.SearchActiveNavigation {
				model.Model.SearchActive = true
				model.Model.SearchInputActive = true
			}

			result := model.View()

			for _, component := range tt.expectedComponents {
				if !strings.Contains(result, component) {
					t.Errorf("View() should contain %q\nActual: %s", component, result)
				}
			}
		})
	}
}

func TestView_LayoutSwitching(t *testing.T) {
	tests := []struct {
		name               string
		width              int
		height             int
		expectedLayout     string
		expectedComponents []string
	}{
		{
			name:           "Narrow layout (1 column)",
			width:          60,
			height:         30,
			expectedLayout: "Narrow",
			expectedComponents: []string{
				"MCP Manager", // Single column title
			},
		},
		{
			name:           "Medium layout (2 columns)",
			width:          100,
			height:         30,
			expectedLayout: "Medium",
			expectedComponents: []string{
				"MCPs",             // Left column
				"Status & Details", // Right column
			},
		},
		{
			name:           "Wide layout (4 columns/grid)",
			width:          150,
			height:         40,
			expectedLayout: "Grid (4-column MCP)",
			expectedComponents: []string{
				"MCP Inventory", // Grid header
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.Model.Width = tt.width
			model.Model.Height = tt.height

			// Update layout based on width (simulating layout service)
			if tt.width >= types.WIDE_LAYOUT_MIN {
				model.Model.ColumnCount = 4
			} else if tt.width >= types.MEDIUM_LAYOUT_MIN {
				model.Model.ColumnCount = 2
			} else {
				model.Model.ColumnCount = 1
			}

			result := model.View()

			// Check layout name in header
			if !strings.Contains(result, tt.expectedLayout) {
				t.Errorf("View() should contain layout name %q\nActual: %s", tt.expectedLayout, result)
			}

			// Check layout-specific components
			for _, component := range tt.expectedComponents {
				if !strings.Contains(result, component) {
					t.Errorf("View() should contain layout component %q\nActual: %s", component, result)
				}
			}
		})
	}
}

func TestView_StateTransitions(t *testing.T) {
	tests := []struct {
		name              string
		state             types.AppState
		searchActive      bool
		searchInputActive bool
		searchQuery       string
		expectedShortcuts []string
		expectedFooter    []string
	}{
		{
			name:  "MainNavigation shows main shortcuts",
			state: types.MainNavigation,
			expectedShortcuts: []string{
				"A=Add", "D=Delete", "E=Edit", "/=Search", "ESC=Exit",
			},
			expectedFooter: []string{
				"Claude CLI:",
			},
		},
		{
			name:  "SearchMode shows search shortcuts",
			state: types.SearchMode,
			expectedShortcuts: []string{
				"Type to search", "Enter=Apply", "ESC=Cancel",
			},
		},
		{
			name:              "SearchActiveNavigation with input active",
			state:             types.SearchActiveNavigation,
			searchActive:      true,
			searchInputActive: true,
			searchQuery:       "test",
			expectedShortcuts: []string{
				"Type to search", "Tab=Navigate Mode",
			},
			expectedFooter: []string{
				"Search:", "test_", "[INPUT MODE]",
			},
		},
		{
			name:              "SearchActiveNavigation with input inactive",
			state:             types.SearchActiveNavigation,
			searchActive:      true,
			searchInputActive: false,
			searchQuery:       "test",
			expectedShortcuts: []string{
				"Navigate Mode", "Tab=Input Mode",
			},
			expectedFooter: []string{
				"Search:", "test", "[NAVIGATION MODE]",
			},
		},
		{
			name:  "ModalActive shows modal shortcuts",
			state: types.ModalActive,
			expectedShortcuts: []string{
				"ESC=Cancel",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.Model.Width = 120
			model.Model.Height = 40
			model.Model.State = tt.state
			model.Model.SearchActive = tt.searchActive
			model.Model.SearchInputActive = tt.searchInputActive
			model.Model.SearchQuery = tt.searchQuery

			result := model.View()

			// Check shortcuts in header
			for _, shortcut := range tt.expectedShortcuts {
				if !strings.Contains(result, shortcut) {
					t.Errorf("View() should contain shortcut %q\nActual: %s", shortcut, result)
				}
			}

			// Check footer content
			for _, footerItem := range tt.expectedFooter {
				if !strings.Contains(result, footerItem) {
					t.Errorf("View() should contain footer item %q\nActual: %s", footerItem, result)
				}
			}
		})
	}
}

func TestView_ComponentIntegration(t *testing.T) {
	t.Run("Header shows correct MCP count", func(t *testing.T) {
		model := NewModel()
		model.Model.Width = 120
		model.Model.Height = 40

		// Set specific MCPs with known active states
		model.Model.MCPItems = []types.MCPItem{
			{Name: "active1", Active: true},
			{Name: "inactive1", Active: false},
			{Name: "active2", Active: true},
		}

		result := model.View()

		if !strings.Contains(result, "MCPs: 2/3 Active") {
			t.Errorf("View() should show correct MCP count in header")
		}
	})

	t.Run("Footer shows search results count", func(t *testing.T) {
		model := NewModel()
		model.Model.Width = 120
		model.Model.Height = 40
		model.Model.SearchQuery = "active"
		model.Model.SearchActive = false // Not actively searching, but has query

		model.Model.MCPItems = []types.MCPItem{
			{Name: "active1", Active: true},
			{Name: "inactive1", Active: false},
			{Name: "active2", Active: true},
		}

		result := model.View()

		if !strings.Contains(result, "Found") || !strings.Contains(result, "matching 'active'") {
			t.Errorf("View() should show search results count in footer\nActual: %s", result)
		}
	})

	t.Run("Body shows filtered MCPs", func(t *testing.T) {
		model := NewModel()
		model.Model.Width = 120
		model.Model.Height = 40
		model.Model.ColumnCount = 4 // Grid layout
		model.Model.SearchQuery = "github"

		model.Model.MCPItems = []types.MCPItem{
			{Name: "github-mcp", Active: true},
			{Name: "docker-mcp", Active: false},
			{Name: "github-api", Active: true},
		}

		result := model.View()

		// Should show filtered items
		if !strings.Contains(result, "github-mcp") {
			t.Errorf("View() should show filtered MCP 'github-mcp'")
		}
		if !strings.Contains(result, "github-api") {
			t.Errorf("View() should show filtered MCP 'github-api'")
		}
		// Should not show unfiltered items
		if strings.Contains(result, "docker-mcp") {
			t.Errorf("View() should not show unfiltered MCP 'docker-mcp'")
		}
	})
}

func TestView_ResponsiveLayout(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		height         int
		columnCount    int
		expectedLayout string
	}{
		{
			name:           "Breakpoint 1: Narrow layout",
			width:          70,
			height:         25,
			columnCount:    1,
			expectedLayout: "Narrow",
		},
		{
			name:           "Breakpoint 2: Medium layout",
			width:          90,
			height:         30,
			columnCount:    2,
			expectedLayout: "Medium",
		},
		{
			name:           "Breakpoint 3: Wide layout",
			width:          140,
			height:         45,
			columnCount:    4,
			expectedLayout: "Grid (4-column MCP)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.Model.Width = tt.width
			model.Model.Height = tt.height
			model.Model.ColumnCount = tt.columnCount

			result := model.View()

			if !strings.Contains(result, tt.expectedLayout) {
				t.Errorf("View() should show layout %q for %dx%d",
					tt.expectedLayout, tt.width, tt.height)
			}
		})
	}
}

func TestView_EdgeCases(t *testing.T) {
	t.Run("Zero dimensions shows loading", func(t *testing.T) {
		model := NewModel()
		model.Model.Width = 0
		model.Model.Height = 0

		result := model.View()

		if result != "Loading..." {
			t.Errorf("View() should show 'Loading...' for zero dimensions, got: %s", result)
		}
	})

	t.Run("Empty MCP list", func(t *testing.T) {
		model := NewModel()
		model.Model.Width = 120
		model.Model.Height = 40
		model.Model.MCPItems = []types.MCPItem{} // Empty list

		result := model.View()

		if !strings.Contains(result, "MCPs: 0/0 Active") {
			t.Errorf("View() should handle empty MCP list in header")
		}
	})

	t.Run("Very small dimensions", func(t *testing.T) {
		model := NewModel()
		model.Model.Width = 10
		model.Model.Height = 5
		model.Model.ColumnCount = 1

		result := model.View()

		// Should not panic and should render something
		if result == "" {
			t.Errorf("View() should not return empty string for small dimensions")
		}
	})

	t.Run("Large dimensions", func(t *testing.T) {
		model := NewModel()
		model.Model.Width = 200
		model.Model.Height = 80
		model.Model.ColumnCount = 4

		result := model.View()

		// Should handle large dimensions gracefully
		if !strings.Contains(result, "MCP Manager v1.0") {
			t.Errorf("View() should render properly for large dimensions")
		}
	})
}

func TestView_VerticalComposition(t *testing.T) {
	t.Run("Components appear in correct order", func(t *testing.T) {
		model := NewModel()
		model.Model.Width = 120
		model.Model.Height = 40

		result := model.View()

		// Find positions of components
		headerPos := strings.Index(result, "MCP Manager v1.0")
		bodyPos := strings.Index(result, "Debug: MCPs:")
		footerPos := strings.Index(result, "📁") // Look for project context icon

		if headerPos == -1 {
			t.Errorf("View() should contain header")
		}
		if bodyPos == -1 {
			t.Errorf("View() should contain body")
		}
		if footerPos == -1 {
			t.Errorf("View() should contain footer with project context")
		}

		// Check order: header should come before body, body before footer
		if headerPos > bodyPos {
			t.Errorf("Header should come before body")
		}
		if bodyPos > footerPos {
			t.Errorf("Body should come before footer")
		}
	})
}
