package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/jmoiron/sqlx"

	. "github.com/titpetric/etl/internal/model"
	"github.com/titpetric/etl/internal/repository"
)

// outputRestore reads a Commit from stdin, fetches the outputs for that commit ID, and restores the folder
func outputRestore(ctx context.Context, command *Command) error {
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
		if err := os.WriteFile(path.Join(config.Folder, output.Filename), []byte(output.Contents), os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
