package services

import (
	"context"
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

// ToggleMCPStatus toggles the active status of the currently selected MCP with enhanced error handling
func ToggleMCPStatus(model types.Model) types.Model {
	selectedMCP := GetSelectedMCP(model)
	if selectedMCP == nil {
		return model
	}

	// Prevent toggle if Claude CLI is not available
	if !model.ClaudeAvailable {
		model.ToggleState = types.ToggleError
		model.ToggleError = "Claude CLI not available. Install Claude CLI to manage MCP activation."
		model.ToggleMCPName = selectedMCP.Name
		return model
	}

	// Set loading state immediately for visual feedback
	model.ToggleState = types.ToggleLoading
	model.ToggleMCPName = selectedMCP.Name
	model.ToggleError = ""
	model.ToggleRetrying = false

	return model
}

// EnhancedToggleMCPStatus performs the actual enhanced toggle operation with Claude CLI
func EnhancedToggleMCPStatus(model types.Model, mcpName string, activate bool) types.Model {
	// Find the MCP configuration in the model
	var mcpConfig *types.MCPItem
	for i := range model.MCPItems {
		if model.MCPItems[i].Name == mcpName {
			mcpConfig = &model.MCPItems[i]
			break
		}
	}

	// Create Claude service and perform toggle
	claudeService := NewClaudeService()
	ctx := context.Background()

	result, err := claudeService.ToggleMCPStatus(ctx, mcpName, activate, mcpConfig)
	if err != nil {
		// Unexpected error from service
		model.ToggleState = types.ToggleError
		model.ToggleError = "Internal error during MCP toggle operation"
		model.ToggleMCPName = mcpName
		return model
	}

	if result.Success {
		// Update local MCP status
		for i := range model.MCPItems {
			if model.MCPItems[i].Name == mcpName {
				model.MCPItems[i].Active = activate
				break
			}
		}

		// Save to storage
		if err := SaveInventory(model.MCPItems); err != nil {
			model.ToggleState = types.ToggleError
			model.ToggleError = "MCP toggled but failed to save to storage"
		} else {
			model.ToggleState = types.ToggleSuccess
		}
		model.ToggleMCPName = mcpName
	} else {
		// Handle different error types
		if result.Retryable && !model.ToggleRetrying {
			model.ToggleState = types.ToggleRetrying
			model.ToggleRetrying = true
		} else {
			model.ToggleState = types.ToggleError
			model.ToggleRetrying = false
		}
		model.ToggleError = result.ErrorMsg
		model.ToggleMCPName = mcpName
	}

	return model
}

// LegacyToggleMCPStatus provides backward compatibility for MCP toggle operations
func LegacyToggleMCPStatus(model types.Model) types.Model {
	filteredMCPs := GetFilteredMCPs(model)

	// Use appropriate index based on search state
	selectedIndex := model.SelectedItem
	if model.SearchQuery != "" {
		selectedIndex = model.FilteredSelectedIndex
	}

	if selectedIndex < len(filteredMCPs) {
		// Find the original item and toggle it
		selectedItem := filteredMCPs[selectedIndex]
		for i := range model.MCPItems {
			if model.MCPItems[i].Name == selectedItem.Name {
				model.MCPItems[i].Active = !model.MCPItems[i].Active

				// Save to storage immediately after change
				if err := SaveInventory(model.MCPItems); err != nil {
					// Log error but don't fail the operation
					// Error is already logged in SaveInventory
					// Intentionally empty - MCP status change should succeed even if save fails
					_ = err // Acknowledge error but continue
				}
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
	filteredMCPs := GetFilteredMCPs(model)

	// Use appropriate index based on search state
	selectedIndex := model.SelectedItem
	if model.SearchQuery != "" {
		selectedIndex = model.FilteredSelectedIndex
	}

	if selectedIndex < 0 || selectedIndex >= len(filteredMCPs) {
		return nil
	}

	selectedItem := filteredMCPs[selectedIndex]

	// Return a copy to avoid accidental modifications
	return &selectedItem
}
