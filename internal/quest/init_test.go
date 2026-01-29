package quest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jovanpet/quest/internal/types"
)

func TestFlattenTasks(t *testing.T) {
	tests := []struct {
		name          string
		plan          *types.Plan
		expectedCount int
	}{
		{
			name: "empty plan",
			plan: &types.Plan{
				Chapters: []types.Chapter{},
			},
			expectedCount: 0,
		},
		{
			name: "single chapter with single quest and single task",
			plan: &types.Plan{
				Chapters: []types.Chapter{
					{
						Title: "Chapter 1",
						Quests: []types.Quest{
							{
								Title: "Quest 1",
								Tasks: []types.Task{
									{ID: "task1", Title: "Task 1"},
								},
							},
						},
					},
				},
			},
			expectedCount: 1,
		},
		{
			name: "multiple chapters with multiple quests and tasks",
			plan: &types.Plan{
				Chapters: []types.Chapter{
					{
						Title: "Chapter 1",
						Quests: []types.Quest{
							{
								Title: "Quest 1",
								Tasks: []types.Task{
									{ID: "task1", Title: "Task 1"},
									{ID: "task2", Title: "Task 2"},
								},
							},
							{
								Title: "Quest 2",
								Tasks: []types.Task{
									{ID: "task3", Title: "Task 3"},
								},
							},
						},
					},
					{
						Title: "Chapter 2",
						Quests: []types.Quest{
							{
								Title: "Quest 3",
								Tasks: []types.Task{
									{ID: "task4", Title: "Task 4"},
									{ID: "task5", Title: "Task 5"},
								},
							},
						},
					},
				},
			},
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks := FlattenTasks(tt.plan)
			if len(tasks) != tt.expectedCount {
				t.Errorf("Expected %d tasks, got %d", tt.expectedCount, len(tasks))
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "empty slice",
			slice:    []string{},
			item:     "test",
			expected: false,
		},
		{
			name:     "item exists",
			slice:    []string{"a", "b", "c"},
			item:     "b",
			expected: true,
		},
		{
			name:     "item doesn't exist",
			slice:    []string{"a", "b", "c"},
			item:     "d",
			expected: false,
		},
		{
			name:     "item at start",
			slice:    []string{"test", "foo", "bar"},
			item:     "test",
			expected: true,
		},
		{
			name:     "item at end",
			slice:    []string{"foo", "bar", "test"},
			item:     "test",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCountFilesMatching(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := []string{
		"main.go",
		"server.go",
		"utils.go",
		"main_test.go",
		"server_test.go",
		"readme.md",
	}

	for _, file := range files {
		filePath := filepath.Join(tmpDir, file)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	tests := []struct {
		name          string
		pattern       string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "match all go files",
			pattern:       filepath.Join(tmpDir, "*.go"),
			expectedCount: 5,
			expectError:   false,
		},
		{
			name:          "match test files only",
			pattern:       filepath.Join(tmpDir, "*_test.go"),
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "match markdown files",
			pattern:       filepath.Join(tmpDir, "*.md"),
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "no matches",
			pattern:       filepath.Join(tmpDir, "*.txt"),
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := CountFilesMatching(tt.pattern)
			if (err != nil) != tt.expectError {
				t.Errorf("CountFilesMatching() error = %v, expectError %v", err, tt.expectError)
			}
			if count != tt.expectedCount {
				t.Errorf("Expected count %d, got %d", tt.expectedCount, count)
			}
		})
	}
}

func TestCountTasks(t *testing.T) {
	tests := []struct {
		name     string
		plan     types.Plan
		expected int
	}{
		{
			name: "empty plan",
			plan: types.Plan{
				Chapters: []types.Chapter{},
			},
			expected: 0,
		},
		{
			name: "single task",
			plan: types.Plan{
				Chapters: []types.Chapter{
					{
						Quests: []types.Quest{
							{
								Tasks: []types.Task{
									{ID: "task1"},
								},
							},
						},
					},
				},
			},
			expected: 1,
		},
		{
			name: "multiple chapters and quests",
			plan: types.Plan{
				Chapters: []types.Chapter{
					{
						Quests: []types.Quest{
							{
								Tasks: []types.Task{
									{ID: "task1"},
									{ID: "task2"},
								},
							},
							{
								Tasks: []types.Task{
									{ID: "task3"},
								},
							},
						},
					},
					{
						Quests: []types.Quest{
							{
								Tasks: []types.Task{
									{ID: "task4"},
									{ID: "task5"},
									{ID: "task6"},
								},
							},
						},
					},
				},
			},
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := countTasks(tt.plan)
			if count != tt.expected {
				t.Errorf("Expected %d tasks, got %d", tt.expected, count)
			}
		})
	}
}

func TestSelectTemplateForSpec(t *testing.T) {
	tests := []struct {
		name     string
		spec     *types.ProjectSpec
		expected string
	}{
		{
			name:     "nil spec",
			spec:     nil,
			expected: "go-web-api",
		},
		{
			name: "explicit template",
			spec: &types.ProjectSpec{
				Template: "custom-template",
			},
			expected: "custom-template",
		},
		{
			name: "CLI build type",
			spec: &types.ProjectSpec{
				BuildType: types.BuildCLI,
			},
			expected: "go-cli-tool",
		},
		{
			name: "Automation build type",
			spec: &types.ProjectSpec{
				BuildType: types.BuildAutomation,
			},
			expected: "go-cli-tool",
		},
		{
			name: "Stream processor build type",
			spec: &types.ProjectSpec{
				BuildType: types.BuildStreamProcessor,
			},
			expected: "go-concurrency",
		},
		{
			name: "Worker build type",
			spec: &types.ProjectSpec{
				BuildType: types.BuildWorker,
			},
			expected: "go-concurrency",
		},
		{
			name: "Service with todo theme and normal difficulty",
			spec: &types.ProjectSpec{
				BuildType:  types.BuildService,
				Theme:      types.Todo,
				Difficulty: types.DifficultyNormal,
			},
			expected: "go-todo-api",
		},
		{
			name: "Service with deep difficulty",
			spec: &types.ProjectSpec{
				BuildType:  types.BuildService,
				Difficulty: types.DifficultyDeep,
			},
			expected: "go-auth-system",
		},
		{
			name: "Service default",
			spec: &types.ProjectSpec{
				BuildType:  types.BuildService,
				Difficulty: types.DifficultyQuick,
			},
			expected: "go-web-api",
		},
		{
			name: "Unknown build type with deep difficulty",
			spec: &types.ProjectSpec{
				BuildType:  "unknown",
				Difficulty: types.DifficultyDeep,
			},
			expected: "go-auth-system",
		},
		{
			name: "Unknown build type with normal difficulty",
			spec: &types.ProjectSpec{
				BuildType:  "unknown",
				Difficulty: types.DifficultyNormal,
			},
			expected: "go-web-api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := selectTemplateForSpec(tt.spec)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestIfCompletedPlan(t *testing.T) {
	tests := []struct {
		name     string
		plan     types.Plan
		state    types.State
		expected bool
	}{
		{
			name: "no tasks completed",
			plan: types.Plan{
				NumberOfTasks: 5,
			},
			state: types.State{
				CompletedTaskIDs: []string{},
			},
			expected: false,
		},
		{
			name: "some tasks completed",
			plan: types.Plan{
				NumberOfTasks: 5,
			},
			state: types.State{
				CompletedTaskIDs: []string{"task1", "task2"},
			},
			expected: false,
		},
		{
			name: "all tasks completed",
			plan: types.Plan{
				NumberOfTasks: 3,
			},
			state: types.State{
				CompletedTaskIDs: []string{"task1", "task2", "task3"},
			},
			expected: true,
		},
		{
			name: "zero tasks",
			plan: types.Plan{
				NumberOfTasks: 0,
			},
			state: types.State{
				CompletedTaskIDs: []string{},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ifCompletedPlan(tt.plan, tt.state)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCurrentChapter(t *testing.T) {
	plan := types.Plan{
		Chapters: []types.Chapter{
			{
				Title: "Chapter 1",
				Quests: []types.Quest{
					{
						Tasks: []types.Task{
							{ID: "task1"},
							{ID: "task2"},
						},
					},
				},
			},
			{
				Title: "Chapter 2",
				Quests: []types.Quest{
					{
						Tasks: []types.Task{
							{ID: "task3"},
							{ID: "task4"},
							{ID: "task5"},
						},
					},
				},
			},
			{
				Title: "Chapter 3",
				Quests: []types.Quest{
					{
						Tasks: []types.Task{
							{ID: "task6"},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name            string
		taskIndex       int
		expectedChapter int
	}{
		{
			name:            "first task",
			taskIndex:       0,
			expectedChapter: 0,
		},
		{
			name:            "second task in chapter 1",
			taskIndex:       1,
			expectedChapter: 0,
		},
		{
			name:            "first task in chapter 2",
			taskIndex:       2,
			expectedChapter: 1,
		},
		{
			name:            "middle task in chapter 2",
			taskIndex:       3,
			expectedChapter: 1,
		},
		{
			name:            "last task in chapter 2",
			taskIndex:       4,
			expectedChapter: 1,
		},
		{
			name:            "task in chapter 3",
			taskIndex:       5,
			expectedChapter: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := currentChapter(tt.taskIndex, plan)
			if result != tt.expectedChapter {
				t.Errorf("Expected chapter %d, got %d", tt.expectedChapter, result)
			}
		})
	}
}

func TestCurrentQuest(t *testing.T) {
	plan := types.Plan{
		Chapters: []types.Chapter{
			{
				Quests: []types.Quest{
					{
						Tasks: []types.Task{
							{ID: "task1"},
							{ID: "task2"},
						},
					},
					{
						Tasks: []types.Task{
							{ID: "task3"},
						},
					},
				},
			},
			{
				Quests: []types.Quest{
					{
						Tasks: []types.Task{
							{ID: "task4"},
							{ID: "task5"},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name          string
		taskIndex     int
		expectedQuest int
	}{
		{
			name:          "first task in first quest",
			taskIndex:     0,
			expectedQuest: 0,
		},
		{
			name:          "second task in first quest",
			taskIndex:     1,
			expectedQuest: 0,
		},
		{
			name:          "task in second quest",
			taskIndex:     2,
			expectedQuest: 1,
		},
		{
			name:          "first task in third quest",
			taskIndex:     3,
			expectedQuest: 2,
		},
		{
			name:          "second task in third quest",
			taskIndex:     4,
			expectedQuest: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := currentQuest(tt.taskIndex, plan)
			if result != tt.expectedQuest {
				t.Errorf("Expected quest %d, got %d", tt.expectedQuest, result)
			}
		})
	}
}

func TestTaskPerChapter(t *testing.T) {
	tests := []struct {
		name     string
		plan     types.Plan
		expected []int
	}{
		{
			name: "single chapter",
			plan: types.Plan{
				Chapters: []types.Chapter{
					{
						Quests: []types.Quest{
							{
								Tasks: []types.Task{
									{ID: "task1"},
									{ID: "task2"},
								},
							},
						},
					},
				},
			},
			expected: []int{2},
		},
		{
			name: "multiple chapters",
			plan: types.Plan{
				Chapters: []types.Chapter{
					{
						Quests: []types.Quest{
							{
								Tasks: []types.Task{
									{ID: "task1"},
									{ID: "task2"},
								},
							},
							{
								Tasks: []types.Task{
									{ID: "task3"},
								},
							},
						},
					},
					{
						Quests: []types.Quest{
							{
								Tasks: []types.Task{
									{ID: "task4"},
									{ID: "task5"},
									{ID: "task6"},
								},
							},
						},
					},
				},
			},
			expected: []int{3, 3},
		},
		{
			name: "empty chapters",
			plan: types.Plan{
				Chapters: []types.Chapter{
					{
						Quests: []types.Quest{},
					},
				},
			},
			expected: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := taskPerChapter(tt.plan)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("At index %d: expected %d, got %d", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestTaskPerQuest(t *testing.T) {
	tests := []struct {
		name     string
		plan     types.Plan
		expected []int
	}{
		{
			name: "single quest",
			plan: types.Plan{
				Chapters: []types.Chapter{
					{
						Quests: []types.Quest{
							{
								Tasks: []types.Task{
									{ID: "task1"},
									{ID: "task2"},
								},
							},
						},
					},
				},
			},
			expected: []int{2},
		},
		{
			name: "multiple quests in multiple chapters",
			plan: types.Plan{
				Chapters: []types.Chapter{
					{
						Quests: []types.Quest{
							{
								Tasks: []types.Task{
									{ID: "task1"},
								},
							},
							{
								Tasks: []types.Task{
									{ID: "task2"},
									{ID: "task3"},
								},
							},
						},
					},
					{
						Quests: []types.Quest{
							{
								Tasks: []types.Task{
									{ID: "task4"},
									{ID: "task5"},
									{ID: "task6"},
								},
							},
						},
					},
				},
			},
			expected: []int{1, 2, 3},
		},
		{
			name: "empty quests",
			plan: types.Plan{
				Chapters: []types.Chapter{
					{
						Quests: []types.Quest{
							{Tasks: []types.Task{}},
						},
					},
				},
			},
			expected: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := taskPerQuest(tt.plan)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("At index %d: expected %d, got %d", i, tt.expected[i], result[i])
				}
			}
		})
	}
}
