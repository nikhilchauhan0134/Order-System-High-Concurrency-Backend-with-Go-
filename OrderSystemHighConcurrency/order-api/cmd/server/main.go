package main

import (
	"OrderSystemHighConcurrency/order-api/internal/infrastructure/kafka"
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
	// 1️ Load configuration (env-based)
	// ------------------------------------------------
	kafkaBrokers := []string{"localhost:9092"}
	kafkaTopic := "orders"

	httpAddr := ":8080"
	http.Handle("/metrics", promhttp.Handler())

	// ------------------------------------------------
	// 2 Initialize Kafka Producer
	// ------------------------------------------------
	producer, err := kafka.NewKafkaProducer(kafkaBrokers, kafkaTopic)
	if err != nil {
		log.Fatalf("failed to init kafka producer: %v", err)
	}

	// ------------------------------------------------
	// 3️ Initialize Order Service
	// ------------------------------------------------
	orderService := services.NewOrderService(producer)

	// ------------------------------------------------
	// 4️ Initialize HTTP Handler
	// ------------------------------------------------
	orderHandler := handlers.NewOrderHandler(orderService)

	// ------------------------------------------------
	// 5️ Rate Limiter Middleware
	// ------------------------------------------------
	rateLimiter := ratelimit.NewIPRateLimiter(100, time.Minute)

	mux := http.NewServeMux()
	mux.Handle("/orders", rateLimiter.Middleware(orderHandler))

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	// ------------------------------------------------
	// 6️ Graceful Shutdown
	// ------------------------------------------------
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Printf("Order API running on %s", httpAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-ctx.Done() // wait for signal
	log.Println("shutting down order-api...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("order-api stopped cleanly")
}
