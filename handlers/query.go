package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/titpetric/etl/drivers"
	"github.com/titpetric/etl/internal"
	"github.com/titpetric/etl/model"
)

func Query(ctx context.Context, command *model.Command, _ io.Reader) error {
	driver, err := drivers.New(command.DB)
	if err != nil {
		return err
	}

	flagSet := model.NewFlagSet("Query")
	if err := flagSet.Parse(command.Args); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}
	args := flagSet.Args()

	query, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	stmts := internal.Statements(query)

	for idx, stmt := range stmts {
		if command.Verbose {
			log.Printf("-- %s %#v\n", stmt, args[1:])
		}

		result, err := driver.Query(stmt, args[1:]...)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Statement #%d OK", idx)
				continue
			}
			return err
		}

		return json.NewEncoder(os.Stdout).Encode(result)
	}
	if len(stmts) == 0 {
		return sql.ErrNoRows
	}
	return nil
}
