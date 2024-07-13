package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	. "github.com/titpetric/etl/internal/model"
	"github.com/titpetric/etl/internal/repository"
)

// commitCreate reads commits from stdin and inserts them into the database.
func commitCreate(ctx context.Context, command *Command) error {
	// Read from stdin
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var commits []*Commit
	if err := json.Unmarshal(input, &commits); err != nil {
		return err
	}

	tx, err := command.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	repo := repository.NewCommitRepository(tx)
	count := 0

	for _, commit := range commits {
		commit.SetCreatedAt(time.Now())
		if err := repo.Create(ctx, commit); err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				continue // Ignore unique constraint errors
			}
			return err
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		return tx.Commit()
	}

	fmt.Printf("Commits created successfully: %d rows added\n", count)
	return nil
}
