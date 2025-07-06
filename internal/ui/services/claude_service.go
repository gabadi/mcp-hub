// Package services provides Claude CLI integration and MCP management services.
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"mcp-hub/internal/platform"
	"mcp-hub/internal/ui/types"
)

// Use types.ClaudeStatus for consistency

const (
	// ClaudeCommand is the command name for the Claude CLI
	ClaudeCommand = "claude"

	// TestActiveStatus represents the active status string for testing
	TestActiveStatus = "active"
)

// allowedCommands defines the commands that are allowed to be executed
var allowedCommands = map[string]bool{
	"claude": true,
}

// ToggleResult represents the result of an MCP toggle operation
type ToggleResult struct {
	Success   bool
	MCPName   string
	NewState  string
	ErrorType string
	ErrorMsg  string
	Retryable bool
	Retrying  bool
	Duration  time.Duration
}

// Error type constants for toggle operations
const (
	ErrorTypeClaudeUnavailable = "CLAUDE_UNAVAILABLE"
	// Service string constants
	UnknownStatus             = "Unknown"
	ErrorTypeNetworkTimeout   = "NETWORK_TIMEOUT"
	ErrorTypePermissionError  = "PERMISSION_ERROR"
	ErrorTypeMCPAlreadyExists = "MCP_ALREADY_EXISTS"
	ErrorTypeMCPNotFound      = "MCP_NOT_FOUND"
	ErrorTypeInvalidCommand   = "INVALID_COMMAND"
	ErrorTypeUnknownError     = "UNKNOWN_ERROR"
)

// ErrorMessages provides templates for user-friendly error messages
var ErrorMessages = map[string]string{
	ErrorTypeClaudeUnavailable: "Claude CLI not available. Install Claude CLI to manage MCP activation.",
	ErrorTypeNetworkTimeout:    "MCP toggle operation timed out. Retrying...",
	ErrorTypePermissionError:   "Permission denied. Check Claude CLI authentication.",
	ErrorTypeMCPAlreadyExists:  "MCP is already active in Claude CLI.",
	ErrorTypeMCPNotFound:       "MCP not found in Claude CLI configuration.",
	ErrorTypeInvalidCommand:    "Invalid MCP command or configuration. Check MCP settings.",
	ErrorTypeUnknownError:      "MCP toggle failed. Press 'R' to refresh and try again.",
}

// ClaudeService handles Claude CLI interactions with platform abstraction
type ClaudeService struct {
	timeout         time.Duration
	platformService platform.PlatformService
}

// NewClaudeService creates a new Claude service instance with platform abstraction
func NewClaudeService(platformService platform.PlatformService) *ClaudeService {
	return &ClaudeService{
		timeout:         10 * time.Second, // 10 second timeout for commands
		platformService: platformService,
	}
}

// DetectClaudeCLI checks if Claude CLI is available on the system
func (cs *ClaudeService) DetectClaudeCLI(ctx context.Context) types.ClaudeStatus {
	status := types.ClaudeStatus{
		Available: false,
		LastCheck: time.Now(),
	}

	// Try to detect Claude CLI using platform-specific command detection
	var cmd *exec.Cmd
	detectionCmd := cs.platformService.GetCommandDetectionCommand()
	cmd = exec.CommandContext(ctx, detectionCmd, "claude")

	// Run with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, cs.timeout)
	defer cancel()
	//nolint:gosec // G204: Args are from a command we just created, not user input
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
	cmd := exec.CommandContext(ctx, ClaudeCommand, "--version")
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

	cmd := exec.CommandContext(timeoutCtx, ClaudeCommand, "mcp", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to query active MCPs: %w", err)
	}

	// Parse the output to extract active MCP names
	activeMCPs := cs.parseActiveMCPs(string(output))
	return activeMCPs, nil
}

// parseActiveMCPs parses the output of 'claude mcp list' command
func (cs *ClaudeService) parseActiveMCPs(output string) []string {
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

	return activeMCPs
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
	if active, exists := mcp[TestActiveStatus]; exists {
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

	// Check for "name: command" format from claude mcp list
	if strings.Contains(line, ":") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) >= 1 {
			name := strings.TrimSpace(parts[0])
			if name != "" && !cs.shouldExcludeLine(name) {
				return name
			}
		}
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
	switch cs.platformService.GetPlatform() {
	case platform.PlatformDarwin:
		return "Install Claude CLI:\n• Download from: https://claude.ai/cli\n• Or use Homebrew: brew install claude-cli\n• Ensure it's in your PATH"
	case platform.PlatformWindows:
		return "Install Claude CLI:\n• Download from: https://claude.ai/cli\n• Add to your system PATH\n• Restart your terminal after installation"
	case platform.PlatformLinux:
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
// In Claude's paradigm, an MCP is active if it's added to Claude CLI configuration
func SyncMCPStatus(model types.Model, activeMCPs []string) types.Model {
	// Create a map for quick lookup of active MCPs from Claude
	claudeActiveMCPs := make(map[string]bool)
	for _, mcpName := range activeMCPs {
		claudeActiveMCPs[mcpName] = true
	}

	// Update local MCP items based on Claude's active MCPs
	// This implements the mapping: Claude "added" = Local TestActiveStatus
	for i := range model.MCPItems {
		mcpName := model.MCPItems[i].Name

		// Set active status based on Claude's configuration
		// If MCP is in Claude's list, it's active (added)
		// If MCP is not in Claude's list, it's inactive (not added/removed)
		model.MCPItems[i].Active = claudeActiveMCPs[mcpName]
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

// ToggleMCPStatus toggles the active status of an MCP using Claude CLI add/remove commands
// This method maps the toggle concept to Claude's add/remove paradigm
func (cs *ClaudeService) ToggleMCPStatus(ctx context.Context, mcpName string, activate bool, mcpConfig *types.MCPItem) (*ToggleResult, error) {
	start := time.Now()
	result := cs.initializeToggleResult(mcpName, activate)

	// First check if Claude CLI is available
	if !cs.validateClaudeAvailability(ctx, result, start) {
		return result, nil
	}

	// Attempt the toggle operation with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, cs.timeout)
	defer cancel()

	cmd, err := cs.buildToggleCommand(timeoutCtx, mcpName, activate, mcpConfig, result, start)
	if err != nil {
		return result, err
	}

	return cs.executeToggleCommand(cmd, result, start)
}

func (cs *ClaudeService) initializeToggleResult(mcpName string, activate bool) *ToggleResult {
	result := &ToggleResult{
		MCPName:  mcpName,
		NewState: "deactivating",
	}

	if activate {
		result.NewState = "activating"
	}

	return result
}

func (cs *ClaudeService) validateClaudeAvailability(ctx context.Context, result *ToggleResult, start time.Time) bool {
	status := cs.DetectClaudeCLI(ctx)
	if !status.Available {
		result.Success = false
		result.ErrorType = ErrorTypeClaudeUnavailable
		result.ErrorMsg = ErrorMessages[ErrorTypeClaudeUnavailable]
		result.Retryable = false
		result.Duration = time.Since(start)
		return false
	}
	return true
}

func (cs *ClaudeService) buildToggleCommand(ctx context.Context, mcpName string, activate bool, mcpConfig *types.MCPItem, result *ToggleResult, start time.Time) (*exec.Cmd, error) {
	if activate {
		return cs.buildActivateCommand(ctx, mcpConfig, result, start)
	}
	return cs.buildDeactivateCommand(ctx, mcpName, result), nil
}

func (cs *ClaudeService) buildActivateCommand(ctx context.Context, mcpConfig *types.MCPItem, result *ToggleResult, start time.Time) (*exec.Cmd, error) {
	if mcpConfig == nil {
		result.Success = false
		result.ErrorType = ErrorTypeUnknownError
		result.ErrorMsg = "MCP configuration required for activation"
		result.Duration = time.Since(start)
		return nil, nil
	}

	cmd, cmdErr := cs.buildAddCommand(ctx, mcpConfig)
	if cmdErr != nil {
		result.Success = false
		result.ErrorType = ErrorTypeUnknownError
		result.ErrorMsg = "Invalid MCP configuration: " + cmdErr.Error()
		result.Duration = time.Since(start)
		return nil, cmdErr
	}

	result.NewState = TestActiveStatus
	return cmd, nil
}

func (cs *ClaudeService) buildDeactivateCommand(ctx context.Context, mcpName string, result *ToggleResult) *exec.Cmd {
	result.NewState = "inactive"
	return exec.CommandContext(ctx, ClaudeCommand, "mcp", "remove", mcpName)
}

func (cs *ClaudeService) executeToggleCommand(cmd *exec.Cmd, result *ToggleResult, start time.Time) (*ToggleResult, error) {
	output, err := cmd.CombinedOutput()
	result.Duration = time.Since(start)

	if err != nil {
		return cs.handleToggleError(err, string(output), result)
	}

	result.Success = true
	return result, nil
}

func (cs *ClaudeService) handleToggleError(err error, output string, result *ToggleResult) (*ToggleResult, error) {
	result.Success = false
	result.ErrorType = cs.classifyError(err, output)
	result.ErrorMsg = ErrorMessages[result.ErrorType]
	result.Retryable = (result.ErrorType == ErrorTypeNetworkTimeout)

	// If retryable and within time budget, mark as retrying
	// The actual retry will be handled by the UI layer with proper delay
	if result.Retryable && result.Duration < 8*time.Second {
		result.Retrying = true
		result.ErrorMsg = "MCP toggle timed out. Retrying..."
	}

	return result, nil
}

// retryToggleOperation performs a single retry of the toggle operation
func (cs *ClaudeService) retryToggleOperation(ctx context.Context, mcpName string, activate bool, mcpConfig *types.MCPItem, originalStart time.Time) (*ToggleResult, error) {
	result := &ToggleResult{
		MCPName:  mcpName,
		NewState: "deactivating",
	}

	if activate {
		result.NewState = "activating"
	}

	// Check remaining time budget (max 20 seconds total)
	elapsed := time.Since(originalStart)
	if elapsed > 18*time.Second {
		result.Success = false
		result.ErrorType = ErrorTypeNetworkTimeout
		result.ErrorMsg = "Operation timed out"
		result.Duration = elapsed
		return result, fmt.Errorf("operation timeout")
	}

	// Perform retry with remaining time
	remainingTime := 20*time.Second - elapsed
	timeoutCtx, cancel := context.WithTimeout(ctx, remainingTime)
	defer cancel()

	var cmd *exec.Cmd
	if activate {
		// Use 'claude mcp add' command with proper configuration
		if mcpConfig == nil {
			result.Success = false
			result.ErrorType = ErrorTypeUnknownError
			result.ErrorMsg = "MCP configuration required for activation"
			result.Duration = time.Since(originalStart)
			return result, fmt.Errorf("missing MCP configuration")
		}
		var cmdErr error
		cmd, cmdErr = cs.buildAddCommand(timeoutCtx, mcpConfig)
		if cmdErr != nil {
			result.Success = false
			result.ErrorType = ErrorTypeUnknownError
			result.ErrorMsg = "Invalid MCP configuration: " + cmdErr.Error()
			result.Duration = time.Since(originalStart)
			return result, cmdErr
		}
		result.NewState = TestActiveStatus
	} else {
		// Use 'claude mcp remove' command
		cmd = exec.CommandContext(timeoutCtx, ClaudeCommand, "mcp", "remove", mcpName)
		result.NewState = "inactive"
	}

	_, err := cmd.CombinedOutput()
	result.Duration = time.Since(originalStart)

	if err != nil {
		result.Success = false
		result.ErrorType = ErrorTypeUnknownError
		result.ErrorMsg = ErrorMessages[ErrorTypeUnknownError]
		return result, err
	}

	result.Success = true
	return result, nil
}

// validateMCPConfig validates MCP configuration to prevent command injection
func (cs *ClaudeService) validateMCPConfig(mcpConfig *types.MCPItem) error {
	// Validate MCP name - should only contain safe characters
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !nameRegex.MatchString(mcpConfig.Name) {
		return fmt.Errorf("invalid MCP name: contains unsafe characters")
	}

	// Validate environment variable keys
	envKeyRegex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	for key := range mcpConfig.Environment {
		if !envKeyRegex.MatchString(key) {
			return fmt.Errorf("invalid environment variable key: %s", key)
		}
	}

	return nil
}

// buildAddCommand constructs the claude mcp add command based on MCP configuration
func (cs *ClaudeService) buildAddCommand(ctx context.Context, mcpConfig *types.MCPItem) (*exec.Cmd, error) {
	// Validate configuration to prevent command injection
	if err := cs.validateMCPConfig(mcpConfig); err != nil {
		return nil, err
	}

	args := []string{"mcp", "add", mcpConfig.Name}

	// Add the command or URL
	if mcpConfig.URL != "" {
		args = append(args, mcpConfig.URL)
	} else {
		args = append(args, mcpConfig.Command)
	}

	// Add command arguments
	if len(mcpConfig.Args) > 0 {
		args = append(args, mcpConfig.Args...)
	}

	// Determine transport type based on MCP type
	switch strings.ToUpper(mcpConfig.Type) {
	case "SSE":
		args = append(args, "-t", "sse")
	case "HTTP":
		args = append(args, "-t", "http")
	default:
		// stdio is the default, no need to specify
	}

	// Add environment variables if present
	if len(mcpConfig.Environment) > 0 {
		for key, value := range mcpConfig.Environment {
			args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Validate command name to prevent command injection
	if ClaudeCommand != "claude" {
		return nil, fmt.Errorf("invalid command name")
	}

	// Use a secure command executor
	return createSecureCommand(ctx, ClaudeCommand, args...)
}

// createSecureCommand creates a command only for allowed commands
func createSecureCommand(ctx context.Context, cmdName string, args ...string) (*exec.Cmd, error) {
	// Check if command is in allowlist
	if !allowedCommands[cmdName] {
		return nil, fmt.Errorf("command not allowed: %s", cmdName)
	}

	// For extra security, ensure we only allow the specific command we expect
	switch cmdName {
	case "claude":
		return exec.CommandContext(ctx, "claude", args...), nil
	default:
		return nil, fmt.Errorf("unknown command: %s", cmdName)
	}
}

// classifyError classifies the error type based on the error and output
func (cs *ClaudeService) classifyError(err error, output string) string {
	outputLower := strings.ToLower(output)
	errMsg := strings.ToLower(err.Error())

	if cs.isTimeoutError(errMsg) {
		return ErrorTypeNetworkTimeout
	}

	if cs.isClaudeUnavailableError(errMsg, outputLower) {
		return ErrorTypeClaudeUnavailable
	}

	if cs.isMCPAlreadyExistsError(outputLower) {
		return ErrorTypeMCPAlreadyExists
	}

	if cs.isMCPNotFoundError(outputLower) {
		return ErrorTypeMCPNotFound
	}

	if cs.isInvalidCommandError(outputLower) {
		return ErrorTypeInvalidCommand
	}

	if cs.isPermissionError(outputLower) {
		return ErrorTypePermissionError
	}

	return ErrorTypeUnknownError
}

// isTimeoutError checks for timeout-related errors
func (cs *ClaudeService) isTimeoutError(errMsg string) bool {
	return strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline exceeded")
}

// isClaudeUnavailableError checks for Claude CLI availability issues
func (cs *ClaudeService) isClaudeUnavailableError(errMsg, outputLower string) bool {
	return strings.Contains(errMsg, "executable file not found") ||
		strings.Contains(outputLower, "command not found") ||
		strings.Contains(outputLower, "not recognized")
}

// isMCPAlreadyExistsError checks for MCP already exists errors
func (cs *ClaudeService) isMCPAlreadyExistsError(outputLower string) bool {
	return strings.Contains(outputLower, "already exists") ||
		strings.Contains(outputLower, "already configured") ||
		strings.Contains(outputLower, "duplicate")
}

// isMCPNotFoundError checks for MCP not found errors
func (cs *ClaudeService) isMCPNotFoundError(outputLower string) bool {
	return (strings.Contains(outputLower, "mcp server") && strings.Contains(outputLower, "not found")) ||
		strings.Contains(outputLower, "no mcp server found") ||
		strings.Contains(outputLower, "does not exist") ||
		strings.Contains(outputLower, "unknown server")
}

// isInvalidCommandError checks for invalid command/configuration errors
func (cs *ClaudeService) isInvalidCommandError(outputLower string) bool {
	return strings.Contains(outputLower, "invalid command") ||
		strings.Contains(outputLower, "invalid configuration") ||
		strings.Contains(outputLower, "malformed") ||
		strings.Contains(outputLower, "invalid transport")
}

// isPermissionError checks for permission-related errors
func (cs *ClaudeService) isPermissionError(outputLower string) bool {
	return strings.Contains(outputLower, "permission denied") ||
		strings.Contains(outputLower, "unauthorized") ||
		strings.Contains(outputLower, "authentication")
}

// GetProjectContext returns the current project context information
func GetProjectContext(model types.Model) types.ProjectContext {
	currentPath, err := os.Getwd()
	if err != nil {
		currentPath = UnknownStatus
	}

	// Calculate display path (truncated for UI)
	displayPath := FormatPathForDisplay(currentPath, 50)

	// Count active MCPs
	activeMCPs := 0
	for _, mcp := range model.MCPItems {
		if mcp.Active {
			activeMCPs++
		}
	}

	// Determine sync status
	syncStatus := GetSyncStatus(model)
	syncStatusText := FormatSyncStatusText(syncStatus)

	return types.ProjectContext{
		CurrentPath:    currentPath,
		LastSyncTime:   model.LastClaudeSync,
		ActiveMCPs:     activeMCPs,
		TotalMCPs:      len(model.MCPItems),
		SyncStatus:     syncStatus,
		DisplayPath:    displayPath,
		SyncStatusText: syncStatusText,
	}
}

// FormatPathForDisplay formats a path for display by truncating if necessary
func FormatPathForDisplay(path string, maxLength int) string {
	if len(path) <= maxLength {
		return path
	}

	// Try to show relative path from home directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		if relPath, err := filepath.Rel(homeDir, path); err == nil && !strings.HasPrefix(relPath, "..") {
			shortPath := "~/" + relPath
			if len(shortPath) <= maxLength {
				return shortPath
			}
		}
	}

	// If still too long, truncate with ellipsis
	if len(path) > maxLength {
		return "..." + path[len(path)-maxLength+3:]
	}

	return path
}

// GetSyncStatus determines the sync status between local and Claude state
func GetSyncStatus(model types.Model) types.SyncStatus {
	if !model.ClaudeAvailable {
		return types.SyncStatusError
	}

	if model.ClaudeSyncError != "" {
		return types.SyncStatusError
	}

	if model.LastClaudeSync.IsZero() {
		return types.SyncStatusUnknown
	}

	// Check if local active MCPs match Claude's active MCPs
	localActiveMCPs := make(map[string]bool)
	for _, mcp := range model.MCPItems {
		if mcp.Active {
			localActiveMCPs[mcp.Name] = true
		}
	}

	claudeActiveMCPs := make(map[string]bool)
	for _, mcpName := range model.ClaudeStatus.ActiveMCPs {
		claudeActiveMCPs[mcpName] = true
	}

	// Compare local and Claude active MCPs
	if len(localActiveMCPs) != len(claudeActiveMCPs) {
		return types.SyncStatusOutOfSync
	}

	for mcpName := range localActiveMCPs {
		if !claudeActiveMCPs[mcpName] {
			return types.SyncStatusOutOfSync
		}
	}

	return types.SyncStatusInSync
}

// FormatSyncStatusText formats the sync status as human-readable text
func FormatSyncStatusText(syncStatus types.SyncStatus) string {
	switch syncStatus {
	case types.SyncStatusInSync:
		return "In Sync"
	case types.SyncStatusOutOfSync:
		return "Out of Sync"
	case types.SyncStatusError:
		return "Error"
	case types.SyncStatusUnknown:
		return UnknownStatus
	default:
		return UnknownStatus
	}
}

// UpdateProjectContext updates the project context in the model
func UpdateProjectContext(model types.Model) types.Model {
	model.ProjectContext = GetProjectContext(model)
	return model
}

// HasDirectoryChanged checks if the current directory has changed
func HasDirectoryChanged(currentPath string) bool {
	if currentPath == "" {
		return false
	}

	actualPath, err := os.Getwd()
	if err != nil {
		return false
	}

	return currentPath != actualPath
}
