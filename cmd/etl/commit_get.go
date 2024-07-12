package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/etl/internal/repository"
)

// commitList retrieves and prints the last 10 commits from the database.
func commitGet(ctx context.Context, command *Command) error {
	if len(os.Args) < 3 {
		return errors.New("usage: etl commit get <id> [options]")
	}

	db, err := sqlx.Open("sqlite3", config.GetDSN())
	if err != nil {
		return err
	}
	defer db.Close()

	repo := repository.NewCommitRepositoryReader(db)

	commit, err := repo.Get(ctx, command.Args[0])
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(commit, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}
