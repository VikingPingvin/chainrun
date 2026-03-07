package shell

import (
	"context"

	"github.com/vikingpingvin/chainrun/internal/types"
)

// Executor runs shell commands.
type Executor struct{}

// New returns a new shell Executor.
func New() *Executor {
	return &Executor{}
}

// Execute runs the shell command defined in step.
func (e *Executor) Execute(ctx context.Context, step types.StepDef, run *types.RunContext) (types.StepResult, error) {
	panic("not implemented")
}
