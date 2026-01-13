package models

import "time"

type OrderStatus string

const (
	OrderStatusCreated    OrderStatus = "CREATED"
	OrderStatusQueued     OrderStatus = "QUEUED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusFailed     OrderStatus = "FAILED"
)

type Order struct {
	OrderID    string            `json:"order_id"`
	UserID     string            `json:"user_id"`
	Amount     float64           `json:"amount"`
	Currency   string            `json:"currency"`
	Status     OrderStatus       `json:"status"`
	Source     string            `json:"source"` // web, pos, mobile
	RetryCount int               `json:"retry_count"`
	Metadata   map[string]string `json:"metadata"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type OrderEvent struct {
	OrderID string `json:"order_id"`
	Type    string `json:"type"` // ORDER_CREATED
	Payload Order  `json:"payload"`
}
