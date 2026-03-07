package config

import "github.com/vikingpingvin/chainrun/internal/types"

// Loader loads and validates workflow configuration.
type Loader interface {
	Load(path string) (*types.Config, error)
	Validate(cfg *types.Config) []ValidationError
}

// ValidationError describes a single configuration validation failure.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
