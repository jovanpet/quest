package copilot_helper

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jovanpet/quest/internal/types"
)

type CheckAnnotation struct {
	File    string
	Line    int
	Comment string
	Type    string // "success", "warning", "error"
}

// GenerateCheckAnnotations uses AI to analyze code and create contextual comments
func GenerateCheckAnnotations(task string, rules []types.Rule, files []string) ([]CheckAnnotation, error) {
	// Build prompt with check results context
	var passedChecks []string
	var failedChecks []string

	for _, rule := range rules {
		checkDesc := fmt.Sprintf("- %s", rule.Name)
		if rule.LastState != nil && *rule.LastState == types.Pass {
			passedChecks = append(passedChecks, checkDesc)
		} else {
			failedChecks = append(failedChecks, checkDesc)
		}
	}

	fileList := strings.Join(files, "\n- ")
	passedList := strings.Join(passedChecks, "\n")
	failedList := strings.Join(failedChecks, "\n")

	prompt := fmt.Sprintf(`Task: %s

Files to review:
- %s

Checks that PASSED:
%s

Checks that FAILED:
%s

Review the student's actual code and provide inline feedback comments.

For PASSED checks:
- Add encouraging comments explaining WHY their code works
- Highlight good practices they used
- Use "// ✓ GOOD: <explanation>"

For FAILED checks:
- Point out EXACTLY what's missing or wrong in their code
- Be specific about what needs to be added/changed
- Use "// ✗ ERROR: <specific issue>" or "// ⚠ WARNING: <suggestion>"

Output ONLY in this format (one comment per line):
<file>:<line>: // <type>: <comment>

Examples:
main.go:5: // ✓ GOOD: Using http.HandleFunc correctly
main.go:12: // ✗ ERROR: Missing http.ListenAndServe call to start the server
utils.go:20: // ⚠ WARNING: Consider adding error handling here

Rules:
- Maximum 5 comments total across all files
- Be specific to their actual code
- Focus on what they wrote or didn't write
- Keep comments concise (under 80 chars)

Generate comments now:`, task, fileList, passedList, failedList)

	// Execute Copilot CLI
	cmd := exec.Command("copilot")
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Env = os.Environ()

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("copilot execution failed: %w", err)
	}

	output := outBuf.String()
	return parseAnnotations(output)
}

// parseAnnotations extracts annotations from AI output
func parseAnnotations(output string) ([]CheckAnnotation, error) {
	// Regex: filepath:linenumber: // comment
	re := regexp.MustCompile(`^([^:]+):(\d+):\s*//\s*(.+)$`)

	var annotations []CheckAnnotation
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) != 4 {
			continue
		}

		file := matches[1]
		lineNum := 0
		fmt.Sscanf(matches[2], "%d", &lineNum)
		comment := matches[3]

		// Determine type based on comment prefix
		annotationType := "info"
		if strings.Contains(comment, "✓") || strings.Contains(comment, "GOOD") {
			annotationType = "success"
		} else if strings.Contains(comment, "✗") || strings.Contains(comment, "ERROR") {
			annotationType = "error"
		} else if strings.Contains(comment, "⚠") || strings.Contains(comment, "WARNING") {
			annotationType = "warning"
		}

		annotations = append(annotations, CheckAnnotation{
			File:    file,
			Line:    lineNum,
			Comment: comment,
			Type:    annotationType,
		})
	}

	return annotations, nil
}

// ApplyCheckAnnotations adds inline comments to student code
func ApplyCheckAnnotations(annotations []CheckAnnotation, workDir string) error {
	// Group by file
	annotationsByFile := make(map[string][]CheckAnnotation)
	for _, ann := range annotations {
		annotationsByFile[ann.File] = append(annotationsByFile[ann.File], ann)
	}

	// Apply to each file
	for file, fileAnnotations := range annotationsByFile {
		fullPath := filepath.Join(workDir, file)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			continue
		}

		if err := applyAnnotationsToFile(fullPath, fileAnnotations); err != nil {
			return fmt.Errorf("failed to annotate %s: %w", file, err)
		}
	}

	return nil
}

// applyAnnotationsToFile inserts comments into a file
func applyAnnotationsToFile(filePath string, annotations []CheckAnnotation) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	// Remove any existing check comments
	lines = removeCheckComments(lines)

	// Sort annotations by line number (descending) to avoid offset issues
	for i := len(annotations) - 1; i >= 0; i-- {
		ann := annotations[i]

		if ann.Line < 1 || ann.Line > len(lines) {
			continue
		}

		// Insert comment above target line
		idx := ann.Line - 1
		indent := getIndentation(lines[idx])
		comment := fmt.Sprintf("%s// %s", indent, ann.Comment)

		lines = append(lines[:idx], append([]string{comment}, lines[idx:]...)...)
	}

	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}

// removeCheckComments removes existing check annotations
func removeCheckComments(lines []string) []string {
	checkCommentRegex := regexp.MustCompile(`^\s*//\s*[✓✗⚠]\s*(GOOD|ERROR|WARNING):`)

	var cleaned []string
	for _, line := range lines {
		if !checkCommentRegex.MatchString(line) {
			cleaned = append(cleaned, line)
		}
	}
	return cleaned
}

// getIndentation returns the leading whitespace
func getIndentation(line string) string {
	return line[:len(line)-len(strings.TrimLeft(line, " \t"))]
}
