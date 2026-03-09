package shell

import (
	"context"
	"strings"
	"testing"

	"github.com/vikingpingvin/chainrun/internal/types"
)

func TestExecute_Success(t *testing.T) {
	e := New()
	step := types.StepDef{ID: "greet", Shell: &types.ShellStep{Command: "echo hello"}}

	result, err := e.Execute(context.Background(), step, &types.RunContext{})
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if result.Status != types.StepStatusSuccess {
		t.Errorf("status: want %q, got %q", types.StepStatusSuccess, result.Status)
	}
	if !strings.Contains(result.Stdout, "hello") {
		t.Errorf("stdout: want %q to contain %q", result.Stdout, "hello")
	}
	if result.StepID != "greet" {
		t.Errorf("StepID: want %q, got %q", "greet", result.StepID)
	}
	if result.StartedAt.IsZero() {
		t.Error("StartedAt should not be zero")
	}
	if result.EndedAt.IsZero() {
		t.Error("EndedAt should not be zero")
	}
}

func TestExecute_NonZeroExit(t *testing.T) {
	e := New()
	step := types.StepDef{ID: "fail", Shell: &types.ShellStep{Command: "exit 1"}}

	result, err := e.Execute(context.Background(), step, &types.RunContext{})
	if err != nil {
		t.Fatalf("non-zero exit must not return a Go error, got: %v", err)
	}
	if result.Status != types.StepStatusFailed {
		t.Errorf("status: want %q, got %q", types.StepStatusFailed, result.Status)
	}
	if result.ExitCode != 1 {
		t.Errorf("ExitCode: want 1, got %d", result.ExitCode)
	}
}
