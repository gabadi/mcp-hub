// Package types provides shared data types and constants for the MCP manager UI application.
// This package contains UI-specific types, constants, and message structures.
package types

const (
	// ColumnWidth defines the fixed width for grid columns to maintain alignment
	ColumnWidth = 28

	// WideLayoutMin defines the minimum width for 4-column layout
	WideLayoutMin = 120
	// MediumLayoutMin defines the minimum width for 2-column layout
	MediumLayoutMin = 80

	// WideColumns defines the 4-column grid for wide layout
	WideColumns = 4
	// MediumColumns defines the 2-column layout for medium screens
	MediumColumns = 2
	// NarrowColumns defines the single column for narrow screens
	NarrowColumns = 1

	// BulletChar represents the bullet character used in UI displays
	BulletChar = "◦"

	// KeyCommandC represents the command+c keyboard shortcut
	KeyCommandC = "command+c"
	// KeyCommandV represents the command+v keyboard shortcut
	KeyCommandV = "command+v"

	// TestCMD represents the CMD type identifier for testing
	TestCMD = "CMD"
	// TestStr represents a generic test string constant
	TestStr = "test"
)
