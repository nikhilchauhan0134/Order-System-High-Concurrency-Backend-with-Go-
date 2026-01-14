package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"OrderSystemHighConcurrency/order-processor/internal/config"
	"OrderSystemHighConcurrency/order-processor/internal/infrastructure/db"
	"OrderSystemHighConcurrency/order-processor/internal/infrastructure/dlq"
	"OrderSystemHighConcurrency/order-processor/internal/infrastructure/kafka"
	"OrderSystemHighConcurrency/order-processor/internal/services"
	"OrderSystemHighConcurrency/order-processor/internal/worker"
)

func main() {
	// ------------------------------------------------
	// 1️⃣ Load Configuration
	// ------------------------------------------------
	cfg := config.LoadConfig()

	// ------------------------------------------------
	// 2️⃣ Context & Graceful Shutdown
	// ------------------------------------------------
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// ------------------------------------------------
	// 3️⃣ Database Connection (SQL Server)
	// ------------------------------------------------
	database, err := db.NewDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}
	defer database.Close()

	// ------------------------------------------------
	// 4️⃣ Repository
	// ------------------------------------------------
	repository := db.NewOrderRepository(database)

	// ------------------------------------------------
	// 5️⃣ Batch Service
	// ------------------------------------------------
	batchService := services.NewBatchService(
		repository,
		cfg.BatchSize,
		cfg.BatchFlushInterval,
	)

	// ------------------------------------------------
	// 6️⃣ Retry Service
	// ------------------------------------------------
	retryService := services.NewRetryService(cfg.MaxRetries)

	// ------------------------------------------------
	// 7️⃣ DLQ Producer
	// ------------------------------------------------
	dlqPublisher, err := dlq.NewDLQProducer(
		cfg.KafkaBrokers,
		"orders-dlq",
	)
	if err != nil {
		log.Fatalf("failed to init DLQ producer: %v", err)
	}

	// ------------------------------------------------
	// 8️⃣ Order Processor (Core Logic)
	// ------------------------------------------------
	orderProcessor := services.NewOrderProcessor(
		batchService,
		retryService,
		dlqPublisher,
	)

	// ------------------------------------------------
	// 9️⃣ Worker Pool
	// ------------------------------------------------
	queueSize := 1000 // or any number of pending orders you want to buffer

	workerPool := worker.NewWorkerPool(cfg.WorkerCount, queueSize, orderProcessor)
	workerPool.Start(ctx)
	defer workerPool.Stop()

	// ------------------------------------------------
	// 🔟 Kafka Consumer
	// ------------------------------------------------
	consumer, err := kafka.NewOrderConsumer(
		cfg.KafkaBrokers,
		cfg.ConsumerGroup,
		cfg.KafkaTopic,
		workerPool,
	)
	if err != nil {
		log.Fatalf("failed to init kafka consumer: %v", err)
	}
	defer consumer.Close()

	// ------------------------------------------------
	// 1️⃣1️⃣ Start Consumer
	// ------------------------------------------------
	go func() {
		log.Println("order-processor started")
		if err := consumer.Start(ctx); err != nil {
			log.Printf("consumer stopped: %v", err)
			stop()
		}
	}()

	// ------------------------------------------------
	// 1️⃣2️⃣ Wait for shutdown
	// ------------------------------------------------
	<-ctx.Done()
	log.Println("shutting down order-processor...")

	// Allow final batch flush
	time.Sleep(2 * time.Second)

	log.Println("order-processor stopped cleanly")
}
