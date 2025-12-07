package loader

import (
	"os"
	"path"
	"strings"

	"github.com/titpetric/etl/server/config"
)

// Loader is responsible for loading the configuration file using the specified cache.
type Loader struct {
	cache Cache
}

// New creates a new Loader with the given cache.
func New(cache Cache) *Loader {
	return &Loader{cache: cache}
}

// Load loads the configuration file using the specified cache.
func (l *Loader) Load(filename string) (*Config, error) {
	return l.cache.Get(filename)
}

// Load will not cache the result by default and read files each time.
func Load(filename string) (*config.Config, error) {
	if !strings.Contains(filename, "/") {
		filename = "./" + filename
	}

	file := path.Base(filename)
	dir := os.DirFS(path.Dir(filename))
	cache := NewCacheNone(dir)

	return New(cache).Load(file)
}
