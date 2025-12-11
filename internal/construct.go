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
	var result *PRGenerationResult
	err = RunSpinnerWithTask("Generating PR content", func() error {
		var err error
		result, err = GeneratePRContentWithProvider(config, diff, backgroundInfo)
		return err
	})
	if err != nil {
		ShowError("Failed to generate PR content", err)
		return
	}

	// Track session ID for conversation continuity
	title, body, sessionID := result.Title, result.Body, result.SessionID

	// Display generated content and handle refinement loop
	for {
		ShowGeneratedContent(title, body)

		// Ask user what they want to do
		choice := AskRefinementOrAccept()

		switch choice {
		case ChoiceAccept:
			// User accepted, break out of refinement loop
			break
		case ChoiceRefine:
			// Get feedback and refine
			feedback := AskRefinementFeedback()
			if feedback == "" {
				// No feedback provided, show content again
				continue
			}

			// Refine the PR content using the same session for conversation continuity
			refinement := &RefinementContext{
				SessionID: sessionID,
				Feedback:  feedback,
			}

			err = RunSpinnerWithTask("Refining PR content", func() error {
				var err error
				result, err = RefinePRContentWithProvider(config, diff, backgroundInfo, refinement)
				return err
			})
			if err != nil {
				ShowError("Failed to refine PR content", err)
				return
			}

			// Update with refined content (session ID should remain the same)
			title, body, sessionID = result.Title, result.Body, result.SessionID

			// Loop continues to show refined content
			continue
		case ChoiceCancel:
			fmt.Println(infoStyle.Render("ℹ️  PR creation cancelled by user"))
			return
		}

		// If we got here via ChoiceAccept, break the loop
		break
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
