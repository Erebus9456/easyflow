# Development Guide

This guide covers contributing to EasyFlow, including development setup, coding standards, and contribution guidelines.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Building](#building)
- [Contributing](#contributing)
- [Release Process](#release-process)

---

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- GitHub CLI (gh)
- Basic understanding of Go and terminal UI development

### Fork and Clone

```bash
# Fork the repository on GitHub
# Clone your fork
git clone https://github.com/your-username/easyflow.git
cd easyflow

# Add upstream remote
git remote add upstream https://github.com/Erebus9456/easyflow.git
```

---

## Development Setup

### Install Dependencies

```bash
# Install Go dependencies
go mod download

# Verify dependencies
go mod verify
```

### Build the Project

```bash
# Build the binary
go build -o easyflow

# Run the binary
./easyflow
```

### Development Workflow

```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Make your changes
# ... edit files ...

# Test your changes
go build -o easyflow
./easyflow

# Commit your changes
git add .
git commit -m "Add your feature description"

# Push to your fork
git push origin feature/your-feature-name
```

---

## Project Structure

```
easyflow/
├── main.go                  # Application entry point
├── cmd/
│   └── root.go              # CLI command definition
├── internal/
│   ├── ui/                  # UI layer
│   │   ├── model.go         # Application state
│   │   ├── update.go        # Event handling
│   │   ├── view.go          # UI rendering
│   │   ├── menu.go          # Menu definitions
│   │   └── styles.go        # Visual styling
│   ├── workflow/            # Workflow orchestration
│   │   ├── workflow.go      # Workflow engine
│   │   └── state.go         # State definitions
│   ├── github/              # GitHub integration
│   │   ├── issues.go        # Issue operations
│   │   └── pr.go            # PR operations
│   ├── git/                 # Git operations
│   │   ├── branch.go        # Branch management
│   │   ├── commit.go        # Commit operations
│   │   └── push.go          # Push operations
│   └── config/              # Configuration
│       └── config.go        # Layout settings
├── utils/                   # Utilities
│   ├── shell.go             # Command execution
│   └── errors.go            # Error definitions
├── docs/                    # Documentation
├── go.mod                   # Go module file
├── go.sum                   # Go dependencies
└── README.md                # Project README
```

### Module Responsibilities

- **`internal/ui`**: Terminal UI using Bubble Tea
- **`internal/workflow`**: State machine and workflow orchestration
- **`internal/github`**: GitHub API integration via CLI
- **`internal/git`**: Git repository operations
- **`internal/config`**: Configuration and layout settings
- **`utils`**: Shared utilities (shell, errors)

---

## Coding Standards

### Go Conventions

Follow standard Go conventions as defined in [Effective Go](https://golang.org/doc/effective_go).

#### Naming

- Use `camelCase` for variables and functions
- Use `PascalCase` for exported types and functions
- Use `UPPER_SNAKE_CASE` for constants
- Use descriptive names that convey purpose

```go
// Good
func CreateBranch(name string) error {
    // ...
}

// Bad
func cb(n string) error {
    // ...
}
```

#### Package Names

- Use short, lowercase package names
- Avoid package name repetition in identifiers
- Use descriptive package names

```go
// Good
package github

func FetchIssues() ([]Issue, error) {
    // ...
}

// Bad
package github

func GitHubFetchIssues() ([]GitHubIssue, error) {
    // ...
}
```

#### Error Handling

- Always handle errors explicitly
- Use descriptive error messages
- Wrap errors with context

```go
// Good
func CreateBranch(name string) error {
    if name == "" {
        return fmt.Errorf("branch name cannot be empty")
    }
    _, err := utils.ExecuteCommand("git", "checkout", "-b", name)
    if err != nil {
        return fmt.Errorf("failed to create branch %s: %w", name, err)
    }
    return nil
}

// Bad
func CreateBranch(name string) error {
    utils.ExecuteCommand("git", "checkout", "-b", name)
    return nil
}
```

### Code Organization

#### File Organization

- Keep related functions in the same file
- Separate concerns into different files
- Use clear file names that describe contents

#### Function Length

- Keep functions focused and concise
- Aim for functions under 50 lines
- Extract complex logic into helper functions

#### Comments

- Comment exported functions and types
- Comment complex logic
- Use godoc format for documentation

```go
// FetchOpenIssues retrieves the latest 30 open issues from the current context repository.
// It uses the GitHub CLI to fetch issues and returns them as a slice of Issue structs.
// Returns an error if the GitHub CLI command fails or if the response cannot be parsed.
func FetchOpenIssues() ([]Issue, error) {
    // Implementation
}
```

### Bubble Tea Patterns

#### Model-View-Update (MVU)

Follow the MVU pattern consistently:

```go
// Model holds state
type AppModel struct {
    // State fields
}

// Update handles events
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle message
    return m, cmd
}

// View renders UI
func (m AppModel) View() string {
    // Render UI
    return "UI string"
}
```

#### Async Operations

Use Bubble Tea's async message pattern:

```go
// Define message type
type issuesMsg []github.Issue

// Return command that sends message
return func() tea.Msg {
    issues, err := github.FetchOpenIssues()
    if err != nil {
        return errMsg(err)
    }
    return issuesMsg(issues)
}

// Handle message in Update
case issuesMsg:
    m.Issues = msg
    m.Loading = false
    return m, nil
```

### Git and GitHub Operations

#### Command Execution

Always use the `utils.ExecuteCommand` wrapper:

```go
stdout, stderr, err := utils.ExecuteCommand("git", "branch", "--list")
if err != nil {
    return fmt.Errorf("failed to list branches: %s %w", stderr, err)
}
```

#### Error Handling

Handle Git and GitHub CLI errors gracefully:

```go
if err != nil {
    if strings.Contains(stderr, "not a git repository") {
        return utils.ErrNotAGitRepository
    }
    return fmt.Errorf("git error: %s %w", stderr, err)
}
```

---

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./internal/ui
```

### Writing Tests

Follow standard Go testing conventions:

```go
package git

import (
    "testing"
)

func TestSanitizeBranchName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "simple name",
            input:    "feature-auth",
            expected: "feature-auth",
        },
        {
            name:     "with spaces",
            input:    "feature auth",
            expected: "feature-auth",
        },
        {
            name:     "with special chars",
            input:    "feature*auth",
            expected: "featureauth",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := SanitizeBranchName(tt.input)
            if result != tt.expected {
                t.Errorf("SanitizeBranchName(%q) = %q, want %q", tt.input, result, tt.expected)
            }
        })
    }
}
```

### Test Coverage

Aim for >70% code coverage:

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out
```

---

## Building

### Development Build

```bash
# Build for current platform
go build -o easyflow

# Build with debug info
go build -gcflags="all=-N -l" -o easyflow
```

### Production Build

```bash
# Build for multiple platforms
GOOS=darwin GOARCH=amd64 go build -o easyflow-darwin-amd64
GOOS=linux GOARCH=amd64 go build -o easyflow-linux-amd64
GOOS=windows GOARCH=amd64 go build -o easyflow-windows-amd64.exe

# Build with optimizations
go build -ldflags="-s -w" -o easyflow
```

### Cross-Compilation

```bash
# macOS
GOOS=darwin GOARCH=amd64 go build -o easyflow-mac

# Linux
GOOS=linux GOARCH=amd64 go build -o easyflow-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o easyflow.exe
```

---

## Contributing

### Contribution Guidelines

1. **Fork the repository**
2. **Create a feature branch**
3. **Make your changes**
4. **Test thoroughly**
5. **Update documentation**
6. **Submit a pull request**

### Pull Request Process

1. **Describe your changes** in the PR description
2. **Link related issues** if applicable
3. **Ensure tests pass**
4. **Update documentation** if needed
5. **Request review** from maintainers

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] Tests added/updated for new features
- [ ] Documentation updated
- [ ] All tests pass
- [ ] No linting errors

### Commit Messages

Follow conventional commit format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build process or auxiliary tool changes

Examples:
```
feat(ui): add dark theme support

Add dark theme option to configuration and update
color palette for better visibility in low-light environments.

Closes #123
```

```
fix(git): handle branch name sanitization edge cases

Fix issue where special characters in branch names
were not properly sanitized, causing git commands to fail.

Fixes #456
```

---

## Documentation

### Updating Documentation

When making changes, update relevant documentation:

- **API changes**: Update `docs/api.md`
- **New features**: Update `docs/workflow.md` and `docs/quickstart.md`
- **Configuration changes**: Update `docs/configuration.md`
- **Architecture changes**: Update `docs/architecture.md`

### Documentation Style

- Use clear, concise language
- Include code examples
- Use mermaid diagrams for visual explanations
- Keep documentation up to date with code

---

## Debugging

### Debug Mode

Enable debug mode for detailed output:

```bash
export EASYFLOW_DEBUG=1
./easyflow
```

### Using Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug with delve
dlv debug ./main.go
```

### Logging

Add temporary logging for debugging:

```go
fmt.Printf("DEBUG: Current state: %v\n", m.Engine.Ctx.CurrentStep)
fmt.Printf("DEBUG: Issue cursor: %d\n", m.IssueCursor)
```

---

## Release Process

### Versioning

EasyFlow follows semantic versioning (SemVer):

- `MAJOR.MINOR.PATCH`
- Increment MAJOR for breaking changes
- Increment MINOR for new features
- Increment PATCH for bug fixes

### Release Checklist

- [ ] Update version in code
- [ ] Update CHANGELOG.md
- [ ] Tag release
- [ ] Build release binaries
- [ ] Create GitHub release
- [ ] Update documentation

### Creating a Release

```bash
# Update version
# Update CHANGELOG.md

# Commit changes
git add .
git commit -m "chore: release v1.0.0"

# Tag release
git tag -a v1.0.0 -m "Release v1.0.0"

# Push tag
git push origin v1.0.0

# Build release binaries
# Create GitHub release with binaries
```

---

## Code Review Guidelines

### Reviewing Code

When reviewing PRs:

1. **Check functionality**: Does it work as intended?
2. **Check style**: Does it follow project conventions?
3. **Check tests**: Are tests adequate?
4. **Check documentation**: Is documentation updated?
5. **Check for edge cases**: Are edge cases handled?

### Providing Feedback

- Be constructive and specific
- Explain the "why" behind suggestions
- Offer solutions, not just problems
- Be respectful and collaborative

---

## Performance Considerations

### UI Performance

- Minimize view rendering complexity
- Use efficient string building
- Avoid unnecessary allocations
- Cache expensive computations

### Command Execution

- Minimize shell command calls
- Use async operations for long-running commands
- Handle timeouts appropriately
- Cache results when possible

### Memory Management

- Avoid memory leaks in long-running processes
- Clean up resources properly
- Use efficient data structures
- Monitor memory usage

---

## Security Considerations

### Command Injection

Always use parameterized command execution:

```go
// Good
utils.ExecuteCommand("git", "checkout", branchName)

// Bad (vulnerable to injection)
cmd := fmt.Sprintf("git checkout %s", branchName)
utils.ExecuteCommand("sh", "-c", cmd)
```

### Error Messages

Avoid exposing sensitive information in error messages:

```go
// Good
return fmt.Errorf("failed to execute command")

// Bad (may expose sensitive data)
return fmt.Errorf("failed: %s", stderr)
```

### Authentication

Never store credentials in code:
- Use environment variables
- Use system keychain
- Use GitHub CLI authentication

---

## Resources

### Go Resources

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Doc](https://golang.org/doc/)

### Bubble Tea Resources

- [Bubble Tea GitHub](https://github.com/charmbracelet/bubbletea)
- [Bubble Tea Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)
- [Lip Gloss Documentation](https://github.com/charmbracelet/lipgloss)

### GitHub CLI Resources

- [GitHub CLI Manual](https://cli.github.com/manual/)
- [GitHub CLI API](https://github.com/cli/cli/blob/trunk/docs/api.md)

---

## Getting Help

If you need help with development:

1. Check existing documentation
2. Review similar code in the codebase
3. Ask questions in GitHub Discussions
4. Open an issue for bugs or feature requests

---

**Related Documentation**:
- [Architecture Overview](architecture.md) - System design
- [API Reference](api.md) - Complete API documentation
- [Configuration](configuration.md) - Configuration options
