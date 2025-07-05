package components

import (
	"fmt"
	"strings"

	"cc-mcp-manager/internal/ui/services"
	"cc-mcp-manager/internal/ui/types"

	"github.com/charmbracelet/lipgloss"
)

// Modal string constants
const (
	// Button/instruction text
	EscCancelText = "ESC=Cancel"
	EditInstructionText = "[Tab] Next Field • [Enter] Update • ESC Cancel"
	AddInstructionText = "[Tab] Next Field • [Enter] Add • ESC Cancel"
	
	// Form field labels
	NameRequiredLabel = "Name: (required)"
	EnvironmentOptionalLabel = "Environment: (optional)"
)

// OverlayModal renders a modal on top of existing content
func OverlayModal(model types.Model, width, height int, _ string) string {
	modalWidth, modalHeight := calculateModalDimensions(model.ActiveModal, width, height)
	modalStyle := createModalStyle(modalWidth, modalHeight)
	title, content, footer := getModalContent(model)
	modalContent := buildModalContent(title, content, footer)
	modal := modalStyle.Render(modalContent)
	
	return centerModal(modal, width, height)
}

func calculateModalDimensions(activeModal types.ModalType, width, height int) (int, int) {
	modalWidth := 60
	modalHeight := 20

	// Adjust dimensions for different modal types
	switch activeModal {
	case types.NoModal:
		// Default modal size, no adjustment needed
	case types.AddModal:
		// Default modal size, no adjustment needed
	case types.AddMCPTypeSelection:
		modalHeight = 18 // Smaller for type selection
	case types.AddCommandForm:
		modalHeight = 22 // Slightly larger for 3 fields
	case types.AddSSEForm:
		modalHeight = 20 // Standard form size
	case types.AddJSONForm:
		modalHeight = 25 // Larger for JSON text area
	case types.EditModal:
		modalHeight = 15 // Smaller for edit confirmation
	case types.DeleteModal:
		modalHeight = 12 // Smaller for delete confirmation
	}

	if modalWidth > width-10 {
		modalWidth = width - 10
	}
	if modalHeight > height-10 {
		modalHeight = height - 10
	}

	return modalWidth, modalHeight
}

func createModalStyle(modalWidth, modalHeight int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(1).
		Width(modalWidth).
		Height(modalHeight).
		Background(lipgloss.Color("#1E1E2E")).
		Foreground(lipgloss.Color("#FFFFFF"))
}

func getModalContent(model types.Model) (string, string, string) {
	var title, content, footer string

	switch model.ActiveModal {
	case types.NoModal:
		// For legacy support - treat NoModal as basic AddModal
		title = "Add New MCP"
		content = "Select type of MCP to add"
		footer = EscCancelText
	case types.AddModal:
		title = "Add New MCP"
		content = "Select type of MCP to add"
		footer = EscCancelText
	case types.AddMCPTypeSelection:
		title = "Add New MCP - Select Type"
		content = renderTypeSelectionContent(model)
		footer = "[1-3] Select • ESC Cancel"
	case types.AddCommandForm:
		title, footer = getFormTitleAndFooter("Command/Binary", model.EditMode, model.EditMCPName)
		content = renderCommandFormContent(model)
	case types.AddSSEForm:
		title, footer = getFormTitleAndFooter("SSE Server", model.EditMode, model.EditMCPName)
		content = renderSSEFormContent(model)
	case types.AddJSONForm:
		title, footer = getFormTitleAndFooter("JSON Configuration", model.EditMode, model.EditMCPName)
		content = renderJSONFormContent(model)
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
		footer = EscCancelText
	}

	return title, content, footer
}

func getFormTitleAndFooter(formType string, editMode bool, editMCPName string) (string, string) {
	if editMode {
		title := fmt.Sprintf("Edit MCP - %s: %s", formType, editMCPName)
		return title, EditInstructionText
	}
	title := fmt.Sprintf("Add New MCP - %s", formType)
	return title, AddInstructionText
}

func buildModalContent(title, content, footer string) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		MarginBottom(1)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginTop(1)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		content,
		footerStyle.Render(footer),
	)
}

func centerModal(modal string, width, height int) string {
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
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
	nameLabel := NameRequiredLabel
	nameValue := model.FormData.Name
	if model.FormData.ActiveField == 0 {
		nameValue += "_"  // Show cursor
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
		commandValue += "_"
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
		argsValue += "_"
		argsLabel = "> " + argsLabel
	}
	lines = append(lines, argsLabel)
	lines = append(lines, fmt.Sprintf("[%s]", argsValue))
	lines = append(lines, "")

	// Environment Variables field
	envLabel := EnvironmentOptionalLabel
	envValue := model.FormData.Environment
	if model.FormData.ActiveField == 3 {
		envValue += "_"
		envLabel = "> " + envLabel
	}
	lines = append(lines, envLabel)
	lines = append(lines, fmt.Sprintf("[%s]", envValue))
	lines = append(lines, "Format: KEY1=value1,KEY2=value2")

	return strings.Join(lines, "\n")
}

// renderSSEFormContent renders the SSE Server MCP form
func renderSSEFormContent(model types.Model) string {
	var lines []string

	// Name field
	nameLabel := NameRequiredLabel
	nameValue := model.FormData.Name
	if model.FormData.ActiveField == 0 {
		nameValue += "_"
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
		urlValue += "_"
		urlLabel = "> " + urlLabel
	}
	lines = append(lines, urlLabel)
	lines = append(lines, fmt.Sprintf("[%s]", urlValue))
	if err, exists := model.FormErrors["url"]; exists {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("  Error: "+err))
	}
	lines = append(lines, "")

	// Environment Variables field
	envLabel := EnvironmentOptionalLabel
	envValue := model.FormData.Environment
	if model.FormData.ActiveField == 2 {
		envValue += "_"
		envLabel = "> " + envLabel
	}
	lines = append(lines, envLabel)
	lines = append(lines, fmt.Sprintf("[%s]", envValue))
	lines = append(lines, "Format: KEY1=value1,KEY2=value2")
	lines = append(lines, "")
	lines = append(lines, "Enter a valid HTTP/HTTPS URL for the SSE server.")

	return strings.Join(lines, "\n")
}

// renderJSONFormContent renders the JSON Configuration MCP form
func renderJSONFormContent(model types.Model) string {
	var lines []string

	// Name field
	nameLabel := NameRequiredLabel
	nameValue := model.FormData.Name
	if model.FormData.ActiveField == 0 {
		nameValue += "_"
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
		jsonValue += "_"
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
	lines = append(lines, "")

	// Environment Variables field
	envLabel := EnvironmentOptionalLabel
	envValue := model.FormData.Environment
	if model.FormData.ActiveField == 2 {
		envValue += "_"
		envLabel = "> " + envLabel
	}
	lines = append(lines, envLabel)
	lines = append(lines, fmt.Sprintf("[%s]", envValue))
	lines = append(lines, "Format: KEY1=value1,KEY2=value2")

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
