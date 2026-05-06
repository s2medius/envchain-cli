package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Backend represents a secret backend configuration.
type Backend struct {
	Type   string            `json:"type"`
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
}

// Config holds the top-level envchain configuration.
type Config struct {
	Version  int       `json:"version"`
	Backends []Backend `json:"backends"`
}

// Load reads and parses the config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that the config has required fields.
func (c *Config) Validate() error {
	if c.Version == 0 {
		return errors.New("config: missing or zero version field")
	}
	for i, b := range c.Backends {
		if b.Type == "" {
			return fmt.Errorf("config: backend[%d] missing type", i)
		}
		if b.Name == "" {
			return fmt.Errorf("config: backend[%d] missing name", i)
		}
	}
	return nil
}
