// Package ui provides the TUI application model and business logic for the MCP manager.
package ui

import (
	"fmt"
	"os"
	"time"

	"mcp-hub/internal/platform"
	"mcp-hub/internal/ui/handlers"
	"mcp-hub/internal/ui/services"
	"mcp-hub/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

// Model is a wrapper around the types.Model to provide UI-specific methods
type Model struct {
	types.Model
	PlatformService platform.PlatformService
}

// NewModel creates a new application model with inventory loaded from storage
func NewModel() Model {
	// Create platform service
	platformService := platform.NewPlatformServiceFactoryDefault().CreatePlatformService()
	
	// Try to load inventory from storage
	mcpItems, err := services.LoadInventory(platformService)
	var model Model

	switch {
	case err != nil:
		// Fall back to default model if loading fails
		model = Model{
			Model: types.NewModel(platformService),
			PlatformService: platformService,
		}
	case len(mcpItems) == 0:
		// First-time setup: save defaults to storage
		defaultModel := types.NewModel(platformService)
		if saveErr := services.SaveInventory(defaultModel.MCPItems, platformService); saveErr != nil {
			// Log error but continue - the app should still work
			// Error is already logged in SaveInventory
			// Intentionally empty - we don't want to fail app startup due to save issues
			_ = saveErr // Acknowledge error but continue
		}
		model = Model{
			Model: defaultModel,
			PlatformService: platformService,
		}
	default:
		// Use loaded inventory
		model = Model{
			Model: types.NewModelWithMCPs(mcpItems, platformService),
			PlatformService: platformService,
		}
	}

	// Initialize project context
	model.Model = services.UpdateProjectContext(model.Model)

	return model
}

// Init initializes the application and returns initial commands
func (m Model) Init() tea.Cmd {
	// Start startup loading overlay
	m.StartLoadingOverlay(types.LoadingStartup)

	// Return batch of commands for startup
	return tea.Batch(
		handlers.StartupLoadingCmd(),
		handlers.StartupLoadingTimerCmd(0),
		handlers.LoadingSpinnerCmd(types.LoadingStartup),
		ProjectContextCheckCmd(),        // Start project context monitoring
		DelayedClaudeStatusRefreshCmd(), // Auto-detect Claude CLI after UI loads
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg), nil
	case tea.KeyMsg:
		// Block input when loading overlay is active, except for exit keys
		if m.IsLoadingOverlayActive() {
			// Allow exit even during loading
			if msg.String() == "esc" || msg.String() == "ctrl+c" {
				return m.handleKeyMsg(msg)
			}
			return m, nil
		}
		return m.handleKeyMsg(msg)
	case handlers.SuccessMsg:
		return m.handleSuccessMsg(msg), nil
	case handlers.ClaudeStatusMsg:
		return m.handleClaudeStatusMsg(msg)
	case handlers.ToggleResultMsg:
		return m.handleToggleResultMsg(msg)
	case types.TimerTickMsg:
		return m.handleTimerTickMsg(msg)
	case types.LoadingProgressMsg:
		return m.handleLoadingProgressMsg(msg)
	case types.LoadingStepMsg:
		return m.handleLoadingStepMsg(msg)
	case types.LoadingSpinnerMsg:
		return m.handleLoadingSpinnerMsg(msg)
	case types.ProjectContextCheckMsg:
		return m.handleProjectContextCheckMsg(msg)
	case types.DirectoryChangeMsg:
		return m.handleDirectoryChangeMsg(msg)
	case StartClaudeDetectionMsg:
		return m.handleStartClaudeDetectionMsg(msg)
	}
	return m, nil
}

// handleWindowSizeMsg handles window resize messages
func (m Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) Model {
	m.Width = msg.Width
	m.Height = msg.Height
	m.Model = services.UpdateLayout(m.Model)
	// Update project context on window resize
	m.Model = services.UpdateProjectContext(m.Model)
	return m
}

// handleKeyMsg handles keyboard input messages
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Model, cmd = handlers.HandleKeyPress(m.Model, msg)
	return m, cmd
}

// handleSuccessMsg handles success messages
func (m Model) handleSuccessMsg(msg handlers.SuccessMsg) Model {
	m.SuccessMessage = msg.Message
	return m
}

// handleClaudeStatusMsg handles Claude status update messages
func (m Model) handleClaudeStatusMsg(msg handlers.ClaudeStatusMsg) (tea.Model, tea.Cmd) {
	// Stop Claude detection loading overlay if active
	if m.LoadingOverlay != nil && m.LoadingOverlay.Active && m.LoadingOverlay.Type == types.LoadingClaude {
		m.StopLoadingOverlay()
	}

	// Update model with Claude status
	m.Model = services.UpdateModelWithClaudeStatus(m.Model, msg.Status)

	// Sync MCP status if Claude is available and has active MCPs
	switch {
	case msg.Status.Available && len(msg.Status.ActiveMCPs) > 0:
		m.Model = services.SyncMCPStatus(m.Model, msg.Status.ActiveMCPs)
		// Save updated inventory after sync
		if err := services.SaveModelInventory(m.Model, m.PlatformService); err != nil {
			// Set error message but don't fail
			m.SuccessMessage = fmt.Sprintf("Claude status updated, but failed to save inventory: %v", err)
			m.SuccessTimer = 240 // Show error for 4 seconds
		} else {
			m.SuccessMessage = "Claude status refreshed and MCPs synced"
			m.SuccessTimer = 120 // Show success for 2 seconds
		}
	case msg.Status.Available:
		m.SuccessMessage = "Claude status refreshed"
		m.SuccessTimer = 120
	default:
		m.SuccessMessage = "Claude CLI not available"
		m.SuccessTimer = 180 // Show message for 3 seconds
	}

	// Update project context after Claude status update
	m.Model = services.UpdateProjectContext(m.Model)

	// Start timer for success message countdown (not toggle-specific, so use general timer)
	return m, handlers.TimerCmd("success_timer")
}

// handleToggleResultMsg handles toggle operation result messages
func (m Model) handleToggleResultMsg(msg handlers.ToggleResultMsg) (tea.Model, tea.Cmd) {
	// Handle enhanced toggle operation results (Epic 2 Story 2)
	var cmd tea.Cmd
	if msg.Success {
		m, cmd = m.handleToggleSuccess(msg)
	} else {
		m, cmd = m.handleToggleError(msg)
	}
	m.ToggleMCPName = msg.MCPName
	return m, cmd
}

// handleToggleSuccess handles successful toggle operations
func (m Model) handleToggleSuccess(msg handlers.ToggleResultMsg) (Model, tea.Cmd) {
	// Update local MCP status and save
	for i := range m.MCPItems {
		if m.MCPItems[i].Name == msg.MCPName {
			m.MCPItems[i].Active = msg.Activate
			break
		}
	}

	if err := services.SaveInventory(m.MCPItems, m.PlatformService); err != nil {
		m.ToggleState = types.ToggleError
		m.ToggleError = "MCP toggled but failed to save to storage"
		m.SuccessTimer = 240
		// Start timer for error state
		return m, handlers.TimerCmd("success_timer")
	}

	m.ToggleState = types.ToggleSuccess
	activationState := "deactivated"
	if msg.Activate {
		activationState = "activated"
	}
	m.SuccessMessage = fmt.Sprintf("MCP '%s' %s successfully", msg.MCPName, activationState)
	m.SuccessTimer = 120
	// Update project context after successful toggle
	m.Model = services.UpdateProjectContext(m.Model)
	// Start timer for success state
	return m, handlers.TimerCmd("success_timer")
}

// handleToggleError handles failed toggle operations
func (m Model) handleToggleError(msg handlers.ToggleResultMsg) (Model, tea.Cmd) {
	if msg.Retrying {
		m.ToggleState = types.ToggleRetrying
		m.SuccessMessage = fmt.Sprintf("MCP toggle failed, retrying: %s", msg.Error)
		m.SuccessTimer = 180
	} else {
		m.ToggleState = types.ToggleError
		m.SuccessMessage = fmt.Sprintf("MCP toggle failed: %s", msg.Error)
		m.SuccessTimer = 240
	}
	m.ToggleError = msg.Error
	// Start timer for error/retry state
	return m, handlers.TimerCmd("success_timer")
}

// handleTimerTickMsg handles timer tick messages for countdown functionality
func (m Model) handleTimerTickMsg(msg types.TimerTickMsg) (tea.Model, tea.Cmd) {
	// Only handle success timer ticks
	if msg.ID == "success_timer" && m.SuccessTimer > 0 {
		m.SuccessTimer--

		// If timer reaches 0, reset toggle state and clear success message
		if m.SuccessTimer <= 0 {
			// Reset toggle state and clear toggle MCP name
			m.ToggleState = types.ToggleIdle
			m.ToggleMCPName = ""
			m.ToggleError = ""

			// Clear success message
			m.SuccessMessage = ""

			// Timer has expired, don't continue
			return m, nil
		}

		// Continue timer countdown
		return m, handlers.TimerCmd("success_timer")
	}

	return m, nil
}

// handleLoadingProgressMsg handles loading progress messages
func (m Model) handleLoadingProgressMsg(msg types.LoadingProgressMsg) (tea.Model, tea.Cmd) {
	if msg.Done {
		// Loading is complete
		m.StopLoadingOverlay()
		return m, nil
	}

	// Update loading message
	m.UpdateLoadingMessage(msg.Message)
	return m, nil
}

// handleLoadingStepMsg handles loading step progression messages
func (m Model) handleLoadingStepMsg(msg types.LoadingStepMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case types.LoadingStartup:
		// Progress to next step and set timer for next progression
		return m, tea.Batch(
			handlers.StartupLoadingProgressCmd(msg.Step),
			handlers.StartupLoadingTimerCmd(msg.Step),
		)
	case types.LoadingRefresh:
		// Progress to next step and set timer for next progression
		return m, tea.Batch(
			handlers.RefreshLoadingProgressCmd(msg.Step),
			handlers.RefreshLoadingTimerCmd(msg.Step),
		)
	case types.LoadingClaude:
		// Progress to next step for Claude loading
		return m, tea.Batch(
			handlers.RefreshLoadingProgressCmd(msg.Step), // Reuse refresh loading for Claude
			handlers.RefreshLoadingTimerCmd(msg.Step),
		)
	}
	return m, nil
}

// handleLoadingSpinnerMsg handles loading spinner animation messages
func (m Model) handleLoadingSpinnerMsg(msg types.LoadingSpinnerMsg) (tea.Model, tea.Cmd) {
	if m.IsLoadingOverlayActive() {
		// Advance spinner animation
		m.AdvanceSpinner()

		// Continue spinner animation
		return m, handlers.LoadingSpinnerCmd(msg.Type)
	}

	return m, nil
}

// All key handling has been moved to handlers package

// Getter methods for testing

// GetColumnCount returns the current number of columns
func (m Model) GetColumnCount() int {
	return m.ColumnCount
}

// GetActiveColumn returns the currently active column index
func (m Model) GetActiveColumn() int {
	return m.ActiveColumn
}

// GetSelectedItem returns the currently selected item index
func (m Model) GetSelectedItem() int {
	return m.SelectedItem
}

// GetState returns the current application state
func (m Model) GetState() types.AppState {
	return m.State
}

// GetSearchQuery returns the current search query
func (m Model) GetSearchQuery() string {
	return m.SearchQuery
}

// GetSearchActive returns whether search is currently active
func (m Model) GetSearchActive() bool {
	return m.SearchActive
}

// GetSearchInputActive returns whether search input is currently active
func (m Model) GetSearchInputActive() bool {
	return m.SearchInputActive
}

// GetFilteredMCPs returns MCPs filtered by search query
func (m Model) GetFilteredMCPs() []types.MCPItem {
	return services.GetFilteredMCPs(m.Model)
}

// handleProjectContextCheckMsg handles periodic project context checks
func (m Model) handleProjectContextCheckMsg(_ types.ProjectContextCheckMsg) (tea.Model, tea.Cmd) {
	// Check if directory has changed
	if services.HasDirectoryChanged(m.ProjectContext.CurrentPath) {
		// Directory has changed, trigger directory change message
		newPath, err := os.Getwd()
		if err == nil {
			return m, DirectoryChangeCmd(newPath)
		}
	}

	// Update project context regardless to refresh sync status and timestamps
	m.Model = services.UpdateProjectContext(m.Model)

	// Schedule next check in 5 seconds
	return m, ProjectContextCheckCmd()
}

// handleDirectoryChangeMsg handles directory change events
func (m Model) handleDirectoryChangeMsg(_ types.DirectoryChangeMsg) (tea.Model, tea.Cmd) {
	// Update project context with new directory
	m.Model = services.UpdateProjectContext(m.Model)

	// Optionally trigger a Claude status refresh to sync with new directory
	// This ensures the MCP status is accurate for the new project context
	if m.ClaudeAvailable {
		return m, RefreshClaudeStatusCmd()
	}

	return m, nil
}

// ProjectContextCheckCmd returns a command to check project context
func ProjectContextCheckCmd() tea.Cmd {
	return tea.Tick(time.Second*5, func(_ time.Time) tea.Msg {
		return types.ProjectContextCheckMsg{}
	})
}

// DirectoryChangeCmd returns a command to signal directory change
func DirectoryChangeCmd(newPath string) tea.Cmd {
	return func() tea.Msg {
		return types.DirectoryChangeMsg{NewPath: newPath}
	}
}

// RefreshClaudeStatusCmd returns a command to refresh Claude status
func RefreshClaudeStatusCmd() tea.Cmd {
	return handlers.RefreshClaudeStatusCmd()
}

// DelayedClaudeStatusRefreshCmd returns a command to refresh Claude status after a short delay
// This allows the UI to load first before detecting Claude CLI
func DelayedClaudeStatusRefreshCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(_ time.Time) tea.Msg {
		// After 500ms, start the Claude detection process
		return StartClaudeDetectionMsg{}
	})
}

// StartClaudeDetectionMsg signals the start of Claude CLI detection
type StartClaudeDetectionMsg struct{}

// handleStartClaudeDetectionMsg handles the start of Claude detection
func (m Model) handleStartClaudeDetectionMsg(_ StartClaudeDetectionMsg) (tea.Model, tea.Cmd) {
	// Start Claude detection loading overlay
	m.StartLoadingOverlay(types.LoadingClaude)

	// Return command to refresh Claude status and spinner
	return m, tea.Batch(
		RefreshClaudeStatusCmd(),
		handlers.LoadingSpinnerCmd(types.LoadingClaude),
	)
}

// All layout and navigation logic has been moved to services and handlers packages
