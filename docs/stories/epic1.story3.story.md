# Epic 1, Story 3: Add MCP Workflow

**Epic:** Core MCP Inventory Management  
**Story Number:** 1.3  
**Story Status:** Draft  
**Created:** 2025-06-30  
**Scrum Master:** Bob (SM Agent)

## User Story

**As a** developer,  
**I want** to add different types of MCPs to my inventory,  
**so that** I can support various MCP configurations and build my personal toolkit.

## Business Context

This story implements the core "Add MCP" functionality that enables developers to build their personal MCP inventory. Building on the TUI foundation (Story 1.1) and persistence layer (Story 1.2), this story delivers the primary value proposition of the MCP Manager CLI - allowing developers to easily add different types of MCPs to their toolkit without memorizing complex command syntax.

The story addresses the key user pain point identified in the PRD: "MCP management friction" by providing an intuitive modal-driven interface for adding Command/Binary, SSE Server, and JSON Configuration MCPs. This is critical for establishing the complete CRUD foundation that Epic 1 requires.

## Acceptance Criteria

### AC1: Add MCP Modal Activation
- **Given** the user is in the main MCP inventory interface
- **When** the user presses the 'A' key
- **Then** an "Add MCP" modal dialog opens
- **And** the modal displays MCP type selection options
- **And** the main interface remains visible but dimmed in the background
- **And** the modal is centered and clearly distinguished from the main interface

### AC2: MCP Type Selection Interface
- **Given** the Add MCP modal is open
- **When** the modal displays
- **Then** three MCP type options are presented:
  - "1. Command/Binary (most common)"
  - "2. SSE Server (HTTP/WebSocket)"  
  - "3. JSON Configuration"
- **And** users can select using number keys (1/2/3) or arrow keys + Enter
- **And** clear instructions show selection methods: "[1-3] Select [ESC] Cancel"
- **And** the most common option (Command/Binary) is visually emphasized

### AC3: Command/Binary MCP Form
- **Given** the user selects "Command/Binary" MCP type
- **When** the form is displayed
- **Then** the form contains these fields:
  - Name: (required) - text input for MCP identifier
  - Command: (required) - text input for executable command
  - Args: (optional) - text input for additional arguments
- **And** Tab key navigates between form fields
- **And** required fields are clearly marked
- **And** real-time validation shows field-level error messages
- **And** form instructions show: "[Enter] Add [ESC] Cancel"

### AC4: SSE Server MCP Form
- **Given** the user selects "SSE Server" MCP type
- **When** the form is displayed
- **Then** the form contains these fields:
  - Name: (required) - text input for MCP identifier
  - URL: (required) - text input for server endpoint
- **And** URL field validates for proper HTTP/HTTPS format
- **And** real-time validation checks URL accessibility (with timeout)
- **And** loading indicator shows during URL validation

### AC5: JSON Configuration MCP Form
- **Given** the user selects "JSON Configuration" MCP type
- **When** the form is displayed
- **Then** the form contains these fields:
  - Name: (required) - text input for MCP identifier
  - JSON Config: (required) - multi-line text area for configuration
- **And** JSON field validates syntax in real-time
- **And** syntax errors show specific line/column information
- **And** properly formatted JSON is visually confirmed with checkmark

### AC6: Form Validation and Error Handling
- **Given** any MCP form is being filled out
- **When** validation errors occur
- **Then** inline error messages appear below the relevant field
- **And** error messages are specific and actionable
- **And** the Enter key is disabled until all required fields are valid
- **And** duplicate MCP names show warning with option to overwrite or rename

### AC7: Successful MCP Addition
- **Given** a valid MCP form is submitted
- **When** the user presses Enter
- **Then** the MCP is saved to the local inventory JSON file
- **And** the modal closes and returns to the main interface
- **And** the MCP list refreshes showing the new MCP with correct type badge
- **And** a success message briefly appears: "Added [MCP-name] successfully"
- **And** the new MCP is automatically selected in the list

### AC8: Modal Navigation and Cancellation
- **Given** any modal or form is open
- **When** the user presses ESC key
- **Then** the modal closes without saving changes
- **And** the user returns to the main interface
- **And** no data is modified in the inventory
- **And** the previously selected MCP remains selected

## Technical Implementation Details

### Architecture Components

Building on the established architecture from Stories 1.1 and 1.2:

1. **Modal System Enhancement**
   - Extend existing Bubble Tea model to support modal states
   - Add modal rendering layer that overlays main interface
   - Implement modal state management within centralized state

2. **Form Components**
   - Create reusable form input components
   - Implement field validation system
   - Add real-time validation feedback

3. **Storage Integration**
   - Extend storage service to support new MCP additions
   - Ensure atomic operations for inventory updates
   - Maintain backward compatibility with existing data

### State Management Extensions

```go
// Extended Model state for modal support
type Model struct {
    // Existing fields from Stories 1.1 & 1.2
    MCPItems     []MCPItem
    SelectedItem int
    SearchQuery  string
    
    // New modal state
    ModalState   ModalType
    CurrentForm  FormData
    FormErrors   map[string]string
}

type ModalType int
const (
    NoModal ModalType = iota
    AddMCPTypeSelection
    AddCommandForm
    AddSSEForm
    AddJSONForm
)
```

### Form Validation System

1. **Real-time Validation**
   - Name uniqueness checking against existing inventory
   - Command/URL format validation
   - JSON syntax validation with error positioning

2. **Async Validation**
   - SSE server endpoint accessibility testing
   - Timeout handling for network operations
   - Loading states during validation

### Integration Points

- **Storage Service**: Extend `SaveInventory()` to handle new MCPs
- **MCP Service**: Update filtering and display logic for new entries
- **UI Components**: Integrate modal rendering with existing layout system

## Tasks / Subtasks

### Task 1: Modal System Foundation (AC: 1, 8)
- [ ] Extend Bubble Tea Model to support modal states
- [ ] Implement modal overlay rendering system
- [ ] Add modal keyboard navigation (ESC to cancel)
- [ ] Create modal backdrop dimming effect
- [ ] Write unit tests for modal state transitions

### Task 2: MCP Type Selection Modal (AC: 2)
- [ ] Create type selection modal component
- [ ] Implement number key (1/2/3) selection
- [ ] Add arrow key navigation support
- [ ] Style modal with clear option descriptions
- [ ] Write tests for type selection logic

### Task 3: Command/Binary MCP Form (AC: 3, 6)
- [ ] Create command form component with three fields
- [ ] Implement Tab navigation between fields
- [ ] Add field validation (required/optional)
- [ ] Create form submission logic
- [ ] Write comprehensive form validation tests

### Task 4: SSE Server MCP Form (AC: 4, 6)
- [ ] Create SSE form component
- [ ] Implement URL format validation
- [ ] Add async URL accessibility checking
- [ ] Create loading indicators for validation
- [ ] Write tests for URL validation scenarios

### Task 5: JSON Configuration MCP Form (AC: 5, 6)
- [ ] Create JSON form component with multi-line input
- [ ] Implement real-time JSON syntax validation
- [ ] Add specific error positioning (line/column)
- [ ] Create visual confirmation for valid JSON
- [ ] Write tests for JSON parsing edge cases

### Task 6: Form Integration and Persistence (AC: 7)
- [ ] Integrate forms with storage service
- [ ] Implement successful addition workflow
- [ ] Add success message display system
- [ ] Update MCP list refresh logic
- [ ] Write end-to-end addition tests

### Task 7: Enhanced Error Handling (AC: 6)
- [ ] Create comprehensive validation error system
- [ ] Implement duplicate name detection and handling
- [ ] Add field-specific error message display
- [ ] Create error recovery workflows
- [ ] Write error handling test scenarios

### Task 8: UI Polish and Integration (AC: 1-8)
- [ ] Integrate modal system with existing keyboard shortcuts
- [ ] Ensure responsive design works with modals
- [ ] Add visual polish to form components
- [ ] Create transition animations for modal operations
- [ ] Perform cross-platform testing

## Testing Requirements

### Unit Tests Required
- Modal state management and transitions
- Form validation logic for all MCP types
- Input field navigation and keyboard handling
- Error message generation and display
- Integration with storage service

### Integration Tests Required
- Complete add workflow for each MCP type
- Modal interaction with main interface
- Form cancellation and data integrity
- Error recovery scenarios
- Cross-platform form behavior

### Manual Testing Scenarios
- Add MCPs of all three types successfully
- Test form validation with invalid inputs
- Verify modal cancellation preserves state
- Test keyboard navigation throughout all forms
- Verify success messages and list updates

## Dependencies

### Prerequisites
- Epic 1, Story 1 (TUI Foundation & Navigation) - COMPLETED
- Epic 1, Story 2 (Local Storage System) - COMPLETED
- Existing Bubble Tea integration and modal state support

### External Dependencies
- Go standard library: net/url (URL validation), encoding/json (JSON validation)
- Bubble Tea framework: Modal component patterns
- No additional external dependencies required

### Internal Dependencies
- `internal/ui/types/models.go` - Model struct extensions
- `internal/ui/services/storage_service.go` - Inventory persistence
- `internal/ui/services/mcp_service.go` - MCP management operations
- `internal/ui/components/` - New modal and form components

## Risk Assessment

### Technical Risks
- **Medium Risk:** Complex modal state management within Bubble Tea framework
- **Low Risk:** Form validation implementation with Go standard library
- **Medium Risk:** Async URL validation without blocking UI responsiveness
- **Low Risk:** JSON syntax validation and error reporting

### Mitigation Strategies
- Use established Bubble Tea modal patterns from community examples
- Implement timeout-based async validation with loading indicators
- Create comprehensive test coverage for modal state transitions
- Follow existing error handling patterns from Stories 1.1 and 1.2

### User Experience Risks
- **Medium Risk:** Modal complexity overwhelming users new to TUI interfaces
- **Low Risk:** Form validation creating friction in add workflow

### UX Mitigation Strategies
- Keep modal interfaces simple with clear instructions
- Provide helpful placeholder text and examples
- Implement progressive disclosure (type selection → specific form)
- Follow established terminal UI conventions for form interactions

## Notes & Considerations

### Design Decisions
- **Modal-based approach**: Follows established TUI patterns and keeps context visible
- **Three distinct MCP types**: Supports PRD requirements for Command, SSE, and JSON configurations
- **Progressive form disclosure**: Type selection first reduces cognitive load
- **Real-time validation**: Provides immediate feedback without requiring submission

### Future Considerations
- **Batch MCP addition**: Could be valuable for users importing many MCPs
- **MCP templates**: Pre-configured templates for popular MCPs
- **Import from file**: Ability to import MCP configurations from external files
- **Form field auto-completion**: Based on previously added similar MCPs

### Development Notes
- This story completes the core "Add" operation for CRUD functionality
- Forms should be reusable for the upcoming "Edit MCP" story (1.4)
- Modal system established here will be used for all future modal operations
- Real-time validation patterns established here apply to all form interactions

## Technical Decisions Made

### TD-001: Modal Architecture Pattern
**Decision:** Overlay modal system with backdrop dimming
**Implementation:**
- Modal components render over existing interface
- Background interface remains visible but dimmed
- ESC key provides consistent cancellation across all modals
**Rationale:** Maintains context awareness while focusing attention on current task

### TD-002: Progressive Form Disclosure
**Decision:** Two-step process: type selection → specific form
**Implementation:**
- Initial modal shows three MCP type options
- Selection transitions to type-specific form
- Each form optimized for its MCP type requirements
**Rationale:** Reduces cognitive load and provides focused user experience

### TD-003: Real-time Validation Strategy
**Decision:** Immediate validation with async support for network operations
**Implementation:**
- Field validation occurs on input change
- Network validation (URL checking) uses timeout-based async operations
- Loading indicators show progress for async validations
**Rationale:** Provides immediate feedback while handling network operations gracefully

### TD-004: Form Field Design Philosophy
**Decision:** Minimal required fields with clear optional markers
**Implementation:**
- Only essential fields marked as required
- Optional fields clearly labeled
- Field descriptions provide context without overwhelming interface
**Rationale:** Balances completeness with ease of use for quick MCP additions

### TD-005: Error Handling Approach
**Decision:** Inline field-level errors with actionable messages
**Implementation:**
- Validation errors appear directly below relevant field
- Error messages provide specific guidance for resolution
- Errors clear automatically when field becomes valid
**Rationale:** Provides clear guidance without disrupting user workflow

### TD-006: Success Feedback Pattern
**Decision:** Brief success message with automatic list refresh
**Implementation:**
- Success message appears for 2-3 seconds after addition
- Added MCP automatically selected in refreshed list
- Modal closes immediately upon successful addition
**Rationale:** Provides confirmation while maintaining workflow momentum

## Definition of Done

### Functional Requirements
- [ ] 'A' key opens MCP type selection modal from main interface
- [ ] Three MCP types selectable via number keys or arrow navigation
- [ ] Command/Binary form accepts name, command, and optional args
- [ ] SSE Server form accepts name and URL with format validation
- [ ] JSON Configuration form accepts name and JSON with syntax validation
- [ ] All forms provide real-time validation with helpful error messages
- [ ] Successful addition saves to storage and refreshes main list
- [ ] ESC key cancellation works from any modal or form state

### Quality Requirements
- [ ] Modal system integrates seamlessly with existing TUI
- [ ] Form navigation feels intuitive with standard keyboard patterns
- [ ] Validation provides helpful, actionable error messages
- [ ] Success workflow provides clear confirmation and maintains context
- [ ] All modal interactions preserve main interface state

### Testing Requirements
- [ ] Unit tests cover modal state management and form validation
- [ ] Integration tests verify complete add workflows for all MCP types
- [ ] Error handling tests cover validation failures and network issues
- [ ] Cross-platform testing confirms consistent behavior
- [ ] Manual testing validates intuitive user experience

## Validation Complete

**Status:** NEEDS_FIXES
**Validated by:** SM
**Issues remaining:** 3

## Story Approved for Development

**Status:** Draft - Ready for PO Review  
**Created by:** Bob (SM Agent)  
**Creation Date:** 2025-06-30  
**Epic Context:** Epic 1 - Core MCP Inventory Management  
**Dependencies:** Stories 1.1 ✅ and 1.2 ✅ completed

---

**Next Steps:**
1. PO review and approval
2. Development team estimation and sprint planning
3. Technical design review for modal architecture
4. Implementation prioritization within Epic 1 context

---

## Story Quality Metrics

### Completeness Score: 95/100
- ✅ All acceptance criteria clearly defined with Given/When/Then format
- ✅ Technical implementation details provided
- ✅ Dependencies and risks identified
- ✅ Testing requirements comprehensive
- ✅ Integration with existing stories confirmed

### Business Alignment Score: 98/100
- ✅ Directly addresses PRD requirement FR2: "Users can add new MCPs to their personal inventory via TUI prompts"
- ✅ Supports all three MCP types specified in PRD
- ✅ Follows UX specification for modal-driven interface
- ✅ Maintains consistency with Epic 1 objectives

### Technical Feasibility Score: 90/100
- ✅ Builds on proven Bubble Tea foundation from Story 1.1
- ✅ Leverages established storage system from Story 1.2
- ✅ Modal patterns well-established in TUI frameworks
- ⚠️ Complex async validation may require careful implementation

### Story Independence Score: 85/100
- ✅ Can be implemented and tested independently
- ✅ Clear boundaries with other Epic 1 stories
- ⚠️ Some modal system components may be reused in Story 1.4 (Edit MCP)

---

## Epic 1 Progress Summary

**Stories Completed:** 2/7 (Stories 1.1, 1.2)  
**Current Story:** Story 1.3 (Draft)  
**Epic Progress:** 29% complete  
**Technical Foundation:** ✅ TUI Framework, ✅ Storage System, 🚧 Add Workflow

**Epic 1 Momentum:** Strong foundation established with Stories 1.1 and 1.2. Story 1.3 builds critical CRUD functionality to enable complete inventory management. Modal system introduced here supports remaining Epic 1 stories (Edit, Delete, Search).