package cmd

import (
	"github.com/jovanpet/quest/internal/quest"
	"github.com/spf13/cobra"
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Get AI explanations and hints for the current task",
	Long: `Analyzes your code and adds helpful comment hints to guide you.

The AI will review your work and insert TODO/HINT comments at specific 
lines where you might need to consider edge cases or fix issues.

Example:
  quest explain`,
	Run: quest.RunExplain,
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
