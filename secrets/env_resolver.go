package secrets

import (
	"os"
	"strings"
)

type EnvResolver struct{}

func NewEnvResolver() *EnvResolver {
	return &EnvResolver{}
}

func (r *EnvResolver) Resolve(key string) (string, bool) {
	v := os.Getenv(key)
	return v, v != ""
}

func (r *EnvResolver) ResolveAll() map[string]string {
	m := map[string]string{}
	for _, e := range os.Environ() {
		k, v, _ := strings.Cut(e, "=")
		m[k] = v
	}
	return m
}
