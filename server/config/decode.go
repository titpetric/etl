package config

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/huandu/go-clone"
)

// Decode unmarshals YAML data into a Config struct.
// If fsys is provided, it will also process any included files.
// Environment variables (ETL_DB_DRIVER, ETL_DB_DSN) override YAML config.
func Decode(fsys fs.FS, data []byte) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	// Apply environment variable overrides for storage config
	applyStorageEnvOverrides(cfg)

	// Process includes if fsys is provided
	if fsys != nil && len(cfg.Include) > 0 {
		for _, includePath := range cfg.Include {
			// Load the included file from the filesystem
			includeData, err := fs.ReadFile(fsys, includePath)
			if err != nil {
				return nil, fmt.Errorf("failed to load include file %s: %w", includePath, err)
			}

			// Decode the included config recursively
			includeCfg, err := Decode(fsys, includeData)
			if err != nil {
				return nil, fmt.Errorf("failed to decode include file %s: %w", includePath, err)
			}

			// Merge endpoints
			cfg.Endpoints = append(cfg.Endpoints, includeCfg.Endpoints...)

			// Merge server config (later includes override)
			if includeCfg.Server.HttpAddr != "" {
				cfg.Server.HttpAddr = includeCfg.Server.HttpAddr
			}
			if includeCfg.Server.GrpcAddr != "" {
				cfg.Server.GrpcAddr = includeCfg.Server.GrpcAddr
			}
			if includeCfg.Server.Features != nil {
				if cfg.Server.Features == nil {
					cfg.Server.Features = make(map[string]bool)
				}
				for k, v := range includeCfg.Server.Features {
					cfg.Server.Features[k] = v
				}
			}

			// Merge storage (later includes override)
			if includeCfg.Storage != nil {
				if cfg.Storage == nil {
					cfg.Storage = clone.Clone(includeCfg.Storage).(*Storage)
				} else {
					if includeCfg.Storage.Driver != "" {
						cfg.Storage.Driver = includeCfg.Storage.Driver
					}
					if includeCfg.Storage.DSN != "" {
						cfg.Storage.DSN = includeCfg.Storage.DSN
					}
				}
			}
		}
	}

	return cfg, nil
}

// applyStorageEnvOverrides applies environment variable overrides for storage configuration.
// Environment variables: ETL_DB_DRIVER, ETL_DB_DSN
func applyStorageEnvOverrides(cfg *Config) {
	if cfg.Storage == nil {
		cfg.Storage = &Storage{}
	}

	if driver := os.Getenv("ETL_DB_DRIVER"); driver != "" {
		cfg.Storage.Driver = driver
	}

	if dsn := os.Getenv("ETL_DB_DSN"); dsn != "" {
		cfg.Storage.DSN = dsn
	}
}
