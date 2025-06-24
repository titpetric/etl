package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/titpetric/etl/server/config"
	"github.com/titpetric/etl/server/internal"
	"github.com/titpetric/etl/server/internal/handler/model"

	_ "github.com/titpetric/etl/server/internal/handler/query"
	_ "github.com/titpetric/etl/server/internal/handler/sql"
)

// Server sets up the HTTP router based on the given endpoints and returns an http.Handler.
func Server(conf *config.Config) (http.Handler, error) {
	return ServerWithEndpoints(conf, conf.Endpoints)
}

// ServerWithEndpoints uses the []*config.Endpoint arguments to return a http.Handler.
func ServerWithEndpoints(conf *config.Config, endpoints []*config.Endpoint) (http.Handler, error) {
	router := mux.NewRouter()
	registeredHandlers := model.Handlers()

	for idx, endpoint := range endpoints {
		handlerType := endpoint.Handler.Type
		factory := registeredHandlers[handlerType]

		if factory == nil {
			log.Printf("unknown handler type: %s, skipping endpoint %d", endpoint.Handler.Type, idx)
			continue
		}

		handler, err := factory.Handler(conf, endpoint)
		if err != nil {
			return nil, fmt.Errorf("error in handler %s: %w", endpoint.Handler.Type, err)
		}

		if len(endpoint.Paths) > 0 {
			// individual path matches and methods
			for _, p := range endpoint.Paths {
				methods := "ANY"
				if len(p.Methods) > 0 {
					methods = strings.Join(p.Methods, ", ")
				}

				// Log route information
				log.Printf("%s (methods: %s, handler: %s, properties: %s)", p.Path, methods, handlerType, string(internal.Marshal(handler)))

				router.Handle(p.Path, handler).Methods(p.Methods...)
			}
		} else {
			// Log route information
			log.Printf("%s (methods: ANY, handler: %s, properties: %s)", endpoint.Path, handlerType, string(internal.Marshal(handler)))

			// full path match
			router.Handle(endpoint.Path, handler)
		}
	}

	return router, nil
}
