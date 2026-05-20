package db

import (
	"context"
	"database/sql"
	"fmt"
	"my-go-server/internal/models"
	"time"
	"errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	var userID string
	err = r.db.GetContext(ctx, &userID, registerUserQuery, login, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &models.User{
		ID: userID,
		Login: login,
		Password: string(hashedPassword),
		CreatedAt: time.Now(),
	}, nil
}

func (r *DBRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	user := &models.User{}
	err := r.db.GetContext(ctx, user, getUserByLoginQuery, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
 
func (r *DBRepository) CreateSession(ctx context.Context, userID string) (*models.Session, error) {
	session := &models.Session{}

	_, err := r.db.ExecContext(ctx, createSessionQuery, session, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}


