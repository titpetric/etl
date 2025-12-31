package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/etl/model"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := start(ctx); err != nil {
		log.Fatal(err)
	}
}

func start(ctx context.Context) error {
	if len(os.Args) < 2 {
		return errors.New("usage: etl <command> <tableName> [options]")
	}

	config := model.NewConfig()
	args, err := config.ParseFlags()
	if err != nil {
		return err
	}

	db, err := sqlx.Open(config.GetDriver(), config.GetDSN())
	if err != nil {
		return err
	}
	defer db.Close()

	command := model.Command{
		DB:      db,
		Name:    args[0],
		Args:    args[1:],
		Verbose: config.Verbose,
		Quiet:   config.Quiet,
	}

	return HandleCommand(ctx, &command, getInput())
}

func getInput() io.Reader {
	fallback := strings.NewReader("{}")

	// Check if stdin is a terminal
	stat, err := os.Stdin.Stat()
	if err != nil {
		return fallback
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return fallback
	}

	return os.Stdin
}
