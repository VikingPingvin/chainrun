package file

import (
	"context"
	"os"

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
	switch step.FileOp {
	case "write":
		return e.writeFile(ctx, step, step.FilePath, step.FileContent)
	}

	return types.StepResult{
		StepID: step.ID,
		Status: types.StepStatusFailed,
		Error:  "unsupported file operation",
	}, nil
}

func (e *Executor) writeFile(ctx context.Context, step types.StepDef, path string, fileContent string) (types.StepResult, error) {
	err := os.WriteFile(path, []byte(fileContent), 0644)
	if err != nil {
		return types.StepResult{
			StepID: step.ID,
			Status: types.StepStatusFailed,
			Error:  err.Error(),
		}, err
	}
	return types.StepResult{
		StepID: step.ID,
		Status: types.StepStatusSuccess,
	}, nil
}
