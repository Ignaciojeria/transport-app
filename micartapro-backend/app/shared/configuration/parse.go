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

func loadEnvOnce() {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			slog.Warn(".env not found, loading environment variables from system.")
		} else {
			slog.Info("Environment variables loaded from .env file.")
		}
	})
}

func Parse[T any]() (T, error) {
	loadEnvOnce()
	var conf T
	if err := env.Parse(&conf); err != nil {
		return conf, fmt.Errorf("failed to parse configuration: %w", err)
	}
	return conf, nil
}

func Getenv(key string) string {
	loadEnvOnce()
	return os.Getenv(key)
}
