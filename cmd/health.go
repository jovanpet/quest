package cmd

import (
	"github.com/jovanpet/quest/internal/quest"
	"github.com/spf13/cobra"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check the health status of the quest",
	Long: `This command checks the health status of the quest system and reports any issues found.
	Checks include: checking the .quest folder, loading and unloading configurations, and verifying system components.`,
	Run: quest.RunHealthCheck,
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
