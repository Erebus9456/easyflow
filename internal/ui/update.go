package ui

import (
	"fmt"

	"github.com/Erebus9456/easyflow/internal/git"
	"github.com/Erebus9456/easyflow/internal/github"
	"github.com/Erebus9456/easyflow/internal/workflow"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Async status response messages wrapper tokens
type issuesMsg []github.Issue
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
			issueTitle := m.TextInput.Value()
			if issueTitle == "" {
				return m, nil
			}
			m.Loading = true
			return m, func() tea.Msg {
				num, err := github.CreateIssue(issueTitle)
				if err != nil {
					return errMsg(err)
				}

				m.Engine.Ctx.ActiveIssueNumber = num
				m.Engine.Ctx.ActiveIssueTitle = issueTitle
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
			m.Engine.Advance(workflow.StateWorking)
			return m, func() tea.Msg {
				_, err := git.CreateAndCheckoutBranch(branchVal)
				if err != nil {
					return errMsg(err)
				}
				m.Engine.Ctx.BranchName = git.SanitizeBranchName(branchVal)
				return actionSuccessMsg("Local workspace branch checked out cleanly!")
			}
		}
		return m, cmd

	case workflow.StateWorking:
		// Developer presses enter on the info screen once they are finished writing code
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
		// Automatically trigger push routine upon state landing entry
		m.Loading = true
		return m, func() tea.Msg {
			_, err := git.PushToRemote()
			if err != nil {
				return errMsg(err)
			}

			// If inside continuous loop mode, route forward. Otherwise drop back cleanly.
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

	switch m.Cursor {
	case 0: // Start continuous end-to-end pipeline loop
		m.Engine.Advance(workflow.StateSelectIssue)
		m.Engine.Ctx.PipelineMode = true // Lock in loop configuration
		m.Loading = true
		return m, func() tea.Msg {
			issues, err := github.FetchOpenIssues()
			if err != nil {
				return errMsg(err)
			}
			return issuesMsg(issues)
		}
	case 1: // One-off Manual Staging Commit
		m.Engine.Advance(workflow.StateCommitReady)
		m.Engine.Ctx.PipelineMode = false
		m.TextInput.Placeholder = "Write standard local commit message..."
		m.TextInput.SetValue("")
	case 2: // One-off Push Sync
		m.Engine.Advance(workflow.StatePushing)
		m.Engine.Ctx.PipelineMode = false
		return m, nil
	case 3: // Reset context state engine
		m.Engine.Reset()
	}
	return m, nil
}
