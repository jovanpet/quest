package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "quest",
	Short: "Interactive CLI tool for learning Go through hands-on quests",
	Long: `Quest is an interactive learning platform that teaches Go through hands-on coding challenges.

Choose from curated templates covering REST APIs, CLI tools, concurrency patterns, and more.
Each quest guides you through tasks with automated validation and AI-powered hints.

Perfect for beginners learning Go fundamentals or experienced developers exploring new patterns.`,
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
	// Future: Add config file support
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.quest.yaml)")
}


