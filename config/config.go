package config

import (
	"ethereum-parser/shared"
	"os"
)

type Config = shared.Config

func LoadConfig() Config {
	return shared.Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		RPCURL:     getEnv("RPC_URL", "https://cloudflare-eth.com"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
