package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/etl/server/config"
)

// TestNewModule verifies that NewModule creates a module with the provided config.
func TestNewModule(t *testing.T) {
	cfg := &config.Config{
		Server: config.Server{
			HttpAddr: ":8080",
		},
	}

	module := NewModule(cfg)

	require.NotNil(t, module)
	require.Equal(t, cfg, module.config)
}

// TestModuleName verifies that the module returns the correct name.
func TestModuleName(t *testing.T) {
	cfg := &config.Config{}
	module := NewModule(cfg)

	require.Equal(t, "etl", module.Name())
}

// TestModuleStart verifies that Start returns nil (no-op).
func TestModuleStart(t *testing.T) {
	cfg := &config.Config{}
	module := NewModule(cfg)

	err := module.Start(context.Background())
	require.NoError(t, err)
}

// TestModuleStartWithCancelledContext verifies that Start handles cancelled context.
func TestModuleStartWithCancelledContext(t *testing.T) {
	cfg := &config.Config{}
	module := NewModule(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := module.Start(ctx)
	require.NoError(t, err)
}

// TestModuleStop verifies that Stop returns nil (no-op).
func TestModuleStop(t *testing.T) {
	cfg := &config.Config{}
	module := NewModule(cfg)

	err := module.Stop(context.Background())
	require.NoError(t, err)
}

// TestModuleStopWithCancelledContext verifies that Stop handles cancelled context.
func TestModuleStopWithCancelledContext(t *testing.T) {
	cfg := &config.Config{}
	module := NewModule(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := module.Stop(ctx)
	require.NoError(t, err)
}

// TestModuleWithNilConfig verifies that NewModule can be created with nil config.
func TestModuleWithNilConfig(t *testing.T) {
	module := NewModule(nil)
	require.NotNil(t, module)
	require.Nil(t, module.config)
}

// TestModuleWithEmptyConfig verifies that NewModule works with empty config.
func TestModuleWithEmptyConfig(t *testing.T) {
	cfg := &config.Config{}
	module := NewModule(cfg)

	require.Equal(t, cfg, module.config)
	require.Equal(t, "etl", module.Name())
}

// TestModuleConfigPreserved verifies that the module preserves the exact config instance.
func TestModuleConfigPreserved(t *testing.T) {
	cfg := &config.Config{
		Server: config.Server{
			HttpAddr: ":9090",
		},
	}
	module := NewModule(cfg)

	require.Equal(t, ":9090", module.config.Server.HttpAddr)
}

// TestModuleNameConsistency verifies that Name returns the same value on multiple calls.
func TestModuleNameConsistency(t *testing.T) {
	cfg := &config.Config{}
	module := NewModule(cfg)

	name1 := module.Name()
	name2 := module.Name()

	require.Equal(t, name1, name2)
	require.Equal(t, "etl", name1)
}
