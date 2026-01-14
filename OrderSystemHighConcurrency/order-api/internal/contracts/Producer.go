package contracts

import (
	"OrderSystemHighConcurrency/shared/models"
	"context"
)

// Producer defines the contract for publishing orders to a message queue.
type Producer interface {
	// Publish sends the order to the message queue (Kafka, RabbitMQ, etc.)
	Publish(ctx context.Context, order *models.Order) error
	Close() error
}
