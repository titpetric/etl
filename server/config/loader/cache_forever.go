package loader

import (
	"os"
	"sync"
)

// CacheForever holds the cached configuration data indefinitely.
type CacheForever struct {
	data []byte
}

// CacheForeverManager manages a cache that caches configuration data indefinitely.
type CacheForeverManager struct {
	mu      sync.RWMutex
	entries map[string]CacheForever
}

// NewCacheForeverManager creates a new CacheForeverManager, which caches the configuration data indefinitely.
func NewCacheForeverManager() *CacheForeverManager {
	return &CacheForeverManager{
		entries: make(map[string]CacheForever),
	}
}

// Get retrieves the cached configuration if available, or loads it if not.
func (c *CacheForeverManager) Get(filename string) (*Config, error) {
	c.mu.RLock()
	cacheEntry, exists := c.entries[filename]
	c.mu.RUnlock()

	if exists {
		return Decode(cacheEntry.data)
	}

	return c.loadAndSet(filename)
}

// Set stores the configuration in the cache.
func (c *CacheForeverManager) set(filename string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[filename] = CacheForever{
		data: data,
	}
	return nil
}

// loadAndSet loads the configuration file and updates the cache.
func (c *CacheForeverManager) loadAndSet(filename string) (*Config, error) {
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
func (c *CacheForeverManager) String() string {
	return "CacheForever"
}
