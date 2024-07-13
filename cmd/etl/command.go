package main

import "github.com/jmoiron/sqlx"

type Command struct {
	Key    string
	SubKey string
	Args   []string
	DB     *sqlx.DB
}
