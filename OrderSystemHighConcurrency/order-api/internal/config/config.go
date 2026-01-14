package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configurable settings for order-api
type Config struct {
	// HTTP Server
	HTTPPort string

	// Kafka
	KafkaBrokers []string
	KafkaTopic   string

	// Rate Limiter
	RateLimitRequests int           // requests
	RateLimitInterval time.Duration // per interval
}

// LoadConfig loads configuration from environment variables or defaults
func LoadConfig() *Config {
	cfg := &Config{}

	// HTTP Port
	cfg.HTTPPort = getEnv("HTTP_PORT", "8080")

	// Kafka Brokers (comma-separated)
	brokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	cfg.KafkaBrokers = splitAndTrim(brokers, ",")

	// Kafka Topic
	cfg.KafkaTopic = getEnv("KAFKA_TOPIC", "orders")

	// Rate Limiter
	cfg.RateLimitRequests = getEnvAsInt("RATE_LIMIT_REQUESTS", 100)                // default 100 requests
	cfg.RateLimitInterval = getEnvAsDuration("RATE_LIMIT_INTERVAL", 1*time.Minute) // default 1 min

	return cfg
}

// --------------- Helpers ----------------
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

func splitAndTrim(str, sep string) []string {
	var res []string
	parts := split(str, sep)
	for _, p := range parts {
		if p != "" {
			res = append(res, trim(p))
		}
	}
	return res
}

// Use built-in functions to avoid extra imports
func split(str, sep string) []string {
	var result []string
	curr := ""
	for i := 0; i < len(str); i++ {
		if string(str[i]) == sep {
			result = append(result, curr)
			curr = ""
		} else {
			curr += string(str[i])
		}
	}
	result = append(result, curr)
	return result
}

func trim(str string) string {
	start, end := 0, len(str)-1
	for start <= end && (str[start] == ' ' || str[start] == '\t') {
		start++
	}
	for end >= start && (str[end] == ' ' || str[end] == '\t') {
		end--
	}
	if start > end {
		return ""
	}
	return str[start : end+1]
}
