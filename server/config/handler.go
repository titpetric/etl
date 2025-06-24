package config

import "errors"

// Handler represents the handler configuration for an endpoint.
type Handler struct {
	// Type is the handler type to invoke for the endpoint.
	Type string `yaml:"type"`

	// Command mimics os.Args for the handler.
	Command []string `yaml:"command,omitempty"`

	// Library is an optional library parameter for a handler.
	Library string `yaml:"library,omitempty"`

	// Function is an optional function name parameter for a handler.
	Function string `yaml:"function,omitempty"`

	decoder func(interface{}) error
}

// UnmarshalYAML stores the decoder func for later use.
func (e *Handler) UnmarshalYAML(decoder func(interface{}) error) error {
	e.decoder = decoder
	type alias Handler
	return decoder((*alias)(e))
}

// Decode provides extended config functionality to handlers.
func (e *Handler) Decode(out any) error {
	if e.decoder != nil {
		return e.decoder(out)
	}
	return errors.New("raw node decoder unset")
}
