package trigger

import (
	"context"

	"github.com/vikingpingvin/chainrun/internal/types"
)

// Source is a trigger that emits TriggerEvents on a channel.
type Source interface {
	Start(ctx context.Context) (<-chan types.TriggerEvent, error)
	Stop() error
	Name() string
}

// SourceFactory constructs a Source from a TriggerDef.
type SourceFactory func(def types.TriggerDef) (Source, error)
