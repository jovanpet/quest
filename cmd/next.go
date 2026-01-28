/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/jovanpet/quest/internal/quest"
	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Move to the next task in your quest",
	Long:  `Advances to the next task and displays what you need to work on.`,
	Run:   quest.RunNext,
}

func init() {
	rootCmd.AddCommand(nextCmd)
}
