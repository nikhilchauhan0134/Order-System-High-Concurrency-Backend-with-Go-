package config

import (
	"os"
	"strings"
)

type Config struct {
	GRPCPort     string
	KafkaBrokers []string
	KafkaTopic   string
}

func LoadConfig() *Config {
	return &Config{
		GRPCPort:     getEnv("GRPC_PORT", "50051"),
		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "orders"),
	}
}

func getEnv(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultVal
}
