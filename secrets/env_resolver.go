package secrets

// EnvResolver resolves secrets from environment variables.
type EnvResolver struct{}

// NewEnvResolver returns a new EnvResolver.
func NewEnvResolver() *EnvResolver {
	return &EnvResolver{}
}

// Resolve returns the value of the environment variable named by key.
func (r *EnvResolver) Resolve(key string) (string, bool) {
	panic("not implemented")
}

// ResolveAll returns all environment variables as a map.
func (r *EnvResolver) ResolveAll() map[string]string {
	panic("not implemented")
}
