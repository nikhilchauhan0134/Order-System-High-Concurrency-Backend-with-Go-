package worker

import (
	"OrderSystemHighConcurrency/order-processor/internal/contracts"
	"OrderSystemHighConcurrency/shared/models"
	"context"
	"log"
	"sync"
)

// WorkerPool controls concurrent order processing
type WorkerPool struct {
	workerCount int
	jobs        chan *models.Order
	processor   contracts.OrderProcessor
	wg          sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workerCount int, bufferSize int, processor contracts.OrderProcessor) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		jobs:        make(chan *models.Order, bufferSize),
		processor:   processor,
	}
}

// Start initializes worker goroutines
func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx, i)
	}
}

// Submit sends an order to the worker pool
func (wp *WorkerPool) Submit(order *models.Order) {
	wp.jobs <- order
}

// worker processes jobs from the channel
func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d shutting down", id)
			return

		case order := <-wp.jobs:
			if order == nil {
				continue
			}

			if err := wp.processor.Process(ctx, order); err != nil {
				log.Printf("worker %d failed to process order %s: %v", id, order.OrderID, err)
			}
		}
	}
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
}
