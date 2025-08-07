package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/adapter/out/storjbucket"
	"transport-app/app/shared/chunker"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newAgentOptimizationRequested,
		natsconn.NewJetStream,
		observability.NewObservability,
		configuration.NewConf,
		storjbucket.NewTransportAppBucket,
	)
}

func newAgentOptimizationRequested(
	js jetstream.JetStream,
	obs observability.Observability,
	conf configuration.Conf,
	storjBucket *storjbucket.TransportAppBucket,
) (jetstream.ConsumeContext, error) {
	// Validación para verificar si el nombre de la suscripción está vacío
	if conf.AGENT_OPTIMIZATION_REQUESTED_SUBSCRIPTION == "" {
		obs.Logger.Warn("Agent optimization requested subscription name is empty, skipping consumer initialization")
		// Retornar nil para indicar que no hay consumidor activo
		return nil, nil
	}

	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, conf.TRANSPORT_APP_TOPIC, jetstream.ConsumerConfig{
		Name:          fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.WEBHOOK_SUBMITTED_SUBSCRIPTION),
		Durable:       fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.WEBHOOK_SUBMITTED_SUBSCRIPTION),
		FilterSubject: conf.TRANSPORT_APP_TOPIC + "." + conf.ENVIRONMENT + ".*.*.agentOptimizationRequested",
		MaxAckPending: 5,
		// Configuración de reintentos: 3 reintentos con intervalos de 2 segundos
		MaxDeliver: 4, // 1 intento inicial + 3 reintentos = 4 total
		BackOff:    []time.Duration{2 * time.Second, 2 * time.Second, 2 * time.Second},
	})

	if err != nil {
		return nil, err
	}

	return consumer.Consume(func(msg jetstream.Msg) {
		// Deserializar el mensaje como pubsub.Message
		var pubsubMsg pubsub.Message
		if err := json.Unmarshal(msg.Data(), &pubsubMsg); err != nil {
			obs.Logger.Error("Error deserializando mensaje NATS", "error", err)
			msg.Ack()
			return
		}

		// Extraer contexto de OpenTelemetry
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(pubsubMsg.Attributes))
		ctx = sharedcontext.WithAccessToken(ctx, msg.Headers().Get("X-Access-Token"))

		// Intentar deserializar el payload como arreglo de IDs de chunks
		var chunkIDs []string
		if err := json.Unmarshal(pubsubMsg.Data, &chunkIDs); err == nil && len(chunkIDs) > 0 {
			// Es un mensaje chunked, reconstruir el mensaje original
			var chunks []chunker.Chunk //a4c36c63-4c6e-4b0b-9bc7-3b452e76fa9d -a4c36c63-4c6e-4b0b-9bc7-3b452e76fa9d
			for idx, id := range chunkIDs {
				token := msg.Headers().Get("X-Bucket-Token")
				entry, err := storjBucket.DownloadWithToken(ctx, token, id)
				if err != nil {
					obs.Logger.ErrorContext(ctx, "Error obteniendo chunk del bucket", "chunkID", id, "error", err)
					msg.Ack()
					return
				}
				chunks = append(chunks, chunker.Chunk{
					ID:   uuidMustParse(id),
					Data: entry,
					Idx:  idx,
				})
			}

			data, err := chunker.ReconstructBytes(chunks)
			if err != nil {
				obs.Logger.ErrorContext(ctx, "Error reconstruyendo mensaje desde chunks", "error", err)
				msg.Ack()
				return
			}

			// Deserializar el payload reconstruido como OptimizeFleetRequest
			var input interface{}

			if err := json.Unmarshal(data, &input); err != nil {
				obs.Logger.ErrorContext(ctx, "Error deserializando payload de agent optimization request", "error", err)
				msg.Ack()
				return
			}

			obs.Logger.InfoContext(ctx, "Agent optimization request received", "input", input)
		}

		msg.Ack()
	})
}

// uuidMustParse es un helper para parsear UUID y hacer panic si falla (solo para uso interno seguro)
func uuidMustParse(id string) uuid.UUID {
	u, err := uuid.Parse(id)
	if err != nil {
		panic(err)
	}
	return u
}
