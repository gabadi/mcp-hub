# MCP Manager CLI

> Terminal-based MCP (Model Context Protocol) inventory management tool for Claude Code developers

## Project Status

**Current:** Epic 1, Story 1 ✅ Complete  
**Foundation:** TUI with responsive navigation, search, and keyboard shortcuts  
**Next:** Story 1.2 - Local Storage System  

## Quick Start

```bash
# Build and run
task build && task run

# Development mode (auto-rebuild)
task dev

# Run tests
task test
```

## Architecture

### Framework & Dependencies
- **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea) v1.3.5
- **Styling:** [Lipgloss](https://github.com/charmbracelet/lipgloss) v1.1.0
- **Language:** Go 1.23.5
- **Build System:** [Task](https://taskfile.dev/) (Taskfile.yml)

### Project Structure
```
├── cmd/mcp-manager/          # Application entry point
├── internal/
│   ├── tui/                  # TUI components and logic
│   ├── models/               # Data models (MCP servers, config)
│   ├── config/               # Configuration management
│   └── errors/               # Error handling system
├── pkg/mcp/                  # MCP-specific functionality
├── docs/                     # Requirements and specifications
│   ├── prd.md               # Product Requirements
│   ├── front-end-spec.md    # UI/UX Specification
│   └── stories/             # BMAD story files
└── Taskfile.yml            # Development workflow
```

### Current Implementation

#### TUI Foundation (Epic 1, Story 1)
- **Responsive Layout:** 3-column (≥80 cols), 2-column (60-79), 1-column (<60)
- **Navigation:** Arrow keys + Vim keys (hjkl) + Tab (search) + ESC/Q (quit)
- **Search Mode:** Tab to activate, real-time query input, ESC to exit
- **Performance:** <100ms startup, <50ms navigation response
- **Testing:** 95.5% coverage, comprehensive unit/integration/acceptance tests

#### Architecture Patterns
- **MVC:** Clean separation of Model, View, Update logic
- **State Management:** Centralized immutable state with Bubble Tea patterns
- **Component System:** Modular TUI components (ready for Epic 1 expansion)
- **Error Handling:** Robust error management with user-friendly messages

## Development

### Prerequisites
- Go 1.21+
- [Task](https://taskfile.dev/installation/) (replaces make)
- Terminal with 256 colors (most modern terminals)

### Common Commands
```bash
# Core Development
task build              # Build binary
task run                # Build and run
task dev                # Auto-rebuild on changes
task test               # Run all tests
task test-coverage      # Tests with coverage report

# Code Quality
task check              # Run all quality checks
task fmt                # Format code
task lint               # Lint code (requires golangci-lint)
task vet                # Run go vet

# Testing Variants
task test-race          # Tests with race detection
task test-bench         # Benchmark tests
task test-unit          # Unit tests only
task test-integration   # Integration tests only

# Dependencies
task deps               # Install/tidy dependencies
task deps-update        # Update all dependencies

# Terminal Compatibility
task term-test          # Test terminal compatibility
```

### Project Workflow (BMAD Method)
This project follows Business-Motivated Agile Development (BMAD) with story-driven development:

1. **Epic → Stories:** Features broken into implementable stories
2. **Story Implementation:** TUI foundation → Storage → Features
3. **Quality Gates:** Architecture + Business + Process + QA + UX reviews
4. **Learning Extraction:** Technical debt and improvements tracked systematically

## Epic 1: Core MCP Inventory Management

### Completed Stories
- ✅ **Story 1.1:** TUI Foundation & Navigation (Done - Delivered)

### Planned Stories  
- 🔄 **Story 1.2:** Local Storage System (JSON-based MCP inventory)
- 📋 **Story 1.3:** Add MCP Workflow (TUI prompts for MCP types)
- 📋 **Story 1.4:** Edit MCP Capability (Modify existing MCPs)
- 📋 **Story 1.5:** Remove MCP Operation (Delete with confirmation)
- 📋 **Story 1.6:** Search & Filter (Real filtering, not just display)
- 📋 **Story 1.7:** Seed Data Integration (Bootstrap with common MCPs)

## Technical Decisions & Learning Items

### High Priority (Next Sprint)
1. **Code Organization:** Split `internal/tui/model.go` (353 lines) into focused modules
2. **Component Abstractions:** Extract HeaderComponent, ColumnComponent, FooterComponent

### Medium Priority (Epic 1.2-1.6)
3. **Configuration System:** Extract hard-coded layout breakpoints and dimensions
4. **Search Enhancement:** Implement real filtering logic (currently UI-only)

### Low Priority (Epic 2+)
5. **Interface Abstractions:** Plugin system foundation for extensibility

## Documentation

### For LLMs/AI Development
- **PRD:** `/docs/prd.md` - Product requirements and business context
- **UI Spec:** `/docs/front-end-spec.md` - TUI patterns and user flows  
- **Story Files:** `/docs/stories/` - BMAD story documentation
- **Architecture:** Clean MVC with Bubble Tea, comprehensive test coverage
- **Tech Stack:** Go + Bubble Tea + Task, terminal-native design

### Key Implementation Notes
- **Bubble Tea Framework:** MVC architecture with proper state management
- **Responsive Design:** Adapts to terminal width (3/2/1 column layouts)
- **Testing Strategy:** Unit + Integration + Acceptance tests (95.5% coverage)
- **Performance:** Exceeds requirements (startup <100ms vs 500ms target)
- **Error Handling:** Comprehensive error management with recovery patterns

## Contributing

1. **Review Requirements:** Read `/docs/prd.md` and `/docs/front-end-spec.md`
2. **Follow BMAD Process:** Implement features as documented stories
3. **Maintain Quality:** All tests must pass (`task test`)
4. **Code Standards:** Run quality checks (`task check`)
5. **Architecture:** Follow established MVC patterns with Bubble Tea

## License

MIT License - see LICENSE file for details.

---

**Built with:** Go + Bubble Tea | **Development:** Task + BMAD Method | **Status:** Epic 1.1 Complete ✅