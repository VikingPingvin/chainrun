package secrets

// Resolver resolves secret values by key.
type Resolver interface {
	Resolve(key string) (value string, found bool)
	ResolveAll() map[string]string
}
