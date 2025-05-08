package gcppubsub

import (
	"context"
	"errors"
	"transport-app/app/shared/configuration"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewClient, configuration.NewConf)
}

func NewClient(conf configuration.Conf) (*pubsub.Client, error) {
	c, err := pubsub.NewClient(context.Background(), conf.GOOGLE_PROJECT_ID)
	if conf.GOOGLE_PROJECT_ID == "" {
		return &pubsub.Client{}, errors.New("GOOGLE_PROJECT_ID is not present")
	}
	if err != nil {
		return &pubsub.Client{}, err
	}
	return c, nil
}
