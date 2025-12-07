package server

import (
	"context"

	"github.com/titpetric/platform"

	"github.com/titpetric/etl/server/config"
	"github.com/titpetric/etl/server/config/loader"
)

// Start will load the config and start a HTTP server using the platform.
func Start(ctx context.Context) error {
	conf, err := NewConfig()
	if err != nil {
		return err
	}

	opts := platform.NewOptions()
	opts.ServerAddr = conf.Server.HttpAddr

	svc := platform.New(opts)
	svc.Register(NewModule(conf))

	return svc.Start(ctx)
}

// NewConfig will load the config from etl.yml.
func NewConfig() (*config.Config, error) {
	return loader.Load("etl.yml")
}
