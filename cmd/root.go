/*
Copyright © 2025 Lukas Nakamura lugen4ro@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lugen4ro/prgen/internal"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "prgen",
	Short: "A CLI tool for generating customized GitHub PRs with the power of AI.",
	Long: `Hello there. 
prgen is a CLI tool that allows users to generate and create PRs on GitHub in seconds. 
Customize your experience to get exactly the PR you want.

Your personal config files such as templates are stored under ~/.config/prgen/`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		configFlag, _ := cmd.Flags().GetBool("config")
		if configFlag {
			openConfigFile()
			return
		}
		internal.Construct()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.prgen.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("config", "c", false, "Open the main config file in the default editor")
}

// openConfigFile opens the main config file with the default editor
func openConfigFile() {
	config, err := internal.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	configPath := config.GetConfigPath()

	// Try to get the default editor from environment variables
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Default to common editors based on OS
		editor = "vi" // Unix default
	}

	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error opening config file with %s: %v\n", editor, err)
		fmt.Printf("Config file location: %s\n", configPath)
		os.Exit(1)
	}
}
