package action

import (
	"context"

	"github.com/vikingpingvin/chainrun/internal/types"
)

// Executor runs a single workflow step.
//
// Contracts:
//   - error != nil  → hard failure (executor could not run)
//   - error == nil, result.Status == Failed → step ran but produced a failure
//     (e.g. non-zero exit code, HTTP 4xx/5xx)
type Executor interface {
	Execute(ctx context.Context, step types.StepDef, run *types.RunContext) (types.StepResult, error)
}

// ExecutorFactory is a constructor function for an Executor.
type ExecutorFactory func() Executor
