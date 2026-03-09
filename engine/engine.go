package engine

import (
	"context"
	"fmt"

	"github.com/vikingpingvin/chainrun/action"
	"github.com/vikingpingvin/chainrun/config"
	"github.com/vikingpingvin/chainrun/internal/types"
	"github.com/vikingpingvin/chainrun/secrets"
	"github.com/vikingpingvin/chainrun/template"
	"github.com/vikingpingvin/chainrun/trigger"
)

// Engine orchestrates workflow execution.
type Engine interface {
	RunOnce(ctx context.Context, workflowName string, event types.TriggerEvent) (*types.RunContext, error)
	DryRun(ctx context.Context, workflowName string, event types.TriggerEvent) error
	StartDaemon(ctx context.Context) error
	Workflows() []string
}

// Deps bundles all dependencies required to construct an engine.
type Deps struct {
	Workflows       []types.WorkflowDef
	ActionRegistry  action.Registry
	TriggerRegistry trigger.Registry
	Renderer        template.Renderer
	Secrets         secrets.Resolver
	Loader          config.Loader
	Logger          types.Logger
}

// engine is the concrete Engine implementation.
type engine struct {
	deps      Deps
	sequencer *sequencer
}

// New constructs a new Engine from the provided Deps.
func New(deps Deps) Engine {
	return &engine{
		deps:      deps,
		sequencer: newSequencer(deps),
	}
}

// RunOnce executes the named workflow once with the given trigger event.
func (e *engine) RunOnce(ctx context.Context, workflowName string, event types.TriggerEvent) (*types.RunContext, error) {
	// slices.Contains(e.deps.Workflows, workflowName)
	for _, wf := range e.deps.Workflows {
		if wf.Name == workflowName {
			runCtx := Build(wf, event, e.deps)
			err := e.sequencer.Run(ctx, wf, runCtx)
			return runCtx, err
		}
	}
	return nil, fmt.Errorf("workflow %q not found", workflowName)
}

// DryRun renders all template fields for the named workflow without executing any steps.
func (e *engine) DryRun(ctx context.Context, workflowName string, event types.TriggerEvent) error {
	for _, wf := range e.deps.Workflows {
		if wf.Name == workflowName {
			runCtx := Build(wf, event, e.deps)
			return e.sequencer.dryRun(ctx, wf, runCtx)
		}
	}
	return fmt.Errorf("workflow %q not found", workflowName)
}

// StartDaemon starts all workflow triggers and blocks until ctx is cancelled.
func (e *engine) StartDaemon(ctx context.Context) error {
	panic("not implemented")
}

// Workflows returns the names of all loaded workflows.
func (e *engine) Workflows() []string {
	wfNames := make([]string, len(e.deps.Workflows))
	for i := range e.deps.Workflows {
		wfNames[i] = e.deps.Workflows[i].Name
	}
	return wfNames
}
