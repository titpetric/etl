package loader

// ConfigLoader is responsible for loading the configuration file using the specified cache.
type ConfigLoader struct {
	cache Cache
}

// NewConfigLoader creates a new ConfigLoader with the given cache.
func NewConfigLoader(cache Cache) *ConfigLoader {
	return &ConfigLoader{cache: cache}
}

// Load loads the configuration file using the specified cache.
func (loader *ConfigLoader) Load(filename string) (*Config, error) {
	return loader.cache.Get(filename)
}

// Load loads the configuration file using the global cache.
func Load(filename string) (*Config, error) {
	return GetCache().Get(filename)
}
