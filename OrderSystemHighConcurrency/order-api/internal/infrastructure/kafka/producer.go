package kafka

import (
	"OrderSystemHighConcurrency/order-api/internal/contracts"
	"OrderSystemHighConcurrency/shared/models"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/IBM/sarama"
)

// kafkaProducer implements contracts.Producer
type kafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string, topic string) (contracts.Producer, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers required")
	}

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &kafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

// Publish sends an order to Kafka
func (k *kafkaProducer) Publish(ctx context.Context, order *models.Order) error {
	if order == nil {
		return errors.New("order is nil")
	}

	payload, err := json.Marshal(order)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.ByteEncoder(payload),
	}

	// Respect context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, _, err = k.producer.SendMessage(msg)
		return err
	}
}
