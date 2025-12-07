package order

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAscOrder verifies that Asc creates an ascending order with the correct field and order value.
func TestAscOrder(t *testing.T) {
	result := Asc("name")
	require.Equal(t, "name", result.Field)
	require.Equal(t, "ASC", result.Order)
}

// TestAscOrderMultipleFields verifies that Asc works correctly with different field names.
func TestAscOrderMultipleFields(t *testing.T) {
	tests := []struct {
		field string
	}{
		{"id"},
		{"email"},
		{"created_at"},
		{"user_name"},
	}

	for _, tt := range tests {
		result := Asc(tt.field)
		require.Equal(t, tt.field, result.Field)
		require.Equal(t, "ASC", result.Order)
	}
}

// TestDescOrder verifies that Desc creates a descending order with the correct field and order value.
func TestDescOrder(t *testing.T) {
	result := Desc("age")
	require.Equal(t, "age", result.Field)
	require.Equal(t, "DESC", result.Order)
}

// TestDescOrderMultipleFields verifies that Desc works correctly with different field names.
func TestDescOrderMultipleFields(t *testing.T) {
	tests := []struct {
		field string
	}{
		{"id"},
		{"email"},
		{"updated_at"},
		{"score"},
	}

	for _, tt := range tests {
		result := Desc(tt.field)
		require.Equal(t, tt.field, result.Field)
		require.Equal(t, "DESC", result.Order)
	}
}

// TestOrderFieldPreservation verifies that the field name is preserved without modification.
func TestOrderFieldPreservation(t *testing.T) {
	field := "my_custom_field_123"
	asc := Asc(field)
	require.Equal(t, field, asc.Field)

	desc := Desc(field)
	require.Equal(t, field, desc.Field)
}

// TestOrderEmptyField verifies that empty fields are handled.
func TestOrderEmptyField(t *testing.T) {
	asc := Asc("")
	require.Equal(t, "", asc.Field)
	require.Equal(t, "ASC", asc.Order)

	desc := Desc("")
	require.Equal(t, "", desc.Field)
	require.Equal(t, "DESC", desc.Order)
}
