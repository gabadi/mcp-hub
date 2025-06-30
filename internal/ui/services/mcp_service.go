package services

import (
	"strings"

	"cc-mcp-manager/internal/ui/types"
)

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

// ToggleMCPStatus toggles the active status of the currently selected MCP
func ToggleMCPStatus(model types.Model) types.Model {
	filteredMCPs := GetFilteredMCPs(model)
	if model.SelectedItem < len(filteredMCPs) {
		// Find the original item and toggle it
		selectedItem := filteredMCPs[model.SelectedItem]
		for i := range model.MCPItems {
			if model.MCPItems[i].Name == selectedItem.Name {
				model.MCPItems[i].Active = !model.MCPItems[i].Active
				break
			}
		}
	}
	return model
}

// GetActiveMCPCount returns the number of active MCPs
func GetActiveMCPCount(model types.Model) int {
	count := 0
	for _, item := range model.MCPItems {
		if item.Active {
			count++
		}
	}
	return count
}

// GetSelectedMCP returns the currently selected MCP item, or nil if none selected
func GetSelectedMCP(model types.Model) *types.MCPItem {
	if model.SelectedItem >= len(model.MCPItems) {
		return nil
	}
	return &model.MCPItems[model.SelectedItem]
}