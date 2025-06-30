package main

import (
	"log"

	"cc-mcp-manager/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create the initial model
	model := ui.NewModel()

	// Create the program with alt screen for full terminal control
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		log.Fatal("Error running MCP Manager:", err)
	}
}
