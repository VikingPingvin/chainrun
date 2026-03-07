package template

import (
	"testing"

	"github.com/vikingpingvin/chainrun/internal/types"
)

func newCtx(env map[string]string, steps map[string]types.StepResult) *types.RunContext {
	return &types.RunContext{Env: env, Steps: steps}
}

func TestRender_PlainString(t *testing.T) {
	r := NewRenderer()
	ctx := newCtx(nil, nil)

	got, err := r.Render("echo hello", ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "echo hello" {
		t.Errorf("want %q, got %q", "echo hello", got)
	}
}

func TestRender_EnvSubstitution(t *testing.T) {
	r := NewRenderer()
	ctx := newCtx(map[string]string{"USER": "viking"}, nil)

	got, err := r.Render("echo {{.Env.USER}}", ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "echo viking"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}
