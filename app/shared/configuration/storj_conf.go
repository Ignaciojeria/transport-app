package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewStorjConfiguration)
}

type StorjConfiguration struct {
	STORJ_ACCESS_GRANT         string `env:"STORJ_ACCESS_GRANT"`
	STORJ_S3_ACCESS_KEY_ID     string `env:"STORJ_S3_ACCESS_KEY_ID"`
	STORJ_S3_SECRET_ACCESS_KEY string `env:"STORJ_S3_SECRET_ACCESS_KEY"`
	STORJ_S3_ENDPOINT          string `env:"STORJ_S3_ENDPOINT" envDefault:"https://gateway.storjshare.io"`
}

func NewStorjConfiguration() (StorjConfiguration, error) {
	return Parse[StorjConfiguration]()
}
