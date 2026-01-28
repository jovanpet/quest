package cmd

import (
	"github.com/jovanpet/quest/internal/quest"
	"github.com/spf13/cobra"
)

// beginCmd represents the begin command
var beginCmd = &cobra.Command{
	Use:   "begin",
	Short: "Start a new coding quest",
	Long: `Start a new coding quest with three options:

1. Pick a Legendary Path - Choose from curated templates (REST APIs, CLI tools, concurrency, etc.)
2. Forge Your Own Quest - Customize and generate a quest with AI assistance
3. Seek a Mystery Quest - Get a surprise AI-generated quest

Each quest contains chapters with tasks, automated validation, and AI-powered hints.`,
	Run: quest.RunBegin,
}

func init() {
	rootCmd.AddCommand(beginCmd)
}
