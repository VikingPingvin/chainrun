package cron

import (
	"context"

	"github.com/vikingpingvin/chainrun/internal/types"
)

// Source fires TriggerEvents on a cron schedule using robfig/cron/v3.
type Source struct {
	def types.TriggerDef
}

// New constructs a cron Source from a TriggerDef.
func New(def types.TriggerDef) (*Source, error) {
	panic("not implemented")
}

// Start begins the cron schedule and returns the event channel.
func (s *Source) Start(ctx context.Context) (<-chan types.TriggerEvent, error) {
	panic("not implemented")
}

// Stop halts the cron scheduler.
func (s *Source) Stop() error {
	panic("not implemented")
}

// Name returns "cron".
func (s *Source) Name() string {
	return "cron"
}
