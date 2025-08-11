package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/titpetric/etl/server/config"
	"github.com/titpetric/etl/server/config/loader"
	"github.com/titpetric/etl/server/internal/handler"
)

// Start will load the config and start a HTTP server.
func Start(ctx context.Context) error {
	conf, err := NewConfig()
	if err != nil {
		return err
	}

	return Server(ctx, conf)
}

// NewConfig will load the config from etl.yml.
func NewConfig() (*config.Config, error) {
	return loader.Load("etl.yml")
}

// NewHandler is here so it can be used in other routers.
func NewHandler() (http.Handler, error) {
	conf, err := NewConfig()
	if err != nil {
		return nil, err
	}
	return handler.Server(conf)
}

// Server starts a HTTP server using the provided config.
func Server(ctx context.Context, conf *config.Config) error {
	handler, err := handler.Server(conf)
	if err != nil {
		return fmt.Errorf("error creating server handler: %w", err)
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
