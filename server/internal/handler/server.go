package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/titpetric/platform"

	"github.com/titpetric/etl/server/config"
	"github.com/titpetric/etl/server/internal"
	"github.com/titpetric/etl/server/internal/handler/model"

	_ "github.com/titpetric/etl/server/internal/handler/query"
	_ "github.com/titpetric/etl/server/internal/handler/request"
	_ "github.com/titpetric/etl/server/internal/handler/sql"
)

// Mount uses the []*config.Endpoint arguments to populate routes.
func Mount(router platform.Router, conf *config.Config, endpoints []*config.Endpoint) error {
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
			return fmt.Errorf("error in handler %s: %w", endpoint.Handler.Type, err)
		}

		methods := strings.Join(endpoint.Methods, ", ")
		if methods == "" {
			methods = "ANY"
		}

		log.Printf("%s (methods: %s, handler: %s, properties: %s)", endpoint.Path, methods, handlerType, string(internal.Marshal(handler)))

		if len(endpoint.Methods) == 0 {
			router.Handle(endpoint.Path, handler)
		} else {
			for _, method := range endpoint.Methods {
				router.Method(method, endpoint.Path, handler)
			}
		}
	}

	return nil
}
