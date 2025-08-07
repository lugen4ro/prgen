package internal

import (
	"fmt"
	"log"
)

// Construct creates the PR proposal using LLMs.
// This is the main entrypoint for PR generation.
func Construct() {
	fmt.Println("Starting PR proposal construction...")

	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Config loaded from: %s\n", config.ConfigDir)
	fmt.Printf("LLM Provider: %v\n", config.MainConfig["llm_provider"])
	fmt.Printf("Model: %v\n", config.MainConfig["model"])

	// Get git diff
	diff, err := GetDiff()
	if err != nil {
		fmt.Printf("Failed to get diff: %v\n", err)
		return
	}

	if diff == "" {
		fmt.Println("No changes detected. Nothing to generate PR for.")
		return
	}

	fmt.Printf("Found changes (%d characters)\n", len(diff))
	
	// Display loaded templates and instructions
	fmt.Println("\n--- Configuration Summary ---")
	fmt.Printf("Body Template: %d characters\n", len(config.BodyTemplate))
	fmt.Printf("Title Template: %d characters\n", len(config.TitleTemplate))
	fmt.Printf("Body Instructions: %d characters\n", len(config.BodyInstructions))
	fmt.Printf("Title Instructions: %d characters\n", len(config.TitleInstructions))
	
	fmt.Println("\nReady to generate PR with LLM!")
}
