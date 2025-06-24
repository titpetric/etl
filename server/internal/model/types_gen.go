package etl

import (
	sqlxTypes "github.com/jmoiron/sqlx/types"
)

// GithubCommit generated for db table `github_commit`
//
// GitHub commit details
type GithubCommit struct {
	// Unique identifier
	ID int32 `db:"id" json:"id"`

	// Foreign key to github_org
	OrgID int32 `db:"org_id" json:"org_id"`

	// A copy of github_org name
	OrgName string `db:"org_name" json:"org_name"`

	// Foreign key to github_repository
	RepositoryID int32 `db:"repository_id" json:"repository_id"`

	// Name of github_repository
	RepositoryName string `db:"repository_name" json:"repository_name"`

	// Git commit hash
	CommitSha string `db:"commit_sha" json:"commit_sha"`

	// Commit date
	CommitDate string `db:"commit_date" json:"commit_date"`

	// Commit author
	Author string `db:"author" json:"author"`

	// Commit message
	Message string `db:"message" json:"message"`

	// Commit URL
	URL string `db:"url" json:"url"`

	// Lines added
	LinesAdded int32 `db:"lines_added" json:"lines_added"`

	// Lines removed
	LinesRemoved int32 `db:"lines_removed" json:"lines_removed"`

	// Lines modified
	LinesModified int32 `db:"lines_modified" json:"lines_modified"`
}

// GithubCommitTable is the name of the table in the DB
const GithubCommitTable = "`github_commit`"

// GithubCommitFields are all the field names in the DB table
var GithubCommitFields = []string{"id", "org_id", "org_name", "repository_id", "repository_name", "commit_sha", "commit_date", "author", "message", "url", "lines_added", "lines_removed", "lines_modified"}

// GithubCommitPrimaryFields are the primary key fields in the DB table
var GithubCommitPrimaryFields = []string{"id"}

// GithubOrg generated for db table `github_org`
//
// GitHub organization details
type GithubOrg struct {
	// Unique identifier
	ID int32 `db:"id" json:"id"`

	// Organization name
	Name string `db:"name" json:"name"`

	// Organization URL
	URL string `db:"url" json:"url"`

	// Details JSON
	JSON sqlxTypes.JSONText `db:"json" json:"json"`
}

// GithubOrgTable is the name of the table in the DB
const GithubOrgTable = "`github_org`"

// GithubOrgFields are all the field names in the DB table
var GithubOrgFields = []string{"id", "name", "url", "json"}

// GithubOrgPrimaryFields are the primary key fields in the DB table
var GithubOrgPrimaryFields = []string{"id"}

// GithubRepository generated for db table `github_repository`
//
// GitHub repository details
type GithubRepository struct {
	// Unique identifier
	ID int32 `db:"id" json:"id"`

	// Foreign key to github_org
	OrgID int32 `db:"org_id" json:"org_id"`

	// A copy of github_org name
	OrgName string `db:"org_name" json:"org_name"`

	// Repository flags for processing
	Flags string `db:"flags" json:"flags"`

	// Repository name
	Name string `db:"name" json:"name"`

	// Repository description
	Description string `db:"description" json:"description"`

	// Repository URL
	URL string `db:"url" json:"url"`

	// Repository Go import path
	ImportPath string `db:"import_path" json:"import_path"`

	// Creation date
	CreatedAt string `db:"createdAt" json:"createdAt"`

	// Last update date
	UpdatedAt string `db:"updatedAt" json:"updatedAt"`
}

// GithubRepositoryTable is the name of the table in the DB
const GithubRepositoryTable = "`github_repository`"

// GithubRepositoryFields are all the field names in the DB table
var GithubRepositoryFields = []string{"id", "org_id", "org_name", "flags", "name", "description", "url", "import_path", "createdAt", "updatedAt"}

// GithubRepositoryPrimaryFields are the primary key fields in the DB table
var GithubRepositoryPrimaryFields = []string{"id"}

// GithubTag generated for db table `github_tag`
//
// Github tags on repository
type GithubTag struct {
	// Unique identifier
	ID int32 `db:"id" json:"id"`

	// Foreign key to github_repository
	RepositoryID int32 `db:"repository_id" json:"repository_id"`

	// Name of github_repository
	RepositoryName string `db:"repository_name" json:"repository_name"`

	// Git tag name
	Name string `db:"name" json:"name"`

	// Git commit hash
	CommitSha string `db:"commit_sha" json:"commit_sha"`
}

// GithubTagTable is the name of the table in the DB
const GithubTagTable = "`github_tag`"

// GithubTagFields are all the field names in the DB table
var GithubTagFields = []string{"id", "repository_id", "repository_name", "name", "commit_sha"}

// GithubTagPrimaryFields are the primary key fields in the DB table
var GithubTagPrimaryFields = []string{"id"}

// GoMod generated for db table `go_mod`
//
// Go module details
type GoMod struct {
	// Unique identifier
	ID int32 `db:"id" json:"id"`

	// Foreign key to github_repository
	RepositoryID int32 `db:"repository_id" json:"repository_id"`

	// File name
	Name string `db:"name" json:"name"`

	// Contents of go.mod file
	Content string `db:"content" json:"content"`

	// Timestamp of when the go.mod was fetched
	FetchedAt string `db:"fetchedAt" json:"fetchedAt"`
}

// GoModTable is the name of the table in the DB
const GoModTable = "`go_mod`"

// GoModFields are all the field names in the DB table
var GoModFields = []string{"id", "repository_id", "name", "content", "fetchedAt"}

// GoModPrimaryFields are the primary key fields in the DB table
var GoModPrimaryFields = []string{"id"}

// GoModule generated for db table `go_module`
//
// Go module import details
type GoModule struct {
	// Module path
	Path string `db:"path" json:"path"`

	// Module version
	Version string `db:"version" json:"version"`

	// Source repository
	Source string `db:"source" json:"source"`
}

// GoModuleTable is the name of the table in the DB
const GoModuleTable = "`go_module`"

// GoModuleFields are all the field names in the DB table
var GoModuleFields = []string{"path", "version", "source"}

// GoModulePrimaryFields are the primary key fields in the DB table
var GoModulePrimaryFields = []string{"path", "version", "source"}

// Migrations generated for db table `migrations`
//
// Migration log of applied migrations
type Migrations struct {
	// Microservice or project name
	Project string `db:"project" json:"project"`

	// yyyy-mm-dd-HHMMSS.sql
	Filename string `db:"filename" json:"filename"`

	// Statement number from SQL file
	StatementIndex int32 `db:"statement_index" json:"statement_index"`

	// ok or full error message
	Status string `db:"status" json:"status"`
}

// MigrationsTable is the name of the table in the DB
const MigrationsTable = "`migrations`"

// MigrationsFields are all the field names in the DB table
var MigrationsFields = []string{"project", "filename", "statement_index", "status"}

// MigrationsPrimaryFields are the primary key fields in the DB table
var MigrationsPrimaryFields = []string{"project", "filename"}
