package contracts

import (
	"OrderSystemHighConcurrency/shared/models"
	"context"
)

type BatchService interface {
	Add(ctx context.Context, order *models.Order) error
	Flush(ctx context.Context) error
}
