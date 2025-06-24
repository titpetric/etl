package handlers

import (
	"context"
	"io"

	"github.com/titpetric/etl/model"
	"github.com/titpetric/etl/server"
)

func Server(ctx context.Context, command *model.Command, _ io.Reader) error {
	return server.Start(ctx)
}
