package query

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/titpetric/etl/server/config"
	"github.com/titpetric/etl/server/internal/handler/query/model"
)

// Handler represents a handler that executes a SQL query based on configuration.
type Handler struct {
	Storage *config.Storage

	Filename string
	Query    string

	Parameters map[string]any
}

// NewHandler creates a new Handler instance.
func NewHandler() *Handler {
	return &Handler{}
}

// ServeHTTP handles the HTTP request, executes the SQL query, and writes the response.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Load the configuration from the provided file
	conf, err := model.Load(h.Filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error loading config", http.StatusInternalServerError)
		return
	}

	// Process the SQL queries and gather results
	queryParams := h.prepareQueryParams(r)
	results, err := h.eval(conf, queryParams)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error executing SQL query", http.StatusInternalServerError)
		return
	}

	// Write the results as a JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Println(err)
		http.Error(w, "Error encoding query results to JSON", http.StatusInternalServerError)
	}
}

// Type returns the type of the handler.
func (h *Handler) Type() string {
	return "query"
}

// Handler creates a new instance of Handler and returns the http.Handler.
func (h *Handler) Handler(conf *config.Config, endpoint *config.Endpoint) (http.Handler, error) {
	handle := NewHandler()
	err := endpoint.Handler.Decode(handle)
	if err != nil {
		return nil, fmt.Errorf("error decoding config: %w", err)
	}

	// Copy global storage config if endpoint has it unset.
	if handle.Storage == nil {
		handle.Storage = conf.Storage
	}

	return handle, nil
}

func (h *Handler) eval(conf *model.Config, queryParams map[string]any) (map[string]any, error) {
	db, err := sqlx.Open(h.Storage.Driver, h.Storage.DSN)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}
	defer db.Close()

	results := make(map[string]any)

	for _, response := range conf.Response {
		// Execute the SQL query with proper parameter handling
		rows, err := db.NamedQuery(response.Query.SQL, queryParams)
		if err != nil {
			return nil, fmt.Errorf("error executing SQL query: %w", err)
		}
		defer rows.Close()

		var result []map[string]string
		for rows.Next() {
			row := make(map[string]any)
			if err := rows.MapScan(row); err != nil {
				return nil, fmt.Errorf("error processing query results: %w", err)
			}

			res := make(map[string]string, len(row))
			for k, v := range row {
				res[strings.ToLower(k)] = dbValue(v)
			}

			result = append(result, res)
		}

		results[response.Produces] = result
	}

	return results, nil
}

// prepareQueryParams prepares the query parameters from the request.
func (h *Handler) prepareQueryParams(r *http.Request) map[string]any {
	queryParams := make(map[string]any)
	for k, v := range h.Parameters {
		queryParams[k] = v
	}

	// Get path parameters from chi
	rctx := chi.RouteContext(r.Context())
	if rctx != nil {
		for i, key := range rctx.URLParams.Keys {
			if i < len(rctx.URLParams.Values) {
				queryParams[key] = rctx.URLParams.Values[i]
			}
		}
	}

	for k, v := range r.URL.Query() {
		queryParams[k] = v[0]
	}

	// If the method is not GET and the request body is not nil, scan the body into queryParams
	if r.Method != http.MethodGet && r.Body != nil {
		defer r.Body.Close()
		bodyParams := map[string]any{}
		if err := json.NewDecoder(r.Body).Decode(&bodyParams); err == nil {
			for k, v := range bodyParams {
				queryParams[k] = v
			}
		}
	}

	return queryParams
}
