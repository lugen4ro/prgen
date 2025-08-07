/*
Copyright © 2025 Lukas Nakamura lugen4ro@gmail.com
*/
package cmd

import (
	"os"

	"github.com/lugen4ro/prgen/internal"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "prgen",
	Short: "A CLI tool for generating customized GitHub PRs with the power of AI.",
	Long: `Hello there. 
prgen is a CLI tool that allows users to generate and create PRs on GitHub in seconds. 
Customize your experience to get exactly the PR you want.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
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
}
