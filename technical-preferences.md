# Technical Preferences

## Development MCPs (Seed Data)

### Purpose
These MCPs are recommended for development and testing of the MCP Manager CLI. They can be used as optional seed data during development and provide specific functionality for validation.

### Recommended MCPs

#### Context7
- **Name:** context7
- **Type:** SSE
- **URL:** https://mcp.context7.com/sse
- **Purpose:** Library version lookup and best practices research
- **Usage:** Development-time validation of latest library versions and implementation patterns
- **When to use:** When researching latest versions of dependencies, implementation examples, or best practices

#### HT-MCP
- **Name:** ht-mcp  
- **Type:** Command
- **Command:** ht-mcp
- **Purpose:** Manual testing and validation of MCP functionality
- **Usage:** QA and development testing to verify MCP operations work correctly
- **When to use:** For manual testing of the MCP Manager CLI functionality, validating enable/disable operations

### Implementation Notes

- These MCPs should be defined in a seed data configuration file for optional initialization
- Users should be able to start with a clean inventory if preferred
- Seed data should be clearly marked as optional/development-focused
- Production users may choose to skip seed data entirely

### Technology Stack

#### Core Technologies
- **Language:** Go (latest stable version)
- **TUI Framework:** Bubble Tea (latest stable version)
- **Configuration:** JSON files for local storage
- **CLI Integration:** Shell command execution for Claude Code integration

#### Project Structure
- Standard Go project layout
- Separation of UI (Bubble Tea) and business logic
- Clean interfaces for Claude Code integration
- Simple JSON-based local storage

#### Testing Approach
- Unit tests for core business logic
- Integration tests for Claude Code CLI interaction  
- Manual testing using ht-mcp for TUI validation
- No automated UI testing required for MVP