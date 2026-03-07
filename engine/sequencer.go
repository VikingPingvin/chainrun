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

	executor, ok := s.deps.ActionRegistry.Get(rendered.Type)
	if !ok {
		return fmt.Errorf("step %q: unknown action type %q", step.ID, rendered.Type)
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

// renderStep returns a copy of step with all template fields resolved against runCtx.
func (s *sequencer) renderStep(step types.StepDef, runCtx *types.RunContext) (types.StepDef, error) {
	r := s.deps.Renderer

	fields := []*string{
		&step.Command, &step.URL, &step.Body,
		&step.Prompt, &step.System,
		&step.Message, &step.FilePath, &step.FileContent,
	}
	for _, f := range fields {
		rendered, err := r.Render(*f, runCtx)
		if err != nil {
			return step, err
		}
		*f = rendered
	}

	if len(step.Headers) > 0 {
		rendered, err := r.RenderMap(step.Headers, runCtx)
		if err != nil {
			return step, err
		}
		step.Headers = rendered
	}
	if len(step.Env) > 0 {
		rendered, err := r.RenderMap(step.Env, runCtx)
		if err != nil {
			return step, err
		}
		step.Env = rendered
	}

	return step, nil
}
