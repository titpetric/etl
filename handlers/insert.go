package handlers

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/titpetric/etl/drivers"
	"github.com/titpetric/etl/internal"
	"github.com/titpetric/etl/model"
)

func buildInsertQuery(table string, data model.RecordInput) (string, []any) {
	// Step 1: List the keys
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}

	// Step 2: Construct placeholders string
	placeholders := strings.Repeat("?, ", len(keys))
	placeholders = strings.TrimSuffix(placeholders, ", ")

	// Step 3: Create a slice of values
	values := make([]any, 0, len(data))
	for _, key := range keys {
		values = append(values, data[key])
	}

	// Step 4: Construct ON DUPLICATE KEY UPDATE clause
	updates := make([]string, 0, len(data))
	for _, key := range keys {
		updates = append(updates, fmt.Sprintf("%s = VALUES(%s)", key, key))
	}

	// Step 5: Create the query string
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		table, strings.Join(keys, ", "), placeholders, csv(updates))

	return query, values
}

func Insert(ctx context.Context, command *model.Command, r io.Reader) error {
	driver, err := drivers.New(command.Driver, command.DB)
	if err != nil {
		return err
	}

	records, err := internal.DecodeRecords(r)
	if err != nil {
		return err
	}

	args := command.Args
	table := args[0]

	_, err = driver.Insert(table, records, args[1:]...)
	return err
}
