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

	_ "github.com/mattn/go-sqlite3"
)

type Command struct {
	Key    string
	SubKey string
	Args   []string
}

type Config struct {
	DSN    string
	Folder string
}

func (c *Config) GetDSN() string {
	return c.DSN + "?parseTime=true"
}

var config = &Config{}

func main() {
	pflag.StringVar(&config.DSN, "db-dsn", "file:etl.db", "Database DSN")
	pflag.StringVarP(&config.Folder, "folder", "f", "output", "Folder with outputs")
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

	index := map[string]map[string]func(context.Context, *Command) error{
		"commit": {
			"create": commitCreate,
			"list":   commitList,
		},
		"output": {
			"list":    outputList,
			"save":    outputSave,
			"restore": outputRestore,
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

	return action(ctx, &command)
}
