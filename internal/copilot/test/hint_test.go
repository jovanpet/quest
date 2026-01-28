package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	copilot_helper "github.com/jovanpet/quest/internal/copilot"
)

func TestParseHints(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []copilot_helper.Hint
	}{
		{
			name:  "single hint",
			input: `main.go:15: // TODO: Consider edge case for empty input`,
			expected: []copilot_helper.Hint{
				{File: "main.go", Line: 15, Comment: "TODO: Consider edge case for empty input"},
			},
		},
		{
			name: "multiple hints",
			input: `main.go:15: // TODO: Consider edge case
utils.go:42: // HINT: Check for nil pointer`,
			expected: []copilot_helper.Hint{
				{File: "main.go", Line: 15, Comment: "TODO: Consider edge case"},
				{File: "utils.go", Line: 42, Comment: "HINT: Check for nil pointer"},
			},
		},
		{
			name: "hints with extra whitespace",
			input: `  main.go:20: //   Handle error case  
`,
			expected: []copilot_helper.Hint{
				{File: "main.go", Line: 20, Comment: "Handle error case"},
			},
		},
		{
			name: "mixed valid and invalid lines",
			input: `main.go:15: // Valid hint
Some random text
utils.go:not-a-number: // Invalid
auth.go:30: // Another valid hint`,
			expected: []copilot_helper.Hint{
				{File: "main.go", Line: 15, Comment: "Valid hint"},
				{File: "auth.go", Line: 30, Comment: "Another valid hint"},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []copilot_helper.Hint{},
		},
		{
			name:     "no valid hints",
			input:    "Just some random text\nNo hints here",
			expected: []copilot_helper.Hint{},
		},
		{
			name: "filepath with directory",
			input: `src/main.go:10: // TODO: Add validation
internal/utils/helper.go:25: // HINT: Optimize this loop`,
			expected: []copilot_helper.Hint{
				{File: "src/main.go", Line: 10, Comment: "TODO: Add validation"},
				{File: "internal/utils/helper.go", Line: 25, Comment: "HINT: Optimize this loop"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := copilot_helper.ParseHints(tt.input)
			if err != nil {
				t.Fatalf("parseHints returned error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d hints, got %d", len(tt.expected), len(result))
			}

			for i, hint := range result {
				if i >= len(tt.expected) {
					break
				}
				exp := tt.expected[i]
				if hint.File != exp.File || hint.Line != exp.Line || hint.Comment != exp.Comment {
					t.Errorf("hint %d mismatch:\nexpected: %+v\ngot: %+v", i, exp, hint)
				}
			}
		})
	}
}

func TestGetIndentation(t *testing.T) {
	tests := []struct {
		line     string
		expected string
	}{
		{"func main() {", ""},
		{"    fmt.Println()", "    "},
		{"\t\treturn nil", "\t\t"},
		{"  \t  mixed", "  \t  "},
		{"", ""},
		{"   ", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			result := copilot_helper.GetIndentation(tt.line)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestApplyHints(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Create test file
	testFile := "test.go"
	testContent := `package main

func main() {
    x := 10
    y := 20
    result := x + y
    println(result)
}
`
	testPath := filepath.Join(tempDir, testFile)
	if err := os.WriteFile(testPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Apply hints
	hints := []copilot_helper.Hint{
		{File: testFile, Line: 4, Comment: "TODO: Validate x is not zero"},
		{File: testFile, Line: 6, Comment: "HINT: Consider overflow"},
	}

	if err := copilot_helper.ApplyHints(hints, tempDir); err != nil {
		t.Fatalf("ApplyHints failed: %v", err)
	}

	// Read modified file
	modifiedContent, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("failed to read modified file: %v", err)
	}

	modified := string(modifiedContent)

	// Check that hints were added
	expectedLines := []string{
		"// TODO: Validate x is not zero",
		"// HINT: Consider overflow",
	}

	for _, expected := range expectedLines {
		if !strings.Contains(modified, expected) {
			t.Errorf("modified file missing expected line: %s\n\nFull content:\n%s", expected, modified)
		}
	}

	// Verify hints are on correct lines (with proper indentation)
	lines := strings.Split(modified, "\n")

	// Line 4 should now have the comment before "x := 10"
	// With the first hint inserted, line 4 becomes line 5
	if !strings.Contains(lines[3], "// TODO: Validate x is not zero") {
		t.Errorf("hint not at expected position. Line 4: %s", lines[3])
	}

	// Verify original code is still present
	if !strings.Contains(modified, "x := 10") {
		t.Error("original code 'x := 10' is missing")
	}
	if !strings.Contains(modified, "result := x + y") {
		t.Error("original code 'result := x + y' is missing")
	}
}

func TestApplyHintsPreservesIndentation(t *testing.T) {
	tempDir := t.TempDir()

	testFile := "indented.go"
	testContent := `package main

func nested() {
    if true {
        x := 42
    }
}
`
	testPath := filepath.Join(tempDir, testFile)
	if err := os.WriteFile(testPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	hints := []copilot_helper.Hint{
		{File: testFile, Line: 5, Comment: "Check this value"},
	}

	if err := copilot_helper.ApplyHints(hints, tempDir); err != nil {
		t.Fatalf("ApplyHints failed: %v", err)
	}

	modifiedContent, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("failed to read modified file: %v", err)
	}

	modified := string(modifiedContent)
	lines := strings.Split(modified, "\n")

	// Find the comment line
	var commentLine string
	for _, line := range lines {
		if strings.Contains(line, "Check this value") {
			commentLine = line
			break
		}
	}

	if commentLine == "" {
		t.Fatal("comment not found in modified file")
	}

	// Should have 8 spaces of indentation (matching "        x := 42")
	if !strings.HasPrefix(commentLine, "        // Check this value") {
		t.Errorf("indentation not preserved. Got: %q", commentLine)
	}
}

func TestApplyHintsInvalidLineNumber(t *testing.T) {
	tempDir := t.TempDir()

	testFile := "small.go"
	testContent := `package main
func main() {}
`
	testPath := filepath.Join(tempDir, testFile)
	if err := os.WriteFile(testPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	hints := []copilot_helper.Hint{
		{File: testFile, Line: 0, Comment: "Line 0 should be skipped"},
		{File: testFile, Line: 999, Comment: "Line 999 should be skipped"},
		{File: testFile, Line: 2, Comment: "Valid hint"},
	}

	if err := copilot_helper.ApplyHints(hints, tempDir); err != nil {
		t.Fatalf("ApplyHints failed: %v", err)
	}

	modifiedContent, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("failed to read modified file: %v", err)
	}

	modified := string(modifiedContent)

	// Only valid hint should be present
	if !strings.Contains(modified, "Valid hint") {
		t.Error("valid hint was not applied")
	}

	if strings.Contains(modified, "Line 0") || strings.Contains(modified, "Line 999") {
		t.Error("invalid hints were applied")
	}
}
