package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Erebus9456/easyflow/internal/git"
	"github.com/Erebus9456/easyflow/internal/github"
	"github.com/Erebus9456/easyflow/internal/workflow"
	"github.com/Erebus9456/easyflow/utils"

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
		case "tab": // Dual-Panel Tab Focus Toggling
			if m.Engine.Ctx.PanelFocus == "left" {
				m.Engine.Ctx.PanelFocus = "right"
			} else {
				m.Engine.Ctx.PanelFocus = "left"
			}
			return m, nil
		case "esc":
			m.Engine.Reset()
			m.Engine.Ctx.SearchFilter = ""
			m.Engine.Ctx.PanelFocus = "left"
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
		// 🆕 Append a virtual fallback navigation choice at the bottom of the list
		m.Issues = append(m.Issues, github.Issue{Number: -1, Title: "➕ Create a brand new issue instead"})
		return m, nil

	case branchesMsg:
		m.Loading = false
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

	// 🆕 Capture Alphanumeric characters on selection lists for Live Searching
	// Added a preventative check ensuring key routing stays clear of swallowing shortcuts on isolated setups
	if m.Engine.Ctx.CurrentStep == workflow.StateSelectIssue && len(m.Issues) <= 1 {
		// Let the StateSelectIssue block exclusively process raw keystrokes if the menu holds no real issues
	} else if m.Engine.Ctx.CurrentStep == workflow.StateSelectIssue || m.Engine.Ctx.CurrentStep == workflow.StateListBranches {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			kStr := keyMsg.String()
			if len(kStr) == 1 && kStr >= " " && kStr <= "~" {
				m.Engine.Ctx.SearchFilter += kStr
				m.IssueCursor = 0 // Snap back to top row
				return m, nil
			} else if kStr == "backspace" && len(m.Engine.Ctx.SearchFilter) > 0 {
				m.Engine.Ctx.SearchFilter = m.Engine.Ctx.SearchFilter[:len(m.Engine.Ctx.SearchFilter)-1]
				m.IssueCursor = 0
				return m, nil
			}
		}
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

	case workflow.StateUnsavedChangesWarning: // Safety Intercept Screen
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch strings.ToLower(keyMsg.String()) {
			case "y": // Bypass safety shield and skip to target view
				m.Engine.Advance(workflow.StateManageBranches)
				m.IssueCursor = 0
			case "n", "enter": // Bounce back to dashboard safety
				m.Engine.Reset()
			}
		}

	case workflow.StateSettingsMenu: // Interactive Runtime Settings Controller
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			subOptions := GetSubMenuOptions("settings")
			switch keyMsg.String() {
			case "up", "k":
				if m.IssueCursor > 0 {
					m.IssueCursor--
				}
			case "down", "j":
				if m.IssueCursor < len(subOptions)-1 {
					m.IssueCursor++
				}
			case "enter": // Cycle through settings variables in memory
				switch m.IssueCursor {
				case 0: // Spacing Loop 1-4
					m.Config.Layout.MenuSpacing++
					if m.Config.Layout.MenuSpacing > 4 {
						m.Config.Layout.MenuSpacing = 1
					}
				case 1: // Width Loop 40-70
					m.Config.Layout.ColumnWidth += 10
					if m.Config.Layout.ColumnWidth > 70 {
						m.Config.Layout.ColumnWidth = 40
					}
				case 2: // Merge Rules
					switch m.Config.Workflow.MergePolicy {
					case "Standard Merge":
						m.Config.Workflow.MergePolicy = "Squash Merge"
					case "Squash Merge":
						m.Config.Workflow.MergePolicy = "Rebase"
					default:
						m.Config.Workflow.MergePolicy = "Standard Merge"
					}
				case 3: // Toggle Safety Check Boolean
					m.Config.Workflow.SafetyShield = !m.Config.Workflow.SafetyShield
				}
			}
		}

	case workflow.StateManageStash: // Stash Commands Control Layer
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			subOptions := GetSubMenuOptions("stash")
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
				m.Loading = true
				switch m.IssueCursor {
				case 0: // Save snapshot
					return m, func() tea.Msg {
						if _, _, err := utils.ExecuteCommand("git", "stash"); err != nil {
							return errMsg(err)
						}
						return actionSuccessMsg("Current workspace modifications safely stashed away!")
					}
				case 1: // Pop entry
					return m, func() tea.Msg {
						if _, _, err := utils.ExecuteCommand("git", "stash", "pop"); err != nil {
							return errMsg(err)
						}
						return actionSuccessMsg("Workspace state restored from stash stack successfully!")
					}
				case 2: // Clear shelf
					return m, func() tea.Msg {
						if _, _, err := utils.ExecuteCommand("git", "stash", "clear"); err != nil {
							return errMsg(err)
						}
						return actionSuccessMsg("Stash history wiped clean!")
					}
				}
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
					m.Engine.Ctx.SearchFilter = ""
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
					m.Engine.Advance(workflow.StateCreateIssue)
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
					m.Engine.Ctx.SearchFilter = ""
					m.Loading = true
					m.Cursor = m.IssueCursor // Track Checkout (0) vs Delete (2)
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
					m.Engine.Ctx.SearchFilter = ""
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
			kStr := keyMsg.String()

			// 1. If there are no active repository issues, back up 'n' as a direct forward trigger
			if (len(m.Issues) == 0 || (len(m.Issues) == 1 && m.Issues[0].Number == -1)) && kStr == "n" {
				m.Engine.Advance(workflow.StateCreateIssue)
				m.TextInput.Placeholder = "Enter new issue title..."
				m.TextInput.SetValue("")
				m.Engine.Ctx.SearchFilter = ""
				return m, nil
			}

			// 2. Standard Arrow/Navigation keys
			switch kStr {
			case "up", "k":
				if m.IssueCursor > 0 {
					m.IssueCursor--
				}
				return m, nil
			case "down", "j":
				if m.IssueCursor < len(m.Issues)-1 {
					m.IssueCursor++
				}
				return m, nil
			case "enter":
				if len(m.Issues) > 0 {
					selected := m.Issues[m.IssueCursor]

					// 🆕 Arrow Navigation selection override for the creation fallback option
					if selected.Number == -1 {
						m.Engine.Advance(workflow.StateCreateIssue)
						m.TextInput.Placeholder = "Enter new issue title..."
						m.TextInput.SetValue("")
						m.Engine.Ctx.SearchFilter = ""
						return m, nil
					}

					// Standard valid tracking issue loop setup routing path
					m.Engine.Ctx.ActiveIssueNumber = selected.Number
					m.Engine.Ctx.ActiveIssueTitle = selected.Title
					m.Engine.Ctx.SearchFilter = ""
					m.Engine.Advance(workflow.StateCreateBranch)
					m.TextInput.SetValue(fmt.Sprintf("issue-%d", selected.Number))
				}
				return m, nil
			}
		}

	case workflow.StateCreateIssue:
		m.TextInput, cmd = m.TextInput.Update(msg)
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			inputVal := m.TextInput.Value()
			if inputVal == "" {
				return m, nil
			}

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
	m.IssueCursor = 0

	// Phase 1 Safety Shield Check: Intercept branch modification routes if workspace is dirty
	if m.Config.Workflow.SafetyShield && (m.Cursor == 1 || m.Cursor == 3) {
		stdout, _, _ := utils.ExecuteCommand("git", "status", "--porcelain")
		if strings.TrimSpace(stdout) != "" {
			m.Engine.Advance(workflow.StateUnsavedChangesWarning)
			return m, nil
		}
	}

	switch m.Cursor {
	case 0: // Manage Issues Menu
		m.Engine.Advance(workflow.StateManageIssues)
	case 1: // Manage Branches Menu
		m.Engine.Advance(workflow.StateManageBranches)
	case 2: // Manage Commits Menu
		m.Engine.Advance(workflow.StateManageCommits)
	case 3: // Manage Stash Menu
		m.Engine.Advance(workflow.StateManageStash)
	case 4: // Start continuous pipeline loop
		m.Engine.Advance(workflow.StateSelectIssue)
		m.Engine.Ctx.PipelineMode = true
		m.Engine.Ctx.SearchFilter = ""
		m.Loading = true
		return m, func() tea.Msg {
			issues, err := github.FetchOpenIssues()
			if err != nil {
				return errMsg(err)
			}
			return issuesMsg(issues)
		}
	case 5: // Manual Stage & Commit
		m.Engine.Advance(workflow.StateCommitReady)
		m.Engine.Ctx.PipelineMode = false
		m.TextInput.Placeholder = "Write standard local commit message..."
		m.TextInput.SetValue("")
	case 6: // Manual Push Sync
		m.Engine.Advance(workflow.StatePushing)
		m.Engine.Ctx.PipelineMode = false
	case 7: // Reset state engine
		m.Engine.Reset()
	case 8: // App Settings Configuration Menu
		m.Engine.Advance(workflow.StateSettingsMenu)
	}
	return m, nil
}
