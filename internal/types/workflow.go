package types

// WorkflowDef is the top-level workflow definition loaded from a config file.
type WorkflowDef struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Trigger     TriggerDef        `yaml:"trigger"`
	Steps       []StepDef         `yaml:"steps"`
	OnError     ErrorPolicy       `yaml:"on_error"`
	Env         map[string]string `yaml:"env"`
}

// TriggerDef describes how a workflow is triggered.
type TriggerDef struct {
	Type     string   `yaml:"type"`     // "cron" | "manual" | "webhook" | "watch" | "git"
	Schedule string   `yaml:"schedule"` // cron only
	Path     string   `yaml:"path"`     // watch only
	Events   []string `yaml:"events"`   // watch: ["create","modify","delete"]
	Port     int      `yaml:"port"`     // webhook only
	Route    string   `yaml:"route"`    // webhook only
	Hook     string   `yaml:"hook"`     // git only
}

// StepDef describes a single step in a workflow.
type StepDef struct {
	ID              string            `yaml:"id"`
	Type            string            `yaml:"type"` // "shell" | "http" | "llm" | "file" | "notify"
	ContinueOnError bool              `yaml:"continue_on_error"`
	Retry           RetryPolicy       `yaml:"retry"`
	Timeout         string            `yaml:"timeout"` // e.g. "30s", "5m"
	Env             map[string]string `yaml:"env"`

	// Shell
	Command string `yaml:"command"`
	Shell   string `yaml:"shell"` // "sh" | "bash" | "powershell"

	// HTTP
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`

	// LLM
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
	Prompt   string `yaml:"prompt"`
	System   string `yaml:"system"`

	// File
	FileOp      string `yaml:"file_op"` // "read" | "write" | "append"
	FilePath    string `yaml:"file_path"`
	FileContent string `yaml:"file_content"`

	// Notify
	Channel string `yaml:"channel"`
	Message string `yaml:"message"`
	To      string `yaml:"to"`
}

// RetryPolicy controls retry behaviour for a step.
type RetryPolicy struct {
	Attempts int    `yaml:"attempts"`
	Delay    string `yaml:"delay"`   // e.g. "5s"
	Backoff  string `yaml:"backoff"` // "fixed" | "exponential"
}

// ErrorPolicy controls what happens when a workflow-level error occurs.
type ErrorPolicy struct {
	Strategy string `yaml:"strategy"`  // "stop" (default) | "continue" | "notify"
	NotifyID string `yaml:"notify_id"` // step ID to invoke on error
}

// Config is the root configuration structure loaded from a config file.
type Config struct {
	Workflows []WorkflowDef `yaml:"workflows"`
}
