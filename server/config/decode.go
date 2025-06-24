package config

import (
	"github.com/goccy/go-yaml"
)

// Decode unmarshals YAML data into a Config struct.
func Decode(data []byte) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
