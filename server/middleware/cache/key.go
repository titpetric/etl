package cache

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// KeyBuilder defines the interface for cache key generation
type KeyBuilder interface {
	// BuildKey creates a cache key from the request
	BuildKey(r *http.Request) string
}

// DefaultKeyBuilder creates cache keys using request method, URL, and headers
type DefaultKeyBuilder struct {
	// Pattern allows customization of key generation
	Pattern string
	// IncludeHeaders specifies which headers to include in the key
	IncludeHeaders []string
	// IncludeQuery specifies which query parameters to include in the key
	IncludeQuery []string
}

// NewDefaultKeyBuilder creates a new default key builder
func NewDefaultKeyBuilder() *DefaultKeyBuilder {
	return &DefaultKeyBuilder{
		Pattern: "default",
	}
}

// WithPattern sets the key pattern
func (kb *DefaultKeyBuilder) WithPattern(pattern string) *DefaultKeyBuilder {
	kb.Pattern = pattern
	return kb
}

// WithHeaders sets which headers to include in the key
func (kb *DefaultKeyBuilder) WithHeaders(headers ...string) *DefaultKeyBuilder {
	kb.IncludeHeaders = headers
	return kb
}

// WithQuery sets which query parameters to include in the key
func (kb *DefaultKeyBuilder) WithQuery(params ...string) *DefaultKeyBuilder {
	kb.IncludeQuery = params
	return kb
}

// BuildKey creates a cache key from the request
func (kb *DefaultKeyBuilder) BuildKey(r *http.Request) string {
	var parts []string

	// Add pattern
	parts = append(parts, kb.Pattern)

	// Add method and path
	parts = append(parts, r.Method)
	parts = append(parts, r.RequestURI)

	// Add specified headers if any
	if len(kb.IncludeHeaders) > 0 {
		headerParts := kb.extractHeaders(r)
		parts = append(parts, headerParts...)
	}

	// Add specified query parameters if any
	if len(kb.IncludeQuery) > 0 {
		queryParts := kb.extractQuery(r)
		parts = append(parts, queryParts...)
	}

	// Create MD5 hash of the combined parts
	keyData := strings.Join(parts, ":")
	hash := md5.Sum([]byte(keyData))
	return fmt.Sprintf("cache:%x", hash)
}

// extractHeaders extracts specified headers from the request
func (kb *DefaultKeyBuilder) extractHeaders(r *http.Request) []string {
	var parts []string
	for _, headerName := range kb.IncludeHeaders {
		value := r.Header.Get(headerName)
		if value != "" {
			parts = append(parts, fmt.Sprintf("h:%s=%s", headerName, value))
		}
	}
	return parts
}

// extractQuery extracts specified query parameters from the request
func (kb *DefaultKeyBuilder) extractQuery(r *http.Request) []string {
	var parts []string
	query := r.URL.Query()

	// Sort params for consistent key generation
	var keys []string
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		// Check if this parameter should be included
		shouldInclude := false
		for _, includeKey := range kb.IncludeQuery {
			if key == includeKey {
				shouldInclude = true
				break
			}
		}

		if shouldInclude {
			values := query[key]
			for _, value := range values {
				parts = append(parts, fmt.Sprintf("q:%s=%s", key, value))
			}
		}
	}

	return parts
}

// CustomKeyBuilder allows for custom cache key generation via a function
type CustomKeyBuilder struct {
	BuildFunc func(r *http.Request) string
}

// BuildKey uses the custom function to build the key
func (ckb *CustomKeyBuilder) BuildKey(r *http.Request) string {
	return ckb.BuildFunc(r)
}

// NewCustomKeyBuilder creates a new custom key builder
func NewCustomKeyBuilder(fn func(r *http.Request) string) *CustomKeyBuilder {
	return &CustomKeyBuilder{BuildFunc: fn}
}
