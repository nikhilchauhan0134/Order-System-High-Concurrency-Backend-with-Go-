package services

import (
	"OrderSystemHighConcurrency/order-processor/internal/contracts"
	"OrderSystemHighConcurrency/shared/models"
	"context"
	"fmt"
)

const maxRetries = 3

// RetryService handles retry and DLQ logic
type RetryService struct {
	dlq contracts.DLQPublisher
}

// NewRetryService creates a RetryService
func NewRetryService(dlq contracts.DLQPublisher) *RetryService {
	return &RetryService{dlq: dlq}
}

// HandleFailure decides retry or DLQ
func (r *RetryService) HandleFailure(ctx context.Context, order *models.Order, cause error) error {
	order.RetryCount++

	if order.RetryCount > maxRetries {
		order.Status = models.OrderStatusFailed
		_ = r.dlq.Publish(ctx, order, cause.Error())
		return fmt.Errorf("order sent to DLQ: %s", order.OrderID)
	}

	return cause
}
