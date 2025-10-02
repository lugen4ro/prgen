package internal

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// GeneratePRContentWithClaude generates both PR title and body using Claude Code CLI
func GeneratePRContentWithClaude(config *Config, diff, background string) (title, body string, err error) {
	// Check if Claude Code CLI is available
	if _, err := exec.LookPath("claude"); err != nil {
		return "", "", fmt.Errorf("claude CLI not found. Please install Claude Code CLI first")
	}

	// Filter and summarize the diff to manage token usage
	filteredDiff := diff
	if estimateTokens(diff) > 6000 { // Use constant from diff_filter.go
		summary, err := FilterDiff(diff)
		if err != nil {
			return "", "", fmt.Errorf("failed to filter diff: %w", err)
		}
		filteredDiff = summary.FilteredDiff
	}

	// Build combined prompt for both title and body
	combinedPrompt := buildCombinedPrompt(config, filteredDiff, background)

	// Check token limit for combined prompt
	if estimateTokens(combinedPrompt) > MaxInputTokens {
		return "", "", fmt.Errorf("combined prompt too large (%d estimated tokens, max %d)", estimateTokens(combinedPrompt), MaxInputTokens)
	}

	response, err := callClaudeCLI(combinedPrompt)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate PR content: %w", err)
	}

	// Parse the response to extract title and body
	title, body, err = parseCombinedResponse(response)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse Claude response: %w", err)
	}

	return title, body, nil
}

// callClaudeCLI executes the Claude Code CLI with the given prompt
func callClaudeCLI(prompt string) (string, error) {
	ctx := context.Background()

	// Validate prompt is not empty
	if strings.TrimSpace(prompt) == "" {
		return "", fmt.Errorf("empty prompt provided to Claude CLI")
	}

	// Use Claude Code CLI in print mode with text output format
	cmd := exec.CommandContext(ctx, "claude", "-p", "--output-format", "text")
	cmd.Stdin = strings.NewReader(prompt)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude CLI error (exit code %d): %s", exitError.ExitCode(), string(output))
		}
		return "", fmt.Errorf("failed to execute claude CLI: %w, output: %s", err, string(output))
	}

	response := strings.TrimSpace(string(output))
	if response == "" {
		return "", fmt.Errorf("empty response from Claude CLI")
	}

	return response, nil
}

// buildCombinedPrompt constructs a single prompt for generating both PR title and body
func buildCombinedPrompt(config *Config, diff, background string) string {
	prompt := "Please generate both a PR title and PR body based on the following requirements and git diff.\n\n"

	if strings.TrimSpace(background) != "" {
		prompt += "BACKGROUND INFORMATION:\n"
		prompt += background + "\n\n"
	}

	prompt += "TITLE REQUIREMENTS:\n"
	prompt += config.TitleInstructions + "\n\n"

	if config.TitleExample != "" {
		prompt += "TITLE EXAMPLE:\n" + config.TitleExample + "\n\n"
	}

	prompt += "BODY REQUIREMENTS:\n"
	prompt += config.BodyInstructions + "\n\n"

	if config.BodyExample != "" {
		prompt += "BODY EXAMPLE:\n" + config.BodyExample + "\n\n"
	}

	prompt += "GIT DIFF:\n" + diff + "\n\n"

	prompt += "Please respond with the following format:\n"
	prompt += "TITLE: [your generated title]\n"
	prompt += "BODY:\n[your generated body]"

	return prompt
}

// parseCombinedResponse parses the Claude response to extract title and body
func parseCombinedResponse(response string) (title, body string, err error) {
	lines := strings.Split(response, "\n")

	var titleFound bool
	var bodyLines []string
	var collectingBody bool

	for _, line := range lines {
		if strings.HasPrefix(line, "TITLE:") {
			title = strings.TrimSpace(strings.TrimPrefix(line, "TITLE:"))
			titleFound = true
		} else if strings.HasPrefix(line, "BODY:") {
			collectingBody = true
		} else if collectingBody {
			bodyLines = append(bodyLines, line)
		}
	}

	if !titleFound || title == "" {
		return "", "", fmt.Errorf("could not extract title from Claude response")
	}

	body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
	if body == "" {
		return "", "", fmt.Errorf("could not extract body from Claude response")
	}

	return title, body, nil
}
