package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewStorjConfiguration)
}

type StorjConfiguration struct {
	STORJ_ACCESS_GRANT string `env:"STORJ_ACCESS_GRANT"`
}

func NewStorjConfiguration() (StorjConfiguration, error) {
	return Parse[StorjConfiguration]()
}
