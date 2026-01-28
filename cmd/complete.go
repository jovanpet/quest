package cmd

import (
	"github.com/jovanpet/quest/internal/quest"
	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark the current task as complete and advance to the next one",
	Long: `Mark the current task as complete and automatically advance to the next task.

Use this command after you've finished implementing a task and validated it with 'quest check'.
Your progress will be saved and the next task will be displayed.`,
	Run: quest.RunCompete,
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
