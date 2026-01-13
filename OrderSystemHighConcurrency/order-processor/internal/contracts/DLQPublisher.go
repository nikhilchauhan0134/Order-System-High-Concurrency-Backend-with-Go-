package contracts

import (
	"OrderSystemHighConcurrency/shared/models"
	"context"
)

// DLQPublisher defines how failed orders are sent to a dead letter queue.
type DLQPublisher interface {
	// Publish sends the failed order to DLQ with a reason.
	Publish(ctx context.Context, order *models.Order, reason string) error
}
