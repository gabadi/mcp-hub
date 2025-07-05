// Package main provides a debug tool for testing key bindings and clipboard functionality.
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"mcp-hub/internal/ui/services"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	keys            []string
	clipboardTest   string
	clipboardError  string
	diagnosticInfo  map[string]interface{}
	enhancedTesting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		keyStr := keyMsg.String()
		m.keys = append(m.keys, keyStr)

		// Enhanced logging
		f, err := os.OpenFile("key_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err == nil {
			_, _ = fmt.Fprintf(f, "Key: '%s' Type: %d Alt: %t Runes: %v\n", keyStr, keyMsg.Type, keyMsg.Alt, keyMsg.Runes)
			_, _ = fmt.Fprintf(f, "Terminal: %s\n", os.Getenv("TERM_PROGRAM"))
			_ = f.Close()
		}

		// Test clipboard functionality when cmd+v or ctrl+v is pressed
		if keyStr == "cmd+v" || keyStr == "ctrl+v" || keyStr == "⌘v" || keyStr == "command+v" {
			// Test with enhanced clipboard service
			clipboardService := services.NewClipboardService()
			m.diagnosticInfo = clipboardService.GetDiagnosticInfo()

			content, clipErr := clipboardService.EnhancedPaste()
			if clipErr != nil {
				m.clipboardError = clipErr.Error()
				m.clipboardTest = ""
			} else {
				m.clipboardError = ""
				m.clipboardTest = content
			}
			m.enhancedTesting = true
		}

		// Test basic clipboard on 'c' key
		if keyStr == "c" {
			content, clipErr := clipboard.ReadAll()
			if clipErr != nil {
				m.clipboardError = "Basic test: " + clipErr.Error()
				m.clipboardTest = ""
			} else {
				m.clipboardError = ""
				m.clipboardTest = "Basic test: " + content
			}
			m.enhancedTesting = false
		}

		if keyStr == "ctrl+c" || keyStr == "q" {
			return m, tea.Quit
		}

		// Keep only last 10 keys
		if len(m.keys) > 10 {
			m.keys = m.keys[len(m.keys)-10:]
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "🔍 Enhanced Key & Clipboard Debug Tool for mcp-hub\n\n"
	s += "System Info:\n"
	s += fmt.Sprintf("• OS: %s\n", runtime.GOOS)
	s += fmt.Sprintf("• Terminal: %s\n", os.Getenv("TERM_PROGRAM"))
	s += fmt.Sprintf("• Terminal Version: %s\n", os.Getenv("TERM_PROGRAM_VERSION"))
	s += fmt.Sprintf("• TERM: %s\n", os.Getenv("TERM"))
	s += "\n"

	s += "Recent keys:\n"
	for i, key := range m.keys {
		s += fmt.Sprintf("%d: '%s'\n", i+1, key)
	}

	s += "\n📋 Clipboard Test Results:\n"
	switch {
	case m.clipboardError != "":
		s += fmt.Sprintf("❌ Error: %s\n", m.clipboardError)
	case m.clipboardTest != "":
		s += fmt.Sprintf("✅ Content: %q\n", m.clipboardTest)
	default:
		s += "⏳ Press Cmd+V or Ctrl+V to test enhanced clipboard\n"
		s += "⏳ Press 'c' to test basic clipboard\n"
	}

	if m.enhancedTesting && len(m.diagnosticInfo) > 0 {
		s += "\n🔍 Enhanced Clipboard Diagnostics:\n"
		for key, value := range m.diagnosticInfo {
			s += fmt.Sprintf("• %s: %v\n", key, value)
		}
	}

	s += "\n📋 Test these keys:\n"
	s += "• Ctrl+V (should show 'ctrl+v')\n"
	s += "• Cmd+V (should show 'cmd+v' on macOS)\n"
	s += "• ⌘v (alternative Command+V representation)\n"
	s += "• command+v (another Command+V representation)\n"
	s += "• c (test basic clipboard without key combos)\n"
	s += "\nPress 'q' or Ctrl+C to quit\n"
	s += "\n💡 Tips for macOS:\n"
	s += "• Grant Warp accessibility permissions in System Settings\n"
	s += "• Try restarting Warp after permission changes\n"
	s += "• Check if pbpaste/pbcopy work: 'echo test | pbcopy && pbpaste'\n"
	return s
}

func main() {
	_ = os.Remove("key_debug.log")
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
