package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/vikingpingvin/chainrun/internal/types"
)

// sequencer executes the steps of a workflow in order, applying retry and
// error policy as configured.
type sequencer struct {
	deps Deps
}

func newSequencer(deps Deps) *sequencer {
	return &sequencer{deps: deps}
}

// Run executes all steps in def sequentially against runCtx.
// If a step fails and ContinueOnError is set, the failure is logged and
// execution proceeds. Otherwise the first failure aborts the workflow.
func (s *sequencer) Run(ctx context.Context, def types.WorkflowDef, runCtx *types.RunContext) error {
	for _, step := range def.Steps {
		if err := s.runStep(ctx, step, runCtx); err != nil {
			if step.ContinueOnError {
				runCtx.Logger.Warn("step %q failed, continuing: %v", step.ID, err)
				continue
			}
			return err
		}
	}
	return nil
}

// runStep executes a single step: applies the timeout, renders templates,
// dispatches to the executor, and commits the result to runCtx.
func (s *sequencer) runStep(ctx context.Context, step types.StepDef, runCtx *types.RunContext) error {
	stepCtx, cancel := s.stepContext(ctx, step.Timeout)
	defer cancel()

	rendered, err := s.renderStep(step, runCtx)
	if err != nil {
		return fmt.Errorf("step %q: render failed: %w", step.ID, err)
	}

	executor, ok := s.deps.ActionRegistry.Get(rendered.ActionType())
	if !ok {
		return fmt.Errorf("step %q: unknown action type %q", step.ID, rendered.ActionType())
	}

	runCtx.Steps[step.ID] = types.StepResult{
		StepID:    step.ID,
		Status:    types.StepStatusRunning,
		StartedAt: time.Now(),
	}

	result, err := executor.Execute(stepCtx, rendered, runCtx)
	result.StepID = step.ID
	result.EndedAt = time.Now()
	runCtx.Steps[step.ID] = result

	if err != nil {
		return fmt.Errorf("step %q: %w", step.ID, err)
	}
	if result.Status == types.StepStatusFailed {
		return fmt.Errorf("step %q failed: %s", step.ID, result.Error)
	}
	return nil
}

// stepContext returns a context with a timeout if one is configured, or the
// parent context unchanged. The caller must always call the returned cancel.
func (s *sequencer) stepContext(ctx context.Context, timeout string) (context.Context, context.CancelFunc) {
	if timeout == "" {
		return ctx, func() {}
	}
	d, err := time.ParseDuration(timeout)
	if err != nil {
		// Invalid timeout string — fall back to parent context rather than
		// failing hard here; the step inherits the workflow-level deadline.
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, d)
}

// dryRun renders all template fields for each step without executing any of them.
func (s *sequencer) dryRun(_ context.Context, def types.WorkflowDef, runCtx *types.RunContext) error {
	for _, step := range def.Steps {
		if _, err := s.renderStep(step, runCtx); err != nil {
			return fmt.Errorf("step %q: render failed: %w", step.ID, err)
		}
	}
	return nil
}

// renderStep returns a copy of step with all template fields resolved against runCtx.
func (s *sequencer) renderStep(step types.StepDef, runCtx *types.RunContext) (types.StepDef, error) {
	rendered := step

	if step.Shell != nil {
		shell := *step.Shell
		var err error
		if shell.Command, err = s.deps.Renderer.Render(shell.Command, runCtx); err != nil {
			return rendered, err
		}
		rendered.Shell = &shell
	}

	if step.HTTP != nil {
		h := *step.HTTP
		var err error
		if h.URL, err = s.deps.Renderer.Render(h.URL, runCtx); err != nil {
			return rendered, err
		}
		if h.Body, err = s.deps.Renderer.Render(h.Body, runCtx); err != nil {
			return rendered, err
		}
		renderedHeaders := make(map[string]string, len(h.Headers))
		for k, v := range h.Headers {
			rv, err := s.deps.Renderer.Render(v, runCtx)
			if err != nil {
				return rendered, err
			}
			renderedHeaders[k] = rv
		}
		h.Headers = renderedHeaders
		rendered.HTTP = &h
	}

	if step.LLM != nil {
		l := *step.LLM
		var err error
		if l.Prompt, err = s.deps.Renderer.Render(l.Prompt, runCtx); err != nil {
			return rendered, err
		}
		if l.System, err = s.deps.Renderer.Render(l.System, runCtx); err != nil {
			return rendered, err
		}
		rendered.LLM = &l
	}

	if step.File != nil {
		f := *step.File
		var err error
		if f.Path, err = s.deps.Renderer.Render(f.Path, runCtx); err != nil {
			return rendered, err
		}
		if f.Content, err = s.deps.Renderer.Render(f.Content, runCtx); err != nil {
			return rendered, err
		}
		rendered.File = &f
	}

	if step.Notify != nil {
		n := *step.Notify
		var err error
		if n.Message, err = s.deps.Renderer.Render(n.Message, runCtx); err != nil {
			return rendered, err
		}
		rendered.Notify = &n
	}

	// render Env map
	renderedEnv := make(map[string]string, len(step.Env))
	for k, v := range step.Env {
		rv, err := s.deps.Renderer.Render(v, runCtx)
		if err != nil {
			return rendered, err
		}
		renderedEnv[k] = rv
	}
	rendered.Env = renderedEnv

	return rendered, nil
}
