package cache_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/titpetric/etl/server/middleware/cache"
)

func TestCacheMiddlewareHitMiss(t *testing.T) {
	store := cache.NewMemoryStore()
	keyBuilder := cache.NewDefaultKeyBuilder()
	cacheMiddleware := cache.NewMiddleware(store, keyBuilder).
		WithTTL(5 * time.Minute)

	callCount := 0
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"count":%d}`, callCount)
	})

	handler := cacheMiddleware.Wrap(baseHandler)

	// First request (cache miss)
	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("X-Cache") != "MISS" {
		t.Errorf("Expected X-Cache: MISS, got %s", w.Header().Get("X-Cache"))
	}
	if callCount != 1 {
		t.Errorf("Expected handler called once, got %d", callCount)
	}

	// Second request (cache hit)
	req = httptest.NewRequest("GET", "/api/data", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("X-Cache") != "HIT" {
		t.Errorf("Expected X-Cache: HIT, got %s", w.Header().Get("X-Cache"))
	}
	if callCount != 1 {
		t.Errorf("Expected handler called once (cached), got %d", callCount)
	}
}

func TestCacheMiddlewareDisabled(t *testing.T) {
	store := cache.NewMemoryStore()
	keyBuilder := cache.NewDefaultKeyBuilder()
	cacheMiddleware := cache.NewMiddleware(store, keyBuilder).
		WithEnabled(false)

	callCount := 0
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	handler := cacheMiddleware.Wrap(baseHandler)

	// First request
	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second request (should not be cached)
	req = httptest.NewRequest("GET", "/api/data", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestCacheMiddlewareOnlyGET(t *testing.T) {
	store := cache.NewMemoryStore()
	keyBuilder := cache.NewDefaultKeyBuilder()
	cacheMiddleware := cache.NewMiddleware(store, keyBuilder).
		WithTTL(5 * time.Minute)

	callCount := 0
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	handler := cacheMiddleware.Wrap(baseHandler)

	// POST request should not be cached
	req := httptest.NewRequest("POST", "/api/data", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if callCount != 1 {
		t.Errorf("Expected 1 call for POST, got %d", callCount)
	}

	// Second POST request should not hit cache
	req = httptest.NewRequest("POST", "/api/data", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if callCount != 2 {
		t.Errorf("Expected 2 calls for POST (no caching), got %d", callCount)
	}
}

func TestCacheMiddlewareExpiration(t *testing.T) {
	store := cache.NewMemoryStore()
	keyBuilder := cache.NewDefaultKeyBuilder()
	cacheMiddleware := cache.NewMiddleware(store, keyBuilder).
		WithTTL(100 * time.Millisecond)

	callCount := 0
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	handler := cacheMiddleware.Wrap(baseHandler)

	// First request (cache miss)
	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second request (cache hit)
	req = httptest.NewRequest("GET", "/api/data", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if callCount != 1 {
		t.Errorf("Expected 1 call (cached), got %d", callCount)
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Third request (cache expired, miss)
	req = httptest.NewRequest("GET", "/api/data", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if callCount != 2 {
		t.Errorf("Expected 2 calls (expired), got %d", callCount)
	}
}

func TestCacheKeyBuilder(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple GET",
			path:     "/api/users",
			expected: "cache:",
		},
		{
			name:     "with query params",
			path:     "/api/users?id=1",
			expected: "cache:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyBuilder := cache.NewDefaultKeyBuilder()
			req := httptest.NewRequest("GET", tt.path, nil)
			key := keyBuilder.BuildKey(req)
			if len(key) < len(tt.expected) {
				t.Errorf("Expected key to start with %s, got %s", tt.expected, key)
			}
		})
	}
}

func TestMemoryStoreGetSet(t *testing.T) {
	store := cache.NewMemoryStore()
	ctx := context.Background()

	// Test Set
	entry := &cache.Entry{
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: []byte(`{"test":"data"}`),
	}
	err := store.Set(ctx, "test-key", entry, 5*time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test Get
	retrieved, err := store.Get(ctx, "test-key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved == nil {
		t.Fatal("Expected to retrieve entry")
	}
	if string(retrieved.Body) != `{"test":"data"}` {
		t.Errorf("Expected body %q, got %q", `{"test":"data"}`, string(retrieved.Body))
	}

	// Test Get non-existent
	retrieved, err = store.Get(ctx, "non-existent")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved != nil {
		t.Error("Expected nil for non-existent key")
	}
}

func TestMemoryStoreDelete(t *testing.T) {
	store := cache.NewMemoryStore()
	ctx := context.Background()

	// Set an entry
	entry := &cache.Entry{Body: []byte("test")}
	store.Set(ctx, "test-key", entry, 5*time.Minute)

	// Delete it
	err := store.Delete(ctx, "test-key")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify it's gone
	retrieved, err := store.Get(ctx, "test-key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved != nil {
		t.Error("Expected nil after delete")
	}
}

func TestMemoryStoreClear(t *testing.T) {
	store := cache.NewMemoryStore()
	ctx := context.Background()

	// Set multiple entries
	for i := 0; i < 5; i++ {
		entry := &cache.Entry{Body: []byte(fmt.Sprintf("test-%d", i))}
		store.Set(ctx, fmt.Sprintf("key-%d", i), entry, 5*time.Minute)
	}

	// Clear
	err := store.Clear(ctx)
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify all are gone
	for i := 0; i < 5; i++ {
		retrieved, err := store.Get(ctx, fmt.Sprintf("key-%d", i))
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if retrieved != nil {
			t.Errorf("Expected nil after clear, got value for key-%d", i)
		}
	}
}

func TestMemoryStoreCleanupExpired(t *testing.T) {
	store := cache.NewMemoryStore()
	ctx := context.Background()

	// Set entry with short TTL
	entry := &cache.Entry{Body: []byte("test")}
	store.Set(ctx, "test-key", entry, 1*time.Millisecond)

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Cleanup
	err := store.CleanupExpired(ctx)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Verify expired entry was removed
	retrieved, err := store.Get(ctx, "test-key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved != nil {
		t.Error("Expected expired entry to be removed")
	}
}

func TestCustomKeyBuilder(t *testing.T) {
	customKeyBuilder := cache.NewCustomKeyBuilder(func(r *http.Request) string {
		return fmt.Sprintf("custom:%s", r.URL.Path)
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	key := customKeyBuilder.BuildKey(req)

	if key != "custom:/api/users" {
		t.Errorf("Expected 'custom:/api/users', got %s", key)
	}
}

func TestCacheMiddlewareResponseRecorder(t *testing.T) {
	store := cache.NewMemoryStore()
	keyBuilder := cache.NewDefaultKeyBuilder()
	cacheMiddleware := cache.NewMiddleware(store, keyBuilder).
		WithTTL(5 * time.Minute)

	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom", "value")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"status":"success"}`)
	})

	handler := cacheMiddleware.Wrap(baseHandler)

	// First request (cache miss)
	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Verify response from first request
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type: application/json, got %s", w.Header().Get("Content-Type"))
	}
	if w.Header().Get("X-Custom") != "value" {
		t.Errorf("Expected X-Custom: value, got %s", w.Header().Get("X-Custom"))
	}
	if w.Header().Get("X-Cache") != "MISS" {
		t.Errorf("Expected X-Cache: MISS, got %s", w.Header().Get("X-Cache"))
	}
	if !strings.Contains(w.Body.String(), "success") {
		t.Errorf("Expected body to contain 'success', got %s", w.Body.String())
	}

	// Second request (cache hit)
	req = httptest.NewRequest("GET", "/api/data", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Verify response from cache
	if w.Code != http.StatusCreated {
		t.Errorf("Expected cached status %d, got %d", http.StatusCreated, w.Code)
	}
	if w.Header().Get("X-Cache") != "HIT" {
		t.Errorf("Expected X-Cache: HIT, got %s", w.Header().Get("X-Cache"))
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected cached Content-Type: application/json, got %s", w.Header().Get("Content-Type"))
	}
}

func TestCacheMiddlewareNoCacheOnError(t *testing.T) {
	store := cache.NewMemoryStore()
	keyBuilder := cache.NewDefaultKeyBuilder()
	cacheMiddleware := cache.NewMiddleware(store, keyBuilder).
		WithTTL(5 * time.Minute)

	callCount := 0
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"server error"}`)
	})

	handler := cacheMiddleware.Wrap(baseHandler)

	// First request with error
	req := httptest.NewRequest("GET", "/api/error", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected error status, got %d", w.Code)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second request should not use cache
	req = httptest.NewRequest("GET", "/api/error", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if callCount != 2 {
		t.Errorf("Expected 2 calls (error not cached), got %d", callCount)
	}
}

func TestCacheMiddlewareHeadRequest(t *testing.T) {
	store := cache.NewMemoryStore()
	keyBuilder := cache.NewDefaultKeyBuilder()
	cacheMiddleware := cache.NewMiddleware(store, keyBuilder).
		WithTTL(5 * time.Minute)

	callCount := 0
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Length", "100")
		w.WriteHeader(http.StatusOK)
		if r.Method != http.MethodHead {
			fmt.Fprintf(w, "response body")
		}
	})

	handler := cacheMiddleware.Wrap(baseHandler)

	// HEAD request
	req := httptest.NewRequest("HEAD", "/api/resource", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected OK status, got %d", w.Code)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second HEAD request
	req = httptest.NewRequest("HEAD", "/api/resource", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if callCount != 1 {
		t.Errorf("Expected 1 call total (cached), got %d", callCount)
	}
}

func TestCacheMiddlewareCacheBoundaries(t *testing.T) {
	keyBuilder := cache.NewDefaultKeyBuilder()

	tests := []struct {
		name        string
		status      int
		shouldCache bool
	}{
		{"199 Not cached", 199, false},
		{"200 Cached", http.StatusOK, true},
		{"201 Cached", http.StatusCreated, true},
		{"299 Cached", 299, true},
		{"300 Not cached", http.StatusMultipleChoices, false},
		{"301 Not cached", http.StatusMovedPermanently, false},
		{"400 Not cached", http.StatusBadRequest, false},
		{"404 Not cached", http.StatusNotFound, false},
		{"500 Not cached", http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.NewMemoryStore()
			cacheMiddleware := cache.NewMiddleware(store, keyBuilder).
				WithTTL(5 * time.Minute)

			callCount := 0
			baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount++
				w.WriteHeader(tt.status)
			})

			handler := cacheMiddleware.Wrap(baseHandler)

			// First request
			req := httptest.NewRequest("GET", "/api/test", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Second request
			req = httptest.NewRequest("GET", "/api/test", nil)
			w = httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if tt.shouldCache && callCount != 1 {
				t.Errorf("Expected response to be cached, but handler was called %d times", callCount)
			} else if !tt.shouldCache && callCount != 2 {
				t.Errorf("Expected response not to be cached, but handler was called %d times", callCount)
			}
		})
	}
}

func TestCacheMiddlewareMultipleHeaders(t *testing.T) {
	store := cache.NewMemoryStore()
	keyBuilder := cache.NewDefaultKeyBuilder()
	cacheMiddleware := cache.NewMiddleware(store, keyBuilder).
		WithTTL(5 * time.Minute)

	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Set-Cookie", "cookie1=value1")
		w.Header().Add("Set-Cookie", "cookie2=value2")
		w.Header().Add("Accept-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
	})

	handler := cacheMiddleware.Wrap(baseHandler)

	// First request (cache miss)
	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if len(w.Header()["Set-Cookie"]) != 2 {
		t.Errorf("Expected 2 Set-Cookie headers, got %d", len(w.Header()["Set-Cookie"]))
	}

	// Second request (cache hit)
	req = httptest.NewRequest("GET", "/api/data", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if len(w.Header()["Set-Cookie"]) != 2 {
		t.Errorf("Expected cached 2 Set-Cookie headers, got %d", len(w.Header()["Set-Cookie"]))
	}
	if len(w.Header()["Accept-Encoding"]) != 1 {
		t.Errorf("Expected 1 Accept-Encoding header, got %d", len(w.Header()["Accept-Encoding"]))
	}
}

func TestCacheMiddlewareEmptyBody(t *testing.T) {
	store := cache.NewMemoryStore()
	keyBuilder := cache.NewDefaultKeyBuilder()
	cacheMiddleware := cache.NewMiddleware(store, keyBuilder).
		WithTTL(5 * time.Minute)

	callCount := 0
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("X-Processed", "true")
		w.WriteHeader(http.StatusNoContent)
	})

	handler := cacheMiddleware.Wrap(baseHandler)

	// First request
	req := httptest.NewRequest("GET", "/api/empty", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected NoContent status, got %d", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Errorf("Expected empty body, got %s", w.Body.String())
	}

	// Second request (should hit cache even with empty body)
	req = httptest.NewRequest("GET", "/api/empty", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if callCount != 1 {
		t.Errorf("Expected 1 call (cached), got %d", callCount)
	}
	if w.Header().Get("X-Processed") != "true" {
		t.Errorf("Expected cached header X-Processed: true")
	}
}
