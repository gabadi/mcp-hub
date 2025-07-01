package components

import (
	"fmt"
	"strings"

	"cc-mcp-manager/internal/ui/services"
	"cc-mcp-manager/internal/ui/types"

	"github.com/charmbracelet/lipgloss"
)

// OverlayModal renders a modal on top of existing content
func OverlayModal(model types.Model, width, height int, backgroundContent string) string {
	// Modal dimensions based on modal type
	modalWidth := 60
	modalHeight := 20

	// Adjust dimensions for different modal types
	switch model.ActiveModal {
	case types.AddJSONForm:
		modalHeight = 25 // Larger for JSON text area
	case types.AddCommandForm:
		modalHeight = 22 // Slightly larger for 3 fields
	}

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

	// Title and content based on modal type
	var title string
	var content string
	var footer string

	switch model.ActiveModal {
	case types.AddMCPTypeSelection:
		title = "Add New MCP - Select Type"
		content = renderTypeSelectionContent(model)
		footer = "[1-3] Select • ESC Cancel"
	case types.AddCommandForm:
		title = "Add New MCP - Command/Binary"
		content = renderCommandFormContent(model)
		footer = "[Tab] Next Field • [Enter] Add • ESC Cancel"
	case types.AddSSEForm:
		title = "Add New MCP - SSE Server"
		content = renderSSEFormContent(model)
		footer = "[Tab] Next Field • [Enter] Add • ESC Cancel"
	case types.AddJSONForm:
		title = "Add New MCP - JSON Configuration"
		content = renderJSONFormContent(model)
		footer = "[Tab] Next Field • [Enter] Add • ESC Cancel"
	case types.EditModal:
		title = "Edit MCP"
		content = renderEditModalContent(model)
		footer = "Enter=Confirm • ESC=Cancel"
	case types.DeleteModal:
		title = "Delete MCP"
		content = renderDeleteModalContent(model)
		footer = "Enter=Confirm • ESC=Cancel"
	default:
		title = "Unknown Modal"
		content = "Unknown modal type"
		footer = "ESC=Cancel"
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

// renderTypeSelectionContent renders the MCP type selection interface
func renderTypeSelectionContent(model types.Model) string {
	selectedOption := model.FormData.ActiveField
	if selectedOption == 0 {
		selectedOption = 1 // Default to first option
	}

	var lines []string
	lines = append(lines, "Choose the type of MCP to add:")
	lines = append(lines, "")

	// Option 1 - Command/Binary
	option1Style := lipgloss.NewStyle()
	if selectedOption == 1 {
		option1Style = option1Style.Background(lipgloss.Color("#7C3AED")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	}
	lines = append(lines, option1Style.Render("1. Command/Binary (most common)"))
	lines = append(lines, "   Execute MCP as a command or binary")
	lines = append(lines, "")

	// Option 2 - SSE Server
	option2Style := lipgloss.NewStyle()
	if selectedOption == 2 {
		option2Style = option2Style.Background(lipgloss.Color("#7C3AED")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	}
	lines = append(lines, option2Style.Render("2. SSE Server (HTTP/WebSocket)"))
	lines = append(lines, "   Connect to an SSE server endpoint")
	lines = append(lines, "")

	// Option 3 - JSON Configuration
	option3Style := lipgloss.NewStyle()
	if selectedOption == 3 {
		option3Style = option3Style.Background(lipgloss.Color("#7C3AED")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	}
	lines = append(lines, option3Style.Render("3. JSON Configuration"))
	lines = append(lines, "   Add MCP with custom JSON configuration")
	lines = append(lines, "")

	lines = append(lines, "Use number keys (1-3), arrow keys, or Enter to select.")

	return strings.Join(lines, "\n")
}

// renderCommandFormContent renders the Command/Binary MCP form
func renderCommandFormContent(model types.Model) string {
	var lines []string

	// Name field
	nameLabel := "Name: (required)"
	nameValue := model.FormData.Name
	if model.FormData.ActiveField == 0 {
		nameValue = nameValue + "_"  // Show cursor
		nameLabel = "> " + nameLabel // Show focus
	}
	lines = append(lines, nameLabel)
	lines = append(lines, fmt.Sprintf("[%s]", nameValue))
	if err, exists := model.FormErrors["name"]; exists {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("  Error: "+err))
	}
	lines = append(lines, "")

	// Command field
	commandLabel := "Command: (required)"
	commandValue := model.FormData.Command
	if model.FormData.ActiveField == 1 {
		commandValue = commandValue + "_"
		commandLabel = "> " + commandLabel
	}
	lines = append(lines, commandLabel)
	lines = append(lines, fmt.Sprintf("[%s]", commandValue))
	if err, exists := model.FormErrors["command"]; exists {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("  Error: "+err))
	}
	lines = append(lines, "")

	// Args field
	argsLabel := "Args: (optional)"
	argsValue := model.FormData.Args
	if model.FormData.ActiveField == 2 {
		argsValue = argsValue + "_"
		argsLabel = "> " + argsLabel
	}
	lines = append(lines, argsLabel)
	lines = append(lines, fmt.Sprintf("[%s]", argsValue))

	return strings.Join(lines, "\n")
}

// renderSSEFormContent renders the SSE Server MCP form
func renderSSEFormContent(model types.Model) string {
	var lines []string

	// Name field
	nameLabel := "Name: (required)"
	nameValue := model.FormData.Name
	if model.FormData.ActiveField == 0 {
		nameValue = nameValue + "_"
		nameLabel = "> " + nameLabel
	}
	lines = append(lines, nameLabel)
	lines = append(lines, fmt.Sprintf("[%s]", nameValue))
	if err, exists := model.FormErrors["name"]; exists {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("  Error: "+err))
	}
	lines = append(lines, "")

	// URL field
	urlLabel := "URL: (required)"
	urlValue := model.FormData.URL
	if model.FormData.ActiveField == 1 {
		urlValue = urlValue + "_"
		urlLabel = "> " + urlLabel
	}
	lines = append(lines, urlLabel)
	lines = append(lines, fmt.Sprintf("[%s]", urlValue))
	if err, exists := model.FormErrors["url"]; exists {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("  Error: "+err))
	}
	lines = append(lines, "")
	lines = append(lines, "Enter a valid HTTP/HTTPS URL for the SSE server.")

	return strings.Join(lines, "\n")
}

// renderJSONFormContent renders the JSON Configuration MCP form
func renderJSONFormContent(model types.Model) string {
	var lines []string

	// Name field
	nameLabel := "Name: (required)"
	nameValue := model.FormData.Name
	if model.FormData.ActiveField == 0 {
		nameValue = nameValue + "_"
		nameLabel = "> " + nameLabel
	}
	lines = append(lines, nameLabel)
	lines = append(lines, fmt.Sprintf("[%s]", nameValue))
	if err, exists := model.FormErrors["name"]; exists {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("  Error: "+err))
	}
	lines = append(lines, "")

	// JSON Config field
	jsonLabel := "JSON Configuration: (required)"
	jsonValue := model.FormData.JSONConfig
	if model.FormData.ActiveField == 1 {
		jsonValue = jsonValue + "_"
		jsonLabel = "> " + jsonLabel
	}
	lines = append(lines, jsonLabel)

	// Show JSON in a box with multiple lines
	jsonLines := strings.Split(jsonValue, "\n")
	for _, line := range jsonLines {
		lines = append(lines, fmt.Sprintf("  %s", line))
	}

	if err, exists := model.FormErrors["json"]; exists {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("  Error: "+err))
	} else if jsonValue != "" {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#51CF66")).Render("  ✓ Valid JSON"))
	}

	return strings.Join(lines, "\n")
}

func renderEditModalContent(model types.Model) string {
	// Get selected MCP if available
	filteredMCPs := services.GetFilteredMCPs(model)
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
	filteredMCPs := services.GetFilteredMCPs(model)
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
		"Are you sure you want to delete this MCP?",
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
