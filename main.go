// Package main provides the mcp-hub CLI application for managing Claude MCP configurations.
package main

import (
	"log"
	"os"
	"path/filepath"

	"mcp-hub/internal/platform"
	"mcp-hub/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := runApp(); err != nil {
		os.Exit(1)
	}
}

func runApp() error {
	// Create platform service for dynamic path resolution
	platformService := platform.NewPlatformServiceFactoryDefault().CreatePlatformService()

	// Redirect log output to a platform-specific file to prevent interference with TUI
	var logFile *os.File
	logPath := filepath.Join(platformService.GetLogPath(), "mcp-hub.log")
	
	// Ensure log directory exists
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, platformService.GetDefaultDirectoryPermissions()); err == nil {
		if file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, platformService.GetDefaultFilePermissions()); err == nil {
			logFile = file
			log.SetOutput(logFile)
			defer func() {
				_ = logFile.Close()
			}()
		}
	}

	model := ui.NewModel()

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Printf("Error running MCP Manager: %v", err)
		return err
	}
	return nil
}
