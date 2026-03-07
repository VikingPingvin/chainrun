package shell

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"runtime"
	"time"

	"github.com/vikingpingvin/chainrun/internal/types"
)

// Executor runs shell commands.
type Executor struct{}

// New returns a new shell Executor.
func New() *Executor {
	return &Executor{}
}

// Execute runs the shell command defined in step.
// A non-zero exit code is reflected in result.Status and result.ExitCode but
// does not produce a Go error. A Go error is only returned for hard failures
// (binary not found, context cancelled before the process starts, etc.).
func (e *Executor) Execute(ctx context.Context, step types.StepDef, run *types.RunContext) (types.StepResult, error) {
	binary, flag := resolveShell(step.Shell)

	cmd := exec.CommandContext(ctx, binary, flag, step.Command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	result := types.StepResult{
		StepID:    step.ID,
		StartedAt: time.Now(),
	}

	runErr := cmd.Run()

	result.EndedAt = time.Now()
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			result.ExitCode = exitErr.ExitCode()
			result.Status = types.StepStatusFailed
			return result, nil
		}
		result.Status = types.StepStatusFailed
		result.Error = runErr.Error()
		return result, runErr
	}

	result.Status = types.StepStatusSuccess
	return result, nil
}

// resolveShell returns the shell binary and command flag for the given shell name.
// If name is empty, the OS default is used.
func resolveShell(name string) (binary, flag string) {
	switch name {
	case "bash":
		return "bash", "-c"
	case "sh":
		return "sh", "-c"
	case "powershell":
		return "powershell", "-Command"
	default:
		if runtime.GOOS == "windows" {
			return "cmd", "/C"
		}
		return "sh", "-c"
	}
}
