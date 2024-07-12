package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/etl/internal/repository"
)

// commitList retrieves and prints the last 10 commits from the database.
func commitList(ctx context.Context, command *Command) error {
	db, err := sqlx.Open("sqlite3", config.GetDSN())
	if err != nil {
		return err
	}
	defer db.Close()

	repo := repository.NewCommitRepositoryReader(db)

	commits, err := repo.List(ctx, repository.ListOptions{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		return err
	}

	for _, commit := range commits {
		data, err := json.MarshalIndent(commit, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	}

	return nil
}
