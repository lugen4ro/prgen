package internal

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetDiff executes 'git diff' between the current branch and main branch.
// It returns the diff output showing changes that would be included in a PR.
// Returns an empty string and nil error if there are no changes,
// or returns an error if the git command fails.
func GetDiff() (string, error) {
	cmd := exec.Command("git", "diff", "main...HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetRemoteOrigin gets the origin remote URL
func GetRemoteOrigin() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCurrentBranch gets the current git branch name
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// PushCurrentBranch pushes the current branch to origin
func PushCurrentBranch() error {
	// Get current branch name
	branch, err := GetCurrentBranch()
	if err != nil {
		return err
	}

	// Push current branch to origin with upstream tracking
	cmd := exec.Command("git", "push", "-u", "origin", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to push branch %s: %s", branch, string(output))
	}

	return nil
}
