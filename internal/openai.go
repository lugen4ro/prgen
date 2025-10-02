package internal

import (
	"context"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/sashabaranov/go-openai"
)

const (
	// MaxInputTokens is the maximum number of tokens we'll send to OpenAI
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

// GeneratePRContent generates both PR title and body using OpenAI
func GeneratePRContent(config *Config, diff, background string) (title, body string, err error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", "", fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	client := openai.NewClient(apiKey)
	ctx := context.Background()

	// Filter and summarize the diff to manage token usage
	filteredDiff := diff
	if estimateTokens(diff) > 6000 { // Use constant from diff_filter.go
		summary, err := FilterDiff(diff)
		if err != nil {
			return "", "", fmt.Errorf("failed to filter diff: %w", err)
		}
		filteredDiff = summary.FilteredDiff
	}

	// Get model from config, default to gpt-4
	model := "gpt-4"
	if modelValue, ok := config.MainConfig["model"].(string); ok && modelValue != "" {
		model = modelValue
	}

	// Generate PR title
	titlePrompt := buildTitlePrompt(config, filteredDiff, background)

	// Check token limit for title prompt
	if estimateTokens(titlePrompt) > MaxInputTokens {
		return "", "", fmt.Errorf("title prompt too large (%d estimated tokens, max %d)", estimateTokens(titlePrompt), MaxInputTokens)
	}

	titleResp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: titlePrompt,
			},
		},
		MaxTokens:   100,
		Temperature: 0.3,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate title: %w", err)
	}

	if len(titleResp.Choices) == 0 {
		return "", "", fmt.Errorf("no title generated")
	}
	title = titleResp.Choices[0].Message.Content

	// Generate PR body
	bodyPrompt := buildBodyPrompt(config, filteredDiff, background)

	// Check token limit for body prompt
	if estimateTokens(bodyPrompt) > MaxInputTokens {
		return "", "", fmt.Errorf("body prompt too large (%d estimated tokens, max %d)", estimateTokens(bodyPrompt), MaxInputTokens)
	}

	bodyResp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: bodyPrompt,
			},
		},
		MaxTokens:   2000,
		Temperature: 0.3,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate body: %w", err)
	}

	if len(bodyResp.Choices) == 0 {
		return "", "", fmt.Errorf("no body generated")
	}
	body = bodyResp.Choices[0].Message.Content

	return title, body, nil
}

// buildTitlePrompt constructs the prompt for generating PR title
func buildTitlePrompt(config *Config, diff, background string) string {
	prompt := config.TitleInstructions + "\n\n"

	if strings.TrimSpace(background) != "" {
		prompt += "Background information:\n" + background + "\n\n"
	}

	if config.TitleExample != "" {
		prompt += "Example:\n" + config.TitleExample + "\n\n"
	}

	prompt += "Git diff:\n" + diff + "\n\n"
	prompt += "Generate a concise PR title based on the above information:"

	return prompt
}

// buildBodyPrompt constructs the prompt for generating PR body
func buildBodyPrompt(config *Config, diff, background string) string {
	prompt := config.BodyInstructions + "\n\n"

	if strings.TrimSpace(background) != "" {
		prompt += "Background information:\n" + background + "\n\n"
	}

	if config.BodyExample != "" {
		prompt += "Example:\n" + config.BodyExample + "\n\n"
	}

	prompt += "Git diff:\n" + diff + "\n\n"
	prompt += "Generate a comprehensive PR body based on the above information:"

	return prompt
}
