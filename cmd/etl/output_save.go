package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/titpetric/etl/internal/model"
	"github.com/titpetric/etl/internal/repository"
)

// outputSave reads files from `config.Folder` and processes commits from stdin
func outputSave(ctx context.Context, command *Command) error {
	// Read from stdin
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var commit Commit
	if err := json.Unmarshal(input, &commit); err != nil {
		return err
	}

	repo := repository.NewOutputRepository(command.DB)
	count, skipped := 0, 0

	files, err := ioutil.ReadDir(config.Folder)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(config.Folder, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		commitOutput := &CommitOutput{
			CommitID:    commit.ID,
			CreatedWith: "input",
			Filename:    file.Name(),
			Contents:    string(content),
			CreatedAt:   nil,
		}

		if err := repo.Create(ctx, commitOutput); err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				skipped++
				continue // Ignore unique constraint errors
			}
			return err
		}
		count++
	}

	fmt.Printf("Outputs created successfully: %d files added, %d skipped\n", count, skipped)
	return nil
}
