// Package main provides the MCP Manager CLI application.
//
// MCP Manager is a Terminal User Interface (TUI) application for managing
// Model Context Protocol (MCP) servers. It provides an intuitive, responsive
// interface for browsing, configuring, and monitoring MCP servers.
//
// # Features
//
//   - Responsive TUI that adapts to different terminal sizes
//   - Intuitive keyboard navigation with vim-style bindings
//   - Real-time search and filtering
//   - Comprehensive error handling with recovery options
//   - Configuration management with validation
//   - Status notifications and user feedback
//
// # Usage
//
// Basic usage:
//   mcp-manager
//
// The application will launch a full-screen TUI interface. Use the keyboard
// shortcuts displayed in the header to navigate and interact with MCP servers.
//
// # Keyboard Shortcuts
//
// Navigation:
//   - Arrow keys or h/j/k/l: Navigate between items and columns
//   - Tab: Enter search mode
//   - Esc: Exit current mode or quit application
//   - Q: Quit application
//
// Actions:
//   - A: Add new MCP server
//   - E: Edit selected MCP server
//   - D: Delete selected MCP server
//   - Space: Toggle server status
//   - R: Refresh server list
//
// Error Handling:
//   - R: Attempt error recovery (when available)
//   - Enter/Esc/Space: Dismiss error dialog
//   - Q: Quit on critical errors
//
// # Configuration
//
// The application looks for configuration files in the following locations:
//   - ~/.config/mcp-manager/config.json
//   - ~/.mcp/config.json
//   - ./mcp-config/config.json
//
// If no configuration file is found, default settings are used. The configuration
// file is automatically created with default values when first saving settings.
//
// # Architecture
//
// The application follows a modular architecture:
//
//   - cmd/mcp-manager: Entry point and CLI setup
//   - internal/tui: Terminal User Interface implementation
//   - internal/config: Configuration management
//   - internal/errors: Error handling and user messages
//   - internal/models: Data models for MCP servers
//   - pkg/mcp: MCP protocol implementation
//
// The TUI is built using the Bubble Tea framework and follows the Elm Architecture
// pattern for predictable state management.
//
// # Error Handling
//
// The application provides comprehensive error handling with:
//   - User-friendly error messages
//   - Automatic recovery options for common issues
//   - Graceful degradation for non-critical errors
//   - Clear guidance for resolution steps
//
// # Terminal Compatibility
//
// The application supports all ANSI-compatible terminals and requires:
//   - Minimum width: 40 columns (configurable)
//   - Minimum height: 10 lines
//   - Color support: 256 colors (recommended)
//   - Unicode support: For icons and special characters
//
// # Development
//
// For development and testing:
//   make build    # Build the application
//   make test     # Run all tests
//   make run      # Build and run
//   make dev      # Development mode with auto-rebuild
//
// # License
//
// This project is licensed under the MIT License.
package main