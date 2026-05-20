package db

import (
	"github.com/jmoiron/sqlx"
	"context"
)

type DBRepository struct {
	db *sqlx.DB
}

func NewDBRepository(db *sqlx.DB) *DBRepository {
	return &DBRepository{db: db}
}

func (r *DBRepository) SaveMessage(ctx context.Context, value string) (int, error) {
	var id int
	query := "INSERT INTO db_test (value) values ($1) RETURNING id"
	
	row := r.db.QueryRowContext(ctx, query, value)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}