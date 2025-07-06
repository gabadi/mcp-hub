package components

import (
	"strings"
	"testing"
	"time"

	"mcp-hub/internal/testutil"
	"mcp-hub/internal/ui/types"
)

func TestRenderFooter(t *testing.T) {
	tests := getFooterTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := buildFooterTestModel(tt)
			result := RenderFooter(model)
			assertFooterContains(t, result, tt.expectedContains)
			assertFooterNotContains(t, result, tt.expectedNotContains)
		})
	}
}

type footerTestCase struct {
	name                string
	searchActive        bool
	searchQuery         string
	searchInputActive   bool
	state               types.AppState
	width               int
	height              int
	expectedContains    []string
	expectedNotContains []string
}

func getFooterTestCases() []footerTestCase {
	cases := []footerTestCase{}
	cases = append(cases, getSearchActiveTestCases()...)
	cases = append(cases, getSearchInactiveTestCases()...)
	cases = append(cases, getDefaultFooterTestCases()...)
	return cases
}

func getSearchActiveTestCases() []footerTestCase {
	return []footerTestCase{
		{
			name:             "Search active shows search input with cursor",
			searchActive:     true,
			searchQuery:      "test",
			width:            120,
			height:           40,
			expectedContains: []string{"Search:", "test"},
		},
		{
			name:              "Search active navigation with input active shows INPUT MODE",
			searchActive:      true,
			searchQuery:       "query",
			searchInputActive: true,
			state:             types.SearchActiveNavigation,
			width:             120,
			height:            40,
			expectedContains:  []string{"Search:", "query_", "[INPUT MODE]"},
		},
		{
			name:              "Search active navigation with input inactive shows NAVIGATION MODE",
			searchActive:      true,
			searchQuery:       "query",
			searchInputActive: false,
			state:             types.SearchActiveNavigation,
			width:             120,
			height:            40,
			expectedContains:  []string{"Search:", "query", "[NAVIGATION MODE]"},
		},
	}
}

func getSearchInactiveTestCases() []footerTestCase {
	return []footerTestCase{
		{
			name:             "Search inactive with query shows search results info",
			searchActive:     false,
			searchQuery:      "github",
			width:            120,
			height:           40,
			expectedContains: []string{"Found", "matching 'github'", "ESC to clear", "Terminal: 120x40"},
		},
	}
}

func getDefaultFooterTestCases() []footerTestCase {
	return []footerTestCase{
		{
			name:             "No search shows default footer with project context",
			searchActive:     false,
			searchQuery:      "",
			width:            100,
			height:           30,
			expectedContains: []string{"📁", "MCPs", "R=Retry Claude"},
		},
	}
}

func buildFooterTestModel(tt footerTestCase) types.Model {
	model := testutil.NewTestModel().
		WithWindowSize(tt.width, tt.height).
		WithState(tt.state).
		WithSearchQuery(tt.searchQuery).
		WithSearchActive(tt.searchActive).
		WithSearchInputActive(tt.searchInputActive).
		Build()

	// Add mock MCPs for filtering tests
	model.MCPItems = testutil.MockMCPItems()
	return model
}

func assertFooterContains(t *testing.T, result string, expectedContains []string) {
	for _, expected := range expectedContains {
		if !strings.Contains(result, expected) {
			t.Errorf("RenderFooter() should contain %q\nActual: %s", expected, result)
		}
	}
}

func assertFooterNotContains(t *testing.T, result string, expectedNotContains []string) {
	for _, notExpected := range expectedNotContains {
		if strings.Contains(result, notExpected) {
			t.Errorf("RenderFooter() should not contain %q\nActual: %s", notExpected, result)
		}
	}
}

func TestGetFilteredMCPs(t *testing.T) {
	mcpItems := getTestMCPItems()
	tests := getFilteredMCPTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := buildFilteredMCPTestModel(tt.searchQuery, mcpItems)
			result := GetFilteredMCPs(model)
			assertFilteredMCPResults(t, result, tt.expected, tt.expectedNames)
		})
	}
}

type filteredMCPTestCase struct {
	name          string
	searchQuery   string
	expected      int
	expectedNames []string
}

func getTestMCPItems() []types.MCPItem {
	return []types.MCPItem{
		{Name: "github-mcp", Type: "CMD", Active: true},
		{Name: "docker-mcp", Type: "CMD", Active: false},
		{Name: "context7", Type: "SSE", Active: true},
		{Name: "filesystem", Type: "CMD", Active: false},
	}
}

func getFilteredMCPTestCases() []filteredMCPTestCase {
	return []filteredMCPTestCase{
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
}

func buildFilteredMCPTestModel(searchQuery string, mcpItems []types.MCPItem) types.Model {
	model := testutil.NewTestModel().
		WithSearchQuery(searchQuery).
		Build()
	model.MCPItems = mcpItems
	return model
}

func assertFilteredMCPResults(t *testing.T, result []types.MCPItem, expected int, expectedNames []string) {
	if len(result) != expected {
		t.Errorf("GetFilteredMCPs() returned %d items, expected %d", len(result), expected)
	}

	resultNames := make(map[string]bool)
	for _, item := range result {
		resultNames[item.Name] = true
	}

	for _, expectedName := range expectedNames {
		if !resultNames[expectedName] {
			t.Errorf("GetFilteredMCPs() should include %s", expectedName)
		}
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

func TestRenderFooter_ResponsiveProjectContext(t *testing.T) {
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

			// Should contain project context elements regardless of terminal size
			expectedElements := []string{"📁", "MCPs", "R="}
			for _, element := range expectedElements {
				if !strings.Contains(result, element) {
					t.Errorf("RenderFooter() should contain project context element %q for %dx%d\nActual: %s",
						element, tt.width, tt.height, result)
				}
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

		if !strings.Contains(result, "📁") {
			t.Errorf("RenderFooter() should handle zero dimensions gracefully and show project context")
		}
	})
}

// Epic 2 Story 5 Tests - Project Context Display

func TestRenderFooter_ProjectContext(t *testing.T) {
	tests := getProjectContextTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := buildProjectContextTestModel(tt.mcpItems, tt.projectContext)
			result := RenderFooter(model)
			assertFooterContains(t, result, tt.expectedContains)
		})
	}
}

type projectContextTestCase struct {
	name             string
	mcpItems         []types.MCPItem
	projectContext   types.ProjectContext
	expectedContains []string
}

func getProjectContextTestCases() []projectContextTestCase {
	return []projectContextTestCase{
		{
			name: "Shows project context with active MCPs",
			mcpItems: []types.MCPItem{
				{Name: "github-mcp", Active: true},
				{Name: "docker-mcp", Active: false},
				{Name: "context7", Active: true},
			},
			projectContext: types.ProjectContext{
				DisplayPath:    "~/test-project",
				ActiveMCPs:     2,
				TotalMCPs:      3,
				SyncStatusText: "In Sync",
			},
			expectedContains: []string{
				"📁 ~/test-project",
				"2/3 MCPs",
				"In Sync",
			},
		},
		{
			name: "Shows out of sync status",
			mcpItems: []types.MCPItem{
				{Name: "github-mcp", Active: true},
				{Name: "docker-mcp", Active: false},
			},
			projectContext: types.ProjectContext{
				DisplayPath:    "/long/path/to/project",
				ActiveMCPs:     1,
				TotalMCPs:      2,
				SyncStatusText: "Out of Sync",
			},
			expectedContains: []string{
				"📁 /long/path/to/project",
				"1/2 MCPs",
				"Out of Sync",
			},
		},
		{
			name: "Shows error status",
			mcpItems: []types.MCPItem{
				{Name: "github-mcp", Active: false},
			},
			projectContext: types.ProjectContext{
				DisplayPath:    "~/project",
				ActiveMCPs:     0,
				TotalMCPs:      1,
				SyncStatusText: "Error",
			},
			expectedContains: []string{
				"📁 ~/project",
				"0/1 MCPs",
				"Error",
			},
		},
	}
}

func buildProjectContextTestModel(mcpItems []types.MCPItem, projectContext types.ProjectContext) types.Model {
	model := testutil.NewTestModel().
		WithWindowSize(120, 40).
		WithState(types.MainNavigation).
		Build()

	model.MCPItems = mcpItems
	model.ProjectContext = projectContext
	return model
}

func TestRenderFooter_ProjectContextWithSyncTime(t *testing.T) {
	tests := []struct {
		name         string
		lastSyncTime time.Time
		expectedText string
	}{
		{
			name:         "Recent sync (30 seconds ago)",
			lastSyncTime: time.Now().Add(-30 * time.Second),
			expectedText: "Last sync: 30s ago",
		},
		{
			name:         "Recent sync (2 minutes ago)",
			lastSyncTime: time.Now().Add(-2 * time.Minute),
			expectedText: "Last sync: 2m ago",
		},
		{
			name:         "Old sync (2 hours ago)",
			lastSyncTime: time.Now().Add(-2 * time.Hour),
			expectedText: "Last sync:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testutil.NewTestModel().
				WithWindowSize(120, 40).
				WithState(types.MainNavigation).
				Build()

			model.ProjectContext = types.ProjectContext{
				DisplayPath:    "~/test",
				ActiveMCPs:     1,
				TotalMCPs:      2,
				SyncStatusText: "In Sync",
				LastSyncTime:   tt.lastSyncTime,
			}

			result := RenderFooter(model)

			if !strings.Contains(result, tt.expectedText) {
				t.Errorf("RenderFooter() should contain %q for sync time\nActual: %s", tt.expectedText, result)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "Seconds",
			duration: 30 * time.Second,
			expected: "30s",
		},
		{
			name:     "Minutes",
			duration: 2 * time.Minute,
			expected: "2m",
		},
		{
			name:     "Hours",
			duration: 3 * time.Hour,
			expected: "3h",
		},
		{
			name:     "Zero duration",
			duration: 0,
			expected: "0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("formatDuration(%v) = %s, expected %s", tt.duration, result, tt.expected)
			}
		})
	}
}

func TestRenderFooter_ProjectContextPriority(t *testing.T) {
	// Test that project context is shown when not in search mode
	t.Run("Project context shown in main navigation", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithState(types.MainNavigation).
			WithSearchActive(false).
			WithSearchQuery("").
			Build()

		model.ProjectContext = types.ProjectContext{
			DisplayPath:    "~/project",
			ActiveMCPs:     1,
			TotalMCPs:      2,
			SyncStatusText: "In Sync",
		}

		result := RenderFooter(model)

		expectedElements := []string{
			"📁 ~/project",
			"1/2 MCPs",
			"In Sync",
		}

		for _, expected := range expectedElements {
			if !strings.Contains(result, expected) {
				t.Errorf("RenderFooter() should contain project context element %q\nActual: %s", expected, result)
			}
		}
	})

	// Test that search takes priority over project context
	t.Run("Search takes priority over project context", func(t *testing.T) {
		model := testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithState(types.MainNavigation).
			WithSearchActive(true).
			WithSearchQuery("test").
			Build()

		model.ProjectContext = types.ProjectContext{
			DisplayPath:    "~/project",
			ActiveMCPs:     1,
			TotalMCPs:      2,
			SyncStatusText: "In Sync",
		}

		result := RenderFooter(model)

		// Should show search, not project context
		if !strings.Contains(result, "Search:") {
			t.Errorf("RenderFooter() should show search when active")
		}
		if strings.Contains(result, "📁") {
			t.Errorf("RenderFooter() should not show project context when search is active\nActual: %s", result)
		}
	})
}
