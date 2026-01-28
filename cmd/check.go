/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/jovanpet/quest/internal/quest"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if your current task is complete",
	Long:  `Validates that you have completed the requirements for the current task.`,
	Run:   quest.RunCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolP("annotate", "a", false, "Add inline comments to code showing check results")
}
