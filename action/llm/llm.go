package llm

import (
	"context"
	"net/http"

	"github.com/vikingpingvin/chainrun/internal/types"
	"github.com/vikingpingvin/chainrun/secrets"
	"github.com/vikingpingvin/chainrun/template"
)

// Executor calls an LLM provider (Anthropic for MVP).
type Executor struct {
	client   *http.Client
	secrets  secrets.Resolver
	renderer template.Renderer
}

// New returns a new LLM Executor.
func New(client *http.Client, secretsRes secrets.Resolver, renderer template.Renderer) *Executor {
	return &Executor{client: client, secrets: secretsRes, renderer: renderer}
}

// Execute calls the configured LLM provider with the prompt in step.
func (e *Executor) Execute(ctx context.Context, step types.StepDef, run *types.RunContext) (types.StepResult, error) {
	panic("not implemented")
}
