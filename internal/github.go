package internal

import (
	"fmt"
	"os/exec"
	"strings"
)

// CreateGitHubPR creates a pull request using gh CLI
func CreateGitHubPR(title, body string) (string, error) {
	// Check if gh CLI is available
	if err := checkGHCLI(); err != nil {
		return "", err
	}

	// Use gh CLI to create PR in draft state
	// gh pr create --title "title" --body "body" --base main --draft
	cmd := exec.Command("gh", "pr", "create", "--title", title, "--body", body, "--base", "main", "--draft")

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("gh pr create failed: %s", string(exitError.Stderr))
		}
		return "", fmt.Errorf("failed to execute gh pr create: %w", err)
	}

	// Return the PR URL from gh output
	prURL := strings.TrimSpace(string(output))
	return prURL, nil
}

// checkGHCLI verifies that gh CLI is installed and authenticated
func checkGHCLI() error {
	// Check if gh is installed
	cmd := exec.Command("gh", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gh CLI is not installed. Please install it from https://cli.github.com/")
	}

	// Check if user is authenticated
	cmd = exec.Command("gh", "auth", "status")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not authenticated with GitHub. Please run 'gh auth login'")
	}

	return nil
}

// OpenPRInBrowser opens the current PR in the web browser
func OpenPRInBrowser() error {
	cmd := exec.Command("gh", "pr", "view", "--web")
	return cmd.Run()
}
