package kafka

import (
	"OrderSystemHighConcurrency/order-processor/internal/contracts"
	"OrderSystemHighConcurrency/order-processor/internal/worker"
	"OrderSystemHighConcurrency/shared/models"
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

// orderConsumer implements contracts.Consumer
type orderConsumer struct {
	consumerGroup sarama.ConsumerGroup
	topic         string
	workerPool    *worker.WorkerPool
}

// NewOrderConsumer creates a new Kafka consumer
func NewOrderConsumer(
	brokers []string,
	groupID string,
	topic string,
	workerPool *worker.WorkerPool,
) (contracts.Consumer, error) {

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Return.Errors = true

	cg, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &orderConsumer{
		consumerGroup: cg,
		topic:         topic,
		workerPool:    workerPool,
	}, nil
}

// Start begins consuming Kafka messages
func (c *orderConsumer) Start(ctx context.Context) error {
	handler := &consumerHandler{
		workerPool: c.workerPool,
	}

	for {
		if err := c.consumerGroup.Consume(ctx, []string{c.topic}, handler); err != nil {
			log.Printf("kafka consume error: %v", err)
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// Close shuts down consumer
func (c *orderConsumer) Close() error {
	return c.consumerGroup.Close()
}

type consumerHandler struct {
	workerPool *worker.WorkerPool
}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for msg := range claim.Messages() {
		var order models.Order

		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("failed to unmarshal order: %v", err)
			session.MarkMessage(msg, "")
			continue
		}

		// Send order to worker pool (async)
		h.workerPool.Submit(&order)

		// Mark message as consumed
		session.MarkMessage(msg, "")
	}

	return nil
}
