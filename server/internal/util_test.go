package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMarshal verifies that Marshal correctly marshals valid data to JSON.
func TestMarshal(t *testing.T) {
	input := map[string]interface{}{"key": "value"}
	result := Marshal(input)
	require.NotNil(t, result)
	require.Equal(t, `{"key":"value"}`, string(result))
}

// TestMarshalString verifies that Marshal works with string values.
func TestMarshalString(t *testing.T) {
	result := Marshal("test string")
	require.Equal(t, `"test string"`, string(result))
}

// TestMarshalNumber verifies that Marshal works with numeric values.
func TestMarshalNumber(t *testing.T) {
	result := Marshal(42)
	require.Equal(t, `42`, string(result))
}

// TestMarshalSlice verifies that Marshal works with slice data.
func TestMarshalSlice(t *testing.T) {
	input := []int{1, 2, 3}
	result := Marshal(input)
	require.Equal(t, `[1,2,3]`, string(result))
}

// TestMarshalStruct verifies that Marshal works with struct data.
func TestMarshalStruct(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	input := Person{Name: "John", Age: 30}
	result := Marshal(input)
	require.NotNil(t, result)
	require.Contains(t, string(result), "John")
	require.Contains(t, string(result), "30")
}

// TestMarshalNull verifies that Marshal works with nil values.
func TestMarshalNull(t *testing.T) {
	result := Marshal(nil)
	require.Equal(t, `null`, string(result))
}

// TestMarshalPanic verifies that Marshal panics on invalid types.
func TestMarshalPanic(t *testing.T) {
	defer func() {
		recover := recover()
		require.NotNil(t, recover)
	}()

	// Create an unmarshalable type (like a channel)
	ch := make(chan int)
	Marshal(ch)
	t.Error("Expected panic, but none occurred")
}

// TestMarshalIndent verifies that MarshalIndent formats JSON with indentation.
func TestMarshalIndent(t *testing.T) {
	input := map[string]interface{}{"key": "value", "nested": map[string]interface{}{"inner": "data"}}
	result := MarshalIndent(input)
	require.NotNil(t, result)
	resultStr := string(result)
	require.Contains(t, resultStr, "  ") // Should contain indentation
	require.Contains(t, resultStr, "\n") // Should contain newlines
}

// TestMarshalIndentMultipleLines verifies that MarshalIndent creates proper multiline output.
func TestMarshalIndentMultipleLines(t *testing.T) {
	input := []map[string]string{
		{"name": "Alice"},
		{"name": "Bob"},
	}
	result := MarshalIndent(input)
	resultStr := string(result)
	lines := len(string(result))
	require.Greater(t, lines, 20) // Should be relatively large with indentation
	require.Contains(t, resultStr, "Alice")
	require.Contains(t, resultStr, "Bob")
}

// TestMarshalIndentWithNumber verifies that MarshalIndent works with numbers.
func TestMarshalIndentWithNumber(t *testing.T) {
	input := map[string]interface{}{"count": 42, "value": 3.14}
	result := MarshalIndent(input)
	resultStr := string(result)
	require.Contains(t, resultStr, "42")
	require.Contains(t, resultStr, "3.14")
}

// TestMarshalIndentIndentationFormat verifies that MarshalIndent uses two-space indentation.
func TestMarshalIndentIndentationFormat(t *testing.T) {
	input := map[string]interface{}{"outer": map[string]interface{}{"inner": "value"}}
	result := MarshalIndent(input)
	resultStr := string(result)
	// Should contain "  " for two-space indentation
	require.Contains(t, resultStr, "  ")
	// Should start with {
	require.True(t, resultStr[0] == '{')
	// Should end with } (allow for trailing newline)
	lastNonNewline := len(resultStr) - 1
	for lastNonNewline >= 0 && (resultStr[lastNonNewline] == '\n' || resultStr[lastNonNewline] == '\r') {
		lastNonNewline--
	}
	require.True(t, lastNonNewline >= 0 && resultStr[lastNonNewline] == '}')
}

// TestMarshalIndentPanic verifies that MarshalIndent panics on invalid types.
func TestMarshalIndentPanic(t *testing.T) {
	defer func() {
		recover := recover()
		require.NotNil(t, recover)
	}()

	// Create an unmarshalable type
	ch := make(chan struct{})
	MarshalIndent(ch)
	t.Error("Expected panic, but none occurred")
}
