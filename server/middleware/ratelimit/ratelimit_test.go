package ratelimit_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/titpetric/etl/server/middleware/ratelimit"
)

func TestRateLimitMiddlewareBasic(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	keyBuilder := ratelimit.NewDefaultKeyBuilder()
	// 2 req/sec, burst of 3
	rateLimitMiddleware := ratelimit.NewMiddleware(store, keyBuilder, 2, 3)

	callCount := 0
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	handler := rateLimitMiddleware.Wrap(baseHandler)

	// Make 3 requests (within burst)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/api/data", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d", i+1, w.Code)
		}
		if callCount != i+1 {
			t.Errorf("Request %d: expected call count %d, got %d", i+1, i+1, callCount)
		}
	}

	// Fourth request should be rate limited
	req := httptest.NewRequest("GET", "/api/data", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w.Code)
	}
	if callCount != 3 {
		t.Errorf("Expected call count 3 (limited), got %d", callCount)
	}
}

func TestRateLimitMiddlewareDisabled(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	keyBuilder := ratelimit.NewDefaultKeyBuilder()
	rateLimitMiddleware := ratelimit.NewMiddleware(store, keyBuilder, 2, 3).
		WithEnabled(false)

	callCount := 0
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	handler := rateLimitMiddleware.Wrap(baseHandler)

	// Make requests beyond the limit
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/api/data", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: expected status 200 (disabled), got %d", i+1, w.Code)
		}
	}

	if callCount != 5 {
		t.Errorf("Expected all 5 calls when disabled, got %d", callCount)
	}
}

func TestRateLimitHeaders(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	keyBuilder := ratelimit.NewDefaultKeyBuilder()
	rateLimitMiddleware := ratelimit.NewMiddleware(store, keyBuilder, 10, 15)

	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rateLimitMiddleware.Wrap(baseHandler)

	req := httptest.NewRequest("GET", "/api/data", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Check headers
	if w.Header().Get("X-RateLimit-Limit") != "10" {
		t.Errorf("Expected X-RateLimit-Limit: 10, got %s", w.Header().Get("X-RateLimit-Limit"))
	}

	if remaining := w.Header().Get("X-RateLimit-Remaining"); remaining == "" {
		t.Error("Expected X-RateLimit-Remaining header")
	}

	if reset := w.Header().Get("X-RateLimit-Reset"); reset == "" {
		t.Error("Expected X-RateLimit-Reset header")
	}
}

func TestRateLimitPerIP(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	keyBuilder := ratelimit.NewDefaultKeyBuilder()
	// 2 req/sec, burst of 3
	rateLimitMiddleware := ratelimit.NewMiddleware(store, keyBuilder, 2, 3)

	callCount := 0
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	handler := rateLimitMiddleware.Wrap(baseHandler)

	// IP 1: make 3 requests (within burst)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/api/data", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("IP1 Request %d failed: status %d", i+1, w.Code)
		}
	}

	// IP 2: should have separate rate limit bucket
	req := httptest.NewRequest("GET", "/api/data", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("IP2 Request should succeed, got status %d", w.Code)
	}

	if callCount != 4 {
		t.Errorf("Expected 4 calls total, got %d", callCount)
	}
}

func TestRateLimitKeyBuilder(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		pattern    string
	}{
		{
			name:       "basic IP",
			remoteAddr: "127.0.0.1:12345",
			pattern:    "ratelimit",
		},
		{
			name:       "IPv6",
			remoteAddr: "[::1]:12345",
			pattern:    "ratelimit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyBuilder := ratelimit.NewDefaultKeyBuilder().WithPattern(tt.pattern)
			req := httptest.NewRequest("GET", "/api/data", nil)
			req.RemoteAddr = tt.remoteAddr
			key := keyBuilder.BuildKey(req)

			if len(key) == 0 {
				t.Error("Expected non-empty key")
			}
		})
	}
}

func TestRateLimitXForwardedFor(t *testing.T) {
	keyBuilder := ratelimit.NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "/api/data", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 198.51.100.1")
	req.RemoteAddr = "127.0.0.1:12345"

	key := keyBuilder.BuildKey(req)
	if len(key) == 0 {
		t.Error("Expected non-empty key with X-Forwarded-For")
	}
}

func TestMemoryStoreInc(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	ctx := context.Background()

	// Inc should return 1 on first call
	count, err := store.Inc(ctx, "test-key")
	if err != nil {
		t.Fatalf("Inc failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Inc should return 2 on second call
	count, err = store.Inc(ctx, "test-key")
	if err != nil {
		t.Fatalf("Inc failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestMemoryStoreRate(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	ctx := context.Background()

	// Rate on non-existent key should return 0
	count, err := store.Rate(ctx, "non-existent")
	if err != nil {
		t.Fatalf("Rate failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0 for non-existent, got %d", count)
	}

	// Inc then Rate should return same count
	store.Inc(ctx, "test-key")
	store.Inc(ctx, "test-key")
	count, err = store.Rate(ctx, "test-key")
	if err != nil {
		t.Fatalf("Rate failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestMemoryStoreReset(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	ctx := context.Background()

	// Inc then reset
	store.Inc(ctx, "test-key")
	store.Inc(ctx, "test-key")

	err := store.Reset(ctx, "test-key")
	if err != nil {
		t.Fatalf("Reset failed: %v", err)
	}

	// Rate should return 0 after reset
	count, err := store.Rate(ctx, "test-key")
	if err != nil {
		t.Fatalf("Rate failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0 after reset, got %d", count)
	}
}

func TestMemoryStoreClear(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	ctx := context.Background()

	// Inc multiple keys
	for i := 0; i < 5; i++ {
		store.Inc(ctx, fmt.Sprintf("key-%d", i))
	}

	// Clear
	err := store.Clear(ctx)
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// All should be 0 now
	for i := 0; i < 5; i++ {
		count, err := store.Rate(ctx, fmt.Sprintf("key-%d", i))
		if err != nil {
			t.Fatalf("Rate failed: %v", err)
		}
		if count != 0 {
			t.Errorf("Expected count 0 after clear, got %d for key-%d", count, i)
		}
	}
}

func TestMemoryStoreCleanupExpired(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	ctx := context.Background()

	// Inc a key
	store.Inc(ctx, "test-key")

	// Wait for reset time to pass (default 1 second)
	time.Sleep(1100 * time.Millisecond)

	// Cleanup
	err := store.CleanupExpired(ctx)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// The entry should be cleaned up (count resets)
	count, err := store.Rate(ctx, "test-key")
	if err != nil {
		t.Fatalf("Rate failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0 after cleanup of expired, got %d", count)
	}
}

func TestCustomRateLimitKeyBuilder(t *testing.T) {
	customKeyBuilder := ratelimit.NewCustomKeyBuilder(func(r *http.Request) string {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			apiKey = "anonymous"
		}
		return fmt.Sprintf("api:%s", apiKey)
	})

	req := httptest.NewRequest("GET", "/api/data", nil)
	req.Header.Set("X-API-Key", "key-123")
	key := customKeyBuilder.BuildKey(req)

	if key != "api:key-123" {
		t.Errorf("Expected 'api:key-123', got %s", key)
	}

	// Without header
	req = httptest.NewRequest("GET", "/api/data", nil)
	key = customKeyBuilder.BuildKey(req)

	if key != "api:anonymous" {
		t.Errorf("Expected 'api:anonymous', got %s", key)
	}
}

func TestRateLimitContextCancellation(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := store.Inc(ctx, "test-key")
	if err == nil {
		t.Error("Expected error for cancelled context")
	}
}
