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

	GetSessionByUserID(ctx context.Context, userID string) (*models.Session, error)
	UpdateSessionExpiry(ctx context.Context, sessionID string) error

	PublishNewOrder(ctx context.Context, orderId string) error
	ChangeOrderStatus(ctx context.Context, orderID string, status string) error
	StartStatusChangeConsumer(ctx context.Context)
}

type Repository interface {
	GetMessage(ctx context.Context) (string, error)
	SaveMessage(ctx context.Context, value string) (int, error)
	RegisterUser(ctx context.Context, login, password string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	CreateSession(ctx context.Context, userID string) (*models.Session, error)

	CreateOrder(ctx context.Context, userID string) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID string, activeOnly bool) ([]*models.Order, error)

	GetSessionByUserID(ctx context.Context, userID string) (*models.Session, error)
	UpdateSessionExpiry(ctx context.Context, sessionID string) error

	PublishNewOrder(ctx context.Context, orderId string) error
	ChangeOrderStatus(ctx context.Context, orderID string, status string) error
	StartStatusChangeConsumer(ctx context.Context, handler func(orderID, status string) error)
}