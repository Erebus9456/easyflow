package workflow

// State represents an explicit milestone marker in the EasyFlow dev loop
type State int

const (
	StateDashboard State = iota
	StateSelectIssue
	StateCreateIssue
	StateCreateBranch
	StateWorking
	StateCommitReady
	StatePushing
	StatePRPending
	StateMerging
	StateCompleted

	// CRUD Sub-Category Submenu View States
	StateManageIssues
	StateManageBranches
	StateManageCommits

	// Explicit Action States
	StateListBranches // Used for branch selection/checkout & deletion lists
	StateViewCommits  // Used for displaying git log history entries

	// 🆕 ROADMAP STATES
	StateSettingsMenu          // Interactive Configuration Dashboard
	StateUnsavedChangesWarning // Safety Shield confirmation intercept view
	StateConflictResolution    // Graceful error state panel for merge blocks
	StateManageStash           // Workspace stash shelf management panel
)

// RuntimeContext retains persistent cross-state memory metrics
type RuntimeContext struct {
	ActiveIssueNumber int
	ActiveIssueTitle  string
	BranchName        string
	PullRequestURL    string
	CurrentStep       State
	PipelineMode      bool // Tracks if the user is in the continuous end-to-end loop

	// 🆕 ROADMAP TRACKING VARIABLES
	SearchFilter string // Stores text input for on-the-fly list filtering
	PanelFocus   string // Panel focus toggle state: "left" or "right"
}

// NewRuntimeContext builds a default tracking state
func NewRuntimeContext() *RuntimeContext {
	return &RuntimeContext{
		CurrentStep:  StateDashboard,
		PipelineMode: false,
		SearchFilter: "",
		PanelFocus:   "left", // Default keyboard interaction focuses left panel
	}
}
