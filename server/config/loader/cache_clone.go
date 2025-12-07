package loader

import (
	"io/fs"
	"sync"

	clone "github.com/huandu/go-clone"
)

// CacheClone holds the cached configuration data as a decoded value.
type CacheClone struct {
	config *Config
}

// CacheCloneManager manages a cache that clones the configuration data on retrieval.
type CacheCloneManager struct {
	mu      sync.RWMutex
	storage fs.FS
	entries map[string]CacheClone
}

// NewCacheCloneManager creates a new CacheCloneManager, which caches the configuration data and clones it on retrieval.
func NewCacheCloneManager(storage fs.FS) *CacheCloneManager {
	return &CacheCloneManager{
		storage: storage,
		entries: make(map[string]CacheClone),
	}
}

// Get retrieves the cached configuration if available, or loads it if not, and returns a cloned copy.
func (c *CacheCloneManager) Get(filename string) (*Config, error) {
	c.mu.RLock()
	cacheEntry, exists := c.entries[filename]
	if exists {
		defer c.mu.RUnlock()
		return clone.Clone(cacheEntry.config).(*Config), nil
	}
	c.mu.RUnlock()

	return c.loadAndSet(filename)
}

// Set stores the configuration in the cache.
func (c *CacheCloneManager) set(filename string, data []byte) error {
	cfg, err := Decode(c.storage, data)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[filename] = CacheClone{
		config: cfg,
	}
	return nil
}

// loadAndSet loads the configuration file and updates the cache.
func (c *CacheCloneManager) loadAndSet(filename string) (*Config, error) {
	data, err := fs.ReadFile(c.storage, filename)
	if err != nil {
		return nil, err
	}

	if err := c.set(filename, data); err != nil {
		return nil, err
	}

	return c.Get(filename)
}

// String returns the name of the cache implementation.
func (c *CacheCloneManager) String() string {
	return "CacheCloneManager"
}
