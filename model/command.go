package model

import "github.com/jmoiron/sqlx"

// Command is a cli entrypoint for etl commands.
type Command struct {
	Name string
	Args []string

	DB     *sqlx.DB
	Driver string

	Verbose bool
	Quiet   bool
}
