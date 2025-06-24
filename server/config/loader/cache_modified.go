package loader

import (
	"os"
	"sync"
	"time"
)

// CacheModified holds the cached configuration data based on file modification time.
type CacheModified struct {
	modTime time.Time
	data    []byte
}

// CacheModifiedManager manages a cache that invalidates entries based on file modification time.
type CacheModifiedManager struct {
	mu      sync.RWMutex
	entries map[string]CacheModified
}

// NewCacheModifiedManager creates a new CacheModifiedManager, which caches the configuration data and invalidates it based on file modification time.
func NewCacheModifiedManager() *CacheModifiedManager {
	return &CacheModifiedManager{
		entries: make(map[string]CacheModified),
	}
}

// Get retrieves the cached configuration if available and valid based on file modification time, or loads it if not.
func (c *CacheModifiedManager) Get(filename string) (*Config, error) {
	c.mu.RLock()
	cacheEntry, exists := c.entries[filename]
	c.mu.RUnlock()

	if exists {
		stat, err := os.Stat(filename)
		if err != nil {
			return nil, err
		}

		if stat.ModTime().After(cacheEntry.modTime) {
			return c.loadAndSet(filename)
		}

		return Decode(cacheEntry.data)
	}

	return c.loadAndSet(filename)
}

// Set stores the configuration in the cache with the file modification time.
func (c *CacheModifiedManager) set(filename string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[filename] = CacheModified{
		modTime: time.Now(),
		data:    data,
	}

	return nil
}

// loadAndSet loads the configuration file and updates the cache.
func (c *CacheModifiedManager) loadAndSet(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := c.set(filename, data); err != nil {
		return nil, err
	}

	return Decode(data)
}

// String returns the name of the cache implementation.
func (c *CacheModifiedManager) String() string {
	return "CacheModifiedManager"
}
