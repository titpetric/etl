package sql

import (
	handlermodel "github.com/titpetric/etl/server/internal/handler/model"
)

var (
	dbValue = handlermodel.DBValue
)

func init() {
	handlermodel.Register(NewHandler())
}
