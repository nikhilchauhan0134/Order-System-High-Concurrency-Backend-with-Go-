package contracts

import "context"

// Consumer defines a message consumer (Kafka, RabbitMQ, etc.)
type Consumer interface {
	// Start begins consuming messages and pushing them to workers.
	Start(ctx context.Context) error

	// Close shuts down the consumer gracefully.
	Close() error
}
