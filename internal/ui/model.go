package ui

import (
	"fmt"
	"log"

	"cc-mcp-manager/pkg/config"
	"cc-mcp-manager/pkg/models"
	"cc-mcp-manager/pkg/storage"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the TUI application state
type Model struct {
	config    *config.Config
	storage   storage.InventoryStorage
	inventory *models.MCPInventory
	ready     bool
	err       error
}

// NewModel creates a new TUI model with storage initialization
func NewModel() *Model {
	// Initialize configuration
	cfg, err := config.New()
	if err != nil {
		log.Printf("Failed to initialize config: %v", err)
		return &Model{err: err}
	}

	// Initialize storage
	storageImpl := storage.NewJSONStorage(cfg)

	// Load inventory on startup (AC2: Inventory loads automatically on startup)
	inventory, err := storageImpl.LoadWithRecovery()
	if err != nil {
		log.Printf("Failed to load inventory: %v", err)
		// Create new inventory if loading fails
		inventory = models.NewInventory()
	}

	// Save initial inventory to ensure JSON file is created (AC1: JSON file created in appropriate config directory)
	if !storageImpl.Exists() {
		if err := storageImpl.Save(inventory); err != nil {
			log.Printf("Failed to save initial inventory: %v", err)
		} else {
			log.Printf("Created initial config file: %s", storageImpl.GetPath())
		}
	}

	log.Printf("Initialized MCP Manager with %d MCPs", len(inventory.MCPs))

	return &Model{
		config:    cfg,
		storage:   storageImpl,
		inventory: inventory,
		ready:     true,
	}
}

// Init initializes the TUI model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles TUI events
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the TUI
func (m *Model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n\nPress 'q' to quit.\n"
	}

	if !m.ready {
		return "Loading...\n"
	}

	s := "MCP Manager - Local Storage System\n\n"
	s += "=== MCP Inventory ===\n\n"
	
	if len(m.inventory.MCPs) == 0 {
		s += "No MCPs configured yet.\n"
	} else {
		for _, mcp := range m.inventory.MCPs {
			s += "• " + mcp.Name + " (" + string(mcp.Type) + ")\n"
		}
	}
	
	s += "\n=== Configuration ===\n"
	s += "Config File: " + m.config.GetConfigFile() + "\n"
	s += "MCPs Count: " + fmt.Sprintf("%d", len(m.inventory.MCPs)) + "\n"
	s += "Version: " + m.inventory.Version + "\n"
	
	s += "\nPress 'q' to quit.\n"
	
	return s
}

// GetInventory returns the current inventory
func (m *Model) GetInventory() *models.MCPInventory {
	return m.inventory
}

// GetStorage returns the storage implementation
func (m *Model) GetStorage() storage.InventoryStorage {
	return m.storage
}