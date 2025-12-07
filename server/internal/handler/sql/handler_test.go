package sql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestShouldUseTransactionInsert verifies that INSERT queries are identified as write operations.
func TestShouldUseTransactionInsert(t *testing.T) {
	h := NewHandler()
	require.True(t, h.shouldUseTransaction("INSERT INTO users (name) VALUES ('John')"))
}

// TestShouldUseTransactionUpdate verifies that UPDATE queries are identified as write operations.
func TestShouldUseTransactionUpdate(t *testing.T) {
	h := NewHandler()
	require.True(t, h.shouldUseTransaction("UPDATE users SET name = 'Jane' WHERE id = 1"))
}

// TestShouldUseTransactionDelete verifies that DELETE queries are identified as write operations.
func TestShouldUseTransactionDelete(t *testing.T) {
	h := NewHandler()
	require.True(t, h.shouldUseTransaction("DELETE FROM users WHERE id = 1"))
}

// TestShouldUseTransactionSelect verifies that SELECT queries are not write operations.
func TestShouldUseTransactionSelect(t *testing.T) {
	h := NewHandler()
	require.False(t, h.shouldUseTransaction("SELECT * FROM users"))
}

// TestShouldUseTransactionWithLeadingWhitespace verifies that leading whitespace is handled.
func TestShouldUseTransactionWithLeadingWhitespace(t *testing.T) {
	h := NewHandler()
	require.True(t, h.shouldUseTransaction("  INSERT INTO users (name) VALUES ('John')"))
}

// TestShouldUseTransactionWithTrailingWhitespace verifies that trailing whitespace is handled.
func TestShouldUseTransactionWithTrailingWhitespace(t *testing.T) {
	h := NewHandler()
	require.True(t, h.shouldUseTransaction("INSERT INTO users (name) VALUES ('John')  "))
}

// TestShouldUseTransactionLowercase verifies that lowercase queries are handled.
func TestShouldUseTransactionLowercase(t *testing.T) {
	h := NewHandler()
	require.True(t, h.shouldUseTransaction("insert into users (name) values ('John')"))
}

// TestShouldUseTransactionMixedCase verifies that mixed case queries are handled.
func TestShouldUseTransactionMixedCase(t *testing.T) {
	h := NewHandler()
	require.True(t, h.shouldUseTransaction("InSeRt INTO users (name) values ('John')"))
}

// TestNewHandlerInitialization verifies that NewHandler creates a handler with zero values.
func TestNewHandlerInitialization(t *testing.T) {
	h := NewHandler()
	require.NotNil(t, h)
	require.Nil(t, h.Storage)
	require.Equal(t, "", h.Query)
	require.False(t, h.Single)
}

// TestHandlerType verifies that the handler type is "sql".
func TestHandlerType(t *testing.T) {
	h := NewHandler()
	require.Equal(t, "sql", h.Type())
}
