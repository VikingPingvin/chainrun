package types

// WorkflowDef is the top-level workflow definition loaded from a config file.
type WorkflowDef struct {
	Name        string
	Description string
	Trigger     TriggerDef
	Steps       []StepDef
	OnError     ErrorPolicy
	Env         map[string]string
}

// TriggerDef describes how a workflow is triggered.
type TriggerDef struct {
	Type     string   // "cron" | "manual" | "webhook" | "watch" | "git"
	Schedule string   // cron only
	Path     string   // watch only
	Events   []string // watch: ["create","modify","delete"]
	Port     int      // webhook only
	Route    string   // webhook only
	Hook     string   // git only
}

// StepDef describes a single step in a workflow.
type StepDef struct {
	ID              string
	Type            string // "shell" | "http" | "llm" | "file" | "notify"
	ContinueOnError bool
	Retry           RetryPolicy
	Timeout         string            // e.g. "30s", "5m"
	Env             map[string]string

	// Shell
	Command string
	Shell   string // "sh" | "bash" | "powershell"

	// HTTP
	URL     string
	Method  string
	Headers map[string]string
	Body    string

	// LLM
	Provider string
	Model    string
	Prompt   string
	System   string

	// File
	FileOp      string // "read" | "write" | "append"
	FilePath    string
	FileContent string

	// Notify
	Channel string
	Message string
	To      string
}

// RetryPolicy controls retry behaviour for a step.
type RetryPolicy struct {
	Attempts int
	Delay    string // e.g. "5s"
	Backoff  string // "fixed" | "exponential"
}

// ErrorPolicy controls what happens when a workflow-level error occurs.
type ErrorPolicy struct {
	Strategy string // "stop" (default) | "continue" | "notify"
	NotifyID string // step ID to invoke on error
}

// Config is the root configuration structure loaded from a config file.
type Config struct {
	Workflows []WorkflowDef
}
