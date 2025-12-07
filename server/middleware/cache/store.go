package cache

import (
	"context"
	"sync"
	"time"
)

// Entry represents a cached response with its metadata
type Entry struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
	ExpiresAt  time.Time
}

// Store defines the interface for cache storage backends
type Store interface {
	// Get retrieves a cached entry by key
	Get(ctx context.Context, key string) (*Entry, error)

	// Set stores an entry in the cache with optional TTL
	Set(ctx context.Context, key string, entry *Entry, ttl time.Duration) error

	// Delete removes an entry from the cache
	Delete(ctx context.Context, key string) error

	// Clear removes all entries from the cache
	Clear(ctx context.Context) error
}

// MemoryStore is a simple in-memory cache implementation using map+mutex
type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]*Entry
}

// NewMemoryStore creates a new in-memory cache store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string]*Entry),
	}
}

// Get retrieves a cached entry if it exists and hasn't expired
func (ms *MemoryStore) Get(ctx context.Context, key string) (*Entry, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	ms.mu.RLock()
	defer ms.mu.RUnlock()

	entry, exists := ms.store[key]
	if !exists {
		return nil, nil
	}

	// Check expiration
	if time.Now().After(entry.ExpiresAt) {
		return nil, nil
	}

	return entry, nil
}

// Set stores an entry in the cache with a TTL
func (ms *MemoryStore) Set(ctx context.Context, key string, entry *Entry, ttl time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if ttl == 0 {
		ttl = 5 * time.Minute // default TTL
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	entry.ExpiresAt = time.Now().Add(ttl)
	ms.store[key] = entry

	return nil
}

// Delete removes an entry from the cache
func (ms *MemoryStore) Delete(ctx context.Context, key string) error {
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

// Clear removes all entries from the cache
func (ms *MemoryStore) Clear(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.store = make(map[string]*Entry)
	return nil
}

// CleanupExpired removes all expired entries from the cache
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
	for key, entry := range ms.store {
		if now.After(entry.ExpiresAt) {
			delete(ms.store, key)
		}
	}

	return nil
}
