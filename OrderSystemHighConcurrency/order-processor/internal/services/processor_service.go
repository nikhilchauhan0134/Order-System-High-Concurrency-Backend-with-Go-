package services

import (
	"OrderSystemHighConcurrency/order-processor/internal/contracts"
	"OrderSystemHighConcurrency/shared/models"
	"context"
	"errors"
)

// processorService implements OrderProcessor
type processorService struct {
	batchService contracts.Repository
	retryService *RetryService
}

// NewOrderProcessor creates a new OrderProcessor service
func NewOrderProcessor(
	repo contracts.Repository,
	dlq contracts.DLQPublisher,
) contracts.OrderProcessor {
	return &processorService{
		batchService: repo,
		retryService: NewRetryService(dlq),
	}
}

// Process processes a single order
func (p *processorService) Process(ctx context.Context, order *models.Order) error {
	if order == nil {
		return errors.New("order is nil")
	}

	order.Status = models.OrderStatusProcessing

	// Add order to batch
	err := p.batchService.SaveBatch(ctx, []*models.Order{order})
	if err != nil {
		return p.retryService.HandleFailure(ctx, order, err)
	}

	order.Status = models.OrderStatusCompleted
	return nil
}
