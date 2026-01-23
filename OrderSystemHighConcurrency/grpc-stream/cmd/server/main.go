package main

import (
	"OrderSystemHighConcurrency/grpc-stream/internal/config"
	servicescontract "OrderSystemHighConcurrency/grpc-stream/internal/contracts"
	"OrderSystemHighConcurrency/grpc-stream/internal/services"
	sharedkafa "OrderSystemHighConcurrency/shared/kafka"
	"OrderSystemHighConcurrency/shared/models"
	"context"

	pb "OrderSystemHighConcurrency/grpc-stream/internal/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedOrderStreamServer
	streamService servicescontract.StreamService
}

func (s *server) StreamOrders(stream pb.OrderStream_StreamOrdersServer) error {
	for {
		orderProto, err := stream.Recv()
		if err != nil {
			log.Printf("stream closed: %v", err)
			return err
		}

		order := &models.Order{
			OrderID:   orderProto.Id,
			Amount:    orderProto.Amount,
			UserID:    orderProto.CustomerId,
			CreatedAt: orderProto.CreatedAt.AsTime(),
		}

		_ = s.streamService.PublishOrder(context.Background(), order)

		if err := stream.Send(&pb.StreamResponse{Status: "received"}); err != nil {
			log.Printf("failed to send response: %v", err)
		}
	}
}

func main() {
	cfg := config.LoadConfig()

	producer, err := sharedkafa.NewKafkaProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("failed to init kafka producer: %v", err)
	}

	streamService := services.NewStreamService(producer)

	grpcServer := grpc.NewServer()
	pb.RegisterOrderStreamServer(grpcServer, &server{streamService: streamService})

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on port %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
