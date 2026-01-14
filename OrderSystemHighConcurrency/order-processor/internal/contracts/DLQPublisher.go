package contracts

import (
	"OrderSystemHighConcurrency/shared/models"
	"context"
)

// DLQPublisher defines how failed orders are sent to a dead letter queue.
type DLQPublisher interface {
	Publish(ctx context.Context, order *models.Order, reason string) error
}
