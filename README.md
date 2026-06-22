# EasyFlow 🚀

EasyFlow is a **terminal-first GitHub workflow automation tool** built in Go. It replaces repetitive Git + GitHub UI actions with a **guided interactive terminal application** powered by Bubble Tea.

Instead of manually switching between browser, terminal, and git commands, you run:

```bash
easyflow
```

And complete your entire development workflow inside one UI.

---

## 🧠 Core Idea

A single command drives your entire daily development loop:

```
Issue → Branch → Code → Commit → Push → PR → Merge → Close Issue
```

Everything is interactive, fast, and keyboard-driven.

---

## 🛠 Tech Stack

| Tool | Role |
|------|------|
| **Go** | Core language |
| **Bubble Tea** | UI engine |
| **Bubbles** | Inputs, lists, spinners |
| **Lip Gloss** | UI styling |
| **Cobra** | CLI entry |
| **GitHub CLI (`gh`)** | GitHub API operations |
| **Git** | Local repository operations |

---

## 📁 Project Structure

```
easyflow/
│
├── main.go
├── cmd/
│   └── root.go
│
├── internal/
│   ├── ui/
│   ├── workflow/
│   ├── github/
│   ├── git/
│   └── config/
│
├── utils/
├── scripts/
└── docs/
```

---

## 📌 File-by-File Breakdown

### 🟢 `main.go`

**Purpose:** Application entry point.

**Responsibilities:**
- Starts the CLI
- Initializes Cobra root command
- Launches Bubble Tea UI

**Flow:** `main.go` → `cmd/root.go` → Bubble Tea app

---

### 🟢 `cmd/root.go`

**Purpose:** CLI bootstrap layer.

**Responsibilities:**
- Defines the `easyflow` command
- Starts the Bubble Tea program
- Handles global flags (future)

> This file is the bridge between the terminal command and the UI system.

---

## 🎨 UI Layer (Bubble Tea)

Everything inside `internal/ui/` controls the terminal experience.

### 🟣 `internal/ui/model.go`

**Purpose:** Core application state.

**Stores:**
- Current selected menu option
- Active issue
- Branch name
- Workflow step
- GitHub context

> Think of it as the "brain state" of the UI.

---

### 🟣 `internal/ui/update.go`

**Purpose:** Handles all keyboard input & logic updates.

**Responsibilities:**
- Arrow key navigation
- Enter selection
- State transitions
- Trigger workflow actions

**Example:** User presses `ENTER` → moves from menu → branch creation

---

### 🟣 `internal/ui/view.go`

**Purpose:** Renders the terminal UI.

**Responsibilities:**
- Draws menu
- Displays workflow status
- Shows current step
- Renders success/error messages

**Output:** Pretty terminal dashboard UI

---

### 🟣 `internal/ui/menu.go`

**Purpose:** Defines all available actions.

**Menu options:**
- Start Work
- Create Branch
- Commit Changes
- Push
- Create PR
- Merge PR
- Close Issue

---

### 🟣 `internal/ui/styles.go`

**Purpose:** Visual styling system.

**Responsibilities:**
- Colors
- Layout spacing
- Highlight styles
- Selection cursor styles

**Uses:** Lip Gloss styling engine

---

## 🔄 Workflow Engine

### 🟡 `internal/workflow/workflow.go`

**Purpose:** Core automation logic.

**Responsibilities:**
- Controls full GitHub lifecycle
- Orchestrates Git + GitHub CLI
- Moves step-by-step through the workflow

**Flow:** `Issue → Branch → Commit → Push → PR → Merge → Close`

---

### 🟡 `internal/workflow/state.go`

**Purpose:** Stores workflow runtime state.

**Stores:**
- Active issue ID
- Branch name
- PR URL
- Current step
- Repository context

---

## 🐙 GitHub Integration Layer

Uses GitHub CLI (`gh`).

### 🔵 `internal/github/issues.go`

**Purpose:** Manages GitHub Issues.

**Commands used:** `gh issue list`, `gh issue view`, `gh issue create`, `gh issue close`

**Responsibilities:**
- Fetch issues
- Select issue for work
- Close issue after merge

---

### 🔵 `internal/github/pr.go`

**Purpose:** Manages Pull Requests.

**Commands used:** `gh pr create`, `gh pr merge`, `gh pr view`

**Responsibilities:**
- Create PR from branch
- Merge PR
- Auto-delete branch after merge

---

## 🌿 Git Layer

### 🟠 `internal/git/branch.go`

**Purpose:** Branch management.

**Commands used:** `git checkout -b`, `git branch`

**Responsibilities:**
- Create new branch from issue
- Switch branches

---

### 🟠 `internal/git/commit.go`

**Purpose:** Commit management.

**Commands used:** `git add .`, `git commit -m`

**Responsibilities:**
- Stage changes
- Create commits
- Support multiple commits per issue

---

### 🟠 `internal/git/push.go`

**Purpose:** Push local changes.

**Commands used:** `git push -u origin HEAD`

**Responsibilities:**
- Push current branch
- Set upstream automatically

---

## ⚙️ Configuration Layer

### ⚪ `internal/config/config.go`

**Purpose:** Stores user settings.

**Future options:**
- Default branch name
- Preferred commit style
- GitHub username
- Workflow preferences

---

## 🔧 Utilities

### ⚫ `utils/shell.go`

**Purpose:** Executes system commands.

**Responsibilities:**
- Run `git` commands
- Run `gh` commands
- Capture output
- Handle errors

---

### ⚫ `utils/errors.go`

**Purpose:** Centralized error handling.

**Responsibilities:**
- Standard error format
- Logging failures
- Debug support

---

## 🚀 How the App Works

**Step 1** — Run the command:

```bash
easyflow
```

**Step 2** — Interactive menu appears:

```
> Start Work
  Create Branch
  Commit
  Push
  PR
  Merge
  Close Issue
```

**Step 3** — Select **Start Work**

**Step 4** — Guided flow begins:

```
Select Issue → Create Branch → Commit → Push → PR → Merge → Close
```

---

## 🔮 Roadmap

| Version | Features |
|---------|----------|
| **v1** | Full workflow automation, GitHub CLI integration, interactive UI |
| **v2** | Multi-repo support, config profiles, keyboard shortcuts |
| **v3** | AI commit messages, AI PR descriptions, sprint planning mode |

---

## 🎯 Goal

Replace:
- GitHub Web UI
- Manual git commands
- Context switching

With **a single terminal application** that runs your entire dev workflow.