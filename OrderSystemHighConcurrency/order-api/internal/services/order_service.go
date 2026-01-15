package services

import (
	"OrderSystemHighConcurrency/order-api/internal/contracts"
	sharedkafa "OrderSystemHighConcurrency/shared/contracts"
	"OrderSystemHighConcurrency/shared/models"
	"context"
	"errors"
	"time"
)

// orderService implements OrderService contract
type orderService struct {
	producer sharedkafa.Producer
}

// NewOrderService creates a new OrderService
func NewOrderService(producer sharedkafa.Producer) contracts.OrderService {
	return &orderService{
		producer: producer,
	}
}

// CreateOrder handles order creation business logic
func (s *orderService) CreateOrder(ctx context.Context, order *models.Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	// Generate OrderID if not present (idempotency support)
	if order.OrderID == "" {
		return errors.New("amount must be greater than zero")
	}

	// Basic validation
	if order.UserID == "" {
		return errors.New("user_id is required")
	}
	if order.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	// Set initial order state
	now := time.Now().UTC()
	order.Status = models.OrderStatusQueued
	order.RetryCount = 0
	order.CreatedAt = now
	order.UpdatedAt = now

	// Publish order to message queue (Kafka via Producer)
	if err := s.producer.Publish(ctx, order); err != nil {
		return err
	}

	return nil
}
