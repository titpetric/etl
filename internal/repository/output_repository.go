package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/etl/internal/model"
)

// OutputRepository handles the interactions with the commit_output table
type OutputRepository struct {
	DB *sqlx.DB
}

// NewOutputRepository creates a new instance of OutputRepository
func NewOutputRepository(db *sqlx.DB) *OutputRepository {
	return &OutputRepository{DB: db}
}

// ListByCommitID returns all the commit outputs for a given commit ID from the commit_output table
func (r *OutputRepository) ListByCommitID(ctx context.Context, commitID int64) ([]*model.CommitOutput, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE commit_id = ?", model.CommitOutputTable)

	var results []*model.CommitOutput
	err := r.DB.SelectContext(ctx, &results, query, commitID)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Create inserts a new commit output into the commit_output table
func (r *OutputRepository) Create(ctx context.Context, commitOutput *model.CommitOutput) error {
	commitOutput.SetCreatedAt(time.Now())

	query := sqlInsert(model.CommitOutputTable, model.CommitOutputFields[1:]) // Exclude the "id" field

	_, err := r.DB.NamedExecContext(ctx, query, commitOutput)
	return err
}
