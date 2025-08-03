package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/out/storjbucket"
	canonicaljson "transport-app/app/shared/caonincaljson"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase"
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
		storjbucket.NewTransportAppBucket,
		workers.NewFleetOptimizer,
		usecase.NewOptimizeFleetWorkflow,
		observability.NewObservability,
		configuration.NewConf,
		usecase.NewStoreDataInBucketWorkflow,
		usecase.NewPublishWebhookWorkflow,
	)
}

func newOptimizationRequestedConsumer(
	js jetstream.JetStream,
	kv jetstream.KeyValue,
	storjBucket *storjbucket.TransportAppBucket,
	optimize workers.FleetOptimizer,
	optimizeFleetWorkflow usecase.OptimizeFleetWorkflow,
	obs observability.Observability,
	conf configuration.Conf,
	storeDataInBucketWorkflow usecase.StoreDataInBucketWorkflow,
	publishWebhookWorkflow usecase.PublishWebhookWorkflow,
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

			// Orquestación usando workflows
			key, err := canonicaljson.HashKey(ctx, "optimize_fleet", input)
			if err != nil {
				obs.Logger.ErrorContext(ctx, "Error procesando headers de orden", "error", err)
				msg.Ack()
				return
			}
			optimizeFleetWorkflowCtx := sharedcontext.WithIdempotencyKey(ctx, key)
			optimizeFleetWorkflowCtx = sharedcontext.WithAccessToken(optimizeFleetWorkflowCtx, msg.Headers().Get("X-Access-Token"))
			optimizeFleetWorkflowCtx = sharedcontext.WithBucketToken(optimizeFleetWorkflowCtx, msg.Headers().Get("X-Bucket-Token"))
			routeRequests, err := optimizeFleetWorkflow(optimizeFleetWorkflowCtx, input.Map())
			if err != nil {
				obs.Logger.ErrorContext(ctx, "Error procesando optimización", "error", err)
				msg.Ack()
				return
			}

			for _, routeRequest := range routeRequests {
				routeRequestBytes, err := json.Marshal(routeRequest)
				if err != nil {
					obs.Logger.ErrorContext(ctx, "Error serializando ruta de optimización", "error", err)
					msg.Ack()
					return
				}
				routeKey, err := canonicaljson.HashKey(ctx, "store_route_in_bucket", routeRequest.ReferenceID)
				storeDataInBucketWorkflowCtx := sharedcontext.WithIdempotencyKey(ctx, routeKey)
				storeDataInBucketWorkflowCtx = sharedcontext.WithBucketToken(storeDataInBucketWorkflowCtx, msg.Headers().Get("X-Bucket-Token"))
				err = storeDataInBucketWorkflow(storeDataInBucketWorkflowCtx, routeRequest.ReferenceID, routeRequestBytes)
				if err != nil {
					obs.Logger.ErrorContext(ctx, "Error almacenando ruta de optimización en bucket", "error", err)
					msg.Ack()
					return
				}
			}

			webhookKey, err := canonicaljson.HashKey(ctx, "publish_webhook", input)
			publishWebhookWorkflowCtx := sharedcontext.WithAccessToken(ctx, msg.Headers().Get("X-Access-Token"))
			publishWebhookWorkflowCtx = sharedcontext.WithBucketToken(publishWebhookWorkflowCtx, msg.Headers().Get("X-Bucket-Token"))
			publishWebhookWorkflowCtx = sharedcontext.WithIdempotencyKey(publishWebhookWorkflowCtx, webhookKey)
			type fleetOptimizedWebhook struct {
				Plan   string   `json:"plan"`
				Routes []string `json:"routes"`
			}
			var webhookBody fleetOptimizedWebhook

			for _, routeRequest := range routeRequests {
				webhookBody.Routes = append(webhookBody.Routes, routeRequest.ReferenceID)
			}
			if len(routeRequests) > 0 {
				webhookBody.Plan = routeRequests[0].PlanReferenceID
			}

			if err := publishWebhookWorkflow(publishWebhookWorkflowCtx, webhookBody, "fleet-optimized"); err != nil {
				obs.Logger.ErrorContext(ctx, "Error publicando webhook", "error", err)
				msg.Ack()
				return
			}

			obs.Logger.InfoContext(ctx, "Optimización procesada exitosamente desde NATS (reconstruido)",
				"eventType", "optimizationRequested")
			msg.Ack()
			return
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
