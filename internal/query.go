package internal

import (
	"database/sql"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/titpetric/etl/model"
)

func DecodeRecords(r io.Reader) ([]model.RecordInput, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if len(input) == 0 {
		return nil, sql.ErrNoRows
	}

	multi := false
	if string(input[0:1]) == "[" {
		multi = true
	}

	if !multi {
		var record model.RecordInput
		if err := json.Unmarshal(input, &record); err != nil {
			return nil, err
		}
		return []model.RecordInput{record}, nil
	}

	var records []model.RecordInput
	if err := json.Unmarshal(input, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func DecodeQuery(args []string) (model.RecordInput, error) {
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
