package loader

import (
	"fmt"
	"sync"
)

// Cache is the interface any loader must implement.
type Cache interface {
	fmt.Stringer
	Get(filename string) (*Config, error)
}

type internalCache interface {
	set(filename string, data []byte) error
}

var (
	cache   Cache = NewCacheForeverManager()
	cacheMu sync.Mutex
)

// SetCache sets the global config cache.
func SetCache(c Cache) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	cache = c
}

// GetCache returns the global config cache.
func GetCache() Cache {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	return cache
}
