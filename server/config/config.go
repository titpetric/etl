// Package config provides the configuration structure and parsing for the ETL server.
//
// # Configuration Structure
//
// The top-level configuration defines the server settings, optional storage configuration,
// and a list of endpoints. Each endpoint routes requests to a specific handler.
//
// # Example
//
//	```yaml
//	server:
//	  http: ":3000"
//	  grpc: ":50051"
//	  features:
//	    feature_flag: true
//	storage:
//	  driver: pgx
//	  dsn: "postgres://localhost/db"
//	endpoints:
//	  - path: /api
//	    methods: [GET, POST]
//	    handler:
//	      driver: http
//	      url: "http://backend:8080"
//	```
package config

// Config represents the overall configuration structure, which includes a server and multiple endpoints.
type Config struct {
	// Server is used to configure the service.
	Server Server `yaml:"server"`

	// Storage configures the storage default.
	Storage *Storage `yaml:"storage"`

	// Include specifies additional config files to include and merge.
	Include []string `yaml:"include,omitempty"`

	// Endpoints contain a list of configured endpoints.
	Endpoints []*Endpoint `yaml:"endpoints"`
}

// Endpoint represents an endpoint configuration with a path and handler.
type Endpoint struct {
	// Path contains the path passed to the handler. It's a request path e.g. `/api`.
	Path string `yaml:"path"`

	// Name describes what the endpoint returns.
	Name string `yaml:"name"`

	// Methods contains the HTTP request methods. If omitted, all methods are considered.
	Methods []string `yaml:"methods"`

	// Handler contains configuration related to the particular handler for the request.
	Handler Handler `yaml:"handler"`
}

// Path is the combination of Method and Path.
type Path struct {
	// Methods contains the HTTP request methods. If omitted, all methods are considered.
	Methods []string `yaml:"methods"`

	// Path contains the request path for the router, e.g. `/users/{id}`.
	Path string `yaml:"path"`
}

// Server represents the server configuration with an address.
type Server struct {
	// HttpAddr contains the address the server should listen on. Example: ":3123".
	HttpAddr string `yaml:"http"`

	// GrpcAddr contains the address the server should listen on. Example: ":50051".
	GrpcAddr string `yaml:"grpc"`

	// Features contains feature flags available for conditional query execution.
	Features map[string]bool `yaml:"features"`
}

// Storage type configures database connection DSN.
// The driver is automatically derived from the DSN connection string.
type Storage struct {
	// DSN configures the connection string for the database.
	// Supports mysql://, postgres://, postgresql://, sqlite://, and driver-specific formats.
	DSN string `yaml:"dsn"`
}
