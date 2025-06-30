package services

import (
	"testing"

	"cc-mcp-manager/internal/ui/types"
)

func TestUpdateLayout(t *testing.T) {
	tests := []struct {
		name                string
		width               int
		expectedColumns     int
		expectedColumnCount int
		expectedTitle       string
	}{
		{
			name:                "Wide layout - 4 columns",
			width:               150,
			expectedColumns:     4,
			expectedColumnCount: types.WIDE_COLUMNS,
			expectedTitle:       "MCPs Column 1",
		},
		{
			name:                "Wide layout boundary - exactly 120",
			width:               120,
			expectedColumns:     4,
			expectedColumnCount: types.WIDE_COLUMNS,
			expectedTitle:       "MCPs Column 1",
		},
		{
			name:                "Medium layout - 2 columns",
			width:               100,
			expectedColumns:     2,
			expectedColumnCount: types.MEDIUM_COLUMNS,
			expectedTitle:       "MCPs",
		},
		{
			name:                "Medium layout boundary - exactly 80",
			width:               80,
			expectedColumns:     2,
			expectedColumnCount: types.MEDIUM_COLUMNS,
			expectedTitle:       "MCPs",
		},
		{
			name:                "Narrow layout - 1 column",
			width:               70,
			expectedColumns:     1,
			expectedColumnCount: types.NARROW_COLUMNS,
			expectedTitle:       "MCP Manager",
		},
		{
			name:                "Very narrow layout",
			width:               40,
			expectedColumns:     1,
			expectedColumnCount: types.NARROW_COLUMNS,
			expectedTitle:       "MCP Manager",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := types.Model{
				Width:        tt.width,
				ActiveColumn: 0,
			}

			result := UpdateLayout(model)

			// Test column count
			if result.ColumnCount != tt.expectedColumnCount {
				t.Errorf("UpdateLayout() ColumnCount = %d, expected %d",
					result.ColumnCount, tt.expectedColumnCount)
			}

			// Test actual columns array length
			if len(result.Columns) != tt.expectedColumns {
				t.Errorf("UpdateLayout() len(Columns) = %d, expected %d",
					len(result.Columns), tt.expectedColumns)
			}

			// Test first column title
			if len(result.Columns) > 0 && result.Columns[0].Title != tt.expectedTitle {
				t.Errorf("UpdateLayout() first column title = %s, expected %s",
					result.Columns[0].Title, tt.expectedTitle)
			}

			// Test that ActiveColumn is within bounds
			if result.ActiveColumn >= result.ColumnCount {
				t.Errorf("UpdateLayout() ActiveColumn %d should be < ColumnCount %d",
					result.ActiveColumn, result.ColumnCount)
			}
		})
	}
}

func TestUpdateLayoutColumnWidths(t *testing.T) {
	tests := []struct {
		name          string
		width         int
		expectedWidth int
		layout        string
	}{
		{
			name:          "Wide layout column width calculation",
			width:         140,
			expectedWidth: (140 - 10) / 4, // 32
			layout:        "wide",
		},
		{
			name:          "Medium layout column width calculation",
			width:         100,
			expectedWidth: (100 - 6) / 2, // 47
			layout:        "medium",
		},
		{
			name:          "Narrow layout column width calculation",
			width:         60,
			expectedWidth: 60 - 4, // 56
			layout:        "narrow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := types.Model{
				Width:        tt.width,
				ActiveColumn: 0,
			}

			result := UpdateLayout(model)

			if len(result.Columns) > 0 {
				actualWidth := result.Columns[0].Width
				if actualWidth != tt.expectedWidth {
					t.Errorf("UpdateLayout() column width = %d, expected %d",
						actualWidth, tt.expectedWidth)
				}
			}
		})
	}
}

func TestUpdateLayoutActiveColumnReset(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		initialColumn  int
		expectedColumn int
		shouldReset    bool
	}{
		{
			name:           "ActiveColumn within bounds - no reset",
			width:          150, // 4 columns
			initialColumn:  2,   // valid for 4 columns
			expectedColumn: 2,
			shouldReset:    false,
		},
		{
			name:           "ActiveColumn out of bounds - reset to last",
			width:          100, // 2 columns
			initialColumn:  3,   // invalid for 2 columns
			expectedColumn: 1,   // reset to last valid (2-1)
			shouldReset:    true,
		},
		{
			name:           "ActiveColumn way out of bounds",
			width:          60, // 1 column
			initialColumn:  5,  // way out of bounds
			expectedColumn: 0,  // reset to 0
			shouldReset:    true,
		},
		{
			name:           "ActiveColumn at exact boundary",
			width:          150, // 4 columns
			initialColumn:  3,   // exactly at last valid index
			expectedColumn: 3,
			shouldReset:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := types.Model{
				Width:        tt.width,
				ActiveColumn: tt.initialColumn,
			}

			result := UpdateLayout(model)

			if result.ActiveColumn != tt.expectedColumn {
				t.Errorf("UpdateLayout() ActiveColumn = %d, expected %d",
					result.ActiveColumn, tt.expectedColumn)
			}
		})
	}
}

func TestUpdateLayoutBreakpoints(t *testing.T) {
	// Test exact boundary conditions
	boundaryTests := []struct {
		name     string
		width    int
		expected int
	}{
		{"Just below wide threshold", types.WIDE_LAYOUT_MIN - 1, types.MEDIUM_COLUMNS},
		{"Exactly at wide threshold", types.WIDE_LAYOUT_MIN, types.WIDE_COLUMNS},
		{"Just above wide threshold", types.WIDE_LAYOUT_MIN + 1, types.WIDE_COLUMNS},
		{"Just below medium threshold", types.MEDIUM_LAYOUT_MIN - 1, types.NARROW_COLUMNS},
		{"Exactly at medium threshold", types.MEDIUM_LAYOUT_MIN, types.MEDIUM_COLUMNS},
		{"Just above medium threshold", types.MEDIUM_LAYOUT_MIN + 1, types.MEDIUM_COLUMNS},
	}

	for _, tt := range boundaryTests {
		t.Run(tt.name, func(t *testing.T) {
			model := types.Model{Width: tt.width}
			result := UpdateLayout(model)

			if result.ColumnCount != tt.expected {
				t.Errorf("UpdateLayout() at width %d: ColumnCount = %d, expected %d",
					tt.width, result.ColumnCount, tt.expected)
			}
		})
	}
}

func TestUpdateLayoutColumnTitles(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		expectedTitles []string
	}{
		{
			name:  "Wide layout column titles",
			width: 150,
			expectedTitles: []string{
				"MCPs Column 1",
				"MCPs Column 2",
				"MCPs Column 3",
				"MCPs Column 4",
			},
		},
		{
			name:  "Medium layout column titles",
			width: 100,
			expectedTitles: []string{
				"MCPs",
				"Status & Details",
			},
		},
		{
			name:  "Narrow layout column titles",
			width: 60,
			expectedTitles: []string{
				"MCP Manager",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := types.Model{Width: tt.width}
			result := UpdateLayout(model)

			if len(result.Columns) != len(tt.expectedTitles) {
				t.Errorf("UpdateLayout() columns length = %d, expected %d",
					len(result.Columns), len(tt.expectedTitles))
				return
			}

			for i, expectedTitle := range tt.expectedTitles {
				if result.Columns[i].Title != expectedTitle {
					t.Errorf("UpdateLayout() column %d title = %s, expected %s",
						i, result.Columns[i].Title, expectedTitle)
				}
			}
		})
	}
}

// Edge case and error condition tests
func TestUpdateLayoutEdgeCases(t *testing.T) {
	t.Run("Zero width", func(t *testing.T) {
		model := types.Model{Width: 0}
		result := UpdateLayout(model)

		// Should default to narrow layout
		if result.ColumnCount != types.NARROW_COLUMNS {
			t.Errorf("UpdateLayout() with zero width should use narrow layout")
		}
	})

	t.Run("Negative width", func(t *testing.T) {
		model := types.Model{Width: -10}
		result := UpdateLayout(model)

		// Should default to narrow layout
		if result.ColumnCount != types.NARROW_COLUMNS {
			t.Errorf("UpdateLayout() with negative width should use narrow layout")
		}
	})

	t.Run("Very large width", func(t *testing.T) {
		model := types.Model{Width: 10000}
		result := UpdateLayout(model)

		// Should use wide layout
		if result.ColumnCount != types.WIDE_COLUMNS {
			t.Errorf("UpdateLayout() with very large width should use wide layout")
		}

		// Column widths should be reasonable
		expectedWidth := (10000 - 10) / 4
		if len(result.Columns) > 0 && result.Columns[0].Width != expectedWidth {
			t.Errorf("UpdateLayout() column width calculation incorrect for large width")
		}
	})
}

// Test layout consistency
func TestUpdateLayoutConsistency(t *testing.T) {
	t.Run("Multiple updates preserve state", func(t *testing.T) {
		model := types.Model{
			Width:        150,
			ActiveColumn: 2,
		}

		// Apply layout update twice
		result1 := UpdateLayout(model)
		result2 := UpdateLayout(result1)

		// Results should be identical
		if result1.ColumnCount != result2.ColumnCount {
			t.Errorf("Multiple UpdateLayout() calls should be consistent")
		}

		if result1.ActiveColumn != result2.ActiveColumn {
			t.Errorf("ActiveColumn should be preserved across multiple updates")
		}
	})

	t.Run("Layout transitions", func(t *testing.T) {
		model := types.Model{
			Width:        150, // Start wide
			ActiveColumn: 3,   // Last column in wide layout
		}

		// Update to medium layout
		model.Width = 100
		result := UpdateLayout(model)

		// ActiveColumn should be reset to valid range
		if result.ActiveColumn >= result.ColumnCount {
			t.Errorf("Layout transition should reset ActiveColumn to valid range")
		}

		// Should be medium layout
		if result.ColumnCount != types.MEDIUM_COLUMNS {
			t.Errorf("Layout transition should update to correct layout")
		}
	})
}
