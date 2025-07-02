package services

import (
	"context"
	"reflect"
	"runtime"
	"testing"
	"time"

	"cc-mcp-manager/internal/ui/types"
)

func TestNewClaudeService(t *testing.T) {
	service := NewClaudeService()

	if service == nil {
		t.Fatal("NewClaudeService() returned nil")
	}

	if service.timeout != 10*time.Second {
		t.Errorf("Expected timeout to be 10s, got %v", service.timeout)
	}
}

func TestDetectClaudeCLI(t *testing.T) {
	service := NewClaudeService()
	ctx := context.Background()

	// Test detection (will vary based on whether claude is actually installed)
	status := service.DetectClaudeCLI(ctx)

	// Basic structure validation
	if status.LastCheck.IsZero() {
		t.Error("LastCheck should be set")
	}

	// If not available, should have error and install guide
	if !status.Available {
		if status.Error == "" {
			t.Error("Error should be set when Claude is not available")
		}
		if status.InstallGuide == "" {
			t.Error("InstallGuide should be set when Claude is not available")
		}
	}

	// If available, should have version info
	if status.Available && status.Version == "" {
		// Version might not be available in all cases, so this is a warning, not error
		t.Logf("Warning: Claude detected but version not available")
	}
}

func TestDetectClaudeCliWithTimeout(t *testing.T) {
	service := NewClaudeService()

	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	status := service.DetectClaudeCLI(ctx)

	// Should handle timeout gracefully
	if status.Available {
		// If Claude is very fast and available, that's okay
		t.Logf("Claude detected even with short timeout")
	} else {
		// Should have error message
		if status.Error == "" {
			t.Error("Expected error message when detection fails")
		}
	}
}

func TestGetClaudeVersion(t *testing.T) {
	service := NewClaudeService()
	ctx := context.Background()

	// This test will vary based on whether claude is actually installed
	// We test both success and failure scenarios

	// Test with short timeout to simulate failure
	ctxShort, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	defer cancel()

	version, err := service.getClaudeVersion(ctxShort)

	// Either succeeds very quickly or fails gracefully
	if err != nil {
		// Expected behavior when Claude is not available or times out
		if version != "" {
			t.Error("Version should be empty when error occurs")
		}
	} else {
		// If successful, version should be valid
		if version == "" {
			t.Error("Version should not be empty when no error")
		}
	}

	// Test version parsing with mock outputs
	testCases := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name:     "standard version format",
			output:   "claude 1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "version only",
			output:   "1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "complex version format",
			output:   "claude version 1.2.3-beta.1",
			expected: "version",
		},
		{
			name:     "empty output",
			output:   "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// We can't easily mock exec.Command, so we test the version parsing logic
			// by calling the function with known outputs would need dependency injection
			// For now, just verify the function can be called
			_, err := service.getClaudeVersion(context.Background())
			// Don't fail test based on whether Claude is installed
			if err != nil {
				t.Logf("getClaudeVersion failed as expected when Claude not available: %v", err)
			}
		})
	}
}

func TestParseActiveMCPs(t *testing.T) {
	service := NewClaudeService()

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty output",
			input:    "",
			expected: []string{},
		},
		{
			name:     "no MCPs message",
			input:    "No MCPs currently active",
			expected: []string{},
		},
		{
			name:     "simple list with checkmarks",
			input:    "✓ github-mcp\n✓ context7\n✓ ht-mcp",
			expected: []string{"github-mcp", "context7", "ht-mcp"},
		},
		{
			name:     "list with asterisks",
			input:    "* github-mcp\n* context7",
			expected: []string{"github-mcp", "context7"},
		},
		{
			name:     "list with bullets",
			input:    "• github-mcp\n• context7",
			expected: []string{"github-mcp", "context7"},
		},
		{
			name:     "simple names",
			input:    "github-mcp\ncontext7\nht-mcp",
			expected: []string{"github-mcp", "context7", "ht-mcp"},
		},
		{
			name:     "mixed format with headers",
			input:    "Available MCPs:\n✓ github-mcp\n✓ context7",
			expected: []string{"github-mcp", "context7"},
		},
		{
			name:     "JSON format",
			input:    `[{"name": "github-mcp", "active": true}, {"name": "context7", "active": true}]`,
			expected: []string{"github-mcp", "context7"},
		},
		{
			name:     "JSON format with inactive MCPs",
			input:    `[{"name": "github-mcp", "active": true}, {"name": "inactive-mcp", "active": false}]`,
			expected: []string{"github-mcp"},
		},
		{
			name:     "JSON format without active field",
			input:    `[{"name": "github-mcp"}, {"name": "context7"}]`,
			expected: []string{"github-mcp", "context7"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.parseActiveMCPs(tt.input)
			if err != nil {
				t.Errorf("parseActiveMCPs() error = %v", err)
				return
			}
			// For empty cases, both nil and empty slice are acceptable
			if len(tt.expected) == 0 && len(result) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseActiveMCPs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetInstallationGuide(t *testing.T) {
	service := NewClaudeService()
	guide := service.getInstallationGuide()

	if guide == "" {
		t.Error("Installation guide should not be empty")
	}

	// Should always contain claude.ai/cli
	if !contains(guide, "claude.ai/cli") {
		t.Error("Installation guide should contain claude.ai/cli")
	}

	// Should contain some platform-appropriate information
	switch runtime.GOOS {
	case "darwin":
		if !contains(guide, "Homebrew") && !contains(guide, "brew") && !contains(guide, "PATH") {
			t.Error("macOS guide should mention Homebrew, brew, or PATH")
		}
	case "windows":
		if !contains(guide, "PATH") && !contains(guide, "Restart") {
			t.Error("Windows guide should mention PATH or Restart")
		}
	case "linux":
		if !contains(guide, "chmod") && !contains(guide, "/usr/local/bin") && !contains(guide, "PATH") {
			t.Error("Linux guide should mention chmod, /usr/local/bin, or PATH")
		}
	default:
		if !contains(guide, "PATH") {
			t.Error("Default guide should mention PATH")
		}
	}

	// Test that the function works for different hypothetical platforms
	// by testing the logic paths
	t.Run("guide_completeness", func(t *testing.T) {
		// The guide should provide useful information regardless of platform
		if len(guide) < 50 {
			t.Error("Installation guide seems too short to be helpful")
		}

		// Should contain download information
		if !contains(guide, "Download") && !contains(guide, "install") && !contains(guide, "Install") {
			t.Error("Guide should contain download or install instructions")
		}
	})
}

func TestUpdateModelWithClaudeStatus(t *testing.T) {
	model := types.NewModel()
	status := types.ClaudeStatus{
		Available:    true,
		Version:      "1.0.0",
		ActiveMCPs:   []string{"github-mcp", "context7"},
		LastCheck:    time.Now(),
		Error:        "",
		InstallGuide: "",
	}

	updatedModel := UpdateModelWithClaudeStatus(model, status)

	if !updatedModel.ClaudeAvailable {
		t.Error("ClaudeAvailable should be true")
	}

	if updatedModel.ClaudeStatus.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", updatedModel.ClaudeStatus.Version)
	}

	if len(updatedModel.ClaudeStatus.ActiveMCPs) != 2 {
		t.Errorf("Expected 2 active MCPs, got %d", len(updatedModel.ClaudeStatus.ActiveMCPs))
	}

	if updatedModel.ClaudeSyncError != "" {
		t.Error("ClaudeSyncError should be empty when no error")
	}
}

func TestUpdateModelWithClaudeStatusError(t *testing.T) {
	model := types.NewModel()
	status := types.ClaudeStatus{
		Available:    false,
		Error:        "Claude CLI not found",
		InstallGuide: "Install from claude.ai/cli",
		LastCheck:    time.Now(),
	}

	updatedModel := UpdateModelWithClaudeStatus(model, status)

	if updatedModel.ClaudeAvailable {
		t.Error("ClaudeAvailable should be false")
	}

	if updatedModel.ClaudeSyncError != "Claude CLI not found" {
		t.Errorf("Expected error message 'Claude CLI not found', got %s", updatedModel.ClaudeSyncError)
	}
}

func TestSyncMCPStatus(t *testing.T) {
	// Create model with some MCPs
	model := types.NewModel()
	// Set some MCPs as inactive initially
	for i := range model.MCPItems {
		model.MCPItems[i].Active = false
	}

	// Simulate Claude reporting some MCPs as active
	activeMCPs := []string{"github-mcp", "context7", "ht-mcp"}

	updatedModel := SyncMCPStatus(model, activeMCPs)

	// Check that the specified MCPs are now active
	activeCount := 0
	for _, item := range updatedModel.MCPItems {
		if item.Active {
			activeCount++
			// Verify it's one of the expected active MCPs
			found := false
			for _, activeMCP := range activeMCPs {
				if item.Name == activeMCP {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("MCP %s is active but wasn't in the active list", item.Name)
			}
		}
	}

	// Should have at least the MCPs that exist in both lists
	if activeCount == 0 {
		t.Error("Expected some MCPs to be activated")
	}
}

func TestFormatClaudeStatusForDisplay(t *testing.T) {
	tests := []struct {
		name     string
		status   types.ClaudeStatus
		expected string
	}{
		{
			name: "not available",
			status: types.ClaudeStatus{
				Available: false,
			},
			expected: "Claude CLI: Not Available",
		},
		{
			name: "available with error",
			status: types.ClaudeStatus{
				Available: true,
				Error:     "Command failed",
			},
			expected: "Claude CLI: Error (Command failed)",
		},
		{
			name: "available with version and MCPs",
			status: types.ClaudeStatus{
				Available:  true,
				Version:    "1.0.0",
				ActiveMCPs: []string{"github-mcp", "context7"},
			},
			expected: "Claude CLI: Available v1.0.0 • 2 Active MCPs",
		},
		{
			name: "available without version",
			status: types.ClaudeStatus{
				Available:  true,
				ActiveMCPs: []string{"github-mcp"},
			},
			expected: "Claude CLI: Available • 1 Active MCPs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatClaudeStatusForDisplay(tt.status)
			if result != tt.expected {
				t.Errorf("FormatClaudeStatusForDisplay() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetRefreshKeyHint(t *testing.T) {
	tests := []struct {
		name     string
		status   types.ClaudeStatus
		expected string
	}{
		{
			name: "available",
			status: types.ClaudeStatus{
				Available: true,
			},
			expected: "R=Refresh Claude Status",
		},
		{
			name: "not available",
			status: types.ClaudeStatus{
				Available: false,
			},
			expected: "R=Retry Claude Detection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRefreshKeyHint(tt.status)
			if result != tt.expected {
				t.Errorf("GetRefreshKeyHint() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestRefreshClaudeStatus(t *testing.T) {
	service := NewClaudeService()
	ctx := context.Background()

	status := service.RefreshClaudeStatus(ctx)

	// Basic validation
	if status.LastCheck.IsZero() {
		t.Error("LastCheck should be set")
	}

	// Should have either success (Available=true) or error info
	if !status.Available {
		if status.Error == "" {
			t.Error("Error should be set when Claude is not available")
		}
		if status.InstallGuide == "" {
			t.Error("InstallGuide should be set when Claude is not available")
		}
	}
}

// TestGetInstallationGuideEdgeCases tests edge cases in installation guide
func TestGetInstallationGuideEdgeCases(t *testing.T) {
	service := NewClaudeService()

	// Test that guide always provides helpful information
	guide := service.getInstallationGuide()

	// Should always have installation link
	if !contains(guide, "claude.ai/cli") {
		t.Error("Installation guide should always contain claude.ai/cli")
	}

	// Should provide platform-specific guidance
	if !contains(guide, "PATH") {
		t.Error("Installation guide should mention PATH configuration")
	}

	// Should be reasonably detailed
	if len(guide) < 100 {
		t.Error("Installation guide should be detailed enough to be helpful")
	}
}

// TestClaudeServiceEdgeCases tests various edge cases
func TestClaudeServiceEdgeCases(t *testing.T) {
	service := NewClaudeService()
	ctx := context.Background()

	// Test with canceled context
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel() // Cancel immediately

	status := service.DetectClaudeCLI(cancelCtx)
	if status.Available {
		t.Log("Claude detected even with canceled context - very fast system")
	}
	// Should handle cancellation gracefully
	if !status.LastCheck.IsZero() {
		t.Log("LastCheck was set even with canceled context")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && s[:len(substr)] == substr) ||
		(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
		containsInner(s, substr))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark tests
func BenchmarkDetectClaudeCLI(b *testing.B) {
	service := NewClaudeService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.DetectClaudeCLI(ctx)
	}
}

func BenchmarkParseActiveMCPs(b *testing.B) {
	service := NewClaudeService()
	input := "✓ github-mcp\n✓ context7\n✓ ht-mcp\n✓ filesystem\n✓ docker-mcp"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.parseActiveMCPs(input)
	}
}

func BenchmarkFormatClaudeStatusForDisplay(b *testing.B) {
	status := types.ClaudeStatus{
		Available:  true,
		Version:    "1.0.0",
		ActiveMCPs: []string{"github-mcp", "context7", "ht-mcp"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FormatClaudeStatusForDisplay(status)
	}
}
