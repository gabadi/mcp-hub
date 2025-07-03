package components

import (
	"testing"

	"cc-mcp-manager/internal/ui/types"
)

func TestGetEnhancedStatusIndicator_ToggleStates(t *testing.T) {
	// Create test MCP item
	mcpItem := types.MCPItem{
		Name:   "test-mcp",
		Active: false,
	}

	t.Run("Default inactive state", func(t *testing.T) {
		model := types.Model{
			ToggleState:   types.ToggleIdle,
			ToggleMCPName: "",
		}
		
		icon := getEnhancedStatusIndicator(model, mcpItem)
		expected := "○"
		if icon != expected {
			t.Errorf("Expected %s for inactive state, got %s", expected, icon)
		}
	})

	t.Run("Default active state", func(t *testing.T) {
		activeMCP := mcpItem
		activeMCP.Active = true
		
		model := types.Model{
			ToggleState:   types.ToggleIdle,
			ToggleMCPName: "",
		}
		
		icon := getEnhancedStatusIndicator(model, activeMCP)
		expected := "●"
		if icon != expected {
			t.Errorf("Expected %s for active state, got %s", expected, icon)
		}
	})

	t.Run("Loading state", func(t *testing.T) {
		model := types.Model{
			ToggleState:   types.ToggleLoading,
			ToggleMCPName: "test-mcp",
		}
		
		icon := getEnhancedStatusIndicator(model, mcpItem)
		expected := "⏳"
		if icon != expected {
			t.Errorf("Expected %s for loading state, got %s", expected, icon)
		}
	})

	t.Run("Retry state", func(t *testing.T) {
		model := types.Model{
			ToggleState:   types.ToggleRetrying,
			ToggleMCPName: "test-mcp",
		}
		
		icon := getEnhancedStatusIndicator(model, mcpItem)
		expected := "🔄"
		if icon != expected {
			t.Errorf("Expected %s for retry state, got %s", expected, icon)
		}
	})

	t.Run("Success state - activation", func(t *testing.T) {
		activeMCP := mcpItem
		activeMCP.Active = true
		
		model := types.Model{
			ToggleState:   types.ToggleSuccess,
			ToggleMCPName: "test-mcp",
		}
		
		icon := getEnhancedStatusIndicator(model, activeMCP)
		expected := "✅"
		if icon != expected {
			t.Errorf("Expected %s for successful activation, got %s", expected, icon)
		}
	})

	t.Run("Success state - deactivation", func(t *testing.T) {
		inactiveMCP := mcpItem
		inactiveMCP.Active = false
		
		model := types.Model{
			ToggleState:   types.ToggleSuccess,
			ToggleMCPName: "test-mcp",
		}
		
		icon := getEnhancedStatusIndicator(model, inactiveMCP)
		expected := "◦"
		if icon != expected {
			t.Errorf("Expected %s for successful deactivation, got %s", expected, icon)
		}
	})

	t.Run("Error state", func(t *testing.T) {
		model := types.Model{
			ToggleState:   types.ToggleError,
			ToggleMCPName: "test-mcp",
		}
		
		icon := getEnhancedStatusIndicator(model, mcpItem)
		expected := "✗"
		if icon != expected {
			t.Errorf("Expected %s for error state, got %s", expected, icon)
		}
	})

	t.Run("Different MCP not affected by toggle state", func(t *testing.T) {
		differentMCP := types.MCPItem{
			Name:   "different-mcp",
			Active: true,
		}
		
		model := types.Model{
			ToggleState:   types.ToggleSuccess,
			ToggleMCPName: "test-mcp", // Different MCP is being toggled
		}
		
		icon := getEnhancedStatusIndicator(model, differentMCP)
		expected := "●" // Should show normal active state
		if icon != expected {
			t.Errorf("Expected %s for different MCP (normal active state), got %s", expected, icon)
		}
	})
}

func TestToggleVisualFix_ConfusingCheckmark(t *testing.T) {
	t.Run("Removed MCP shows clear deactivation success icon", func(t *testing.T) {
		// This test specifically addresses the issue described in the bug report
		// where removed MCPs showed a confusing checkmark
		
		removedMCP := types.MCPItem{
			Name:   "github-mcp",
			Active: false, // MCP was successfully removed/deactivated
		}
		
		model := types.Model{
			ToggleState:   types.ToggleSuccess,
			ToggleMCPName: "github-mcp",
		}
		
		icon := getEnhancedStatusIndicator(model, removedMCP)
		
		// Should NOT be a confusing checkmark "✓"
		if icon == "✓" {
			t.Error("BUG: Removed MCP shows confusing checkmark ✓ - this was the original issue")
		}
		
		// Should be the new clear deactivation success icon
		expected := "◦"
		if icon != expected {
			t.Errorf("Expected %s for successful removal/deactivation, got %s", expected, icon)
		}
		
		// Verify it's visually distinct from regular inactive state
		regularInactiveModel := types.Model{
			ToggleState:   types.ToggleIdle,
			ToggleMCPName: "",
		}
		regularIcon := getEnhancedStatusIndicator(regularInactiveModel, removedMCP)
		
		if icon == regularIcon {
			t.Error("Success deactivation icon should be visually distinct from regular inactive icon")
		}
	})
}