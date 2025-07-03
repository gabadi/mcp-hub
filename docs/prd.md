# MCP Manager CLI Product Requirements Document (PRD)

## Goals and Background Context

### Goals

- Create a personal TUI tool to easily manage local MCP inventory without context overload
- Enable quick switching of MCPs for focused Claude Code sessions
- Eliminate the need to remember complex MCP installation/removal commands
- Provide a clean, intuitive Bubble Tea interface for MCP management

### Background Context

Developers using Claude Code frequently face MCP management friction. The current `claude mcp` CLI commands are functional but require remembering syntax and don't provide an easy way to manage a personal "toolkit" of MCPs. More importantly, having too many MCPs active simultaneously can overwhelm the LLM with unnecessary context, reducing effectiveness.

This tool addresses the gap between MCP discovery (done elsewhere) and daily MCP usage by providing a personal inventory system with visual switching capabilities. The focus is on post-discovery workflow optimization rather than marketplace functionality.

### Change Log

| Date | Version | Description | Author |
| :--- | :------ | :---------- | :----- |
| 2025-06-29 | 1.0 | Initial PRD creation | John (PM) |

## Requirements

### Functional

- FR1: The TUI shall display a list of user's personal MCP inventory
- FR2: Users can add new MCPs to their personal inventory via TUI prompts
- FR3: Users can remove MCPs from their personal inventory
- FR4: Users can enable/disable MCPs for Claude Code directly from the TUI
- FR5: The TUI shall show visual indicators for which MCPs are currently active in Claude Code
- FR6: The system shall maintain state consistency between local inventory and Claude Code configuration
- FR7: All MCP operations shall provide immediate visual feedback in the TUI

### Non Functional

- NFR1: The TUI must be responsive and feel snappy for daily developer use
- NFR2: Local inventory storage should use simple JSON format for transparency and debuggability
- NFR3: Error handling must gracefully recover from Claude Code CLI failures
- NFR4: The interface should follow standard TUI navigation patterns (keyboard-driven)

## User Interface Design Goals

### Overall UX Vision

Clean, minimal TUI focused on speed and clarity. The interface should feel like a developer tool - functional over flashy, with clear visual hierarchy and instant feedback. Navigation should be intuitive for anyone familiar with terminal interfaces.

### Key Interaction Paradigms

- **Keyboard-first navigation**: Arrow keys, Enter, Escape standard patterns
- **List-based interface**: Primary view is a navigable list of MCPs
- **Modal prompts**: Add/remove operations use simple modal dialogs
- **Status indicators**: Clear visual differentiation between active/inactive MCPs
- **Instant feedback**: Operations show immediate results with success/error states

### Core Screens and Views

- **Main Inventory Screen**: List of personal MCPs with active status indicators
- **Add MCP Modal**: Simple form to add new MCP (name + command/URL)
- **Confirmation Dialogs**: For destructive operations (remove, disable)
- **Status/Feedback Area**: Show operation results and current state

### Accessibility: None

MVP focuses on core functionality. Future versions may add accessibility features.

### Branding

Minimal, professional developer tool aesthetic. No special branding requirements - clean terminal interface with standard TUI conventions.

### Target Device and Platforms

Terminal/CLI environment on macOS, Linux, and Windows. Single binary distribution.

## Technical Assumptions

### Repository Structure: Single Repository

Simple Go project structure with standard organization.

### Service Architecture

Single binary CLI application with TUI interface. No server components or external services required.

Architecture: TUI (Bubble Tea) → Local Storage (JSON) → Claude Code CLI integration

### Testing requirements

- Unit tests for core logic (inventory management, Claude Code integration)
- Integration tests for Claude Code CLI interaction
- Manual testing for TUI functionality
- No automated UI testing required for MVP

### Additional Technical Assumptions and Requests

- **Language**: Go for cross-platform compatibility and single binary distribution
- **TUI Framework**: Bubble Tea for rich terminal interface
- **Local Storage**: JSON file in user's home directory or XDG config location
- **Claude Code Integration**: Shell command execution of `claude mcp` commands
- **Error Handling**: Graceful degradation when Claude Code CLI is unavailable
- **Configuration**: Minimal config needs, defaults should work out of box

## Epics

### Epic List

1. **Core MCP Inventory Management**: Establish foundational TUI and complete CRUD operations for personal MCP inventory
2. **Claude Code Integration**: Implement bidirectional sync between local inventory and Claude Code for project-specific activation

## Epic 1: Core MCP Inventory Management

Establish the foundational TUI interface with Bubble Tea and implement complete MCP inventory management with full CRUD operations. This epic delivers a working TUI that can manage a local collection of MCPs with multiple types and search capabilities.

### Story 1.1: TUI Foundation & Navigation

As a developer,
I want a responsive TUI interface with intuitive navigation,
so that I can efficiently manage my MCP inventory.

#### Acceptance Criteria

- 1: Application launches with Bubble Tea TUI showing 3-column layout
- 2: Arrow keys navigate within and between columns
- 3: Tab key jumps to search field
- 4: ESC key exits application or clears current mode
- 5: Responsive layout adapts to terminal width (3/2/1 columns)
- 6: Header shows keyboard shortcuts and current context

### Story 1.2: Local Storage System

As a developer,
I want my MCP inventory to persist between sessions,
so that I don't lose my configuration.

#### Acceptance Criteria

- 1: JSON file created in appropriate config directory
- 2: Inventory loads automatically on startup
- 3: File format supports multiple MCP types (CMD/SSE/JSON/HTTP)
- 4: Graceful handling of missing or corrupted config files
- 5: Config file location logged for user reference

### Story 1.3: Add MCP Workflow

As a developer,
I want to add different types of MCPs to my inventory,
so that I can support various MCP configurations.

#### Acceptance Criteria

- 1: 'A' key opens MCP type selection modal
- 2: Modal shows options: Command/Binary, SSE Server, JSON Configuration
- 3: Each type presents appropriate form fields
- 4: Form validation ensures required fields are provided
- 5: ESC cancels operation, Enter confirms
- 6: List refreshes with new MCP and type indicator
- 7: Success message confirms addition

### Story 1.4: Edit MCP Capability

As a developer,
I want to modify existing MCP details,
so that I can update configurations without re-adding.

#### Acceptance Criteria

- 1: 'E' key opens edit modal for selected MCP
- 2: Modal pre-populates with current MCP details
- 3: Same validation as add workflow
- 4: Changes save to local inventory
- 5: List refreshes with updated information

### Story 1.5: Remove MCP Operation

As a developer,
I want to remove MCPs I no longer need,
so that my inventory stays clean and relevant.

#### Acceptance Criteria

- 1: 'D' key opens confirmation dialog for selected MCP
- 2: Confirmation shows MCP name and type being deleted
- 3: ESC cancels, Enter confirms deletion
- 4: MCP removed from local inventory
- 5: List refreshes with MCP removed

### Story 1.6: Search & Filter

As a developer,
I want to quickly find specific MCPs in my inventory,
so that I can locate tools efficiently.

#### Acceptance Criteria

- 1: '/' key or typing in search field activates search mode
- 2: Real-time filtering as user types
- 3: Matching text highlighted in results
- 4: ESC clears search and shows full list
- 5: Search works across MCP names
- 6: Status shows number of matches found

### Story 1.7: Seed Data Integration

As a developer,
I want option to initialize with recommended development MCPs,
so that I can start with useful tools pre-configured.

#### Acceptance Criteria

- 1: Application detects technical-preferences.md configuration
- 2: Option to initialize with context7 and ht-mcp on first run
- 3: User can choose to start with empty inventory
- 4: Seed MCPs clearly marked as development/testing tools
- 5: Normal add/edit/remove operations work on seed MCPs

## Epic 2: Claude Code Integration

Implement bidirectional synchronization between the local MCP inventory and Claude Code's MCP configuration, enabling project-specific MCP activation and status visualization.

### Story 2.1: Claude Status Detection

As a developer,
I want the application to detect Claude Code CLI availability and current state,
so that I know what MCPs are currently active.

#### Acceptance Criteria

- 1: Application detects if `claude` CLI is available
- 2: Startup queries `claude mcp list` to get current active MCPs
- 3: Graceful handling when Claude CLI is not available
- 4: Error messages provide helpful installation guidance
- 5: Manual refresh option (R key) to re-query status

### Story 2.2: MCP Activation Toggle

As a developer,
I want to enable and disable MCPs for the current project,
so that I can focus my Claude Code context.

#### Acceptance Criteria

- 1: SPACE key toggles enable/disable for selected MCP
- 2: Enable calls appropriate `claude mcp add` command based on MCP type
- 3: Disable calls `claude mcp remove <name>`
- 4: Loading indicators during async operations
- 5: Success/failure feedback for each operation
- 6: Operations only affect current project directory

### Story 2.3: Status Visualization

As a developer,
I want clear visual indicators of which MCPs are active,
so that I understand my current Claude Code configuration.

#### Acceptance Criteria

- 1: Active MCPs show ● (green dot) indicator
- 2: Inactive MCPs show ○ (gray circle) indicator
- 3: Status updates in real-time after toggle operations
- 4: Visual distinction between available and active MCPs
- 5: Type badges ([CMD], [SSE], [JSON]) always visible

### Story 2.4: Error Handling & Recovery

As a developer,
I want graceful error handling when Claude operations fail,
so that I can understand and resolve issues.

#### Acceptance Criteria

- 1: Network timeouts handled with retry options
- 2: Permission errors show specific resolution guidance
- 3: Invalid MCP configurations detected and reported
- 4: Claude CLI version compatibility checks
- 5: Error states don't crash the application

### Story 2.5: Project Context Display

As a developer,
I want to see my current project context and sync status,
so that I know which project's MCPs I'm managing.

#### Acceptance Criteria

- 1: Status bar shows current project directory path
- 2: Last sync timestamp displayed
- 3: Active MCP count shown (e.g., "3/30 active")
- 4: Project context updates when changing directories
- 5: Clear indication when out of sync with Claude Code

### Story 2.6: Enhanced Loading State Feedback System

As a developer,
I want clear feedback during all loading and background operations,
so that I understand system status and can respond appropriately to delays or issues.

#### Acceptance Criteria

- 1: Application startup shows progressive loading messages with current operation status
- 2: MCP activation/deactivation operations display loading spinners with operation descriptions
- 3: Long-running operations provide periodic progress updates and estimated time remaining
- 4: All loading states can be cancelled with ESC key, returning to previous stable state
- 5: Background sync operations show unobtrusive progress indicators in status bar
- 6: Loading state transitions are smooth with clear visual feedback for completion or failure

## Checklist Results Report

[To be populated after checklist execution]

## Next Steps

### Design Architect Prompt

Not applicable for this MVP - the TUI design is straightforward and can be handled during development.

### Architect Prompt

Please review this PRD and create a technical architecture document for the MCP Manager CLI. Focus on:

1. Go project structure and organization
2. Bubble Tea TUI architecture and state management  
3. Local storage design and data models
4. Claude Code CLI integration patterns
5. Error handling and user feedback strategies
6. Build and distribution approach

The architecture should emphasize simplicity and maintainability for this developer-focused MVP.