package main

import (
	"log"
	"os"

	"cc-mcp-manager/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Redirect log output to a file to prevent interference with TUI
	logFile, err := os.OpenFile("/tmp/cc-mcp-manager.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		log.SetOutput(logFile)
		defer logFile.Close()
	}

	model := ui.NewModel()

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal("Error running MCP Manager:", err)
	}
}
