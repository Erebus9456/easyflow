# API Reference

This document provides a complete API reference for all public functions, types, and constants in the EasyFlow codebase.

## Table of Contents

- [UI Module (`internal/ui`)](#ui-module-internalui)
- [Workflow Module (`internal/workflow`)](#workflow-module-internalworkflow)
- [GitHub Module (`internal/github`)](#github-module-internalgithub)
- [Git Module (`internal/git`)](#git-module-internalgit)
- [Config Module (`internal/config`)](#config-module-internalconfig)
- [Utilities Module (`utils`)](#utilities-module-utils)

---

## UI Module (`internal/ui`)

### Types

#### `AppModel`

Application state model for Bubble Tea.

```go
type AppModel struct {
    Engine      *workflow.Engine      // Workflow orchestration engine
    RepoCtx     *git.RepoContext     // Repository context information
    MenuItems   []MainMenuItem       // Available menu options
    Cursor      int                   // Current menu selection index
    Issues      []github.Issue        // Fetched issues list
    IssueCursor int                   // Current issue selection index
    Layout      config.LayoutConfig   // UI layout configuration
    TextInput   textinput.Model       // Text input component
    Spinner     spinner.Model         // Loading spinner component
    Loading     bool                  // Loading state flag
    ErrorMessage string               // Error message to display
    SuccessMsg   string               // Success message to display
}
```

#### `MainMenuItem`

Menu item structure.

```go
type MainMenuItem struct {
    Title       string  // Menu item title
    Description string  // Menu item description
}
```

### Functions

#### `InitialModel`

Creates initial application state.

```go
func InitialModel(repo *git.RepoContext) AppModel
```

**Parameters**:
- `repo`: Repository context from `git.GetCurrentRepo()`

**Returns**:
- `AppModel`: Initialized application model

**Example**:
```go
repoCtx, err := git.GetCurrentRepo()
model := ui.InitialModel(repoCtx)
```

#### `GetMainMenuOptions`

Returns main menu items.

```go
func GetMainMenuOptions() []MainMenuItem
```

**Returns**:
- `[]MainMenuItem`: List of main menu options

#### `GetSubMenuOptions`

Returns submenu items for a category.

```go
func GetSubMenuOptions(category string) []MainMenuItem
```

**Parameters**:
- `category`: One of "issues", "branches", "commits"

**Returns**:
- `[]MainMenuItem`: List of submenu options

**Example**:
```go
issuesMenu := GetSubMenuOptions("issues")
branchesMenu := GetSubMenuOptions("branches")
commitsMenu := GetSubMenuOptions("commits")
```

### Methods

#### `Init`

Initializes Bubble Tea subcomponents.

```go
func (m AppModel) Init() tea.Cmd
```

**Returns**:
- `tea.Cmd`: Initialization command for text input and spinner

#### `Update`

Main update loop for handling events.

```go
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)
```

**Parameters**:
- `msg`: Bubble Tea message (key press, async result, etc.)

**Returns**:
- `tea.Model`: Updated model
- `tea.Cmd`: Command to execute

#### `View`

Renders the terminal UI.

```go
func (m AppModel) View() string
```

**Returns**:
- `string`: Rendered UI string

---

## Workflow Module (`internal/workflow`)

### Types

#### `Engine`

Workflow orchestration engine.

```go
type Engine struct {
    Ctx *RuntimeContext  // Runtime context
}
```

#### `RuntimeContext`

Runtime context for workflow state.

```go
type RuntimeContext struct {
    ActiveIssueNumber int     // Currently selected issue ID
    ActiveIssueTitle  string  // Currently selected issue title
    BranchName        string  // Current working branch
    PullRequestURL    string  // Created PR URL
    CurrentStep       State   // Current workflow state
    PipelineMode      bool    // Continuous pipeline mode flag
}
```

#### `State`

Workflow state enumeration.

```go
type State int
```

**Constants**:
```go
const (
    StateDashboard      State = iota  // Main menu
    StateSelectIssue                  // Issue selection
    StateCreateIssue                  // Issue creation
    StateCreateBranch                 // Branch creation
    StateWorking                      // Code editing
    StateCommitReady                  // Commit preparation
    StatePushing                      // Pushing to remote
    StatePRPending                    // PR creation
    StateMerging                      // Merge authorization
    StateCompleted                    // Pipeline complete
    StateManageIssues                 // Issue management submenu
    StateManageBranches              // Branch management submenu
    StateManageCommits               // Commit management submenu
    StateListBranches                // Branch list display
    StateViewCommits                 // Commit log display
)
```

### Functions

#### `NewEngine`

Creates new workflow engine.

```go
func NewEngine() *Engine
```

**Returns**:
- `*Engine`: Initialized workflow engine

**Example**:
```go
engine := workflow.NewEngine()
```

#### `NewRuntimeContext`

Creates default runtime context.

```go
func NewRuntimeContext() *RuntimeContext
```

**Returns**:
- `*RuntimeContext`: Initialized runtime context

### Methods

#### `Advance`

Transitions to next state with validation.

```go
func (e *Engine) Advance(next State) error
```

**Parameters**:
- `next`: Target state to transition to

**Returns**:
- `error`: Validation error if transition is invalid

**Validation Rules**:
- `StateCreateBranch`: Requires `ActiveIssueNumber` to be set
- `StateCommitReady`: Requires `BranchName` to be set
- `StatePRPending`: Requires `BranchName` to be set
- `StateMerging`: Requires `PullRequestURL` to be set

**Example**:
```go
err := engine.Advance(workflow.StateCreateBranch)
if err != nil {
    return fmt.Errorf("invalid transition: %w", err)
}
```

#### `Reset`

Clears context and returns to dashboard.

```go
func (e *Engine) Reset()
```

**Example**:
```go
engine.Reset()
```

---

## GitHub Module (`internal/github`)

### Types

#### `Issue`

GitHub issue representation.

```go
type Issue struct {
    Number int    `json:"number"`  // Issue number
    Title  string `json:"title"`   // Issue title
    Body   string `json:"body"`    // Issue body
}
```

### Functions

#### `FetchOpenIssues`

Fetches open issues from repository.

```go
func FetchOpenIssues() ([]Issue, error)
```

**Returns**:
- `[]Issue`: List of open issues
- `error`: Error if fetch fails

**Example**:
```go
issues, err := github.FetchOpenIssues()
if err != nil {
    return fmt.Errorf("failed to fetch issues: %w", err)
}
for _, issue := range issues {
    fmt.Printf("#%d: %s\n", issue.Number, issue.Title)
}
```

#### `CreateIssue`

Creates a new issue on GitHub.

```go
func CreateIssue(title string) (int, error)
```

**Parameters**:
- `title`: Issue title (required, non-empty)

**Returns**:
- `int`: Created issue number
- `error`: Error if creation fails

**Example**:
```go
issueNum, err := github.CreateIssue("Fix authentication bug")
if err != nil {
    return fmt.Errorf("failed to create issue: %w", err)
}
fmt.Printf("Created issue #%d\n", issueNum)
```

#### `CloseIssue`

Closes an issue by number.

```go
func CloseIssue(number int) error
```

**Parameters**:
- `number`: Issue number to close

**Returns**:
- `error`: Error if close fails

**Example**:
```go
err := github.CloseIssue(42)
if err != nil {
    return fmt.Errorf("failed to close issue: %w", err)
}
```

#### `CreatePullRequest`

Creates a pull request.

```go
func CreatePullRequest(title, body string) (string, error)
```

**Parameters**:
- `title`: PR title (required, non-empty)
- `body`: PR description body

**Returns**:
- `string`: PR URL
- `error`: Error if creation fails

**Example**:
```go
url, err := github.CreatePullRequest("Fix auth bug", "Closes #123")
if err != nil {
    return fmt.Errorf("failed to create PR: %w", err)
}
fmt.Printf("PR created: %s\n", url)
```

#### `MergeAndCleanupPR`

Merges PR and deletes branch.

```go
func MergeAndCleanupPR() (string, error)
```

**Returns**:
- `string`: Merge output
- `error`: Error if merge fails

**Behavior**:
- Uses merge commit strategy
- Automatically deletes remote branch after merge

**Example**:
```go
output, err := github.MergeAndCleanupPR()
if err != nil {
    return fmt.Errorf("failed to merge PR: %w", err)
}
fmt.Println("PR merged successfully")
```

---

## Git Module (`internal/git`)

### Types

#### `RepoContext`

Repository context information.

```go
type RepoContext struct {
    Owner          string  // Repository owner
    RepositoryName string  // Repository name
    CurrentBranch  string  // Current branch name
}
```

### Functions

#### `GetCurrentRepo`

Gets current repository context.

```go
func GetCurrentRepo() (*RepoContext, error)
```

**Returns**:
- `*RepoContext`: Repository context
- `error`: Error if not in a git repository

**Example**:
```go
repoCtx, err := git.GetCurrentRepo()
if err != nil {
    return fmt.Errorf("not in a git repository: %w", err)
}
fmt.Printf("Repo: %s/%s\n", repoCtx.Owner, repoCtx.RepositoryName)
```

#### `SanitizeBranchName`

Cleans branch name for Git.

```go
func SanitizeBranchName(input string) string
```

**Parameters**:
- `input`: Raw branch name

**Returns**:
- `string`: Sanitized branch name

**Sanitization Rules**:
- Converts to lowercase
- Replaces spaces with hyphens
- Removes invalid characters: `*`, `?`, `~`, `^`, `:`, `\`

**Example**:
```go
clean := git.SanitizeBranchName("Feature/New Auth")
// Returns: "feature/new-auth"
```

#### `CreateAndCheckoutBranch`

Creates and switches to branch.

```go
func CreateAndCheckoutBranch(name string) (string, error)
```

**Parameters**:
- `name`: Branch name

**Returns**:
- `string`: Command output
- `error`: Error if creation fails

**Example**:
```go
output, err := git.CreateAndCheckoutBranch("feature-auth")
if err != nil {
    return fmt.Errorf("failed to create branch: %w", err)
}
```

#### `ListLocalBranches`

Lists all local branches.

```go
func ListLocalBranches() ([]string, error)
```

**Returns**:
- `[]string`: List of branch names
- `error`: Error if list fails

**Example**:
```go
branches, err := git.ListLocalBranches()
if err != nil {
    return fmt.Errorf("failed to list branches: %w", err)
}
for _, branch := range branches {
    fmt.Println(branch)
}
```

#### `CheckoutBranch`

Switches to existing branch.

```go
func CheckoutBranch(name string) error
```

**Parameters**:
- `name`: Branch name to checkout

**Returns**:
- `error`: Error if checkout fails

**Example**:
```go
err := git.CheckoutBranch("main")
if err != nil {
    return fmt.Errorf("failed to checkout: %w", err)
}
```

#### `DeleteLocalBranch`

Deletes a local branch.

```go
func DeleteLocalBranch(name string) error
```

**Parameters**:
- `name`: Branch name to delete

**Returns**:
- `error`: Error if deletion fails

**Behavior**:
- Uses safe deletion flag (`-d`)
- Prevents deletion if branch has unmerged changes

**Example**:
```go
err := git.DeleteLocalBranch("old-feature")
if err != nil {
    return fmt.Errorf("failed to delete branch: %w", err)
}
```

#### `StageAllChanges`

Stages all changes.

```go
func StageAllChanges() error
```

**Returns**:
- `error`: Error if staging fails

**Behavior**:
- Executes `git add .`
- Stages all tracked and untracked files

**Example**:
```go
err := git.StageAllChanges()
if err != nil {
    return fmt.Errorf("failed to stage: %w", err)
}
```

#### `CreateCommit`

Creates commit with message.

```go
func CreateCommit(message string) (string, error)
```

**Parameters**:
- `message`: Commit message (required, non-empty)

**Returns**:
- `string`: Command output
- `error`: Error if commit fails

**Example**:
```go
output, err := git.CreateCommit("Fix authentication bug")
if err != nil {
    return fmt.Errorf("failed to commit: %w", err)
}
```

#### `GetLocalCommitLog`

Gets recent commit history.

```go
func GetLocalCommitLog(count int) ([]string, error)
```

**Parameters**:
- `count`: Number of commits to retrieve

**Returns**:
- `[]string`: List of commit messages
- `error`: Error if log fails

**Example**:
```go
logs, err := git.GetLocalCommitLog(5)
if err != nil {
    return fmt.Errorf("failed to get log: %w", err)
}
for _, log := range logs {
    fmt.Println(log)
}
```

#### `UndoLastCommit`

Undoes last commit (soft reset).

```go
func UndoLastCommit() error
```

**Returns**:
- `error`: Error if undo fails

**Behavior**:
- Executes `git reset --soft HEAD~1`
- Preserves local changes in working directory

**Example**:
```go
err := git.UndoLastCommit()
if err != nil {
    return fmt.Errorf("failed to undo: %w", err)
}
```

#### `PushToRemote`

Pushes current branch to remote.

```go
func PushToRemote() (string, error)
```

**Returns**:
- `string`: Command output
- `error`: Error if push fails

**Behavior**:
- Pushes current branch to origin
- Sets upstream tracking automatically

**Example**:
```go
output, err := git.PushToRemote()
if err != nil {
    return fmt.Errorf("failed to push: %w", err)
}
fmt.Println("Push successful")
```

---

## Config Module (`internal/config`)

### Types

#### `Config`

Application configuration.

```go
type Config struct {
    DefaultBranch string  // Default branch name (future use)
}
```

#### `LayoutConfig`

UI layout configuration.

```go
type LayoutConfig struct {
    MenuSpacing int  // Number of newlines between menu items
    ColumnWidth int  // Width of UI columns
}
```

### Functions

#### `DefaultLayout`

Returns default layout configuration.

```go
func DefaultLayout() LayoutConfig
```

**Returns**:
- `LayoutConfig`: Default layout settings

**Default Values**:
- `MenuSpacing`: 4 (relaxed spacing)
- `ColumnWidth`: 50 (standard column width)

**Example**:
```go
layout := config.DefaultLayout()
fmt.Printf("Menu spacing: %d\n", layout.MenuSpacing)
fmt.Printf("Column width: %d\n", layout.ColumnWidth)
```

---

## Utilities Module (`utils`)

### Errors

#### Standard Errors

```go
var (
    ErrGitNotFound         = errors.New("git binary not found on your system PATH")
    ErrGitHubCLINotFound  = errors.New("github cli (gh) binary not found on your system PATH")
    ErrNotAGitRepository   = errors.New("not a git repository (or any of the parent directories)")
    ErrNoRemoteOrigin      = errors.New("git remote 'origin' is missing or not configured")
    ErrGitHubAuthMissing   = errors.New("gh cli is not authenticated. run 'gh auth login' first")
)
```

### Functions

#### `ExecuteCommand`

Executes shell command with output capture.

```go
func ExecuteCommand(name string, args ...string) (string, string, error)
```

**Parameters**:
- `name`: Command name (e.g., "git", "gh")
- `args`: Command arguments

**Returns**:
- `string`: Standard output (trimmed)
- `string`: Standard error (trimmed)
- `error`: Execution error if any

**Example**:
```go
stdout, stderr, err := utils.ExecuteCommand("git", "branch", "--list")
if err != nil {
    fmt.Printf("Error: %s\n", stderr)
    return err
}
fmt.Println("Branches:", stdout)
```

---

## Usage Examples

### Complete Pipeline Example

```go
package main

import (
    "fmt"
    "github.com/Erebus9456/easyflow/internal/git"
    "github.com/Erebus9456/easyflow/internal/github"
    "github.com/Erebus9456/easyflow/internal/workflow"
)

func main() {
    // Initialize workflow engine
    engine := workflow.NewEngine()
    
    // Get repository context
    repoCtx, err := git.GetCurrentRepo()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // Fetch issues
    issues, err := github.FetchOpenIssues()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // Select first issue
    if len(issues) > 0 {
        issue := issues[0]
        engine.Ctx.ActiveIssueNumber = issue.Number
        engine.Ctx.ActiveIssueTitle = issue.Title
        
        // Advance to branch creation
        err = engine.Advance(workflow.StateCreateBranch)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }
        
        // Create branch
        branchName := fmt.Sprintf("issue-%d", issue.Number)
        _, err = git.CreateAndCheckoutBranch(branchName)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }
        
        fmt.Printf("Created branch: %s\n", branchName)
    }
}
```

### Git Operations Example

```go
package main

import (
    "fmt"
    "github.com/Erebus9456/easyflow/internal/git"
)

func main() {
    // Stage all changes
    err := git.StageAllChanges()
    if err != nil {
        fmt.Printf("Error staging: %v\n", err)
        return
    }
    
    // Create commit
    output, err := git.CreateCommit("Add new feature")
    if err != nil {
        fmt.Printf("Error committing: %v\n", err)
        return
    }
    fmt.Println(output)
    
    // Push to remote
    output, err = git.PushToRemote()
    if err != nil {
        fmt.Printf("Error pushing: %v\n", err)
        return
    }
    fmt.Println(output)
}
```

### GitHub Operations Example

```go
package main

import (
    "fmt"
    "github.com/Erebus9456/easyflow/internal/github"
)

func main() {
    // Create issue
    issueNum, err := github.CreateIssue("Fix authentication bug")
    if err != nil {
        fmt.Printf("Error creating issue: %v\n", err)
        return
    }
    fmt.Printf("Created issue #%d\n", issueNum)
    
    // Create PR
    url, err := github.CreatePullRequest("Fix auth bug", "Closes #123")
    if err != nil {
        fmt.Printf("Error creating PR: %v\n", err)
        return
    }
    fmt.Printf("PR created: %s\n", url)
    
    // Merge PR
    output, err := github.MergeAndCleanupPR()
    if err != nil {
        fmt.Printf("Error merging PR: %v\n", err)
        return
    }
    fmt.Println(output)
    
    // Close issue
    err = github.CloseIssue(issueNum)
    if err != nil {
        fmt.Printf("Error closing issue: %v\n", err)
        return
    }
    fmt.Println("Issue closed")
}
```

---

**Related Documentation**:
- [Module Documentation](modules.md) - Detailed module documentation
- [Architecture Overview](architecture.md) - System architecture
- [Workflow Guide](workflow.md) - Workflow automation guide
