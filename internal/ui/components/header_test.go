package components

import (
	"strings"
	"testing"

	"cc-mcp-manager/internal/testutil"
	"cc-mcp-manager/internal/ui/types"
)

func TestRenderHeader(t *testing.T) {
	tests := []struct {
		name              string
		state             types.AppState
		searchActive      bool
		searchInputActive bool
		expectedShortcuts string
	}{
		{
			name:              "MainNavigation state shows main shortcuts",
			state:             types.MainNavigation,
			expectedShortcuts: "A=Add • D=Delete • E=Edit • /=Search • Tab=Focus Search • ESC=Exit • ↑↓←→=Navigate",
		},
		{
			name:              "SearchMode state shows search shortcuts",
			state:             types.SearchMode,
			expectedShortcuts: "Type to search • Enter=Apply • ESC=Cancel",
		},
		{
			name:              "SearchActiveNavigation with input active shows input shortcuts",
			state:             types.SearchActiveNavigation,
			searchActive:      true,
			searchInputActive: true,
			expectedShortcuts: "Type to search • Tab=Navigate Mode • ↑↓←→=Navigate • Space=Toggle • Enter=Apply • ESC=Cancel",
		},
		{
			name:              "SearchActiveNavigation with input inactive shows navigation shortcuts",
			state:             types.SearchActiveNavigation,
			searchActive:      true,
			searchInputActive: false,
			expectedShortcuts: "Navigate Mode • Tab=Input Mode • ↑↓←→=Navigate • Space=Toggle • Enter=Apply • ESC=Cancel",
		},
		{
			name:              "ModalActive state shows modal shortcuts",
			state:             types.ModalActive,
			expectedShortcuts: "Enter=Confirm • ESC=Cancel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testutil.NewTestModel().
				WithWindowSize(120, 40).
				WithState(tt.state).
				WithSearchActive(tt.searchActive).
				WithSearchInputActive(tt.searchInputActive).
				Build()

			result := RenderHeader(model)

			if !strings.Contains(result, tt.expectedShortcuts) {
				t.Errorf("RenderHeader() shortcuts mismatch\nExpected to contain: %s\nActual result: %s",
					tt.expectedShortcuts, result)
			}

			// Verify header contains title
			if !strings.Contains(result, "MCP Manager v1.0") {
				t.Errorf("RenderHeader() should contain title")
			}
		})
	}
}

func TestRenderHeader_ContextInfo(t *testing.T) {
	tests := []struct {
		name         string
		mcpItems     []types.MCPItem
		columnCount  int
		expectedText string
	}{
		{
			name: "Shows active MCP count and layout",
			mcpItems: []types.MCPItem{
				{Name: "active1", Active: true},
				{Name: "inactive1", Active: false},
				{Name: "active2", Active: true},
			},
			columnCount:  4,
			expectedText: "MCPs: 2/3 Active • Layout: Grid (4-column MCP)",
		},
		{
			name: "Shows no active MCPs",
			mcpItems: []types.MCPItem{
				{Name: "inactive1", Active: false},
				{Name: "inactive2", Active: false},
			},
			columnCount:  2,
			expectedText: "MCPs: 0/2 Active • Layout: Medium",
		},
		{
			name:         "Shows empty inventory",
			mcpItems:     []types.MCPItem{},
			columnCount:  1,
			expectedText: "MCPs: 0/0 Active • Layout: Narrow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testutil.NewTestModel().
				WithWindowSize(120, 40).
				WithState(types.MainNavigation).
				Build()

			// Override MCPItems and ColumnCount
			model.MCPItems = tt.mcpItems
			model.ColumnCount = tt.columnCount

			result := RenderHeader(model)

			if !strings.Contains(result, tt.expectedText) {
				t.Errorf("RenderHeader() context info mismatch\nExpected to contain: %s\nActual result: %s",
					tt.expectedText, result)
			}
		})
	}
}

func TestGetLayoutName(t *testing.T) {
	tests := []struct {
		name         string
		columnCount  int
		expectedName string
	}{
		{
			name:         "1 column returns Narrow",
			columnCount:  1,
			expectedName: "Narrow",
		},
		{
			name:         "2 columns returns Medium",
			columnCount:  2,
			expectedName: "Medium",
		},
		{
			name:         "3 columns returns Wide (3-panel)",
			columnCount:  3,
			expectedName: "Wide (3-panel)",
		},
		{
			name:         "4 columns returns Grid (4-column MCP)",
			columnCount:  4,
			expectedName: "Grid (4-column MCP)",
		},
		{
			name:         "Unknown column count returns Unknown",
			columnCount:  5,
			expectedName: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testutil.NewTestModel().
				WithWindowSize(120, 40).
				Build()

			model.ColumnCount = tt.columnCount

			result := GetLayoutName(model)

			if result != tt.expectedName {
				t.Errorf("GetLayoutName() = %s, expected %s", result, tt.expectedName)
			}
		})
	}
}

func TestRenderHeader_ResponsiveWidth(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "Narrow screen",
			width:  60,
			height: 20,
		},
		{
			name:   "Medium screen",
			width:  100,
			height: 30,
		},
		{
			name:   "Wide screen",
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

			result := RenderHeader(model)

			// Verify header renders without panic
			if result == "" {
				t.Errorf("RenderHeader() should not return empty string")
			}

			// Header should contain all required elements
			requiredElements := []string{
				"MCP Manager v1.0",
				"Layout:",
				"MCPs:",
			}

			for _, element := range requiredElements {
				if !strings.Contains(result, element) {
					t.Errorf("RenderHeader() should contain %s for width %d", element, tt.width)
				}
			}
		})
	}
}

func TestRenderHeader_StateTransitions(t *testing.T) {
	model := testutil.NewTestModel().
		WithWindowSize(120, 40).
		Build()

	// Test that each state produces different shortcut content
	states := []types.AppState{
		types.MainNavigation,
		types.SearchMode,
		types.SearchActiveNavigation,
		types.ModalActive,
	}

	results := make(map[types.AppState]string)

	for _, state := range states {
		model.State = state
		if state == types.SearchActiveNavigation {
			model.SearchInputActive = true
		}
		results[state] = RenderHeader(model)
	}

	// Verify each state produces unique content
	for i, state1 := range states {
		for j, state2 := range states {
			if i != j {
				if results[state1] == results[state2] {
					t.Errorf("States %v and %v should produce different header content", state1, state2)
				}
			}
		}
	}
}
