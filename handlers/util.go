package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/etl/model"
)

func scanRecord(rows *sqlx.Rows) (model.Record, error) {
	row := model.RecordInput{}
	if err := rows.MapScan(row); err != nil {
		return nil, fmt.Errorf("error scanning result: %w", err)
	}
	return row.Record(), nil
}

func scanAllRecords(rows *sqlx.Rows) ([]model.Record, error) {
	var columns []model.Record

	for rows.Next() {
		row, err := scanRecord(rows)
		if err != nil {
			return nil, err
		}
		columns = append(columns, row)
	}

	return columns, nil
}

func decodeQueryParameters(args []string) (model.RecordInput, error) {
	result := model.RecordInput{}
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			parts[1] = strings.Trim(parts[1], "'\"")

			if strings.HasPrefix(parts[1], "@") {
				contents, err := os.ReadFile(parts[1][1:])
				if err != nil {
					return nil, err
				}
				parts[1] = string(contents)
			}

			k, v := parts[0], parts[1]
			result[k] = v
			continue
		}
	}

	return result, nil
}

func csv(in []string) string {
	return strings.Join(in, ",")
}
