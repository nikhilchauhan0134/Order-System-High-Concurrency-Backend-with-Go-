package contracts

import (
	"OrderSystemHighConcurrency/shared/models"
	"context"
)

type StreamService interface {
	PublishOrder(ctx context.Context, order *models.Order) error
}
