package config

import "github.com/vikingpingvin/chainrun/internal/types"

// YAMLLoader loads workflow configuration from a YAML file.
type YAMLLoader struct{}

// NewYAMLLoader returns a new YAMLLoader.
func NewYAMLLoader() *YAMLLoader {
	return &YAMLLoader{}
}

// Load reads and parses a YAML config file at the given path.
func (l *YAMLLoader) Load(path string) (*types.Config, error) {
	panic("not implemented")
}

// Validate checks a loaded Config for logical errors.
func (l *YAMLLoader) Validate(cfg *types.Config) []ValidationError {
	panic("not implemented")
}
