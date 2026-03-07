package engine

import (
	"fmt"
	"time"

	"github.com/vikingpingvin/chainrun/internal/types"
)

// Build constructs a RunContext for a workflow execution.
// Env is merged in priority order: secrets < workflow.Env.
// os.Environ() is not included — executors inherit the process environment automatically.
func Build(def types.WorkflowDef, event types.TriggerEvent, deps Deps) *types.RunContext {
	return &types.RunContext{
		WorkflowName: def.Name,
		RunID:        resolveRunID(),
		Env:          mergeEnv(def, deps),
		TriggerEvent: event,
		Steps:        make(map[string]types.StepResult),
		Logger:       deps.Logger,
		StartedAt:    time.Now(),
	}
}

// mergeEnv builds the environment map for a run.
// Secrets form the base layer; workflow-level env overrides on top.
func mergeEnv(def types.WorkflowDef, deps Deps) map[string]string {
	secrets := deps.Secrets.ResolveAll()
	env := make(map[string]string, len(secrets)+len(def.Env))
	for k, v := range secrets {
		env[k] = v
	}
	for k, v := range def.Env {
		env[k] = v
	}
	return env
}

func resolveRunID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
