package domain

import (
	"context"
	"net/http"
)

// type DBRepositoryMessage interface {
// 	InitTable() error
// 	Save(value string) (int, error)
// }

// type DBServiceMessage interface {
// 	Save(value string) (int, error)
// }

// type MessageRepository interface {
// 	GetMessage() string
// }

// type MessageService interface {
// 	GetMessage() string
// }

type Handler interface {
	Handle() http.HandlerFunc
	HandleDBTest() http.HandlerFunc
}

type UseCase interface {
	GetMessage(ctx context.Context) (string, error)
	SaveMessage(ctx context.Context, value string) (int, error) 
}

type Repository interface {
	GetMessage(ctx context.Context) (string, error)
	SaveMessage(ctx context.Context, value string) (int, error)
}