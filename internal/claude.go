package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"unicode"
)

// claudeJSONResponse represents the JSON output from Claude CLI
type claudeJSONResponse struct {
	Type      string `json:"type"`
	Result    string `json:"result"`
	SessionID string `json:"session_id"`
	IsError   bool   `json:"is_error"`
}

const (
	// MaxInputTokens is the maximum number of tokens we'll send to the LLM
	// This helps control costs by preventing very large inputs
	MaxInputTokens = 8000
)

// estimateTokens provides a rough estimate of token count
// Accounts for both English and Japanese text patterns:
// - English: ~4 characters per token
// - Japanese: Each character (hiragana, katakana, kanji) ≈ 1 token
// - Code/symbols: More conservative estimation
func estimateTokens(text string) int {
	var japaneseChars, englishChars, otherChars int

	for _, r := range text {
		switch {
		case unicode.In(r, unicode.Hiragana, unicode.Katakana, unicode.Han):
			japaneseChars++
		case unicode.IsLetter(r) && r < 128: // ASCII letters
			englishChars++
		default:
			otherChars++
		}
	}

	// Conservative estimation:
	// - Japanese chars: 1 token each
	// - English chars: 1 token per 3 chars (more conservative than 4)
	// - Other chars (code, symbols): 1 token per 2 chars
	estimatedTokens := japaneseChars + englishChars/3 + otherChars/2

	// Add word-based estimation for English text structure
	wordCount := len(strings.Fields(text))

	// Return the higher of the two estimates for safety
	if wordCount > estimatedTokens {
		return wordCount
	}
	return estimatedTokens
}

// PRGenerationResult holds the result of PR content generation
type PRGenerationResult struct {
	Title     string
	Body      string
	SessionID string
}

// GeneratePRContentWithClaude generates both PR title and body using Claude Code CLI
// If refinement is provided, it will refine the previous output based on user feedback
// The session ID from refinement context is used to continue the conversation
func GeneratePRContentWithClaude(config *Config, diff, background string, refinement *RefinementContext) (*PRGenerationResult, error) {
	// Check if Claude Code CLI is available
	if _, err := exec.LookPath("claude"); err != nil {
		return nil, fmt.Errorf("claude CLI not found. Please install Claude Code CLI first")
	}

	// Filter and summarize the diff to manage token usage
	filteredDiff := diff
	if estimateTokens(diff) > 6000 { // Use constant from diff_filter.go
		summary, err := FilterDiff(diff)
		if err != nil {
			return nil, fmt.Errorf("failed to filter diff: %w", err)
		}
		filteredDiff = summary.FilteredDiff
	}

	var combinedPrompt string
	var sessionID string

	if refinement != nil {
		// For refinement, we just send the feedback since we're continuing the session
		combinedPrompt = buildRefinementPrompt(refinement)
		sessionID = refinement.SessionID
	} else {
		// Build full prompt for initial generation
		combinedPrompt = buildCombinedPrompt(config, filteredDiff, background)
	}

	// Check token limit for combined prompt
	if estimateTokens(combinedPrompt) > MaxInputTokens {
		return nil, fmt.Errorf("combined prompt too large (%d estimated tokens, max %d)", estimateTokens(combinedPrompt), MaxInputTokens)
	}

	// Call Claude CLI (either new session or resume existing)
	response, newSessionID, err := callClaudeCLI(combinedPrompt, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PR content: %w", err)
	}

	// Parse the response to extract title and body
	title, body, err := parseCombinedResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Claude response: %w", err)
	}

	return &PRGenerationResult{
		Title:     title,
		Body:      body,
		SessionID: newSessionID,
	}, nil
}

// callClaudeCLI executes the Claude Code CLI with the given prompt
// If sessionID is provided, it resumes that session; otherwise starts a new one
// Returns the response text and the session ID for future continuation
func callClaudeCLI(prompt string, sessionID string) (response string, newSessionID string, err error) {
	ctx := context.Background()

	// Validate prompt is not empty
	if strings.TrimSpace(prompt) == "" {
		return "", "", fmt.Errorf("empty prompt provided to Claude CLI")
	}

	// Build command arguments
	args := []string{"-p", "--output-format", "json"}
	if sessionID != "" {
		// Resume existing session
		args = append(args, "--resume", sessionID)
	}

	cmd := exec.CommandContext(ctx, "claude", args...)
	cmd.Stdin = strings.NewReader(prompt)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", "", fmt.Errorf("claude CLI error (exit code %d): %s", exitError.ExitCode(), string(output))
		}
		return "", "", fmt.Errorf("failed to execute claude CLI: %w, output: %s", err, string(output))
	}

	// Parse JSON response
	var jsonResp claudeJSONResponse
	if err := json.Unmarshal(output, &jsonResp); err != nil {
		return "", "", fmt.Errorf("failed to parse Claude JSON response: %w, raw output: %s", err, string(output))
	}

	if jsonResp.IsError {
		return "", "", fmt.Errorf("Claude returned an error: %s", jsonResp.Result)
	}

	if jsonResp.Result == "" {
		return "", "", fmt.Errorf("empty response from Claude CLI")
	}

	return jsonResp.Result, jsonResp.SessionID, nil
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

// buildRefinementPrompt constructs a prompt for refining a previously generated PR
// Since we're continuing the session, Claude already has context from the previous exchange
func buildRefinementPrompt(refinement *RefinementContext) string {
	prompt := "Please refine the PR title and body based on my feedback:\n\n"
	prompt += refinement.Feedback + "\n\n"
	prompt += "Please respond in the same format as before:\n"
	prompt += "TITLE: [your refined title]\n"
	prompt += "BODY:\n[your refined body]"

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
