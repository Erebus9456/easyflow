package git

import (
	"fmt"
	"strings"

	"github.com/Erebus9456/easyflow/utils"
)

// SanitizeBranchName cleans up a user string to be valid for git
func SanitizeBranchName(input string) string {
	clean := strings.ToLower(input)
	clean = strings.ReplaceAll(clean, " ", "-")
	// Strip basic problematic characters
	for _, char := range []string{"*", "?", "~", "^", ":", "\\"} {
		clean = strings.ReplaceAll(clean, char, "")
	}
	return clean
}

// CreateAndCheckoutBranch generates a local git branch and switches to it
func CreateAndCheckoutBranch(name string) (string, error) {
	safeName := SanitizeBranchName(name)

	// Execute local branch creation checkout
	stdout, stderr, err := utils.ExecuteCommand("git", "checkout", "-b", safeName)
	if err != nil {
		return "", fmt.Errorf("failed to create branch %s: %s %w", safeName, stderr, err)
	}

	return stdout, nil
}
