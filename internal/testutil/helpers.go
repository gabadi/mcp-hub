package testutil

import (
	"strings"

	"mcp-hub/internal/ui/types"
)

// GetLayoutType determines layout based on terminal width
func GetLayoutType(width int) string {
	if width >= types.WideLayoutMin {
		return "wide"
	} else if width >= types.MediumLayoutMin {
		return "medium"
	}
	return "narrow"
}

// GetExpectedColumns returns expected column count for given width
func GetExpectedColumns(width int) int {
	if width >= types.WideLayoutMin {
		return types.WideColumns
	} else if width >= types.MediumLayoutMin {
		return types.MediumColumns
	}
	return types.NarrowColumns
}

// GetFilteredMCPs returns MCPs filtered by search query
func GetFilteredMCPs(model types.Model) []types.MCPItem {
	// If no search query, return all MCPs
	if model.SearchQuery == "" {
		return model.MCPItems
	}

	// Filter MCPs by search query directly
	var filtered []types.MCPItem
	query := strings.ToLower(model.SearchQuery)
	for _, item := range model.MCPItems {
		if strings.Contains(strings.ToLower(item.Name), query) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// MockMCPItems returns a smaller set of MCPs for testing
func MockMCPItems() []types.MCPItem {
	return []types.MCPItem{
		{Name: "context7", Type: "SSE", Active: true, Command: "npx @context7/mcp-server"},
		{Name: "github-mcp", Type: "CMD", Active: true, Command: "github-mcp"},
		{Name: "ht-mcp", Type: "CMD", Active: true, Command: "ht-mcp"},
		{Name: "filesystem", Type: "CMD", Active: false, Command: "filesystem-mcp"},
		{Name: "docker-mcp", Type: "CMD", Active: false, Command: "docker-mcp"},
	}
}
