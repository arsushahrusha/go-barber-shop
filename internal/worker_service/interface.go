package workerservice

import "context"

type Repository interface {
	StartNewOrdersConsumer(ctx context.Context, handler func(orderID string) error) 
	PublishOrderStatus(ctx context.Context, orderID, status string) error
}

type Service interface {
	Run(ctx context.Context)
}