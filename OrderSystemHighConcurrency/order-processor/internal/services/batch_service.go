package services

import (
	"OrderSystemHighConcurrency/order-processor/internal/contracts"
	"OrderSystemHighConcurrency/shared/models"
	"context"
	"sync"
	"time"
)

// BatchService handles order batching
type BatchService struct {
	repo      contracts.Repository
	batchSize int
	timeout   time.Duration

	mu     sync.Mutex
	buffer []*models.Order
}

// NewBatchService creates a batch service
func NewBatchService(repo contracts.Repository, size int, timeout time.Duration) *BatchService {
	return &BatchService{
		repo:      repo,
		batchSize: size,
		timeout:   timeout,
		buffer:    make([]*models.Order, 0),
	}
}

// Add adds order to batch and flushes if needed
func (b *BatchService) Add(ctx context.Context, order *models.Order) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buffer = append(b.buffer, order)

	if len(b.buffer) >= b.batchSize {
		return b.flush(ctx)
	}

	return nil
}

// flush writes batch to repository
func (b *BatchService) flush(ctx context.Context) error {
	if len(b.buffer) == 0 {
		return nil
	}

	err := b.repo.SaveBatch(ctx, b.buffer)
	b.buffer = make([]*models.Order, 0)
	return err
}
