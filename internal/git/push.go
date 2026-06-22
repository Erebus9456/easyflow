package git

import (
	"fmt"

	"github.com/Erebus9456/easyflow/utils"
)

// PushToRemote pushes the active branch to origin and establishes upstream tracking
func PushToRemote() (string, error) {
	stdout, stderr, err := utils.ExecuteCommand("git", "push", "-u", "origin", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to push changes to remote: %s %w", stderr, err)
	}
	return stdout, nil
}
