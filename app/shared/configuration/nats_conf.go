package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewNatsConfiguration)
}

type NatsConfiguration struct {
	NATS_CONNECTION_URL               string `env:"NATS_CONNECTION_URL"`
	NATS_CONNECTION_CREDS_FILEPATH    string `env:"NATS_CONNECTION_CREDS_FILEPATH" envDefault:"nats_credentials"`
	NATS_CONNECTION_CREDS_FILECONTENT string `env:"NATS_CONNECTION_CREDS_FILECONTENT"`
}

func NewNatsConfiguration() (NatsConfiguration, error) {
	return Parse[NatsConfiguration]()
}
