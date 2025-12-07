package ratelimit

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// KeyBuilder defines the interface for rate limit key generation
type KeyBuilder interface {
	// BuildKey creates a rate limit key from the request
	BuildKey(r *http.Request) string
}

// DefaultKeyBuilder creates rate limit keys using IP address and request path
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
		Pattern: "ratelimit",
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

// BuildKey creates a rate limit key from the request
// By default, it uses the client IP address
func (kb *DefaultKeyBuilder) BuildKey(r *http.Request) string {
	var parts []string

	// Add pattern
	parts = append(parts, kb.Pattern)

	// Add client IP address (primary identifier)
	clientIP := kb.getClientIP(r)
	parts = append(parts, clientIP)

	// Add path
	parts = append(parts, r.URL.Path)

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

	// Create MD5 hash of the combined parts for consistency
	keyData := strings.Join(parts, ":")
	hash := md5.Sum([]byte(keyData))
	return fmt.Sprintf("ratelimit:%x", hash)
}

// getClientIP extracts the client IP address from the request
func (kb *DefaultKeyBuilder) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header (proxy)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Use RemoteAddr
	if r.RemoteAddr != "" {
		// Remove port if present
		if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
			return r.RemoteAddr[:idx]
		}
		return r.RemoteAddr
	}

	return "unknown"
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

// CustomKeyBuilder allows for custom rate limit key generation via a function
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
