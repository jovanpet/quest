package quest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jovanpet/quest/internal/types"
)

func TestCheckExistenceOfFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	// Test file doesn't exist
	exists, err := checkExistenceOfFile(testFile)
	if exists {
		t.Error("Expected file to not exist")
	}
	if err == nil {
		t.Error("Expected error when file doesn't exist")
	}

	// Create file
	err = os.WriteFile(testFile, []byte("package main"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test file exists
	exists, err = checkExistenceOfFile(testFile)
	if !exists {
		t.Error("Expected file to exist")
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCountFilesBasedOnRegex(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		files         []string
		pattern       string
		expectedCount int
	}{
		{
			name:          "no files",
			files:         []string{},
			pattern:       "*.go",
			expectedCount: 0,
		},
		{
			name:          "single match",
			files:         []string{"main.go"},
			pattern:       "*.go",
			expectedCount: 1,
		},
		{
			name:          "multiple matches",
			files:         []string{"main.go", "test.go", "helper.go"},
			pattern:       "*.go",
			expectedCount: 3,
		},
		{
			name:          "no matches",
			files:         []string{"main.go", "test.go"},
			pattern:       "*.txt",
			expectedCount: 0,
		},
		{
			name:          "test files only",
			files:         []string{"main.go", "main_test.go", "helper_test.go"},
			pattern:       "*_test.go",
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test directory for this test
			testDir := filepath.Join(tmpDir, tt.name)
			err := os.MkdirAll(testDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create test dir: %v", err)
			}

			// Create test files
			for _, file := range tt.files {
				filePath := filepath.Join(testDir, file)
				err := os.WriteFile(filePath, []byte("package main"), 0644)
				if err != nil {
					t.Fatalf("Failed to create file %s: %v", file, err)
				}
			}

			// Test count
			pattern := filepath.Join(testDir, tt.pattern)
			count, err := countFilesBasedOnRegex(pattern)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if count != tt.expectedCount {
				t.Errorf("Expected count %d, got %d", tt.expectedCount, count)
			}
		})
	}
}

func TestCheckFileContainsSpecificSnippets(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		fileContent string
		patterns    []string
		expected    bool
	}{
		{
			name:        "exact match",
			fileContent: "package main\n\nfunc main() {\n\thttp.ListenAndServe(\":8080\", nil)\n}",
			patterns:    []string{"http.ListenAndServe"},
			expected:    true,
		},
		{
			name:        "no match",
			fileContent: "package main\n\nfunc main() {}",
			patterns:    []string{"http.ListenAndServe"},
			expected:    false,
		},
		{
			name:        "multiple patterns - first matches",
			fileContent: "package main\n\nimport \"net/http\"",
			patterns:    []string{"net/http", "fmt"},
			expected:    true,
		},
		{
			name:        "multiple patterns - second matches",
			fileContent: "package main\n\nimport \"fmt\"",
			patterns:    []string{"net/http", "fmt"},
			expected:    true,
		},
		{
			name:        "regex pattern",
			fileContent: "func handleRequest() {}",
			patterns:    []string{"func \\w+\\("},
			expected:    true,
		},
		{
			name:        "regex pattern no match",
			fileContent: "const value = 10",
			patterns:    []string{"func \\w+\\("},
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tmpDir, tt.name+".go")
			err := os.WriteFile(testFile, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Test
			result, err := checkFileContainsSpecificSnippets(testFile, tt.patterns)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCheckFileContainsSpecificSnippetsWithGlob(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple files
	files := map[string]string{
		"main.go":   "package main\n\nfunc main() {}",
		"server.go": "package main\n\nfunc startServer() {\n\thttp.ListenAndServe(\":8080\", nil)\n}",
		"utils.go":  "package main\n\nfunc helper() {}",
	}

	for name, content := range files {
		filePath := filepath.Join(tmpDir, name)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", name, err)
		}
	}

	tests := []struct {
		name     string
		glob     string
		patterns []string
		expected bool
	}{
		{
			name:     "pattern found in one file",
			glob:     filepath.Join(tmpDir, "*.go"),
			patterns: []string{"http.ListenAndServe"},
			expected: true,
		},
		{
			name:     "pattern not found in any file",
			glob:     filepath.Join(tmpDir, "*.go"),
			patterns: []string{"database.Connect"},
			expected: false,
		},
		{
			name:     "specific file glob",
			glob:     filepath.Join(tmpDir, "server.go"),
			patterns: []string{"startServer"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := checkFileContainsSpecificSnippets(tt.glob, tt.patterns)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCheckRule(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tmpDir, "main.go")
	err := os.WriteFile(testFile, []byte("package main\n\nfunc main() {\n\thttp.ListenAndServe(\":8080\", nil)\n}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		rule     types.Rule
		expected bool
		wantErr  bool
	}{
		{
			name: "exists - file exists",
			rule: types.Rule{
				Type: types.TypeExists,
				Path: testFile,
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "exists - file doesn't exist",
			rule: types.Rule{
				Type: types.TypeExists,
				Path: filepath.Join(tmpDir, "nonexistent.go"),
			},
			expected: false,
			wantErr:  true,
		},
		{
			name: "glob_count_min - passes",
			rule: types.Rule{
				Type: types.TypeGlobCountMin,
				Glob: filepath.Join(tmpDir, "*.go"),
				Min:  1,
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "glob_count_min - fails",
			rule: types.Rule{
				Type: types.TypeGlobCountMin,
				Glob: filepath.Join(tmpDir, "*.go"),
				Min:  5,
			},
			expected: false,
			wantErr:  false,
		},
		{
			name: "file_contains_any - passes",
			rule: types.Rule{
				Type: types.TypeFileContainsAny,
				Path: testFile,
				Any:  []string{"http.ListenAndServe"},
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "file_contains_any - fails",
			rule: types.Rule{
				Type: types.TypeFileContainsAny,
				Path: testFile,
				Any:  []string{"database.Connect"},
			},
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CheckRule(tt.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("CheckRule() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCreateArtifact(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		artifact   string
		expectFile bool
	}{
		{
			name:       "create file with parent dir",
			artifact:   "handlers/health.go",
			expectFile: true,
		},
		{
			name:       "create file in root",
			artifact:   "main.go",
			expectFile: true,
		},
		{
			name:       "create nested file structure",
			artifact:   "internal/handlers/health.go",
			expectFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a subdirectory for each test to avoid conflicts
			testDir := filepath.Join(tmpDir, tt.name)
			os.MkdirAll(testDir, 0755)
			
			artifactPath := filepath.Join(testDir, tt.artifact)
			
			err := createArtifact(artifactPath)
			if err != nil {
				t.Fatalf("createArtifact() error = %v", err)
			}

			// Check if file was created
			if tt.expectFile {
				info, err := os.Stat(artifactPath)
				if err != nil {
					t.Errorf("File was not created: %v", err)
					return
				}
				if info.IsDir() {
					t.Error("Expected file, got directory")
				}

				// Check parent directory exists
				parentDir := filepath.Dir(artifactPath)
				if _, err := os.Stat(parentDir); err != nil {
					t.Errorf("Parent directory was not created: %v", err)
				}
			}

			// Test idempotency - calling again should not error
			err = createArtifact(artifactPath)
			if err != nil {
				t.Errorf("createArtifact() should be idempotent, got error: %v", err)
			}
		})
	}
}
