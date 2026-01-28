package cmd

import (
	"github.com/jovanpet/quest/internal/quest"
	"github.com/spf13/cobra"
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "View your quest progress and completion status",
	Long: `Display a comprehensive overview of your quest progress.

Shows all chapters, tasks, and your completion status across the entire quest.
Use this to track what you've accomplished and what's left to complete.`,
	Run: quest.RunSummary,
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}
