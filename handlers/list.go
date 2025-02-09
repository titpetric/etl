package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/etl/model"
)

func List(ctx context.Context, command *model.Command, _ io.Reader) error {
	var offset, limit int
	var order, sortBy string

	flagSet := model.NewFlagSet("List")
	flagSet.StringVar(&sortBy, "sort-by", "id", "Sort by field")
	flagSet.StringVar(&order, "order", "desc", "Order")
	flagSet.IntVar(&offset, "offset", 0, "Offset for the results")
	flagSet.IntVar(&limit, "limit", 1000, "Limit the number of results")
	if err := flagSet.Parse(command.Args); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}
	args := flagSet.Args()

	if order != "asc" {
		order = "desc"
	}

	var err error
	var rows *sqlx.Rows

	table := args[0]
	if len(args) > 1 {
		query := fmt.Sprintf("SELECT * FROM %s WHERE id=?", table)
		rows, err = command.DB.Queryx(query, command.Args[1])
	} else {
		query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s %s LIMIT %d OFFSET %d", table, sortBy, order, limit, offset)
		rows, err = command.DB.Queryx(query)
	}
	if err != nil {
		return err
	}
	defer rows.Close()

	results, err := scanAllRecords(rows)
	if err != nil {
		return err
	}

	var output []byte
	if len(command.Args) > 1 {
		var res *model.Record
		if len(results) > 0 {
			res = &results[0]
		}
		output, err = json.Marshal(res)
	} else {
		output, err = json.Marshal(results)
	}

	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
