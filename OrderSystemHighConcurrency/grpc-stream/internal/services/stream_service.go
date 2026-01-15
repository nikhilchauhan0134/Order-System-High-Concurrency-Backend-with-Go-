package services

import (
	streamContracts "OrderSystemHighConcurrency/grpc-stream/internal/contracts"
	sharedContracts "OrderSystemHighConcurrency/shared/contracts"
	"OrderSystemHighConcurrency/shared/models"
	"context"
	"log"
)

type streamService struct {
	producer sharedContracts.Producer
}

func NewStreamService(producer sharedContracts.Producer) streamContracts.StreamService {
	return &streamService{producer: producer}
}

func (s *streamService) PublishOrder(ctx context.Context, order *models.Order) error {
	if order == nil {
		return nil
	}

	if err := s.producer.Publish(ctx, order); err != nil {
		log.Printf("[ERROR] failed to publish order: %v", err)
		return err
	}

	log.Printf("[INFO] order %s published to Kafka via gRPC", order.OrderID)
	return nil
}
