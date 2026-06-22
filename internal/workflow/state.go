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

	// New CRUD Sub-Category Submenu View States
	StateManageIssues
	StateManageBranches
	StateManageCommits

	// New Explicit Action States
	StateListBranches // Used for branch selection/checkout & deletion lists
	StateViewCommits  // Used for displaying git log history entries
)

// RuntimeContext retains persistent cross-state memory metrics
type RuntimeContext struct {
	ActiveIssueNumber int
	ActiveIssueTitle  string
	BranchName        string
	PullRequestURL    string
	CurrentStep       State
	PipelineMode      bool // Tracks if the user is in the continuous end-to-end loop
}

// NewRuntimeContext builds a default tracking state
func NewRuntimeContext() *RuntimeContext {
	return &RuntimeContext{
		CurrentStep:  StateDashboard,
		PipelineMode: false,
	}
}
