package utils

import "errors"

var (
	// ErrGitNotFound indicates git is missing from the system PATH
	ErrGitNotFound = errors.New("git binary not found on your system PATH")

	// ErrGitHubCLINotFound indicates gh is missing from the system PATH
	ErrGitHubCLINotFound = errors.New("github cli (gh) binary not found on your system PATH")

	// ErrNotAGitRepository implies easyflow was executed outside a tracking environment
	ErrNotAGitRepository = errors.New("not a git repository (or any of the parent directories)")

	// ErrNoRemoteOrigin means git is initialized, but no upstream remote exists
	ErrNoRemoteOrigin = errors.New("git remote 'origin' is missing or not configured")

	// ErrGitHubAuthMissing means the gh CLI tool is logged out
	ErrGitHubAuthMissing = errors.New("gh cli is not authenticated. run 'gh auth login' first")
)
