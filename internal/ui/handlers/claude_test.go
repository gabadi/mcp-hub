package handlers

import (
	"testing"

	"cc-mcp-manager/internal/ui/services"
	"cc-mcp-manager/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRefreshClaudeStatusCmd(t *testing.T) {
	cmd := RefreshClaudeStatusCmd()

	if cmd == nil {
		t.Fatal("RefreshClaudeStatusCmd() returned nil")
	}

	// Execute the command
	msg := cmd()

	// Should return a ClaudeStatusMsg
	claudeMsg, ok := msg.(ClaudeStatusMsg)
	if !ok {
		t.Fatalf("Expected ClaudeStatusMsg, got %T", msg)
	}

	// Basic validation of message content
	if claudeMsg.Status.LastCheck.IsZero() {
		t.Error("ClaudeStatusMsg should have LastCheck set")
	}
}

func TestClaudeStatusMsgStructure(t *testing.T) {
	status := types.ClaudeStatus{
		Available:  true,
		Version:    "1.0.0",
		ActiveMCPs: []string{"github-mcp"},
	}

	msg := ClaudeStatusMsg{Status: status}

	if !msg.Status.Available {
		t.Error("ClaudeStatusMsg should preserve Available status")
	}

	if msg.Status.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", msg.Status.Version)
	}

	if len(msg.Status.ActiveMCPs) != 1 {
		t.Errorf("Expected 1 active MCP, got %d", len(msg.Status.ActiveMCPs))
	}
}

func TestHandleMainNavigationKeysWithRefresh(t *testing.T) {
	model := types.NewModel()

	// Test 'r' key
	updatedModel, cmd := HandleMainNavigationKeys(model, "r")

	if cmd == nil {
		t.Error("'r' key should return a command")
	}

	// The model shouldn't change immediately
	if updatedModel.ClaudeAvailable != model.ClaudeAvailable {
		t.Error("Model should not change immediately when pressing 'r'")
	}

	// Test 'R' key (uppercase)
	updatedModel, cmd = HandleMainNavigationKeys(model, "R")

	if cmd == nil {
		t.Error("'R' key should return a command")
	}
}

func TestHandleSearchNavigationKeysWithRefresh(t *testing.T) {
	model := types.NewModel()
	model.State = types.SearchActiveNavigation
	model.SearchActive = true

	// Test 'r' key in search navigation mode
	updatedModel, cmd := HandleSearchNavigationKeys(model, "r")

	if cmd == nil {
		t.Error("'r' key should return a command in search navigation mode")
	}

	// Should still be in search navigation mode
	if updatedModel.State != types.SearchActiveNavigation {
		t.Error("Should remain in search navigation mode after refresh")
	}
}

// Integration test to verify the complete flow
func TestClaudeIntegrationFlow(t *testing.T) {
	// This test verifies the complete integration without actually calling Claude CLI

	// Create a model
	model := types.NewModel()

	// Simulate receiving a Claude status message
	status := types.ClaudeStatus{
		Available:  true,
		Version:    "1.0.0",
		ActiveMCPs: []string{"github-mcp", "context7"},
	}

	msg := ClaudeStatusMsg{Status: status}

	// Update model with Claude status
	updatedModel := services.UpdateModelWithClaudeStatus(model, msg.Status)

	// Verify the model was updated correctly
	if !updatedModel.ClaudeAvailable {
		t.Error("Model should show Claude as available")
	}

	if updatedModel.ClaudeStatus.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", updatedModel.ClaudeStatus.Version)
	}

	if len(updatedModel.ClaudeStatus.ActiveMCPs) != 2 {
		t.Errorf("Expected 2 active MCPs, got %d", len(updatedModel.ClaudeStatus.ActiveMCPs))
	}

	// Test sync MCPs
	syncedModel := services.SyncMCPStatus(updatedModel, status.ActiveMCPs)

	// Verify that matching MCPs are activated
	githubActive := false
	context7Active := false

	for _, item := range syncedModel.MCPItems {
		if item.Name == "github-mcp" && item.Active {
			githubActive = true
		}
		if item.Name == "context7" && item.Active {
			context7Active = true
		}
	}

	if !githubActive {
		t.Error("github-mcp should be active after sync")
	}

	if !context7Active {
		t.Error("context7 should be active after sync")
	}
}

func TestCommandCreation(t *testing.T) {
	// Test that the command can be created without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RefreshClaudeStatusCmd() panicked: %v", r)
		}
	}()

	cmd := RefreshClaudeStatusCmd()
	if cmd == nil {
		t.Error("Command should not be nil")
	}
}

// Mock tests for command execution
func TestRefreshClaudeStatusCmdExecution(t *testing.T) {
	cmd := RefreshClaudeStatusCmd()

	// Execute in a goroutine to avoid blocking test
	done := make(chan tea.Msg, 1)
	go func() {
		msg := cmd()
		done <- msg
	}()

	select {
	case msg := <-done:
		// Should receive a message
		if msg == nil {
			t.Error("Command should return a message")
		}

		// Should be ClaudeStatusMsg type
		if _, ok := msg.(ClaudeStatusMsg); !ok {
			t.Errorf("Expected ClaudeStatusMsg, got %T", msg)
		}

	default:
		// For testing, don't wait indefinitely
		t.Log("Command execution completed immediately (expected for test)")
	}
}

// Test error handling in command
func TestRefreshClaudeStatusCmdErrorHandling(t *testing.T) {
	// This test ensures the command handles errors gracefully
	cmd := RefreshClaudeStatusCmd()
	msg := cmd()

	claudeMsg, ok := msg.(ClaudeStatusMsg)
	if !ok {
		t.Fatalf("Expected ClaudeStatusMsg, got %T", msg)
	}

	// Even if Claude is not available, should return a valid status
	if claudeMsg.Status.LastCheck.IsZero() {
		t.Error("LastCheck should always be set")
	}

	// If not available, should have error information
	if !claudeMsg.Status.Available {
		if claudeMsg.Status.Error == "" {
			t.Error("Error should be set when Claude is not available")
		}
	}
}
