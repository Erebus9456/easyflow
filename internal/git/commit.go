package git

import (
	"fmt"
	"strings"

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

func GetLocalCommitLog(count int) ([]string, error) {
	countStr := fmt.Sprintf("-n %d", count)
	stdout, _, err := utils.ExecuteCommand("git", "log", countStr, "--oneline")
	if err != nil {
		return nil, err
	}

	rawLines := strings.Split(stdout, "\n")
	var logs []string
	for _, line := range rawLines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			logs = append(logs, trimmed)
		}
	}
	return logs, nil
}

// UndoLastCommit executes a soft reset back one revision, saving local changes intact
func UndoLastCommit() error {
	_, _, err := utils.ExecuteCommand("git", "reset", "--soft", "HEAD~1")
	return err
}
