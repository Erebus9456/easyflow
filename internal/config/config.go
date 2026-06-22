package config

// Config aggregates all runtime system preferences, safety rules, and layout specs
type Config struct {
	DefaultBranch string
	Layout        LayoutConfig
	Workflow      WorkflowConfig
}

// LayoutConfig manages the spacing and dimension variables for the TUI engine
type LayoutConfig struct {
	MenuSpacing int // Number of newlines between list choices (e.g., 2 for double space)
	ColumnWidth int // Width of the side-by-side layout columns
}

// WorkflowConfig handles behavioral safety parameters and automation defaults
type WorkflowConfig struct {
	MergePolicy  string // Options: "Standard Merge", "Squash Merge", "Rebase"
	SafetyShield bool   // If true, intercepts dirty workspaces with strict safety prompts
}

// NewDefaultConfig yields the baseline unified configuration suite
func NewDefaultConfig() Config {
	return Config{
		DefaultBranch: "main",
		Layout: LayoutConfig{
			MenuSpacing: 2,  // Balanced spacing
			ColumnWidth: 50, // Standard split view column width
		},
		Workflow: WorkflowConfig{
			MergePolicy:  "Standard Merge",
			SafetyShield: true, // Safety checks active by default
		},
	}
}
