package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDefaultKeyBuilderBuildKey verifies that BuildKey generates a valid rate limit key.
func TestDefaultKeyBuilderBuildKey(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "http://example.com/path", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "ratelimit:"))
	require.Equal(t, 42, len(key)) // "ratelimit:" (10 chars) + 32 char hex MD5
}

// TestDefaultKeyBuilderWithPattern verifies that WithPattern modifies the pattern.
func TestDefaultKeyBuilderWithPattern(t *testing.T) {
	kb := NewDefaultKeyBuilder().WithPattern("custom_pattern")
	req := httptest.NewRequest("GET", "http://example.com/path", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "ratelimit:"))
}

// TestDefaultKeyBuilderConsistency verifies that identical requests produce identical keys.
func TestDefaultKeyBuilderConsistency(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "http://example.com/path", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	key1 := kb.BuildKey(req)
	key2 := kb.BuildKey(req)

	require.Equal(t, key1, key2)
}

// TestDefaultKeyBuilderDifferentIPs verifies that different IPs produce different keys.
func TestDefaultKeyBuilderDifferentIPs(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req1 := httptest.NewRequest("GET", "http://example.com/path", nil)
	req1.RemoteAddr = "192.168.1.1:12345"

	req2 := httptest.NewRequest("GET", "http://example.com/path", nil)
	req2.RemoteAddr = "192.168.1.2:12345"

	key1 := kb.BuildKey(req1)
	key2 := kb.BuildKey(req2)

	require.NotEqual(t, key1, key2)
}

// TestDefaultKeyBuilderDifferentPaths verifies that different paths produce different keys.
func TestDefaultKeyBuilderDifferentPaths(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req1 := httptest.NewRequest("GET", "http://example.com/path1", nil)
	req1.RemoteAddr = "192.168.1.1:12345"

	req2 := httptest.NewRequest("GET", "http://example.com/path2", nil)
	req2.RemoteAddr = "192.168.1.1:12345"

	key1 := kb.BuildKey(req1)
	key2 := kb.BuildKey(req2)

	require.NotEqual(t, key1, key2)
}

// TestDefaultKeyBuilderGetClientIPFromRemoteAddr verifies that IP extraction from RemoteAddr works.
func TestDefaultKeyBuilderGetClientIPFromRemoteAddr(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "http://example.com/path", nil)
	req.RemoteAddr = "192.168.1.100:54321"

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "ratelimit:"))
}

// TestDefaultKeyBuilderGetClientIPFromXForwardedFor verifies that X-Forwarded-For header is used.
func TestDefaultKeyBuilderGetClientIPFromXForwardedFor(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "http://example.com/path", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 198.51.100.2")

	key1 := kb.BuildKey(req)

	req2 := httptest.NewRequest("GET", "http://example.com/path", nil)
	req2.RemoteAddr = "127.0.0.1:12345"
	req2.Header.Set("X-Forwarded-For", "203.0.113.2, 198.51.100.2")

	key2 := kb.BuildKey(req2)

	require.NotEqual(t, key1, key2)
}

// TestDefaultKeyBuilderGetClientIPFromXRealIP verifies that X-Real-IP header is used.
func TestDefaultKeyBuilderGetClientIPFromXRealIP(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "http://example.com/path", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("X-Real-IP", "203.0.113.100")

	key1 := kb.BuildKey(req)

	req2 := httptest.NewRequest("GET", "http://example.com/path", nil)
	req2.RemoteAddr = "127.0.0.1:12345"
	req2.Header.Set("X-Real-IP", "203.0.113.101")

	key2 := kb.BuildKey(req2)

	require.NotEqual(t, key1, key2)
}

// TestDefaultKeyBuilderWithHeaders verifies that headers can be included in the key.
func TestDefaultKeyBuilderWithHeaders(t *testing.T) {
	kb := NewDefaultKeyBuilder().WithHeaders("Authorization")
	req := httptest.NewRequest("GET", "http://example.com/path", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Header.Set("Authorization", "Bearer token123")

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "ratelimit:"))
}

// TestDefaultKeyBuilderWithHeadersDifferentValues verifies that different header values produce different keys.
func TestDefaultKeyBuilderWithHeadersDifferentValues(t *testing.T) {
	kb := NewDefaultKeyBuilder().WithHeaders("Authorization")

	req1 := httptest.NewRequest("GET", "http://example.com/path", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	req1.Header.Set("Authorization", "Bearer token1")

	req2 := httptest.NewRequest("GET", "http://example.com/path", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	req2.Header.Set("Authorization", "Bearer token2")

	key1 := kb.BuildKey(req1)
	key2 := kb.BuildKey(req2)

	require.NotEqual(t, key1, key2)
}

// TestDefaultKeyBuilderWithQuery verifies that query parameters can be included in the key.
func TestDefaultKeyBuilderWithQuery(t *testing.T) {
	kb := NewDefaultKeyBuilder().WithQuery("user_id")
	req := httptest.NewRequest("GET", "http://example.com/path?user_id=123&other=value", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "ratelimit:"))
}

// TestDefaultKeyBuilderWithQueryOrdering verifies that query parameter order doesn't affect the key when using WithQuery.
func TestDefaultKeyBuilderWithQueryOrdering(t *testing.T) {
	kb := NewDefaultKeyBuilder().WithQuery("a", "b")

	req1 := httptest.NewRequest("GET", "http://example.com/path?a=1&b=2", nil)
	req1.RemoteAddr = "192.168.1.1:12345"

	req2 := httptest.NewRequest("GET", "http://example.com/path?b=2&a=1", nil)
	req2.RemoteAddr = "192.168.1.1:12345"

	key1 := kb.BuildKey(req1)
	key2 := kb.BuildKey(req2)

	// Keys should be the same because path and IP are the same, query params are sorted consistently
	require.Equal(t, key1, key2)
}

// TestDefaultKeyBuilderIPWithoutPort verifies that IP without port is handled correctly.
func TestDefaultKeyBuilderIPWithoutPort(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "http://example.com/path", nil)
	req.RemoteAddr = "192.168.1.1"

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "ratelimit:"))
}

// TestDefaultKeyBuilderUnknownIP verifies that unknown IP is handled gracefully.
func TestDefaultKeyBuilderUnknownIP(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "http://example.com/path", nil)
	req.RemoteAddr = ""

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "ratelimit:"))
}

// TestCustomKeyBuilder verifies that CustomKeyBuilder uses the provided function.
func TestCustomKeyBuilder(t *testing.T) {
	ckb := NewCustomKeyBuilder(func(r *http.Request) string {
		return "custom_limit_" + r.Method
	})
	req := httptest.NewRequest("GET", "http://example.com/path", nil)

	key := ckb.BuildKey(req)

	require.Equal(t, "custom_limit_GET", key)
}

// TestDefaultKeyBuilderIPWithIPv6 verifies that IPv6 addresses are handled.
func TestDefaultKeyBuilderIPWithIPv6(t *testing.T) {
	kb := NewDefaultKeyBuilder()
	req := httptest.NewRequest("GET", "http://example.com/path", nil)
	req.RemoteAddr = "[::1]:8080"

	key := kb.BuildKey(req)

	require.NotEmpty(t, key)
	require.True(t, strings.HasPrefix(key, "ratelimit:"))
}

// TestDefaultKeyBuilderXForwardedForPriority verifies that X-Forwarded-For takes priority over RemoteAddr.
func TestDefaultKeyBuilderXForwardedForPriority(t *testing.T) {
	kb := NewDefaultKeyBuilder()

	req1 := httptest.NewRequest("GET", "http://example.com/path", nil)
	req1.RemoteAddr = "127.0.0.1:12345"
	req1.Header.Set("X-Forwarded-For", "203.0.113.1")

	req2 := httptest.NewRequest("GET", "http://example.com/path", nil)
	req2.RemoteAddr = "203.0.113.1:12345"

	key1 := kb.BuildKey(req1)
	key2 := kb.BuildKey(req2)

	require.Equal(t, key1, key2)
}
