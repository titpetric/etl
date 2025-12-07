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

	// Methods contains the HTTP request methods. If omitted, all methods are considered.
	Methods []string `yaml:"methods"`

	// Paths contains multiple paths and methods for the handler.
	Paths []Path `yaml:"paths"`

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

// Storage type configures database Driver and DSN values.
type Storage struct {
	// Driver configures the database driver to use (sqlite, mysql, pgx).
	Driver string `yaml:"driver"`

	// DSN configures the connection string for the driver.
	DSN string `yaml:"dsn"`
}
