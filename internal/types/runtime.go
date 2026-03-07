package types

import "time"

// RunContext holds the runtime state for a single workflow execution.
type RunContext struct {
	WorkflowName string
	RunID        string // UUID
	StartedAt    time.Time
	TriggerEvent TriggerEvent
	Steps        map[string]StepResult // populated as steps complete
	Env          map[string]string     // merged: os env + secrets + workflow env
	Logger       Logger
}

// TriggerEvent describes the event that fired a workflow.
type TriggerEvent struct {
	Type      string
	FiredAt   time.Time
	Payload   map[string]interface{} // webhook body, etc.
	FilePath  string
	FileEvent string
}

// StepResult holds the outcome of a single step execution.
type StepResult struct {
	StepID     string
	Status     StepStatus
	StartedAt  time.Time
	EndedAt    time.Time
	Stdout     string
	Stderr     string
	ExitCode   int
	StatusCode int               // HTTP
	Body       string            // HTTP response / file read
	Headers    map[string]string
	Output     string // LLM / general text
	Error      string
}

// StepStatus is the execution state of a step.
type StepStatus string

const (
	StepStatusPending StepStatus = "pending"
	StepStatusRunning StepStatus = "running"
	StepStatusSuccess StepStatus = "success"
	StepStatusFailed  StepStatus = "failed"
	StepStatusSkipped StepStatus = "skipped"
)

// Logger is the logging interface used throughout the engine.
type Logger interface {
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}
