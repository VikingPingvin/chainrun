package file

import (
	"context"

	"github.com/vikingpingvin/chainrun/internal/types"
	"github.com/vikingpingvin/chainrun/template"
)

// Executor performs file operations (read / write / append).
type Executor struct {
	renderer template.Renderer
}

// New returns a new file Executor.
func New(renderer template.Renderer) *Executor {
	return &Executor{renderer: renderer}
}

// Execute performs the file operation defined in step.
func (e *Executor) Execute(ctx context.Context, step types.StepDef, run *types.RunContext) (types.StepResult, error) {
	panic("not implemented")
}
