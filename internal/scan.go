package internal

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/etl/model"
)

// Scan scanns a single row and returns it.
func Scan(rows *sqlx.Rows) (model.Record, error) {
	row := model.RecordInput{}
	if err := rows.MapScan(row); err != nil {
		return nil, fmt.Errorf("error scanning result: %w", err)
	}
	return row.Record(), nil
}

// ScanAll scans all the rows and returns them.
func ScanAll(rows *sqlx.Rows) ([]model.Record, error) {
	var columns []model.Record

	for rows.Next() {
		row, err := Scan(rows)
		if err != nil {
			return nil, err
		}
		columns = append(columns, row)
	}

	if len(columns) == 0 {
		return nil, sql.ErrNoRows
	}

	return columns, nil
}
