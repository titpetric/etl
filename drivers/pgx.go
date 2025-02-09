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

type Pgx struct {
	db *sqlx.DB
}

func NewPgx(db *sqlx.DB) (*Pgx, error) {
	return &Pgx{
		db: db,
	}, nil
}

func (d *Pgx) Tables() ([]model.Record, error) {
	return d.Query("SELECT table_name FROM information_schema.tables where table_schema=current_schema()")
}

func (d *Pgx) Insert(table string, records []model.RecordInput, params ...string) (int64, error) {
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

		result, err := d.insertQueryNamed(table, r)
		if err != nil {
			return count, err
		}

		rowsAffected, _ := result.RowsAffected()
		count += rowsAffected
	}

	log.Printf("Done processing %d rows, %d affected", len(records), count)

	return count, nil
}

func (d *Pgx) insertQueryNamed(table string, data model.RecordInput) (sql.Result, error) {
	template := "INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO NOTHING"

	keys := maps.Keys(data)
	names := strings.Join(keys, ",")
	values := ":" + strings.Join(keys, ", :")

	query := fmt.Sprintf(template, table, names, values)

	return d.db.NamedExec(query, data)
}

func (d *Pgx) Query(sql string, params ...string) ([]model.Record, error) {
	args, err := internal.DecodeQuery(params)
	if err != nil {
		return nil, err
	}

	rows, err := d.db.NamedQuery(sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return internal.ScanAll(rows)
}
