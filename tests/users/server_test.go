//go:build integration
// +build integration

package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestUsersAPI(t *testing.T) {
	// Generate random database filename
	dbFile := filepath.Join(t.TempDir(), fmt.Sprintf("test_%d.db", rand.Int63()))

	// Cleanup is automatic with t.TempDir() but we'll be explicit
	t.Cleanup(func() {
		os.Remove(dbFile)
	})

	// Setup database
	cmd := exec.Command("etl", "query", "users.sql")
	cmd.Env = append(os.Environ(),
		"ETL_DB_DRIVER=sqlite",
		fmt.Sprintf("ETL_DB_DSN=file:%s", dbFile),
	)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to setup database: %v", err)
	}

	// Start server
	serverCmd := exec.Command("etl", "server")
	serverCmd.Env = append(os.Environ(),
		"ETL_DB_DRIVER=sqlite",
		fmt.Sprintf("ETL_DB_DSN=file:%s", dbFile),
	)
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	if err := serverCmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	t.Cleanup(func() {
		serverCmd.Process.Kill()
		serverCmd.Wait()
	})

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Run venom tests
	venomCmd := exec.Command("venom", "run", "users-test.yml")
	venomCmd.Stdout = os.Stdout
	venomCmd.Stderr = os.Stderr

	if err := venomCmd.Run(); err != nil {
		t.Errorf("venom tests failed: %v", err)
	}
}
