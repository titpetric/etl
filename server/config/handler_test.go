package config

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testConfigLoad(tb testing.TB) *Config {
	tb.Helper()

	data, err := fs.ReadFile(testConfig.Storage, testConfig.Path)
	assert.NoError(tb, err)

	cfg, err := Decode(testConfig.Storage, data)
	assert.NoError(tb, err)
	assert.NotNil(tb, cfg)
	assert.True(tb, "/api/ping" == cfg.Endpoints[0].Path)
	assert.True(tb, "ping" == cfg.Endpoints[0].Handler.Type)

	return cfg
}

func TestHandler_Decode(t *testing.T) {
	cfg := testConfigLoad(t)

	data := map[string]string{}
	endpoint := cfg.Endpoints[2]

	assert.NoError(t, endpoint.Handler.Decode(&data))

	val, ok := data["network"]
	assert.True(t, ok)
	assert.Equal(t, "tcp", val)
}
