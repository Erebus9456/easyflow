package git

import (
	"strconv"
	"strings"

	"github.com/Erebus9456/easyflow/utils"
)

// RepoContext encapsulates the parsed ownership data for the current workspace
type RepoContext struct {
	Owner          string
	RepositoryName string
	CurrentBranch  string
}

// StatusMetrics tracks raw numeric quantities of uncommitted adjustments
type StatusMetrics struct {
	Modified  int
	Untracked int
	Added     int
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

// 🆕 IsWorkspaceDirty performs a fast-fail check for local modifications
func IsWorkspaceDirty() (bool, error) {
	stdout, _, err := utils.ExecuteCommand("git", "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(stdout) != "", nil
}

// 🆕 GetStatusMetrics parses porcelain lines to group uncommitted changes
func GetStatusMetrics() (StatusMetrics, error) {
	var metrics StatusMetrics
	stdout, _, err := utils.ExecuteCommand("git", "status", "--porcelain")
	if err != nil {
		return metrics, err
	}

	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 2 {
			continue
		}
		prefix := line[:2]
		switch {
		case strings.Contains(prefix, "M"):
			metrics.Modified++
		case strings.Contains(prefix, "??"):
			metrics.Untracked++
		case strings.Contains(prefix, "A"):
			metrics.Added++
		}
	}
	return metrics, nil
}

// 🆕 GetUpstreamDivergence calculates exact ahead/behind metrics relative to origin
// GetUpstreamDivergence calculates exact ahead/behind metrics relative to origin
func GetUpstreamDivergence() (int, int, error) { // 👈 Fixed: Removed parameter names to fix syntax error
	// Query symlink matrix counts: HEAD versus upstream tracking definitions
	stdout, _, err := utils.ExecuteCommand("git", "rev-list", "--left-right", "--count", "HEAD...@{u}")
	if err != nil {
		// If there is no upstream tracking configured yet, return clean zeroes gracefully
		return 0, 0, nil
	}

	fields := strings.Fields(stdout)
	if len(fields) >= 2 {
		aheadCount, _ := strconv.Atoi(fields[0])
		behindCount, _ := strconv.Atoi(fields[1])
		return aheadCount, behindCount, nil
	}

	return 0, 0, nil
}
