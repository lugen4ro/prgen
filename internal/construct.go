package internal

import (
	"fmt"
)

// Construct creates the PR proposal using LLMs.
// This is the main entrypoint for PR generation.
func Construct() {
	fmt.Println("Starting PR proposal construction...")

	diff, err := GetDiff()
	if err != nil {
		fmt.Printf("Failed to get diff: %v\n", err)
		return
	}

	fmt.Printf("got %s\n", diff)
}
