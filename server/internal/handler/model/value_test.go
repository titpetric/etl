package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDBValueString verifies that DBValue converts byte slices to strings.
func TestDBValueString(t *testing.T) {
	input := []byte("test data")
	result := DBValue(input)
	require.Equal(t, "test data", result)
}

// TestDBValueNil verifies that DBValue returns empty string for nil values.
func TestDBValueNil(t *testing.T) {
	result := DBValue(nil)
	require.Equal(t, "", result)
}

// TestDBValueString_ verifies that DBValue converts string values directly.
func TestDBValueStringDirect(t *testing.T) {
	input := "hello world"
	result := DBValue(input)
	require.Equal(t, "hello world", result)
}

// TestDBValueInteger verifies that DBValue converts integers to string representation.
func TestDBValueInteger(t *testing.T) {
	result := DBValue(42)
	require.Equal(t, "42", result)
}

// TestDBValueFloat verifies that DBValue converts floats to string representation.
func TestDBValueFloat(t *testing.T) {
	result := DBValue(3.14)
	require.Equal(t, "3.14", result)
}

// TestDBValueBool verifies that DBValue converts booleans to string representation.
func TestDBValueBool(t *testing.T) {
	resultTrue := DBValue(true)
	require.Equal(t, "true", resultTrue)

	resultFalse := DBValue(false)
	require.Equal(t, "false", resultFalse)
}

// TestDBValueEmptyBytes verifies that DBValue handles empty byte slices correctly.
func TestDBValueEmptyBytes(t *testing.T) {
	result := DBValue([]byte{})
	require.Equal(t, "", result)
}

// TestDBValueZero verifies that DBValue converts zero values correctly.
func TestDBValueZero(t *testing.T) {
	resultInt := DBValue(0)
	require.Equal(t, "0", resultInt)

	resultFloat := DBValue(0.0)
	require.Equal(t, "0", resultFloat)
}

// TestDBValueNegativeNumber verifies that DBValue handles negative numbers correctly.
func TestDBValueNegativeNumber(t *testing.T) {
	result := DBValue(-123)
	require.Equal(t, "-123", result)
}

// TestDBValueLargeNumber verifies that DBValue handles large numbers correctly.
func TestDBValueLargeNumber(t *testing.T) {
	result := DBValue(9999999999)
	require.Equal(t, "9999999999", result)
}

// TestDBValueByteSliceWithSpecialChars verifies that DBValue preserves special characters in byte slices.
func TestDBValueByteSliceWithSpecialChars(t *testing.T) {
	input := []byte("test@!#$%^&*()")
	result := DBValue(input)
	require.Equal(t, "test@!#$%^&*()", result)
}

// TestDBValueByteSliceWithUTF8 verifies that DBValue preserves UTF-8 characters in byte slices.
func TestDBValueByteSliceWithUTF8(t *testing.T) {
	input := []byte("hello 世界 مرحبا")
	result := DBValue(input)
	require.Equal(t, "hello 世界 مرحبا", result)
}

// TestDBValueByteSliceTypePriority verifies that byte slice takes priority over other types.
func TestDBValueByteSliceTypePriority(t *testing.T) {
	// This test verifies the type check order: []byte is checked first
	input := []byte("123")
	result := DBValue(input)
	require.Equal(t, "123", result)
	require.NotEqual(t, 123, result) // Should be string "123", not integer 123
}
