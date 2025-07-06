# MCP Hub - Claude MCP Configuration Manager

> **Transform MCP Management from Frustration to Flow**

A powerful Terminal User Interface (TUI) for managing Claude Model Context Protocol (MCP) configurations. Built for developers who want to streamline their Claude Code workflow without context overload.

## 🎯 What is MCP Hub?

**The Problem**: Managing Claude MCPs is painful. You forget commands, lose track of configurations, and get overwhelmed by too many active MCPs cluttering your Claude Code context.

**The Solution**: MCP Hub provides an intuitive terminal interface that transforms MCP management from a chore into a seamless workflow. Think of it as your personal MCP inventory manager - discover, organize, and activate MCPs with surgical precision.

## ✨ Key Benefits

### 🚀 **Instant Productivity**
- Launch with `npx mcp-hub-tui` - no installation required
- Visual MCP inventory with real-time status indicators
- One-key MCP activation/deactivation (Space bar)

### 🧠 **Context Control**
- Focus your Claude Code context by enabling only relevant MCPs
- Prevent context overload that reduces LLM effectiveness
- Project-specific MCP management

### 🎨 **Developer-Friendly Experience**
- Beautiful, responsive TUI that adapts to your terminal
- Intuitive keyboard navigation (arrows, tab, space)
- Search and filter your MCP collection instantly

## 📦 Quick Start

### Installation & Usage

```bash
# Run immediately (no installation needed)
npx mcp-hub-tui

# Or install globally
npm install -g mcp-hub-tui
mcp-hub-tui
```

### First Launch

1. **Launch the TUI**: Run `npx mcp-hub-tui`
2. **Add Your First MCP**: Press `A` to add an MCP
3. **Choose MCP Type**: Select Command, SSE Server, or JSON Config
4. **Configure**: Fill in the details (name, command, etc.)
5. **Activate**: Press `Space` to enable/disable MCPs for Claude Code

## 🎮 Core Features

### 📋 **MCP Inventory Management**
- **Add MCPs**: Support for Command/Binary, SSE Server, and JSON configurations
- **Edit & Delete**: Full CRUD operations with confirmation dialogs
- **Search & Filter**: Real-time search across your MCP collection
- **Visual Status**: Clear indicators for active/inactive MCPs

### 🔄 **Claude Code Integration**
- **Status Detection**: Automatically detect Claude CLI availability
- **One-Key Toggle**: Space bar to enable/disable MCPs instantly
- **Project Context**: See current project and active MCP count
- **Sync Status**: Visual indicators for configuration sync state

### 🎯 **Smart Interface**
- **Responsive Layout**: Adapts from 1-4 columns based on terminal width
- **Keyboard Navigation**: Full keyboard control with intuitive shortcuts
- **Modal Workflows**: Progressive forms for complex operations
- **Loading States**: Clear feedback for all background operations

## 🔧 How It Works

### MCP Types Supported

**Command/Binary MCPs**
```
Name: github-mcp
Command: github-mcp
Args: --config ~/.config/github-mcp.json
```

**SSE Server MCPs**
```
Name: web-search
URL: http://localhost:3001/sse
```

**JSON Configuration MCPs**
```
Name: custom-mcp
Config: { "server": "localhost", "port": 3000 }
```

### Workflow

1. **Discover**: Find MCPs from various sources
2. **Inventory**: Add MCPs to your personal collection
3. **Focus**: Enable only the MCPs you need for your current task
4. **Code**: Use Claude Code with perfectly focused context
5. **Switch**: Change MCP configuration for different projects

## ⌨️ Keyboard Shortcuts

### Navigation
- `↑↓` or `j/k` - Navigate items
- `←→` or `h/l` - Switch columns
- `Tab` - Jump to search

### Actions
- `A` - Add new MCP
- `E` - Edit selected MCP
- `D` - Delete selected MCP
- `Space` - Toggle MCP active/inactive
- `/` - Search MCPs
- `R` - Refresh status
- `q` or `Esc` - Exit/Cancel

## 🏗️ Technical Architecture

### Built With
- **Go** - High-performance, cross-platform binary
- **Bubble Tea** - Reactive TUI framework
- **NPX Distribution** - Zero-friction installation

### Data Storage
- **Local JSON** - `~/.config/mcp-hub/inventory.json`
- **Atomic Operations** - Safe, concurrent access
- **Version Management** - Forward-compatible configuration

### Platform Support
- **macOS** - ARM64 & Intel
- **Linux** - ARM64 & x86_64
- **Windows** - ARM64 & x86_64

## 🎨 UI/UX Design

### Responsive Layout
- **Wide (120+ chars)**: 4-column layout for maximum information
- **Medium (80-120 chars)**: 3-column balanced view
- **Narrow (<80 chars)**: 2-column optimized layout

### Visual Language
- **Status Indicators**: ● (active) ○ (inactive)
- **Type Badges**: [CMD] [SSE] [JSON]
- **Loading States**: Animated spinners with progress messages
- **Color Coding**: Consistent visual hierarchy

## 🔐 Security & Privacy

### Data Protection
- **Local Storage**: All data stays on your machine
- **File Permissions**: Config files secured (0600)
- **Input Validation**: Sanitized user input
- **No Telemetry**: Zero data collection

### Command Safety
- **Validation**: MCP commands validated before execution
- **Sandboxing**: Safe execution patterns
- **User Control**: Explicit confirmation for destructive actions

## 🚀 Performance

### Optimizations
- **Fast Startup**: < 200ms launch time
- **Efficient Rendering**: Minimal screen updates
- **Memory Management**: Bounded resource usage
- **Responsive Input**: Sub-100ms key response

### Benchmarks
- **Large Inventories**: 1000+ MCPs supported
- **Search Performance**: Real-time filtering
- **Terminal Compatibility**: Works in any ANSI terminal

## 🛠️ Development

### Project Structure
```
mcp-hub/
├── bin/           # Platform-specific binaries
├── internal/      # Core application logic
│   └── ui/        # TUI components and services
├── docs/          # Architecture and specifications
└── index.js       # NPX entry point
```

### Testing
- **Unit Tests**: 85%+ coverage
- **Integration Tests**: End-to-end workflows
- **Platform Tests**: Cross-platform compatibility
- **Performance Tests**: Load and stress testing

## 🤝 Contributing

### Development Setup
```bash
# Clone repository
git clone https://github.com/gabadi/cc-mcp-manager.git
cd cc-mcp-manager

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o mcp-hub main.go
```

### Architecture
- **Clean Architecture**: Separation of concerns
- **Service Layer**: Business logic isolation
- **Component System**: Reusable UI components
- **State Management**: Centralized application state

## 📄 License

MIT License - See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- **Bubble Tea** - Excellent TUI framework
- **Claude Team** - MCP specification and tooling
- **Go Community** - Robust language and ecosystem

---

**Ready to transform your MCP management experience?**

```bash
npx mcp-hub-tui
```

*Built with ❤️ for developers who love efficient workflows*