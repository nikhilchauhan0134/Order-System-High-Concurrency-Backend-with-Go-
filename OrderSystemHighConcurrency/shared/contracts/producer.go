package contracts

import (
	"OrderSystemHighConcurrency/shared/models"
	"context"
)

type Producer interface {
	Publish(ctx context.Context, order *models.Order) error
	Close() error
}
