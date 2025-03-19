package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewTiDBConfiguration)
}

type DBConfiguration struct {
	DB_STRATEGY       string `env:"DB_STRATEGY,required"`
	DB_HOSTNAME       string `env:"DB_HOSTNAME,required"`
	DB_PORT           string `env:"DB_PORT,required"`
	DB_SSL_MODE       string `env:"DB_SSL_MODE,required"`
	DB_NAME           string `env:"DB_NAME,required"`
	DB_USERNAME       string `env:"DB_USERNAME,required"`
	DB_PASSWORD       string `env:"DB_PASSWORD,required"`
	DB_RUN_MIGRATIONS string `env:"DB_RUN_MIGRATIONS"`
}

func NewTiDBConfiguration() (DBConfiguration, error) {
	return Parse[DBConfiguration]()
}
