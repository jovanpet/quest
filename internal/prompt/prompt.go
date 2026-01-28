package prompt

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jovanpet/quest/internal/format"
	"github.com/jovanpet/quest/internal/types"
)

// WizardMode represents the initial choice
type WizardMode string

const (
	ModeTemplate WizardMode = "template"
	ModeForge    WizardMode = "forge"
	ModeSurprise WizardMode = "surprise"
)

// SelectWizardMode prompts for the initial wizard choice
func SelectWizardMode() (WizardMode, bool, error) {
	lines := []string{
		fmt.Sprintf("%s[1]%s %sPick a Legendary Path%s",
			format.ColorPink, format.ColorReset,
			format.ColorBold, format.ColorReset),
		fmt.Sprintf("    %sChoose from pre-made templates with guided tasks%s", format.ColorDim, format.ColorReset),
		"",
		fmt.Sprintf("%s[2]%s %sForge Your Own Quest%s",
			format.ColorPink, format.ColorReset,
			format.ColorBold, format.ColorReset),
		fmt.Sprintf("    %sCustomize difficulty, build type, and theme%s", format.ColorDim, format.ColorReset),
		"",
		fmt.Sprintf("%s[3]%s %sSeek a Mystery Quest%s",
			format.ColorPink, format.ColorReset,
			format.ColorBold, format.ColorReset),
		fmt.Sprintf("    %sRandom selection - surprise me with something fun!%s", format.ColorDim, format.ColorReset),
	}
	
	format.Box("✨ Begin Your Journey", lines)
	
	reader := bufio.NewReader(os.Stdin)
	format.Prompt("Choose your destiny (or 'q' to retreat)")
	
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", false, err
	}
	
	input = strings.TrimSpace(input)
	
	if input == "q" || input == "Q" {
		format.Newline()
		return "", true, nil
	}
	
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > 3 {
		format.Warning("Invalid choice")
		return "", true, nil
	}
	
	switch choice {
	case 1:
		return ModeTemplate, false, nil
	case 2:
		return ModeForge, false, nil
	case 3:
		return ModeSurprise, false, nil
	default:
		return "", true, nil
	}
}

// ForgeQuestWizard guides user through creating custom quest
func ForgeQuestWizard() (*types.ProjectSpec, bool, error) {
	spec := &types.ProjectSpec{}
	
	// Step 1: Select BuildType
	buildType, cancelled, err := selectBuildType()
	if err != nil || cancelled {
		return nil, cancelled, err
	}
	spec.BuildType = buildType
	
	// Step 2: Select Difficulty
	difficulty, cancelled, err := selectDifficulty()
	if err != nil || cancelled {
		return nil, cancelled, err
	}
	spec.Difficulty = difficulty
	
	// Step 3: Optional Theme
	theme, cancelled, err := selectTheme()
	if err != nil || cancelled {
		return nil, cancelled, err
	}
	spec.Theme = theme
	
	return spec, false, nil
}

// SurpriseQuest generates random quest spec
func SurpriseQuest() *types.ProjectSpec {
	buildTypes := []types.BuildType{
		types.BuildService, types.BuildCLI, types.BuildWorker, types.BuildLibrary,
		types.BuildScheduler, types.BuildProxy, types.BuildStreamProcessor,
	}
	
	difficulties := []types.Difficulty{
		types.DifficultyQuick, types.DifficultyNormal, types.DifficultyDeep,
	}
	
	themes := []types.Theme{
		types.Todo, types.Notes, types.Bookmarks, types.Expenses,
		types.Contacts, types.Events, types.Inventory,
	}
	
	rand.Seed(time.Now().UnixNano())
	
	return &types.ProjectSpec{
		BuildType:  buildTypes[rand.Intn(len(buildTypes))],
		Difficulty: difficulties[rand.Intn(len(difficulties))],
		Theme:      themes[rand.Intn(len(themes))],
		SurpriseMe: true,
	}
}

func selectBuildType() (types.BuildType, bool, error) {
	format.Newline()
	
	buildTypes := []struct {
		Type        types.BuildType
		Description string
	}{
		{types.BuildService, "Web API / REST backend with HTTP endpoints"},
		{types.BuildCLI, "Command-line tool with flags and arguments"},
		{types.BuildPipeline, "Data pipeline / ETL for batch processing"},
		{types.BuildWorker, "Background job processor / message consumer"},
		{types.BuildLibrary, "Reusable Go package / module"},
		{types.BuildScheduler, "Task scheduler / cron system"},
		{types.BuildProxy, "Reverse proxy / API gateway"},
		{types.BuildStreamProcessor, "Real-time event / stream processor"},
	}
	
	lines := []string{}
	for i, bt := range buildTypes {
		line := fmt.Sprintf("%s[%d]%s %s%s%s",
			format.ColorPink, i+1, format.ColorReset,
			format.ColorBold, bt.Type, format.ColorReset)
		lines = append(lines, line)
		lines = append(lines, fmt.Sprintf("    %s%s%s", format.ColorDim, bt.Description, format.ColorReset))
		if i < len(buildTypes)-1 {
			lines = append(lines, "")
		}
	}
	
	format.Box("🔨 Choose Your Creation", lines)
	
	reader := bufio.NewReader(os.Stdin)
	format.Prompt("Select your craft (or 'q' to retreat)")
	
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", false, err
	}
	
	input = strings.TrimSpace(input)
	if input == "q" || input == "Q" {
		return "", true, nil
	}
	
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(buildTypes) {
		format.Warning("Invalid choice")
		return "", true, nil
	}
	
	return buildTypes[choice-1].Type, false, nil
}

func selectDifficulty() (types.Difficulty, bool, error) {
	format.Newline()
	
	lines := []string{
		fmt.Sprintf("%s[1]%s %sApprentice (Quick)%s",
			format.ColorPink, format.ColorReset,
			format.ColorBold, format.ColorReset),
		fmt.Sprintf("    %s3-10 tasks, about 30 minutes%s",
			format.ColorDim, format.ColorReset),
		"",
		fmt.Sprintf("%s[2]%s %sAdventurer (Normal)%s",
			format.ColorPink, format.ColorReset,
			format.ColorBold, format.ColorReset),
		fmt.Sprintf("    %s10-20 tasks, 1-2 hours to complete%s",
			format.ColorDim, format.ColorReset),
		"",
		fmt.Sprintf("%s[3]%s %sMaster (Deep)%s",
			format.ColorPink, format.ColorReset,
			format.ColorBold, format.ColorReset),
		fmt.Sprintf("    %s20-35 tasks, comprehensive 3-4 hour journey%s",
			format.ColorDim, format.ColorReset),
	}
	
	format.Box("⚔️ Select Your Challenge", lines)
	
	reader := bufio.NewReader(os.Stdin)
	format.Prompt("Choose difficulty rank (or 'q' to retreat)")
	
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", false, err
	}
	
	input = strings.TrimSpace(input)
	if input == "q" || input == "Q" {
		return "", true, nil
	}
	
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > 3 {
		format.Warning("Invalid choice")
		return "", true, nil
	}
	
	switch choice {
	case 1:
		return types.DifficultyQuick, false, nil
	case 2:
		return types.DifficultyNormal, false, nil
	case 3:
		return types.DifficultyDeep, false, nil
	default:
		return "", true, nil
	}
}

func selectTheme() (types.Theme, bool, error) {
	format.Newline()
	
	themeDescriptions := []struct {
		Theme types.Theme
		Desc  string
	}{
		{types.Todo, "Task management system"},
		{types.Notes, "Note-taking application"},
		{types.Bookmarks, "Bookmark manager"},
		{types.Expenses, "Expense tracker"},
		{types.Contacts, "Contact management"},
		{types.Events, "Event scheduler"},
		{types.Inventory, "Inventory system"},
	}
	
	lines := []string{}
	for i, td := range themeDescriptions {
		line := fmt.Sprintf("%s[%d]%s %s%s%s - %s%s%s",
			format.ColorPink, i+1, format.ColorReset,
			format.ColorBold, td.Theme, format.ColorReset,
			format.ColorDim, td.Desc, format.ColorReset)
		lines = append(lines, line)
	}
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("%s[0]%s %sNo theme - pure technical focus%s",
		format.ColorDim, format.ColorReset,
		format.ColorDim, format.ColorReset))
	
	format.Box("🎨 Choose Your Realm", lines)
	
	reader := bufio.NewReader(os.Stdin)
	format.Prompt("Select theme (or 0 to skip)")
	
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", false, err
	}
	
	input = strings.TrimSpace(input)
	if input == "q" || input == "Q" {
		return "", true, nil
	}
	
	if input == "0" || input == "" {
		return "", false, nil // No theme
	}
	
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 0 || choice > len(themeDescriptions) {
		format.Warning("Invalid choice")
		return "", true, nil
	}
	
	return themeDescriptions[choice-1].Theme, false, nil
}

// SelectTemplate prompts user to select a template
func SelectTemplate(templates []string) (string, bool, error) {
	format.Newline()
	
	// Create box content
	lines := []string{}
	for i, template := range templates {
		parts := strings.SplitN(template, " - ", 2)
		name := parts[0]
		title := ""
		if len(parts) > 1 {
			title = parts[1]
		}
		
		line := fmt.Sprintf("%s[%d]%s %s%s%s",
			format.ColorCyan, i+1, format.ColorReset,
			format.ColorBold, name, format.ColorReset)
		lines = append(lines, line)
		
		if title != "" {
			lines = append(lines, fmt.Sprintf("    %s%s%s", format.ColorDim, title, format.ColorReset))
		}
		
		if i < len(templates)-1 {
			lines = append(lines, "") // Spacer
		}
	}
	
	format.Box("🎮 Choose Your Quest", lines)
	
	reader := bufio.NewReader(os.Stdin)
	format.Prompt("Enter number (or 'q' to cancel)")
	
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", false, err
	}
	
	input = strings.TrimSpace(input)
	
	if input == "q" || input == "Q" {
		format.Newline()
		return "", true, nil
	}
	
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(templates) {
		format.Warning("Invalid choice")
		return "", true, nil
	}
	
	// Extract just the template name (before the " - ")
	selected := templates[choice-1]
	parts := strings.SplitN(selected, " - ", 2)
	templateName := parts[0]
	
	// Show a quick loading animation
	format.Newline()
	showLoading("Loading template", 500)
	
	return templateName, false, nil
}

// ConfirmComplete prompts user to confirm quest completion
func ConfirmComplete(completedTasks, totalTasks int, allCompleted bool) (bool, bool, error) {
	format.Newline()
	
	var lines []string
	if allCompleted {
		lines = append(lines, fmt.Sprintf("%s🎉 All %d tasks completed!%s", 
			format.ColorGreen, totalTasks, format.ColorReset))
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("%sThis will remove the .quest folder%s", 
			format.ColorDim, format.ColorReset))
	} else {
		lines = append(lines, fmt.Sprintf("%sProgress: %d/%d tasks completed%s", 
			format.ColorYellow, completedTasks, totalTasks, format.ColorReset))
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("%sThis will remove the .quest folder%s", 
			format.ColorDim, format.ColorReset))
	}
	
	title := "Complete Quest"
	if !allCompleted {
		title = "⚠️  Complete Quest Early?"
	}
	
	format.Box(title, lines)
	
	format.Prompt("Confirm completion? (y/N)")
	
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, false, err
	}
	
	input = strings.TrimSpace(strings.ToLower(input))
	
	if input == "y" || input == "yes" {
		format.Newline()
		showLoading("Cleaning up", 400)
		return true, false, nil
	}
	
	format.Newline()
	return false, true, nil
}

// showLoading shows a quick loading animation
func showLoading(msg string, durationMs int) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	duration := time.Duration(durationMs) * time.Millisecond
	frameDelay := 80 * time.Millisecond
	iterations := int(duration / frameDelay)
	
	for i := 0; i < iterations; i++ {
		frame := frames[i%len(frames)]
		format.Printf("\r  %s%s%s %s...%s", format.ColorPink, frame, format.ColorReset, msg, format.ColorReset)
		time.Sleep(frameDelay)
	}
	format.Printf("\r  %s%s%s %s    \n", format.ColorGreen, format.BoxCheck, format.ColorReset, msg)
}

