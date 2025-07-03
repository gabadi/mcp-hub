# Story 2.5: Project Context Display

## Status: In Progress

## Story Approved for Development

**Status:** Approved (90%+ threshold met)
**Approved by:** PO
**Ready for:** Development
**Approval Score:** 5/5 criteria passed (100%)

## Story

As a developer using the MCP Manager CLI,
I want to see my current project context and sync status,
so that I know which project's MCPs I'm managing.

## Acceptance Criteria (ACs)

1. Status bar shows current project directory path
2. Last sync timestamp displayed
3. Active MCP count shown (e.g., "3/30 active")
4. Project context updates when changing directories
5. Clear indication when out of sync with Claude Code

## Tasks / Subtasks

- [ ] Task 1: Status Bar Project Path Display (AC: 1)
  - [ ] Implement current working directory detection
  - [ ] Add project path to status bar component
  - [ ] Handle long path truncation for display
  - [ ] Show relative path from home directory when appropriate

- [ ] Task 2: Last Sync Timestamp Implementation (AC: 2)
  - [ ] Track last successful Claude CLI sync time
  - [ ] Display human-readable timestamp format
  - [ ] Update timestamp on successful sync operations
  - [ ] Handle cases where no sync has occurred

- [ ] Task 3: Active MCP Count Display (AC: 3)
  - [ ] Calculate active vs total MCP counts
  - [ ] Format display as "active/total active" pattern
  - [ ] Update count in real-time with MCP state changes
  - [ ] Handle edge cases (no MCPs, all active, etc.)

- [ ] Task 4: Directory Change Detection (AC: 4)
  - [ ] Implement directory change monitoring
  - [ ] Update project context when directory changes
  - [ ] Refresh Claude CLI status for new directory
  - [ ] Handle cases where user changes directories externally

- [ ] Task 5: Sync Status Indication (AC: 5)
  - [ ] Implement sync status tracking
  - [ ] Display visual indicators for sync state
  - [ ] Show warnings when local vs Claude state differs
  - [ ] Provide clear guidance for resolving sync issues

- [ ] Task 6: Integration and Testing
  - [ ] Update existing tests for new status components
  - [ ] Add unit tests for context display logic
  - [ ] Create integration tests for directory monitoring
  - [ ] Test various project contexts and edge cases

## Dev Notes

### Previous Story Insights
Story 2.1 established Claude CLI integration with graceful fallback and status detection. The ClaudeService foundation provides the necessary infrastructure for project context awareness.

### Technical Context
[Source: docs/architecture.md#Service Layer]

**Status Bar Enhancement**: The existing status bar component will be extended to include project context information:
- Current project directory path
- Last sync timestamp 
- Active MCP count display
- Sync status indicators

**Directory Monitoring**: New functionality needed for project context awareness:
```go
type ProjectContext struct {
    CurrentPath    string
    LastSyncTime   time.Time
    ActiveMCPs     int
    TotalMCPs      int
    SyncStatus     SyncStatus
}

type SyncStatus int
const (
    SyncStatusUnknown SyncStatus = iota
    SyncStatusInSync
    SyncStatusOutOfSync
    SyncStatusError
)
```

**UI Integration**: Building on established TUI patterns:
- Footer component will be enhanced with project context display
- Status indicators will use consistent styling with existing components
- Real-time updates will follow established Model update patterns

### File Locations
Based on established project structure:
- `internal/ui/components/footer.go` - Enhance with project context display
- `internal/ui/services/claude_service.go` - Add project context tracking
- `internal/ui/types/models.go` - Add ProjectContext types
- `internal/ui/model.go` - Integrate project context state management

### Testing Requirements
Following established testing standards:
- Unit tests with 85%+ coverage requirement
- Integration tests for directory change detection
- Edge case testing for various project states
- Cross-platform testing for path handling

### Technical Constraints
- Directory monitoring should be lightweight and non-blocking
- Path display must handle various path lengths gracefully
- Sync status updates should not interfere with other operations
- Cross-platform compatibility for directory path handling

## Epic Integration
This story builds upon Story 2.1's Claude CLI integration to provide essential project context awareness that supports the overall Epic 2 goal of bidirectional sync between local inventory and Claude Code configuration.