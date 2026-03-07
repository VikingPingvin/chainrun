package engine

import (
	"context"

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

// Run executes all steps in def against the provided RunContext.
// Steps are run sequentially; retry and error policies are honoured.
func (s *sequencer) Run(ctx context.Context, def types.WorkflowDef, runCtx *types.RunContext) error {
	panic("not implemented")
}
