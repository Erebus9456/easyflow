package git

import (
	"strings"

	"github.com/Erebus9456/easyflow/utils"
)

// RepoContext encapsulates the parsed ownership data for the current workspace
type RepoContext struct {
	Owner          string
	RepositoryName string
	CurrentBranch  string
}

// GetCurrentRepo inspects the localized workspace context dynamically
func GetCurrentRepo() (*RepoContext, error) {
	// 1. Verify local git state and extract remote origin URL
	out, _, err := utils.ExecuteCommand("git", "config", "--get", "remote.origin.url")
	if err != nil {
		return nil, utils.ErrNotAGitRepository
	}

	if out == "" {
		return nil, utils.ErrNoRemoteOrigin
	}

	// 2. Extract current active local branch
	branchOut, _, err := utils.ExecuteCommand("git", "branch", "--show-current")
	if err != nil {
		branchOut = "main" // Safe fallback default
	}

	// 3. Parse GitHub URL formats (Handles both SSH and HTTPS variants)
	// Example SSH:   git@github.com:owner/repo.git or git@github.com:owner/repo
	// Example HTTPS: https://github.com/owner/repo.git or https://github.com/owner/repo
	urlClean := strings.TrimSuffix(out, ".git")
	var pathSegments []string

	if strings.Contains(urlClean, "git@github.com:") {
		splitSSH := strings.Split(urlClean, "git@github.com:")
		if len(splitSSH) > 1 {
			pathSegments = strings.Split(splitSSH[1], "/")
		}
	} else if strings.Contains(urlClean, "https://github.com/") {
		splitHTTPS := strings.Split(urlClean, "https://github.com/")
		if len(splitHTTPS) > 1 {
			pathSegments = strings.Split(splitHTTPS[1], "/")
		}
	}

	if len(pathSegments) < 2 {
		return nil, utils.ErrNoRemoteOrigin
	}

	return &RepoContext{
		Owner:          pathSegments[0],
		RepositoryName: pathSegments[1],
		CurrentBranch:  branchOut,
	}, nil
}
