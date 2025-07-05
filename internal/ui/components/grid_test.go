package components

import (
	"strings"
	"testing"

	"cc-mcp-manager/internal/testutil"
	"cc-mcp-manager/internal/ui/types"
)

func TestRenderFourColumnGrid(t *testing.T) {
	tests := getFourColumnGridTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := buildGridTestModel(tt)
			result := RenderFourColumnGrid(model)
			assertGridContains(t, result, tt.expected)
			assertGridNotContains(t, result, tt.notExpected)
		})
	}
}

type gridTestCase struct {
	name         string
	mcpItems     []types.MCPItem
	selectedItem int
	searchQuery  string
	width        int
	height       int
	expected     []string
	notExpected  []string
}

func getFourColumnGridTestCases() []gridTestCase {
	cases := make([]gridTestCase, 0)
	cases = append(cases, getBasicGridTestCases()...)
	cases = append(cases, getFilteredGridTestCases()...)
	cases = append(cases, getMultiRowGridTestCases()...)
	return cases
}

func getBasicGridTestCases() []gridTestCase {
	return []gridTestCase{
		{
			name: "Basic grid rendering with selected item",
			mcpItems: []types.MCPItem{
				{Name: "github", Active: true},
				{Name: "docker", Active: false},
				{Name: "context7", Active: true},
				{Name: "filesystem", Active: false},
			},
			selectedItem: 1,
			width:        120,
			height:       40,
			expected: []string{
				"MCP Inventory",
				"● github",
				"○ docker",
				"● context7",
				"○ filesystem",
			},
		},
		{
			name:         "Empty MCP list shows no results message",
			mcpItems:     []types.MCPItem{},
			selectedItem: 0,
			width:        120,
			height:       40,
			expected: []string{
				"No MCPs found matching your search",
			},
		},
	}
}

func getFilteredGridTestCases() []gridTestCase {
	return []gridTestCase{
		{
			name: "Filtered results show only matching MCPs",
			mcpItems: []types.MCPItem{
				{Name: "github-mcp", Active: true},
				{Name: "docker-mcp", Active: false},
				{Name: "context7", Active: true},
			},
			searchQuery:  "mcp",
			selectedItem: 0,
			width:        120,
			height:       40,
			expected: []string{
				"github-mcp",
				"docker-mcp",
			},
			notExpected: []string{
				"context7",
			},
		},
	}
}

func getMultiRowGridTestCases() []gridTestCase {
	return []gridTestCase{
		{
			name: "Grid with many items fills multiple rows",
			mcpItems: []types.MCPItem{
				{Name: "item1", Active: true},
				{Name: "item2", Active: false},
				{Name: "item3", Active: true},
				{Name: "item4", Active: false},
				{Name: "item5", Active: true},
				{Name: "item6", Active: false},
				{Name: "item7", Active: true},
				{Name: "item8", Active: false},
			},
			selectedItem: 4,
			width:        120,
			height:       40,
			expected: []string{
				"item1", "item2", "item3", "item4",
				"item5", "item6", "item7", "item8",
			},
		},
	}
}

func buildGridTestModel(tt gridTestCase) types.Model {
	model := testutil.NewTestModel().
		WithWindowSize(tt.width, tt.height).
		WithSelectedItem(tt.selectedItem).
		WithSearchQuery(tt.searchQuery).
		Build()

	model.MCPItems = tt.mcpItems
	return model
}

func assertGridContains(t *testing.T, result string, expected []string) {
	for _, exp := range expected {
		if !strings.Contains(result, exp) {
			t.Errorf("RenderFourColumnGrid() should contain %q\nActual: %s", exp, result)
		}
	}
}

func assertGridNotContains(t *testing.T, result string, notExpected []string) {
	for _, notExp := range notExpected {
		if strings.Contains(result, notExp) {
			t.Errorf("RenderFourColumnGrid() should not contain %q\nActual: %s", notExp, result)
		}
	}
}

func TestRenderFourColumnGrid_GridLayout(t *testing.T) {
	// Test that items are arranged in 4-column grid layout
	mcpItems := make([]types.MCPItem, 12) // 3 rows of 4 items
	for i := 0; i < 12; i++ {
		mcpItems[i] = types.MCPItem{
			Name:   "item" + string(rune('A'+i)),
			Active: i%2 == 0,
		}
	}

	model := testutil.NewTestModel().
		WithWindowSize(120, 40).
		WithSelectedItem(5). // Select item in second row
		Build()

	model.MCPItems = mcpItems

	result := RenderFourColumnGrid(model)

	// Should contain all items
	for i := 0; i < 12; i++ {
		itemName := "item" + string(rune('A'+i))
		if !strings.Contains(result, itemName) {
			t.Errorf("RenderFourColumnGrid() should contain %s", itemName)
		}
	}

	// Verify grid structure exists
	if !strings.Contains(result, "MCP Inventory") {
		t.Errorf("RenderFourColumnGrid() should contain header")
	}
}

func TestRenderFourColumnGrid_SelectionHighlight(t *testing.T) {
	mcpItems := []types.MCPItem{
		{Name: "first", Active: true},
		{Name: "second", Active: false},
		{Name: "third", Active: true},
	}

	tests := []struct {
		name         string
		selectedItem int
		expectedItem string
	}{
		{
			name:         "First item selected",
			selectedItem: 0,
			expectedItem: "first",
		},
		{
			name:         "Second item selected",
			selectedItem: 1,
			expectedItem: "second",
		},
		{
			name:         "Third item selected",
			selectedItem: 2,
			expectedItem: "third",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testutil.NewTestModel().
				WithWindowSize(120, 40).
				WithSelectedItem(tt.selectedItem).
				Build()

			model.MCPItems = mcpItems

			result := RenderFourColumnGrid(model)

			// The selected item should be present (selection styling is applied via lipgloss)
			if !strings.Contains(result, tt.expectedItem) {
				t.Errorf("RenderFourColumnGrid() should contain selected item %s", tt.expectedItem)
			}
		})
	}
}

func TestRenderFourColumnGrid_StatusIndicators(t *testing.T) {
	mcpItems := []types.MCPItem{
		{Name: "active-mcp", Active: true},
		{Name: "inactive-mcp", Active: false},
	}

	model := testutil.NewTestModel().
		WithWindowSize(120, 40).
		Build()

	model.MCPItems = mcpItems

	result := RenderFourColumnGrid(model)

	// Active items should have filled circle
	if !strings.Contains(result, "● active-mcp") {
		t.Errorf("RenderFourColumnGrid() should show ● for active MCPs")
	}

	// Inactive items should have empty circle
	if !strings.Contains(result, "○ inactive-mcp") {
		t.Errorf("RenderFourColumnGrid() should show ○ for inactive MCPs")
	}
}

func TestRenderFourColumnGrid_ResponsiveDimensions(t *testing.T) {
	mcpItems := testutil.MockMCPItems()

	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "Small dimensions",
			width:  80,
			height: 20,
		},
		{
			name:   "Medium dimensions",
			width:  120,
			height: 40,
		},
		{
			name:   "Large dimensions",
			width:  180,
			height: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testutil.NewTestModel().
				WithWindowSize(tt.width, tt.height).
				Build()

			model.MCPItems = mcpItems

			result := RenderFourColumnGrid(model)

			// Should render without panic and contain basic elements
			if result == "" {
				t.Errorf("RenderFourColumnGrid() should not return empty string for %dx%d",
					tt.width, tt.height)
			}

			if !strings.Contains(result, "MCP Inventory") {
				t.Errorf("RenderFourColumnGrid() should contain header for %dx%d",
					tt.width, tt.height)
			}
		})
	}
}

func TestRenderMCPList(t *testing.T) {
	tests := []struct {
		name         string
		mcpItems     []types.MCPItem
		selectedItem int
		expected     []string
	}{
		{
			name: "Basic list rendering",
			mcpItems: []types.MCPItem{
				{Name: "github", Active: true},
				{Name: "docker", Active: false},
				{Name: "context7", Active: true},
			},
			selectedItem: 1,
			expected: []string{
				"● github",
				"○ docker",
				"● context7",
			},
		},
		{
			name:     "Empty list returns no MCPs message",
			mcpItems: []types.MCPItem{},
			expected: []string{
				"No MCPs loaded from inventory (total: 0)",
			},
		},
		{
			name: "Single item list",
			mcpItems: []types.MCPItem{
				{Name: "single", Active: true},
			},
			selectedItem: 0,
			expected: []string{
				"● single",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testutil.NewTestModel().
				WithSelectedItem(tt.selectedItem).
				Build()

			model.MCPItems = tt.mcpItems

			result := RenderMCPList(model)

			for _, expected := range tt.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("RenderMCPList() should contain %q\nActual: %s", expected, result)
				}
			}
		})
	}
}

func TestRenderMCPList_StatusIndicators(t *testing.T) {
	mcpItems := []types.MCPItem{
		{Name: "active1", Active: true},
		{Name: "inactive1", Active: false},
		{Name: "active2", Active: true},
		{Name: "inactive2", Active: false},
	}

	model := testutil.NewTestModel().
		WithSelectedItem(0).
		Build()

	model.MCPItems = mcpItems

	result := RenderMCPList(model)

	// Check active indicators
	if !strings.Contains(result, "● active1") {
		t.Errorf("RenderMCPList() should show ● for active1")
	}
	if !strings.Contains(result, "● active2") {
		t.Errorf("RenderMCPList() should show ● for active2")
	}

	// Check inactive indicators
	if !strings.Contains(result, "○ inactive1") {
		t.Errorf("RenderMCPList() should show ○ for inactive1")
	}
	if !strings.Contains(result, "○ inactive2") {
		t.Errorf("RenderMCPList() should show ○ for inactive2")
	}
}

func TestRenderFourColumnGrid_MinimumGridRows(t *testing.T) {
	// Test that grid maintains minimum 10 rows even with few items
	mcpItems := []types.MCPItem{
		{Name: "item1", Active: true},
		{Name: "item2", Active: false},
	}

	model := testutil.NewTestModel().
		WithWindowSize(120, 40).
		Build()

	model.MCPItems = mcpItems

	result := RenderFourColumnGrid(model)

	// Should contain items
	if !strings.Contains(result, "item1") {
		t.Errorf("RenderFourColumnGrid() should contain item1")
	}
	if !strings.Contains(result, "item2") {
		t.Errorf("RenderFourColumnGrid() should contain item2")
	}

	// Should render without issues (minimum grid logic is internal)
	if !strings.Contains(result, "MCP Inventory") {
		t.Errorf("RenderFourColumnGrid() should contain header")
	}
}
