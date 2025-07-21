package natsconn

import (
	"context"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
)

func init() {
	ioc.Registry(NewKeyValue, NewJetStream)
}

func NewKeyValue(js jetstream.JetStream) (jetstream.KeyValue, error) {
	kv, err := js.KeyValue(context.Background(), "transport-app-events-bucket")
	if err != nil {
		return nil, err
	}
	return kv, nil
}
