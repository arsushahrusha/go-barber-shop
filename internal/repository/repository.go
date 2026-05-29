package repository

import (
	"context"
	"my-go-server/internal/models"
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

func (r *Repository) RegisterUser(ctx context.Context, login, password string) (*models.User, error) {
	return r.db.RegisterUser(ctx, login, password)
}

func (r *Repository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	return r.db.GetUserByLogin(ctx, login)
}

func (r *Repository) CreateSession(ctx context.Context, userID string) (*models.Session, error) {
	return r.db.CreateSession(ctx, userID)
}

func (r *Repository) CreateOrder(ctx context.Context, userID string) (*models.Order, error) {
	return r.db.CreateOrder(ctx, userID)
}

func (r *Repository) GetOrdersByUserID(ctx context.Context, userID string, activeOnly bool) ([]*models.Order, error) {
	return r.db.GetOrdersByUserID(ctx, userID, activeOnly)
}

func (r *Repository) GetSessionByUserID(ctx context.Context, userID string) (*models.Session, error) {
	return r.db.GetSessionsByUserID(ctx, userID)
}

func (r *Repository) UpdateSessionExpiry(ctx context.Context, sessionID string) (error) {
	return r.db.UpdateSessionExpiry(ctx, sessionID)
}