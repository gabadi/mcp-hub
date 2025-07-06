// Package main provides the mcp-hub CLI application for managing Claude MCP configurations.
package main

import (
	"log"
	"os"

	"mcp-hub/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := runApp(); err != nil {
		os.Exit(1)
	}
}

func runApp() error {
	// Redirect log output to a file to prevent interference with TUI
	var logFile *os.File
	if file, err := os.OpenFile("/tmp/mcp-hub.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err == nil {
		logFile = file
		log.SetOutput(logFile)
		defer func() {
			_ = logFile.Close()
		}()
	}

	model := ui.NewModel()

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Printf("Error running MCP Manager: %v", err)
		return err
	}
	return nil
}
