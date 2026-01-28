package quest

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	copilot_helper "github.com/jovanpet/quest/internal/copilot"
	"github.com/jovanpet/quest/internal/format"
	"github.com/jovanpet/quest/internal/prompt"
	"github.com/jovanpet/quest/internal/template"
	"github.com/jovanpet/quest/internal/types"
	"github.com/spf13/cobra"
)

func FlattenTasks(p *types.Plan) []types.Task {
	var tasks []types.Task
	for _, chapter := range p.Chapters {
		for _, quest := range chapter.Quests {
			tasks = append(tasks, quest.Tasks...)
		}
	}
	return tasks
}

func RunBegin(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(QuestFolderName); err == nil {
		format.Warning("Quest started and .quest directory already exists.")
		return
	}

	err := runBegin()
	if err != nil {
		format.Warning("Trying to remove .quest folder due to error")
		err = os.RemoveAll(QuestFolderName)
		if err != nil {
			format.ErrorWithTip("Failed to remove .quest folder", err, "Try manually deleting the .quest folder")
		}
		return
	}
}

func runBegin() error {

	format.Header("Starting a new quest")

	// Epic description
	fmt.Printf("  %sEmbark on an interactive coding adventure!%s\n", format.ColorDim, format.ColorReset)
	fmt.Printf("  %sLearn by building with guided tasks, AI hints, and instant validation.%s\n", format.ColorDim, format.ColorReset)
	fmt.Println()

	err := os.MkdirAll(QuestFolderName, 0755)
	if err != nil {
		format.ErrorWithTip("Error creating directory for the quest", err, "Check folder permissions")
		return err
	}

	// Step 1: Select wizard mode
	mode, cancelled, err := prompt.SelectWizardMode()
	if err != nil {
		format.Error("Failed to get wizard selection", err)
		return err
	}
	if cancelled {
		format.Warning(" Cancelled. No changes made.")
		return fmt.Errorf("")
	}

	var plan *types.Plan

	switch mode {
	case prompt.ModeTemplate:
		// Get available templates
		templates := template.ListTemplates()
		if len(templates) == 0 {
			return fmt.Errorf("no templates found")
		}

		// Prompt user to select template
		templateName, cancelled, err := prompt.SelectTemplate(templates)
		if err != nil {
			format.Error("Failed to get template selection", err)
			return err
		}
		if cancelled {
			format.Warning(" Cancelled. No changes made.")
			return fmt.Errorf("")
		}

		plan, err = template.Load(templateName)
		if err != nil {
			format.ErrorWithTip("Error loading template", err, "Check that the template exists")
			return err
		}

	case prompt.ModeForge:
		spec, cancelled, err := prompt.ForgeQuestWizard()
		if err != nil {
			format.Error("Failed to forge quest", err)
			return err
		}
		if cancelled {
			format.Warning("⚠ Cancelled. No changes made.")
			return fmt.Errorf("")
		}

		// Generate plan with AI
		spinner := format.NewSpinner("⚡ AI forging your custom quest...")
		spinner.Start()

		plan, err = copilot_helper.GeneratePlan(*spec)

		spinner.Stop()

		if err != nil {
			templateName := selectTemplateForSpec(spec)
			format.Warning(fmt.Sprintf("AI plan generation failed (%s). Falling back to template '%s'.", err.Error(), templateName))
			plan, err = template.Load(templateName)
			if err != nil {
				return err
			}
		} else {
			format.Success(fmt.Sprintf("Quest forged! %s %s with %d tasks",
				spec.BuildType, spec.Difficulty, countTasks(*plan)))
		}

	case prompt.ModeSurprise:
		// Random quest
		spec := prompt.SurpriseQuest()
		format.Success(fmt.Sprintf("🎲 Surprise! Building a %s (%s difficulty) with theme: %s",
			spec.BuildType, spec.Difficulty, spec.Theme))

		// Generate plan with AI
		spinner := format.NewSpinner("✨ AI crafting your mystery quest...")
		spinner.Start()

		plan, err = copilot_helper.GeneratePlan(*spec)

		spinner.Stop()

		if err != nil {
			format.Warning(fmt.Sprintf("AI generation failed (%s). Using random template instead.", err.Error()))
			templates := template.ListTemplates()
			randomIdx := time.Now().UnixNano() % int64(len(templates))
			templateName := strings.Split(templates[randomIdx], " - ")[0]
			plan, err = template.Load(templateName)
			if err != nil {
				return err
			}
		} else {
			format.Success(fmt.Sprintf("Mystery quest revealed! %d tasks await you",
				countTasks(*plan)))
		}
	}

	plan.NumberOfTasks = countTasks(*plan)

	firstState := &types.State{
		Version:          1,
		CurrentTaskIndex: 0,
		CompletedTaskIDs: []string{},
		LastCheck:        nil,
		QuestStarted:     false,
	}

	err = UploadStateAndPlan(firstState, plan)
	if err != nil {
		return err
	}

	// Show quest ready message
	format.Newline()
	format.Line(fmt.Sprintf("%s%s%s", format.ColorBold, plan.Journey.Name, format.ColorReset))
	if plan.Journey.Description != "" {
		format.Line(fmt.Sprintf("%s%s%s", format.ColorDim, plan.Journey.Description, format.ColorReset))
	}
	format.Newline()

	lines := []string{
		fmt.Sprintf("%s✨ Quest initialized successfully!%s", format.ColorGreen, format.ColorReset),
		"",
		fmt.Sprintf("%sFiles created:%s", format.ColorDim, format.ColorReset),
		fmt.Sprintf("  %s• %s%s", format.ColorCyan, PlanFilePath, format.ColorReset),
		fmt.Sprintf("  %s• %s%s", format.ColorCyan, StateFilePath, format.ColorReset),
	}

	format.Box("Quest Ready", lines)
	format.CommandHint("Ready to start? Run", "quest next")

	return nil
}

func RunNext(cmd *cobra.Command, args []string) {
	state, plan, err := LoadStateAndPlan()
	if err != nil {
		format.ErrorWithTip("Failed to load quest data", err, "Run 'quest begin' to start a new quest")
		return
	}

	if state.LastCheck != nil || state.LastCheck.Status != types.CheckPass {
		format.Warning("The previous check didn't pass.")
		format.Print("Are you sure you want to continue without passing the check? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			format.Info("Run 'quest check' to validate your progress first.")
			return
		}
		format.Info("Continuing without check validation...")
	}

	if state.QuestStarted {
		state.CurrentTaskIndex++
		// Reset explain count when moving to new task
		state.ExplainCount = 0
	} else {
		state.QuestStarted = true
	}

	// Check bounds after increment
	if state.CurrentTaskIndex >= plan.NumberOfTasks {
		format.Success("You've reached the final task! Run 'quest complete' when done.")
		return
	}

	tasks := FlattenTasks(plan)
	currentTask := tasks[state.CurrentTaskIndex]

	for _, file := range currentTask.Files {
		err := createArtifact(file)
		if err != nil {
			format.ErrorWithTip("Error creating artifact", err, fmt.Sprintf("Check permissions for %s", file))
			return
		}
	}

	state.LastCheck = nil
	err = UploadState(state)
	if err != nil {
		format.ErrorWithTip("Failed to save state", err, "Check folder permissions")
		return
	}

	// Display the new task
	format.TaskHeader(state.CurrentTaskIndex+1, currentTask.Title)

	format.Dim(currentTask.Objective)
	format.Newline()

	if len(currentTask.Steps) > 0 {
		format.Line(fmt.Sprintf("%sSteps:%s", format.ColorPink, format.ColorReset))
		for i, step := range currentTask.Steps {
			format.Step(i+1, step)
		}
		format.Newline()
	}

	if len(currentTask.Artifacts) > 0 {
		format.Line(fmt.Sprintf("%sFiles to create:%s", format.ColorPink, format.ColorReset))
		for _, artifact := range currentTask.Artifacts {
			format.File(artifact)
		}
		format.Newline()
	}

	format.CommandHint("When ready, run", "quest check")
}

func RunCheck(cmd *cobra.Command, args []string) {
	state, plan, err := LoadStateAndPlan()
	if err != nil {
		format.ErrorWithTip("Failed to load quest data", err, "Run 'quest begin' to start a new quest")
		return
	}
	tasks := FlattenTasks(plan)

	// Bounds check for safety
	if state.CurrentTaskIndex >= len(tasks) {
		format.Error("Invalid task index", fmt.Errorf("index out of bounds"))
		return
	}

	currentTask := tasks[state.CurrentTaskIndex]
	currentTaskValidation := &currentTask.Validation

	format.CheckHeader(state.CurrentTaskIndex+1, currentTask.Title)

	lastCheck := &types.CheckResult{
		TaskID:    state.CurrentTaskIndex,
		Timestamp: time.Now(),
	}
	passState := types.Pass
	failState := types.Fail

	passedCount := 0
	failedCount := 0

	for i, rule := range currentTaskValidation.Rules {
		check, err := CheckRule(rule)
		if check {
			currentTaskValidation.Rules[i].LastState = &passState
			passedCount++

			// Show success with reason
			ruleName := rule.Name
			if ruleName == "" {
				ruleName = fmt.Sprintf("Rule %d", i+1)
			}

			// Build success message based on rule type
			successReason := ""
			switch rule.Type {
			case types.TypeExists:
				successReason = fmt.Sprintf("- found '%s'", rule.Path)
			case types.TypeGlobCountMin:
				count, _ := CountFilesMatching(rule.Glob)
				successReason = fmt.Sprintf("- found %d file(s) matching '%s'", count, rule.Glob)
			case types.TypeFileContainsAny:
				if len(rule.Any) == 1 {
					successReason = fmt.Sprintf("- %s - contains '%s'", rule.Glob, rule.Any[0])
				} else {
					successReason = fmt.Sprintf("- %s - contains required '%v'", rule.Glob, rule.Any)
				}
			}

			format.CheckPass(ruleName, successReason)
		} else {
			failedCount++

			currentTaskValidation.Rules[i].LastState = &failState

			// Show failure - name + reason
			ruleName := rule.Name
			if ruleName == "" {
				ruleName = fmt.Sprintf("Rule %d", i+1)
			}

			// Show the failure reason on the same line
			failureReason := ""
			if err != nil {
				failureReason = err.Error()
			} else if rule.Description != "" {
				failureReason = rule.Description
			}

			format.CheckFail(ruleName, failureReason)
		}
	}

	// Set final status
	if failedCount == 0 {
		lastCheck.Status = types.CheckPass
		lastCheck.Message = "All checks passed"
		format.CheckSummaryPass(passedCount)

		// Mark task as complete
		taskID := currentTask.ID
		if !contains(state.CompletedTaskIDs, taskID) {
			state.CompletedTaskIDs = append(state.CompletedTaskIDs, taskID)
		}

		format.CommandHint("Ready for next task? Run", "quest next")
		format.CommandHint("To Run an AI check", "quest check --annotate")
	} else {
		lastCheck.Status = types.CheckFail
		lastCheck.Message = "Some validations failed"
		format.CheckSummaryFail(passedCount, failedCount)
		format.CommandHint("Fix the issues and run", "quest check")
	}

	// Check if user wants AI-powered annotations
	annotate, _ := cmd.Flags().GetBool("annotate")
	if annotate {
		workDir, err := os.Getwd()
		if err == nil {
			fmt.Println()

			// Collect files to analyze
			var filesToAnalyze []string
			for _, artifact := range currentTask.Artifacts {
				fullPath := filepath.Join(workDir, artifact)
				if _, err := os.Stat(fullPath); err == nil {
					filesToAnalyze = append(filesToAnalyze, artifact)
				}
			}

			if len(filesToAnalyze) > 0 {
				// Start spinner while AI analyzes code
				spinner := format.NewSpinner("AI analyzing your code for contextual feedback...")
				spinner.Start()

				annotations, err := copilot_helper.GenerateCheckAnnotations(
					currentTask.Title,
					currentTaskValidation.Rules,
					filesToAnalyze,
				)

				spinner.Stop()

				if err != nil {
					format.Warning("Failed to generate annotations: " + err.Error())
				} else if len(annotations) > 0 {
					err = copilot_helper.ApplyCheckAnnotations(annotations, workDir)
					if err != nil {
						format.Warning("Failed to add annotations: " + err.Error())
					} else {
						format.AnnotationSummary(len(annotations))

						issueLevel := 0
						for _, ann := range annotations {
							switch ann.Type {
							case "error":
								lastCheck.Status = types.CheckFail
							case "warning":
								if issueLevel < 2 {
									lastCheck.Status = types.CheckWarn
								}
							case "info":
								if issueLevel < 1 {
									lastCheck.Status = types.CheckPass
								}
							}
						}
						// Show what was added
						for _, ann := range annotations {
							format.FileLocation(ann.File, ann.Line)
						}
						format.Newline()
					}
				} else {
					format.AnnotationSummary(0)
				}
			}
		}
	}

	state.LastCheck = lastCheck

	err = UploadStateAndPlan(state, plan)
	if err != nil {
		format.ErrorWithTip("Failed to save quest data", err, "Check folder permissions")
		return
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// CountFilesMatching returns the number of files matching a glob pattern
func CountFilesMatching(pattern string) (int, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0, err
	}
	return len(matches), nil
}

func RunCompete(cmd *cobra.Command, args []string) {
	state, plan, err := LoadStateAndPlan()
	if err != nil {
		format.ErrorWithTip("Failed to load quest data", err, "Run 'quest begin' to start a new quest")
		return
	}

	// Count completed vs total tasks
	totalTasks := countTasks(*plan)
	completedTasks := len(state.CompletedTaskIDs)
	allCompleted := ifCompletedPlan(*plan, *state)

	// Confirm with user
	confirmed, cancelled, err := prompt.ConfirmComplete(completedTasks, totalTasks, allCompleted)
	if err != nil {
		format.Error("Error reading input", err)
		return
	}

	if cancelled || !confirmed {
		format.Warning("Quest completion cancelled")
		return
	}

	// User confirmed - remove quest folder
	err = os.RemoveAll(QuestFolderName)
	if err != nil {
		format.ErrorWithTip("Failed to remove .quest folder", err, "Try manually deleting the .quest folder")
		return
	}

	if allCompleted {
		format.Success("🎉 Quest completed! All files cleaned up")
	} else {
		format.Success("Quest marked as complete. Files cleaned up")
		format.Printf("You completed %d of %d tasks.\n\n", completedTasks, totalTasks)
	}
}

func countTasks(plan types.Plan) int {
	countTasks := 0
	for _, chapter := range plan.Chapters {
		for _, quest := range chapter.Quests {
			countTasks += len(quest.Tasks)
		}
	}
	return countTasks
}

func selectTemplateForSpec(spec *types.ProjectSpec) string {
	if spec == nil {
		return "go-web-api"
	}
	if spec.Template != "" {
		return spec.Template
	}

	switch spec.BuildType {
	case types.BuildCLI, types.BuildAutomation, types.BuildFormatter, types.BuildLinter, types.BuildScaffolder, types.BuildProtocol, types.BuildClientSDK:
		return "go-cli-tool"
	case types.BuildStreamProcessor, types.BuildWorker, types.BuildPipeline, types.BuildScheduler, types.BuildAggregator, types.BuildIndexer:
		return "go-concurrency"
	case types.BuildService, types.BuildProxy, types.BuildLoadBalancer:
		if spec.Theme == types.Todo && spec.Difficulty != types.DifficultyQuick {
			return "go-todo-api"
		}
		if spec.Difficulty == types.DifficultyDeep {
			return "go-auth-system"
		}
		return "go-web-api"
	default:
		if spec.Difficulty == types.DifficultyDeep {
			return "go-auth-system"
		}
		return "go-web-api"
	}
}

func RunJumpTo(cmd *cobra.Command, args []string) {
	state, plan, err := LoadStateAndPlan()
	if err != nil {
		format.ErrorWithTip("Failed to load quest data", err, "Run 'quest begin' to start a new quest")
		return
	}

	var taskIndex int
	lastComplete, _ := cmd.Flags().GetBool("last-complete")

	if lastComplete {
		// Find the last completed task index
		if len(state.CompletedTaskIDs) == 0 {
			format.Warning("No completed tasks yet")
			return
		}

		// Find highest completed task index
		allTasks := FlattenTasks(plan)
		maxCompletedIndex := -1
		for i, task := range allTasks {
			if contains(state.CompletedTaskIDs, task.ID) && i > maxCompletedIndex {
				maxCompletedIndex = i
			}
		}

		if maxCompletedIndex < 0 {
			format.Warning("Could not find last completed task")
			return
		}
		taskIndex = maxCompletedIndex
	} else {
		// Parse task index from args
		if len(args) == 0 {
			format.ErrorWithTip("Missing task index", nil, "Usage: quest jump-to <task-number>")
			return
		}

		taskNum, err := strconv.Atoi(args[0])
		if err != nil {
			format.ErrorWithTip("Invalid task index", err, "Task index must be a number")
			return
		}

		// Convert from 1-based to 0-based index
		taskIndex = taskNum - 1
	}

	// Validate bounds
	if taskIndex < 0 || taskIndex >= plan.NumberOfTasks {
		format.Error(fmt.Sprintf("Task index out of range (1-%d)", plan.NumberOfTasks), nil)
		return
	}

	// Update state
	state.CurrentTaskIndex = taskIndex

	// Save state
	err = UploadState(state)
	if err != nil {
		format.ErrorWithTip("Failed to save state", err, "Changes not persisted")
		return
	}

	// Show confirmation
	allTasks := FlattenTasks(plan)
	currentTask := allTasks[taskIndex]
	format.Success(fmt.Sprintf("Jumped to Task %d: %s", taskIndex+1, currentTask.Title))
	format.CommandHint("View task details with", "quest next")
}

func RunSummary(cmd *cobra.Command, args []string) {
	state, plan, err := LoadStateAndPlan()
	if err != nil {
		format.ErrorWithTip("Failed to load quest data", err, "Run 'quest begin' to start a new quest")
		return
	}
	printSummary(*plan, *state)
}

func RunHealthCheck(cmd *cobra.Command, args []string) {
	format.Header("Quest Health Check")
	format.Newline()

	allHealthy := true

	// Check .quest folder
	questFolder, err := checkExistenceOfFile(QuestFolderName)
	if questFolder {
		format.CheckPass(".quest folder", "exists")
	} else {
		format.CheckFail(".quest folder", "not found - run 'quest begin' to start")
		allHealthy = false
	}

	// Check plan file
	planFile, err := checkExistenceOfFile(PlanFilePath)
	if planFile {
		format.CheckPass("plan.json", "exists")
	} else {
		format.CheckFail("plan.json", "not found")
		allHealthy = false
	}

	// Check state file
	stateFile, err := checkExistenceOfFile(StateFilePath)
	if stateFile {
		format.CheckPass("state.json", "exists")
	} else {
		format.CheckFail("state.json", "not found")
		allHealthy = false
	}

	// Check if plan can be loaded
	plan, err := LoadPlan()
	planLoad := plan != nil
	if err != nil {
		format.CheckFail("plan loading", fmt.Sprintf("failed to load: %s", err.Error()))
		allHealthy = false
	} else if planLoad {
		format.CheckPass("plan loading", fmt.Sprintf("loaded successfully (%d tasks)", plan.NumberOfTasks))
	} else {
		format.CheckFail("plan loading", "plan is nil")
		allHealthy = false
	}

	// Check if state can be loaded
	state, err := LoadState()
	stateLoad := state != nil
	if err != nil {
		format.CheckFail("state loading", fmt.Sprintf("failed to load: %s", err.Error()))
		allHealthy = false
	} else if stateLoad {
		format.CheckPass("state loading", fmt.Sprintf("loaded successfully (task %d/%d)", state.CurrentTaskIndex+1, plan.NumberOfTasks))
	} else {
		format.CheckFail("state loading", "state is nil")
		allHealthy = false
	}

	format.Newline()
	format.Divider()
	format.Newline()

	// Warning checks (non-critical)
	format.Bold("Environment Checks:")
	format.Newline()

	// Check if in a git repository
	_, err = os.Stat(".git")
	if err != nil {
		format.Warning("⚠ Not in a git repository - version control recommended")
	} else {
		format.Info("✓ Git repository detected")
	}

	// Check if GitHub Copilot CLI is available
	_, err = os.Stat(os.ExpandEnv("$HOME/.copilot"))
	if err != nil {
		// Try alternative check via command
		format.Warning("⚠ GitHub Copilot CLI not detected - enhanced features may be limited")
	} else {
		format.Info("✓ GitHub Copilot CLI detected")
	}

	format.Newline()
	format.Divider()
	format.Newline()

	// Final summary
	if allHealthy {
		format.Success("🎉 All health checks passed!")
		format.CommandHint("Continue your quest with", "quest next")
	} else {
		format.Error("❌ Some health checks failed", nil)
		format.CommandHint("Try running", "quest begin")
	}
}

func printSummary(plan types.Plan, state types.State) {
	ifCompleted := ifCompletedPlan(plan, state)
	currentChapter := currentChapter(state.CurrentTaskIndex, plan)
	currentQuest := currentQuest(state.CurrentTaskIndex, plan)

	// Header
	format.Header("Quest Progress")

	// Journey name, description, and tier
	format.Line(fmt.Sprintf("%s%s%s", format.ColorBold, plan.Journey.Name, format.ColorReset))
	if plan.Journey.Description != "" {
		format.Line(fmt.Sprintf("%s%s%s", format.ColorDim, plan.Journey.Description, format.ColorReset))
	}

	tierLabel := ""
	if plan.NumberOfTasks <= 3 {
		tierLabel = "Quick"
	} else if plan.NumberOfTasks <= 10 {
		tierLabel = "Standard"
	} else {
		tierLabel = "Extended"
	}

	focusStr := strings.Join(plan.Journey.Focus, ", ")
	format.Line(fmt.Sprintf("%s%s • %s%s",
		format.ColorDim, tierLabel, focusStr, format.ColorReset))
	format.Newline()

	// If not completed, show current task context
	if !ifCompleted {
		allTasks := FlattenTasks(&plan)
		currentTask := allTasks[state.CurrentTaskIndex]
		chapter := plan.Chapters[currentChapter]
		quest := chapter.Quests[currentQuest]

		// Show all chapters with status (if multiple chapters)
		if len(plan.Chapters) > 1 {
			format.Bold("Chapters:")
			for chIdx, ch := range plan.Chapters {
				// Determine chapter status
				chapterTasks := []types.Task{}
				for _, q := range ch.Quests {
					chapterTasks = append(chapterTasks, q.Tasks...)
				}

				allChapterComplete := true
				for _, task := range chapterTasks {
					if !contains(state.CompletedTaskIDs, task.ID) {
						allChapterComplete = false
						break
					}
				}

				if allChapterComplete {
					format.Line(fmt.Sprintf("  %s✔%s Chapter %d: %s",
						format.ColorGreen, format.ColorReset, chIdx+1, ch.Title))
				} else if chIdx == currentChapter {
					format.Line(fmt.Sprintf("  %s→%s Chapter %d: %s",
						format.ColorPink, format.ColorReset, chIdx+1, ch.Title))
				} else {
					format.Line(fmt.Sprintf("  %s○%s Chapter %d: %s",
						format.ColorDim, format.ColorReset, chIdx+1, ch.Title))
				}
			}
			format.Newline()
		}

		// Show current chapter details
		format.Bold("Current Chapter:")
		format.Line(fmt.Sprintf("  %sChapter %d: %s%s",
			format.ColorPink, currentChapter+1, chapter.Title, format.ColorReset))
		format.Line(fmt.Sprintf("    %sQuest: %s%s",
			format.ColorDim, quest.Title, format.ColorReset))
		format.Newline()

		// Show all tasks in current quest with status
		for _, task := range quest.Tasks {
			taskID := task.ID
			if contains(state.CompletedTaskIDs, taskID) {
				// Completed task
				format.Line(fmt.Sprintf("      %s✔%s %s",
					format.ColorGreen, format.ColorReset, task.Title))
			} else if task.ID == currentTask.ID {
				// Current task
				format.Line(fmt.Sprintf("      %s→%s %s",
					format.ColorPink, format.ColorReset, task.Title))
			} else {
				// Pending task
				format.Line(fmt.Sprintf("      %s○%s %s",
					format.ColorDim, format.ColorReset, task.Title))
			}
		}
		format.Newline()
	}

	// Footer separator
	format.Divider()

	// Summary stats
	completedCount := len(state.CompletedTaskIDs)
	format.Line(fmt.Sprintf("%sCompleted: %d / %d tasks%s",
		format.ColorDim, completedCount, plan.NumberOfTasks, format.ColorReset))

	if !ifCompleted {
		format.Line(fmt.Sprintf("%sNext action:%s Run %squest check%s",
			format.ColorDim, format.ColorReset, format.ColorCyan, format.ColorReset))
	} else {
		format.Line(fmt.Sprintf("%s🎉 All tasks completed!%s Run %squest complete%s",
			format.ColorGreen, format.ColorReset, format.ColorCyan, format.ColorReset))
	}
	format.Newline()
}

func ifCompletedPlan(plan types.Plan, state types.State) bool {
	if plan.NumberOfTasks == len(state.CompletedTaskIDs) {
		return true
	}
	return false
}

func currentChapter(taskID int, p types.Plan) int {
	taskPerChapter := taskPerChapter(p)
	current_task_count := taskID + 1
	total_count := 0
	for i, count := range taskPerChapter {
		total_count += count
		if total_count >= current_task_count {
			return i
		}
	}
	return 0
}

func currentQuest(taskID int, p types.Plan) int {
	taskPerQuest := taskPerQuest(p)
	current_task_count := taskID + 1
	total_count := 0
	for i, count := range taskPerQuest {
		total_count += count
		if total_count >= current_task_count {
			return i
		}
	}
	return 0
}

func taskPerChapter(plan types.Plan) []int {
	taskPerChapter := []int{}
	for _, ch := range plan.Chapters {
		num_tasks := 0
		for _, q := range ch.Quests {
			num_tasks += len(q.Tasks)
		}
		taskPerChapter = append(taskPerChapter, num_tasks)
	}
	return taskPerChapter
}

func taskPerQuest(plan types.Plan) []int {
	taskPerQuest := []int{}
	for _, ch := range plan.Chapters {
		for _, q := range ch.Quests {
			taskPerQuest = append(taskPerQuest, len(q.Tasks))
		}
	}
	return taskPerQuest
}
