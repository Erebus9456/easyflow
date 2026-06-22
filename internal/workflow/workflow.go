package workflow

import "fmt"

// Engine orchestrates legal pipeline advancements
type Engine struct {
	Ctx *RuntimeContext
}

// NewEngine instantiates a managed workflow context engine
func NewEngine() *Engine {
	return &Engine{
		Ctx: NewRuntimeContext(),
	}
}

// Advance transitions safely to the next milestone loop step
func (e *Engine) Advance(next State) error {
	// Expanded state validation rules
	switch next {
	case StateCreateBranch:
		if e.Ctx.ActiveIssueNumber == 0 {
			return fmt.Errorf("cannot initialize branch mapping: no issue selected")
		}
	case StateCommitReady:
		if e.Ctx.BranchName == "" {
			return fmt.Errorf("cannot prepare commit: no working branch active")
		}
	case StatePRPending:
		if e.Ctx.BranchName == "" {
			return fmt.Errorf("cannot initiate PR build sequence: missing working branch")
		}
	case StateMerging:
		if e.Ctx.PullRequestURL == "" {
			return fmt.Errorf("cannot merge: no pull request URL detected")
		}
	}

	e.Ctx.CurrentStep = next
	return nil
}

// Reset clears out transactional context flags back to the main dashboard
func (e *Engine) Reset() {
	e.Ctx.ActiveIssueNumber = 0
	e.Ctx.ActiveIssueTitle = ""
	e.Ctx.BranchName = ""
	e.Ctx.PullRequestURL = ""
	e.Ctx.CurrentStep = StateDashboard
	e.Ctx.PipelineMode = false // Clean up loop configuration tracking
}
