package server

import (
	"context"

	"github.com/titpetric/platform"

	"github.com/titpetric/etl/server/config"
	"github.com/titpetric/etl/server/internal/handler"
)

var _ platform.Module = (*Module)(nil)

// Module implements platform.Module for the ETL server.
type Module struct {
	config *config.Config
}

// NewModule creates a new ETL server module.
func NewModule(conf *config.Config) *Module {
	return &Module{
		config: conf,
	}
}

// Name returns the module name.
func (m *Module) Name() string {
	return "etl"
}

// Start initializes the module (no-op for ETL server as mounting handles setup).
func (m *Module) Start(ctx context.Context) error {
	return nil
}

// Stop cleans up the module (no-op for ETL server).
func (m *Module) Stop(ctx context.Context) error {
	return nil
}

// Mount registers the ETL routes on the router.
func (m *Module) Mount(ctx context.Context, r platform.Router) error {
	etlHandler, err := handler.ServerWithEndpoints(m.config, m.config.Endpoints)
	if err != nil {
		return err
	}

	// Mount the ETL handler at the root path
	r.Mount("/", etlHandler)

	return nil
}
