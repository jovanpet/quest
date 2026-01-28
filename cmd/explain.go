package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	copilot "github.com/jovanpet/quest/internal/copilot"
	"github.com/jovanpet/quest/internal/format"
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
	Run: func(cmd *cobra.Command, args []string) {
		if err := runExplain(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}

func runExplain() error {
	// Load current state
	state, err := quest.LoadState()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Load plan
	plan, err := quest.LoadPlan()
	if err != nil {
		return fmt.Errorf("failed to load plan: %w", err)
	}

	// Get current task
	tasks := quest.FlattenTasks(plan)
	if state.CurrentTaskIndex >= len(tasks) {
		return fmt.Errorf("no active task")
	}

	currentTask := tasks[state.CurrentTaskIndex]
	
	// Increment explain count for current task
	state.ExplainCount++
	
	// Show context based on attempt count
	format.ExplainHeader(state.ExplainCount, currentTask.Title)

	// Collect student files from artifacts
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	var files []string
	for _, artifact := range currentTask.Artifacts {
		// Check if file exists
		fullPath := filepath.Join(workDir, artifact)
		if _, err := os.Stat(fullPath); err == nil {
			files = append(files, artifact)
		}
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found to analyze. Expected: %v", currentTask.Artifacts)
	}

	// Start spinner while generating hints
	spinner := format.NewSpinner("AI analyzing your code...")
	spinner.Start()

	// Generate hints via Copilot with attempt count
	hints, err := copilot.GenerateHints(
		currentTask.Title,
		currentTask.Objective,
		files,
		state.ExplainCount,
	)

	// Stop spinner before showing results
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to generate hints: %w", err)
	}

	if len(hints) == 0 {
		format.Printf("  %s✨ Your code looks good! No hints needed%s\n\n", 
			format.ColorGreen, format.ColorReset)
		// Save state even if no hints
		quest.UploadState(state)
		return nil
	}

	// Apply hints
	if err := copilot.ApplyHints(hints, workDir); err != nil {
		return fmt.Errorf("failed to apply hints: %w", err)
	}

	// Show what was added
	format.HintSummary(len(hints))
	
	for _, hint := range hints {
		format.Printf("    %s%s:%d%s - %s\n", 
			format.ColorDim, hint.File, hint.Line, format.ColorReset, hint.Comment)
	}

	format.CommandHint("\nCheck your files for the new comments!", "")
	
	// Save updated state with explain count
	if err := quest.UploadState(state); err != nil {
		format.Warning("Failed to save state")
	}
	
	return nil
}
