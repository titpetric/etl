package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"

	"github.com/titpetric/vuego"

	"github.com/titpetric/etl/server/config"
	"github.com/titpetric/etl/server/internal/handler/model"
	"github.com/titpetric/etl/server/internal/handler/sql"
)

// Handler represents a handler that invokes an upstream request and formats the response.
type Handler struct {
	RequestPath string
	Request     []*config.Request
	Response    *config.Response
	RateLimit   *config.RateLimit
	Cache       *config.Cache

	// Internal state
	template vuego.Template
	router   http.Handler
}

// NewHandler creates a new Handler.
func NewHandler() *Handler {
	return &Handler{
		template: vuego.New(),
	}
}

func (h *Handler) EvaluateRequest(r *http.Request) (map[string]any, error) {
	ctx := r.Context()
	result := make(map[string]any)

	for _, reqSpec := range h.Request {
		var body io.Reader

		key := reqSpec.As
		method := reqSpec.Method
		if method == "" {
			method = http.MethodGet
		}
		if reqSpec.Body != "" {
			body = bytes.NewReader([]byte(reqSpec.Body))
		}

		upstreamPath := h.buildUpstreamPath(r, reqSpec.Path)
		log.Println("request:", method, upstreamPath)

		// --- Updated: use NewRequestWithContext ---
		upstreamReq, err := http.NewRequestWithContext(ctx, method, upstreamPath, body)
		if err != nil {
			return nil, err
		}

		for k, v := range reqSpec.Headers {
			upstreamReq.Header.Set(k, v)
		}

		resp, err := http.DefaultClient.Do(upstreamReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var responseBody any
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			return nil, err
		}

		result[key] = responseBody
	}

	return result, nil
}

// ServeHTTP invokes the upstream request handler and returns the response.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check rate limiting if enabled
	if h.RateLimit != nil && h.RateLimit.Enabled {
		// Rate limiting would be applied here if needed
		// For now, we delegate to the upstream handler
	}

	// Build the upstream request URL by substituting path parameters
	data, err := h.EvaluateRequest(r)
	if err != nil {
		log.Println("error creating upstream request", err)
		http.Error(w, "error creating upstream request", http.StatusInternalServerError)
		return
	}

	// Add custom response headers if configured
	if h.Response != nil && h.Response.Headers != nil {
		for key, value := range h.Response.Headers {
			w.Header().Set(key, value)
		}
	}

	// Process template if configured
	if h.Response != nil && h.Response.Template != "" {
		err := h.renderTemplateResponse(ctx, w, data)
		if err != nil {
			log.Printf("error rendering template: %v", err)
			http.Error(w, "Template rendering failed", http.StatusInternalServerError)
			return
		}
		return
	}

	// Otherwise, write captured response as-is
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding: %v", err)
		http.Error(w, "Failed encoding response", http.StatusInternalServerError)
		return
	}
}

// renderTemplateResponse parses the JSON response and renders it using VueGo template
func (h *Handler) renderTemplateResponse(ctx context.Context, w io.Writer, data map[string]any) error {
	var buf bytes.Buffer
	if err := h.template.Fill(data).RenderString(ctx, &buf, h.Response.Template); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())
	return err
}

// buildUpstreamPath constructs the upstream request path by substituting path parameters.
// Supports brace syntax: /api/users/{id}
func (h *Handler) buildUpstreamPath(r *http.Request, pattern string) string {
	vars := mux.Vars(r)

	// Replace {param} with actual values
	for key, value := range vars {
		placeholder := "{" + key + "}"
		pattern = strings.ReplaceAll(pattern, placeholder, value)
	}

	// Also support :param syntax for consistency
	re := regexp.MustCompile(`:(\w+)`)
	pattern = re.ReplaceAllStringFunc(pattern, func(match string) string {
		paramName := match[1:] // Remove leading ':'
		if value, ok := vars[paramName]; ok {
			return value
		}
		return match
	})

	return pattern
}

// Type returns the handler type.
func (h *Handler) Type() string {
	return "request"
}

// Handler creates a new handler instance from endpoint config.
func (h *Handler) Handler(conf *config.Config, endpoint *config.Endpoint) (http.Handler, error) {
	handle := NewHandler()
	err := endpoint.Handler.Decode(handle)
	if err != nil {
		return nil, fmt.Errorf("error decoding config: %w", err)
	}

	// Copy handler configuration
	handle.RequestPath = endpoint.Path
	handle.Response = endpoint.Handler.Response
	handle.RateLimit = endpoint.Handler.RateLimit
	handle.Cache = endpoint.Handler.Cache

	if handle.RequestPath == "" {
		return nil, fmt.Errorf("request handler requires 'request' field")
	}

	// Build a temporary router for the upstream handler
	// We need to get the SQL handler factory and invoke it with a dummy endpoint
	sqlFactory := &sql.Handler{}
	handler, err := sqlFactory.Handler(conf, endpoint)
	if err != nil {
		return nil, fmt.Errorf("error creating upstream handler: %w", err)
	}

	handle.router = handler

	return handle, nil
}

func init() {
	model.Register(&Handler{})
}
