package db

import (
	"context"
	"database/sql"
	"fmt"
	"my-go-server/internal/models"
	"errors"
	"github.com/jmoiron/sqlx"
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

func (r *DBRepository) RegisterUser(ctx context.Context, login, password string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, registerUserQuery, login, password)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &user, nil
}

func (r *DBRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, getUserByLoginQuery, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
 
func (r *DBRepository) CreateSession(ctx context.Context, userID string) (*models.Session, error) {
	var session models.Session

	err := r.db.GetContext(ctx, &session, createSessionQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &session, nil
}

func (r *DBRepository) CreateOrder(ctx context.Context, userID string) (*models.Order, error) {
	var order models.Order

	err := r.db.GetContext(ctx, &order, createOrderQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return &order, nil
}

func (r *DBRepository) GetOrdersByUserID(ctx context.Context, userID string, activeOnly bool) ([]*models.Order, error) {
	var orders []*models.Order

	query := getOrdersByUserIDQuery

	if activeOnly {
		query = getActiveOrdersByUserIDQuery
	}

	err := r.db.GetContext(ctx, &orders, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return orders, nil
}


