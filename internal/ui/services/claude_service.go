package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"cc-mcp-manager/internal/ui/types"
)

// Use types.ClaudeStatus for consistency

// ClaudeService handles Claude CLI interactions
type ClaudeService struct {
	timeout time.Duration
}

// NewClaudeService creates a new Claude service instance
func NewClaudeService() *ClaudeService {
	return &ClaudeService{
		timeout: 10 * time.Second, // 10 second timeout for commands
	}
}

// DetectClaudeCLI checks if Claude CLI is available on the system
func (cs *ClaudeService) DetectClaudeCLI(ctx context.Context) types.ClaudeStatus {
	status := types.ClaudeStatus{
		Available: false,
		LastCheck: time.Now(),
	}

	// Try to detect Claude CLI using which/where command
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "where", "claude")
	default:
		cmd = exec.CommandContext(ctx, "which", "claude")
	}

	// Run with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, cs.timeout)
	defer cancel()
	// nolint:gosec // G204: Args are from a command we just created, not user input
	cmd = exec.CommandContext(timeoutCtx, cmd.Args[0], cmd.Args[1:]...)

	output, err := cmd.Output()
	if err != nil {
		status.Error = "Claude CLI not found in system PATH"
		status.InstallGuide = cs.getInstallationGuide()
		return status
	}

	// If we found claude, try to get version
	claudePath := strings.TrimSpace(string(output))
	if claudePath != "" {
		status.Available = true
		version, err := cs.getClaudeVersion(timeoutCtx)
		if err != nil {
			status.Error = fmt.Sprintf("Found Claude CLI but failed to get version: %v", err)
		} else {
			status.Version = version
		}
	}

	return status
}

// getClaudeVersion attempts to get the Claude CLI version
func (cs *ClaudeService) getClaudeVersion(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "claude", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get claude version: %w", err)
	}

	version := strings.TrimSpace(string(output))
	// Handle different version output formats
	if strings.Contains(version, "claude") {
		parts := strings.Fields(version)
		if len(parts) >= 2 {
			return parts[1], nil
		}
	}
	return version, nil
}

// QueryActiveMCPs queries Claude CLI for currently active MCPs
func (cs *ClaudeService) QueryActiveMCPs(ctx context.Context) ([]string, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, cs.timeout)
	defer cancel()

	cmd := exec.CommandContext(timeoutCtx, "claude", "mcp", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to query active MCPs: %w", err)
	}

	// Parse the output to extract active MCP names
	return cs.parseActiveMCPs(string(output))
}

// parseActiveMCPs parses the output of 'claude mcp list' command
func (cs *ClaudeService) parseActiveMCPs(output string) ([]string, error) {
	var activeMCPs []string
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try JSON parsing first
		if jsonMCPs := cs.tryParseJSONLine(line); jsonMCPs != nil {
			activeMCPs = append(activeMCPs, jsonMCPs...)
			continue
		}

		// Parse as plain text format
		if mcpName := cs.parsePlainTextLine(line); mcpName != "" {
			activeMCPs = append(activeMCPs, mcpName)
		}
	}

	return activeMCPs, nil
}

// tryParseJSONLine attempts to parse a line as JSON and extract active MCPs
func (cs *ClaudeService) tryParseJSONLine(line string) []string {
	if !strings.HasPrefix(line, "{") && !strings.HasPrefix(line, "[") {
		return nil
	}

	var jsonMCPs []map[string]interface{}
	if err := json.Unmarshal([]byte(line), &jsonMCPs); err != nil {
		return nil
	}

	var activeMCPs []string
	for _, mcp := range jsonMCPs {
		if name := cs.extractMCPNameFromJSON(mcp); name != "" {
			activeMCPs = append(activeMCPs, name)
		}
	}

	return activeMCPs
}

// extractMCPNameFromJSON extracts MCP name from JSON object if it's active
func (cs *ClaudeService) extractMCPNameFromJSON(mcp map[string]interface{}) string {
	name, ok := mcp["name"].(string)
	if !ok {
		return ""
	}

	// Check if MCP is active
	if active, exists := mcp["active"]; exists {
		if isActive, ok := active.(bool); ok && isActive {
			return name
		}
		return ""
	}

	// If no active field, assume it's listed because it's active
	return name
}

// parsePlainTextLine parses a plain text line to extract MCP name
func (cs *ClaudeService) parsePlainTextLine(line string) string {
	// Look for patterns like "✓ mcp-name" or "* mcp-name" or just "mcp-name"
	if cs.hasActiveIndicator(line) {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			return parts[1]
		}
		return ""
	}

	// Treat the line as an MCP name if it doesn't contain exclusion patterns
	if cs.shouldExcludeLine(line) {
		return ""
	}

	return line
}

// hasActiveIndicator checks if line has an active MCP indicator
func (cs *ClaudeService) hasActiveIndicator(line string) bool {
	return strings.Contains(line, "✓") || strings.Contains(line, "*") || strings.Contains(line, "•")
}

// shouldExcludeLine checks if line should be excluded from parsing
func (cs *ClaudeService) shouldExcludeLine(line string) bool {
	return strings.Contains(line, "No MCPs") || strings.Contains(line, "Available")
}

// RefreshClaudeStatus performs a complete refresh of Claude status
func (cs *ClaudeService) RefreshClaudeStatus(ctx context.Context) types.ClaudeStatus {
	status := cs.DetectClaudeCLI(ctx)

	if status.Available && status.Error == "" {
		// Query active MCPs if Claude is available
		activeMCPs, err := cs.QueryActiveMCPs(ctx)
		if err != nil {
			status.Error = fmt.Sprintf("Failed to query active MCPs: %v", err)
		} else {
			status.ActiveMCPs = activeMCPs
		}
	}

	return status
}

// getInstallationGuide returns platform-specific installation guidance
func (cs *ClaudeService) getInstallationGuide() string {
	switch runtime.GOOS {
	case "darwin":
		return "Install Claude CLI:\n• Download from: https://claude.ai/cli\n• Or use Homebrew: brew install claude-cli\n• Ensure it's in your PATH"
	case "windows":
		return "Install Claude CLI:\n• Download from: https://claude.ai/cli\n• Add to your system PATH\n• Restart your terminal after installation"
	case "linux":
		return "Install Claude CLI:\n• Download from: https://claude.ai/cli\n• Make executable: chmod +x claude\n• Add to PATH: sudo mv claude /usr/local/bin/\n• Or use package manager if available"
	default:
		return "Install Claude CLI:\n• Download from: https://claude.ai/cli\n• Follow platform-specific installation instructions\n• Ensure it's available in your system PATH"
	}
}

// UpdateModelWithClaudeStatus updates the model with Claude status information
func UpdateModelWithClaudeStatus(model types.Model, status types.ClaudeStatus) types.Model {
	model.ClaudeAvailable = status.Available
	model.ClaudeStatus = status
	model.LastClaudeSync = status.LastCheck
	if status.Error != "" {
		model.ClaudeSyncError = status.Error
	} else {
		model.ClaudeSyncError = ""
	}
	return model
}

// SyncMCPStatus synchronizes local MCP status with Claude's active MCPs
func SyncMCPStatus(model types.Model, activeMCPs []string) types.Model {
	// Create a map for quick lookup of active MCPs from Claude
	claudeActiveMCPs := make(map[string]bool)
	for _, mcpName := range activeMCPs {
		claudeActiveMCPs[mcpName] = true
	}

	// Update local MCP items based on Claude's active MCPs
	for i := range model.MCPItems {
		mcpName := model.MCPItems[i].Name
		if claudeActiveMCPs[mcpName] {
			model.MCPItems[i].Active = true
		}
		// Note: We don't set inactive status here to avoid disabling
		// MCPs that might be active but not reported by Claude
	}

	return model
}

// FormatClaudeStatusForDisplay formats Claude status for UI display
func FormatClaudeStatusForDisplay(status types.ClaudeStatus) string {
	if !status.Available {
		return "Claude CLI: Not Available"
	}

	if status.Error != "" {
		return fmt.Sprintf("Claude CLI: Error (%s)", status.Error)
	}

	activeMCPCount := len(status.ActiveMCPs)
	versionInfo := ""
	if status.Version != "" {
		versionInfo = fmt.Sprintf(" v%s", status.Version)
	}

	return fmt.Sprintf("Claude CLI: Available%s • %d Active MCPs", versionInfo, activeMCPCount)
}

// GetRefreshKeyHint returns the refresh key hint for the UI
func GetRefreshKeyHint(status types.ClaudeStatus) string {
	if status.Available {
		return "R=Refresh Claude Status"
	}
	return "R=Retry Claude Detection"
}
