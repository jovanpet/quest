package copilot_helper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jovanpet/quest/internal/types"
)

// GeneratePlan creates a quest plan from user specifications
func GeneratePlan(spec types.ProjectSpec) (*types.Plan, error) {
	return GeneratePlanWithExecutor(spec, DefaultCopilotExecutor)
}

// GeneratePlanWithExecutor allows custom executor (for testing)
func GeneratePlanWithExecutor(spec types.ProjectSpec, executor CopilotExecutor) (*types.Plan, error) {
	// Build the AI prompt
	prompt := buildPlanPrompt(spec)

	// Execute via Copilot
	output, err := executor(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate plan: %w", err)
	}

	// Parse JSON response into Plan
	plan, err := parsePlanJSON(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated plan: %w", err)
	}

	return plan, nil
}

func buildPlanPrompt(spec types.ProjectSpec) string {
	// Determine task count based on difficulty
	taskRange := ""
	switch spec.Difficulty {
	case types.DifficultyQuick:
		taskRange = "3-10 tasks (aim for ~5-7)"
	case types.DifficultyNormal:
		taskRange = "10-20 tasks (aim for ~12-15)"
	case types.DifficultyDeep:
		taskRange = "20-35 tasks (aim for ~25-30)"
	}

	themeContext := ""
	if spec.Theme != "" {
		themeContext = fmt.Sprintf("\nTheme: %s - incorporate this domain into the project (e.g., build a %s %s)",
			spec.Theme, spec.Theme, spec.BuildType)
	}

	return fmt.Sprintf(`You are a coding quest designer. Generate a structured learning quest for a student.

PROJECT SPECIFICATION:
- Build Type: %s
- Difficulty: %s (%s)%s
- Language: Go

YOUR TASK:
Create a complete quest plan in JSON format following this EXACT structure. The quest should teach best practices and gradually increase in complexity.

STRUCTURE GUIDELINES:
1. Organize into logical chapters (1-3 chapters depending on difficulty)
2. Each chapter has 1-3 quests
3. Each quest has multiple tasks (bite-sized, achievable steps)
4. Total tasks: %s
5. Each task must have:
   - Unique ID (e.g., "task-1-1-1")
   - Clear title (e.g., "Set up Go module")
   - Objective (what student will learn/build)
   - Steps (array of actionable instructions)
   - Files (these are the files that the program will automatically create)
   - Artifacts (all files involved in this task that the student will work on)
   - Validation rules (how to check success)

VALIDATION RULE TYPES:
1. "exists" - Check if file/folder exists
   Example: {"type": "exists", "name": "File exists", "path": "main.go"}

2. "glob_count_min" - Count files matching pattern
   Example: {"type": "glob_count_min", "name": "Has test files", "glob": "*_test.go", "min": 1}

3. "file_contains_any" - Check file contains specific strings
   Example: {"type": "file_contains_any", "name": "Uses HTTP handler", "glob": "*.go", "any": ["http.HandleFunc", "http.Handler"]}

OUTPUT FORMAT (strict JSON):
{
  "version": 1,
  "journey": {
    "name": "Build a [Type]",
	"description": "Description of the project",
    "language": "Go",
    "focus": ["concept1", "concept2", "concept3"]
  },
  "chapters": [
    {
      "id": "chapter-1",
      "title": "Chapter Title",
      "quests": [
        {
          "id": "quest-1-1",
          "title": "Quest Title",
          "tasks": [
            {
              "id": "task-1-1-1",
              "title": "Task Title",
              "objective": "What the student will accomplish",
              "steps": [
                "Step 1 instruction",
                "Step 2 instruction"
              ],
              "files": ["main.go"],
			  "artifacts": ["main.go"]
              "validation": {
                "rules": [
                  {
                    "type": "exists",
                    "name": "Main file created",
                    "path": "main.go"
                  }
                ]
              }
            }
          ]
        }
      ]
    }
  ]
}

IMPORTANT RULES:
- Start simple, increase complexity gradually
- Each task should be completable in 5-10 minutes
- Include realistic validation rules
- Focus on practical, working code
- Include error handling, testing, and best practices
- Make it educational but fun
- Use proper Go project structure (cmd/, internal/, pkg/ if needed)
- Include at least 2 validation rules per task

Generate the complete quest plan now as valid JSON:`,
		spec.BuildType,
		spec.Difficulty,
		taskRange,
		themeContext,
		taskRange)
}

func parsePlanJSON(output string) (*types.Plan, error) {
	// Clean up the output - extract JSON if wrapped in markdown
	jsonStr := extractJSON(output)

	// Decode Unicode escape sequences (like \u0026 -> &)
	jsonStr = unescapeUnicode(jsonStr)

	var plan types.Plan
	err := json.Unmarshal([]byte(jsonStr), &plan)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Validate required fields
	if plan.Journey.Name == "" {
		return nil, fmt.Errorf("plan missing journey name")
	}
	if len(plan.Chapters) == 0 {
		return nil, fmt.Errorf("plan has no chapters")
	}
	return &plan, nil
}

// extractJSON attempts to find and extract JSON from markdown or other wrappers
func extractJSON(output string) string {
	// Remove markdown code blocks if present
	output = strings.TrimSpace(output)

	// Look for ```json ... ``` blocks
	if strings.Contains(output, "```json") {
		start := strings.Index(output, "```json")
		if start != -1 {
			start += 7 // Skip past ```json
			end := strings.Index(output[start:], "```")
			if end != -1 {
				return strings.TrimSpace(output[start : start+end])
			}
		}
	}

	// Look for generic ``` blocks
	if strings.HasPrefix(output, "```") {
		lines := strings.Split(output, "\n")
		if len(lines) > 2 {
			// Skip first line (```json or ```) and last line (```)
			return strings.TrimSpace(strings.Join(lines[1:len(lines)-1], "\n"))
		}
	}

	// Find first { and last }
	start := strings.Index(output, "{")
	end := strings.LastIndex(output, "}")
	if start != -1 && end != -1 && end > start {
		return output[start : end+1]
	}

	return output
}

// unescapeUnicode converts Unicode escape sequences to actual characters
func unescapeUnicode(s string) string {
	// Replace common Unicode escapes that AI might output
	s = strings.ReplaceAll(s, `\u0026`, "&")
	s = strings.ReplaceAll(s, `\u003c`, "<")
	s = strings.ReplaceAll(s, `\u003e`, ">")
	s = strings.ReplaceAll(s, `\u0027`, "'")
	s = strings.ReplaceAll(s, `\u0022`, `"`)
	return s
}
