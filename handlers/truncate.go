package handlers

import (
	"context"
	"fmt"
	"io"

	"github.com/titpetric/etl/model"
)

func Truncate(ctx context.Context, command *model.Command, _ io.Reader) error {
	var skipForeignKeyChecks bool

	flagSet := model.NewFlagSet("Get")
	flagSet.BoolVar(&skipForeignKeyChecks, "skip-foreign-key-checks", false, "Skip foreign key checks")
	if err := flagSet.Parse(command.Args); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}
	args := flagSet.Args()

	if skipForeignKeyChecks {
		if _, err := command.DB.Exec("SET FOREIGN_KEY_CHECKS = 0;"); err != nil {
			fmt.Println("Error disabling foreign key checks: %w", err)
		}
		defer func() {
			if _, err := command.DB.Exec("SET FOREIGN_KEY_CHECKS = 1;"); err != nil {
				fmt.Printf("Error re-enabling foreign key checks: %s\n", err)
			}
		}()
	}

	for _, table := range args {
		query := "TRUNCATE " + table
		if command.Verbose {
			fmt.Printf("-- %s\n", query)
		}

		_, err := command.DB.Exec(query)
		if err != nil {
			return fmt.Errorf("Error truncating: %w", err)
		}
	}

	return nil
}
