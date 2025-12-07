package model

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/etl/server/config"
)

// MockHandler is a test implementation of the Handler interface.
type MockHandler struct {
	TypeValue string
}

// Type returns the handler type.
func (m *MockHandler) Type() string {
	return m.TypeValue
}

// Handler returns a simple HTTP handler.
func (m *MockHandler) Handler(conf *config.Config, endpoint *config.Endpoint) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), nil
}

// TestRegisterHandler verifies that Register adds a handler to the registry.
func TestRegisterHandler(t *testing.T) {
	// Start with a fresh registry
	registeredHandlers = make(map[string]Handler)

	mock := &MockHandler{TypeValue: "test"}
	Register(mock)

	require.Contains(t, registeredHandlers, "test")
	require.Equal(t, mock, registeredHandlers["test"])
}

// TestRegisterMultipleHandlers verifies that multiple handlers can be registered.
func TestRegisterMultipleHandlers(t *testing.T) {
	registeredHandlers = make(map[string]Handler)

	mock1 := &MockHandler{TypeValue: "type1"}
	mock2 := &MockHandler{TypeValue: "type2"}

	Register(mock1)
	Register(mock2)

	require.Len(t, registeredHandlers, 2)
	require.Equal(t, mock1, registeredHandlers["type1"])
	require.Equal(t, mock2, registeredHandlers["type2"])
}

// TestRegisterOverwrite verifies that registering the same type overwrites the previous handler.
func TestRegisterOverwrite(t *testing.T) {
	registeredHandlers = make(map[string]Handler)

	mock1 := &MockHandler{TypeValue: "sameType"}
	mock2 := &MockHandler{TypeValue: "sameType"}

	Register(mock1)
	Register(mock2)

	require.Len(t, registeredHandlers, 1)
	require.Equal(t, mock2, registeredHandlers["sameType"])
}

// TestHandlersReturnsAllHandlers verifies that Handlers returns all registered handlers.
func TestHandlersReturnsAllHandlers(t *testing.T) {
	registeredHandlers = make(map[string]Handler)

	mock1 := &MockHandler{TypeValue: "handler1"}
	mock2 := &MockHandler{TypeValue: "handler2"}
	mock3 := &MockHandler{TypeValue: "handler3"}

	Register(mock1)
	Register(mock2)
	Register(mock3)

	handlers := Handlers()

	require.Len(t, handlers, 3)
	require.Contains(t, handlers, "handler1")
	require.Contains(t, handlers, "handler2")
	require.Contains(t, handlers, "handler3")
}

// TestHandlersReturnsEmptyMap verifies that Handlers returns an empty map when no handlers are registered.
func TestHandlersReturnsEmptyMap(t *testing.T) {
	registeredHandlers = make(map[string]Handler)

	handlers := Handlers()

	require.Len(t, handlers, 0)
	require.NotNil(t, handlers)
}

// TestHandlersReturnsMap verifies that Handlers returns the actual registry map.
func TestHandlersReturnsMap(t *testing.T) {
	registeredHandlers = make(map[string]Handler)

	mock := &MockHandler{TypeValue: "test"}
	Register(mock)

	handlers := Handlers()

	require.NotNil(t, handlers)
	require.IsType(t, map[string]Handler{}, handlers)
	require.Equal(t, mock, handlers["test"])
}

// TestRegisterWithEmptyType verifies that handlers with empty type strings can be registered.
func TestRegisterWithEmptyType(t *testing.T) {
	registeredHandlers = make(map[string]Handler)

	mock := &MockHandler{TypeValue: ""}
	Register(mock)

	handlers := Handlers()
	require.Contains(t, handlers, "")
	require.Equal(t, mock, handlers[""])
}

// TestHandlersNilCheck verifies that Handlers never returns nil.
func TestHandlersNilCheck(t *testing.T) {
	registeredHandlers = make(map[string]Handler)

	handlers := Handlers()
	require.NotNil(t, handlers)
}
