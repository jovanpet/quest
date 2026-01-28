package cmd

import (
	"github.com/jovanpet/quest/internal/quest"
	"github.com/spf13/cobra"
)

// jumpToCmd represents the jumpTo command
var jumpToCmd = &cobra.Command{
	Use:   "jumpTo",
	Short: "Jump to a specific task index or the last completed task",
	Long: `Jump to a specific task index or the last completed task in your quest.

You can specify the task index as an argument or use the --last-complete flag to jump to the last completed task.`,
	Run: quest.RunJumpTo,
}

func init() {
	rootCmd.AddCommand(jumpToCmd)
	jumpToCmd.Flags().BoolP("last-complete", "l", false, "Jump to the last completed task")
}
