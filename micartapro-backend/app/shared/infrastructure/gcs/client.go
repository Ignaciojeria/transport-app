package gcs

import (
	"context"
	"micartapro/app/shared/configuration"

	"cloud.google.com/go/storage"
	ioc "github.com/Ignaciojeria/ioc"
)

func init() {
	ioc.Register(NewClient)
}

func NewClient(conf configuration.Conf) (*storage.Client, error) {
	return storage.NewClient(context.Background())
}
