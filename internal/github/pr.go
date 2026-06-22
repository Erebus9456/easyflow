package github

import (
	"fmt"

	"github.com/Erebus9456/easyflow/utils"
)

// CreatePullRequest opens a new remote PR using the head branch context
func CreatePullRequest(title, body string) (string, error) {
	if title == "" {
		return "", fmt.Errorf("pull request title cannot be empty")
	}

	// Executes pull request creation targeting the main branch default
	stdout, stderr, err := utils.ExecuteCommand("gh", "pr", "create", "--title", title, "--body", body)
	if err != nil {
		return "", fmt.Errorf("failed to create pull request: %s %w", stderr, err)
	}

	return stdout, nil // Returns the generated PR URL
}

// MergeAndCleanupPR finalizes a merge transaction and destroys the upstream remote branch branch tracking reference
func MergeAndCleanupPR() (string, error) {
	// Merges the PR targeting the current branch context automatically
	stdout, stderr, err := utils.ExecuteCommand("gh", "pr", "merge", "--merge", "--delete-branch")
	if err != nil {
		return "", fmt.Errorf("failed to merge pull request: %s %w", stderr, err)
	}

	return stdout, nil
}
