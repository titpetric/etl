package handlers

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/titpetric/etl/drivers"
	"github.com/titpetric/etl/model"
)

// Tables retrieves the list of tables in the current database schema along with their comments.
func Tables(ctx context.Context, command *model.Command, _ io.Reader) error {
	driver, err := drivers.New(command.DB)
	if err != nil {
		return err
	}

	tables, err := driver.Tables()
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(tables)
}
