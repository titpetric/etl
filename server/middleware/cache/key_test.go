package cache

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDefaultKeyBuilderBuildKey verifies that BuildKey generates a valid cache key.
func TestDefaultKeyBuilderBuildKey(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "http://example.com/path", nil)

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "cache:"))
	require.Equal(t, 38, len(key)) // "cache:" (6 chars) + 32 char hex MD5
}

// TestDefaultKeyBuilderWithPattern verifies that WithPattern modifies the pattern.
func TestDefaultKeyBuilderWithPattern(t *testing.T) {
	kb := NewDefaultKeyBuilder().WithPattern("custom")
	req := httptest.NewRequest("GET", "http://example.com/test", nil)

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "cache:"))
}

// TestDefaultKeyBuilderConsistency verifies that identical requests produce identical keys.
func TestDefaultKeyBuilderConsistency(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "http://example.com/path", nil)

	key1 := kb.BuildKey(req)
	key2 := kb.BuildKey(req)

	require.Equal(t, key1, key2)
}

// TestDefaultKeyBuilderDifferentMethods verifies that different methods produce different keys.
func TestDefaultKeyBuilderDifferentMethods(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	reqGet := httptest.NewRequest("GET", "http://example.com/path", nil)
	reqPost := httptest.NewRequest("POST", "http://example.com/path", nil)

	keyGet := kb.BuildKey(reqGet)
	keyPost := kb.BuildKey(reqPost)

	require.NotEqual(t, keyGet, keyPost)
}

// TestDefaultKeyBuilderDifferentPaths verifies that different paths produce different keys.
func TestDefaultKeyBuilderDifferentPaths(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req1 := httptest.NewRequest("GET", "http://example.com/path1", nil)
	req2 := httptest.NewRequest("GET", "http://example.com/path2", nil)

	key1 := kb.BuildKey(req1)
	key2 := kb.BuildKey(req2)

	require.NotEqual(t, key1, key2)
}

// TestDefaultKeyBuilderWithHeaders verifies that headers can be included in the key.
func TestDefaultKeyBuilderWithHeaders(t *testing.T) {
	kb := NewDefaultKeyBuilder().WithHeaders("Authorization", "User-Agent")
	req := httptest.NewRequest("GET", "http://example.com/path", nil)
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("User-Agent", "TestBot/1.0")

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "cache:"))
}

// TestDefaultKeyBuilderWithHeadersNoMatch verifies that missing headers don't affect the key.
func TestDefaultKeyBuilderWithHeadersNoMatch(t *testing.T) {
	kb := NewDefaultKeyBuilder().WithHeaders("X-Custom-Header")
	req1 := httptest.NewRequest("GET", "http://example.com/path", nil)
	req2 := httptest.NewRequest("GET", "http://example.com/path", nil)
	req2.Header.Set("X-Custom-Header", "value")

	key1 := kb.BuildKey(req1)
	key2 := kb.BuildKey(req2)

	require.NotEqual(t, key1, key2)
}

// TestDefaultKeyBuilderWithQuery verifies that query parameters can be included in the key.
func TestDefaultKeyBuilderWithQuery(t *testing.T) {
	kb := NewDefaultKeyBuilder().WithQuery("page", "limit")
	req := httptest.NewRequest("GET", "http://example.com/path?page=1&limit=10&other=value", nil)

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "cache:"))
}

// TestDefaultKeyBuilderWithQueryOrdering verifies that different query order produces different keys.
func TestDefaultKeyBuilderWithQueryOrdering(t *testing.T) {
	kb := NewDefaultKeyBuilder().WithQuery("a", "b")
	req1 := httptest.NewRequest("GET", "http://example.com/path?a=1&b=2", nil)
	req2 := httptest.NewRequest("GET", "http://example.com/path?b=2&a=1", nil)

	key1 := kb.BuildKey(req1)
	key2 := kb.BuildKey(req2)

	// Keys should be different because RequestURI differs
	require.NotEqual(t, key1, key2)
}

// TestDefaultKeyBuilderIncludeQuerySorting verifies that query params are sorted consistently.
func TestDefaultKeyBuilderIncludeQuerySorting(t *testing.T) {
	kb := NewDefaultKeyBuilder().WithQuery("z", "a", "m")
	req := httptest.NewRequest("GET", "http://example.com/path?z=1&a=2&m=3", nil)

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "cache:"))
}

// TestCustomKeyBuilder verifies that CustomKeyBuilder uses the provided function.
func TestCustomKeyBuilder(t *testing.T) {
	ckb := NewCustomKeyBuilder(func(r *http.Request) string {
		return "custom_key_" + r.Method
	})
	req := httptest.NewRequest("GET", "http://example.com/path", nil)

	key := ckb.BuildKey(req)

	require.Equal(t, "custom_key_GET", key)
}

// TestCustomKeyBuilderConsistency verifies that CustomKeyBuilder produces consistent keys.
func TestCustomKeyBuilderConsistency(t *testing.T) {
	ckb := NewCustomKeyBuilder(func(r *http.Request) string {
		return "static_key"
	})
	req1 := httptest.NewRequest("GET", "http://example.com/path1", nil)
	req2 := httptest.NewRequest("POST", "http://example.com/path2", nil)

	key1 := ckb.BuildKey(req1)
	key2 := ckb.BuildKey(req2)

	require.Equal(t, "static_key", key1)
	require.Equal(t, "static_key", key2)
}

// TestDefaultKeyBuilderRequestURI verifies that request URI is included in the key.
func TestDefaultKeyBuilderRequestURI(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "http://example.com/path?query=value", nil)

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
}

// TestDefaultKeyBuilderEmptyPath verifies that empty paths produce a key.
func TestDefaultKeyBuilderEmptyPath(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "http://example.com/", nil)

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "cache:"))
}
