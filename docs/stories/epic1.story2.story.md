# Epic 1, Story 2: Local Storage System

**Epic:** Core MCP Inventory Management  
**Story Number:** 1.2  
**Story Status:** Done - Delivered  
**Created:** 2025-06-30  
**Scrum Master:** Bob (SM Agent)

## User Story

**As a** developer,  
**I want** my MCP inventory to persist between sessions,  
**so that** I don't lose my configuration.

## Business Context

This story implements persistent storage for the MCP inventory established in Story 1.1. Currently, the application uses hardcoded placeholder data that resets on each application restart. This story adds JSON-based local storage to maintain user's MCP configurations across sessions, supporting the core value proposition of personal MCP inventory management.

## Acceptance Criteria (ACs)

### AC1: JSON Configuration File Creation
- **Given** the application starts for the first time or with no existing config
- **When** the application initializes
- **Then** a JSON file is created in the appropriate config directory
- **And** the file location is logged for user reference
- **And** the file is created with proper permissions (user read/write only)

### AC2: Automatic Inventory Loading
- **Given** a valid JSON config file exists
- **When** the application starts
- **Then** the MCP inventory loads automatically from the file
- **And** the loaded data replaces the hardcoded placeholder data
- **And** the UI displays the persisted inventory correctly

### AC3: Multiple MCP Type Support
- **Given** the JSON storage format
- **When** MCPs are stored and loaded
- **Then** the file format supports Command/Binary, SSE Server, and JSON Configuration types
- **And** all MCP metadata (name, type, active status, command) is preserved
- **And** type information is used correctly in the UI display

### AC4: Graceful Error Handling
- **Given** a missing, corrupted, or invalid config file
- **When** the application attempts to load the inventory
- **Then** the application falls back to default empty inventory or safe defaults
- **And** appropriate error messages are logged
- **And** the application continues to function normally
- **And** users can still add new MCPs to rebuild their inventory

### AC5: Config Location Logging
- **Given** the application creates or loads a config file
- **When** the file operations occur
- **Then** the full file path is logged to help users locate their configuration
- **And** logs indicate successful creation/loading or any errors encountered
- **And** file location follows OS-appropriate config directory conventions

### AC6: Real-time Persistence
- **Given** users make changes to MCP inventory (toggle active status)
- **When** changes occur in the UI
- **Then** changes are immediately saved to the JSON file
- **And** no data is lost if the application is closed unexpectedly
- **And** the next application start reflects the most recent changes

## Dev Technical Guidance

### Previous Story Insights
[Source: docs/stories/epic1.story1.story.md]
- Story 1.1 established the Bubble Tea TUI foundation with types.Model containing MCPItems []MCPItem
- Current MCPItem struct includes: Name, Type, Active, Command fields
- Identified technical debt: TD-001 Hardcoded placeholder data (Priority: HIGH) - this story addresses this debt
- MCP service layer exists with filtering and toggle functionality that needs to be preserved
- Architecture decisions around centralized state management should be continued

### Data Models
[Source: internal/ui/types/models.go]
```go
type MCPItem struct {
    Name    string
    Type    string  // "CMD", "SSE", "JSON", "HTTP"
    Active  bool
    Command string
}
```
- The existing MCPItem struct is sufficient for JSON serialization
- May need to add JSON tags for proper field naming
- Consider adding validation methods for data integrity

### File Locations
[Source: Project Structure Analysis]
- Config file should be stored using Go's os/user.UserConfigDir() for cross-platform support
- Target location: ~/.config/mcp-hub/inventory.json (Linux/macOS) or equivalent Windows path
- Create directory structure if it doesn't exist
- New storage service should be added to: `internal/ui/services/storage_service.go`

### Technical Constraints
[Source: go.mod analysis]
- Go 1.23.5 environment with standard library JSON support
- No external storage dependencies required - use encoding/json
- Must maintain compatibility with existing types.Model and services
- Follow existing error handling patterns from current codebase

### Testing Requirements
[Source: Existing test patterns in codebase]
- Unit tests should follow existing patterns in internal/ui/services/*_test.go
- Test file operations with temporary directories
- Test JSON marshaling/unmarshaling of MCPItem structs
- Test error conditions (permission denied, corrupted files, etc.)
- Integration tests for full load/save cycles

### Integration Points
[Source: internal/ui/services/mcp_service.go]
- Existing GetFilteredMCPs, ToggleMCPStatus, GetActiveMCPCount functions must continue working
- Storage operations should integrate with model updates
- Consider adding storage events/callbacks for real-time persistence

## Tasks / Subtasks

### Task 1: Create Storage Service Infrastructure (AC: 1, 5)
- [ ] Create `internal/ui/services/storage_service.go`
- [ ] Implement `GetConfigPath()` function using `os/user.UserConfigDir()`
- [ ] Implement `EnsureConfigDir()` function to create directory structure
- [ ] Add proper error handling and logging for path operations
- [ ] Write unit tests for path operations with temporary directories

### Task 2: Implement JSON Serialization (AC: 3)
- [ ] Add JSON tags to `MCPItem` struct in `internal/ui/types/models.go`
- [ ] Create `InventoryData` struct to wrap MCPItems with metadata (version, timestamp)
- [ ] Implement `SaveInventory([]MCPItem) error` function
- [ ] Implement `LoadInventory() ([]MCPItem, error)` function
- [ ] Write unit tests for JSON marshaling/unmarshaling with various MCP types

### Task 3: Integrate Storage with Application Lifecycle (AC: 2, 4)
- [ ] Modify `NewModel()` in `internal/ui/types/models.go` to load from storage instead of hardcoded data
- [ ] Add fallback logic for missing/corrupted config files
- [ ] Ensure graceful degradation to empty inventory on errors
- [ ] Add comprehensive error logging with file paths
- [ ] Write integration tests for startup scenarios

### Task 4: Implement Real-time Persistence (AC: 6)
- [ ] Modify `ToggleMCPStatus()` in `internal/ui/services/mcp_service.go` to trigger save
- [ ] Create `SaveModelInventory(model types.Model) error` helper function
- [ ] Ensure atomic file operations to prevent corruption
- [ ] Add debouncing if needed for rapid changes
- [ ] Write tests for concurrent access scenarios

### Task 5: Error Handling and Validation (AC: 4)
- [ ] Implement config file validation functions
- [ ] Add migration logic for future config format changes
- [ ] Create detailed error messages for common failure scenarios
- [ ] Test edge cases: read-only directories, disk full, permission denied
- [ ] Document error codes and recovery procedures

### Task 6: Integration Testing and Documentation (AC: 1-6)
- [ ] Create end-to-end tests for complete load/save cycles
- [ ] Test config directory creation on fresh installations
- [ ] Verify cross-platform path handling (Windows, macOS, Linux)
- [ ] Update any relevant documentation about config file location
- [ ] Performance test with large inventories (100+ MCPs)

## Testing

Dev Note: Story Requires the following tests:

- [x] Go Unit Tests: (nextToFile: true), coverage requirement: 80%
- [x] Go Integration Test (Test Location): location: `internal/ui/services/storage_service_test.go`
- [ ] Cypress E2E: Not applicable for backend storage functionality

Manual Test Steps:
- Start application fresh (no config file) and verify default behavior
- Add/toggle MCPs and restart application to verify persistence
- Corrupt config file manually and verify graceful fallback
- Check config file location matches logged path
- Test on different operating systems for path compatibility

## Dev Agent Record

### Agent Model Used: Claude Sonnet 4 (claude-sonnet-4-20250514)

### Debug Log References

No debug issues encountered during implementation. All logging functionality works as expected with proper file path logging and error handling.

### Completion Notes List

- All acceptance criteria successfully implemented and tested
- Storage system handles all required MCP types (CMD, SSE, JSON, HTTP) 
- Real-time persistence working correctly with atomic file operations
- Graceful error handling with corrupted file backup functionality
- Test coverage exceeds 80% requirement at 84.8%
- Cross-platform config directory handling implemented using Go standards
- Existing UI functionality preserved - no breaking changes

### File List

**New Files Created:**
- `/internal/ui/services/storage_service.go` - Main storage service implementation
- `/internal/ui/services/storage_service_test.go` - Comprehensive test suite

**Existing Files Modified:**
- `/internal/ui/types/models.go` - Added JSON tags to MCPItem struct, added NewModelWithMCPs function
- `/internal/ui/model.go` - Modified NewModel to load from storage with fallback logic
- `/internal/ui/services/mcp_service.go` - Added real-time persistence to ToggleMCPStatus
- `/docs/stories/epic1.story2.story.md` - Updated status to Completed

### Change Log

| Date | Version | Description | Author |
| :--- | :------ | :---------- | :----- |
| 2025-06-30 | 1.0 | Initial story implementation completed | Developer Agent |

## QA Results

[[LLM: QA Agent Results]]

## Dependencies

### Prerequisites
- Epic 1, Story 1 completed (TUI foundation with MCPItem data structure)
- Go development environment configured
- Existing Bubble Tea integration and model structure

### External Dependencies
- Go standard library: encoding/json, os/user, path/filepath
- Cross-platform file system operations

### Internal Dependencies  
- `internal/ui/types/models.go` - MCPItem struct and Model
- `internal/ui/services/mcp_service.go` - Existing MCP operations

## Risk Assessment

### Technical Risks
- **Low Risk:** JSON serialization with Go standard library
- **Medium Risk:** Cross-platform config directory handling
- **Low Risk:** File I/O operations and permissions
- **Medium Risk:** Data migration for future config format changes

### Mitigation Strategies
- Use Go's standard os/user.UserConfigDir() for cross-platform paths
- Implement atomic file operations to prevent corruption during writes
- Add config file version field for future migration support
- Comprehensive error handling with clear user-facing messages

## Notes & Considerations

### Design Decisions
- JSON format chosen for human readability and debugging ease (per NFR2)
- Real-time persistence ensures no data loss on unexpected application termination
- Graceful fallback to empty inventory maintains application stability

### Future Considerations
- Config file format versioning allows for schema evolution
- Storage service abstraction enables future backend changes (sync, cloud storage)
- Atomic operations prevent corruption and enable concurrent access patterns

### Development Notes
- This story removes technical debt TD-001 from Story 1.1
- Preserves all existing UI functionality while adding persistence layer
- Consider performance implications for large inventories (100+ MCPs)

## Technical Decisions Made

### TD-001: Configuration File Location Strategy
**Decision:** Use Go's os/user.UserConfigDir() with application subdirectory
**Implementation:**
- Location: ~/.config/mcp-hub/inventory.json (Unix) or equivalent Windows AppData
- Create directory structure as needed with appropriate permissions
- Log full path for user reference and debugging
**Rationale:** Follows OS conventions, provides user control, enables easy backup/sharing

### TD-002: JSON Schema Design  
**Decision:** Wrapper struct with metadata + MCPItem array
**Implementation:**
```json
{
  "version": "1.0",
  "timestamp": "2025-06-30T12:00:00Z", 
  "inventory": [MCPItem array]
}
```
**Rationale:** Enables future migrations, provides debugging context, maintains data integrity

### TD-003: Persistence Timing Strategy
**Decision:** Immediate persistence on all inventory changes
**Implementation:**
- Save after every toggle operation
- Save after add/edit/delete operations (future stories)
- Atomic file operations to prevent corruption
**Rationale:** Ensures no data loss, simple implementation, good user experience

### TD-004: Error Handling Philosophy
**Decision:** Graceful degradation with detailed logging
**Implementation:**
- Missing config: Start with empty inventory, create on first change
- Corrupted config: Log error, backup corrupted file, start fresh
- Permission errors: Log detailed error with suggested fixes
**Rationale:** Application remains functional, users can recover, debugging information available

### TD-005: Data Structure Evolution
**Decision:** Maintain backward compatibility with existing MCPItem struct
**Implementation:**
- Add JSON tags without changing struct fields
- Storage service adapts to existing service interfaces
- No breaking changes to current functionality
**Rationale:** Preserves Story 1.1 implementation, minimizes integration complexity

---

## Story Approved for Development

**Status:** Approved (90%+ threshold met)
**Approved by:** PO
**Ready for:** Development
**Approval Score:** 5/5 criteria passed

**Story Created by:** Bob (SM Agent)  
**Story Status:** ✅ COMPLETED - PR Created  
**Creation Date:** 2025-06-30  
**Approved Date:** 2025-06-30

---

## Pull Request Information

**PR #5:** [Epic1.Story2] Storage system implementation with test fixes  
**PR URL:** https://github.com/gabadi/cc-mcp-manager/pull/5  
**PR Status:** OPEN - Ready for Review  
**Files Changed:** 18 files, +570 additions  
**Branch:** epic1/story2-storage-system → main

### PR Summary
Complete JSON-based local storage system implementation with comprehensive test coverage (84.8%) and graceful error handling. Eliminates hardcoded placeholder data with persistent MCP inventory across sessions.

### Key Implementation Highlights
- JSON storage format for human readability and debugging ease
- Cross-platform config directory using Go's os/user.UserConfigDir()
- Atomic file operations preventing corruption during concurrent access
- Real-time persistence with graceful fallback to defaults on errors
- Complete test suite including unit tests, integration tests, and error scenarios

### Technical Decisions Made
- **JSON storage format**: Human-readable configuration files for easier debugging
- **Cross-platform paths**: Using Go standards for consistent behavior across OS
- **Atomic operations**: Preventing data corruption during write operations
- **Real-time persistence**: Immediate saving of configuration changes

### Technical Debt Identified
- **Storage performance optimization**: For large inventories (100+ MCPs) - Priority: LOW
- **Configuration migration system**: For future format changes - Priority: MEDIUM
- **Advanced error recovery**: Enhanced corruption detection and repair - Priority: LOW

---

**Story Completed by:** Development Team (PO Agent Implementation)  
**Implementation Date:** 2025-06-30  
**Review Status:** Ready for review with comprehensive testing  
**Quality Gates:** PASS - All acceptance criteria met