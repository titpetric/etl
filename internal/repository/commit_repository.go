package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	. "github.com/titpetric/etl/internal/model"
)

// ListOptions holds the options for listing commits.
type ListOptions struct {
	Limit  int
	Offset int
}

// CommitRepository provides access to the `commit` table for transactional operations.
type CommitRepository struct {
	tx *sqlx.Tx
}

// NewCommitRepository creates a new CommitRepository.
func NewCommitRepository(tx *sqlx.Tx) *CommitRepository {
	return &CommitRepository{tx: tx}
}

// Create inserts a new commit into the `commit` table.
func (r *CommitRepository) Create(ctx context.Context, commit *Commit) error {
	commit.SetCreatedAt(time.Now())
	query := sqlInsert(CommitTable, CommitFields[1:]) // exclude "id" field for insert
	_, err := r.tx.NamedExecContext(ctx, query, commit)
	return err
}

// Replace replaces an existing commit in the `commit` table.
func (r *CommitRepository) Replace(ctx context.Context, commit *Commit) error {
	commit.SetCreatedAt(time.Now())
	query := sqlReplace(CommitTable, CommitFields)
	_, err := r.tx.NamedExecContext(ctx, query, commit)
	return err
}

// Update updates an existing commit in the `commit` table.
func (r *CommitRepository) Update(ctx context.Context, commit *Commit) error {
	query := sqlUpdate(CommitTable, CommitFields[1:], CommitPrimaryFields) // exclude "id" field for update
	_, err := r.tx.NamedExecContext(ctx, query, commit)
	return err
}

// CreateOrUpdate creates a new commit or updates an existing one based on the presence of the .ID value.
func (r *CommitRepository) CreateOrUpdate(ctx context.Context, commit *Commit) error {
	if commit.ID == 0 {
		return r.Create(ctx, commit)
	}
	return r.Update(ctx, commit)
}

// Get retrieves a commit by ID from the `commit` table.
func (r *CommitRepository) Get(ctx context.Context, id string) (*Commit, error) {
	var commit Commit
	query := "SELECT * FROM " + CommitTable + " WHERE `id` = ?"
	err := r.tx.GetContext(ctx, &commit, query, id)
	if err != nil {
		return nil, err
	}
	return &commit, nil
}

// List retrieves a list of commits from the `commit` table.
func (r *CommitRepository) List(ctx context.Context, opts ListOptions) ([]*Commit, error) {
	var commits []*Commit
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY `created_at` DESC LIMIT %d OFFSET %d", CommitTable, opts.Limit, opts.Offset)
	err := r.tx.SelectContext(ctx, &commits, query)
	if err != nil {
		return nil, err
	}
	return commits, nil
}
