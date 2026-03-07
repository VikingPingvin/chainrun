package trigger

import "github.com/vikingpingvin/chainrun/internal/types"

// Registry stores and retrieves trigger Source factories by trigger type name.
type Registry interface {
	Register(typeName string, factory SourceFactory)
	Get(typeName string, def types.TriggerDef) (Source, bool)
}

// mapRegistry is the concrete Registry backed by a plain map.
type mapRegistry struct {
	factories map[string]SourceFactory
}

// NewRegistry returns a new empty trigger Registry.
func NewRegistry() Registry {
	return &mapRegistry{
		factories: make(map[string]SourceFactory),
	}
}

// Register associates typeName with the given factory.
func (r *mapRegistry) Register(typeName string, factory SourceFactory) {
	panic("not implemented")
}

// Get constructs and returns a Source for the given typeName and def.
func (r *mapRegistry) Get(typeName string, def types.TriggerDef) (Source, bool) {
	panic("not implemented")
}
