package domain

import (
	"context"
	"my-go-server/internal/models"
	"net/http"
)

type Handler interface {
	Handle() http.HandlerFunc
	HandleDBTest() http.HandlerFunc
	RegisterUser() http.HandlerFunc
	LoginUser() http.HandlerFunc

	AddNewOrder() http.HandlerFunc
	GetOrdersList() http.HandlerFunc
}

type UseCase interface {
	GetMessage(ctx context.Context) (string, error)
	SaveMessage(ctx context.Context, value string) (int, error) 
	RegisterUser(ctx context.Context, login, password string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	CreateSession(ctx context.Context, userID string) (*models.Session, error)

	CreateOrder(ctx context.Context, userID string) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID string, activeOnly bool) ([]*models.Order, error)
}

type Repository interface {
	GetMessage(ctx context.Context) (string, error)
	SaveMessage(ctx context.Context, value string) (int, error)
	RegisterUser(ctx context.Context, login, password string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	CreateSession(ctx context.Context, userID string) (*models.Session, error)

	CreateOrder(ctx context.Context, userID string) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID string, activeOnly bool) ([]*models.Order, error)
}