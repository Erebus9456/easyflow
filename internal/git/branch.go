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

// ListLocalBranches parses local branches for the selector panel
func ListLocalBranches() ([]string, error) {
	stdout, _, err := utils.ExecuteCommand("git", "branch", "--format=%(refname:short)")
	if err != nil {
		return nil, err
	}

	rawLines := strings.Split(stdout, "\n")
	var branches []string
	for _, line := range rawLines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			branches = append(branches, trimmed)
		}
	}
	return branches, nil
}

// CheckoutBranch attempts to jump directly onto an existing branch target
func CheckoutBranch(name string) error {
	if name == "" {
		return fmt.Errorf("target branch name cannot be blank")
	}
	_, _, err := utils.ExecuteCommand("git", "checkout", name)
	return err
}

// DeleteLocalBranch safely deletes a branch using the standard '-d' verification flag
func DeleteLocalBranch(name string) error {
	if name == "" {
		return fmt.Errorf("target branch name cannot be blank")
	}
	_, _, err := utils.ExecuteCommand("git", "branch", "-d", name)
	return err
}
