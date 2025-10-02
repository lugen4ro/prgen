package internal

import (
	"fmt"
)

// Construct creates the PR proposal using LLMs.
// This is the main entrypoint for PR generation.
func Construct() {
	// Initialize beautiful UI
	InitializeUI()
	ShowStartupBanner()

	// Load configuration with spinner
	var config *Config
	err := RunSpinnerWithTask("Loading configuration", func() error {
		var err error
		config, err = LoadConfig()
		return err
	})
	if err != nil {
		ShowError("Failed to load config", err)
		return
	}

	// Show configuration summary
	ShowConfigSummary(config)

	// Get git diff with spinner
	var diff string
	err = RunSpinnerWithTask("Analyzing git changes", func() error {
		var err error
		diff, err = GetDiff()
		return err
	})
	if err != nil {
		ShowError("Failed to get diff", err)
		return
	}

	// Check if there are changes
	ShowDiffInfo(len(diff))
	if diff == "" {
		return
	}

	// Collect background information from user
	backgroundInfo := AskBackgroundInfo()

	// Generate PR content with spinner
	var title, body string
	err = RunSpinnerWithTask("Generating PR content", func() error {
		// DEVELOPMENT MODE: Skip OpenAI to save tokens during testing
		// TODO: Remove this and uncomment the real generation code below
		// 		title = "feat: add OpenAI integration for PR generation"
		// 		body = `## TEST BODY
		// - Implemented OpenAI integration for generating PR titles and bodies
		// - Added token limit protection to control costs
		// `
		// return nil

		// PRODUCTION CODE (commented out for development):
		var err error
		title, body, err = GeneratePRContentWithProvider(config, diff, backgroundInfo)
		return err
	})
	if err != nil {
		ShowError("Failed to generate PR content", err)
		return
	}

	// Display generated content
	ShowGeneratedContent(title, body)

	// Ask for confirmation before creating PR
	if !AskConfirmation("Do you want to create this pull request?") {
		fmt.Println(infoStyle.Render("ℹ️  PR creation cancelled by user"))
		return
	}

	// Push current branch to remote with spinner
	err = RunSpinnerWithTask("Pushing current branch to remote", func() error {
		return PushCurrentBranch()
	})
	if err != nil {
		ShowError("Failed to push branch", err)
		return
	}

	// Create GitHub PR with spinner
	var prURL string
	err = RunSpinnerWithTask("Creating GitHub pull request", func() error {
		var err error
		prURL, err = CreateGitHubPR(title, body)
		return err
	})
	if err != nil {
		ShowError("Failed to create GitHub PR", err)
		return
	}

	// Show success with prominent URL display
	ShowPRSuccess(prURL)

	// Open PR in browser with spinner
	err = RunSpinnerWithTask("Opening PR in browser", func() error {
		return OpenPRInBrowser()
	})
	if err != nil {
		ShowError("Failed to open PR in browser", err)
		// Don't return here - this is not a critical error
	}
}
