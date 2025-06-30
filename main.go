package main

import (
	"cc-mcp-manager/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	model := ui.NewModel()

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal("Error running MCP Manager:", err)
	}
}
