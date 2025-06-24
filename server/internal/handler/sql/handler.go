package sql

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/titpetric/etl/server/config"
)

// Handler represents a handler that executes a SQL query from a file.
type Handler struct {
	Storage *config.Storage

	Filename string
	Query    string
	Single   bool

	Parameters map[string]interface{}
}

// NewHandler creates a new Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// ServeHTTP handles the request, reads the SQL query from a file, and executes it.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the SQL query
	var query []byte
	var err error

	switch {
	case h.Query != "":
		query = []byte(h.Query)
	case h.Filename != "":
		query, err = os.ReadFile(h.Filename)
		if err != nil {
			log.Printf("error reading %s: %s", h.Filename, err)
			http.Error(w, "Error reading SQL query file", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "No SQL query defined for path", http.StatusInternalServerError)
	}

	// Open the database connection
	db, err := sqlx.Open(h.Storage.Driver, h.Storage.DSN)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Execute the SQL query
	queryParams := map[string]interface{}{}
	for k, v := range h.Parameters {
		queryParams[k] = v
	}
	for k, v := range mux.Vars(r) {
		queryParams[k] = v
	}
	for k, v := range r.URL.Query() {
		queryParams[k] = v[0]
	}

	rows, err := db.NamedQuery(string(query), queryParams)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error executing SQL query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Process the rows and write the response
	results := []map[string]string{}
	for rows.Next() {
		row := map[string]interface{}{}
		if err := rows.MapScan(row); err != nil {
			log.Println(err)
			http.Error(w, "Error processing query results", http.StatusInternalServerError)
			return
		}

		result := make(map[string]string, len(row))
		for k, v := range row {
			result[strings.ToLower(k)] = dbValue(v)
		}

		results = append(results, result)
	}

	// Coalesce rows to single row if set.
	var result any
	if h.Single && len(results) > 0 {
		result = results[0]
	} else {
		result = results
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Println("Error encoding json response:", err)
		http.Error(w, "Error encoding query results to JSON", http.StatusInternalServerError)
	}
}

// Type returns the type of the handler.
func (h *Handler) Type() string {
	return "sql"
}

// Handler creates a new instance of Handler and returns the http.Handler.
func (h *Handler) Handler(conf *config.Config, endpoint *config.Endpoint) (http.Handler, error) {
	handle := NewHandler()
	err := endpoint.Handler.Decode(handle)
	if err != nil {
		// Log route information and error
		return nil, fmt.Errorf("error decoding config: %w", err)
	}

	// Copy global storage config if endpoint has it unset.
	if handle.Storage == nil {
		handle.Storage = conf.Storage
	}

	return handle, nil
}
