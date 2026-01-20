package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	NodeEnv  string
	Port     int
	Bind     string
	APIToken string
}

// Load loads configuration from environment variables and .env files
func Load() (*Config, error) {
	// Load .env files in order of precedence
	// Note: errors are ignored as these files are optional
	homeDir, _ := os.UserHomeDir()
	_ = godotenv.Load(filepath.Join(homeDir, ".env"))
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env.secret")

	cfg := &Config{
		NodeEnv:  getEnv("NODE_ENV", "development"),
		Port:     getEnvInt("PORT", 5227),
		Bind:     getEnv("BIND", "0.0.0.0"),
		APIToken: getAPIToken(),
	}

	return cfg, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

// getAPIToken gets the API token from environment variables
func getAPIToken() string {
	// Check LAAS_API_TOKEN first, then API_TOKEN
	if token := os.Getenv("LAAS_API_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("API_TOKEN")
}
