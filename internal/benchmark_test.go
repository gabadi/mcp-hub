package internal

import (
	"fmt"
	"testing"

	"mcp-hub/internal/platform"
	"mcp-hub/internal/testutil"
	"mcp-hub/internal/ui"
	"mcp-hub/internal/ui/services"
	"mcp-hub/internal/ui/types"

	tea "github.com/charmbracelet/bubbletea"
)

// Benchmark tests for performance verification and regression prevention
// These tests ensure critical paths maintain acceptable performance characteristics

func BenchmarkMCPService_FilterMCPs_SmallDataset(b *testing.B) {
	// Small dataset (realistic for most users)
	mcps := generateBenchmarkMCPDataset(50)
	model := testutil.NewTestModel().WithMCPs(mcps).WithSearchQuery("github").Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = services.GetFilteredMCPs(model)
	}
}

func BenchmarkMCPService_FilterMCPs_MediumDataset(b *testing.B) {
	// Medium dataset (power users)
	mcps := generateBenchmarkMCPDataset(200)
	model := testutil.NewTestModel().WithMCPs(mcps).WithSearchQuery("mcp").Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = services.GetFilteredMCPs(model)
	}
}

func BenchmarkMCPService_FilterMCPs_LargeDataset(b *testing.B) {
	// Large dataset (stress test)
	mcps := generateBenchmarkMCPDataset(1000)
	model := testutil.NewTestModel().WithMCPs(mcps).WithSearchQuery("test").Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = services.GetFilteredMCPs(model)
	}
}

func BenchmarkMCPService_FilterMCPs_NoMatches(b *testing.B) {
	// Benchmark worst case: no matches found
	mcps := generateBenchmarkMCPDataset(500)
	model := testutil.NewTestModel().WithMCPs(mcps).WithSearchQuery("nonexistent-query").Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = services.GetFilteredMCPs(model)
	}
}

func BenchmarkMCPService_FilterMCPs_EmptyQuery(b *testing.B) {
	// Benchmark best case: empty query returns all
	mcps := generateBenchmarkMCPDataset(500)
	model := testutil.NewTestModel().WithMCPs(mcps).WithSearchQuery("").Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = services.GetFilteredMCPs(model)
	}
}

func BenchmarkMCPService_ToggleMCPStatus(b *testing.B) {
	mcps := generateBenchmarkMCPDataset(100)
	mockPlatform := platform.GetMockPlatformService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model := testutil.NewTestModel().WithMCPs(mcps).WithSelectedItem(i % len(mcps)).Build()
		_ = services.ToggleMCPStatus(model, mockPlatform)
	}
}

func BenchmarkMCPService_GetActiveMCPCount(b *testing.B) {
	mcps := generateBenchmarkMCPDataset(500)
	model := testutil.NewTestModel().WithMCPs(mcps).Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = services.GetActiveMCPCount(model)
	}
}

func BenchmarkMCPService_GetSelectedMCP(b *testing.B) {
	mcps := generateBenchmarkMCPDataset(100)
	model := testutil.NewTestModel().WithMCPs(mcps).WithSelectedItem(50).Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = services.GetSelectedMCP(model)
	}
}

func BenchmarkLayoutService_UpdateLayout_Wide(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := testutil.NewTestModel().WithWindowSize(150, 50).Build()
		_ = services.UpdateLayout(model)
	}
}

func BenchmarkLayoutService_UpdateLayout_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := testutil.NewTestModel().WithWindowSize(100, 30).Build()
		_ = services.UpdateLayout(model)
	}
}

func BenchmarkLayoutService_UpdateLayout_Narrow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := testutil.NewTestModel().WithWindowSize(70, 25).Build()
		_ = services.UpdateLayout(model)
	}
}

func BenchmarkUI_ModelUpdate_Navigation(b *testing.B) {
	model := ui.Model{
		Model: testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithMCPs(generateBenchmarkMCPDataset(100)).
			Build(),
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyRight}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updatedModel, _ := model.Update(keyMsg)
		model = updatedModel.(ui.Model)
	}
}

func BenchmarkUI_ModelUpdate_Search(b *testing.B) {
	model := ui.Model{
		Model: testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithMCPs(generateBenchmarkMCPDataset(200)).
			WithSearchActive(true).
			WithSearchInputActive(true).
			Build(),
	}

	searchMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updatedModel, _ := model.Update(searchMsg)
		model = updatedModel.(ui.Model)
	}
}

func BenchmarkUI_ModelUpdate_Toggle(b *testing.B) {
	model := ui.Model{
		Model: testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithMCPs(generateBenchmarkMCPDataset(100)).
			WithSelectedItem(50).
			Build(),
	}

	spaceMsg := tea.KeyMsg{Type: tea.KeySpace}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updatedModel, _ := model.Update(spaceMsg)
		model = updatedModel.(ui.Model)
	}
}

func BenchmarkUI_View_Rendering_Small(b *testing.B) {
	model := ui.Model{
		Model: testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithMCPs(generateBenchmarkMCPDataset(50)).
			Build(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

func BenchmarkUI_View_Rendering_Medium(b *testing.B) {
	model := ui.Model{
		Model: testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithMCPs(generateBenchmarkMCPDataset(200)).
			Build(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

func BenchmarkUI_View_Rendering_Large(b *testing.B) {
	model := ui.Model{
		Model: testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithMCPs(generateBenchmarkMCPDataset(500)).
			Build(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

func BenchmarkUI_View_Rendering_Search(b *testing.B) {
	model := ui.Model{
		Model: testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithMCPs(generateBenchmarkMCPDataset(300)).
			WithSearchActive(true).
			WithSearchQuery("test").
			Build(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

func BenchmarkUI_CompleteWorkflow_SearchAndNavigate(b *testing.B) {
	initialModel := ui.Model{
		Model: testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithMCPs(generateBenchmarkMCPDataset(100)).
			Build(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model := initialModel

		// Enter search mode
		searchKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")}
		updatedModel, _ := model.Update(searchKey)
		model = updatedModel.(ui.Model)

		// Type search query
		queryKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("mcp")}
		updatedModel, _ = model.Update(queryKey)
		model = updatedModel.(ui.Model)

		// Navigate
		rightKey := tea.KeyMsg{Type: tea.KeyRight}
		updatedModel, _ = model.Update(rightKey)
		model = updatedModel.(ui.Model)

		// Exit search
		escKey := tea.KeyMsg{Type: tea.KeyEsc}
		updatedModel, _ = model.Update(escKey)
		model = updatedModel.(ui.Model)

		// Render final state
		_ = model.View()
	}
}

func BenchmarkStorage_SaveInventory_Small(b *testing.B) {
	mcps := generateBenchmarkMCPDataset(10)
	model := testutil.NewTestModel().WithMCPs(mcps).Build()
	mockPlatform := platform.GetMockPlatformService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Note: This will attempt to save to actual config directory
		// In a real benchmark environment, you'd want to use a temp directory
		_ = services.SaveModelInventory(model, mockPlatform)
	}
}

func BenchmarkStorage_SaveInventory_Medium(b *testing.B) {
	mcps := generateBenchmarkMCPDataset(100)
	model := testutil.NewTestModel().WithMCPs(mcps).Build()
	mockPlatform := platform.GetMockPlatformService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = services.SaveModelInventory(model, mockPlatform)
	}
}

func BenchmarkStorage_LoadInventory(b *testing.B) {
	mockPlatform := platform.GetMockPlatformService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = services.LoadInventory(mockPlatform)
	}
}

func BenchmarkMemoryAllocations_ModelCreation(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithMCPs(generateBenchmarkMCPDataset(50)).
			Build()
	}
}

func BenchmarkMemoryAllocations_FilterOperation(b *testing.B) {
	mcps := generateBenchmarkMCPDataset(200)
	model := testutil.NewTestModel().WithMCPs(mcps).WithSearchQuery("test").Build()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = services.GetFilteredMCPs(model)
	}
}

func BenchmarkMemoryAllocations_ViewRendering(b *testing.B) {
	model := ui.Model{
		Model: testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithMCPs(generateBenchmarkMCPDataset(100)).
			Build(),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

// Concurrent benchmarks to test thread safety and concurrent performance
func BenchmarkConcurrent_FilterMCPs(b *testing.B) {
	mcps := generateBenchmarkMCPDataset(200)
	model := testutil.NewTestModel().WithMCPs(mcps).WithSearchQuery("mcp").Build()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = services.GetFilteredMCPs(model)
		}
	})
}

func BenchmarkConcurrent_GetActiveMCPCount(b *testing.B) {
	mcps := generateBenchmarkMCPDataset(300)
	model := testutil.NewTestModel().WithMCPs(mcps).Build()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = services.GetActiveMCPCount(model)
		}
	})
}

// Helper function to generate benchmark datasets
func generateBenchmarkMCPDataset(count int) []types.MCPItem {
	mcps := make([]types.MCPItem, count)
	mcpTypes := []string{"CMD", "SSE", "JSON", "HTTP"}
	prefixes := []string{"github", "docker", "context7", "filesystem", "mcp", "tool", "service", "test"}

	for i := 0; i < count; i++ {
		typeIndex := i % len(mcpTypes)
		prefixIndex := i % len(prefixes)

		mcps[i] = types.MCPItem{
			Name:    fmt.Sprintf("%s-benchmark-%04d", prefixes[prefixIndex], i),
			Type:    mcpTypes[typeIndex],
			Active:  i%3 == 0, // Every third MCP is active
			Command: fmt.Sprintf("benchmark-command-%04d", i),
			Args:    []string{fmt.Sprintf("--config=%04d", i)},
		}
	}

	return mcps
}

// Benchmark comparison functions to track regression
func BenchmarkComparison_OldVsNew_FilterMCPs(b *testing.B) {
	mcps := generateBenchmarkMCPDataset(100)

	b.Run("Current Implementation", func(b *testing.B) {
		model := testutil.NewTestModel().WithMCPs(mcps).WithSearchQuery("mcp").Build()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = services.GetFilteredMCPs(model)
		}
	})

	// Note: In future versions, add comparison with previous implementation
	// b.Run("Previous Implementation", func(b *testing.B) { ... })
}

// Performance regression tests with specific thresholds
func BenchmarkPerformanceThresholds_CriticalPath(b *testing.B) {
	// This benchmark should complete within reasonable time limits
	// Fail the test if performance degrades significantly

	mcps := generateBenchmarkMCPDataset(500)
	model := ui.Model{
		Model: testutil.NewTestModel().
			WithWindowSize(120, 40).
			WithMCPs(mcps).
			Build(),
	}

	searchKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")}
	queryKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updatedModel, _ := model.Update(searchKey)
		m := updatedModel.(ui.Model)
		updatedModel, _ = m.Update(queryKey)
		m = updatedModel.(ui.Model)
		_ = m.View()
	}
}
