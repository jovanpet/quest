package template

import (
	"testing"
)

func TestLoadTemplates(t *testing.T) {
	templates := []string{"go-web-api", "go-cli-tool", "go-concurrency", "go-todo-api"}
	
	for _, name := range templates {
		t.Run(name, func(t *testing.T) {
			plan, err := Load(name)
			if err != nil {
				t.Fatalf("failed to load template %s: %v", name, err)
			}
			
			// Validate basic structure
			if plan.Version != 1 {
				t.Errorf("expected version 1, got %d", plan.Version)
			}
			
			if plan.Journey.Name == "" {
				t.Error("journey name is empty")
			}
			
			if plan.Journey.Language != "go" {
				t.Errorf("expected language 'go', got %s", plan.Journey.Language)
			}
			
			if len(plan.Chapters) == 0 {
				t.Error("no chapters found")
			}
			
			// Check structure
			taskCount := 0
			for _, chapter := range plan.Chapters {
				if chapter.ID == "" || chapter.Title == "" {
					t.Error("chapter missing id or title")
				}
				
				for _, quest := range chapter.Quests {
					if quest.ID == "" || quest.Title == "" {
						t.Error("quest missing id or title")
					}
					
					taskCount += len(quest.Tasks)
					
					for _, task := range quest.Tasks {
						if task.ID == "" || task.Title == "" || task.Objective == "" {
							t.Error("task missing required fields")
						}
						
						if len(task.Artifacts) == 0 {
							t.Errorf("task %s has no artifacts", task.ID)
						}
						
						if len(task.Validation.Rules) == 0 {
							t.Errorf("task %s has no validation rules", task.ID)
						}
					}
				}
			}
			
			if taskCount == 0 {
				t.Error("template has no tasks")
			}
		})
	}
}

func TestLoadInvalidTemplate(t *testing.T) {
	_, err := Load("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent template")
	}
}

func TestListTemplates(t *testing.T) {
	templates := List()
	
	if len(templates) < 3 {
		t.Errorf("expected at least 3 templates, got %d", len(templates))
	}
	
	for _, tmpl := range templates {
		if tmpl.Name == "" || tmpl.Title == "" || tmpl.Description == "" {
			t.Errorf("template info incomplete: %+v", tmpl)
		}
		
		if tmpl.Tasks <= 0 {
			t.Errorf("template %s has invalid task count: %d", tmpl.Name, tmpl.Tasks)
		}
	}
}
