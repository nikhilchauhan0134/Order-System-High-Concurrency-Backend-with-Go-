package contracts

import (
	"OrderSystemHighConcurrency/shared/models"
	"context"
)

// Repository defines database operations for orders.
type Repository interface {
	// SaveBatch persists a batch of orders in the database.
	// Implementations must ensure atomicity and performance.
	SaveBatch(ctx context.Context, orders []*models.Order) error
}
