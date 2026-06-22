package ui

import (
	"fmt"
	"strconv"

	"github.com/Erebus9456/easyflow/internal/git"
	"github.com/Erebus9456/easyflow/internal/github"
	"github.com/Erebus9456/easyflow/internal/workflow"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Async status response messages wrapper tokens
type issuesMsg []github.Issue
type branchesMsg []string
type commitsMsg []string
type errMsg error
type actionSuccessMsg string

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			m.Engine.Reset()
			m.Cursor = 0
			m.IssueCursor = 0
			m.ErrorMessage = ""
			m.SuccessMsg = ""
			return m, nil
		}

	case spinner.TickMsg:
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case issuesMsg:
		m.Loading = false
		m.Issues = msg
		return m, nil

	case branchesMsg:
		m.Loading = false
		// Temporarily reuse Issues field or create lightweight mapping if necessary
		// For clean structural compliance, we convert string arrays into a readable format or save to local structures
		m.Issues = nil
		for i, b := range msg {
			m.Issues = append(m.Issues, github.Issue{Number: i + 1, Title: b})
		}
		return m, nil

	case commitsMsg:
		m.Loading = false
		m.Issues = nil
		for i, c := range msg {
			m.Issues = append(m.Issues, github.Issue{Number: i + 1, Title: c})
		}
		return m, nil

	case actionSuccessMsg:
		m.Loading = false
		m.SuccessMsg = string(msg)
		return m, nil

	case errMsg:
		m.Loading = false
		m.ErrorMessage = msg.Error()
		return m, nil
	}

	// Route internal update configurations based on active state steps
	switch m.Engine.Ctx.CurrentStep {
	case workflow.StateDashboard:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "up", "k":
				if m.Cursor > 0 {
					m.Cursor--
				}
			case "down", "j":
				if m.Cursor < len(m.MenuItems)-1 {
					m.Cursor++
				}
			case "enter":
				return m.handleMenuSelection()
			}
		}

	case workflow.StateManageIssues:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			subOptions := GetSubMenuOptions("issues")
			switch keyMsg.String() {
			case "up", "k":
				if m.IssueCursor > 0 {
					m.IssueCursor--
				}
			case "down", "j":
				if m.IssueCursor < len(subOptions)-1 {
					m.IssueCursor++
				}
			case "enter":
				switch m.IssueCursor {
				case 0: // List
					m.Engine.Advance(workflow.StateSelectIssue)
					m.Loading = true
					return m, func() tea.Msg {
						issues, err := github.FetchOpenIssues()
						if err != nil {
							return errMsg(err)
						}
						return issuesMsg(issues)
					}
				case 1: // Create
					m.Engine.Advance(workflow.StateCreateIssue)
					m.TextInput.Placeholder = "Enter new issue title..."
					m.TextInput.SetValue("")
				case 2: // Close by ID
					m.Engine.Advance(workflow.StateCreateIssue) // Reuse issue creation input field view for input matching
					m.TextInput.Placeholder = "Enter raw Issue Number to close (e.g. 42)..."
					m.TextInput.SetValue("")
				}
			}
		}

	case workflow.StateManageBranches:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			subOptions := GetSubMenuOptions("branches")
			switch keyMsg.String() {
			case "up", "k":
				if m.IssueCursor > 0 {
					m.IssueCursor--
				}
			case "down", "j":
				if m.IssueCursor < len(subOptions)-1 {
					m.IssueCursor++
				}
			case "enter":
				switch m.IssueCursor {
				case 0, 2: // List existing for Checkout or Deletion
					m.Engine.Advance(workflow.StateListBranches)
					m.Loading = true
					// Clear previous selection markers
					m.Cursor = m.IssueCursor // Track whether context was Checkout (0) or Delete (2)
					m.IssueCursor = 0
					return m, func() tea.Msg {
						branches, err := git.ListLocalBranches()
						if err != nil {
							return errMsg(err)
						}
						return branchesMsg(branches)
					}
				case 1: // Create Custom
					m.Engine.Advance(workflow.StateCreateBranch)
					m.TextInput.Placeholder = "Enter custom experimental branch name..."
					m.TextInput.SetValue("")
				}
			}
		}

	case workflow.StateManageCommits:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			subOptions := GetSubMenuOptions("commits")
			switch keyMsg.String() {
			case "up", "k":
				if m.IssueCursor > 0 {
					m.IssueCursor--
				}
			case "down", "j":
				if m.IssueCursor < len(subOptions)-1 {
					m.IssueCursor++
				}
			case "enter":
				switch m.IssueCursor {
				case 0: // Mini Git Log History Read
					m.Engine.Advance(workflow.StateViewCommits)
					m.Loading = true
					return m, func() tea.Msg {
						logs, err := git.GetLocalCommitLog(5)
						if err != nil {
							return errMsg(err)
						}
						return commitsMsg(logs)
					}
				case 1: // Stage & Commit
					m.Engine.Advance(workflow.StateCommitReady)
					m.TextInput.Placeholder = "Write standalone commit message..."
					m.TextInput.SetValue("")
				case 2: // Undo Soft Reset
					m.Loading = true
					return m, func() tea.Msg {
						if err := git.UndoLastCommit(); err != nil {
							return errMsg(err)
						}
						return actionSuccessMsg("Executed local 'git reset --soft HEAD~1' successfully!")
					}
				}
			}
		}

	case workflow.StateListBranches:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "up", "k":
				if m.IssueCursor > 0 {
					m.IssueCursor--
				}
			case "down", "j":
				if m.IssueCursor < len(m.Issues)-1 {
					m.IssueCursor++
				}
			case "enter":
				if len(m.Issues) > 0 {
					targetBranch := m.Issues[m.IssueCursor].Title
					m.Loading = true
					if m.Cursor == 0 { // Checkout Mode
						return m, func() tea.Msg {
							if err := git.CheckoutBranch(targetBranch); err != nil {
								return errMsg(err)
							}
							m.Engine.Ctx.BranchName = targetBranch
							m.Engine.Advance(workflow.StateDashboard)
							return actionSuccessMsg(fmt.Sprintf("Switched cleanly onto branch: %s", targetBranch))
						}
					} else { // Delete Mode
						return m, func() tea.Msg {
							if err := git.DeleteLocalBranch(targetBranch); err != nil {
								return errMsg(err)
							}
							m.Engine.Advance(workflow.StateDashboard)
							return actionSuccessMsg(fmt.Sprintf("Safely deleted workspace branch: %s", targetBranch))
						}
					}
				}
			}
		}

	case workflow.StateSelectIssue:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "up", "k":
				if m.IssueCursor > 0 {
					m.IssueCursor--
				}
			case "down", "j":
				if m.IssueCursor < len(m.Issues)-1 {
					m.IssueCursor++
				}
			case "n":
				m.Engine.Advance(workflow.StateCreateIssue)
				m.TextInput.Placeholder = "Enter new issue title..."
				m.TextInput.SetValue("")
				return m, nil
			case "enter":
				if len(m.Issues) > 0 {
					selected := m.Issues[m.IssueCursor]
					m.Engine.Ctx.ActiveIssueNumber = selected.Number
					m.Engine.Ctx.ActiveIssueTitle = selected.Title
					m.Engine.Advance(workflow.StateCreateBranch)
					m.TextInput.SetValue(fmt.Sprintf("issue-%d", selected.Number))
				}
			}
		}

	case workflow.StateCreateIssue:
		m.TextInput, cmd = m.TextInput.Update(msg)
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			inputVal := m.TextInput.Value()
			if inputVal == "" {
				return m, nil
			}

			// Fallback routing: Check if user is actually running "Close Issue by ID" action
			if m.TextInput.Placeholder == "Enter raw Issue Number to close (e.g. 42)..." {
				num, err := strconv.Atoi(inputVal)
				if err != nil {
					return m, nil
				}
				m.Loading = true
				return m, func() tea.Msg {
					if err := github.CloseIssue(num); err != nil {
						return errMsg(err)
					}
					m.Engine.Advance(workflow.StateDashboard)
					return actionSuccessMsg(fmt.Sprintf("GitHub Issue #%d resolved and closed successfully!", num))
				}
			}

			// Otherwise continue with standard creation route
			m.Loading = true
			return m, func() tea.Msg {
				num, err := github.CreateIssue(inputVal)
				if err != nil {
					return errMsg(err)
				}
				m.Engine.Ctx.ActiveIssueNumber = num
				m.Engine.Ctx.ActiveIssueTitle = inputVal
				m.Engine.Advance(workflow.StateCreateBranch)
				m.TextInput.SetValue(fmt.Sprintf("issue-%d", num))
				return actionSuccessMsg(fmt.Sprintf("Issue #%d created successfully!", num))
			}
		}
		return m, cmd

	case workflow.StateCreateBranch:
		m.TextInput, cmd = m.TextInput.Update(msg)
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			branchVal := m.TextInput.Value()
			m.Loading = true

			if m.Engine.Ctx.PipelineMode {
				m.Engine.Advance(workflow.StateWorking)
			} else {
				m.Engine.Advance(workflow.StateDashboard)
			}

			return m, func() tea.Msg {
				_, err := git.CreateAndCheckoutBranch(branchVal)
				if err != nil {
					return errMsg(err)
				}
				m.Engine.Ctx.BranchName = git.SanitizeBranchName(branchVal)
				return actionSuccessMsg("Workspace branch checked out cleanly!")
			}
		}
		return m, cmd

	case workflow.StateWorking:
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			m.Engine.Advance(workflow.StateCommitReady)
			m.TextInput.Placeholder = "Write standard local commit message..."
			m.TextInput.SetValue("")
			return m, nil
		}

	case workflow.StateCommitReady:
		m.TextInput, cmd = m.TextInput.Update(msg)
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			commitMsg := m.TextInput.Value()
			if commitMsg == "" {
				return m, nil
			}
			m.Loading = true
			m.Engine.Advance(workflow.StatePushing)
			return m, func() tea.Msg {
				if err := git.StageAllChanges(); err != nil {
					return errMsg(err)
				}
				if _, err := git.CreateCommit(commitMsg); err != nil {
					return errMsg(err)
				}
				return actionSuccessMsg("Changes committed locally!")
			}
		}
		return m, cmd

	case workflow.StatePushing:
		m.Loading = true
		return m, func() tea.Msg {
			_, err := git.PushToRemote()
			if err != nil {
				return errMsg(err)
			}
			if m.Engine.Ctx.PipelineMode {
				m.Engine.Advance(workflow.StatePRPending)
				m.TextInput.Placeholder = "Enter Pull Request Title..."
				m.TextInput.SetValue(m.Engine.Ctx.ActiveIssueTitle)
				return actionSuccessMsg("Branch pushed! Moving to Pull Request creation...")
			}
			m.Engine.Advance(workflow.StateDashboard)
			return actionSuccessMsg("Changes synced upstream securely!")
		}

	case workflow.StatePRPending:
		m.TextInput, cmd = m.TextInput.Update(msg)
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			prTitle := m.TextInput.Value()
			if prTitle == "" {
				return m, nil
			}
			m.Loading = true
			return m, func() tea.Msg {
				bodyMsg := fmt.Sprintf("Closes #%d. Automated via EasyFlow pipeline.", m.Engine.Ctx.ActiveIssueNumber)
				url, err := github.CreatePullRequest(prTitle, bodyMsg)
				if err != nil {
					return errMsg(err)
				}
				m.Engine.Ctx.PullRequestURL = url
				m.Engine.Advance(workflow.StateMerging)
				return actionSuccessMsg("Pull Request generated successfully!")
			}
		}
		return m, cmd

	case workflow.StateMerging:
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			m.Loading = true
			return m, func() tea.Msg {
				if _, err := github.MergeAndCleanupPR(); err != nil {
					return errMsg(err)
				}
				if m.Engine.Ctx.ActiveIssueNumber != 0 {
					if err := github.CloseIssue(m.Engine.Ctx.ActiveIssueNumber); err != nil {
						return errMsg(err)
					}
				}
				m.Engine.Advance(workflow.StateCompleted)
				return actionSuccessMsg("Pipeline execution loop complete!")
			}
		}
	}

	return m, nil
}

func (m AppModel) handleMenuSelection() (tea.Model, tea.Cmd) {
	m.ErrorMessage = ""
	m.SuccessMsg = ""
	m.IssueCursor = 0 // Clear list selection index for upcoming submenus

	switch m.Cursor {
	case 0: // Manage Issues Menu
		m.Engine.Advance(workflow.StateManageIssues)
	case 1: // Manage Branches Menu
		m.Engine.Advance(workflow.StateManageBranches)
	case 2: // Manage Commits Menu
		m.Engine.Advance(workflow.StateManageCommits)
	case 3: // Start continuous pipeline loop
		m.Engine.Advance(workflow.StateSelectIssue)
		m.Engine.Ctx.PipelineMode = true
		m.Loading = true
		return m, func() tea.Msg {
			issues, err := github.FetchOpenIssues()
			if err != nil {
				return errMsg(err)
			}
			return issuesMsg(issues)
		}
	case 4: // Manual Stage & Commit
		m.Engine.Advance(workflow.StateCommitReady)
		m.Engine.Ctx.PipelineMode = false
		m.TextInput.Placeholder = "Write standard local commit message..."
		m.TextInput.SetValue("")
	case 5: // Manual Push Sync
		m.Engine.Advance(workflow.StatePushing)
		m.Engine.Ctx.PipelineMode = false
	case 6: // Reset state engine
		m.Engine.Reset()
	}
	return m, nil
}
