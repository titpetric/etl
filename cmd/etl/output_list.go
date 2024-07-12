package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jmoiron/sqlx"

	. "github.com/titpetric/etl/internal/model"
	"github.com/titpetric/etl/internal/repository"
)

// outputList reads a Commit from stdin, fetches the outputs for that commit ID, and prints them
func outputList(ctx context.Context, command *Command) error {
	// Read from stdin
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var commit Commit
	if err := json.Unmarshal(input, &commit); err != nil {
		return err
	}

	db, err := sqlx.Open("sqlite3", config.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	repo := repository.NewOutputRepository(db)
	outputs, err := repo.ListByCommitID(ctx, commit.ID)
	if err != nil {
		return err
	}

	for _, output := range outputs {
		fmt.Printf("Filename: %s\nContents:\n%s\n", output.Filename, output.Contents)
	}

	return nil
}
