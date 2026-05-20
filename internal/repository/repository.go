package repository

import (
	"context"
	dbrepo "my-go-server/internal/repository/db"
)

type Repository struct{
	db *dbrepo.DBRepository
}

func NewRepository(db *dbrepo.DBRepository) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetMessage(ctx context.Context) (string, error) {
		return "Hello!", nil
	}

func (r *Repository) SaveMessage(ctx context.Context, value string) (int, error) {
	return r.db.SaveMessage(ctx, value)
}
	