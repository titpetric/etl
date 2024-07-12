package model

import (
	"time"
)

// Migrations generated for db table `migrations`
type Migrations struct {
	// Project
	Project string `db:"project" json:"project"`

	// Filename
	Filename string `db:"filename" json:"filename"`

	// Statement index
	StatementIndex int64 `db:"statement_index" json:"statement_index"`

	// Status
	Status string `db:"status" json:"status"`
}

// MigrationsTable is the name of the table in the DB
const MigrationsTable = "`migrations`"

// MigrationsFields are all the field names in the DB table
var MigrationsFields = []string{"project", "filename", "statement_index", "status"}

// MigrationsPrimaryFields are the primary key fields in the DB table
var MigrationsPrimaryFields = []string{"project"}

// Commit generated for db table `commit`
type Commit struct {
	// Id
	ID int64 `db:"id" json:"id"`

	// Commit id
	CommitID string `db:"commit_id" json:"commit_id"`

	// Repository
	Repository string `db:"repository" json:"repository"`

	// Created at
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

// SetCreatedAt sets CreatedAt which requires a *time.Time
func (c *Commit) SetCreatedAt(stamp time.Time) { c.CreatedAt = &stamp }

// CommitTable is the name of the table in the DB
const CommitTable = "`commit`"

// CommitFields are all the field names in the DB table
var CommitFields = []string{"id", "commit_id", "repository", "created_at"}

// CommitPrimaryFields are the primary key fields in the DB table
var CommitPrimaryFields = []string{"id"}

// CommitOutput generated for db table `commit_output`
type CommitOutput struct {
	// Id
	ID int64 `db:"id" json:"id"`

	// Commit id
	CommitID int64 `db:"commit_id" json:"commit_id"`

	// Created with
	CreatedWith string `db:"created_with" json:"created_with"`

	// Filename
	Filename string `db:"filename" json:"filename"`

	// Contents
	Contents string `db:"contents" json:"contents"`

	// Created at
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

// SetCreatedAt sets CreatedAt which requires a *time.Time
func (c *CommitOutput) SetCreatedAt(stamp time.Time) { c.CreatedAt = &stamp }

// CommitOutputTable is the name of the table in the DB
const CommitOutputTable = "`commit_output`"

// CommitOutputFields are all the field names in the DB table
var CommitOutputFields = []string{"id", "commit_id", "created_with", "filename", "contents", "created_at"}

// CommitOutputPrimaryFields are the primary key fields in the DB table
var CommitOutputPrimaryFields = []string{"id"}
