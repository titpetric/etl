package model

type Driver interface {
	Tables() ([]Record, error)
	Query(sql string, params ...string) ([]Record, error)
	Insert(table string, data []RecordInput, params ...string) (int64, error)
}
