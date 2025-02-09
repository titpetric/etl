package drivers

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/maps"

	"github.com/titpetric/etl/internal"
	"github.com/titpetric/etl/model"
)

type MySQL struct {
	db *sqlx.DB
}

func NewMySQL(db *sqlx.DB) (*MySQL, error) {
	return &MySQL{
		db: db,
	}, nil
}

func (m *MySQL) Tables() ([]model.Record, error) {
	return m.Query("SELECT table_name FROM information_schema.tables where table_schema=database()")
}

func (m *MySQL) Insert(table string, records []model.RecordInput, params ...string) (int64, error) {
	var count int64

	args, err := internal.DecodeQuery(params)
	if err != nil {
		return count, err
	}

	// append with commandline args
	for _, r := range records {
		for k, v := range args {
			r[k] = v
		}

		result, err := m.insertQueryNamed(table, r)
		if err != nil {
			return count, err
		}

		rowsAffected, _ := result.RowsAffected()
		count += rowsAffected
	}

	log.Printf("Done processing %d rows, %d affected", len(records), count)

	return count, nil
}

func (m *MySQL) insertQueryNamed(table string, data model.RecordInput) (sql.Result, error) {
	template := "INSERT IGNORE INTO %s (%s) VALUES (%s)"

	keys := maps.Keys(data)
	names := strings.Join(keys, ",")
	values := ":" + strings.Join(keys, ", :")

	query := fmt.Sprintf(template, table, names, values)

	return m.db.NamedExec(query, data)
}

func (m *MySQL) Query(sql string, params ...string) ([]model.Record, error) {
	args, err := internal.DecodeQuery(params)
	if err != nil {
		return nil, err
	}

	rows, err := m.db.NamedQuery(sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return internal.ScanAll(rows)
}
