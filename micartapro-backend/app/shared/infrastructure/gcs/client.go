package gcs

import (
	"context"
	"micartapro/app/shared/configuration"

	"cloud.google.com/go/storage"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewClient, configuration.NewConf)
}

func NewClient(conf configuration.Conf) (*storage.Client, error) {
	return storage.NewClient(context.Background())
}
