package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewTiDBConfiguration)
}

type TiDBConfiguration struct {
	TIDB_HOSTNAME       string `env:"TIDB_HOSTNAME,required"`
	TIDB_PORT           string `env:"TIDB_PORT,required"`
	TIDB_DATABASE       string `env:"TIDB_DATABASE,required"`
	TIDB_USERNAME       string `env:"TIDB_USERNAME,required"`
	TIDB_PASSWORD       string `env:"TIDB_PASSWORD,required"`
	TIDB_RUN_MIGRATIONS string `env:"TIDB_RUN_MIGRATIONS"`
}

func NewTiDBConfiguration() (TiDBConfiguration, error) {
	return Parse[TiDBConfiguration]()
}
