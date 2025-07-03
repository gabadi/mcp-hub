# Story 2.5: Project Context Display

## Status: Review

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

- [x] Task 1: Status Bar Project Path Display (AC: 1)
  - [x] Implement current working directory detection
  - [x] Add project path to status bar component
  - [x] Handle long path truncation for display
  - [x] Show relative path from home directory when appropriate

- [x] Task 2: Last Sync Timestamp Implementation (AC: 2)
  - [x] Track last successful Claude CLI sync time
  - [x] Display human-readable timestamp format
  - [x] Update timestamp on successful sync operations
  - [x] Handle cases where no sync has occurred

- [x] Task 3: Active MCP Count Display (AC: 3)
  - [x] Calculate active vs total MCP counts
  - [x] Format display as "active/total active" pattern
  - [x] Update count in real-time with MCP state changes
  - [x] Handle edge cases (no MCPs, all active, etc.)

- [x] Task 4: Directory Change Detection (AC: 4)
  - [x] Implement directory change monitoring
  - [x] Update project context when directory changes
  - [x] Refresh Claude CLI status for new directory
  - [x] Handle cases where user changes directories externally

- [x] Task 5: Sync Status Indication (AC: 5)
  - [x] Implement sync status tracking
  - [x] Display visual indicators for sync state
  - [x] Show warnings when local vs Claude state differs
  - [x] Provide clear guidance for resolving sync issues

- [x] Task 6: Integration and Testing
  - [x] Update existing tests for new status components
  - [x] Add unit tests for context display logic
  - [x] Create integration tests for directory monitoring
  - [x] Test various project contexts and edge cases

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

## Story Wrap-up

### Implementation Summary
Successfully implemented all 5 acceptance criteria for project context display:
- **Status bar project path**: Implemented with intelligent truncation and home directory shortening
- **Last sync timestamp**: Human-readable format with relative time display
- **Active MCP count**: Shows "X/Y MCPs" format with real-time updates
- **Directory change detection**: 5-second polling with automatic context updates
- **Sync status indication**: Clear indicators for In Sync, Out of Sync, Error, and Unknown states

### Technical Implementation
- **New types**: ProjectContext struct, SyncStatus enum, and related message types
- **Core functions**: GetProjectContext(), FormatPathForDisplay(), GetSyncStatus()
- **UI integration**: Enhanced footer component with priority system
- **Quality assurance**: Comprehensive test coverage with 100% test pass rate

### Agent Model Used
BMAD story-simple workflow with Dev agent implementation and ht-mcp validation

### Files Modified
- `internal/ui/components/footer.go` - Enhanced with project context display
- `internal/ui/services/claude_service.go` - Added project context tracking
- `internal/ui/types/models.go` - Added ProjectContext and SyncStatus types
- `internal/ui/model.go` - Integrated project context state management
- Test files updated with comprehensive coverage

### Development Notes
- Implementation follows established project patterns and conventions
- All project quality gates passed (build, test, lint)
- Directory monitoring optimized for performance with 5-second polling
- Seamless integration with existing Claude CLI workflows

### Story Status: Ready for Review
All tasks completed, tests passing, and implementation validated with ht-mcp integration.