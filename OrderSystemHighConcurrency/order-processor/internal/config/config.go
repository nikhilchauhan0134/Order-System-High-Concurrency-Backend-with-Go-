package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Kafka
	KafkaBrokers  []string
	KafkaTopic    string
	ConsumerGroup string

	// DB
	DBDSN string

	// Worker Pool
	WorkerCount int

	// Batch Service
	BatchSize          int
	BatchFlushInterval time.Duration

	// Retry
	MaxRetries int
}

// LoadConfig reads env variables and returns Config
func LoadConfig() *Config {
	cfg := &Config{}

	// Kafka
	cfg.KafkaBrokers = []string{"localhost:9092"} // could parse from ENV if needed
	cfg.KafkaTopic = getEnv("KAFKA_TOPIC", "orders")
	cfg.ConsumerGroup = getEnv("KAFKA_GROUP", "order-processor-group")

	// DB
	cfg.DBDSN = getEnv("DB_DSN", "sqlserver://SA:YourStrong@Passw0rd@localhost:1433?database=Orders")

	// Worker Pool
	cfg.WorkerCount = getEnvAsInt("WORKER_COUNT", 20)

	// Batch Service
	cfg.BatchSize = getEnvAsInt("BATCH_SIZE", 1000)
	cfg.BatchFlushInterval = getEnvAsDuration("BATCH_FLUSH_INTERVAL", 5*time.Second)

	// Retry
	cfg.MaxRetries = getEnvAsInt("MAX_RETRIES", 3)

	return cfg
}

// helper functions
func getEnv(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if valStr, ok := os.LookupEnv(key); ok {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if valStr, ok := os.LookupEnv(key); ok {
		if val, err := time.ParseDuration(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
