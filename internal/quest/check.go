package quest

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jovanpet/quest/internal/types"
)

func CheckRule(rule types.Rule) (bool, error) {
	switch rule.Type {
	case types.TypeExists:
		exists, err := checkExistenceOfFile(rule.Path)
		if !exists {
			return false, fmt.Errorf("file '%s' does not exist", rule.Path)
		}
		return exists, err
	case types.TypeGlobCountMin:
		count, err := countFilesBasedOnRegex(rule.Glob)
		if err != nil {
			return false, err
		}
		if count < rule.Min {
			return false, fmt.Errorf("found %d file(s) matching '%s', expected at least %d", count, rule.Glob, rule.Min)
		}
		return true, nil
	case types.TypeFileContainsAny:
		contains, err := checkFileContainsSpecificSnippets(rule.Glob, rule.Any)
		if !contains && err == nil {
			// Build a nice error message showing what patterns were expected
			if len(rule.Any) == 1 {
				return false, fmt.Errorf("%s doesn't contain: '%s'", rule.Glob, rule.Any[0])
			}
			return false, fmt.Errorf("%s doesn't contain any of: %v", rule.Glob, rule.Any)
		}
		return contains, err
	}
	return false, fmt.Errorf("The Type setting is invalid.")
}

func checkExistenceOfFile(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	return false, err
}

func getFilePathsBasedOnRegex(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	return matches, err
}

func countFilesBasedOnRegex(pattern string) (int, error) {
	matches, err := getFilePathsBasedOnRegex(pattern)
	if err != nil {
		return 0, err
	}
	return len(matches), nil
}

func checkFileContainsSpecificSnippets(globPattern string, patterns []string) (bool, error) {
	matches, err := getFilePathsBasedOnRegex(globPattern)
	if err != nil {
		return false, err
	}

	if len(matches) == 0 {
		return false, fmt.Errorf("no files found matching pattern '%s'", globPattern)
	}

	for _, filePath := range matches {
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		contentStr := string(content)

		for _, pattern := range patterns {
			regex, err := regexp.Compile(pattern)
			if err != nil {
				return false, err
			}

			if regex.MatchString(contentStr) {
				return true, nil
			}
		}
	}

	return false, nil
}

func createArtifact(artifactPath string) error {
	// Skip if file already exists
	if _, err := os.Stat(artifactPath); err == nil {
		return nil
	}

	// Create parent directories
	dir := filepath.Dir(artifactPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Create the file
	return os.WriteFile(artifactPath, []byte("// TODO: implement\n"), 0644)
}
