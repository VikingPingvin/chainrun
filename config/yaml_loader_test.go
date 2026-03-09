package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vikingpingvin/chainrun/internal/types"
)

// writeTemp writes content to a temp file and returns its path.
func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), "chainrun.yaml")
	if err := os.WriteFile(f, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return f
}

// ── Load ─────────────────────────────────────────────────────────────────────

func TestLoad_ValidFull(t *testing.T) {
	yaml := `
workflows:
  - name: hello
    description: greet the world
    trigger:
      type: manual
    env:
      FOO: bar
    steps:
      - id: greet
        timeout: 10s
        continue_on_error: true
        retry:
          attempts: 3
          delay: 2s
          backoff: exponential
        shell:
          command: echo hello
`
	cfg, err := NewYAMLLoader().Load(writeTemp(t, yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Workflows) != 1 {
		t.Fatalf("want 1 workflow, got %d", len(cfg.Workflows))
	}
	wf := cfg.Workflows[0]
	if wf.Name != "hello" {
		t.Errorf("name: want %q, got %q", "hello", wf.Name)
	}
	if wf.Env["FOO"] != "bar" {
		t.Errorf("env FOO: want %q, got %q", "bar", wf.Env["FOO"])
	}
	step := wf.Steps[0]
	if step.ID != "greet" {
		t.Errorf("step id: want %q, got %q", "greet", step.ID)
	}
	if !step.ContinueOnError {
		t.Error("continue_on_error: want true")
	}
	if step.Retry.Attempts != 3 {
		t.Errorf("retry.attempts: want 3, got %d", step.Retry.Attempts)
	}
}

func TestLoad_MultipleWorkflows(t *testing.T) {
	yaml := `
workflows:
  - name: alpha
    trigger:
      type: manual
    steps:
      - id: s1
        shell:
          command: echo alpha
  - name: beta
    trigger:
      type: cron
      schedule: "0 * * * *"
    steps:
      - id: s1
        shell:
          command: echo beta
`
	cfg, err := NewYAMLLoader().Load(writeTemp(t, yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Workflows) != 2 {
		t.Fatalf("want 2 workflows, got %d", len(cfg.Workflows))
	}
	if cfg.Workflows[1].Trigger.Schedule != "0 * * * *" {
		t.Errorf("schedule: want %q, got %q", "0 * * * *", cfg.Workflows[1].Trigger.Schedule)
	}
}

func TestLoad_EmptyWorkflows(t *testing.T) {
	cfg, err := NewYAMLLoader().Load(writeTemp(t, "workflows: []\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Workflows) != 0 {
		t.Errorf("want 0 workflows, got %d", len(cfg.Workflows))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := NewYAMLLoader().Load("/no/such/file.yaml")
	if err == nil {
		t.Fatal("want error for missing file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	_, err := NewYAMLLoader().Load(writeTemp(t, "key: [unclosed"))
	if err == nil {
		t.Fatal("want error for invalid YAML, got nil")
	}
}

func TestLoad_WrongShape(t *testing.T) {
	// Valid YAML but unknown keys — yaml.v3 ignores them, producing empty Config.
	cfg, err := NewYAMLLoader().Load(writeTemp(t, "something: else\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Workflows) != 0 {
		t.Errorf("want 0 workflows, got %d", len(cfg.Workflows))
	}
}

// ── Validate ──────────────────────────────────────────────────────────────────

func TestValidate_Valid(t *testing.T) {
	cfg := &types.Config{
		Workflows: []types.WorkflowDef{
			{
				Name:    "hello",
				Trigger: types.TriggerDef{Type: "manual"},
				Steps: []types.StepDef{
					{ID: "s1", Shell: &types.ShellStep{Command: "echo hi"}},
					{ID: "s2", HTTP: &types.HTTPStep{URL: "http://example.com"}},
				},
			},
		},
	}
	errs := NewYAMLLoader().Validate(cfg)
	if len(errs) != 0 {
		t.Errorf("want no errors, got %v", errs)
	}
}

func TestValidate_MissingWorkflowName(t *testing.T) {
	cfg := &types.Config{
		Workflows: []types.WorkflowDef{
			{Trigger: types.TriggerDef{Type: "manual"}, Steps: []types.StepDef{{ID: "s1", Shell: &types.ShellStep{Command: "echo hi"}}}},
		},
	}
	assertFieldError(t, NewYAMLLoader().Validate(cfg), "workflows[0].name")
}

func TestValidate_DuplicateWorkflowName(t *testing.T) {
	wf := types.WorkflowDef{
		Name:    "dup",
		Trigger: types.TriggerDef{Type: "manual"},
		Steps:   []types.StepDef{{ID: "s1", Shell: &types.ShellStep{Command: "echo hi"}}},
	}
	cfg := &types.Config{Workflows: []types.WorkflowDef{wf, wf}}
	assertFieldError(t, NewYAMLLoader().Validate(cfg), "workflows[1].name")
}

func TestValidate_MissingTriggerType(t *testing.T) {
	cfg := &types.Config{
		Workflows: []types.WorkflowDef{
			{Name: "hello", Steps: []types.StepDef{{ID: "s1", Shell: &types.ShellStep{Command: "echo hi"}}}},
		},
	}
	assertFieldError(t, NewYAMLLoader().Validate(cfg), "workflows[0].trigger.type")
}

func TestValidate_MissingStepID(t *testing.T) {
	cfg := &types.Config{
		Workflows: []types.WorkflowDef{
			{
				Name:    "hello",
				Trigger: types.TriggerDef{Type: "manual"},
				Steps:   []types.StepDef{{Shell: &types.ShellStep{Command: "echo hi"}}},
			},
		},
	}
	assertFieldError(t, NewYAMLLoader().Validate(cfg), "workflows[0].steps[0].id")
}

func TestValidate_MissingStepType(t *testing.T) {
	cfg := &types.Config{
		Workflows: []types.WorkflowDef{
			{
				Name:    "hello",
				Trigger: types.TriggerDef{Type: "manual"},
				Steps:   []types.StepDef{{ID: "s1"}},
			},
		},
	}
	assertFieldError(t, NewYAMLLoader().Validate(cfg), "workflows[0].steps[0]")
}

func TestValidate_DuplicateStepID(t *testing.T) {
	cfg := &types.Config{
		Workflows: []types.WorkflowDef{
			{
				Name:    "hello",
				Trigger: types.TriggerDef{Type: "manual"},
				Steps: []types.StepDef{
					{ID: "dup", Shell: &types.ShellStep{Command: "echo hi"}},
					{ID: "dup", HTTP: &types.HTTPStep{URL: "http://example.com"}},
				},
			},
		},
	}
	assertFieldError(t, NewYAMLLoader().Validate(cfg), "workflows[0].steps[1].id")
}

func TestValidate_MultipleErrors(t *testing.T) {
	cfg := &types.Config{
		Workflows: []types.WorkflowDef{
			// missing name, missing trigger type, step missing id and executor block
			{Steps: []types.StepDef{{}}},
		},
	}
	errs := NewYAMLLoader().Validate(cfg)
	if len(errs) < 3 {
		t.Errorf("want at least 3 errors, got %d: %v", len(errs), errs)
	}
}

func TestValidate_EmptyConfig(t *testing.T) {
	errs := NewYAMLLoader().Validate(&types.Config{})
	if len(errs) != 0 {
		t.Errorf("want no errors for empty config, got %v", errs)
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

// assertFieldError fails the test if no ValidationError targets the given field.
func assertFieldError(t *testing.T, errs []ValidationError, field string) {
	t.Helper()
	for _, e := range errs {
		if e.Field == field {
			return
		}
	}
	t.Errorf("want ValidationError for field %q, got %v", field, errs)
}
