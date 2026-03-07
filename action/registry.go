package action

// Registry stores and retrieves Executor factories by action type name.
type Registry interface {
	Register(typeName string, factory ExecutorFactory)
	Get(typeName string) (Executor, bool)
}

// mapRegistry is the concrete Registry backed by a plain map.
type mapRegistry struct {
	factories map[string]ExecutorFactory
}

// NewRegistry returns a new empty Registry.
func NewRegistry() Registry {
	return &mapRegistry{
		factories: make(map[string]ExecutorFactory),
	}
}

// Register associates typeName with the given factory.
func (r *mapRegistry) Register(typeName string, factory ExecutorFactory) {
	r.factories[typeName] = factory
}

// Get returns a new Executor for typeName, or false if not found.
func (r *mapRegistry) Get(typeName string) (Executor, bool) {
	factory, ok := r.factories[typeName]
	if !ok {
		return nil, false
	}
	return factory(), true
}
