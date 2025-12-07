package ratelimit

import (
	"context"
	"sync"
	"time"
)

// RateLimit holds rate limit information for a key
type RateLimit struct {
	Count     int64
	ResetTime time.Time
}

// Store defines the interface for rate limit storage backends
type Store interface {
	// Inc increments the counter for a key and returns the current count
	Inc(ctx context.Context, key string) (int64, error)

	// Rate returns the current count for a key without incrementing
	Rate(ctx context.Context, key string) (int64, error)

	// Reset resets the counter for a key
	Reset(ctx context.Context, key string) error

	// Clear removes all rate limit entries
	Clear(ctx context.Context) error
}

// MemoryStore is a simple in-memory rate limit store using map+mutex
type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]*RateLimit
}

// NewMemoryStore creates a new in-memory rate limit store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string]*RateLimit),
	}
}

// Inc increments the counter for a key and returns the new count
func (ms *MemoryStore) Inc(ctx context.Context, key string) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	rl, exists := ms.store[key]
	if !exists {
		rl = &RateLimit{
			Count:     0,
			ResetTime: time.Now().Add(1 * time.Second),
		}
		ms.store[key] = rl
	}

	// Check if we need to reset the counter
	if time.Now().After(rl.ResetTime) {
		rl.Count = 0
		rl.ResetTime = time.Now().Add(1 * time.Second)
	}

	rl.Count++
	return rl.Count, nil
}

// Rate returns the current count for a key without incrementing
func (ms *MemoryStore) Rate(ctx context.Context, key string) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	ms.mu.RLock()
	defer ms.mu.RUnlock()

	rl, exists := ms.store[key]
	if !exists {
		return 0, nil
	}

	// Check if the rate limit has reset
	if time.Now().After(rl.ResetTime) {
		return 0, nil
	}

	return rl.Count, nil
}

// Reset resets the counter for a key
func (ms *MemoryStore) Reset(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	delete(ms.store, key)
	return nil
}

// Clear removes all rate limit entries
func (ms *MemoryStore) Clear(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.store = make(map[string]*RateLimit)
	return nil
}

// CleanupExpired removes all expired rate limit entries from the store
// This should be called periodically to free up memory
func (ms *MemoryStore) CleanupExpired(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	now := time.Now()
	for key, rl := range ms.store {
		if now.After(rl.ResetTime) {
			delete(ms.store, key)
		}
	}

	return nil
}
