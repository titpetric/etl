package model

import (
	"fmt"
	"io/ioutil"
)

func Load(filepath string) (*Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	cfg := &Config{}
	if err := Decode(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return cfg, nil
}
