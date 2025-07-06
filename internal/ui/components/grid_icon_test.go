package components

import (
	"testing"

	"mcp-hub/internal/ui/types"
)

// Test icon constants
const (
	InactiveIcon = "○"
	ActiveIcon   = "●"
)

func TestGetEnhancedStatusIndicator_ToggleStates(t *testing.T) {
	mcpItem := getTestMCPItem()

	testCases := getStatusIndicatorTestCases()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := buildTestModel(tc.toggleState, tc.toggleMCPName)
			mcp := getMCPForTest(mcpItem, tc.mcpActive, tc.mcpName)

			icon := getEnhancedStatusIndicator(model, mcp)
			assertIconEquals(t, icon, tc.expected, tc.description)
		})
	}
}

func getTestMCPItem() types.MCPItem {
	return types.MCPItem{
		Name:   "test-mcp",
		Active: false,
	}
}

type statusIndicatorTestCase struct {
	name          string
	toggleState   types.ToggleOperationState
	toggleMCPName string
	mcpName       string
	mcpActive     bool
	expected      string
	description   string
}

func getStatusIndicatorTestCases() []statusIndicatorTestCase {
	cases := make([]statusIndicatorTestCase, 0)
	cases = append(cases, getDefaultStateTestCases()...)
	cases = append(cases, getTransitionStateTestCases()...)
	cases = append(cases, getSuccessStateTestCases()...)
	cases = append(cases, getErrorAndEdgeTestCases()...)
	return cases
}

func getDefaultStateTestCases() []statusIndicatorTestCase {
	return []statusIndicatorTestCase{
		{
			name:          "Default inactive state",
			toggleState:   types.ToggleIdle,
			toggleMCPName: "",
			mcpName:       "test-mcp",
			mcpActive:     false,
			expected:      InactiveIcon,
			description:   "inactive state",
		},
		{
			name:          "Default active state",
			toggleState:   types.ToggleIdle,
			toggleMCPName: "",
			mcpName:       "test-mcp",
			mcpActive:     true,
			expected:      ActiveIcon,
			description:   "active state",
		},
	}
}

func getTransitionStateTestCases() []statusIndicatorTestCase {
	return []statusIndicatorTestCase{
		{
			name:          "Loading state",
			toggleState:   types.ToggleLoading,
			toggleMCPName: "test-mcp",
			mcpName:       "test-mcp",
			mcpActive:     false,
			expected:      "⏳",
			description:   "loading state",
		},
		{
			name:          "Retry state",
			toggleState:   types.ToggleRetrying,
			toggleMCPName: "test-mcp",
			mcpName:       "test-mcp",
			mcpActive:     false,
			expected:      "🔄",
			description:   "retry state",
		},
	}
}

func getSuccessStateTestCases() []statusIndicatorTestCase {
	return []statusIndicatorTestCase{
		{
			name:          "Success state - activation",
			toggleState:   types.ToggleSuccess,
			toggleMCPName: "test-mcp",
			mcpName:       "test-mcp",
			mcpActive:     true,
			expected:      "✅",
			description:   "successful activation",
		},
		{
			name:          "Success state - deactivation",
			toggleState:   types.ToggleSuccess,
			toggleMCPName: "test-mcp",
			mcpName:       "test-mcp",
			mcpActive:     false,
			expected:      types.BulletChar,
			description:   "successful deactivation",
		},
	}
}

func getErrorAndEdgeTestCases() []statusIndicatorTestCase {
	return []statusIndicatorTestCase{
		{
			name:          "Error state",
			toggleState:   types.ToggleError,
			toggleMCPName: "test-mcp",
			mcpName:       "test-mcp",
			mcpActive:     false,
			expected:      "✗",
			description:   "error state",
		},
		{
			name:          "Different MCP not affected by toggle state",
			toggleState:   types.ToggleSuccess,
			toggleMCPName: "test-mcp",
			mcpName:       "different-mcp",
			mcpActive:     true,
			expected:      ActiveIcon,
			description:   "different MCP (normal active state)",
		},
	}
}

func buildTestModel(toggleState types.ToggleOperationState, toggleMCPName string) types.Model {
	return types.Model{
		ToggleState:   toggleState,
		ToggleMCPName: toggleMCPName,
	}
}

func getMCPForTest(baseMCP types.MCPItem, active bool, mcpName string) types.MCPItem {
	mcp := baseMCP
	mcp.Active = active
	if mcpName != "" && mcpName != baseMCP.Name {
		mcp.Name = mcpName
	}
	return mcp
}

func assertIconEquals(t *testing.T, actual, expected, description string) {
	if actual != expected {
		t.Errorf("Expected %s for %s, got %s", expected, description, actual)
	}
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
		expected := types.BulletChar
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
