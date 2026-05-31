package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"my-go-server/internal/models"
	"my-go-server/internal/rabbitmq"
)

type RabbitRepository struct {
	rabbit *rabbitmq.Client
	newOrderQueue string
	statusQueue string
}

func NewRabbitRepository(
	rabbit *rabbitmq.Client,
	newOrdersQueue string,
	statusQueue string,
) *RabbitRepository {
	return &RabbitRepository{
		rabbit: rabbit,
		newOrderQueue: newOrdersQueue,
		statusQueue: statusQueue,
	}
}

func (r *RabbitRepository) StartNewOrdersConsumer(
	ctx context.Context,
	handler func(orderID string) error,
) {
	deliveries, err := r.rabbit.Consume(r.newOrderQueue)
	if err != nil {
		fmt.Printf("failed to start new orders consumer: %v\n", err)
		return
	}

	go func() {
		for {
			select {
			case <- ctx.Done():
				fmt.Println("new orders consumer stopped")
				return
			
			case msg, ok := <- deliveries:
				if !ok {
					fmt.Println("new orders channel closed")
					return
				}

				var event models.NewOrderMessage
				if err := json.Unmarshal(msg.Body, &event); err != nil {
					fmt.Printf("invalid new order message: %v\n", err)
					_ = msg.Nack(false, false)
					continue
				}

				if err := handler(event.OrderID); err != nil {
					fmt.Printf("failed to process order %s: %v\n", event.OrderID, err)
					_ = msg.Nack(false, true)
					continue
				}

				_ = msg.Ack(false)
			}
		}
	} ()
}

func (r *RabbitRepository) PublishOrderStatus(ctx context.Context, orderID, status string) error {
	msg := models.OrderStatusChangedMessage{
		OrderID: orderID,
		Status: status,
	}

	body, err := json.Marshal(msg)

	if err != nil {
		return fmt.Errorf("failed to marshal new order message: %w", err)
	}

	if err := r.rabbit.Publish(ctx, r.statusQueue, body); err != nil {
		return fmt.Errorf("failed to publish new order: %w", err)
	}

	return nil
}