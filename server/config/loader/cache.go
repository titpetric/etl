package loader

import (
	"fmt"
)

// Cache is the interface any loader must implement.
type Cache interface {
	fmt.Stringer
	Get(filename string) (*Config, error)
}

type internalCache interface {
	set(filename string, data []byte) error
}
