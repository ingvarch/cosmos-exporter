package config

import "os"

// Config holds the application configuration
type Config struct {
	CosmosNodeAddress string
	Port              string
}

// New creates and returns a new Config instance
func New() *Config {
	return &Config{
		CosmosNodeAddress: getEnv("COSMOS_NODE_ADDRESS", "http://localhost:26657"),
		Port:              getEnv("PORT", "8080"),
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
