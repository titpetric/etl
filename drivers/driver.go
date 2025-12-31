package drivers

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/etl/model"
)

// New creates a driver instance from a database connection.
// The driver type is determined from the database connection's DriverName().
func New(db *sqlx.DB) (model.Driver, error) {
	driver := db.DriverName()
	switch driver {
	case "pgx":
		return NewPgx(driver, db)
	case "mysql":
		return NewMySQL(driver, db)
	case "sqlite":
		return NewSqlite(driver, db)
	default:
		return nil, fmt.Errorf("unknown driver: %s", driver)
	}
}
