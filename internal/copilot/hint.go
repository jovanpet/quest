package copilot_helper

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Hint struct {
	File    string
	Line    int
	Comment string
}

// CopilotExecutor allows injecting custom Copilot command for testing
type CopilotExecutor func(prompt string) (string, error)

// DefaultCopilotExecutor runs the actual Copilot CLI
var DefaultCopilotExecutor CopilotExecutor = func(prompt string) (string, error) {
	cmd := exec.Command("copilot")
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Env = os.Environ()

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("copilot execution failed: %w", err)
	}

	return outBuf.String(), nil
}

// GenerateHints analyzes student code and returns structured hints
func GenerateHints(task, objective string, files []string, attemptCount int) ([]Hint, error) {
	return GenerateHintsWithExecutor(task, objective, files, attemptCount, DefaultCopilotExecutor)
}

// GenerateHintsWithExecutor allows custom executor (for testing)
func GenerateHintsWithExecutor(task, objective string, files []string, attemptCount int, executor CopilotExecutor) ([]Hint, error) {
	// Build context-aware prompt
	fileList := strings.Join(files, "\n- ")
	
	// Adjust hint level based on attempt count
	hintLevel := "MINIMAL"
	maxHints := 3
	guidance := "Do NOT write solutions or give away answers"
	
	if attemptCount == 2 {
		hintLevel = "MORE SPECIFIC"
		maxHints = 4
		guidance = "Be more specific about what's wrong. You can hint at function names or patterns to use"
	} else if attemptCount >= 3 {
		hintLevel = "DIRECT"
		maxHints = 5
		guidance = "Be very direct. Show partial code examples if needed. Student is clearly stuck"
	}

	prompt := fmt.Sprintf(`Task: %s
Objective: %s

Student files:
- %s

This is attempt #%d for hints on this task.

Review the student's code and provide %s guiding hints.

Output ONLY in this exact format (one hint per line):
<file>:<line>: // <hint comment>

Example:
main.go:15: // TODO: Consider what happens if input is empty
utils.go:42: // HINT: Edge case - what if the slice has duplicate values?

Rules for attempt #%d:
- %s
- Keep hints minimal and thought-provoking
- Focus on bugs, edge cases, or missing logic
- Maximum %d hints per file
- If code looks good, output nothing

Generate hints now:`, task, objective, fileList, attemptCount, hintLevel, attemptCount, guidance, maxHints)

	// Execute via injected executor
	output, err := executor(prompt)
	if err != nil {
		return nil, err
	}

	// Parse structured hints from output
	return ParseHints(output)
}

// parseHints extracts hints in format: "file.go:42: // comment"
func ParseHints(output string) ([]Hint, error) {
	// Regex: filepath:linenumber: // comment
	re := regexp.MustCompile(`^([^:]+):(\d+):\s*//\s*(.+)$`)

	var hints []Hint
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) == 4 {
			lineNum, err := strconv.Atoi(matches[2])
			if err != nil {
				continue // skip malformed lines
			}

			hints = append(hints, Hint{
				File:    matches[1],
				Line:    lineNum,
				Comment: strings.TrimSpace(matches[3]),
			})
		}
	}

	return hints, scanner.Err()
}

// ApplyHints inserts comments into student files at specified lines
func ApplyHints(hints []Hint, workDir string) error {
	// Group hints by file
	hintsByFile := make(map[string][]Hint)
	for _, hint := range hints {
		hintsByFile[hint.File] = append(hintsByFile[hint.File], hint)
	}

	// Apply hints to each file
	for file, fileHints := range hintsByFile {
		fullPath := filepath.Join(workDir, file)

		// Read file
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		lines := strings.Split(string(content), "\n")

		// Sort hints by line number (descending) to avoid offset issues
		// Insert from bottom to top
		for i := len(fileHints) - 1; i >= 0; i-- {
			hint := fileHints[i]

			if hint.Line < 1 || hint.Line > len(lines) {
				continue // skip invalid line numbers
			}

			// Insert comment above the target line
			idx := hint.Line - 1
			indent := GetIndentation(lines[idx])
			comment := fmt.Sprintf("%s// %s", indent, hint.Comment)

			// Insert the comment
			lines = append(lines[:idx], append([]string{comment}, lines[idx:]...)...)
		}

		// Write back
		newContent := strings.Join(lines, "\n")
		if err := os.WriteFile(fullPath, []byte(newContent), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", file, err)
		}
	}

	return nil
}
