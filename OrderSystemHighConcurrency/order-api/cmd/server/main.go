package main

import (
	"OrderSystemHighConcurrency/order-api/internal/config"
	"OrderSystemHighConcurrency/order-api/internal/handlers"
	"OrderSystemHighConcurrency/shared/kafka"

	"OrderSystemHighConcurrency/order-api/internal/infrastructure/ratelimit"
	"OrderSystemHighConcurrency/order-api/internal/services"

	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// ------------------------------------------------
	// 1️⃣ Load configuration (env-based)
	// ------------------------------------------------
	cfg := config.LoadConfig()
	kafkaBrokers := cfg.KafkaBrokers
	kafkaTopic := cfg.KafkaTopic

	// ------------------------------------------------
	// 2️⃣ Initialize Kafka Producer
	// ------------------------------------------------
	producer, err := kafka.NewKafkaProducer(kafkaBrokers, kafkaTopic)
	if err != nil {
		log.Fatalf("failed to init Kafka producer: %v", err)
	}
	defer producer.Close() // close producer on shutdown

	// ------------------------------------------------
	// 3️⃣ Initialize Order Service
	// ------------------------------------------------
	orderService := services.NewOrderService(producer)

	// ------------------------------------------------
	// 4️⃣ Initialize HTTP Handler
	// ------------------------------------------------
	orderHandler := handlers.NewOrderHandler(orderService)

	// ------------------------------------------------
	// 5️⃣ Rate Limiter Middleware
	// ------------------------------------------------
	rateLimiter := ratelimit.NewIPRateLimiter(100, time.Minute)

	mux := http.NewServeMux()
	mux.Handle("/orders", rateLimiter.Middleware(orderHandler))
	mux.Handle("/metrics", promhttp.Handler())

	// ------------------------------------------------
	// 6️⃣ HTTP Server
	// ------------------------------------------------
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// ------------------------------------------------
	// 7️⃣ Graceful Shutdown
	// ------------------------------------------------
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Printf("Order API running on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done() // wait for termination signal
	log.Println("shutting down order-api...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	} else {
		log.Println("order-api stopped cleanly")
	}
}
