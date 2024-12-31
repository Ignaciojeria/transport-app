package configuration

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

var (
	once sync.Once
)

// Load environment variables only once
func loadEnvOnce() {
	once.Do(func() {
		// Load variables from the .env file
		if err := godotenv.Load(); err != nil {
			slog.Warn(".env not found, loading environment variables from system.")
		} else {
			slog.Info("Environment variables loaded from .env file.")
		}
	})
}

// Generic method to load configurations
func Parse[T any]() (T, error) {
	// Ensure environment variables are loaded only once
	loadEnvOnce()
	var conf T
	// Use the go-env library to map environment variables to the struct
	if err := env.Parse(&conf); err != nil {
		return conf, fmt.Errorf("failed to parse configuration: %w", err)
	}
	return conf, nil
}

func Getenv(key string) string {
	loadEnvOnce()
	return os.Getenv(key)
}
