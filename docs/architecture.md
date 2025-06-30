# CC MCP Manager - System Architecture

**Document Version:** 1.0  
**Created:** 2025-06-30  
**Last Updated:** 2025-06-30  
**Owner:** Technical Lead  
**Status:** Active

## Executive Summary

The CC MCP Manager is a Terminal User Interface (TUI) application built in Go that enables developers to manage their Model Context Protocol (MCP) inventory. The system follows a clean architecture pattern with clear separation of concerns, leveraging the Bubble Tea framework for reactive UI management and implementing a responsive, multi-column layout system.

### Key Architectural Decisions

1. **TUI-First Design**: Native terminal interface for developer-friendly experience
2. **Reactive State Management**: Event-driven UI using Bubble Tea framework
3. **Local-First Storage**: JSON-based configuration with atomic file operations
4. **Responsive Layout System**: Adaptive column layouts based on terminal width
5. **Modular Component Architecture**: Separation of UI, business logic, and storage layers

## System Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      CC MCP Manager                           │
├─────────────────────────────────────────────────────────────────┤
│                    main.go (Entry Point)                      │
├─────────────────────────────────────────────────────────────────┤
│                     UI Layer (Bubble Tea)                     │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────┐ │
│  │   Model     │ │    View     │ │   Update    │ │ Components│ │
│  │  (State)    │ │ (Rendering) │ │ (Logic)     │ │ (Reusable)│ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └───────────┘ │
├─────────────────────────────────────────────────────────────────┤
│                    Handler Layer                               │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────┐ │
│  │ Navigation  │ │  Keyboard   │ │   Search    │ │   Modal   │ │
│  │  Handler    │ │   Handler   │ │  Handler    │ │  Handler  │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └───────────┘ │
├─────────────────────────────────────────────────────────────────┤
│                   Service Layer                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │
│  │    MCP      │ │   Storage   │ │   Layout    │               │
│  │  Service    │ │  Service    │ │  Service    │               │
│  └─────────────┘ └─────────────┘ └─────────────┘               │
├─────────────────────────────────────────────────────────────────┤
│                    Types Layer                                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │
│  │    Model    │ │    State    │ │ Constants   │               │
│  │  Definitions│ │ Definitions │ │    & Enums  │               │
│  └─────────────┘ └─────────────┘ └─────────────┘               │
├─────────────────────────────────────────────────────────────────┤
│                  Storage Layer                                 │
│              Local JSON Configuration                          │
│         ~/.config/cc-mcp-manager/inventory.json               │
└─────────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. Entry Point (`main.go`)

**Responsibility**: Application initialization and Bubble Tea program execution

**Key Functions**:
- Initialize the UI model with default or loaded MCP data
- Start the Bubble Tea program loop
- Handle graceful shutdown

**Dependencies**:
- `internal/ui` package for model initialization
- Bubble Tea framework for program execution

### 2. UI Layer (`internal/ui/`)

#### 2.1 Model (`model.go`)

**Responsibility**: Central state management and data model

**Core State**:
```go
type Model struct {
    // Layout state
    Width, Height   int
    ColumnCount     int
    ActiveColumn    int
    Columns         []Column
    
    // Data state
    MCPItems        []MCPItem
    SelectedItem    int
    
    // Interaction state
    State           AppState
    SearchQuery     string
    SearchActive    bool
    SearchInputActive bool
}
```

**State Management**:
- Immutable state updates following Bubble Tea patterns
- Centralized state for all UI components
- Type-safe state transitions

#### 2.2 View (`view.go`)

**Responsibility**: UI rendering and layout composition

**Layout System**:
- Responsive design with breakpoints (120+, 80-119, <80 width)
- Dynamic column allocation based on terminal dimensions
- Component composition with consistent styling

**Rendering Pipeline**:
1. Calculate available space
2. Determine layout configuration
3. Render individual components
4. Compose final view with proper spacing

#### 2.3 Update (`update.go`)

**Responsibility**: Event handling and state transitions

**Message Processing**:
- Window resize events
- Keyboard input events
- Application state transitions
- Command generation for side effects

**Event Flow**:
```
User Input → Keyboard Handler → State Update → View Render
                ↓
         Side Effects (Commands)
                ↓
         Service Layer Operations
```

### 3. Component Layer (`internal/ui/components/`)

#### 3.1 Grid Component (`grid.go`)

**Responsibility**: MCP list rendering and selection management

**Features**:
- Multi-column MCP display
- Active/inactive status visualization
- Responsive column allocation
- Selection highlighting

#### 3.2 Header Component (`header.go`)

**Responsibility**: Application title and status information

**Information Display**:
- Application branding
- Current mode indication
- Active MCP count
- Layout adaptation

#### 3.3 Footer Component (`footer.go`)

**Responsibility**: Help text and keyboard shortcut display

**Dynamic Content**:
- Context-sensitive help
- Available actions based on current state
- Search mode indicators

#### 3.4 Modal Component (`modal.go`)

**Responsibility**: Overlay dialogs for user interactions

**Modal Types**:
- Add MCP modal (planned)
- Edit MCP modal (planned)
- Confirmation dialogs (planned)

### 4. Handler Layer (`internal/ui/handlers/`)

#### 4.1 Navigation Handler (`navigation.go`)

**Responsibility**: Column and item navigation logic

**Navigation Features**:
- Left/Right arrow column navigation
- Up/Down arrow item selection
- Wraparound behavior
- Layout-aware constraints

#### 4.2 Keyboard Handler (`keyboard.go`)

**Responsibility**: Global keyboard input processing

**Key Bindings**:
- Navigation keys (arrows, hjkl)
- Action keys (space, enter)
- Mode switches (/, escape)
- Application controls (q, ctrl+c)

#### 4.3 Search Handler (`search.go`)

**Responsibility**: Search functionality and filtering

**Search Features**:
- Real-time filtering
- Case-insensitive matching
- Search mode state management
- Query input handling

### 5. Service Layer (`internal/ui/services/`)

#### 5.1 MCP Service (`mcp_service.go`)

**Responsibility**: MCP business logic and operations

**Core Operations**:
- MCP filtering and search
- Status toggling (active/inactive)
- MCP data validation
- Collection management

**Data Operations**:
```go
// Core service methods
GetFilteredMCPs(model Model) []MCPItem
ToggleMCPStatus(model *Model, index int)
GetSelectedMCP(model Model) *MCPItem
GetActiveMCPCount(model Model) int
```

#### 5.2 Storage Service (`storage_service.go`)

**Responsibility**: Data persistence and configuration management

**Storage Features**:
- Atomic file operations
- JSON serialization with metadata
- Configuration directory management
- Error handling and recovery

**File Structure**:
```json
{
  "version": "1.0",
  "timestamp": "2025-06-30T12:00:00Z",
  "inventory": [
    {
      "name": "github-mcp",
      "type": "CMD",
      "active": true,
      "command": "github-mcp"
    }
  ]
}
```

#### 5.3 Layout Service (`layout_service.go`)

**Responsibility**: Responsive layout calculations

**Layout Logic**:
- Breakpoint determination
- Column count calculation
- Width distribution
- Component sizing

### 6. Types Layer (`internal/ui/types/`)

#### 6.1 Model Definitions (`models.go`)

**Core Data Types**:
```go
type MCPItem struct {
    Name        string
    Type        string  // CMD, SSE, JSON, HTTP
    Active      bool
    Command     string
    Args        string
    URL         string
    JSONConfig  string
}

type AppState int
const (
    MainNavigation AppState = iota
    SearchMode
    SearchActiveNavigation
    ModalActive
)
```

#### 6.2 Constants (`constants.go`)

**Layout Constants**:
- Breakpoint values (WIDE_LAYOUT_MIN, MEDIUM_LAYOUT_MIN)
- Column counts (WIDE_COLUMNS, MEDIUM_COLUMNS, NARROW_COLUMNS)
- UI dimensions and spacing

## Data Flow Architecture

### 1. Application Initialization

```
main() → NewModel() → LoadInventory() → StartProgram()
  ↓
Default MCPs + Loaded MCPs → Initial State
  ↓
First Render
```

### 2. User Interaction Flow

```
User Input → Keyboard Handler → Update Logic
     ↓              ↓               ↓
  Key Event → State Validation → New State
                    ↓               ↓
               Side Effects → Command Execution
                              ↓
                         Service Layer
                              ↓
                         Storage/Operations
```

### 3. State Update Cycle

```
Previous State + Event → New State + Commands
        ↓                        ↓
   View Rendering           Command Execution
        ↓                        ↓
   Terminal Display         Service Operations
                                 ↓
                            Future Events
```

## Storage Architecture

### 1. Configuration Management

**Location**: `~/.config/cc-mcp-manager/inventory.json`

**Benefits**:
- Standard configuration directory
- Cross-platform compatibility
- User-specific settings
- Version control friendly

### 2. Data Persistence

**Atomic Operations**:
1. Write to temporary file
2. Verify write success
3. Atomic rename to target
4. Clean up temporary files

**Error Handling**:
- Corrupted file backup
- Default data fallback
- Graceful degradation
- User notification

### 3. Data Format

**Version Management**:
- Explicit version field
- Migration support
- Backward compatibility
- Schema validation

## Component Communication

### 1. Service Integration

```go
// Handler uses service
func HandleKeyPress(model *Model, key string) {
    switch key {
    case " ":
        services.ToggleMCPStatus(model, model.SelectedItem)
        return SaveInventoryCommand(model)
    }
}

// Service operates on model
func ToggleMCPStatus(model *Model, index int) {
    if index >= 0 && index < len(model.MCPItems) {
        model.MCPItems[index].Active = !model.MCPItems[index].Active
    }
}
```

### 2. Component Composition

```go
// View composes components
func (m Model) View() string {
    header := components.RenderHeader(m)
    grid := components.RenderGrid(m)
    footer := components.RenderFooter(m)
    
    return lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        grid,
        footer,
    )
}
```

## Responsive Design System

### 1. Breakpoint Strategy

**Wide Layout (120+ chars)**:
- 4-column MCP grid
- Maximum information density
- Full feature display

**Medium Layout (80-119 chars)**:
- 3-column layout
- Balanced information display
- Optimized for standard terminals

**Narrow Layout (<80 chars)**:
- 2-column layout
- Essential information only
- Mobile-friendly design

### 2. Adaptive Components

**Grid Adaptation**:
- Dynamic column count
- Proportional width allocation
- Content prioritization

**Component Scaling**:
- Text truncation strategies
- Priority-based information display
- Consistent spacing ratios

## Security Considerations

### 1. Configuration Security

**File Permissions**:
- Configuration directory: 0700 (user only)
- Configuration file: 0600 (user read/write only)
- Temporary files: 0600

**Data Validation**:
- JSON schema validation
- Command injection prevention
- Path traversal protection

### 2. Command Execution

**Future Considerations**:
- MCP command validation
- Sandboxed execution
- User permission prompts

## Performance Architecture

### 1. Rendering Optimization

**Efficient Updates**:
- Minimal re-rendering
- Component-level change detection
- Lazy evaluation where possible

**Memory Management**:
- Bounded slice operations
- String builder usage
- Garbage collection friendly

### 2. Data Operations

**Search Performance**:
- Case-insensitive string matching
- Early termination
- Bounded result sets

**File I/O**:
- Atomic operations
- Minimal disk access
- Error recovery

## Testing Architecture

### 1. Testing Strategy

**Test Categories**:
- Unit tests: Individual component testing
- Integration tests: Cross-component workflows
- End-to-end tests: Complete user scenarios
- Performance tests: Benchmark critical paths

### 2. Test Infrastructure

**Test Utilities** (`internal/testutil/`):
- Model builders for consistent test data
- Helper functions for common operations
- Mock implementations where needed

**Coverage Requirements**:
- Minimum 85% overall coverage
- 90% coverage for new features
- 100% coverage for error paths

## Future Architecture Considerations

### 1. Planned Enhancements

**Modal System**:
- Add/Edit MCP workflows
- Configuration wizards
- Confirmation dialogs

**MCP Operations**:
- Start/stop MCP servers
- Health checking
- Performance monitoring

### 2. Scalability Considerations

**Large Inventories**:
- Pagination support
- Virtual scrolling
- Incremental search

**Cross-Platform Support**:
- Windows compatibility
- Terminal capability detection
- Fallback rendering modes

## Development Guidelines

### 1. Code Organization

**Package Structure**:
- Clear separation of concerns
- Minimal inter-package dependencies
- Consistent naming conventions

**Interface Design**:
- Small, focused interfaces
- Dependency injection friendly
- Testability first

### 2. Contribution Guidelines

**Code Quality**:
- Go fmt compliance
- Comprehensive testing
- Documentation requirements
- Performance considerations

**Architecture Compliance**:
- Follow established patterns
- Maintain separation of concerns
- Consider backward compatibility
- Update documentation

## Deployment Architecture

### 1. Build System

**Single Binary**:
- Static compilation
- No external dependencies
- Cross-platform builds

**Distribution**:
- CLI installation
- Package manager support
- Automated releases

### 2. Configuration Management

**Default Behavior**:
- Graceful degradation
- Sensible defaults
- User customization support

**Migration Support**:
- Version detection
- Automatic upgrades
- Backward compatibility

---

**Document Maintenance**:
This architecture document is a living document that should be updated as the system evolves. Major architectural changes should be reviewed and approved by the technical lead before implementation.

**Related Documentation**:
- [Testing Standards](./testing-standards.md)
- [Story Documentation](./stories/)
- [Front-End Specifications](./front-end-spec.md)
- [Product Requirements](./prd.md)