package http

import (
	"context"
	"net/http"

	"github.com/vikingpingvin/chainrun/internal/types"
	"github.com/vikingpingvin/chainrun/template"
)

// Executor performs HTTP requests.
type Executor struct {
	client   *http.Client
	renderer template.Renderer
}

// New returns a new HTTP Executor.
func New(client *http.Client, renderer template.Renderer) *Executor {
	return &Executor{client: client, renderer: renderer}
}

// Execute performs the HTTP request defined in step.
func (e *Executor) Execute(ctx context.Context, step types.StepDef, run *types.RunContext) (types.StepResult, error) {
	panic("not implemented")
}
