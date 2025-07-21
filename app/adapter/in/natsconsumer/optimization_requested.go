package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase/workers"

	"transport-app/app/shared/chunker"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newOptimizationRequestedConsumer,
		natsconn.NewJetStream,
		natsconn.NewKeyValue,
		workers.NewFleetOptimizer,
		observability.NewObservability,
		configuration.NewConf,
	)
}

func newOptimizationRequestedConsumer(
	js jetstream.JetStream,
	kv jetstream.KeyValue,
	optimize workers.FleetOptimizer,
	obs observability.Observability,
	conf configuration.Conf,
) (jetstream.ConsumeContext, error) {
	// Validación para verificar si el nombre de la suscripción está vacío
	if conf.OPTIMIZATION_REQUESTED_SUBSCRIPTION == "" {
		obs.Logger.Warn("Optimization requested subscription name is empty, skipping consumer initialization")
		// Retornar nil para indicar que no hay consumidor activo
		return nil, nil
	}

	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, conf.TRANSPORT_APP_TOPIC, jetstream.ConsumerConfig{
		Name:          fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.OPTIMIZATION_REQUESTED_SUBSCRIPTION),
		Durable:       fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.OPTIMIZATION_REQUESTED_SUBSCRIPTION),
		FilterSubject: conf.TRANSPORT_APP_TOPIC + "." + conf.ENVIRONMENT + ".*.*.optimizationRequested",
		MaxAckPending: 5,
	})

	if err != nil {
		obs.Logger.Warn("Consumer not found or inaccessible. Ensure it is pre-created and accessible with current permissions.", "error", err)

		consumer, err = js.Consumer(ctx, conf.TRANSPORT_APP_TOPIC, fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.OPTIMIZATION_REQUESTED_SUBSCRIPTION))
		if err != nil {
			obs.Logger.Warn("No se pudo obtener el consumidor preexistente", "error", err)
			return nil, err
		}
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
			var chunks []chunker.Chunk
			for idx, id := range chunkIDs {
				entry, err := kv.Get(ctx, id)
				if err != nil {
					obs.Logger.ErrorContext(ctx, "Error obteniendo chunk del KV store", "chunkID", id, "error", err)
					msg.Ack()
					return
				}
				chunks = append(chunks, chunker.Chunk{
					ID:   uuidMustParse(id),
					Data: entry.Value(),
					Idx:  idx,
				})
			}
			// Ordenar los chunks por Idx por si acaso (aunque el orden del array debería ser correcto)
			// Reconstruir el mensaje original
			data, err := chunker.ReconstructBytes(chunks)
			if err != nil {
				obs.Logger.ErrorContext(ctx, "Error reconstruyendo mensaje desde chunks", "error", err)
				msg.Ack()
				return
			}
			// Deserializar el payload reconstruido como OptimizeFleetRequest
			var input request.OptimizeFleetRequest
			if err := json.Unmarshal(data, &input); err != nil {
				obs.Logger.ErrorContext(ctx, "Error deserializando payload de optimización (reconstruido)", "error", err)
				msg.Ack()
				return
			}
			// Procesar la optimización
			if err := optimize(ctx, input.Map()); err != nil {
				obs.Logger.ErrorContext(ctx, "Error procesando optimización (reconstruido)", "error", err)
				msg.Ack()
				return
			}
			obs.Logger.InfoContext(ctx, "Optimización procesada exitosamente desde NATS (reconstruido)",
				"eventType", "optimizationRequested")
			msg.Ack()
			return
		}

		// Si no es chunked, intentar deserializar el payload como OptimizeFleetRequest (flujo original)
		var input request.OptimizeFleetRequest
		if err := json.Unmarshal(pubsubMsg.Data, &input); err != nil {
			obs.Logger.Error("Error deserializando payload de optimización", "error", err)
			msg.Ack()
			return
		}

		// Procesar la optimización
		if err := optimize(ctx, input.Map()); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando optimización", "error", err)
			msg.Ack()
			return
		}

		obs.Logger.InfoContext(ctx, "Optimización procesada exitosamente desde NATS",
			"eventType", "optimizationRequested")
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
