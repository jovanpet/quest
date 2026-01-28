package template

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/jovanpet/quest/internal/types"
)

//go:embed *.json
var templateFiles embed.FS

func Load(name string) (*types.Plan, error) {
	filename := fmt.Sprintf("%s.json", name)

	data, err := templateFiles.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("template not found: %s", name)
	}

	var plan types.Plan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &plan, nil
}

// List returns all available template names with descriptions
func List() []TemplateInfo {
	return []TemplateInfo{
		{
			Name:        "go-web-api",
			Title:       "REST API with Go",
			Description: "Learn to build HTTP servers, handle JSON, and write tests",
			Tasks:       3,
			Tier:        Quick,
		},
		{
			Name:        "go-cli-tool",
			Title:       "CLI Tool with Cobra",
			Description: "Create command-line tools with subcommands and flags",
			Tasks:       3,
			Tier:        Quick,
		},
		{
			Name:        "go-concurrency",
			Title:       "Go Concurrency Patterns",
			Description: "Master goroutines, channels, and worker pools",
			Tasks:       3,
			Tier:        Quick,
		},
		{
			Name:        "go-auth-system",
			Title:       "User Authentication System",
			Description: "Build complete auth with database, sessions, JWT, and security",
			Tasks:       19,
			Tier:        Advanced,
		},
		{
			Name:        "go-todo-api",
			Title:       "Complete Todo REST API",
			Description: "Build a full-featured REST API with CRUD, testing, and middleware",
			Tasks:       15,
			Tier:        Advanced,
		},
		{
			Name:        "go-isekai-server",
			Title:       "Distributed World Manager",
			Description: "Manage a distributed virtual world with Go",
			Tasks:       20,
			Tier:        Advanced,
		},
		{
			Name:        "go-fairy-garden",
			Title:       "Fairy Worker Service",
			Description: "Build a whimsical worker service with Go",
			Tasks:       10,
			Tier:        Normal,
		},
	}
}

// ListTemplates returns just the template names as strings
func ListTemplates() []string {
	templates := List()
	names := make([]string, len(templates))
	for i, t := range templates {
		names[i] = fmt.Sprintf("%s - %s", t.Name, t.Title)
	}
	return names
}

type TemplateInfo struct {
	Name        string
	Title       string
	Description string
	Tasks       int
	Tier        Tier
}

type Tier string

const (
	Quick    Tier = "Quick"
	Normal   Tier = "Normal"
	Advanced Tier = "Advanced"
)
