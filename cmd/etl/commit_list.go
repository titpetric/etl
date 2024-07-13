package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/titpetric/etl/internal/repository"
)

// commitList retrieves and prints the last 10 commits from the database.
func commitList(ctx context.Context, command *Command) error {
	repo := repository.NewCommitRepositoryReader(command.DB)

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
