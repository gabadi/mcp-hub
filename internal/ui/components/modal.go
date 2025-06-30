package components

import (
	"fmt"
	"strings"

	"cc-mcp-manager/internal/ui/types"
	"github.com/charmbracelet/lipgloss"
)

// ModalType represents the type of modal being displayed
type ModalType int

const (
	AddModal ModalType = iota
	EditModal
	DeleteModal
)

// RenderModal creates a modal overlay for add/edit/delete operations (DEPRECATED - use OverlayModal)
func RenderModal(model types.Model, modalType ModalType, width, height int) string {
	// Modal dimensions
	modalWidth := 60
	modalHeight := 20
	if modalWidth > width-10 {
		modalWidth = width - 10
	}
	if modalHeight > height-10 {
		modalHeight = height - 10
	}

	// Modal style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(1).
		Width(modalWidth).
		Height(modalHeight).
		Background(lipgloss.Color("#1E1E2E"))

	// Title based on modal type
	var title string
	var content string

	switch modalType {
	case AddModal:
		title = "Add New MCP"
		content = renderAddModalContent()
	case EditModal:
		title = "Edit MCP"
		content = renderEditModalContent(model)
	case DeleteModal:
		title = "Delete MCP"
		content = renderDeleteModalContent(model)
	}

	// Title style
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		MarginBottom(1)

	// Footer with instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginTop(1)
	footer := "Enter=Confirm • ESC=Cancel"

	// Combine all elements
	modalContent := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		content,
		footerStyle.Render(footer),
	)

	// Center the modal on screen
	centeredModal := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(modalContent),
	)

	// Create semi-transparent overlay effect by dimming the background
	overlayStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Background(lipgloss.Color("#000000"))

	return overlayStyle.Render(centeredModal)
}

// OverlayModal renders a modal on top of existing content
func OverlayModal(model types.Model, modalType ModalType, width, height int, backgroundContent string) string {
	// For now, just render the modal without the background overlay
	// TODO: Fix overlay to properly show dimmed background

	// Modal dimensions
	modalWidth := 60
	modalHeight := 20
	if modalWidth > width-10 {
		modalWidth = width - 10
	}
	if modalHeight > height-10 {
		modalHeight = height - 10
	}

	// Modal style with solid background to ensure it's opaque
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(1).
		Width(modalWidth).
		Height(modalHeight).
		Background(lipgloss.Color("#1E1E2E")).
		Foreground(lipgloss.Color("#FFFFFF"))

	// Title based on modal type
	var title string
	var content string

	switch modalType {
	case AddModal:
		title = "Add New MCP"
		content = renderAddModalContent()
	case EditModal:
		title = "Edit MCP"
		content = renderEditModalContent(model)
	case DeleteModal:
		title = "Delete MCP"
		content = renderDeleteModalContent(model)
	}

	// Title style
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		MarginBottom(1)

	// Footer with instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginTop(1)
	footer := "Enter=Confirm • ESC=Cancel"

	// Combine all elements
	modalContent := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		content,
		footerStyle.Render(footer),
	)

	// Create the modal
	modal := modalStyle.Render(modalContent)

	// Use lipgloss.Place to center the modal
	centeredModal := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
	
	return centeredModal
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func renderAddModalContent() string {
	fields := []string{
		"Name: [________________]",
		"Type: [SSE/CMD/JSON/HTTP]",
		"Command: [________________________________]",
		"Active: [ ] Yes  [X] No",
		"",
		"Fill in the fields above to add a new MCP.",
	}
	return strings.Join(fields, "\n")
}

func renderEditModalContent(model types.Model) string {
	// Get selected MCP if available
	filteredMCPs := GetFilteredMCPs(model)
	selectedIndex := model.SelectedItem
	if model.SearchQuery != "" {
		selectedIndex = model.FilteredSelectedIndex
	}

	if selectedIndex >= len(filteredMCPs) {
		return "No MCP selected"
	}

	item := filteredMCPs[selectedIndex]
	activeCheck := "[ ]"
	inactiveCheck := "[X]"
	if item.Active {
		activeCheck = "[X]"
		inactiveCheck = "[ ]"
	}

	fields := []string{
		fmt.Sprintf("Name: [%s]", item.Name),
		fmt.Sprintf("Type: [%s]", item.Type),
		fmt.Sprintf("Command: [%s]", item.Command),
		fmt.Sprintf("Active: %s Yes  %s No", activeCheck, inactiveCheck),
		"",
		"Edit the fields above and press Enter to save.",
	}
	return strings.Join(fields, "\n")
}

func renderDeleteModalContent(model types.Model) string {
	// Get selected MCP if available
	filteredMCPs := GetFilteredMCPs(model)
	selectedIndex := model.SelectedItem
	if model.SearchQuery != "" {
		selectedIndex = model.FilteredSelectedIndex
	}

	if selectedIndex >= len(filteredMCPs) {
		return "No MCP selected"
	}

	item := filteredMCPs[selectedIndex]

	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true)

	content := []string{
		fmt.Sprintf("Are you sure you want to delete this MCP?"),
		"",
		fmt.Sprintf("Name: %s", item.Name),
		fmt.Sprintf("Type: %s", item.Type),
		fmt.Sprintf("Command: %s", item.Command),
		"",
		warningStyle.Render("This action cannot be undone!"),
		"",
		"Press Enter to confirm deletion, or ESC to cancel.",
	}
	return strings.Join(content, "\n")
}