package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"my-go-server/internal/models"
	"my-go-server/internal/rabbitmq"
	dbrepo "my-go-server/internal/repository/db"
)

type Repository struct{
	db *dbrepo.DBRepository

	rabbitmq *rabbitmq.Client
	newOrderQueue string
	statusQueue string
}

func NewRepository(
	db *dbrepo.DBRepository,
	rabbitmq *rabbitmq.Client,
	newOrderQueue string,
	statusQueue string,
	) *Repository {
	return &Repository{
		db: db,
		rabbitmq: rabbitmq,
		newOrderQueue: newOrderQueue,
		statusQueue: statusQueue,
	}
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

func (r *Repository) PublishNewOrder(ctx context.Context, orderID string) error {
	msg := models.NewOrderMessage{
		OrderID: orderID,
	}

	body, err := json.Marshal(msg)

	if err != nil {
		return fmt.Errorf("failed to marshal new order message: %w", err)
	}

	if err := r.rabbitmq.Publish(ctx, r.newOrderQueue, body); err != nil {
		return fmt.Errorf("failed to publish new order: %w", err)
	}

	return nil
}

func (r *Repository) StartStatusChangeConsumer(
	ctx context.Context,
	handler func(orderID, status string) error,
) {
	deliveries, err := r.rabbitmq.Consume(r.statusQueue)
	if err != nil {
		fmt.Printf("failed to start status consumer: %v", err)
		return
	}

	go func ()  {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("status consumer stopped")
				return

			case msg, ok := <-deliveries:
				if !ok {
					fmt.Println("status consumer channel closed")
					return
				}

				var event models.OrderStatusChangedMessage

				if err := json.Unmarshal(msg.Body, &event); err != nil {
					fmt.Printf("invalid status message: %v\n", err)
					_ = msg.Nack(false, false) // notifies thar msg should be repeated
					continue
				}

				if err := handler(event.OrderID, event.Status); err != nil {
					fmt.Printf("failed to handle status change: %v\n", err)
					_ = msg.Nack(false, true)
					continue
				}

				_ = msg.Ack(false)
			}
		}
	}()
}

func (r *Repository) ChangeOrderStatus(ctx context.Context, orderID, status string) error {
	return r.db.ChangeOrderStatus(ctx, orderID, status)
}