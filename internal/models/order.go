package models

import "time"

type Order struct {
	ID string `db:"id" json:"id"`
	UserID string `db:"user_id" json:"user_id"`
	Status string `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateOrderRequest struct {
	Amount int `json:"amount"`
}

type CreateOrderResponse struct {
	OrderIDs []string `json:"order_ids"`
}

type OrdersListResponse struct {
	Count int `json:"count"`
	Items []*Order `json:"items"`
}