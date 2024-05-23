package config

import (
	"ethereum-parser/shared"
	"os"
)

// Config type alias from shared package
type Config = shared.Config

// LoadConfig loads the configuration from environment variables with defaults
func LoadConfig() Config {
	return shared.Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		RPCURL:     getEnv("RPC_URL", "https://cloudflare-eth.com"),
	}
}

// getEnv retrieves environment variable value or returns a default value if not set
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
