package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	CosmosNodeAddress string
	Port              string
	MaxRetries        int
	RetryDelay        time.Duration
}

// New creates and returns a new Config instance
func New() *Config {
	return &Config{
		CosmosNodeAddress: getEnv("COSMOS_NODE_ADDRESS", "http://localhost:26657"),
		Port:              getEnv("PORT", "8080"),
		MaxRetries:        getEnvAsInt("MAX_RETRIES", 5),
		RetryDelay:        getEnvAsDuration("RETRY_DELAY", 5*time.Second),
	}
}

// getEnv retrieves the value of the environment variable named by the key
// If the variable is not present, it returns the fallback value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsInt retrieves the value of the environment variable as an integer
func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// getEnvAsDuration retrieves the value of the environment variable as a duration
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}
