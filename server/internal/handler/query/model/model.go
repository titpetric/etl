package model

type Config struct {
	Version  string           `yaml:"version"`
	Config   ServiceConfig    `yaml:"config"`
	Inputs   []Input          `yaml:"inputs"`
	Response []ResponseConfig `yaml:"response"`
}

type ServiceConfig struct {
	Service   string     `yaml:"service"`
	Name      string     `yaml:"name"`
	Title     string     `yaml:"title"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
}

type Input struct {
	Name    string `yaml:"name"`
	Title   string `yaml:"title"`
	Default string `yaml:"default,omitempty"`
}

type ResponseConfig struct {
	Produces string      `yaml:"produces"`
	Query    QueryConfig `yaml:"query"`
	With     string      `yaml:"with,omitempty"`
	Want     string      `yaml:"want,omitempty"`
}

type QueryConfig struct {
	Mode string `yaml:"mode,omitempty"`
	SQL  string `yaml:"sql"`
}
