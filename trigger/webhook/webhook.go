package webhook

import (
	"context"
	"net/http"

	"github.com/vikingpingvin/chainrun/internal/types"
)

// Source listens on an HTTP endpoint and emits TriggerEvents for each request.
type Source struct {
	def    types.TriggerDef
	server *http.Server
}

// New constructs a webhook Source from a TriggerDef.
func New(def types.TriggerDef) (*Source, error) {
	panic("not implemented")
}

// Start starts the HTTP server and returns the event channel.
func (s *Source) Start(ctx context.Context) (<-chan types.TriggerEvent, error) {
	panic("not implemented")
}

// Stop shuts down the HTTP server.
func (s *Source) Stop() error {
	panic("not implemented")
}

// Name returns "webhook".
func (s *Source) Name() string {
	return "webhook"
}
