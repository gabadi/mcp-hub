package components

import (
	"strings"
	"testing"

	"cc-mcp-manager/internal/testutil"
	"cc-mcp-manager/internal/ui/types"
)

func TestRenderFooter(t *testing.T) {
	tests := []struct {
		name                string
		searchActive        bool
		searchQuery         string
		searchInputActive   bool
		state               types.AppState
		width               int
		height              int
		expectedContains    []string
		expectedNotContains []string
	}{
		{
			name:         "Search active shows search input with cursor",
			searchActive: true,
			searchQuery:  "test",
			width:        120,
			height:       40,
			expectedContains: []string{
				"Search:",
				"test",
			},
		},
		{
			name:              "Search active navigation with input active shows INPUT MODE",
			searchActive:      true,
			searchQuery:       "query",
			searchInputActive: true,
			state:             types.SearchActiveNavigation,
			width:             120,
			height:            40,
			expectedContains: []string{
				"Search:",
				"query_",
				"[INPUT MODE]",
			},
		},
		{
			name:              "Search active navigation with input inactive shows NAVIGATION MODE",
			searchActive:      true,
			searchQuery:       "query",
			searchInputActive: false,
			state:             types.SearchActiveNavigation,
			width:             120,
			height:            40,
			expectedContains: []string{
				"Search:",
				"query",
				"[NAVIGATION MODE]",
			},
		},
		{
			name:         "Search inactive with query shows search results info",
			searchActive: false,
			searchQuery:  "github",
			width:        120,
			height:       40,
			expectedContains: []string{
				"Found",
				"matching 'github'",
				"ESC to clear",
				"Terminal: 120x40",
			},
		},
		{
			name:         "No search shows default footer",
			searchActive: false,
			searchQuery:  "",
			width:        100,
			height:       30,
			expectedContains: []string{
				"Terminal: 100x30",
				"Use arrow keys to navigate",
				"Tab or / for search",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testutil.NewTestModel().
				WithWindowSize(tt.width, tt.height).
				WithState(tt.state).
				WithSearchQuery(tt.searchQuery).
				WithSearchActive(tt.searchActive).
				WithSearchInputActive(tt.searchInputActive).
				Build()

			// Add mock MCPs for filtering tests
			model.MCPItems = testutil.MockMCPItems()

			result := RenderFooter(model)

			for _, expected := range tt.expectedContains {
				if !strings.Contains(result, expected) {
					t.Errorf("RenderFooter() should contain %q\nActual: %s", expected, result)
				}
			}

			for _, notExpected := range tt.expectedNotContains {
				if strings.Contains(result, notExpected) {
					t.Errorf("RenderFooter() should not contain %q\nActual: %s", notExpected, result)
				}
			}
		})
	}
}

func TestGetFilteredMCPs(t *testing.T) {
	mcpItems := []types.MCPItem{
		{Name: "github-mcp", Type: "CMD", Active: true},
		{Name: "docker-mcp", Type: "CMD", Active: false},
		{Name: "context7", Type: "SSE", Active: true},
		{Name: "filesystem", Type: "CMD", Active: false},
	}

	tests := []struct {
		name          string
		searchQuery   string
		expected      int
		expectedNames []string
	}{
		{
			name:          "Empty query returns all MCPs",
			searchQuery:   "",
			expected:      4,
			expectedNames: []string{"github-mcp", "docker-mcp", "context7", "filesystem"},
		},
		{
			name:          "Case insensitive search for 'mcp' from custom test data",
			searchQuery:   "mcp",
			expected:      2,
			expectedNames: []string{"github-mcp", "docker-mcp"},
		},
		{
			name:          "Search for 'git' matches github",
			searchQuery:   "git",
			expected:      1,
			expectedNames: []string{"github-mcp"},
		},
		{
			name:          "Search for 'context' matches context7",
			searchQuery:   "context",
			expected:      1,
			expectedNames: []string{"context7"},
		},
		{
			name:          "Search for nonexistent returns empty",
			searchQuery:   "nonexistent",
			expected:      0,
			expectedNames: []string{},
		},
		{
			name:          "Case insensitive search",
			searchQuery:   "DOCKER",
			expected:      1,
			expectedNames: []string{"docker-mcp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testutil.NewTestModel().
				WithSearchQuery(tt.searchQuery).
				Build()

			model.MCPItems = mcpItems

			result := GetFilteredMCPs(model)

			if len(result) != tt.expected {
				t.Errorf("GetFilteredMCPs() returned %d items, expected %d", len(result), tt.expected)
			}

			// Check that all expected names are present
			resultNames := make(map[string]bool)
			for _, item := range result {
				resultNames[item.Name] = true
			}

			for _, expectedName := range tt.expectedNames {
				if !resultNames[expectedName] {
					t.Errorf("GetFilteredMCPs() should include %s", expectedName)
				}
			}
		})
	}
}

func TestRenderFooter_SearchResultsCount(t *testing.T) {
	tests := []struct {
		name          string
		mcpItems      []types.MCPItem
		searchQuery   string
		expectedCount string
	}{
		{
			name: "Multiple matches",
			mcpItems: []types.MCPItem{
				{Name: "github-mcp", Active: true},
				{Name: "gitlab-mcp", Active: false},
				{Name: "docker-mcp", Active: false},
			},
			searchQuery:   "mcp",
			expectedCount: "Found 3 MCPs",
		},
		{
			name: "Single match",
			mcpItems: []types.MCPItem{
				{Name: "github-mcp", Active: true},
				{Name: "context7", Active: false},
			},
			searchQuery:   "github",
			expectedCount: "Found 1 MCPs",
		},
		{
			name: "No matches",
			mcpItems: []types.MCPItem{
				{Name: "github-mcp", Active: true},
				{Name: "context7", Active: false},
			},
			searchQuery:   "nonexistent",
			expectedCount: "Found 0 MCPs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testutil.NewTestModel().
				WithWindowSize(120, 40).
				WithSearchQuery(tt.searchQuery).
				WithSearchActive(false). // Not actively searching, just has query
				Build()

			model.MCPItems = tt.mcpItems

			result := RenderFooter(model)

			if !strings.Contains(result, tt.expectedCount) {
				t.Errorf("RenderFooter() should contain %q\nActual: %s", tt.expectedCount, result)
			}
		})
	}
}

func TestRenderFooter_ResponsiveTerminalInfo(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "Small terminal",
			width:  60,
			height: 20,
		},
		{
			name:   "Medium terminal",
			width:  100,
			height: 30,
		},
		{
			name:   "Large terminal",
			width:  150,
			height: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testutil.NewTestModel().
				WithWindowSize(tt.width, tt.height).
				WithState(types.MainNavigation).
				Build()

			result := RenderFooter(model)

			expectedTerminalInfo := "Terminal: " + strings.ReplaceAll(strings.ReplaceAll("WIDTHxHEIGHT", "WIDTH", "60"), "HEIGHT", "20")
			if tt.width == 60 && tt.height == 20 {
				expectedTerminalInfo = "Terminal: 60x20"
			} else if tt.width == 100 && tt.height == 30 {
				expectedTerminalInfo = "Terminal: 100x30"
			} else if tt.width == 150 && tt.height == 50 {
				expectedTerminalInfo = "Terminal: 150x50"
			}

			if !strings.Contains(result, expectedTerminalInfo) {
				t.Errorf("RenderFooter() should contain terminal info for %dx%d\nActual: %s",
					tt.width, tt.height, result)
			}
		})
	}
}

func TestRenderFooter_EdgeCases(t *testing.T) {
	t.Run("Empty MCP list with search", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithSearchQuery("test").
			WithSearchActive(false).
			Build()

		model.MCPItems = []types.MCPItem{} // Empty list

		result := RenderFooter(model)

		if !strings.Contains(result, "Found 0 MCPs") {
			t.Errorf("RenderFooter() should handle empty MCP list gracefully")
		}
	})

	t.Run("Very long search query", func(t *testing.T) {
		longQuery := "verylongsearchquerythatmightcauseissueswithrendering"
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithSearchQuery(longQuery).
			WithSearchActive(true).
			Build()

		result := RenderFooter(model)

		if !strings.Contains(result, longQuery) {
			t.Errorf("RenderFooter() should handle long search queries")
		}
	})

	t.Run("Zero terminal dimensions", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(0, 0).
			Build()

		result := RenderFooter(model)

		if !strings.Contains(result, "Terminal: 0x0") {
			t.Errorf("RenderFooter() should handle zero dimensions gracefully")
		}
	})
}
