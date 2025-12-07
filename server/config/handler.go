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

	// Request specifies the upstream request handler path to invoke (for request type handlers).
	// Supports path parameters with brace syntax: /api/users/{id}
	Request []*Request `yaml:"request,omitempty"`

	// Query is a single SQL query (deprecated, use Queries instead).
	Query string `yaml:"query,omitempty"`

	// Queries is a list of queries to execute in sequence (pipeline execution).
	Queries []*QueryDef `yaml:"queries,omitempty"`

	// Single returns a single object instead of array (applies to all queries).
	Single bool `yaml:"single"`

	// Parameters are static parameters merged with request parameters.
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`

	// Transaction configures transactional behavior for write operations.
	Transaction *Transaction `yaml:"transaction,omitempty"`

	// Cache configures response caching behavior.
	Cache *Cache `yaml:"cache,omitempty"`

	// RateLimit configures rate limiting for the endpoint.
	RateLimit *RateLimit `yaml:"rateLimit,omitempty"`

	// Response configures the response format and headers.
	Response *Response `yaml:"response,omitempty"`

	decoder func(interface{}) error
}

// QueryDef represents a single query in a query pipeline.
type QueryDef struct {
	// Query is the SQL query to execute.
	Query string `yaml:"query"`

	// As is the path where the result should be stored in the response.
	// If empty, result is merged into the current scope.
	As string `yaml:"as,omitempty"`

	// If is a condition expression evaluated against the current scope.
	// Query only executes if condition is true.
	If string `yaml:"if,omitempty"`

	// For is a loop expression to execute the query for each item.
	// Format: (idx, item) in items
	// Results are placed at the path specified in As.
	For string `yaml:"for,omitempty"`
}

// Transaction configures transactional behavior for write operations.
type Transaction struct {
	// Enabled indicates whether transactions should be used.
	Enabled bool `yaml:"enabled"`

	// Retries specifies the number of times to retry a failed transaction commit.
	Retries int `yaml:"retries"`

	// RetryDelayMs specifies the delay in milliseconds between retry attempts.
	RetryDelayMs int `yaml:"retryDelayMs"`
}

// Cache configures response caching behavior.
type Cache struct {
	// Enabled indicates whether response caching is enabled.
	Enabled bool `yaml:"enabled"`

	// TTLSeconds specifies the cache time-to-live in seconds.
	TTLSeconds int `yaml:"ttlSeconds"`

	// KeyPattern specifies the cache key pattern (supports path and query parameters).
	KeyPattern string `yaml:"keyPattern"`
}

// RateLimit configures rate limiting for the endpoint.
// Example: rate: 100, per: "1m" means 100 requests per minute (default).
// Per defaults to "1m" if not specified, so rate defines requests per minute.
type RateLimit struct {
	// Enabled indicates whether rate limiting is enabled.
	Enabled bool `yaml:"enabled" json:"enabled"`

	// Rate specifies the number of requests allowed in the time period.
	// Defaults to requests per minute (when Per is unset).
	// Example: 100 = 100 requests per minute
	Rate int `yaml:"rate" json:"rate"`

	// Per specifies the time interval for the rate limit (e.g., "30s", "5m", "1h").
	// Supports all time.Duration formats. Defaults to "1m" (one minute) if not specified.
	// Example: "30s" means 30 seconds, "5m" means 5 minutes
	Per string `yaml:"per" json:"per"`
}

// Request configures the request format and headers.
type Request struct {
	// Path is the endpoint url to invoke.
	Path string `yaml:"path"`

	// As is the key for the decoded json value
	As string `yaml:"as"`

	// Method sets the request method. Default to GET.
	Method string `yaml:"method,omitempty"`

	// Headers are custom HTTP request headers.
	Headers map[string]string `yaml:"headers,omitempty"`

	// Body is a custom HTTP request body.
	Body string `yaml:"body,omitempty"`
}

// Response configures the response format and headers.
type Response struct {
	// Headers are custom HTTP response headers.
	// If not specified, Content-Type defaults to application/json for JSON responses
	// or text/html; charset=utf-8 for template responses.
	Headers map[string]string `yaml:"headers,omitempty"`

	// Template is a VueGo template string for formatting the response.
	// If specified, the response will be rendered using this template.
	Template string `yaml:"template,omitempty"`
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
