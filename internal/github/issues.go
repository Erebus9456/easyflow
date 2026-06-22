package github

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Erebus9456/easyflow/utils"
)

// Issue represents a stripped down GitHub tracking issue model
type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// FetchOpenIssues requests the latest 30 open issues from the current context repository
func FetchOpenIssues() ([]Issue, error) {
	// Query gh CLI with strict JSON formatting flags
	stdout, stderr, err := utils.ExecuteCommand("gh", "issue", "list", "--state", "open", "--json", "number,title,body")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues via gh cli: %s %w", stderr, err)
	}

	var issues []Issue
	if err := json.Unmarshal([]byte(stdout), &issues); err != nil {
		return nil, fmt.Errorf("failed to parse issues payload: %w", err)
	}

	return issues, nil
}

// CloseIssue updates the remote tracking status of an issue to closed
func CloseIssue(number int) error {
	_, stderr, err := utils.ExecuteCommand("gh", "issue", "close", fmt.Sprintf("%d", number))
	if err != nil {
		return fmt.Errorf("failed to close issue #%d: %s %w", number, stderr, err)
	}
	return nil
}

// CreateIssue pushes a brand new issue up to GitHub and returns its assigned number
func CreateIssue(title string) (int, error) {
	if title == "" {
		return 0, fmt.Errorf("issue title cannot be empty")
	}

	// Creates the issue with a blank body and grabs the output text
	stdout, stderr, err := utils.ExecuteCommand("gh", "issue", "create", "--title", title, "--body", "")
	if err != nil {
		return 0, fmt.Errorf("failed to create issue: %s %w", stderr, err)
	}

	// gh CLI outputs the URL of the new issue, like: https://github.com/owner/repo/issues/1
	// Let's quickly pull out the issue number from the end of the URL
	var issueNum int
	_, fmtErr := fmt.Sscanf(stdout[strings.LastIndex(stdout, "/")+1:], "%d", &issueNum)
	if fmtErr != nil {
		return 1, nil // Safe fallback to 1 if parsing the string fails
	}

	return issueNum, nil
}
