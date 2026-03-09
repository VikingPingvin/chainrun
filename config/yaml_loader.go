package config

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/vikingpingvin/chainrun/internal/types"
)

// YAMLLoader loads workflow configuration from a YAML file.
type YAMLLoader struct{}

// NewYAMLLoader returns a new YAMLLoader.
func NewYAMLLoader() *YAMLLoader {
	return &YAMLLoader{}
}

// Load reads and parses a YAML config file at the given path.
func (l *YAMLLoader) Load(path string) (*types.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg types.Config
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}

	return &cfg, nil
}

// Validate checks a loaded Config for logical errors.
func (l *YAMLLoader) Validate(cfg *types.Config) []ValidationError {
	var errs []ValidationError

	seen := map[string]bool{}
	for i, wf := range cfg.Workflows {
		prefix := fmt.Sprintf("workflows[%d]", i)

		if wf.Name == "" {
			errs = append(errs, ValidationError{Field: prefix + ".name", Message: "must not be empty"})
		} else if seen[wf.Name] {
			errs = append(errs, ValidationError{Field: prefix + ".name", Message: fmt.Sprintf("duplicate workflow name %q", wf.Name)})
		} else {
			seen[wf.Name] = true
		}

		if wf.Trigger.Type == "" {
			errs = append(errs, ValidationError{Field: prefix + ".trigger.type", Message: "must not be empty"})
		}

		stepIDs := map[string]bool{}
		for j, step := range wf.Steps {
			sp := fmt.Sprintf("%s.steps[%d]", prefix, j)

			if step.ID == "" {
				errs = append(errs, ValidationError{Field: sp + ".id", Message: "must not be empty"})
			} else if stepIDs[step.ID] {
				errs = append(errs, ValidationError{Field: sp + ".id", Message: fmt.Sprintf("duplicate step id %q", step.ID)})
			} else {
				stepIDs[step.ID] = true
			}

			if step.ActionType() == "" {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("workflows[%d].steps[%d]", i, j),
					Message: "step must have exactly one executor block (shell, http, llm, file, notify)",
				})
			}
		}
	}

	return errs
}
