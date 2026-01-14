package dlq

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"OrderSystemHighConcurrency/order-processor/internal/contracts"
	"OrderSystemHighConcurrency/shared/models"

	"github.com/IBM/sarama"
)

type dlqProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewDLQProducer creates a new DLQ producer
func NewDLQProducer(brokers []string, topic string) (contracts.DLQPublisher, error) {
	if len(brokers) == 0 {
		return nil, errors.New("brokers required")
	}

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_8_0_0
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Return.Successes = true
	cfg.Producer.Timeout = 5 * time.Second

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return &dlqProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

// Publish sends failed order to DLQ
func (d *dlqProducer) Publish(
	ctx context.Context,
	order *models.Order,
	reason string,
) error {
	if order == nil {
		return errors.New("order is nil")
	}

	payload, err := json.Marshal(struct {
		Order  *models.Order `json:"order"`
		Reason string        `json:"reason"`
		Time   time.Time     `json:"time"`
	}{
		Order:  order,
		Reason: reason,
		Time:   time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: d.topic,
		Value: sarama.ByteEncoder(payload),
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, _, err := d.producer.SendMessage(msg)
		return err
	}
}
