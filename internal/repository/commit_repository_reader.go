package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	. "github.com/titpetric/etl/internal/model"
)

// CommitRepositoryReader provides access to the `commit` table for read-only operations.
type CommitRepositoryReader struct {
	db *sqlx.DB
}

// NewCommitRepositoryReader creates a new CommitRepositoryReader.
func NewCommitRepositoryReader(db *sqlx.DB) *CommitRepositoryReader {
	return &CommitRepositoryReader{db: db}
}

// Get retrieves a commit by ID from the `commit` table.
func (r *CommitRepositoryReader) Get(ctx context.Context, id string) (*Commit, error) {
	var commit Commit
	query := "SELECT * FROM " + CommitTable + " WHERE `id` = ?"
	err := r.db.GetContext(ctx, &commit, query, id)
	if err != nil {
		return nil, err
	}
	return &commit, nil
}

// List retrieves a list of commits from the `commit` table.
func (r *CommitRepositoryReader) List(ctx context.Context, opts ListOptions) ([]*Commit, error) {
	var commits []*Commit
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY `created_at` DESC LIMIT %d OFFSET %d", CommitTable, opts.Limit, opts.Offset)
	err := r.db.SelectContext(ctx, &commits, query)
	if err != nil {
		return nil, err
	}
	return commits, nil
}
