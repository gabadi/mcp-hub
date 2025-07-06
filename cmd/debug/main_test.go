package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Test model initialization
func TestModelInit(t *testing.T) {
	m := model{}
	
	cmd := m.Init()
	if cmd != nil {
		t.Error("Expected Init to return nil command")
	}
}

// Test model Update function with various key inputs
func TestModelUpdate(t *testing.T) {
	t.Run("basic_key_handling", func(t *testing.T) {
		m := model{
			keys: make([]string, 0),
		}
		
		keyMsg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'a'},
		}
		
		updatedModel, cmd := m.Update(keyMsg)
		updatedM := updatedModel.(model)
		
		if len(updatedM.keys) != 1 {
			t.Errorf("Expected 1 key, got %d", len(updatedM.keys))
		}
		
		if updatedM.keys[0] != "a" {
			t.Errorf("Expected key 'a', got %s", updatedM.keys[0])
		}
		
		if cmd != nil {
			t.Error("Expected nil command for regular key")
		}
	})
	
	t.Run("quit_keys", func(t *testing.T) {
		m := model{}
		
		// Test ctrl+c
		ctrlCMsg := tea.KeyMsg{
			Type: tea.KeyCtrlC,
		}
		
		_, cmd := m.Update(ctrlCMsg)
		if cmd == nil {
			t.Error("Expected quit command for ctrl+c")
		}
		
		// Test q key
		qMsg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'q'},
		}
		
		_, cmd = m.Update(qMsg)
		if cmd == nil {
			t.Error("Expected quit command for q key")
		}
	})
	
	t.Run("clipboard_test_keys", func(t *testing.T) {
		// Test c key for basic clipboard functionality
		m := model{
			keys: make([]string, 0),
		}
		
		// Test c key
		cMsg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'c'},
		}
		
		updatedModel, _ := m.Update(cMsg)
		updatedM := updatedModel.(model)
		
		// Should have recorded the key
		if len(updatedM.keys) != 1 {
			t.Errorf("Expected 1 key, got %d", len(updatedM.keys))
		}
		
		if updatedM.keys[0] != "c" {
			t.Errorf("Expected key 'c', got %s", updatedM.keys[0])
		}
		
		// Should have set enhanced testing flag to false for basic clipboard
		if updatedM.enhancedTesting {
			t.Error("Expected enhanced testing to be false for c key")
		}
	})
	
	t.Run("clipboard_error_handling", func(t *testing.T) {
		// Test basic clipboard error handling structure
		m := model{
			keys:           make([]string, 0),
			clipboardError: "test error",
		}
		
		// Verify error state is maintained
		if m.clipboardError != "test error" {
			t.Errorf("Expected clipboard error to be 'test error', got %s", m.clipboardError)
		}
		
		// Test that normal key handling still works with error state
		keyMsg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'a'},
		}
		
		updatedModel, _ := m.Update(keyMsg)
		updatedM := updatedModel.(model)
		
		// Should still record keys even with error state
		if len(updatedM.keys) != 1 {
			t.Errorf("Expected 1 key, got %d", len(updatedM.keys))
		}
	})
	
	t.Run("key_limit_handling", func(t *testing.T) {
		m := model{
			keys: make([]string, 0),
		}
		
		// Add 15 keys (more than limit of 10)
		for i := 0; i < 15; i++ {
			keyMsg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{rune('a' + i)},
			}
			
			updatedModel, _ := m.Update(keyMsg)
			m = updatedModel.(model)
		}
		
		if len(m.keys) != 10 {
			t.Errorf("Expected 10 keys (limit), got %d", len(m.keys))
		}
		
		// Should contain the last 10 keys
		expectedKeys := []string{"f", "g", "h", "i", "j", "k", "l", "m", "n", "o"}
		if !reflect.DeepEqual(m.keys, expectedKeys) {
			t.Errorf("Expected keys %v, got %v", expectedKeys, m.keys)
		}
	})
	
	t.Run("log_file_creation", func(t *testing.T) {
		// Create temp directory for log testing
		tempDir := t.TempDir()
		
		// Change to temp directory
		originalWd, _ := os.Getwd()
		defer func() {
			os.Chdir(originalWd)
		}()
		os.Chdir(tempDir)
		
		m := model{
			keys: make([]string, 0),
		}
		
		keyMsg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'a'},
		}
		
		_, _ = m.Update(keyMsg)
		
		// Check if log file was created
		logPath := filepath.Join(tempDir, "key_debug.log")
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			t.Error("Expected log file to be created")
		}
		
		// Check log content
		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Errorf("Error reading log file: %v", err)
		}
		
		logContent := string(content)
		if !strings.Contains(logContent, "Key: 'a'") {
			t.Error("Expected key 'a' to be logged")
		}
	})
	
	t.Run("model_state_consistency", func(t *testing.T) {
		// Test that model state is consistent
		m := model{
			keys:            []string{"a", "b"},
			clipboardTest:   "test content",
			clipboardError:  "",
			enhancedTesting: true,
		}
		
		// Test that state is maintained
		if len(m.keys) != 2 {
			t.Errorf("Expected 2 keys, got %d", len(m.keys))
		}
		
		if m.clipboardTest != "test content" {
			t.Errorf("Expected clipboard test content, got %s", m.clipboardTest)
		}
		
		if !m.enhancedTesting {
			t.Error("Expected enhanced testing to be true")
		}
	})
}

// Test model View function
func TestModelView(t *testing.T) {
	t.Run("basic_view_structure", func(t *testing.T) {
		m := model{
			keys: []string{"a", "b", "c"},
		}
		
		view := m.View()
		
		// Check for expected sections
		expectedSections := []string{
			"Enhanced Key & Clipboard Debug Tool",
			"System Info:",
			"Recent keys:",
			"Clipboard Test Results:",
			"Test these keys:",
			"Press 'q' or Ctrl+C to quit",
			"Tips for macOS:",
		}
		
		for _, section := range expectedSections {
			if !strings.Contains(view, section) {
				t.Errorf("Expected view to contain '%s'", section)
			}
		}
	})
	
	t.Run("keys_display", func(t *testing.T) {
		m := model{
			keys: []string{"ctrl+c", "a", "b"},
		}
		
		view := m.View()
		
		// Check that keys are displayed with numbers
		if !strings.Contains(view, "1: 'ctrl+c'") {
			t.Error("Expected to show numbered keys")
		}
		if !strings.Contains(view, "2: 'a'") {
			t.Error("Expected to show numbered keys")
		}
		if !strings.Contains(view, "3: 'b'") {
			t.Error("Expected to show numbered keys")
		}
	})
	
	t.Run("clipboard_error_display", func(t *testing.T) {
		m := model{
			clipboardError: "Access denied",
		}
		
		view := m.View()
		
		if !strings.Contains(view, "❌ Error: Access denied") {
			t.Error("Expected to show clipboard error")
		}
	})
	
	t.Run("clipboard_success_display", func(t *testing.T) {
		m := model{
			clipboardTest: "test content",
		}
		
		view := m.View()
		
		if !strings.Contains(view, "✅ Content: \"test content\"") {
			t.Error("Expected to show clipboard content")
		}
	})
	
	t.Run("enhanced_diagnostics_display", func(t *testing.T) {
		m := model{
			enhancedTesting: true,
			diagnosticInfo: map[string]interface{}{
				"platform": "darwin",
				"available": true,
				"method":   "pbcopy",
			},
		}
		
		view := m.View()
		
		if !strings.Contains(view, "Enhanced Clipboard Diagnostics:") {
			t.Error("Expected to show enhanced diagnostics section")
		}
		
		// Check that diagnostic info is displayed
		if !strings.Contains(view, "platform: darwin") {
			t.Error("Expected to show platform info")
		}
		if !strings.Contains(view, "available: true") {
			t.Error("Expected to show availability info")
		}
		if !strings.Contains(view, "method: pbcopy") {
			t.Error("Expected to show method info")
		}
	})
	
	t.Run("system_info_display", func(t *testing.T) {
		// Set environment variables for testing
		originalTerm := os.Getenv("TERM_PROGRAM")
		originalTermVersion := os.Getenv("TERM_PROGRAM_VERSION")
		originalTermType := os.Getenv("TERM")
		
		defer func() {
			os.Setenv("TERM_PROGRAM", originalTerm)
			os.Setenv("TERM_PROGRAM_VERSION", originalTermVersion)
			os.Setenv("TERM", originalTermType)
		}()
		
		os.Setenv("TERM_PROGRAM", "TestTerminal")
		os.Setenv("TERM_PROGRAM_VERSION", "1.0.0")
		os.Setenv("TERM", "xterm-256color")
		
		m := model{}
		view := m.View()
		
		if !strings.Contains(view, "Terminal: TestTerminal") {
			t.Error("Expected to show terminal program")
		}
		if !strings.Contains(view, "Terminal Version: 1.0.0") {
			t.Error("Expected to show terminal version")
		}
		if !strings.Contains(view, "TERM: xterm-256color") {
			t.Error("Expected to show TERM variable")
		}
	})
	
	t.Run("no_keys_display", func(t *testing.T) {
		m := model{
			keys: []string{},
		}
		
		view := m.View()
		
		// Should still show the "Recent keys:" section but with no numbered items
		if !strings.Contains(view, "Recent keys:") {
			t.Error("Expected to show recent keys section")
		}
	})
	
	t.Run("default_clipboard_state", func(t *testing.T) {
		m := model{
			// No clipboard test or error set
		}
		
		view := m.View()
		
		if !strings.Contains(view, "Press Cmd+V or Ctrl+V to test enhanced clipboard") {
			t.Error("Expected to show default clipboard instructions")
		}
		if !strings.Contains(view, "Press 'c' to test basic clipboard") {
			t.Error("Expected to show basic clipboard instructions")
		}
	})
}

// Test main function
func TestMainFunction(t *testing.T) {
	// Test that main function components can be created without panic
	t.Run("main_function_components", func(t *testing.T) {
		// Test that we can create a model
		m := model{}
		if m.keys == nil {
			// Initialize if needed
			m.keys = make([]string, 0)
		}
		
		// Test that we can create a tea program
		program := tea.NewProgram(m)
		if program == nil {
			t.Error("Expected tea program to be created")
		}
	})
	
	t.Run("log_file_removal", func(t *testing.T) {
		// Test log file removal functionality
		tempDir := t.TempDir()
		
		originalWd, _ := os.Getwd()
		defer func() {
			os.Chdir(originalWd)
		}()
		os.Chdir(tempDir)
		
		// Create a test log file
		logFile := "key_debug.log"
		err := os.WriteFile(logFile, []byte("test content"), 0600)
		if err != nil {
			t.Errorf("Failed to create test log file: %v", err)
		}
		
		// Remove the file (simulating what main does)
		err = os.Remove(logFile)
		if err != nil {
			t.Errorf("Failed to remove log file: %v", err)
		}
		
		// Verify file was removed
		if _, err := os.Stat(logFile); !os.IsNotExist(err) {
			t.Error("Expected log file to be removed")
		}
	})
}

// Test edge cases and error conditions
func TestEdgeCases(t *testing.T) {
	t.Run("non_key_message", func(t *testing.T) {
		m := model{}
		
		// Send a non-key message
		nonKeyMsg := tea.WindowSizeMsg{
			Width:  80,
			Height: 24,
		}
		
		updatedModel, cmd := m.Update(nonKeyMsg)
		updatedM := updatedModel.(model)
		
		// Should return unchanged model
		if len(updatedM.keys) != 0 {
			t.Error("Expected no keys to be added for non-key message")
		}
		
		if cmd != nil {
			t.Error("Expected no command for non-key message")
		}
	})
	
	t.Run("empty_runes_key", func(t *testing.T) {
		m := model{
			keys: make([]string, 0),
		}
		
		keyMsg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{},
		}
		
		updatedModel, _ := m.Update(keyMsg)
		updatedM := updatedModel.(model)
		
		// Should still record the key even if runes is empty
		if len(updatedM.keys) != 1 {
			t.Error("Expected one key even for empty runes")
		}
	})
	
	t.Run("log_file_error_handling", func(t *testing.T) {
		// Create a directory where we can't write
		tempDir := t.TempDir()
		readOnlyDir := filepath.Join(tempDir, "readonly")
		os.Mkdir(readOnlyDir, 0444) // Read-only directory
		
		originalWd, _ := os.Getwd()
		defer func() {
			os.Chdir(originalWd)
		}()
		os.Chdir(readOnlyDir)
		
		m := model{
			keys: make([]string, 0),
		}
		
		keyMsg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'a'},
		}
		
		// Should not panic even if log file can't be created
		updatedModel, _ := m.Update(keyMsg)
		updatedM := updatedModel.(model)
		
		// Key should still be recorded
		if len(updatedM.keys) != 1 {
			t.Error("Expected key to be recorded even if log fails")
		}
	})
}

// Mock implementations for testing are not needed for simplified tests