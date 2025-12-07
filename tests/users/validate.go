package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/titpetric/etl/server/config/loader"
)

func main() {
	// Get the current working directory (where the test files are)
	testDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Load config with includes
	configPath := filepath.Join(testDir, "etl.yml")
	cfg, err := loader.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate the loaded config
	fmt.Printf("Server HTTP Address: %s\n", cfg.Server.HttpAddr)
	fmt.Printf("Storage Driver: %s\n", cfg.Storage.Driver)
	fmt.Printf("Total Endpoints: %d\n", len(cfg.Endpoints))

	// Print each endpoint
	for i, ep := range cfg.Endpoints {
		methods := getMethod(ep.Methods)
		fmt.Printf("  [%d] %s (methods: %s)\n", i+1, ep.Path, methods)
	}

	// Check features
	if cfg.Server.Features != nil {
		fmt.Printf("Features: %v\n", cfg.Server.Features)
	}

	// Expected: 13 endpoints
	// Query side: 4 users + 1 login + 3 orders = 8
	// Command side: 3 users + 2 login + 1 order = 6
	// But CQRS split: 8 read + 6 write = 14 total endpoints
	expectedEndpoints := 14
	if len(cfg.Endpoints) != expectedEndpoints {
		log.Fatalf("Expected %d endpoints, got %d", expectedEndpoints, len(cfg.Endpoints))
	}

	fmt.Println("\n✓ Include validation passed!")
	os.Exit(0)
}

func getMethod(methods []string) string {
	if len(methods) == 0 {
		return "GET"
	}
	result := ""
	for i, m := range methods {
		if i > 0 {
			result += ", "
		}
		result += m
	}
	return result
}
