package loader

import "io/fs"

// CacheNone is a cache implementation that does not cache configuration data.
type CacheNone struct {
	storage fs.FS
}

// NewCacheNone creates a new CacheNone, which does not cache configuration data and reads directly from the file each time.
func NewCacheNone(storage fs.FS) *CacheNone {
	return &CacheNone{
		storage: storage,
	}
}

// Get reads the configuration file directly without caching.
func (c *CacheNone) Get(filename string) (*Config, error) {
	data, err := fs.ReadFile(c.storage, filename)
	if err != nil {
		return nil, err
	}

	return Decode(c.storage, data)
}

// Set is a no-op for CacheNone.
func (c *CacheNone) set(_ string, _ []byte) error {
	// No-op for CacheNone
	return nil
}

// String returns the name of the cache implementation.
func (c *CacheNone) String() string {
	return "CacheNone"
}
