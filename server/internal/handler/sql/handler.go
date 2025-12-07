package sql

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/expr-lang/expr"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/time/rate"

	"github.com/titpetric/vuego"

	"github.com/titpetric/etl/server/config"
)

// Handler represents a handler that executes SQL queries from a pipeline.
type Handler struct {
	Storage    *config.Storage
	Filename   string
	Query      string
	Queries    []*config.QueryDef
	Single     bool
	Parameters map[string]interface{}

	// Extensions for advanced features
	Transaction *config.Transaction
	Cache       *config.Cache
	RateLimit   *config.RateLimit
	Response    *config.Response

	// Server features for conditional execution
	Features map[string]bool

	// Internal state for rate limiting and caching
	limiter    *rate.Limiter
	cacheStore map[string]*cacheEntry
	cacheMutex sync.RWMutex
}

// cacheEntry stores cached response with expiration time
type cacheEntry struct {
	data      interface{}
	expiresAt time.Time
}

// NewHandler creates a new Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// ServeHTTP handles the request and executes the query pipeline.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check rate limiting if enabled
	if h.RateLimit != nil && h.RateLimit.Enabled && h.limiter != nil {
		if !h.limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
	}

	// Build cache key if caching is enabled
	var cacheKey string
	if h.Cache != nil && h.Cache.Enabled {
		cacheKey = h.buildCacheKey(r)
		if cached, ok := h.getFromCache(cacheKey); ok {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			json.NewEncoder(w).Encode(cached)
			return
		}
	}

	// Open the database connection
	db, err := sqlx.Open(h.Storage.Driver, h.Storage.DSN)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Collect parameters from various sources
	queryParams := h.collectParameters(r)

	// Execute query pipeline
	result, execErr := h.executePipeline(db, queryParams)
	if execErr != nil {
		log.Printf("Error executing query pipeline: %v", execErr)
		http.Error(w, "Error executing query pipeline", http.StatusInternalServerError)
		return
	}

	// Cache the result if caching is enabled
	if h.Cache != nil && h.Cache.Enabled && cacheKey != "" {
		h.setInCache(cacheKey, result)
	}

	// Set response headers
	h.setResponseHeaders(w)

	// Add rate limit headers if rate limiting is enabled
	if h.RateLimit != nil && h.RateLimit.Enabled {
		h.setRateLimitHeaders(w)
	}

	// Add cache headers if caching is enabled
	if h.Cache != nil && h.Cache.Enabled && cacheKey != "" {
		w.Header().Set("X-Cache", "MISS")
	}

	// Render response: template if specified, otherwise JSON
	if h.Response != nil && h.Response.Template != "" {
		h.renderTemplateResponse(w, result)
	} else {
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Println("Error encoding json response:", err)
			http.Error(w, "Error encoding query results to JSON", http.StatusInternalServerError)
		}
	}
}

// collectParameters gathers parameters from multiple sources
func (h *Handler) collectParameters(r *http.Request) map[string]interface{} {
	params := make(map[string]interface{})

	// 1. Static parameters from config
	for k, v := range h.Parameters {
		params[k] = v
	}

	// 2. Path parameters
	rctx := chi.RouteContext(r.Context())
	if rctx != nil {
		for i, key := range rctx.URLParams.Keys {
			if i < len(rctx.URLParams.Values) {
				params[key] = rctx.URLParams.Values[i]
			}
		}
	}

	// 3. Query string parameters
	for k, v := range r.URL.Query() {
		params[k] = v[0]
	}

	// 4. Request body parameters (for POST/PUT)
	if r.Method == "POST" || r.Method == "PUT" {
		var bodyParams map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&bodyParams); err == nil {
			for k, v := range bodyParams {
				params[k] = v
			}
		}
	}

	return params
}

// executePipeline executes the query pipeline with conditional and loop support
func (h *Handler) executePipeline(db *sqlx.DB, baseParams map[string]interface{}) (interface{}, error) {
	// Build scope context with features and base parameters
	scope := make(map[string]interface{})
	for k, v := range baseParams {
		scope[k] = v
	}

	// Add features to scope for condition evaluation
	if h.Features != nil {
		scope["features"] = h.Features
	}

	// Separate variable for the result (not merged into scope)
	var result interface{}

	// If using old single query format, convert to pipeline
	queries := h.Queries
	if len(queries) == 0 && h.Query != "" {
		queries = []*config.QueryDef{
			{Query: h.Query},
		}
	}

	// Execute each query in the pipeline
	for _, qdef := range queries {
		// Check condition if specified
		if qdef.If != "" {
			condResult, err := h.evaluateCondition(qdef.If, scope)
			if err != nil {
				return nil, fmt.Errorf("condition evaluation failed: %w", err)
			}
			if !condResult {
				continue
			}
		}

		// Handle loop-based execution
		if qdef.For != "" {
			if err := h.executeLoop(db, qdef, scope); err != nil {
				return nil, err
			}
			continue
		}

		// Regular query execution
		queryResult, err := h.executeQuery(db, qdef.Query, scope)
		if err != nil {
			return nil, err
		}

		// Store result
		if qdef.As != "" {
			// Store at specified path in scope
			if err := h.setAtPath(scope, qdef.As, queryResult); err != nil {
				return nil, err
			}
		} else {
			// Store as root result (not merged into scope)
			result = queryResult
		}
	}

	// Return appropriate response format
	// If no explicit "as" was specified, return just the result
	if result != nil {
		return result, nil
	}

	// If result is nil but we have other data in scope (excluding features and base params),
	// return the scope as-is for complex pipelines
	contextKeys := make(map[string]bool)
	for k := range baseParams {
		contextKeys[k] = true
	}
	contextKeys["features"] = true

	resultData := make(map[string]interface{})
	for k, v := range scope {
		if !contextKeys[k] {
			resultData[k] = v
		}
	}

	if len(resultData) > 0 {
		return resultData, nil
	}

	return nil, nil
}

// evaluateCondition evaluates an expression against the scope
func (h *Handler) evaluateCondition(condition string, scope map[string]interface{}) (bool, error) {
	program, err := expr.Compile(condition)
	if err != nil {
		return false, fmt.Errorf("condition compile error: %w", err)
	}

	result, err := expr.Run(program, scope)
	if err != nil {
		return false, fmt.Errorf("condition evaluation error: %w", err)
	}

	// Convert result to bool
	switch v := result.(type) {
	case bool:
		return v, nil
	case nil:
		return false, nil
	default:
		return true, nil
	}
}

// executeQuery executes a single query and returns results
func (h *Handler) executeQuery(db *sqlx.DB, query string, params map[string]interface{}) (interface{}, error) {
	// Check if it's a write operation
	needsTransaction := h.shouldUseTransaction(query)

	if needsTransaction && h.Transaction != nil && h.Transaction.Enabled {
		return h.executeWithTransaction(db, query, params)
	}

	return h.executeQueryDirect(db, query, params)
}

// executeQueryDirect executes a query directly without transaction
func (h *Handler) executeQueryDirect(db *sqlx.DB, query string, params map[string]interface{}) (interface{}, error) {
	rows, err := db.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the rows
	results := []map[string]string{}
	for rows.Next() {
		row := map[string]interface{}{}
		if err := rows.MapScan(row); err != nil {
			return nil, err
		}

		result := make(map[string]string, len(row))
		for k, v := range row {
			result[strings.ToLower(k)] = dbValue(v)
		}

		results = append(results, result)
	}

	if len(results) == 0 {
		return nil, nil
	}
	if h.Single || len(results) == 1 {
		return results[0], nil
	}
	return results, nil
}

// executeWithTransaction executes a query with transaction and retries
func (h *Handler) executeWithTransaction(db *sqlx.DB, query string, params map[string]interface{}) (interface{}, error) {
	retries := 1
	if h.Transaction != nil && h.Transaction.Retries > 0 {
		retries = h.Transaction.Retries + 1
	}

	var lastErr error
	for attempt := 0; attempt < retries; attempt++ {
		if attempt > 0 {
			delayMs := 100
			if h.Transaction != nil && h.Transaction.RetryDelayMs > 0 {
				delayMs = h.Transaction.RetryDelayMs
			}
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
			log.Printf("Retrying transaction (attempt %d/%d)\n", attempt+1, retries)
		}

		tx, err := db.Beginx()
		if err != nil {
			lastErr = err
			continue
		}

		result, err := h.executeQueryTx(tx, query, params)
		if err != nil {
			tx.Rollback()
			lastErr = err
			continue
		}

		if err := tx.Commit(); err != nil {
			lastErr = err
			continue
		}

		return result, nil
	}

	return nil, fmt.Errorf("transaction failed after %d attempts: %w", retries, lastErr)
}

// executeQueryTx executes a query within a transaction
func (h *Handler) executeQueryTx(tx *sqlx.Tx, query string, params map[string]interface{}) (interface{}, error) {
	rows, err := tx.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]string{}
	for rows.Next() {
		row := map[string]interface{}{}
		if err := rows.MapScan(row); err != nil {
			return nil, err
		}

		result := make(map[string]string, len(row))
		for k, v := range row {
			result[strings.ToLower(k)] = dbValue(v)
		}

		results = append(results, result)
	}

	if len(results) == 0 {
		return nil, nil
	}
	if h.Single || len(results) == 1 {
		return results[0], nil
	}
	return results, nil
}

// executeLoop executes a query for each item in an array
func (h *Handler) executeLoop(db *sqlx.DB, qdef *config.QueryDef, scope map[string]interface{}) error {
	// Parse loop expression: (idx, item) in items
	loopRe := regexp.MustCompile(`\(\s*(\w+)\s*,\s*(\w+)\s*\)\s+in\s+(\w+)`)
	matches := loopRe.FindStringSubmatch(qdef.For)
	if len(matches) != 4 {
		return fmt.Errorf("invalid for expression: %s", qdef.For)
	}

	idxVar := matches[1]
	itemVar := matches[2]
	arrayName := matches[3]

	// Get the array from scope
	arrayVal, ok := scope[arrayName]
	if !ok {
		return fmt.Errorf("array %s not found in scope", arrayName)
	}

	// Convert to array of maps
	var items []interface{}
	switch v := arrayVal.(type) {
	case []interface{}:
		items = v
	case []map[string]string:
		for _, m := range v {
			items = append(items, m)
		}
	default:
		return fmt.Errorf("cannot iterate over %T", arrayVal)
	}

	// Execute query for each item
	for idx, item := range items {
		// Build item scope with loop variables
		itemScope := make(map[string]interface{})
		itemScope[idxVar] = idx
		itemScope[itemVar] = item

		// Execute query with item in scope
		result, err := h.executeQuery(db, qdef.Query, itemScope)
		if err != nil {
			return err
		}

		// Store result at specified path
		if qdef.As != "" {
			if err := h.setAtPath(scope, qdef.As, result); err != nil {
				return err
			}
		}
	}

	return nil
}

// setAtPath stores a value at a path in the scope
func (h *Handler) setAtPath(scope map[string]interface{}, path string, value interface{}) error {
	// Handle array indexing: "items[idx].user"
	if strings.Contains(path, "[") {
		return h.setAtPathWithIndex(scope, path, value)
	}

	// Handle nested paths: "user.profile.name"
	parts := strings.Split(path, ".")
	current := scope

	for i, part := range parts[:len(parts)-1] {
		if current[part] == nil {
			current[part] = make(map[string]interface{})
		}

		next, ok := current[part].(map[string]interface{})
		if !ok {
			return fmt.Errorf("cannot navigate path at %s", strings.Join(parts[:i+1], "."))
		}

		current = next
	}

	current[parts[len(parts)-1]] = value
	return nil
}

// setAtPathWithIndex handles array indexing in paths
func (h *Handler) setAtPathWithIndex(scope map[string]interface{}, path string, value interface{}) error {
	// Parse "items[idx].user" -> items, idx, user
	re := regexp.MustCompile(`(\w+)\[(\w+)\](.*)`)
	matches := re.FindStringSubmatch(path)
	if len(matches) != 4 {
		return fmt.Errorf("invalid path with index: %s", path)
	}

	arrayName := matches[1]
	indexVar := matches[2]
	restPath := strings.TrimPrefix(matches[3], ".")

	// Get array from scope
	arrayVal, ok := scope[arrayName]
	if !ok {
		return fmt.Errorf("array %s not found in scope", arrayName)
	}

	// Get index from scope
	idxVal, ok := scope[indexVar]
	if !ok {
		return fmt.Errorf("index variable %s not found in scope", indexVar)
	}

	idx, ok := idxVal.(int)
	if !ok {
		return fmt.Errorf("index %s is not an integer", indexVar)
	}

	// Get the item array
	items, ok := arrayVal.([]interface{})
	if !ok {
		return fmt.Errorf("cannot index non-array type %T", arrayVal)
	}

	if idx < 0 || idx >= len(items) {
		return fmt.Errorf("index %d out of bounds", idx)
	}

	// Set value in the item
	item, ok := items[idx].(map[string]interface{})
	if !ok {
		// Convert to map if needed
		item = make(map[string]interface{})
		items[idx] = item
	}

	if restPath == "" {
		item[arrayName] = value
		return nil
	}

	return h.setAtPath(item, restPath, value)
}

// shouldUseTransaction determines if a query is a write operation
func (h *Handler) shouldUseTransaction(query string) bool {
	query = strings.ToUpper(strings.TrimSpace(query))
	return strings.HasPrefix(query, "INSERT") ||
		strings.HasPrefix(query, "UPDATE") ||
		strings.HasPrefix(query, "DELETE")
}

// buildCacheKey creates a cache key from the request
func (h *Handler) buildCacheKey(r *http.Request) string {
	pattern := "default"
	if h.Cache != nil && h.Cache.KeyPattern != "" {
		pattern = h.Cache.KeyPattern
	}

	keyData := fmt.Sprintf("%s:%s:%s", pattern, r.Method, r.URL.String())
	hash := md5.Sum([]byte(keyData))
	return fmt.Sprintf("cache:%x", hash)
}

// getFromCache retrieves a cached entry if it hasn't expired
func (h *Handler) getFromCache(key string) (interface{}, bool) {
	h.cacheMutex.RLock()
	defer h.cacheMutex.RUnlock()

	entry, exists := h.cacheStore[key]
	if !exists || time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.data, true
}

// setInCache stores a response in the cache
func (h *Handler) setInCache(key string, data interface{}) {
	h.cacheMutex.Lock()
	defer h.cacheMutex.Unlock()

	ttl := 5 * time.Minute
	if h.Cache != nil && h.Cache.Expire != "" {
		if duration, err := time.ParseDuration(h.Cache.Expire); err == nil {
			ttl = duration
		}
	}

	h.cacheStore[key] = &cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(ttl),
	}
}

// Type returns the handler type
func (h *Handler) Type() string {
	return "sql"
}

// Handler creates a new handler instance from endpoint config
func (h *Handler) Handler(conf *config.Config, endpoint *config.Endpoint) (http.Handler, error) {
	handle := NewHandler()
	err := endpoint.Handler.Decode(handle)
	if err != nil {
		return nil, fmt.Errorf("error decoding config: %w", err)
	}

	// Copy global storage config if endpoint has it unset
	if handle.Storage == nil {
		handle.Storage = conf.Storage
	}

	// Copy handler configuration
	handle.Transaction = endpoint.Handler.Transaction
	handle.Cache = endpoint.Handler.Cache
	handle.RateLimit = endpoint.Handler.RateLimit
	handle.Response = endpoint.Handler.Response
	handle.Query = endpoint.Handler.Query
	handle.Queries = endpoint.Handler.Queries
	handle.Single = endpoint.Handler.Single
	handle.Parameters = endpoint.Handler.Parameters

	// Copy features from server config
	if conf.Server.Features != nil {
		handle.Features = conf.Server.Features
	}

	// Initialize cache store
	handle.cacheStore = make(map[string]*cacheEntry)

	// Initialize rate limiter if enabled
	if handle.RateLimit != nil && handle.RateLimit.Enabled && handle.RateLimit.Rate > 0 {
		// Parse the "per" duration (default to "1m" if unset)
		perDuration := "1m"
		if handle.RateLimit.Per != "" {
			perDuration = handle.RateLimit.Per
		}

		duration, err := time.ParseDuration(perDuration)
		if err != nil {
			log.Printf("Invalid duration for rate limit 'per': %v, using default 1m", err)
			duration = time.Minute
		}

		// Convert requests per duration to requests per second
		// e.g., 100 requests per 1m = 1.67 req/s
		rps := float64(handle.RateLimit.Rate) / duration.Seconds()
		burstSize := handle.RateLimit.Rate

		handle.limiter = rate.NewLimiter(rate.Limit(rps), burstSize)
		// log.Printf("Rate limiting enabled: %d req/%s (%.2f req/s) with burst size %d", handle.RateLimit.Rate, perDuration, rps, burstSize)
	}

	// Log cache configuration if enabled
	if handle.Cache != nil && handle.Cache.Enabled {
		expire := handle.Cache.Expire
		if expire == "" {
			expire = "5m"
		}
		// log.Printf("Cache enabled: TTL %s", expire)
	}

	// Log transaction configuration if enabled
	if handle.Transaction != nil && handle.Transaction.Enabled {
		//retries := handle.Transaction.Retries
		delayMs := handle.Transaction.RetryDelayMs
		if delayMs <= 0 {
			delayMs = 100
		}
		// log.Printf("Transaction support enabled: %d retries with %dms delay", retries, delayMs)
	}

	return handle, nil
}

// setResponseHeaders sets custom response headers
func (h *Handler) setResponseHeaders(w http.ResponseWriter) {
	if h.Response == nil {
		// Default JSON response
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Set custom headers
	if h.Response.Headers != nil {
		for key, value := range h.Response.Headers {
			w.Header().Set(key, value)
		}
	}

	// Set default Content-Type if not provided
	if _, hasContentType := h.Response.Headers["Content-Type"]; !hasContentType {
		if h.Response.Template != "" {
			// Default for template responses
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		} else {
			// Default for JSON responses
			w.Header().Set("Content-Type", "application/json")
		}
	}
}

// setRateLimitHeaders adds rate limit information to response headers
func (h *Handler) setRateLimitHeaders(w http.ResponseWriter) {
	if h.RateLimit == nil || !h.RateLimit.Enabled {
		return
	}

	if h.limiter == nil {
		return
	}

	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", h.RateLimit.Rate))
	perWindow := h.RateLimit.Per
	if perWindow == "" {
		perWindow = "1m"
	}
	w.Header().Set("X-RateLimit-Window", perWindow)
}

// renderTemplateResponse renders the result using VueGo template
func (h *Handler) renderTemplateResponse(w http.ResponseWriter, result interface{}) {
	if h.Response == nil || h.Response.Template == "" {
		// Fallback to JSON if no template
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Println("Error encoding json response:", err)
		}
		return
	}

	// Create data context
	data := map[string]interface{}{
		"data": result,
	}

	// Add individual fields to context if result is a map
	if resultMap, ok := result.(map[string]interface{}); ok {
		for k, v := range resultMap {
			data[k] = v
		}
	}

	// Load the template and render directly from the template string
	tpl := vuego.Load(newTemplateFS("template.vuego"))
	if err := tpl.Fill(data).RenderString(context.Background(), w, h.Response.Template); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Template rendering failed", http.StatusInternalServerError)
		return
	}
}

// templateFS is a simple in-memory filesystem for template content
type templateFS struct {
	content []byte
}

// newTemplateFS creates an in-memory filesystem with template content
func newTemplateFS(content string) fs.FS {
	return &templateFS{
		content: []byte(content),
	}
}

// Open implements fs.FS
func (t *templateFS) Open(name string) (fs.File, error) {
	if name != "template.vuego" {
		return nil, fmt.Errorf("file not found: %s", name)
	}
	return &templateFile{
		content: bytes.NewReader(t.content),
		name:    name,
	}, nil
}

// templateFile is an in-memory file
type templateFile struct {
	content *bytes.Reader
	name    string
}

// Stat implements fs.File
func (f *templateFile) Stat() (fs.FileInfo, error) {
	return &templateFileInfo{name: f.name, size: int64(f.content.Len())}, nil
}

// Read implements fs.File
func (f *templateFile) Read(p []byte) (int, error) {
	return f.content.Read(p)
}

// Close implements fs.File
func (f *templateFile) Close() error {
	return nil
}

// templateFileInfo provides file information
type templateFileInfo struct {
	name string
	size int64
}

// Name implements fs.FileInfo
func (i *templateFileInfo) Name() string {
	return i.name
}

// Size implements fs.FileInfo
func (i *templateFileInfo) Size() int64 {
	return i.size
}

// Mode implements fs.FileInfo
func (i *templateFileInfo) Mode() fs.FileMode {
	return 0444
}

// ModTime implements fs.FileInfo
func (i *templateFileInfo) ModTime() time.Time {
	return time.Now()
}

// IsDir implements fs.FileInfo
func (i *templateFileInfo) IsDir() bool {
	return false
}

// Sys implements fs.FileInfo
func (i *templateFileInfo) Sys() interface{} {
	return nil
}
