package ratelimit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewMemoryStore verifies that NewMemoryStore creates a store with an initialized map.
func TestNewMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	require.NotNil(t, store)
	require.NotNil(t, store.store)
	require.Len(t, store.store, 0)
}

// TestMemoryStoreInc verifies that Inc increments the counter for a key.
func TestMemoryStoreInc(t *testing.T) {
	store := NewMemoryStore()

	count, err := store.Inc(context.Background(), "key1")
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	count, err = store.Inc(context.Background(), "key1")
	require.NoError(t, err)
	require.Equal(t, int64(2), count)
}

// TestMemoryStoreIncMultipleKeys verifies that Inc handles multiple keys independently.
func TestMemoryStoreIncMultipleKeys(t *testing.T) {
	store := NewMemoryStore()

	count1, _ := store.Inc(context.Background(), "key1")
	count2, _ := store.Inc(context.Background(), "key2")
	count1Again, _ := store.Inc(context.Background(), "key1")

	require.Equal(t, int64(1), count1)
	require.Equal(t, int64(1), count2)
	require.Equal(t, int64(2), count1Again)
}

// TestMemoryStoreRate verifies that Rate returns the current count without incrementing.
func TestMemoryStoreRate(t *testing.T) {
	store := NewMemoryStore()

	store.Inc(context.Background(), "key1")
	store.Inc(context.Background(), "key1")

	count, err := store.Rate(context.Background(), "key1")
	require.NoError(t, err)
	require.Equal(t, int64(2), count)

	// Verify Rate doesn't increment
	countAgain, err := store.Rate(context.Background(), "key1")
	require.NoError(t, err)
	require.Equal(t, int64(2), countAgain)
}

// TestMemoryStoreRateNonExistent verifies that Rate returns 0 for non-existent keys.
func TestMemoryStoreRateNonExistent(t *testing.T) {
	store := NewMemoryStore()

	count, err := store.Rate(context.Background(), "nonexistent")
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}

// TestMemoryStoreReset verifies that Reset clears the counter for a key.
func TestMemoryStoreReset(t *testing.T) {
	store := NewMemoryStore()

	store.Inc(context.Background(), "key1")
	store.Inc(context.Background(), "key1")

	err := store.Reset(context.Background(), "key1")
	require.NoError(t, err)

	count, _ := store.Rate(context.Background(), "key1")
	require.Equal(t, int64(0), count)
}

// TestMemoryStoreResetNonExistent verifies that Reset on non-existent key returns no error.
func TestMemoryStoreResetNonExistent(t *testing.T) {
	store := NewMemoryStore()

	err := store.Reset(context.Background(), "nonexistent")
	require.NoError(t, err)
}

// TestMemoryStoreClear verifies that Clear removes all rate limit entries.
func TestMemoryStoreClear(t *testing.T) {
	store := NewMemoryStore()

	store.Inc(context.Background(), "key1")
	store.Inc(context.Background(), "key2")
	store.Inc(context.Background(), "key3")

	err := store.Clear(context.Background())
	require.NoError(t, err)

	require.Len(t, store.store, 0)

	count1, _ := store.Rate(context.Background(), "key1")
	require.Equal(t, int64(0), count1)
}

// TestMemoryStoreContextCancellationInc verifies that Inc respects context cancellation.
func TestMemoryStoreContextCancellationInc(t *testing.T) {
	store := NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := store.Inc(ctx, "key1")
	require.Error(t, err)
	require.Equal(t, context.Canceled, err)
}

// TestMemoryStoreContextCancellationGet verifies that Rate respects context cancellation.
func TestMemoryStoreContextCancellationGet(t *testing.T) {
	store := NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := store.Rate(ctx, "key1")
	require.Error(t, err)
	require.Equal(t, context.Canceled, err)
}

// TestMemoryStoreIncWithoutReset verifies that Inc increments within the same second.
func TestMemoryStoreIncWithoutReset(t *testing.T) {
	store := NewMemoryStore()

	for i := 1; i <= 5; i++ {
		count, _ := store.Inc(context.Background(), "key1")
		require.Equal(t, int64(i), count)
	}
}

// TestMemoryStoreMultipleKeys verifies that store handles multiple keys independently.
func TestMemoryStoreMultipleKeys(t *testing.T) {
	store := NewMemoryStore()

	for i := 1; i <= 5; i++ {
		key := "key" + string(rune('0'+i))
		count, _ := store.Inc(context.Background(), key)
		require.Equal(t, int64(1), count)
	}

	require.Len(t, store.store, 5)
}

// TestMemoryStoreIncConsistency verifies that Inc counts correctly over time.
func TestMemoryStoreIncConsistency(t *testing.T) {
	store := NewMemoryStore()

	for i := 1; i <= 100; i++ {
		count, _ := store.Inc(context.Background(), "key1")
		require.Equal(t, int64(i), count)
	}
}

// TestMemoryStoreClearContextCancellation verifies that Clear respects context cancellation.
func TestMemoryStoreClearContextCancellation(t *testing.T) {
	store := NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := store.Clear(ctx)
	require.Error(t, err)
	require.Equal(t, context.Canceled, err)
}

// TestMemoryStoreResetContextCancellation verifies that Reset respects context cancellation.
func TestMemoryStoreResetContextCancellation(t *testing.T) {
	store := NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := store.Reset(ctx, "key1")
	require.Error(t, err)
	require.Equal(t, context.Canceled, err)
}
