package test

type MessageRepository interface {
	GetMessage() string
}

type MessageService interface {
	GetMessage() string
}