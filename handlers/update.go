package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/titpetric/etl/model"
)

func UpdateRequest(r io.Reader, args []string) ([]model.Record, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	multi := false
	if string(input[0:1]) == "[" {
		multi = true
	}

	var records []model.Record
	if multi {
		if err := json.Unmarshal(input, &records); err != nil {
			return nil, err
		}
	} else {
		var record model.Record
		if err := json.Unmarshal(input, &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func buildUpdateQuery(table string, data model.Record, whereKeys []string) (string, []any) {
	var query = fmt.Sprintf("UPDATE %s SET ", table)

	keys := maps.Keys(data)

	values := []any{}
	for _, key := range keys {
		if slices.Contains(whereKeys, key) {
			continue
		}

		prefix := ", "
		if len(values) == 0 {
			prefix = ""
		}
		query += fmt.Sprintf("%s%s=? ", prefix, key)
		values = append(values, data[key])
	}
	query += "WHERE "
	for i, key := range whereKeys {
		prefix := ", "
		if i == 0 {
			prefix = ""
			query += fmt.Sprintf("%s%s=? ", prefix, key)
		}
		values = append(values, data[key])
	}

	return query, values
}

func Update(ctx context.Context, command *model.Command, r io.Reader) error {
	records, err := UpdateRequest(r, command.Args[1:])
	if err != nil {
		return err
	}

	table := command.Args[0]
	where := command.Args[1:]

	var values []string
	var names []string
	for _, arg := range where {
		if strings.HasPrefix(arg, "-") {
			continue
		}

		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			parts[1] = strings.Trim(parts[1], "'\"")

			if parts[1] == "NULL" {
				names = append(names, parts[0]+" IS NULL")
				continue
			}

			names = append(names, parts[0]+" = ?")
			values = append(values, parts[1])
			continue
		}
		names = append(names, arg)
	}

	if len(names) == 0 {
		return errors.New("no where condition for update")
	}

	tx, err := command.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var updates int64
	for _, record := range records {
		query, params := buildUpdateQuery(table, record, names)
		if command.Verbose {
			fmt.Printf("-- %s %#v\n", query, params)
		}

		result, err := tx.Exec(query, params...)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				continue // Ignore unique constraint errors
			}
			if strings.Contains(err.Error(), "Duplicate entry") {
				continue // Ignore unique constraint errors
			}
			return err
		}
		rowsAffected, _ := result.RowsAffected()
		updates += rowsAffected
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if !command.Quiet {
		fmt.Printf("%s: stored %d records, %d rows affected\n", table, len(records), updates)
	}
	return nil
}
