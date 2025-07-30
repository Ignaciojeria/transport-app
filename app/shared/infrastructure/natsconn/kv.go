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
	// Intentar obtener el bucket existente
	kv, err := js.KeyValue(context.Background(), "transport-app-events-bucket")
	if err != nil {
		// Si no existe, crear el bucket con replicación síncrona
		kv, err = js.CreateKeyValue(context.Background(), jetstream.KeyValueConfig{
			Bucket:       "transport-app-events-bucket",
			Replicas:     1, // Para desarrollo local
			Storage:      jetstream.FileStorage,
			MaxValueSize: 1024 * 1024, // 1MB
		})
		if err != nil {
			return nil, err
		}
	}
	return kv, nil
}
