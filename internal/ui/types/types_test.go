package types

import (
	"reflect"
	"testing"
	"time"
)

func TestNewModel(t *testing.T) {
	model := NewModel()

	// Test initial state
	if model.State != MainNavigation {
		t.Errorf("Expected State to be MainNavigation, got %v", model.State)
	}

	if model.ActiveColumn != 0 {
		t.Errorf("Expected ActiveColumn to be 0, got %d", model.ActiveColumn)
	}

	if model.SelectedItem != 0 {
		t.Errorf("Expected SelectedItem to be 0, got %d", model.SelectedItem)
	}

	if model.SearchQuery != "" {
		t.Errorf("Expected SearchQuery to be empty, got %q", model.SearchQuery)
	}

	if model.SearchActive {
		t.Error("Expected SearchActive to be false")
	}

	if model.SearchInputActive {
		t.Error("Expected SearchInputActive to be false")
	}

	if model.ColumnCount != 1 {
		t.Errorf("Expected ColumnCount to be 1, got %d", model.ColumnCount)
	}

	if len(model.Columns) != 1 {
		t.Errorf("Expected 1 column, got %d", len(model.Columns))
	}

	// Model now starts with default MCPs to facilitate initial testing and user onboarding
	if len(model.MCPItems) != 3 {
		t.Errorf("Expected MCPItems to have 3 default items, got %d items", len(model.MCPItems))
	}

	if model.FormErrors == nil {
		t.Error("Expected FormErrors to be initialized")
	}
}

func TestNewModelWithMCPs(t *testing.T) {
	testMCPs := []MCPItem{
		{Name: "test-mcp", Type: "CMD", Active: true, Command: "test-command"},
		{Name: "another-mcp", Type: "SSE", Active: false, URL: "http://example.com"},
	}

	model := NewModelWithMCPs(testMCPs)

	// Test that provided MCPs are used
	if len(model.MCPItems) != 2 {
		t.Errorf("Expected 2 MCPItems, got %d", len(model.MCPItems))
	}

	if !reflect.DeepEqual(model.MCPItems, testMCPs) {
		t.Error("MCPItems do not match provided items")
	}

	// Test other defaults are still set
	if model.State != MainNavigation {
		t.Errorf("Expected State to be MainNavigation, got %v", model.State)
	}

	if model.FormErrors == nil {
		t.Error("Expected FormErrors to be initialized")
	}
}

func TestGetDefaultMCPs(t *testing.T) {
	defaults := getDefaultMCPs()

	// Default MCPs now include common examples to facilitate user onboarding
	if len(defaults) != 3 {
		t.Errorf("Expected default MCPs to have 3 items for user onboarding, got %d items", len(defaults))
	}

	// Verify it returns a valid slice (not nil)
	if defaults == nil {
		t.Error("Expected getDefaultMCPs to return a valid slice, got nil")
	}
}

func TestMCPItemTypes(t *testing.T) {
	// Test different MCP types
	cmdMCP := MCPItem{
		Name:    "test-cmd",
		Type:    "CMD",
		Active:  true,
		Command: "test-command",
		Args:    []string{"arg1", "arg2"},
		Environment: map[string]string{
			"TEST_VAR": "test_value",
		},
	}

	sseMCP := MCPItem{
		Name:   "test-sse",
		Type:   "SSE",
		Active: false,
		URL:    "http://example.com",
		Environment: map[string]string{
			"API_KEY": "secret",
		},
	}

	jsonMCP := MCPItem{
		Name:       "test-json",
		Type:       "JSON",
		Active:     false,
		JSONConfig: `{"key": "value"}`,
	}

	// Test that all fields are accessible
	if cmdMCP.Name != "test-cmd" {
		t.Error("CMD MCP name not set correctly")
	}
	if cmdMCP.Type != "CMD" {
		t.Error("CMD MCP type not set correctly")
	}
	if !cmdMCP.Active {
		t.Error("CMD MCP active not set correctly")
	}
	if cmdMCP.Command != "test-command" {
		t.Error("CMD MCP command not set correctly")
	}
	if len(cmdMCP.Args) != 2 {
		t.Error("CMD MCP args not set correctly")
	}
	if cmdMCP.Environment["TEST_VAR"] != "test_value" {
		t.Error("CMD MCP environment not set correctly")
	}

	if sseMCP.Name != "test-sse" {
		t.Error("SSE MCP name not set correctly")
	}
	if sseMCP.Type != "SSE" {
		t.Error("SSE MCP type not set correctly")
	}
	if sseMCP.Active {
		t.Error("SSE MCP active not set correctly")
	}
	if sseMCP.URL != "http://example.com" {
		t.Error("SSE MCP URL not set correctly")
	}
	if sseMCP.Environment["API_KEY"] != "secret" {
		t.Error("SSE MCP environment not set correctly")
	}

	if jsonMCP.Name != "test-json" {
		t.Error("JSON MCP name not set correctly")
	}
	if jsonMCP.Type != "JSON" {
		t.Error("JSON MCP type not set correctly")
	}
	if jsonMCP.Active {
		t.Error("JSON MCP active not set correctly")
	}
	if jsonMCP.JSONConfig != `{"key": "value"}` {
		t.Error("JSON MCP config not set correctly")
	}
}

func TestAppStateConstants(t *testing.T) {
	// Test that app state constants are properly defined
	if MainNavigation != 0 {
		t.Errorf("Expected MainNavigation to be 0, got %d", MainNavigation)
	}
	if SearchMode != 1 {
		t.Errorf("Expected SearchMode to be 1, got %d", SearchMode)
	}
	if SearchActiveNavigation != 2 {
		t.Errorf("Expected SearchActiveNavigation to be 2, got %d", SearchActiveNavigation)
	}
	if ModalActive != 3 {
		t.Errorf("Expected ModalActive to be 3, got %d", ModalActive)
	}
}

func TestModalTypeConstants(t *testing.T) {
	// Test that modal type constants are properly defined
	if NoModal != 0 {
		t.Errorf("Expected NoModal to be 0, got %d", NoModal)
	}
	if AddModal != 1 {
		t.Errorf("Expected AddModal to be 1, got %d", AddModal)
	}
	if AddMCPTypeSelection != 2 {
		t.Errorf("Expected AddMCPTypeSelection to be 2, got %d", AddMCPTypeSelection)
	}
	if AddCommandForm != 3 {
		t.Errorf("Expected AddCommandForm to be 3, got %d", AddCommandForm)
	}
	if AddSSEForm != 4 {
		t.Errorf("Expected AddSSEForm to be 4, got %d", AddSSEForm)
	}
	if AddJSONForm != 5 {
		t.Errorf("Expected AddJSONForm to be 5, got %d", AddJSONForm)
	}
	if EditModal != 6 {
		t.Errorf("Expected EditModal to be 6, got %d", EditModal)
	}
	if DeleteModal != 7 {
		t.Errorf("Expected DeleteModal to be 7, got %d", DeleteModal)
	}
}

func TestFormData(t *testing.T) {
	formData := FormData{
		Name:        "test-name",
		Command:     "test-command",
		Args:        "arg1 arg2",
		URL:         "http://example.com",
		JSONConfig:  `{"test": true}`,
		Environment: "TEST_VAR=value",
		ActiveField: 1,
	}

	// Test that all fields are accessible
	if formData.Name != "test-name" {
		t.Error("FormData Name not set correctly")
	}
	if formData.Command != "test-command" {
		t.Error("FormData Command not set correctly")
	}
	if formData.Args != "arg1 arg2" {
		t.Error("FormData Args not set correctly")
	}
	if formData.URL != "http://example.com" {
		t.Error("FormData URL not set correctly")
	}
	if formData.JSONConfig != `{"test": true}` {
		t.Error("FormData JSONConfig not set correctly")
	}
	if formData.Environment != "TEST_VAR=value" {
		t.Error("FormData Environment not set correctly")
	}
	if formData.ActiveField != 1 {
		t.Error("FormData ActiveField not set correctly")
	}
}

func TestColumn(t *testing.T) {
	column := Column{
		Title: "Test Column",
		Items: []string{"item1", "item2", "item3"},
		Width: 25,
	}

	if column.Title != "Test Column" {
		t.Error("Column Title not set correctly")
	}
	if len(column.Items) != 3 {
		t.Error("Column Items not set correctly")
	}
	if column.Width != 25 {
		t.Error("Column Width not set correctly")
	}
}

func TestClaudeStatus(t *testing.T) {
	now := time.Now()
	status := ClaudeStatus{
		Available:    true,
		Version:      "1.0.0",
		ActiveMCPs:   []string{"mcp1", "mcp2"},
		LastCheck:    now,
		Error:        "",
		InstallGuide: "Install guide text",
	}

	if !status.Available {
		t.Error("ClaudeStatus Available not set correctly")
	}
	if status.Version != "1.0.0" {
		t.Error("ClaudeStatus Version not set correctly")
	}
	if len(status.ActiveMCPs) != 2 {
		t.Error("ClaudeStatus ActiveMCPs not set correctly")
	}
	if !status.LastCheck.Equal(now) {
		t.Error("ClaudeStatus LastCheck not set correctly")
	}
	if status.InstallGuide != "Install guide text" {
		t.Error("ClaudeStatus InstallGuide not set correctly")
	}
	if status.Error != "" {
		t.Error("ClaudeStatus Error not set correctly")
	}
}

func TestModelInit(t *testing.T) {
	model := NewModel()
	cmd := model.Init()

	// Init should return nil command for types.Model
	if cmd != nil {
		t.Error("Expected Init() to return nil")
	}
}

func TestLayoutConstants(t *testing.T) {
	// Test layout constants are properly defined
	if COLUMN_WIDTH != 28 {
		t.Errorf("Expected COLUMN_WIDTH to be 28, got %d", COLUMN_WIDTH)
	}
	if WIDE_LAYOUT_MIN != 120 {
		t.Errorf("Expected WIDE_LAYOUT_MIN to be 120, got %d", WIDE_LAYOUT_MIN)
	}
	if MEDIUM_LAYOUT_MIN != 80 {
		t.Errorf("Expected MEDIUM_LAYOUT_MIN to be 80, got %d", MEDIUM_LAYOUT_MIN)
	}
	if WIDE_COLUMNS != 4 {
		t.Errorf("Expected WIDE_COLUMNS to be 4, got %d", WIDE_COLUMNS)
	}
	if MEDIUM_COLUMNS != 2 {
		t.Errorf("Expected MEDIUM_COLUMNS to be 2, got %d", MEDIUM_COLUMNS)
	}
	if NARROW_COLUMNS != 1 {
		t.Errorf("Expected NARROW_COLUMNS to be 1, got %d", NARROW_COLUMNS)
	}
}

func TestModelComplexState(t *testing.T) {
	model := NewModel()

	// Test setting various states
	model.State = SearchActiveNavigation
	model.SearchActive = true
	model.SearchInputActive = true
	model.SearchQuery = "test query"
	model.ActiveModal = AddCommandForm
	model.SuccessMessage = "Success!"
	model.SuccessTimer = 120

	// Test state persistence
	if model.State != SearchActiveNavigation {
		t.Error("State not set correctly")
	}
	if !model.SearchActive {
		t.Error("SearchActive not set correctly")
	}
	if !model.SearchInputActive {
		t.Error("SearchInputActive not set correctly")
	}
	if model.SearchQuery != "test query" {
		t.Error("SearchQuery not set correctly")
	}
	if model.ActiveModal != AddCommandForm {
		t.Error("ActiveModal not set correctly")
	}
	if model.SuccessMessage != "Success!" {
		t.Error("SuccessMessage not set correctly")
	}
	if model.SuccessTimer != 120 {
		t.Error("SuccessTimer not set correctly")
	}
}

func TestMCPItemEnvironmentHandling(t *testing.T) {
	// Test MCP with nil environment
	mcp1 := MCPItem{
		Name:        "test1",
		Type:        "CMD",
		Command:     "test",
		Environment: nil,
	}

	if mcp1.Name != "test1" {
		t.Error("MCP1 name not set correctly")
	}
	if mcp1.Type != "CMD" {
		t.Error("MCP1 type not set correctly")
	}
	if mcp1.Command != "test" {
		t.Error("MCP1 command not set correctly")
	}
	if mcp1.Environment != nil {
		t.Error("Expected nil environment to remain nil")
	}

	// Test MCP with empty environment
	mcp2 := MCPItem{
		Name:        "test2",
		Type:        "CMD",
		Command:     "test",
		Environment: make(map[string]string),
	}

	if mcp2.Name != "test2" {
		t.Error("MCP2 name not set correctly")
	}
	if mcp2.Type != "CMD" {
		t.Error("MCP2 type not set correctly")
	}
	if mcp2.Command != "test" {
		t.Error("MCP2 command not set correctly")
	}
	if mcp2.Environment == nil {
		t.Error("Expected empty environment map to be preserved")
	}
	if len(mcp2.Environment) != 0 {
		t.Error("Expected empty environment map to have length 0")
	}

	// Test MCP with populated environment
	mcp3 := MCPItem{
		Name:    "test3",
		Type:    "CMD",
		Command: "test",
		Environment: map[string]string{
			"VAR1": "value1",
			"VAR2": "value2",
		},
	}

	if mcp3.Name != "test3" {
		t.Error("MCP3 name not set correctly")
	}
	if mcp3.Type != "CMD" {
		t.Error("MCP3 type not set correctly")
	}
	if mcp3.Command != "test" {
		t.Error("MCP3 command not set correctly")
	}
	if len(mcp3.Environment) != 2 {
		t.Error("Expected environment to have 2 variables")
	}
	if mcp3.Environment["VAR1"] != "value1" {
		t.Error("Environment variable VAR1 not set correctly")
	}
	if mcp3.Environment["VAR2"] != "value2" {
		t.Error("Environment variable VAR2 not set correctly")
	}
}

func TestLoadingOverlayMethods(t *testing.T) {
	t.Run("StartLoadingOverlay", func(t *testing.T) {
		model := NewModel()
		
		model.StartLoadingOverlay(LoadingStartup)
		
		if !model.IsLoadingOverlayActive() {
			t.Error("Loading overlay should be active after start")
		}
		if model.LoadingOverlay.Message == "" {
			t.Error("Loading message should be set after start")
		}
		if model.LoadingOverlay.Type != LoadingStartup {
			t.Errorf("Expected type LoadingStartup, got %v", model.LoadingOverlay.Type)
		}
	})
	
	t.Run("UpdateLoadingMessage", func(t *testing.T) {
		model := NewModel()
		model.StartLoadingOverlay(LoadingRefresh)
		oldMessage := model.LoadingOverlay.Message
		
		model.UpdateLoadingMessage("New message")
		
		if model.LoadingOverlay.Message != "New message" {
			t.Errorf("Expected message 'New message', got '%s'", model.LoadingOverlay.Message)
		}
		if !model.IsLoadingOverlayActive() {
			t.Error("Loading overlay should remain active")
		}
		
		// Test that it changed from old message
		if model.LoadingOverlay.Message == oldMessage {
			t.Error("Message should have changed")
		}
	})
	
	t.Run("StopLoadingOverlay", func(t *testing.T) {
		model := NewModel()
		model.StartLoadingOverlay(LoadingStartup)
		
		// Ensure it's active first
		if !model.IsLoadingOverlayActive() {
			t.Error("Setup: Loading overlay should be active before stop")
		}
		
		model.StopLoadingOverlay()
		
		if model.IsLoadingOverlayActive() {
			t.Error("Loading overlay should be inactive after stop")
		}
		if model.LoadingOverlay != nil {
			t.Error("Loading overlay should be nil after stop")
		}
	})
	
	t.Run("AdvanceSpinner", func(t *testing.T) {
		model := NewModel()
		model.StartLoadingOverlay(LoadingStartup)
		
		initialSpinner := model.LoadingOverlay.Spinner
		
		model.AdvanceSpinner()
		
		if model.LoadingOverlay.Spinner == initialSpinner {
			t.Error("Spinner should have advanced")
		}
		
		// Test multiple advances to check wrap around
		for i := 0; i < 10; i++ {
			model.AdvanceSpinner()
		}
		
		// Should still be a valid spinner value (0-3)
		if model.LoadingOverlay.Spinner < 0 || model.LoadingOverlay.Spinner > 3 {
			t.Errorf("Spinner should be between 0-3, got %d", model.LoadingOverlay.Spinner)
		}
	})
	
	t.Run("IsLoadingOverlayActive", func(t *testing.T) {
		model := NewModel()
		
		// Initially inactive
		if model.IsLoadingOverlayActive() {
			t.Error("Loading overlay should be inactive initially")
		}
		
		// Activate
		model.StartLoadingOverlay(LoadingStartup)
		if !model.IsLoadingOverlayActive() {
			t.Error("Loading overlay should be active after start")
		}
		
		// Deactivate
		model.StopLoadingOverlay()
		if model.IsLoadingOverlayActive() {
			t.Error("Loading overlay should be inactive after stop")
		}
	})
}

func TestSpinnerTypes(t *testing.T) {
	// Test that spinner constants are defined
	if SpinnerFrame1 != 0 {
		t.Errorf("Expected SpinnerFrame1 to be 0, got %d", SpinnerFrame1)
	}
	if SpinnerFrame2 != 1 {
		t.Errorf("Expected SpinnerFrame2 to be 1, got %d", SpinnerFrame2)
	}
	if SpinnerFrame3 != 2 {
		t.Errorf("Expected SpinnerFrame3 to be 2, got %d", SpinnerFrame3)
	}
	if SpinnerFrame4 != 3 {
		t.Errorf("Expected SpinnerFrame4 to be 3, got %d", SpinnerFrame4)
	}
}

func TestLoadingMessageGeneration(t *testing.T) {
	// Test loading startup message
	model := NewModel()
	model.StartLoadingOverlay(LoadingStartup)
	
	if model.LoadingOverlay.Message == "" {
		t.Error("Startup loading message should not be empty")
	}
	
	// Test loading refresh message
	model2 := NewModel()
	model2.StartLoadingOverlay(LoadingRefresh)
	
	if model2.LoadingOverlay.Message == "" {
		t.Error("Refresh loading message should not be empty")
	}
	
	// Messages should be different for different types
	if model.LoadingOverlay.Message == model2.LoadingOverlay.Message {
		t.Error("Different loading types should have different messages")
	}
}

func TestMessageTypes(t *testing.T) {
	t.Run("TimerTickMsg", func(t *testing.T) {
		msg := TimerTickMsg{ID: "timer-42"}
		if msg.ID != "timer-42" {
			t.Errorf("Expected ID 'timer-42', got '%s'", msg.ID)
		}
	})
	
	t.Run("LoadingProgressMsg", func(t *testing.T) {
		msg := LoadingProgressMsg{
			Type:    LoadingRefresh,
			Message: "Progress message",
			Done:    false,
		}
		if msg.Type != LoadingRefresh {
			t.Errorf("Expected LoadingRefresh, got %v", msg.Type)
		}
		if msg.Message != "Progress message" {
			t.Errorf("Expected 'Progress message', got '%s'", msg.Message)
		}
		if msg.Done {
			t.Error("Expected Done to be false")
		}
	})
	
	t.Run("LoadingStepMsg", func(t *testing.T) {
		msg := LoadingStepMsg{
			Type: LoadingStartup,
			Step: 3,
		}
		if msg.Step != 3 {
			t.Errorf("Expected step 3, got %d", msg.Step)
		}
		if msg.Type != LoadingStartup {
			t.Errorf("Expected LoadingStartup, got %v", msg.Type)
		}
	})
	
	t.Run("LoadingSpinnerMsg", func(t *testing.T) {
		msg := LoadingSpinnerMsg{Type: LoadingStartup}
		if msg.Type != LoadingStartup {
			t.Errorf("Expected LoadingStartup, got %v", msg.Type)
		}
	})
	
	t.Run("ProjectContextCheckMsg", func(t *testing.T) {
		msg := ProjectContextCheckMsg{}
		// Empty struct, just test that it compiles and exists
		_ = msg
	})
	
	t.Run("DirectoryChangeMsg", func(t *testing.T) {
		msg := DirectoryChangeMsg{NewPath: "/new/path"}
		if msg.NewPath != "/new/path" {
			t.Errorf("Expected '/new/path', got '%s'", msg.NewPath)
		}
	})
}

func TestLoadingTypeConstants(t *testing.T) {
	// Test loading type constants
	if LoadingStartup != 0 {
		t.Errorf("Expected LoadingStartup to be 0, got %d", LoadingStartup)
	}
	if LoadingRefresh != 1 {
		t.Errorf("Expected LoadingRefresh to be 1, got %d", LoadingRefresh)
	}
}

func TestProjectContextStruct(t *testing.T) {
	ctx := ProjectContext{
		CurrentPath:    "/test/project",
		ActiveMCPs:     5,
		TotalMCPs:      10,
		DisplayPath:    "/test/project",
		SyncStatusText: "Synced",
	}
	
	if ctx.CurrentPath != "/test/project" {
		t.Error("ProjectContext CurrentPath not set correctly")
	}
	if ctx.ActiveMCPs != 5 {
		t.Error("ProjectContext ActiveMCPs not set correctly")
	}
	if ctx.TotalMCPs != 10 {
		t.Error("ProjectContext TotalMCPs not set correctly")
	}
	if ctx.DisplayPath != "/test/project" {
		t.Error("ProjectContext DisplayPath not set correctly")
	}
	if ctx.SyncStatusText != "Synced" {
		t.Error("ProjectContext SyncStatusText not set correctly")
	}
}

func TestLoadingOverlayStruct(t *testing.T) {
	overlay := LoadingOverlay{
		Active:      true,
		Message:     "Loading...",
		Spinner:     SpinnerFrame2,
		Cancellable: true,
		Type:        LoadingStartup,
	}
	
	if !overlay.Active {
		t.Error("LoadingOverlay Active not set correctly")
	}
	if overlay.Message != "Loading..." {
		t.Error("LoadingOverlay Message not set correctly")
	}
	if overlay.Spinner != SpinnerFrame2 {
		t.Error("LoadingOverlay Spinner not set correctly")
	}
	if !overlay.Cancellable {
		t.Error("LoadingOverlay Cancellable not set correctly")
	}
	if overlay.Type != LoadingStartup {
		t.Error("LoadingOverlay Type not set correctly")
	}
}

func TestSpinnerCharGeneration(t *testing.T) {
	tests := []struct {
		state    SpinnerState
		expected string
	}{
		{SpinnerFrame1, "◐"},
		{SpinnerFrame2, "◓"},
		{SpinnerFrame3, "◑"},
		{SpinnerFrame4, "◒"},
	}
	
	for _, tt := range tests {
		t.Run("Spinner state", func(t *testing.T) {
			result := tt.state.GetSpinnerChar()
			if result != tt.expected {
				t.Errorf("Expected spinner char %s, got %s", tt.expected, result)
			}
		})
	}
}
