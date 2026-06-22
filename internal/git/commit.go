package git

import (
	"fmt"

	"github.com/Erebus9456/easyflow/utils"
)

// StageAllChanges executes a blanket git add on the project root
func StageAllChanges() error {
	_, stderr, err := utils.ExecuteCommand("git", "add", ".")
	if err != nil {
		return fmt.Errorf("failed to stage changes: %s %w", stderr, err)
	}
	return nil
}

// CreateCommit commits staged items with a designated user summary message
func CreateCommit(message string) (string, error) {
	if message == "" {
		return "", fmt.Errorf("commit message cannot be empty")
	}

	stdout, stderr, err := utils.ExecuteCommand("git", "commit", "-m", message)
	if err != nil {
		return "", fmt.Errorf("failed to execute commit: %s %w", stderr, err)
	}

	return stdout, nil
}
