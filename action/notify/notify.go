package notify

import (
	"context"
	"net/http"

	"github.com/vikingpingvin/chainrun/internal/types"
	"github.com/vikingpingvin/chainrun/template"
)

// Executor sends notifications (Slack webhook for MVP).
type Executor struct {
	client   *http.Client
	renderer template.Renderer
}

// New returns a new notify Executor.
func New(client *http.Client, renderer template.Renderer) *Executor {
	return &Executor{client: client, renderer: renderer}
}

// Execute sends the notification defined in step.
func (e *Executor) Execute(ctx context.Context, step types.StepDef, run *types.RunContext) (types.StepResult, error) {
	panic("not implemented")
}
