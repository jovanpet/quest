package types

import "time"

type Plan struct {
	Version       int       `json:"version"`
	Journey       Journey   `json:"journey"`
	Chapters      []Chapter `json:"chapters"`
	NumberOfTasks int       `json:"numberOfTasks"`
}

type Journey struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Language    string   `json:"language"`
	Focus       []string `json:"focus"`
}

type Chapter struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Quests []Quest `json:"quests"`
}

type Quest struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Objective  string     `json:"objective"`
	Steps      []string   `json:"steps,omitempty"`
	Files      []string   `json:"files,omitempty"`
	Artifacts  []string   `json:"artifacts"`
	Validation Validation `json:"validation"`
}

type Validation struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Type        Type   `json:"type"` // "exists", "glob_count_min", "file_contains_any"
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	// For Type == "exists"
	Path string `json:"path,omitempty"`

	// For Type == "glob_count_min" and "file_contains_any"
	Glob string `json:"glob,omitempty"`

	// For Type == "glob_count_min"
	Min int `json:"min,omitempty"`

	// For Type == "file_contains_any"
	Any []string `json:"any,omitempty"`

	LastState *CheckState `json:"lastState,omitempty"`
}

type State struct {
	Version int `json:"version"`

	// Index into the flattened task list
	CurrentTaskIndex int `json:"currentTaskIndex"`

	// Task IDs that have been completed
	CompletedTaskIDs []string `json:"completedTaskIds"`

	// Result of the most recent check (nil if never run)
	LastCheck *CheckResult `json:"lastCheck"`

	QuestStarted bool `json:"questStarted"`

	// Track how many times explain was called for current task
	ExplainCount int `json:"explainCount"`
}

type CheckResult struct {
	TaskID    int         `json:"taskId"`
	Status    CheckStatus `json:"status"` // "pass", "warn", "fail"
	Timestamp time.Time   `json:"timestamp"`
	Message   string      `json:"message,omitempty"`
}

type CheckStatus string

type Status string

const (
	CheckPass CheckStatus = "pass"
	CheckWarn CheckStatus = "warn"
	CheckFail CheckStatus = "fail"
)

type Type string

const (
	TypeExists          Type = "exists"
	TypeGlobCountMin    Type = "glob_count_min"
	TypeFileContainsAny Type = "file_contains_any"
)

type CheckState string

const (
	Pass CheckState = "pass"
	Fail CheckState = "fail"
)

type BuildType string

// WIZARD SETUP

type Difficulty string

const (
	DifficultyQuick  Difficulty = "quick"  // 3-10 tasks
	DifficultyNormal Difficulty = "normal" // 10-20 tasks
	DifficultyDeep   Difficulty = "deep"   // 20-35 tasks
)

const (
	// Core application shapes
	BuildService  BuildType = "service"  // HTTP / backend service
	BuildCLI      BuildType = "cli"      // Command-line tool
	BuildPipeline BuildType = "pipeline" // ETL / batch processing
	BuildWorker   BuildType = "worker"   // Background jobs / consumers
	BuildDaemon   BuildType = "daemon"   // Long-running process
	BuildLibrary  BuildType = "library"  // Reusable Go package

	// Systems & infra
	BuildScheduler    BuildType = "scheduler"     // Cron-like task scheduler
	BuildController   BuildType = "controller"    // Reconciler / control loop
	BuildAgent        BuildType = "agent"         // Host-level agent
	BuildProxy        BuildType = "proxy"         // Reverse proxy / gateway
	BuildLoadBalancer BuildType = "load-balancer" // Traffic routing / balancing

	// Data & streaming
	BuildStreamProcessor BuildType = "stream-processor" // Event / stream processing
	BuildIndexer         BuildType = "indexer"          // Build & query indexes
	BuildAggregator      BuildType = "aggregator"       // Rollups / summaries

	// Tooling & automation
	BuildAutomation BuildType = "automation" // Ops / scripts / tooling
	BuildScaffolder BuildType = "scaffolder" // Code/project generator
	BuildLinter     BuildType = "linter"     // Static analysis tool
	BuildFormatter  BuildType = "formatter"  // Code formatting tool

	// Networking / protocols
	BuildProtocol  BuildType = "protocol"   // Custom protocol implementation
	BuildClientSDK BuildType = "client-sdk" // API client / SDK
)

type Theme string

const (
	Todo      Theme = "todo"
	Incidents Theme = "incidents"
	Notes     Theme = "notes"
	Bookmarks Theme = "bookmarks"
	Expenses  Theme = "expenses"
	Contacts  Theme = "contacts"
	Events    Theme = "events"
	Inventory Theme = "inventory"
)

type ProjectSpec struct {
	BuildType   BuildType  `json:"buildType"`
	Difficulty  Difficulty `json:"difficulty"`
	Theme       Theme      `json:"theme"`       // optional
	Description string     `json:"description"` // optional: free-text
	Template    string     `json:"template"`    // chosen template name
	SurpriseMe  bool       `json:"surpriseMe"`  // optional: randomize project spec
}
