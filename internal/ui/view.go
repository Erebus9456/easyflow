package ui

import (
	"fmt"
	"strings"

	"github.com/Erebus9456/easyflow/internal/workflow"

	"github.com/charmbracelet/lipgloss"
)

func (m AppModel) View() string {
	var s strings.Builder

	// Render standardized platform application header
	s.WriteString(StyleTitle.Render(" EasyFlow Workflow Dashboard 🚀 "))
	s.WriteString("\n")
	s.WriteString(fmt.Sprintf("📍 Active Scope Context: %s / %s\n", m.RepoCtx.Owner, m.RepoCtx.RepositoryName))
	s.WriteString("--------------------------------------------------\n\n")

	// Render global loading overlays or error diagnostics panels
	if m.Loading {
		s.WriteString(fmt.Sprintf("%s Fetching updates from system execution threads...\n\n", m.Spinner.View()))
	}
	if m.ErrorMessage != "" {
		s.WriteString(StyleErrorBanner.Render(fmt.Sprintf("Error Trace: %s", m.ErrorMessage)) + "\n\n")
	}
	if m.SuccessMsg != "" {
		s.WriteString(StyleSuccessBanner.Render(m.SuccessMsg) + "\n\n")
	}

	// Render interaction frames depending on active state context
	switch m.Engine.Ctx.CurrentStep {
	case workflow.StateDashboard:
		s.WriteString(StyleHeader.Render("Select Workflow Task Command:"))
		s.WriteString("\n")
		for i, item := range m.MenuItems {
			if m.Cursor == i {
				s.WriteString(StyleSelectedOption.Render(fmt.Sprintf("%s -> %s", item.Title, item.Description)))
			} else {
				s.WriteString(StyleUnselectedOption.Render(fmt.Sprintf("%s", item.Title)))
			}
			s.WriteString("\n")
		}

	case workflow.StateSelectIssue:
		s.WriteString(StyleHeader.Render("Select Target Tracking Issue for Pipeline Work:"))
		s.WriteString("\n")
		if len(m.Issues) == 0 && !m.Loading {
			s.WriteString(lipgloss.NewStyle().Foreground(ColorTextMuted).Render("No open tracked issues found inside this repository target.\n\n"))
			s.WriteString(StyleSelectedOption.Render("Press [n] to create a brand new GitHub Issue directly!"))
			s.WriteString("\n")
		} else {
			s.WriteString(StyleHelpText.Render("Tip: Press [n] to create a brand new issue on the fly\n\n"))
			for i, issue := range m.Issues {
				if m.IssueCursor == i {
					s.WriteString(StyleSelectedOption.Render(fmt.Sprintf("#%d: %s", issue.Number, issue.Title)))
				} else {
					s.WriteString(StyleUnselectedOption.Render(fmt.Sprintf("#%d: %s", issue.Number, issue.Title)))
				}
				s.WriteString("\n")
			}
		}

	case workflow.StateCreateIssue:
		s.WriteString(StyleHeader.Render("🆕 Create a New GitHub Issue:"))
		s.WriteString("\n\n")
		s.WriteString(m.TextInput.View())
		s.WriteString("\n")

	case workflow.StateCreateBranch:
		s.WriteString(StyleHeader.Render(fmt.Sprintf("Assign localized Git branch for Issue #%d:", m.Engine.Ctx.ActiveIssueNumber)))
		s.WriteString("\n\n")
		s.WriteString(m.TextInput.View())
		s.WriteString("\n")

	case workflow.StateWorking:
		s.WriteString(StyleHeader.Render("🛠️ Pipeline Branch Established Status: Active"))
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf("Branch tracking target mapped out: %s\n\n", m.Engine.Ctx.BranchName))
		s.WriteString("👉 Go complete your code adjustments in another terminal panel.\n")
		s.WriteString("When you are finished modifying files, hit [ENTER] here to capture your changes.\n\n")
		s.WriteString(StyleSelectedOption.Render("Press [ENTER] to move to staging and local commits."))
		s.WriteString("\n")

	case workflow.StateCommitReady:
		s.WriteString(StyleHeader.Render("💾 Commit Staged Changes:"))
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Foreground(ColorTextMuted).Render("This will perform a structural execution equivalent to running 'git add .'\n\n"))
		s.WriteString(m.TextInput.View())
		s.WriteString("\n")

	case workflow.StatePushing:
		s.WriteString(StyleHeader.Render("📤 Syncing Local Commits Upstream..."))
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf("%s Running underlying process: 'git push -u origin HEAD'\n", m.Spinner.View()))

	case workflow.StatePRPending:
		s.WriteString(StyleHeader.Render("🐙 Create GitHub Pull Request:"))
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Foreground(ColorTextMuted).Render("Specify the title text context for your upstream remote Pull Request merging block:\n\n"))
		s.WriteString(m.TextInput.View())
		s.WriteString("\n")

	case workflow.StateMerging:
		s.WriteString(StyleHeader.Render("🚀 Ship It! Final Pipeline Merge Authorization Step:"))
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf("Pull Request URL Detected: %s\n\n", m.Engine.Ctx.PullRequestURL))
		s.WriteString("This step will execute the following automated actions:\n")
		s.WriteString(fmt.Sprintf("  1. Merge your PR into the remote destination branch\n"))
		s.WriteString(fmt.Sprintf("  2. Delete the remote branch tracker securely\n"))
		s.WriteString(fmt.Sprintf("  3. Resolve and close issue target tracking point #%d\n\n", m.Engine.Ctx.ActiveIssueNumber))
		s.WriteString(StyleSelectedOption.Render("Press [ENTER] to execute full clean and remote resolution sequence."))
		s.WriteString("\n")

	case workflow.StateCompleted:
		s.WriteString(StyleSuccessBanner.Render("🎉 Continuous Development Cycle Successfully Executed!"))
		s.WriteString("\n\n")
		s.WriteString("All workspace contexts clean, PRs combined, and tracking targets resolved.\n")
		s.WriteString("Press [ESC] to transition safely back to the main menu screen dashboard.\n")
	}

	// Navigation Footer block
	s.WriteString(StyleHelpText.Render("\n[↑/↓ or j/k: Nav]  [Enter: Select/Advance]  [Esc: Reset to Dashboard]  [q: Exit Process]"))
	return s.String()
}
