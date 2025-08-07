package internal

import (
	"os/exec"
	"strings"
)

// GetDiff executes 'git diff' and returns the diff output.
// It returns an empty string and nil error if there are no changes,
// or returns an error if the git command fails.
func GetDiff() (string, error) {
	cmd := exec.Command("git", "diff")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
