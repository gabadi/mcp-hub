package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"mcp-manager/internal/tui"
)

func main() {
	// Initialize the TUI model
	m := tui.NewModel()

	// Create the Bubble Tea program with alt screen
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running MCP Manager: %v\n", err)
		os.Exit(1)
	}
}