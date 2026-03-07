package engine

import (
	"github.com/vikingpingvin/chainrun/internal/types"
)

// Build constructs a RunContext for a workflow execution.
// Env is merged in priority order: os environment < secrets < workflow.Env.
func Build(def types.WorkflowDef, event types.TriggerEvent, deps Deps) *types.RunContext {
	panic("not implemented")
}
