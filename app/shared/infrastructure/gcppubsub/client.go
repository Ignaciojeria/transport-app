package gcppubsub

import (
	"context"
	"fmt"
	"transport-app/app/shared/configuration"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewClient, configuration.NewConf)
}

func NewClient(conf configuration.Conf) (*pubsub.Client, error) {
	if conf.GOOGLE_PROJECT_ID == "" {
		fmt.Println("GOOGLE_PROJECT_ID is not present")
		return &pubsub.Client{}, nil
	}
	c, err := pubsub.NewClient(context.Background(), conf.GOOGLE_PROJECT_ID)
	if err != nil {
		return &pubsub.Client{}, err
	}
	return c, nil
}
