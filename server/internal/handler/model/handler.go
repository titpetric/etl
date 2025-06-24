package model

import (
	"net/http"

	"github.com/titpetric/etl/server/config"
)

// Handler is an interface that all endpoint handlers must implement.
type Handler interface {
	// Type will return the registered type name of a handler.
	Type() string

	// Handler needs the global config for settings optional to endpoints.
	// Particularly storage credentials can be customized on endpoints.
	Handler(*config.Config, *config.Endpoint) (http.Handler, error)
}

var registeredHandlers = make(map[string]Handler)

// Register registers a new handler.
func Register(handler Handler) {
	registeredHandlers[handler.Type()] = handler
}

// Handlers returns all registered handlers.
func Handlers() map[string]Handler {
	return registeredHandlers
}
