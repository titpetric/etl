package cache

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestNewMemoryStore verifies that NewMemoryStore creates a store with an initialized map.
func TestNewMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	require.NotNil(t, store)
	require.NotNil(t, store.store)
	require.Len(t, store.store, 0)
}

// TestMemoryStoreSet verifies that Set stores an entry in the cache.
func TestMemoryStoreSet(t *testing.T) {
	store := NewMemoryStore()
	entry := &Entry{
		StatusCode: http.StatusOK,
		Headers:    http.Header{},
		Body:       []byte("test"),
	}

	err := store.Set(context.Background(), "key1", entry, 5*time.Minute)
	require.NoError(t, err)

	stored, ok := store.store["key1"]
	require.True(t, ok)
	require.Equal(t, entry.StatusCode, stored.StatusCode)
	require.Equal(t, entry.Body, stored.Body)
}

// TestMemoryStoreGet verifies that Get retrieves a stored entry.
func TestMemoryStoreGet(t *testing.T) {
	store := NewMemoryStore()
	entry := &Entry{
		StatusCode: http.StatusOK,
		Headers:    http.Header{},
		Body:       []byte("test"),
	}

	store.Set(context.Background(), "key1", entry, 5*time.Minute)

	retrieved, err := store.Get(context.Background(), "key1")
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	require.Equal(t, entry.StatusCode, retrieved.StatusCode)
	require.Equal(t, entry.Body, retrieved.Body)
}

// TestMemoryStoreGetNonExistent verifies that Get returns nil for non-existent keys.
func TestMemoryStoreGetNonExistent(t *testing.T) {
	store := NewMemoryStore()

	retrieved, err := store.Get(context.Background(), "nonexistent")
	require.NoError(t, err)
	require.Nil(t, retrieved)
}

// TestMemoryStoreDelete verifies that Delete removes an entry from the cache.
func TestMemoryStoreDelete(t *testing.T) {
	store := NewMemoryStore()
	entry := &Entry{
		StatusCode: http.StatusOK,
		Headers:    http.Header{},
		Body:       []byte("test"),
	}

	store.Set(context.Background(), "key1", entry, 5*time.Minute)

	err := store.Delete(context.Background(), "key1")
	require.NoError(t, err)

	retrieved, _ := store.Get(context.Background(), "key1")
	require.Nil(t, retrieved)
}

// TestMemoryStoreDeleteNonExistent verifies that Delete on non-existent key returns no error.
func TestMemoryStoreDeleteNonExistent(t *testing.T) {
	store := NewMemoryStore()

	err := store.Delete(context.Background(), "nonexistent")
	require.NoError(t, err)
}

// TestMemoryStoreClear verifies that Clear removes all entries from the cache.
func TestMemoryStoreClear(t *testing.T) {
	store := NewMemoryStore()

	entry := &Entry{
		StatusCode: http.StatusOK,
		Headers:    http.Header{},
		Body:       []byte("test"),
	}

	store.Set(context.Background(), "key1", entry, 5*time.Minute)
	store.Set(context.Background(), "key2", entry, 5*time.Minute)

	err := store.Clear(context.Background())
	require.NoError(t, err)

	require.Len(t, store.store, 0)

	retrieved1, _ := store.Get(context.Background(), "key1")
	require.Nil(t, retrieved1)

	retrieved2, _ := store.Get(context.Background(), "key2")
	require.Nil(t, retrieved2)
}

// TestMemoryStoreSetWithZeroTTL verifies that Set uses default TTL when TTL is 0.
func TestMemoryStoreSetWithZeroTTL(t *testing.T) {
	store := NewMemoryStore()
	entry := &Entry{
		StatusCode: http.StatusOK,
		Headers:    http.Header{},
		Body:       []byte("test"),
	}

	before := time.Now()
	err := store.Set(context.Background(), "key1", entry, 0)
	after := time.Now()

	require.NoError(t, err)

	stored, ok := store.store["key1"]
	require.True(t, ok)

	// Should be set with default 5 minute TTL
	require.True(t, stored.ExpiresAt.After(before))
	require.True(t, stored.ExpiresAt.Before(after.Add(5*time.Minute)))
}

// TestMemoryStoreContextCancellation verifies that operations respect context cancellation.
func TestMemoryStoreContextCancellation(t *testing.T) {
	store := NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	entry := &Entry{
		StatusCode: http.StatusOK,
		Headers:    http.Header{},
		Body:       []byte("test"),
	}

	err := store.Set(ctx, "key1", entry, 5*time.Minute)
	require.Error(t, err)
	require.Equal(t, context.Canceled, err)
}

// TestMemoryStoreGetContextCancellation verifies that Get respects context cancellation.
func TestMemoryStoreGetContextCancellation(t *testing.T) {
	store := NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := store.Get(ctx, "key1")
	require.Error(t, err)
	require.Equal(t, context.Canceled, err)
}

// TestMemoryStoreMultipleEntries verifies that store handles multiple entries correctly.
func TestMemoryStoreMultipleEntries(t *testing.T) {
	store := NewMemoryStore()

	for i := 1; i <= 5; i++ {
		entry := &Entry{
			StatusCode: http.StatusOK,
			Headers:    http.Header{},
			Body:       []byte("test" + string(rune(i))),
		}

		err := store.Set(context.Background(), "key"+string(rune(i+48)), entry, 5*time.Minute)
		require.NoError(t, err)
	}

	require.Len(t, store.store, 5)
}

// TestMemoryStoreOverwrite verifies that setting a key twice overwrites the previous value.
func TestMemoryStoreOverwrite(t *testing.T) {
	store := NewMemoryStore()

	entry1 := &Entry{
		StatusCode: http.StatusOK,
		Headers:    http.Header{},
		Body:       []byte("test1"),
	}

	entry2 := &Entry{
		StatusCode: http.StatusNotFound,
		Headers:    http.Header{},
		Body:       []byte("test2"),
	}

	store.Set(context.Background(), "key1", entry1, 5*time.Minute)
	store.Set(context.Background(), "key1", entry2, 5*time.Minute)

	retrieved, _ := store.Get(context.Background(), "key1")
	require.NotNil(t, retrieved)
	require.Equal(t, http.StatusNotFound, retrieved.StatusCode)
	require.Equal(t, []byte("test2"), retrieved.Body)
}
