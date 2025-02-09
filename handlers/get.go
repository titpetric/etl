package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/titpetric/etl/model"
)

func Get(ctx context.Context, command *model.Command, _ io.Reader) error {
	var (
		all           bool
		offset, limit int
	)

	flagSet := model.NewFlagSet("Get")
	flagSet.IntVar(&offset, "offset", 0, "Offset for the results")
	flagSet.IntVar(&limit, "limit", 1, "Limit the number of results")
	flagSet.BoolVar(&all, "all", false, "Return all records")
	if err := flagSet.Parse(command.Args); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}
	args := flagSet.Args()

	table := args[0]

	var values []any
	var params []string
	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "-") {
			continue
		}

		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			parts[1] = strings.Trim(parts[1], "'\"")

			if strings.EqualFold(parts[1], "NULL") {
				params = append(params, parts[0]+" IS NULL")
				continue
			}

			params = append(params, parts[0]+" = ?")
			values = append(values, parts[1])
			continue
		}
		params = append(params, arg)
	}

	if len(params) == 0 {
		params = append(params, "1=1")
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, strings.Join(params, " "))
	if !all || limit > 1 {
		query = query + fmt.Sprintf(" LIMIT %d, %d", offset, limit)
	}
	if command.Verbose {
		log.Printf("-- %s %#v %#v\n", query, params, values)
	}
	rows, err := command.DB.Queryx(query, values...)
	if err != nil {
		return err
	}
	defer rows.Close()

	results, err := scanAllRecords(rows)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		os.Exit(1)
	}

	var data any = results
	if !all && limit == 1 {
		data = &results[0]
	}

	output, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
