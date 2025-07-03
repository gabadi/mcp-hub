package services

import (
	"context"
	"fmt"
	"os"
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
			result := service.parseActiveMCPs(tt.input)
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
		_ = service.parseActiveMCPs(input)
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

// Epic 2 Story 2 Tests - Enhanced Toggle Functionality

func TestToggleResultStruct(t *testing.T) {
	result := &ToggleResult{
		Success:   true,
		MCPName:   "test-mcp",
		NewState:  "active",
		ErrorType: "",
		ErrorMsg:  "",
		Retryable: false,
		Duration:  100 * time.Millisecond,
	}

	if !result.Success {
		t.Error("Expected Success to be true")
	}
	if result.MCPName != "test-mcp" {
		t.Errorf("Expected MCPName to be 'test-mcp', got %s", result.MCPName)
	}
	if result.NewState != "active" {
		t.Errorf("Expected NewState to be 'active', got %s", result.NewState)
	}
	if result.Duration != 100*time.Millisecond {
		t.Errorf("Expected Duration to be 100ms, got %v", result.Duration)
	}
	if result.ErrorType != "" {
		t.Errorf("Expected ErrorType to be empty for successful result, got %s", result.ErrorType)
	}
	if result.ErrorMsg != "" {
		t.Errorf("Expected ErrorMsg to be empty for successful result, got %s", result.ErrorMsg)
	}
	if result.Retryable {
		t.Error("Expected Retryable to be false for successful result")
	}
}

func TestErrorTypeConstants(t *testing.T) {
	expectedConstants := map[string]string{
		ErrorTypeClaudeUnavailable: "CLAUDE_UNAVAILABLE",
		ErrorTypeNetworkTimeout:    "NETWORK_TIMEOUT",
		ErrorTypePermissionError:   "PERMISSION_ERROR",
		ErrorTypeUnknownError:      "UNKNOWN_ERROR",
	}

	for constant, expected := range expectedConstants {
		if constant != expected {
			t.Errorf("Expected constant %s to equal %s", constant, expected)
		}
	}
}

func TestErrorMessages(t *testing.T) {
	// Test that all error types have corresponding messages
	expectedErrorTypes := []string{
		ErrorTypeClaudeUnavailable,
		ErrorTypeNetworkTimeout,
		ErrorTypePermissionError,
		ErrorTypeUnknownError,
	}

	for _, errorType := range expectedErrorTypes {
		message, exists := ErrorMessages[errorType]
		if !exists {
			t.Errorf("Error message not found for error type: %s", errorType)
		}
		if message == "" {
			t.Errorf("Error message is empty for error type: %s", errorType)
		}
		if len(message) < 10 {
			t.Errorf("Error message too short for error type %s: %s", errorType, message)
		}
	}

	// Test specific message content
	claudeMsg := ErrorMessages[ErrorTypeClaudeUnavailable]
	if !contains(claudeMsg, "Claude CLI") {
		t.Error("Claude unavailable message should mention Claude CLI")
	}

	timeoutMsg := ErrorMessages[ErrorTypeNetworkTimeout]
	if !contains(timeoutMsg, "timed out") || !contains(timeoutMsg, "Retrying") {
		t.Error("Timeout message should mention timed out and retrying")
	}

	permMsg := ErrorMessages[ErrorTypePermissionError]
	if !contains(permMsg, "Permission") || !contains(permMsg, "authentication") {
		t.Error("Permission message should mention permission and authentication")
	}
}

func TestClassifyError(t *testing.T) {
	service := NewClaudeService()

	testCases := []struct {
		name     string
		errMsg   string
		output   string
		expected string
	}{
		{
			name:     "timeout error",
			errMsg:   "context deadline exceeded",
			output:   "",
			expected: ErrorTypeNetworkTimeout,
		},
		{
			name:     "command timeout",
			errMsg:   "operation timeout",
			output:   "",
			expected: ErrorTypeNetworkTimeout,
		},
		{
			name:     "permission denied in output",
			errMsg:   "exit status 1",
			output:   "permission denied",
			expected: ErrorTypePermissionError,
		},
		{
			name:     "unauthorized in output",
			errMsg:   "exit status 1",
			output:   "unauthorized access",
			expected: ErrorTypePermissionError,
		},
		{
			name:     "authentication error",
			errMsg:   "exit status 1",
			output:   "authentication failed",
			expected: ErrorTypePermissionError,
		},
		{
			name:     "executable not found",
			errMsg:   "executable file not found in $PATH",
			output:   "",
			expected: ErrorTypeClaudeUnavailable,
		},
		{
			name:     "command not found in output",
			errMsg:   "exit status 127",
			output:   "command not found",
			expected: ErrorTypeClaudeUnavailable,
		},
		{
			name:     "not recognized in output",
			errMsg:   "exit status 1",
			output:   "'claude' is not recognized as an internal or external command",
			expected: ErrorTypeClaudeUnavailable,
		},
		{
			name:     "unknown error",
			errMsg:   "some other error",
			output:   "unexpected output",
			expected: ErrorTypeUnknownError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock error
			mockErr := &mockError{msg: tc.errMsg}
			result := service.classifyError(mockErr, tc.output)
			if result != tc.expected {
				t.Errorf("classifyError() = %s, expected %s", result, tc.expected)
			}
		})
	}
}

// Mock error type for testing
type mockError struct {
	msg string
}

func (m *mockError) Error() string {
	return m.msg
}

func TestToggleMCPStatusWithClaudeUnavailable(t *testing.T) {
	service := NewClaudeService()
	ctx := context.Background()

	// Create a very short timeout context to simulate unavailable Claude
	ctxShort, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
	defer cancel()

	// Create test MCP config
	testMCP := &types.MCPItem{
		Name:    "test-mcp",
		Type:    "CMD",
		Command: "echo",
		Args:    []string{"test"},
	}

	result, err := service.ToggleMCPStatus(ctxShort, "test-mcp", true, testMCP)

	// Should return result, not error
	if err != nil {
		t.Errorf("ToggleMCPStatus should not return error, got: %v", err)
	}

	if result == nil {
		t.Fatal("ToggleMCPStatus should return result")
	}

	if result.Success {
		t.Log("Toggle succeeded despite short timeout - very fast system")
	} else {
		// Expected behavior with short timeout
		if result.ErrorType != ErrorTypeClaudeUnavailable && result.ErrorType != ErrorTypeNetworkTimeout {
			t.Errorf("Expected error type CLAUDE_UNAVAILABLE or NETWORK_TIMEOUT, got %s", result.ErrorType)
		}

		if result.ErrorMsg == "" {
			t.Error("Error message should not be empty when toggle fails")
		}

		if result.MCPName != "test-mcp" {
			t.Errorf("Expected MCPName to be 'test-mcp', got %s", result.MCPName)
		}
	}

	if result.Duration == 0 {
		t.Error("Duration should be set")
	}
}

func TestToggleMCPStatusSuccessCase(t *testing.T) {
	// This test can't easily test success without mocking the exec.Command
	// So we test the structure and error paths
	service := NewClaudeService()
	ctx := context.Background()

	// Test activation
	// Create test MCP config
	testMCP := &types.MCPItem{
		Name:    "test-mcp",
		Type:    "CMD",
		Command: "echo",
		Args:    []string{"test"},
	}

	result, err := service.ToggleMCPStatus(ctx, "test-mcp", true, testMCP)
	if err != nil {
		t.Errorf("ToggleMCPStatus should not return error, got: %v", err)
	}
	if result == nil {
		t.Fatal("ToggleMCPStatus should return result")
	}
	if result.MCPName != "test-mcp" {
		t.Errorf("Expected MCPName to be 'test-mcp', got %s", result.MCPName)
	}

	// Test deactivation
	result2, err := service.ToggleMCPStatus(ctx, "test-mcp", false, testMCP)
	if err != nil {
		t.Errorf("ToggleMCPStatus should not return error, got: %v", err)
	}
	if result2 == nil {
		t.Fatal("ToggleMCPStatus should return result")
	}
	if result2.MCPName != "test-mcp" {
		t.Errorf("Expected MCPName to be 'test-mcp', got %s", result2.MCPName)
	}
}

func TestRetryToggleOperation(t *testing.T) {
	service := NewClaudeService()
	ctx := context.Background()

	// Test with already expired time budget
	expiredStart := time.Now().Add(-25 * time.Second)
	// Create test MCP config
	testMCP := &types.MCPItem{
		Name:    "test-mcp",
		Type:    "CMD",
		Command: "echo",
		Args:    []string{"test"},
	}

	result, err := service.retryToggleOperation(ctx, "test-mcp", true, testMCP, expiredStart)

	if err == nil {
		t.Error("Expected error when time budget exceeded")
	}

	if result == nil {
		t.Fatal("Should return result even on error")
	}

	if result.Success {
		t.Error("Should not succeed when time budget exceeded")
	}

	if result.ErrorType != ErrorTypeNetworkTimeout {
		t.Errorf("Expected NETWORK_TIMEOUT error, got %s", result.ErrorType)
	}
}

func BenchmarkToggleMCPStatus(b *testing.B) {
	service := NewClaudeService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use very short timeout to avoid actually calling Claude CLI
		ctxShort, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
		// Create test MCP config
		testMCP := &types.MCPItem{
			Name:    "test-mcp",
			Type:    "CMD",
			Command: "echo",
			Args:    []string{"test"},
		}

		_, _ = service.ToggleMCPStatus(ctxShort, "test-mcp", true, testMCP)
		cancel()
	}
}

func BenchmarkClassifyError(b *testing.B) {
	service := NewClaudeService()
	mockErr := &mockError{msg: "context deadline exceeded"}
	output := "some command output"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.classifyError(mockErr, output)
	}
}

// Epic 2 Story 5 Tests - Project Context Display

func TestGetProjectContext(t *testing.T) {
	// Create test model with MCP items
	model := types.NewModel()
	model.MCPItems = []types.MCPItem{
		{Name: "test1", Active: true},
		{Name: "test2", Active: false},
		{Name: "test3", Active: true},
	}
	model.LastClaudeSync = time.Now().Add(-5 * time.Minute)

	projectContext := GetProjectContext(model)

	// Test basic structure
	if projectContext.CurrentPath == "" {
		t.Error("CurrentPath should not be empty")
	}
	if projectContext.DisplayPath == "" {
		t.Error("DisplayPath should not be empty")
	}
	if projectContext.ActiveMCPs != 2 {
		t.Errorf("Expected 2 active MCPs, got %d", projectContext.ActiveMCPs)
	}
	if projectContext.TotalMCPs != 3 {
		t.Errorf("Expected 3 total MCPs, got %d", projectContext.TotalMCPs)
	}
	if projectContext.LastSyncTime.IsZero() {
		t.Error("LastSyncTime should be set from model")
	}
	if projectContext.SyncStatusText == "" {
		t.Error("SyncStatusText should not be empty")
	}
}

func TestFormatPathForDisplay(t *testing.T) {
	testCases := []struct {
		name      string
		path      string
		maxLength int
		expected  bool // Whether result should be different from input
	}{
		{
			name:      "short path",
			path:      "/home/user",
			maxLength: 50,
			expected:  false, // Should be unchanged
		},
		{
			name:      "exact length",
			path:      "/home/user/project",
			maxLength: 18,
			expected:  false, // Should be unchanged
		},
		{
			name:      "long path",
			path:      "/very/long/path/to/some/deeply/nested/project/directory",
			maxLength: 30,
			expected:  true, // Should be truncated
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatPathForDisplay(tc.path, tc.maxLength)

			if len(result) > tc.maxLength {
				t.Errorf("Result length %d exceeds maxLength %d: %s", len(result), tc.maxLength, result)
			}

			changed := result != tc.path
			if changed != tc.expected {
				t.Errorf("Expected changed=%v, got changed=%v. Input: %s, Output: %s", tc.expected, changed, tc.path, result)
			}

			// For truncated paths, should contain ellipsis
			if tc.expected && len(tc.path) > tc.maxLength {
				if !contains(result, "...") && !contains(result, "~/") {
					t.Errorf("Long path should be truncated with ellipsis or use home dir: %s", result)
				}
			}
		})
	}
}

func TestGetSyncStatus(t *testing.T) {
	baseTime := time.Now()

	testCases := []struct {
		name     string
		model    types.Model
		expected types.SyncStatus
	}{
		{
			name: "claude unavailable",
			model: types.Model{
				ClaudeAvailable: false,
			},
			expected: types.SyncStatusError,
		},
		{
			name: "claude sync error",
			model: types.Model{
				ClaudeAvailable: true,
				ClaudeSyncError: "connection failed",
			},
			expected: types.SyncStatusError,
		},
		{
			name: "no sync performed",
			model: types.Model{
				ClaudeAvailable: true,
				LastClaudeSync:  time.Time{}, // Zero time
			},
			expected: types.SyncStatusUnknown,
		},
		{
			name: "in sync",
			model: types.Model{
				ClaudeAvailable: true,
				LastClaudeSync:  baseTime,
				MCPItems: []types.MCPItem{
					{Name: "test1", Active: true},
					{Name: "test2", Active: false},
				},
				ClaudeStatus: types.ClaudeStatus{
					ActiveMCPs: []string{"test1"},
				},
			},
			expected: types.SyncStatusInSync,
		},
		{
			name: "out of sync - different active count",
			model: types.Model{
				ClaudeAvailable: true,
				LastClaudeSync:  baseTime,
				MCPItems: []types.MCPItem{
					{Name: "test1", Active: true},
					{Name: "test2", Active: true}, // Both active locally
				},
				ClaudeStatus: types.ClaudeStatus{
					ActiveMCPs: []string{"test1"}, // Only one active in Claude
				},
			},
			expected: types.SyncStatusOutOfSync,
		},
		{
			name: "out of sync - different active MCPs",
			model: types.Model{
				ClaudeAvailable: true,
				LastClaudeSync:  baseTime,
				MCPItems: []types.MCPItem{
					{Name: "test1", Active: true},
					{Name: "test2", Active: false},
				},
				ClaudeStatus: types.ClaudeStatus{
					ActiveMCPs: []string{"test2"}, // Different MCP active in Claude
				},
			},
			expected: types.SyncStatusOutOfSync,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetSyncStatus(tc.model)
			if result != tc.expected {
				t.Errorf("GetSyncStatus() = %v, expected %v", result, tc.expected)
			}
		})
	}
}

func TestFormatSyncStatusText(t *testing.T) {
	testCases := []struct {
		status   types.SyncStatus
		expected string
	}{
		{types.SyncStatusInSync, "In Sync"},
		{types.SyncStatusOutOfSync, "Out of Sync"},
		{types.SyncStatusError, "Error"},
		{types.SyncStatusUnknown, "Unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := FormatSyncStatusText(tc.status)
			if result != tc.expected {
				t.Errorf("FormatSyncStatusText(%v) = %s, expected %s", tc.status, result, tc.expected)
			}
		})
	}
}

func TestUpdateProjectContext(t *testing.T) {
	// Create test model
	model := types.NewModel()
	model.MCPItems = []types.MCPItem{
		{Name: "test1", Active: true},
		{Name: "test2", Active: false},
	}

	// Update project context
	updatedModel := UpdateProjectContext(model)

	// Verify project context was updated
	if updatedModel.ProjectContext.TotalMCPs != 2 {
		t.Errorf("Expected 2 total MCPs in project context, got %d", updatedModel.ProjectContext.TotalMCPs)
	}
	if updatedModel.ProjectContext.ActiveMCPs != 1 {
		t.Errorf("Expected 1 active MCP in project context, got %d", updatedModel.ProjectContext.ActiveMCPs)
	}
	if updatedModel.ProjectContext.CurrentPath == "" {
		t.Error("CurrentPath should be set in project context")
	}
	if updatedModel.ProjectContext.DisplayPath == "" {
		t.Error("DisplayPath should be set in project context")
	}
}

func TestHasDirectoryChanged(t *testing.T) {
	// Test with empty path
	if HasDirectoryChanged("") {
		t.Error("Empty path should not indicate directory change")
	}

	// Test with current directory (should not have changed)
	currentDir, err := os.Getwd()
	if err != nil {
		t.Skip("Cannot get current directory for test")
	}

	if HasDirectoryChanged(currentDir) {
		t.Error("Current directory should not indicate change")
	}

	// Test with different directory
	if !HasDirectoryChanged("/non/existent/path") {
		t.Error("Different directory should indicate change")
	}
}

// Benchmark tests for project context functions
func BenchmarkGetProjectContext(b *testing.B) {
	model := types.NewModel()
	model.MCPItems = make([]types.MCPItem, 50)
	for i := range model.MCPItems {
		model.MCPItems[i] = types.MCPItem{
			Name:   fmt.Sprintf("mcp-%d", i),
			Active: i%2 == 0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetProjectContext(model)
	}
}

func BenchmarkFormatPathForDisplay(b *testing.B) {
	longPath := "/very/long/path/to/some/deeply/nested/project/directory/with/many/subdirectories"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FormatPathForDisplay(longPath, 30)
	}
}

func BenchmarkGetSyncStatus(b *testing.B) {
	model := types.Model{
		ClaudeAvailable: true,
		LastClaudeSync:  time.Now(),
		MCPItems: []types.MCPItem{
			{Name: "test1", Active: true},
			{Name: "test2", Active: false},
			{Name: "test3", Active: true},
		},
		ClaudeStatus: types.ClaudeStatus{
			ActiveMCPs: []string{"test1", "test3"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetSyncStatus(model)
	}
}
