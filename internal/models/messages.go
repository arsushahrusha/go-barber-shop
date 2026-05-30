package models

type NewOrderMessage struct {
	OrderID string `json:"order_id"`
}

type OrderStatusChangedMessage struct {
	OrderID string `json:"order_id"`
	Status string `json:"status"`
}