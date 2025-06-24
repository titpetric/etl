package config

import "errors"

// Validate checks the validity of the configuration.
func (c *Config) Validate() error {
	if c.Server.HttpAddr == "" {
		return errors.New("http server listen address must be specified")
	}
	if c.Server.GrpcAddr == "" {
		return errors.New("grpc server listen address must be specified")
	}
	return nil
}
