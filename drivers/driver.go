package drivers

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/etl/model"
)

func New(driver string, db *sqlx.DB) (model.Driver, error) {
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
