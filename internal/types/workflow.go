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
	Type     TriggerType `yaml:"type"`     // "cron" | "manual" | "webhook" | "watch" | "git"
	Schedule string      `yaml:"schedule"` // cron only
	Path     string      `yaml:"path"`     // watch only
	Events   []string    `yaml:"events"`   // watch: ["create","modify","delete"]
	Port     int         `yaml:"port"`     // webhook only
	Route    string      `yaml:"route"`    // webhook only
	Hook     string      `yaml:"hook"`     // git only
}

type TriggerType string

const (
	Cron    TriggerType = "cron"
	Manual  TriggerType = "manual"
	Webhook TriggerType = "webhook"
	Watch   TriggerType = "watch"
	Git     TriggerType = "git"
)

// ShellStep holds configuration for a shell executor step.
type ShellStep struct {
	Command string `yaml:"command"`
	Shell   string `yaml:"shell"`
}

// HTTPStep holds configuration for an HTTP executor step.
type HTTPStep struct {
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

// LLMStep holds configuration for an LLM executor step.
type LLMStep struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
	Prompt   string `yaml:"prompt"`
	System   string `yaml:"system"`
}

// FileStep holds configuration for a file executor step.
type FileStep struct {
	Op      FileOperationType `yaml:"op"`
	Path    string            `yaml:"path"`
	Content string            `yaml:"content"`
}

// NotifyStep holds configuration for a notify executor step.
type NotifyStep struct {
	Channel string `yaml:"channel"`
	Message string `yaml:"message"`
	To      string `yaml:"to"`
}

// StepDef describes a single step in a workflow.
type StepDef struct {
	// Global
	ID              string            `yaml:"id"`
	ContinueOnError bool              `yaml:"continue_on_error"`
	Retry           RetryPolicy       `yaml:"retry"`
	Timeout         string            `yaml:"timeout"` // e.g. "30s", "5m"
	Env             map[string]string `yaml:"env"`

	// Executor blocks (exactly one should be set)
	Shell  *ShellStep  `yaml:"shell"`
	HTTP   *HTTPStep   `yaml:"http"`
	LLM    *LLMStep    `yaml:"llm"`
	File   *FileStep   `yaml:"file"`
	Notify *NotifyStep `yaml:"notify"`
}

// ActionType returns the name of the executor block that is set on this step.
func (s StepDef) ActionType() string {
	switch {
	case s.Shell != nil:
		return "shell"
	case s.HTTP != nil:
		return "http"
	case s.LLM != nil:
		return "llm"
	case s.File != nil:
		return "file"
	case s.Notify != nil:
		return "notify"
	default:
		return ""
	}
}

type FileOperationType string

const (
	FileRead   FileOperationType = "read"
	FileWrite  FileOperationType = "write"
	FileAppend FileOperationType = "append"
)

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
