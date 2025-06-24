package loader

import (
	"os"
	"sync"
)

// CacheShared holds the cached configuration data as a decoded value.
type CacheShared struct {
	config *Config
}

// CacheSharedManager manages a cache that clones the configuration data on retrieval.
type CacheSharedManager struct {
	mu      sync.RWMutex
	entries map[string]CacheShared
}

// NewCacheSharedManager creates a new CacheSharedManager, which caches the configuration data and clones it on retrieval.
func NewCacheSharedManager() *CacheSharedManager {
	return &CacheSharedManager{
		entries: make(map[string]CacheShared),
	}
}

// Get retrieves the cached configuration if available, or loads it if not, and returns a cloned copy.
func (c *CacheSharedManager) Get(filename string) (*Config, error) {
	c.mu.RLock()
	cacheEntry, exists := c.entries[filename]
	c.mu.RUnlock()

	if exists {
		return cacheEntry.config, nil
	}

	return c.loadAndSet(filename)
}

// Set stores the configuration in the cache.
func (c *CacheSharedManager) set(filename string, data []byte) error {
	cfg, err := Decode(data)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[filename] = CacheShared{
		config: cfg,
	}
	return nil
}

// loadAndSet loads the configuration file and updates the cache.
func (c *CacheSharedManager) loadAndSet(filename string) (*Config, error) {
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
func (c *CacheSharedManager) String() string {
	return "CacheSharedManager"
}
