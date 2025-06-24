package server

import (
	"context"
	"log"
	"net/http"

	"github.com/titpetric/etl/server/config"
	"github.com/titpetric/etl/server/config/loader"
	"github.com/titpetric/etl/server/internal/handler"
)

// Start will load a config and start the service.
func Start(ctx context.Context) error {
	conf, err := loader.Load("etl.yml")
	if err != nil {
		return err
	}

	return StartWithConfig(ctx, conf)
}

// Start is the entry point for a service lifecycle.
func StartWithConfig(ctx context.Context, conf *config.Config) error {
	handler, err := handler.Server(conf)
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    conf.Server.HttpAddr,
		Handler: handler,
	}

	// Start HTTP server in a separate goroutine
	go func() {
		log.Println("Starting HTTP server on address:", conf.Server.HttpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve HTTP server: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	log.Println("Shutting down servers...")

	// Shutdown HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Servers gracefully stopped")
	return nil
}
