package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewTiDBConfiguration)
}

type DBConfiguration struct {
	DB_STRATEGY       string `env:"DB_STRATEGY"`
	DB_HOSTNAME       string `env:"DB_HOSTNAME"`
	DB_PORT           string `env:"DB_PORT"`
	DB_SSL_MODE       string `env:"DB_SSL_MODE"`
	DB_NAME           string `env:"DB_NAME"`
	DB_USERNAME       string `env:"DB_USERNAME"`
	DB_PASSWORD       string `env:"DB_PASSWORD"`
	DB_RUN_MIGRATIONS string `env:"DB_RUN_MIGRATIONS"`
}

func NewTiDBConfiguration() (DBConfiguration, error) {
	return Parse[DBConfiguration]()
}
