package utils

import (
	"bytes"
	"os/exec"
	"strings"
)

// ExecuteCommand executes an OS-level command and cleanly returns
// its stdout, stderr, and an evaluation error.
func ExecuteCommand(name string, args ...string) (string, string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Return sanitized, trimmed string configurations
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}
