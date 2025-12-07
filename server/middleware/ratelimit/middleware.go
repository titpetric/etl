package ratelimit

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// Middleware provides rate limiting capabilities for HTTP handlers
type Middleware struct {
	// Store is the rate limit storage backend
	Store Store

	// KeyBuilder generates rate limit keys from requests
	KeyBuilder KeyBuilder

	// RequestsPerSecond is the maximum number of requests allowed per second
	RequestsPerSecond int64

	// BurstSize is the maximum number of requests allowed in a burst
	BurstSize int64

	// Enabled indicates whether rate limiting is active
	Enabled bool

	// Logger is an optional logger for rate limit operations
	Logger *log.Logger
}

// NewMiddleware creates a new rate limit middleware
func NewMiddleware(store Store, keyBuilder KeyBuilder, requestsPerSecond, burstSize int64) *Middleware {
	return &Middleware{
		Store:             store,
		KeyBuilder:        keyBuilder,
		RequestsPerSecond: requestsPerSecond,
		BurstSize:         burstSize,
		Enabled:           true,
	}
}

// WithLogger sets the logger
func (m *Middleware) WithLogger(logger *log.Logger) *Middleware {
	m.Logger = logger
	return m
}

// WithEnabled sets whether rate limiting is enabled
func (m *Middleware) WithEnabled(enabled bool) *Middleware {
	m.Enabled = enabled
	return m
}

// Wrap returns a middleware function that wraps an http.Handler
// It applies rate limiting before passing the request to the next handler
func (m *Middleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		// Generate rate limit key
		key := m.KeyBuilder.BuildKey(r)

		// Increment the counter
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		count, err := m.Store.Inc(ctx, key)
		cancel()

		if err != nil {
			if m.Logger != nil {
				m.Logger.Printf("rate limit store error: %v", err)
			}
			// On error, allow the request to proceed
			next.ServeHTTP(w, r)
			return
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", m.RequestsPerSecond))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, m.RequestsPerSecond-count)))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(1*time.Second).Unix()))

		// Check if we've exceeded the limit
		if count > m.BurstSize {
			if m.Logger != nil {
				m.Logger.Printf("rate limit exceeded: %s (count=%d, limit=%d)", key, count, m.BurstSize)
			}
			w.Header().Set("Retry-After", "1")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Request is allowed
		next.ServeHTTP(w, r)
	})
}

// Flush implements http.Flusher if the underlying writer supports it
func (rw *responseWrapper) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack implements http.Hijacker if the underlying writer supports it
func (rw *responseWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return h.Hijack()
}

// responseWrapper would only be needed if we need to wrap the response
type responseWrapper struct {
	http.ResponseWriter
}

// max returns the larger of two integers
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// CleanupExpired removes expired rate limit entries
// This should be called periodically
func (m *Middleware) CleanupExpired(ctx context.Context) error {
	if cleaner, ok := m.Store.(*MemoryStore); ok {
		return cleaner.CleanupExpired(ctx)
	}
	return nil
}
