package main

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/exp/maps"

	"github.com/titpetric/etl/handlers"
	"github.com/titpetric/etl/model"
)

type CommandHandlerFunc func(ctx context.Context, command *model.Command, r io.Reader) error

func HandleCommand(ctx context.Context, command *model.Command, r io.Reader) error {
	commandMap := map[string]CommandHandlerFunc{
		"insert":  handlers.Insert,
		"get":     handlers.Get,
		"list":    handlers.List,
		"tables":  handlers.Tables,
		"update":  handlers.Update,
		"query":   handlers.Query,
		"version": handlers.Version,
	}
	commands := maps.Keys(commandMap)

	if fn, ok := commandMap[command.Name]; ok {
		return fn(ctx, command, r)
	}
	return fmt.Errorf("unknown command: %s, supported %v", command.Name, commands)
}
