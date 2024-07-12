package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
)

type Command struct {
	Key    string
	SubKey string
	Args   []string
}

var dbDSN string

func main() {
	pflag.StringVar(&dbDSN, "db-dsn", "file:etl.db", "Database DSN")
	pflag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := start(ctx); err != nil {
		log.Fatal(err)
	}
}

func start(ctx context.Context) error {
	if len(os.Args) < 3 {
		return errors.New("usage: etl <command> <subcommand> [options]")
	}

	command := Command{
		Key:    os.Args[1],
		SubKey: os.Args[2],
		Args:   os.Args[3:],
	}

	index := map[string]map[string]func(context.Context) error{
		"commit": {
			"create": commitCreate,
			"list":   commitList,
		},
	}

	subcommands, ok := index[command.Key]
	if !ok {
		return fmt.Errorf("unknown command: %s", command.Key)
	}

	action, ok := subcommands[command.SubKey]
	if !ok {
		return fmt.Errorf("unknown subcommand: %s", command.SubKey)
	}

	dbDSN = dbDSN + "?parseTime=true"

	return action(ctx)
}
