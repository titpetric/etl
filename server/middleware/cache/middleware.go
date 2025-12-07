package cache

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

// Middleware provides caching capabilities for HTTP handlers
type Middleware struct {
	// Store is the cache storage backend
	Store Store

	// KeyBuilder generates cache keys from requests
	KeyBuilder KeyBuilder

	// TTL is the cache time-to-live duration
	TTL time.Duration

	// Enabled indicates whether caching is active
	Enabled bool

	// Logger is an optional logger for cache operations
	Logger *log.Logger
}

// NewMiddleware creates a new cache middleware
func NewMiddleware(store Store, keyBuilder KeyBuilder) *Middleware {
	return &Middleware{
		Store:      store,
		KeyBuilder: keyBuilder,
		TTL:        5 * time.Minute,
		Enabled:    true,
	}
}

// WithTTL sets the cache TTL
func (m *Middleware) WithTTL(ttl time.Duration) *Middleware {
	m.TTL = ttl
	return m
}

// WithLogger sets the logger
func (m *Middleware) WithLogger(logger *log.Logger) *Middleware {
	m.Logger = logger
	return m
}

// WithEnabled sets whether caching is enabled
func (m *Middleware) WithEnabled(enabled bool) *Middleware {
	m.Enabled = enabled
	return m
}

// Wrap returns a middleware function that wraps an http.Handler
// It attempts to serve cached responses and caches successful responses
func (m *Middleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		// Only cache GET requests
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			next.ServeHTTP(w, r)
			return
		}

		// Generate cache key
		key := m.KeyBuilder.BuildKey(r)

		// Try to get from cache
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		cached, err := m.Store.Get(ctx, key)
		if err == nil && cached != nil {
			if m.Logger != nil {
				m.Logger.Printf("cache hit: %s", key)
			}

			// Serve from cache
			w.Header().Set("X-Cache", "HIT")
			for headerName, headerValues := range cached.Headers {
				for _, value := range headerValues {
					w.Header().Add(headerName, value)
				}
			}
			w.WriteHeader(cached.StatusCode)
			w.Write(cached.Body)
			return
		}

		// Request was not cached, capture the response using ResponseRecorder
		recorder := httptest.NewRecorder()

		// Serve the request
		next.ServeHTTP(recorder, r)

		// Set cache miss header first
		w.Header().Set("X-Cache", "MISS")

		// Copy response headers from recorder to actual writer
		for key, values := range recorder.Header() {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Write status code to actual response writer
		w.WriteHeader(recorder.Code)

		// Cache the response if it was successful
		if recorder.Code >= 200 && recorder.Code < 300 {
			entry := &Entry{
				StatusCode: recorder.Code,
				Headers:    recorder.Header(),
				Body:       recorder.Body.Bytes(),
			}

			ctx, cancel = context.WithTimeout(r.Context(), 5*time.Second)
			err = m.Store.Set(ctx, key, entry, m.TTL)
			cancel()

			if err != nil && m.Logger != nil {
				m.Logger.Printf("cache store error: %v", err)
			} else if m.Logger != nil {
				m.Logger.Printf("cache miss (stored): %s", key)
			}
		}

		// Write response body to actual response writer
		w.Write(recorder.Body.Bytes())
	})
}

// CleanupExpired removes expired cache entries
// This should be called periodically
func (m *Middleware) CleanupExpired(ctx context.Context) error {
	if cleaner, ok := m.Store.(*MemoryStore); ok {
		return cleaner.CleanupExpired(ctx)
	}
	return nil
}
