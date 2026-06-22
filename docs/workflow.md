# Workflow Guide

This document provides a comprehensive guide to the EasyFlow workflow automation system, including the state machine, usage patterns, and best practices.

## Table of Contents

- [Workflow Overview](#workflow-overview)
- [State Machine](#state-machine)
- [Pipeline Mode](#pipeline-mode)
- [Standalone Mode](#standalone-mode)
- [CRUD Operations](#crud-operations)
- [Usage Examples](#usage-examples)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

---

## Workflow Overview

EasyFlow provides two primary modes of operation:

1. **Pipeline Mode**: Automated end-to-end workflow from issue to merge
2. **Standalone Mode**: Individual operations for specific tasks

### Core Workflow Loop

```mermaid
graph LR
    A[Select Issue] --> B[Create Branch]
    B --> C[Write Code]
    C --> D[Commit Changes]
    D --> E[Push to Remote]
    E --> F[Create PR]
    F --> G[Merge PR]
    G --> H[Close Issue]
    H --> A
```

### Workflow States

```mermaid
stateDiagram-v2
    [*] --> Dashboard
    
    Dashboard --> SelectIssue: Start Pipeline
    Dashboard --> ManageIssues: Manage Issues
    Dashboard --> ManageBranches: Manage Branches
    Dashboard --> ManageCommits: Manage Commits
    Dashboard --> CommitReady: Manual Commit
    Dashboard --> Pushing: Manual Push
    
    SelectIssue --> CreateIssue: Press 'n'
    SelectIssue --> CreateBranch: Select Issue
    
    CreateIssue --> CreateBranch: Issue Created
    
    CreateBranch --> Working: Pipeline Mode
    CreateBranch --> Dashboard: Standalone Mode
    
    Working --> CommitReady: Press Enter
    
    CommitReady --> Pushing: Commit Created
    
    Pushing --> PRPending: Pipeline Mode
    Pushing --> Dashboard: Standalone Mode
    
    PRPending --> Merging: PR Created
    Merging --> Completed: PR Merged
    Completed --> Dashboard: Press ESC
    
    ManageIssues --> SelectIssue: List Issues
    ManageIssues --> CreateIssue: Create Issue
    ManageIssues --> Dashboard: Close Issue
    
    ManageBranches --> ListBranches: Checkout/Delete
    ListBranches --> Dashboard: Action Complete
    
    ManageCommits --> ViewCommits: View Log
    ViewCommits --> Dashboard: Press ESC
    ManageCommits --> CommitReady: Stage & Commit
    ManageCommits --> Dashboard: Undo Commit
```

---

## State Machine

### State Definitions

| State | Description | Required Context |
|-------|-------------|------------------|
| `StateDashboard` | Main menu display | None |
| `StateSelectIssue` | Issue selection for pipeline | None |
| `StateCreateIssue` | New issue creation | None |
| `StateCreateBranch` | Branch naming | `ActiveIssueNumber` |
| `StateWorking` | Code editing phase | `BranchName` |
| `StateCommitReady` | Commit message input | `BranchName` |
| `StatePushing` | Pushing to remote | `BranchName` |
| `StatePRPending` | PR creation input | `BranchName` |
| `StateMerging` | Merge authorization | `PullRequestURL` |
| `StateCompleted` | Pipeline completion | All context |
| `StateManageIssues` | Issue management submenu | None |
| `StateManageBranches` | Branch management submenu | None |
| `StateManageCommits` | Commit management submenu | None |
| `StateListBranches` | Branch selection list | None |
| `StateViewCommits` | Commit log display | None |

### State Transitions

```mermaid
graph TD
    A[Dashboard] -->|Start Pipeline| B[Select Issue]
    A -->|Manage Issues| C[Manage Issues]
    A -->|Manage Branches| D[Manage Branches]
    A -->|Manage Commits| E[Manage Commits]
    A -->|Manual Commit| F[Commit Ready]
    A -->|Manual Push| G[Pushing]
    
    B -->|Select Issue| H[Create Branch]
    B -->|Press 'n'| I[Create Issue]
    
    I -->|Issue Created| H
    
    H -->|Pipeline Mode| J[Working]
    H -->|Standalone Mode| A
    
    J -->|Press Enter| F
    
    F -->|Commit Created| G
    
    G -->|Pipeline Mode| K[PR Pending]
    G -->|Standalone Mode| A
    
    K -->|PR Created| L[Merging]
    
    L -->|PR Merged| M[Completed]
    M -->|Press ESC| A
    
    C -->|List Issues| B
    C -->|Create Issue| I
    C -->|Close Issue| A
    
    D -->|Checkout/Delete| N[List Branches]
    N -->|Action Complete| A
    
    E -->|View Log| O[View Commits]
    O -->|Press ESC| A
    E -->|Stage & Commit| F
    E -->|Undo Commit| A
```

### Validation Rules

The workflow engine enforces validation rules to prevent invalid state transitions:

```go
// Cannot create branch without issue
case StateCreateBranch:
    if e.Ctx.ActiveIssueNumber == 0 {
        return fmt.Errorf("cannot initialize branch mapping: no issue selected")
    }

// Cannot commit without branch
case StateCommitReady:
    if e.Ctx.BranchName == "" {
        return fmt.Errorf("cannot prepare commit: no working branch active")
    }

// Cannot create PR without branch
case StatePRPending:
    if e.Ctx.BranchName == "" {
        return fmt.Errorf("cannot initiate PR build sequence: missing working branch")
    }

// Cannot merge without PR
case StateMerging:
    if e.Ctx.PullRequestURL == "" {
        return fmt.Errorf("cannot merge: no pull request URL detected")
    }
```

---

## Pipeline Mode

Pipeline mode provides an automated end-to-end workflow from issue selection to merge.

### Pipeline Flow

```mermaid
sequenceDiagram
    participant User
    participant Dashboard
    participant SelectIssue
    participant CreateBranch
    participant Working
    participant Commit
    participant Push
    participant PR
    participant Merge
    participant Complete
    
    User->>Dashboard: Select "Start Pipeline"
    Dashboard->>SelectIssue: Show issues
    User->>SelectIssue: Select issue #123
    SelectIssue->>CreateBranch: Pre-fill branch name
    User->>CreateBranch: Confirm branch name
    CreateBranch->>Working: Switch to working state
    User->>Working: Write code in another terminal
    User->>Working: Press Enter when done
    Working->>Commit: Prompt for commit message
    User->>Commit: Enter commit message
    Commit->>Push: Push to remote
    Push->>PR: Prompt for PR title
    User->>PR: Enter PR title
    PR->>Merge: Prompt for merge
    User->>Merge: Confirm merge
    Merge->>Complete: Show completion
    Complete->>Dashboard: Return to dashboard
```

### Pipeline Steps

1. **Select Issue**
   - Browse open issues from repository
   - Press `n` to create new issue on-the-fly
   - Select existing issue to work on

2. **Create Branch**
   - Branch name auto-filled as `issue-{number}`
   - Customize branch name if desired
   - Branch created and checked out automatically

3. **Working**
   - Switch to another terminal to write code
   - Make changes to files
   - Press Enter when ready to commit

4. **Commit**
   - All changes automatically staged (`git add .`)
   - Enter commit message
   - Commit created locally

5. **Push**
   - Changes pushed to remote automatically
   - Upstream tracking set automatically

6. **Create PR**
   - PR title auto-filled with issue title
   - Customize PR title if desired
   - PR body auto-generated with issue reference

7. **Merge**
   - Review merge actions
   - Confirm to merge
   - Branch deleted automatically
   - Issue closed automatically

8. **Complete**
   - Success message displayed
   - Press ESC to return to dashboard

### Pipeline Advantages

- **Automated**: Minimal manual intervention
- **Guided**: Clear step-by-step progression
- **Safe**: Validation prevents errors
- **Fast**: No context switching between tools
- **Consistent**: Standardized workflow

---

## Standalone Mode

Standalone mode allows individual operations without the full pipeline.

### Standalone Operations

#### Manual Commit

```mermaid
graph LR
    A[Dashboard] --> B[Select Stage & Commit]
    B --> C[Enter Commit Message]
    C --> D[Changes Staged & Committed]
    D --> A
```

**Use Case**: Quick commit without full pipeline

**Steps**:
1. Select "Stage & Commit Local Modifications" from dashboard
2. Enter commit message
3. Changes automatically staged and committed
4. Return to dashboard

#### Manual Push

```mermaid
graph LR
    A[Dashboard] --> B[Select Sync Upstream]
    B --> C[Push Current Branch]
    C --> D[Push Complete]
    D --> A
```

**Use Case**: Push existing commits without pipeline

**Steps**:
1. Select "Sync Tracked Upstream Modifications" from dashboard
2. Current branch pushed to remote
3. Upstream tracking set automatically
4. Return to dashboard

#### Reset State

```mermaid
graph LR
    A[Dashboard] --> B[Select Reset Engine]
    B --> C[Clear All Context]
    C --> D[Return to Dashboard]
    D --> A
```

**Use Case**: Clear stuck state or start fresh

**Steps**:
1. Select "Reset Context State Engine" from dashboard
2. All context cleared (issue, branch, PR)
3. Return to clean dashboard state

---

## CRUD Operations

EasyFlow provides CRUD (Create, Read, Update, Delete) operations for issues, branches, and commits.

### Issue Management

```mermaid
graph TD
    A[Manage Issues Menu] --> B[List Issues]
    A --> C[Create Issue]
    A --> D[Close Issue]
    
    B --> E[Select Issue]
    E --> F[View Details]
    F --> A
    
    C --> G[Enter Title]
    G --> H[Issue Created]
    H --> A
    
    D --> I[Enter Issue Number]
    I --> J[Issue Closed]
    J --> A
```

#### List Issues

- Displays all open issues from repository
- Shows issue number and title
- Navigate with ↑/↓ or j/k
- Select issue to work on
- Press `n` to create new issue

#### Create Issue

- Enter issue title
- Issue created on GitHub
- Issue number returned
- Auto-advances to branch creation

#### Close Issue

- Enter issue number
- Issue closed on GitHub
- Confirmation displayed
- Returns to dashboard

### Branch Management

```mermaid
graph TD
    A[Manage Branches Menu] --> B[Checkout Branch]
    A --> C[Create Branch]
    A --> D[Delete Branch]
    
    B --> E[List Branches]
    E --> F[Select Branch]
    F --> G[Branch Checked Out]
    G --> A
    
    C --> H[Enter Branch Name]
    H --> I[Branch Created]
    I --> A
    
    D --> J[List Branches]
    J --> K[Select Branch]
    K --> L[Branch Deleted]
    L --> A
```

#### Checkout Branch

- Lists all local branches
- Navigate with ↑/↓ or j/k
- Select branch to checkout
- Branch switched immediately

#### Create Branch

- Enter custom branch name
- Branch sanitized automatically
- Branch created and checked out
- Returns to dashboard

#### Delete Branch

- Lists all local branches
- Navigate with ↑/↓ or j/k
- Select branch to delete
- Safe deletion with verification

### Commit Management

```mermaid
graph TD
    A[Manage Commits Menu] --> B[View Log]
    A --> C[Stage & Commit]
    A --> D[Undo Last Commit]
    
    B --> E[Show Last 5 Commits]
    E --> F[Press ESC to Return]
    F --> A
    
    C --> G[Enter Commit Message]
    G --> H[Changes Staged & Committed]
    H --> A
    
    D --> I[Execute Soft Reset]
    I --> J[Commit Undone]
    J --> A
```

#### View Log

- Displays last 5 commit messages
- Shows commit hash and message
- Press ESC to return to dashboard

#### Stage & Commit

- All changes automatically staged
- Enter commit message
- Commit created locally
- Returns to dashboard

#### Undo Last Commit

- Executes `git reset --soft HEAD~1`
- Changes preserved in working directory
- Commit removed from history
- Returns to dashboard

---

## Usage Examples

### Example 1: Complete Pipeline

**Scenario**: Fix a bug reported in issue #42

```bash
# Run EasyFlow
easyflow

# In the UI:
1. Select "🚀 Start Pipeline Work Loop"
2. Navigate to issue #42: "Fix authentication timeout"
3. Press Enter to select
4. Branch name auto-filled: "issue-42"
5. Press Enter to confirm branch
6. Switch to another terminal to write code
7. Make changes to fix the bug
8. Return to EasyFlow, press Enter
9. Enter commit message: "Fix authentication timeout error"
10. Wait for push to complete
11. PR title auto-filled: "Fix authentication timeout"
12. Press Enter to create PR
13. Review merge actions
14. Press Enter to merge
15. Success! Issue #42 closed automatically
16. Press ESC to return to dashboard
```

### Example 2: Quick Commit

**Scenario**: Quick commit without full pipeline

```bash
# Run EasyFlow
easyflow

# In the UI:
1. Select "💾 Stage & Commit Local Modifications"
2. Enter commit message: "Update documentation"
3. Changes staged and committed
4. Return to dashboard
```

### Example 3: Branch Management

**Scenario**: Create and switch to a feature branch

```bash
# Run EasyFlow
easyflow

# In the UI:
1. Select "🌿 Manage Branches Menu"
2. Select "Create Custom Local Branch"
3. Enter branch name: "feature-user-auth"
4. Branch created and checked out
5. Press ESC to return to dashboard
```

### Example 4: Issue Management

**Scenario**: Create a new issue for a feature request

```bash
# Run EasyFlow
easyflow

# In the UI:
1. Select "🐛 Manage Issues Menu"
2. Select "Create Tracker Issue"
3. Enter issue title: "Add OAuth2 support"
4. Issue created on GitHub
5. Press ESC to return to dashboard
```

### Example 5: View Commit History

**Scenario**: Check recent commits on current branch

```bash
# Run EasyFlow
easyflow

# In the UI:
1. Select "💾 Manage Commits Menu"
2. Select "View Recent Commit Log"
3. View last 5 commits
4. Press ESC to return to dashboard
```

---

## Best Practices

### Workflow Best Practices

1. **Use Pipeline Mode for Complete Features**
   - Ideal for new features or bug fixes
   - Ensures proper issue tracking
   - Automates cleanup (branch deletion, issue closing)

2. **Use Standalone Mode for Quick Tasks**
   - Quick commits without full workflow
   - Manual pushes when needed
   - Reset state when stuck

3. **Create Descriptive Branch Names**
   - Use issue numbers: `issue-123`
   - Use descriptive names: `feature-auth-fix`
   - Avoid special characters

4. **Write Clear Commit Messages**
   - Describe what and why, not how
   - Use present tense: "Fix bug" not "Fixed bug"
   - Keep messages concise but informative

5. **Review Before Merging**
   - Check PR title and body
   - Ensure all changes are committed
   - Verify issue is properly referenced

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate menus |
| `Enter` | Select option / Advance |
| `Esc` | Return to dashboard / Reset |
| `q` or `Ctrl+C` | Quit application |
| `n` | Create new issue (in issue selection) |

### Error Handling

1. **Git Not Found**
   - Install Git: `brew install git` (macOS)
   - Verify installation: `git --version`

2. **GitHub CLI Not Found**
   - Install GitHub CLI: `brew install gh` (macOS)
   - Verify installation: `gh --version`

3. **Not a Git Repository**
   - Navigate to a git repository
   - Initialize if needed: `git init`

4. **No Remote Origin**
   - Add remote: `git remote add origin <url>`
   - Verify remote: `git remote -v`

5. **GitHub Auth Missing**
   - Authenticate: `gh auth login`
   - Verify auth: `gh auth status`

---

## Troubleshooting

### Common Issues

#### Issue: Stuck in Loading State

**Symptoms**: Spinner keeps spinning, no progress

**Solutions**:
- Press `Esc` to reset to dashboard
- Check network connection
- Verify GitHub CLI is authenticated
- Check for GitHub API rate limits

#### Issue: Branch Creation Fails

**Symptoms**: Error when creating branch

**Solutions**:
- Verify branch name is valid
- Check if branch already exists
- Ensure you have write permissions
- Check Git remote configuration

#### Issue: Push Fails

**Symptoms**: Error when pushing to remote

**Solutions**:
- Check network connection
- Verify remote URL is correct
- Ensure you have push permissions
- Check for merge conflicts

#### Issue: PR Creation Fails

**Symptoms**: Error when creating pull request

**Solutions**:
- Verify branch is pushed to remote
- Check PR title is not empty
- Ensure you have PR creation permissions
- Check for existing PR for branch

#### Issue: Merge Fails

**Symptoms**: Error when merging PR

**Solutions**:
- Verify PR exists and is mergeable
- Check for merge conflicts
- Ensure you have merge permissions
- Check branch protection rules

### Debug Mode

To enable debug output, set environment variable:

```bash
export EASYFLOW_DEBUG=1
easyflow
```

This will print detailed error messages and stack traces.

### Getting Help

If you encounter issues not covered here:

1. Check the [Troubleshooting Guide](troubleshooting.md)
2. Review [Architecture Documentation](architecture.md)
3. Check [Module Documentation](modules.md)
4. Open an issue on GitHub

---

## Advanced Usage

### Custom Layout Configuration

Modify `internal/config/config.go` to customize UI layout:

```go
func DefaultLayout() LayoutConfig {
    return LayoutConfig{
        MenuSpacing: 2,  // Tighter spacing
        ColumnWidth: 60, // Wider columns
    }
}
```

### Adding Custom Menu Items

Add new menu options in `internal/ui/menu.go`:

```go
{
    Title:       "Custom Action",
    Description: "Your custom action description",
}
```

Then handle the selection in `internal/ui/update.go`:

```go
case 7: // Custom action index
    // Your custom logic here
```

### Extending Workflow States

Add new states in `internal/workflow/state.go`:

```go
const (
    // ... existing states ...
    StateCustomAction State = iota
)
```

Add validation in `internal/workflow/workflow.go`:

```go
case StateCustomAction:
    // Your validation logic
```

Add UI handling in `internal/ui/update.go`:

```go
case workflow.StateCustomAction:
    // Your UI handling logic
```

Add view rendering in `internal/ui/view.go`:

```go
case workflow.StateCustomAction:
    // Your view rendering logic
```

---

**Related Documentation**:
- [Architecture Overview](architecture.md) - System architecture
- [Module Documentation](modules.md) - Component details
- [API Reference](api.md) - Complete API documentation
- [Troubleshooting](troubleshooting.md) - Common issues and solutions
