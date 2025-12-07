package main

import (
	"path/filepath"
	"testing"

	"github.com/titpetric/etl/server/config/loader"
)

func TestValidateConfig(t *testing.T) {
	// Load config with includes
	configPath := filepath.Join(".", "etl.yml")
	cfg, err := loader.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Expected: 14 endpoints
	// Query side: 4 users + 1 login + 3 orders = 8
	// Command side: 3 users + 2 login + 1 order = 6
	// CQRS split: 8 read + 6 write = 14 total endpoints
	expectedEndpoints := 14
	if len(cfg.Endpoints) != expectedEndpoints {
		t.Fatalf("Expected %d endpoints, got %d", expectedEndpoints, len(cfg.Endpoints))
	}

	// Validate server config
	if cfg.Server.HttpAddr == "" {
		t.Error("Server HTTP Address is empty")
	}

	if cfg.Storage.Driver == "" {
		t.Error("Storage Driver is empty")
	}

	// Validate each endpoint has methods
	for i, ep := range cfg.Endpoints {
		if ep.Path == "" {
			t.Errorf("Endpoint %d has empty path", i+1)
		}
		if len(ep.Methods) == 0 {
			t.Errorf("Endpoint %d (%s) has no methods", i+1, ep.Path)
		}
	}
}
