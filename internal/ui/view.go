package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Erebus9456/easyflow/internal/workflow"
	"github.com/charmbracelet/lipgloss"
)

func (m AppModel) View() string {
	var leftSide strings.Builder
	var rightSide strings.Builder
	var finalView strings.Builder

	// Dynamically compute layout metrics from config parameters
	spacingStr := strings.Repeat("\n", m.Config.Layout.MenuSpacing)

	// 1. DYNAMIC FOCUS CONTRAST LAYOUT STYLING
	// Dim or enrich text contrast colors automatically depending on which side of the app is active
	leftForeground := lipgloss.Color("#FFFFFF")
	rightForeground := lipgloss.Color("#FFFFFF")
	if m.Engine.Ctx.PanelFocus == "right" {
		leftForeground = ColorTextMuted
	} else {
		rightForeground = ColorTextMuted
	}

	var (
		leftColumn = lipgloss.NewStyle().
				Width(m.Config.Layout.ColumnWidth).
				Foreground(leftForeground).
				PaddingRight(2)
		rightColumn = lipgloss.NewStyle().
				Width(m.Config.Layout.ColumnWidth).
				Foreground(rightForeground).
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(ColorNeutral).
				PaddingLeft(3)
	)

	// =========================================================================
	// 🌲 LEFT PANEL: ASSEMBLE INTERACTIVE ACTIONS
	// =========================================================================
	switch m.Engine.Ctx.CurrentStep {
	case workflow.StateDashboard:
		leftSide.WriteString(StyleHeader.Render("Select Workflow Task Command:"))
		leftSide.WriteString("\n\n")
		for i, item := range m.MenuItems {
			if m.Cursor == i {
				leftSide.WriteString(StyleSelectedOption.Render(fmt.Sprintf("> %s\n  %s", item.Title, item.Description)))
			} else {
				leftSide.WriteString(StyleUnselectedOption.Render(fmt.Sprintf("  %s", item.Title)))
			}
			leftSide.WriteString(spacingStr)
		}

	case workflow.StateUnsavedChangesWarning: // 🆕 RENDER: Phase 1 Safety Shield View
		leftSide.WriteString(StyleErrorBanner.Render("⚠️ WARNING: UNCOMMITTED CHANGES DETECTED"))
		leftSide.WriteString("\n\n")
		leftSide.WriteString("Your workspace contains dirty local updates that could be accidentally overwritten.\n\n")
		leftSide.WriteString("Options:\n")
		leftSide.WriteString(StyleSelectedOption.Render("  [Y] Bypasses warning shield and continue onto tasks anyway\n"))
		leftSide.WriteString(StyleUnselectedOption.Render("  [N / Enter] Return immediately to safety on the dashboard\n"))

	case workflow.StateSettingsMenu: // 🆕 RENDER: Interactive Settings Layout View
		leftSide.WriteString(StyleHeader.Render("⚙️ EasyFlow Configuration Profile Settings:"))
		leftSide.WriteString("\n\n")
		subOptions := GetSubMenuOptions("settings")
		for i, item := range subOptions {
			var currentSettingVal string
			switch i {
			case 0:
				currentSettingVal = fmt.Sprintf("[ %d Lines ]", m.Config.Layout.MenuSpacing)
			case 1:
				currentSettingVal = fmt.Sprintf("[ %d Chars ]", m.Config.Layout.ColumnWidth)
			case 2:
				currentSettingVal = fmt.Sprintf("[ %s ]", m.Config.Workflow.MergePolicy)
			case 3:
				currentSettingVal = fmt.Sprintf("[ Shield: %t ]", m.Config.Workflow.SafetyShield)
			}

			if m.IssueCursor == i {
				leftSide.WriteString(StyleSelectedOption.Render(fmt.Sprintf("> %s : %s\n  %s", item.Title, currentSettingVal, item.Description)))
			} else {
				leftSide.WriteString(StyleUnselectedOption.Render(fmt.Sprintf("  %s : %s", item.Title, currentSettingVal)))
			}
			leftSide.WriteString(spacingStr)
		}
		leftSide.WriteString(StyleHelpText.Render("\n💡 Tip: Press [ENTER] to cycle settings. Press [ESC] to save & return."))

	case workflow.StateManageStash: // 🆕 RENDER: Stash Submenu List View
		leftSide.WriteString(StyleHeader.Render("📦 Stash Shelf Control Actions:"))
		leftSide.WriteString("\n\n")
		subOptions := GetSubMenuOptions("stash")
		for i, item := range subOptions {
			if m.IssueCursor == i {
				leftSide.WriteString(StyleSelectedOption.Render(fmt.Sprintf("> %s\n  %s", item.Title, item.Description)))
			} else {
				leftSide.WriteString(StyleUnselectedOption.Render(fmt.Sprintf("  %s", item.Title)))
			}
			leftSide.WriteString(spacingStr)
		}

	case workflow.StateManageIssues:
		leftSide.WriteString(StyleHeader.Render("🐛 Issue Management Actions:"))
		leftSide.WriteString("\n\n")
		subOptions := GetSubMenuOptions("issues")
		for i, item := range subOptions {
			if m.IssueCursor == i {
				leftSide.WriteString(StyleSelectedOption.Render(fmt.Sprintf("> %s\n  %s", item.Title, item.Description)))
			} else {
				leftSide.WriteString(StyleUnselectedOption.Render(fmt.Sprintf("  %s", item.Title)))
			}
			leftSide.WriteString(spacingStr)
		}

	case workflow.StateManageBranches:
		leftSide.WriteString(StyleHeader.Render("🌿 Branch Management Actions:"))
		leftSide.WriteString("\n\n")
		subOptions := GetSubMenuOptions("branches")
		for i, item := range subOptions {
			if m.IssueCursor == i {
				leftSide.WriteString(StyleSelectedOption.Render(fmt.Sprintf("> %s\n  %s", item.Title, item.Description)))
			} else {
				leftSide.WriteString(StyleUnselectedOption.Render(fmt.Sprintf("  %s", item.Title)))
			}
			leftSide.WriteString(spacingStr)
		}

	case workflow.StateManageCommits:
		leftSide.WriteString(StyleHeader.Render("💾 Commit Management Actions:"))
		leftSide.WriteString("\n\n")
		subOptions := GetSubMenuOptions("commits")
		for i, item := range subOptions {
			if m.IssueCursor == i {
				leftSide.WriteString(StyleSelectedOption.Render(fmt.Sprintf("> %s\n  %s", item.Title, item.Description)))
			} else {
				leftSide.WriteString(StyleUnselectedOption.Render(fmt.Sprintf("  %s", item.Title)))
			}
			leftSide.WriteString(spacingStr)
		}

	case workflow.StateListBranches:
		if m.Cursor == 0 {
			leftSide.WriteString(StyleHeader.Render("🌿 Select Local Branch to Checkout:"))
		} else {
			leftSide.WriteString(StyleHeader.Render("🗑️ Select Local Branch to Delete:"))
		}
		leftSide.WriteString("\n")
		// 🆕 Live Filter Header input box display
		leftSide.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFD2")).Render(fmt.Sprintf("🔍 Live List Search Filter: %s█\n\n", m.Engine.Ctx.SearchFilter)))

		if len(m.Issues) == 0 && !m.Loading {
			leftSide.WriteString(lipgloss.NewStyle().Foreground(ColorTextMuted).Render("No local branches discovered.\n"))
		} else {
			for i, b := range m.Issues {
				// Filter list matching checks before printing line layout
				if m.Engine.Ctx.SearchFilter != "" && !strings.Contains(strings.ToLower(b.Title), strings.ToLower(m.Engine.Ctx.SearchFilter)) {
					continue
				}
				if m.IssueCursor == i {
					leftSide.WriteString(StyleSelectedOption.Render(fmt.Sprintf("> %s", b.Title)))
				} else {
					leftSide.WriteString(StyleUnselectedOption.Render(fmt.Sprintf("  %s", b.Title)))
				}
				leftSide.WriteString(spacingStr)
			}
		}

	case workflow.StateViewCommits:
		leftSide.WriteString(StyleHeader.Render("📋 Recent Local Commit History Log:"))
		leftSide.WriteString("\n\n")
		if len(m.Issues) == 0 && !m.Loading {
			leftSide.WriteString(lipgloss.NewStyle().Foreground(ColorTextMuted).Render("No commit revisions found on this branch.\n"))
		} else {
			for _, logEntry := range m.Issues {
				leftSide.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#E2E2E2")).Render(fmt.Sprintf(" %s", logEntry.Title)))
				leftSide.WriteString(spacingStr)
			}
			leftSide.WriteString(StyleHelpText.Render("\nPress [ESC] to return to dashboard..."))
		}

	case workflow.StateSelectIssue:
		leftSide.WriteString(StyleHeader.Render("Select Target Tracking Issue:"))
		leftSide.WriteString("\n")
		// 🆕 Live Filter Header input box display
		leftSide.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFD2")).Render(fmt.Sprintf("🔍 Live List Search Filter: %s█\n\n", m.Engine.Ctx.SearchFilter)))

		if len(m.Issues) == 0 && !m.Loading {
			leftSide.WriteString(lipgloss.NewStyle().Foreground(ColorTextMuted).Render("No open tracked issues found inside target repository.\n\n"))
			leftSide.WriteString(StyleSelectedOption.Render("Press [n] to create a brand new GitHub Issue directly!"))
			leftSide.WriteString("\n")
		} else {
			leftSide.WriteString(StyleHelpText.Render("Tip: Press [n] to create a brand new issue on the fly\n\n"))
			for i, issue := range m.Issues {
				// Filter list matching checks before printing line layout
				if m.Engine.Ctx.SearchFilter != "" && !strings.Contains(strings.ToLower(issue.Title), strings.ToLower(m.Engine.Ctx.SearchFilter)) {
					continue
				}
				if m.IssueCursor == i {
					leftSide.WriteString(StyleSelectedOption.Render(fmt.Sprintf("> #%d: %s", issue.Number, issue.Title)))
				} else {
					leftSide.WriteString(StyleUnselectedOption.Render(fmt.Sprintf("  #%d: %s", issue.Number, issue.Title)))
				}
				leftSide.WriteString(spacingStr)
			}
		}

	case workflow.StateCreateIssue:
		if m.TextInput.Placeholder == "Enter raw Issue Number to close (e.g. 42)..." {
			leftSide.WriteString(StyleHeader.Render("🗑️ Close a Remote GitHub Issue:"))
		} else {
			leftSide.WriteString(StyleHeader.Render("🆕 Create a New GitHub Issue:"))
		}
		leftSide.WriteString("\n\n")
		leftSide.WriteString(m.TextInput.View())
		leftSide.WriteString("\n")

	case workflow.StateCreateBranch:
		leftSide.WriteString(StyleHeader.Render("🌿 Specify Branch Title Context:"))
		leftSide.WriteString("\n\n")
		leftSide.WriteString(m.TextInput.View())
		leftSide.WriteString("\n")

	case workflow.StateWorking:
		leftSide.WriteString(StyleHeader.Render("🛠️ Working on Active Branch"))
		leftSide.WriteString("\n\n")
		leftSide.WriteString("👉 Go complete your code adjustments in another terminal panel.\n\n")
		leftSide.WriteString("When you are finished modifying files, hit [ENTER] here to capture your changes.\n\n")
		leftSide.WriteString(StyleSelectedOption.Render("Press [ENTER] to stage and commit changes."))
		leftSide.WriteString("\n")

	case workflow.StateCommitReady:
		leftSide.WriteString(StyleHeader.Render("💾 Commit Staged Changes:"))
		leftSide.WriteString("\n\n")
		leftSide.WriteString(lipgloss.NewStyle().Foreground(ColorTextMuted).Render("This will automatically execute 'git add .'\n\n"))
		leftSide.WriteString(m.TextInput.View())
		leftSide.WriteString("\n")

	case workflow.StatePushing:
		leftSide.WriteString(StyleHeader.Render("📤 Syncing Commits Upstream..."))
		leftSide.WriteString("\n\n")
		leftSide.WriteString(fmt.Sprintf("%s Running: 'git push -u origin HEAD'\n", m.Spinner.View()))

	case workflow.StatePRPending:
		leftSide.WriteString(StyleHeader.Render("🐙 Create Pull Request:"))
		leftSide.WriteString("\n\n")
		leftSide.WriteString(lipgloss.NewStyle().Foreground(ColorTextMuted).Render("Specify the title text for your PR merge block:\n\n"))
		leftSide.WriteString(m.TextInput.View())
		leftSide.WriteString("\n")

	case workflow.StateMerging:
		leftSide.WriteString(StyleHeader.Render("🚀 Ship It! Merge Authorization Step:"))
		leftSide.WriteString("\n\n")
		leftSide.WriteString("This step will execute the following automated actions:\n\n")
		leftSide.WriteString(fmt.Sprintf("  1. Merge your PR into upstream via %s\n", m.Config.Workflow.MergePolicy))
		leftSide.WriteString("  2. Delete the remote tracking branch safely\n")
		leftSide.WriteString(fmt.Sprintf("  3. Resolve and close issue target point #%d\n\n", m.Engine.Ctx.ActiveIssueNumber))
		leftSide.WriteString(StyleSelectedOption.Render("Press [ENTER] to execute full workspace resolution."))
		leftSide.WriteString("\n")

	case workflow.StateCompleted:
		leftSide.WriteString(StyleSuccessBanner.Render("🎉 Development Cycle Complete!"))
		leftSide.WriteString("\n\n")
		leftSide.WriteString("All workspace contexts clean, PRs combined, and tracking targets resolved.\n\n")
		leftSide.WriteString("Press [ESC] to return to the main dashboard menu.\n")
	}

	// =========================================================================
	// 📊 RIGHT PANEL: PERSISTENT METRICS & ENVIRONMENT HUB
	// =========================================================================
	rightSide.WriteString(StyleHeader.Render("📋 Current Workspace Status:"))
	rightSide.WriteString("\n\n")

	rightSide.WriteString(fmt.Sprintf("📍 Repo Context : %s / %s", m.RepoCtx.Owner, m.RepoCtx.RepositoryName))
	rightSide.WriteString(spacingStr)

	if m.Engine.Ctx.BranchName != "" {
		rightSide.WriteString(fmt.Sprintf("🌲 Target Branch: %s", m.Engine.Ctx.BranchName))
	} else {
		rightSide.WriteString(fmt.Sprintf("🌲 Target Branch: %s", m.RepoCtx.CurrentBranch))
	}
	rightSide.WriteString(spacingStr)

	if m.Engine.Ctx.ActiveIssueNumber != 0 {
		rightSide.WriteString(fmt.Sprintf("🐛 Linked Issue : #%d - %s", m.Engine.Ctx.ActiveIssueNumber, m.Engine.Ctx.ActiveIssueTitle))
	} else {
		rightSide.WriteString("🐛 Linked Issue : None selected")
	}
	rightSide.WriteString(spacingStr)

	if m.Engine.Ctx.PipelineMode {
		rightSide.WriteString("⚙️ Engine Mode  : Continuous Pipeline")
	} else {
		rightSide.WriteString("⚙️ Engine Mode  : Standalone Commands")
	}
	rightSide.WriteString(spacingStr)

	if m.Engine.Ctx.PullRequestURL != "" {
		rightSide.WriteString(fmt.Sprintf("🐙 PR Remote URL: %s", m.Engine.Ctx.PullRequestURL))
	} else {
		rightSide.WriteString("🐙 PR Remote URL: No pull request open")
	}
	rightSide.WriteString(spacingStr)

	currentTimeString := time.Now().Format("2006-01-02 15:04:05")
	rightSide.WriteString(lipgloss.NewStyle().Foreground(ColorTextMuted).Render(fmt.Sprintf("\n🕒 Live Dashboard Render: %s (Local)\n", currentTimeString)))

	// =========================================================================
	// 🗺️ COMBINE COLUMNS AND RENDER DASHBOARD FRAME
	// =========================================================================
	finalView.WriteString(StyleTitle.Render(" EasyFlow Workflow Dashboard 🚀 "))
	finalView.WriteString("\n\n")

	if m.Loading && m.Engine.Ctx.CurrentStep != workflow.StatePushing {
		finalView.WriteString(fmt.Sprintf("%s Fetching system runtime threads...\n\n", m.Spinner.View()))
	}
	if m.ErrorMessage != "" {
		finalView.WriteString(StyleErrorBanner.Render(fmt.Sprintf("Error Trace: %s", m.ErrorMessage)) + "\n\n")
	}
	if m.SuccessMsg != "" {
		finalView.WriteString(StyleSuccessBanner.Render(m.SuccessMsg) + "\n\n")
	}

	columnsJoined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumn.Render(leftSide.String()),
		rightColumn.Render(rightSide.String()),
	)
	finalView.WriteString(columnsJoined)
	finalView.WriteString("\n")

	// 🆕 Dynamic action key indicators display showing active panel focus toggles
	finalView.WriteString(StyleHelpText.Render("\n[↑/↓: Nav] [Enter: Run] [Tab: Swap Active Panel] [Esc: Reset Menu] [q: Exit]"))

	return finalView.String()
}
