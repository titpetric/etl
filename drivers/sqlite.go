package drivers

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/etl/internal"
	"github.com/titpetric/etl/model"
)

// Sqlite represents a SQLite driver using sqlx.
type Sqlite struct {
	db *sqlx.DB
}

// NewSqlite creates a new SQLite driver instance using the provided *sqlx.DB.
func NewSqlite(db *sqlx.DB) (*Sqlite, error) {
	return &Sqlite{
		db: db,
	}, nil
}

// Insert inserts the given slice of model.RecordInput into the specified table.
// It decodes any additional parameters via internal.DecodeQuery, merges them into each record,
// and performs an INSERT using a named query.
func (s *Sqlite) Insert(table string, records []model.RecordInput, params ...string) (int64, error) {
	var count int64

	// Decode any additional parameters into a map.
	args, err := internal.DecodeQuery(params)
	if err != nil {
		return count, err
	}

	// Merge the additional parameters into each record and insert.
	for _, r := range records {
		// Merge command-line args.
		for k, v := range args {
			r[k] = v
		}

		// Insert using a helper that builds a named query.
		result, err := s.insertQueryNamed(table, r)
		if err != nil {
			return count, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return count, err
		}
		count += rowsAffected
	}

	log.Printf("Done processing %d rows, %d affected", len(records), count)
	return count, nil
}

// insertQueryNamed constructs an INSERT statement with named parameters based on the keys in record
// and executes it using sqlx.NamedExec.
func (s *Sqlite) insertQueryNamed(table string, record model.RecordInput) (sql.Result, error) {
	var columns []string
	for col := range record {
		columns = append(columns, col)
	}
	// Sort the columns to ensure a consistent order.
	sort.Strings(columns)

	var placeholders []string
	for _, col := range columns {
		// Named placeholders in sqlx are prefixed with ":".
		placeholders = append(placeholders, ":"+col)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)
	return s.db.NamedExec(query, record)
}

// Query executes the provided SQL query using named parameters (decoded via internal.DecodeQuery)
// and returns the results as a slice of model.Record.
func (s *Sqlite) Query(sqlQuery string, params ...string) ([]model.Record, error) {
	args, err := internal.DecodeQuery(params)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.NamedQuery(sqlQuery, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return internal.ScanAll(rows)
}

// Tables returns the list of user-defined tables in the SQLite database.
// It queries the sqlite_master table and filters out internal tables.
func (s *Sqlite) Tables() ([]model.Record, error) {
	// The returned column is aliased as "table_name" for consistency.
	return s.Query("SELECT name as table_name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
}
