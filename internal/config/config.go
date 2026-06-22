package config

// Context overrides or user preferences can be defined here for V2
type Config struct {
	DefaultBranch string
}

// LayoutConfig manages the spacing and dimension variables for the TUI engine
type LayoutConfig struct {
	MenuSpacing int // 👈 Number of newlines between list choices (e.g., 2 for double space)
	ColumnWidth int // Width of the side-by-side layout columns
}

// DefaultLayout yields the baseline look and feel settings
func DefaultLayout() LayoutConfig {
	return LayoutConfig{
		MenuSpacing: 4,  // Set this to 1 for tight lines, 2 for relaxed double-spacing
		ColumnWidth: 50, // Standard responsive terminal column width split
	}
}
