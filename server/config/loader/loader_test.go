package loader

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func benchLoad(tb *testing.B, load func(string) (*Config, error)) *Config {
	tb.Helper()

	cfg, err := load(testConfig.Path)
	assert.NoError(tb, err)
	assert.NotNil(tb, cfg)
	return cfg
}

func benchmarkLoaderFunc(cache Cache) func(b *testing.B) {
	return func(b *testing.B) {
		loader := NewConfigLoader(cache)
		b.Run(cache.String(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchLoad(b, loader.Load)
			}
		})
	}
}

func BenchmarkCacheSharedManager(b *testing.B) {
	cache := NewCacheSharedManager(testConfig.Storage)
	benchmarkLoaderFunc(cache)(b)
}

func BenchmarkCacheCloneManager(b *testing.B) {
	cache := NewCacheCloneManager(testConfig.Storage)
	benchmarkLoaderFunc(cache)(b)
}

func BenchmarkCacheModifiedManager(b *testing.B) {
	cache := NewCacheModifiedManager(testConfig.Storage)
	benchmarkLoaderFunc(cache)(b)
}

func BenchmarkCacheExpiryManager(b *testing.B) {
	cache := NewCacheExpiryManager(testConfig.Storage, time.Second)
	benchmarkLoaderFunc(cache)(b)
}

func BenchmarkCacheForeverManager(b *testing.B) {
	cache := NewCacheForeverManager(testConfig.Storage)
	benchmarkLoaderFunc(cache)(b)
}

func BenchmarkCacheNone(b *testing.B) {
	cache := NewCacheNone(testConfig.Storage)
	benchmarkLoaderFunc(cache)(b)
}
