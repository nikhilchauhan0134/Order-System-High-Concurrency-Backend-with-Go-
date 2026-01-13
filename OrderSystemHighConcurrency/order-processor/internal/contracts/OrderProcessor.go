package contracts

import (
	"OrderSystemHighConcurrency/shared/models"
	"context"
)

// OrderProcessor defines how an order is processed.
type OrderProcessor interface {
	// Process handles a single order.
	// It may trigger retries, batching, or DLQ routing.
	Process(ctx context.Context, order *models.Order) error
}
