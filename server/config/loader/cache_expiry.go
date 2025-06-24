package loader

import (
	"os"
	"sync"
	"time"
)

// CacheExpiry holds the cached configuration data with an expiration time.
type CacheExpiry struct {
	expiry time.Time
	data   []byte
}

// CacheExpiryManager manages a cache that invalidates entries based on a time-to-live (TTL) value.
type CacheExpiryManager struct {
	mu      sync.RWMutex
	entries map[string]CacheExpiry
	ttl     time.Duration
}

// NewCacheExpiryManager creates a new CacheExpiryManager with the given TTL.
func NewCacheExpiryManager(ttl time.Duration) *CacheExpiryManager {
	return &CacheExpiryManager{
		entries: make(map[string]CacheExpiry),
		ttl:     ttl,
	}
}

// Get retrieves the cached configuration if available and valid based on expiry time, or loads it if not.
func (c *CacheExpiryManager) Get(filename string) (*Config, error) {
	c.mu.RLock()
	cacheEntry, exists := c.entries[filename]
	c.mu.RUnlock()

	if exists && time.Now().Before(cacheEntry.expiry) {
		return Decode(cacheEntry.data)
	}

	return c.loadAndSet(filename)
}

// Set stores the configuration in the cache with an expiry time.
func (c *CacheExpiryManager) set(filename string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[filename] = CacheExpiry{
		expiry: time.Now().Add(c.ttl),
		data:   data,
	}
	return nil
}

// loadAndSet loads the configuration file and updates the cache.
func (c *CacheExpiryManager) loadAndSet(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := c.set(filename, data); err != nil {
		return nil, err
	}

	return c.Get(filename)
}

// String returns the name of the cache implementation.
func (c *CacheExpiryManager) String() string {
	return "CacheExpiryManager"
}
