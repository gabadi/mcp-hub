package ui

import (
	"strings"
	"testing"

	"cc-mcp-manager/internal/ui/types"
)

// Test string constants
const (
	TestLoadingText = "Loading..."
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
				TestLoadingText,
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
			model.Width = tt.width
			model.Height = tt.height
			model.State = tt.state
			if tt.state == types.SearchActiveNavigation {
				model.SearchActive = true
				model.SearchInputActive = true
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
	tests := getLayoutSwitchingTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validateLayoutSwitching(t, tt)
		})
	}
}

func getLayoutSwitchingTestCases() []struct {
	name               string
	width              int
	height             int
	expectedLayout     string
	expectedComponents []string
} {
	return []struct {
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
}

func validateLayoutSwitching(t *testing.T, tt struct {
	name               string
	width              int
	height             int
	expectedLayout     string
	expectedComponents []string
}) {
	model := NewModel()
	model.Width = tt.width
	model.Height = tt.height

	// Update layout based on width (simulating layout service)
	switch {
	case tt.width >= types.WideLayoutMin:
		model.ColumnCount = 4
	case tt.width >= types.MediumLayoutMin:
		model.ColumnCount = 2
	default:
		model.ColumnCount = 1
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
}

func TestView_StateTransitions(t *testing.T) {
	tests := getStateTransitionTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validateStateTransition(t, tt)
		})
	}
}

func getStateTransitionTestCases() []struct {
	name              string
	state             types.AppState
	searchActive      bool
	searchInputActive bool
	searchQuery       string
	expectedShortcuts []string
	expectedFooter    []string
} {
	basicStates := getBasicStateTransitionCases()
	searchStates := getSearchStateTransitionCases()
	
	var allCases []struct {
		name              string
		state             types.AppState
		searchActive      bool
		searchInputActive bool
		searchQuery       string
		expectedShortcuts []string
		expectedFooter    []string
	}
	
	allCases = append(allCases, basicStates...)
	allCases = append(allCases, searchStates...)
	
	return allCases
}

func getBasicStateTransitionCases() []struct {
	name              string
	state             types.AppState
	searchActive      bool
	searchInputActive bool
	searchQuery       string
	expectedShortcuts []string
	expectedFooter    []string
} {
	return []struct {
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
			name:  "ModalActive shows modal shortcuts",
			state: types.ModalActive,
			expectedShortcuts: []string{
				"ESC=Cancel",
			},
		},
	}
}

func getSearchStateTransitionCases() []struct {
	name              string
	state             types.AppState
	searchActive      bool
	searchInputActive bool
	searchQuery       string
	expectedShortcuts []string
	expectedFooter    []string
} {
	return []struct {
		name              string
		state             types.AppState
		searchActive      bool
		searchInputActive bool
		searchQuery       string
		expectedShortcuts []string
		expectedFooter    []string
	}{
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
	}
}

func validateStateTransition(t *testing.T, tt struct {
	name              string
	state             types.AppState
	searchActive      bool
	searchInputActive bool
	searchQuery       string
	expectedShortcuts []string
	expectedFooter    []string
}) {
	model := NewModel()
	model.Width = 120
	model.Height = 40
	model.State = tt.state
	model.SearchActive = tt.searchActive
	model.SearchInputActive = tt.searchInputActive
	model.SearchQuery = tt.searchQuery

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
}

func TestView_ComponentIntegration(t *testing.T) {
	t.Run("Header shows correct MCP count", func(t *testing.T) {
		validateMCPCountInHeader(t)
	})

	t.Run("Footer shows search results count", func(t *testing.T) {
		validateSearchResultsCount(t)
	})

	t.Run("Body shows filtered MCPs", func(t *testing.T) {
		validateFilteredMCPs(t)
	})
}

func validateMCPCountInHeader(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40

	// Set specific MCPs with known active states
	model.MCPItems = []types.MCPItem{
		{Name: "active1", Active: true},
		{Name: "inactive1", Active: false},
		{Name: "active2", Active: true},
	}

	result := model.View()

	if !strings.Contains(result, "MCPs: 2/3 Active") {
		t.Errorf("View() should show correct MCP count in header")
	}
}

func validateSearchResultsCount(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40
	model.SearchQuery = "active"
	model.SearchActive = false // Not actively searching, but has query

	model.MCPItems = []types.MCPItem{
		{Name: "active1", Active: true},
		{Name: "inactive1", Active: false},
		{Name: "active2", Active: true},
	}

	result := model.View()

	if !strings.Contains(result, "Found") || !strings.Contains(result, "matching 'active'") {
		t.Errorf("View() should show search results count in footer\nActual: %s", result)
	}
}

func validateFilteredMCPs(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40
	model.ColumnCount = 4 // Grid layout
	model.SearchQuery = TestPlatformGithub

	model.MCPItems = []types.MCPItem{
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
			model.Width = tt.width
			model.Height = tt.height
			model.ColumnCount = tt.columnCount

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
		model.Width = 0
		model.Height = 0

		result := model.View()

		if result != TestLoadingText {
			t.Errorf("View() should show 'Loading...' for zero dimensions, got: %s", result)
		}
	})

	t.Run("Empty MCP list", func(t *testing.T) {
		model := NewModel()
		model.Width = 120
		model.Height = 40
		model.MCPItems = []types.MCPItem{} // Empty list

		result := model.View()

		if !strings.Contains(result, "MCPs: 0/0 Active") {
			t.Errorf("View() should handle empty MCP list in header")
		}
	})

	t.Run("Very small dimensions", func(t *testing.T) {
		model := NewModel()
		model.Width = 10
		model.Height = 5
		model.ColumnCount = 1

		result := model.View()

		// Should not panic and should render something
		if result == "" {
			t.Errorf("View() should not return empty string for small dimensions")
		}
	})

	t.Run("Large dimensions", func(t *testing.T) {
		model := NewModel()
		model.Width = 200
		model.Height = 80
		model.ColumnCount = 4

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
		model.Width = 120
		model.Height = 40

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

func TestView_UncoveredMethods(t *testing.T) {
	t.Run("renderHeader", func(t *testing.T) {
		validateRenderHeader(t)
	})

	t.Run("renderFourColumns", func(t *testing.T) {
		validateRenderFourColumns(t)
	})

	t.Run("renderMCPColumnList", func(t *testing.T) {
		validateRenderMCPColumnList(t)
	})

	t.Run("renderFooter", func(t *testing.T) {
		validateRenderFooter(t)
	})

	t.Run("getLayoutName", func(t *testing.T) {
		validateGetLayoutName(t)
	})
}

func validateRenderHeader(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40
	model.State = types.MainNavigation
	
	// Test that the full view includes header content
	result := model.View()
	
	if !strings.Contains(result, "MCP Manager v1.0") {
		t.Errorf("renderHeader should include title")
	}
	if !strings.Contains(result, "MCPs:") {
		t.Errorf("renderHeader should include MCP count")
	}
}

func validateRenderFourColumns(t *testing.T) {
	model := NewModel()
	model.Width = 150
	model.Height = 40
	model.ColumnCount = 4
	
	// Add test MCPs
	model.MCPItems = []types.MCPItem{
		{Name: "test-mcp-1", Active: true, Type: "CMD"},
		{Name: "test-mcp-2", Active: false, Type: "SSE"},
		{Name: "test-mcp-3", Active: true, Type: "JSON"},
		{Name: "test-mcp-4", Active: false, Type: "CMD"},
	}
	
	result := model.View()
	
	// Should use 4-column grid layout
	if !strings.Contains(result, "Grid (4-column MCP)") {
		t.Errorf("renderFourColumns should show grid layout")
	}
	
	// Should show all test MCPs
	if !strings.Contains(result, "test-mcp-1") {
		t.Errorf("renderFourColumns should show test-mcp-1")
	}
	if !strings.Contains(result, "test-mcp-2") {
		t.Errorf("renderFourColumns should show test-mcp-2")
	}
}

func validateRenderMCPColumnList(t *testing.T) {
	model := NewModel()
	model.Width = 100
	model.Height = 30
	model.ColumnCount = 3
	
	// Add test MCPs
	model.MCPItems = []types.MCPItem{
		{Name: "column-test-1", Active: true, Type: "CMD"},
		{Name: "column-test-2", Active: false, Type: "SSE"},
	}
	
	result := model.View()
	
	// Should show MCPs in 3-column layout
	if !strings.Contains(result, "column-test-1") {
		t.Errorf("renderMCPColumnList should show column-test-1")
	}
	if !strings.Contains(result, "column-test-2") {
		t.Errorf("renderMCPColumnList should show column-test-2")
	}
}

func validateRenderFooter(t *testing.T) {
	model := NewModel()
	model.Width = 120
	model.Height = 40
	model.State = types.MainNavigation
	
	result := model.View()
	
	// Should include footer content
	if !strings.Contains(result, "📁") {
		t.Errorf("renderFooter should include project context icon")
	}
	if !strings.Contains(result, "Claude CLI:") {
		t.Errorf("renderFooter should include Claude CLI status")
	}
}

func validateGetLayoutName(t *testing.T) {
	tests := getLayoutNameTestCases()
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.Width = 120
			model.Height = 40
			model.ColumnCount = tt.columnCount
			
			result := model.View()
			
			if !strings.Contains(result, tt.expected) {
				t.Errorf("getLayoutName should return %q for %d columns", tt.expected, tt.columnCount)
			}
		})
	}
}

func getLayoutNameTestCases() []struct {
	name        string
	columnCount int
	expected    string
} {
	return []struct {
		name        string
		columnCount int
		expected    string
	}{
		{"Single column", 1, "Narrow"},
		{"Two columns", 2, "Medium"},
		{"Three columns", 3, "Wide"},
		{"Four columns", 4, "Grid (4-column MCP)"},
	}
}

func TestView_LoadingStates(t *testing.T) {
	t.Run("Loading overlay active", func(t *testing.T) {
		model := NewModel()
		model.Width = 120
		model.Height = 40
		model.StartLoadingOverlay(types.LoadingStartup)
		model.LoadingOverlay.Message = "Test loading..."
		
		result := model.View()
		
		if !strings.Contains(result, "Test loading...") {
			t.Errorf("View should show loading overlay message")
		}
	})

	t.Run("Modal active", func(t *testing.T) {
		model := NewModel()
		model.Width = 120
		model.Height = 40
		model.State = types.ModalActive
		model.ActiveModal = types.AddModal
		
		result := model.View()
		
		// Should show modal overlay
		if !strings.Contains(result, "ESC=Cancel") {
			t.Errorf("View should show modal shortcuts when modal is active")
		}
	})

	t.Run("Alert active", func(t *testing.T) {
		model := NewModel()
		model.Width = 120
		model.Height = 40
		model.SuccessMessage = "Test alert message"
		
		result := model.View()
		
		if !strings.Contains(result, "Test alert message") {
			t.Errorf("View should show alert message when alert is active")
		}
	})
}
