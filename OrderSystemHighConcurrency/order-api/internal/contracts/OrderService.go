package contracts

import (
	"OrderSystemHighConcurrency/shared/models"
	"context"
)

type OrderService interface {
	// CreateOrder validates and creates a new order.
	// It sends the order to a message queue for processing.
	CreateOrder(ctx context.Context, order *models.Order) error

	// Optionally, you can add more future operations:
	// GetOrderStatus(ctx context.Context, orderID string) (models.OrderStatus, error)
}
