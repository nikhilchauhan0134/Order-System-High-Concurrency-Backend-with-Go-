package services

import (
	"context"
	"errors"

	"OrderSystemHighConcurrency/order-processor/internal/contracts"
	"OrderSystemHighConcurrency/shared/models"
)

// processorService implements contracts.OrderProcessor
type processorService struct {
	batchService contracts.BatchService
	retryService contracts.RetryService
	dlq          contracts.DLQPublisher
}

// NewOrderProcessor creates OrderProcessor
func NewOrderProcessor(
	batchService contracts.BatchService,
	retryService contracts.RetryService,
	dlq contracts.DLQPublisher,
) contracts.OrderProcessor {
	return &processorService{
		batchService: batchService,
		retryService: retryService,
		dlq:          dlq,
	}
}

// Process processes a single order
func (p *processorService) Process(ctx context.Context, order *models.Order) error {
	if order == nil {
		return errors.New("order is nil")
	}

	order.Status = models.OrderStatusProcessing

	// Try batching
	err := p.batchService.Add(ctx, order)
	if err != nil {
		// Increment retry count
		order.RetryCount++

		// Check if max retries reached
		if !p.retryService.ShouldRetry(order.RetryCount) {
			// Send to DLQ
			return p.dlq.Publish(ctx, order, err.Error())
		}

		// Return error to retry later
		return err
	}

	order.Status = models.OrderStatusCompleted
	return nil
}
