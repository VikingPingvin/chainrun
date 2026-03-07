package manual

import (
	"context"

	"github.com/vikingpingvin/chainrun/internal/types"
)

// Source fires a single TriggerEvent immediately upon Start.
type Source struct {
	def types.TriggerDef
}

// New constructs a manual Source from a TriggerDef.
func New(def types.TriggerDef) (*Source, error) {
	panic("not implemented")
}

// Start fires one event and closes the channel.
func (s *Source) Start(ctx context.Context) (<-chan types.TriggerEvent, error) {
	panic("not implemented")
}

// Stop is a no-op for the manual trigger.
func (s *Source) Stop() error {
	return nil
}

// Name returns "manual".
func (s *Source) Name() string {
	return "manual"
}
