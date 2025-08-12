package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/natsconsumer/mapper"
	"transport-app/app/adapter/out/natspublisher"
	"transport-app/app/adapter/out/storjbucket"
	"transport-app/app/domain"
	"transport-app/app/shared/chunker"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
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
		natspublisher.NewApplicationEvents,
		usecase.NewKeyNormalizationWorkflow,
	)
}

func newAgentOptimizationRequested(
	js jetstream.JetStream,
	obs observability.Observability,
	conf configuration.Conf,
	storjBucket *storjbucket.TransportAppBucket,
	publish natspublisher.ApplicationEvents,
	keyNormalizationWorkflow usecase.KeyNormalizationWorkflow,
) (jetstream.ConsumeContext, error) {
	// Validación para verificar si el nombre de la suscripción está vacío
	if conf.AGENT_OPTIMIZATION_REQUESTED_SUBSCRIPTION == "" {
		obs.Logger.Warn("Agent optimization requested subscription name is empty, skipping consumer initialization")
		// Retornar nil para indicar que no hay consumidor activo
		return nil, nil
	}

	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, conf.TRANSPORT_APP_TOPIC, jetstream.ConsumerConfig{
		Name:          fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.AGENT_OPTIMIZATION_REQUESTED_SUBSCRIPTION),
		Durable:       fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.AGENT_OPTIMIZATION_REQUESTED_SUBSCRIPTION),
		FilterSubject: conf.TRANSPORT_APP_TOPIC + "." + conf.ENVIRONMENT + ".*.*.agentOptimizationRequested",
		MaxAckPending: 100,
		AckWait:       5 * time.Minute, // 5 minutos para procesamiento de IA
		// Configuración de reintentos: 3 reintentos con intervalos de 2 segundos
		MaxDeliver: 4, // 1 intento inicial + 3 reintentos = 4 total
		BackOff:    []time.Duration{2 * time.Second, 2 * time.Second, 2 * time.Second},
	})

	if err != nil {
		return nil, err
	}

	visitFieldMapper := mapper.NewVisitFieldMapper()
	vehicleFieldMapper := mapper.NewVehicleFieldMapper()

	return consumer.Consume(func(msg jetstream.Msg) {
		var pubsubMsg pubsub.Message
		msg.Ack()
		if err := json.Unmarshal(msg.Data(), &pubsubMsg); err != nil {
			obs.Logger.Error("Error deserializando mensaje NATS", "error", err)
			msg.Ack()
			return
		}

		// Extraer contexto de OpenTelemetry
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(pubsubMsg.Attributes))
		ctx = sharedcontext.WithAccessToken(ctx, msg.Headers().Get("X-Access-Token"))

		var request request.AgentOptimizationRequest

		// Intentar deserializar el payload como arreglo de IDs de chunks
		var chunkIDs []string
		err := json.Unmarshal(pubsubMsg.Data, &chunkIDs)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error deserializando mensaje NATS", "error", err)
			msg.Ack()
			return
		}

		var chunks []chunker.Chunk
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

		// Deserializar el payload reconstruido
		if err := json.Unmarshal(data, &request); err != nil {
			obs.Logger.ErrorContext(ctx, "Error deserializando payload de agent optimization request", "error", err)
			msg.Ack()
			return
		}

		// Validar que haya visitas y flota para procesar
		if len(request.Visits) == 0 || len(request.Fleet) == 0 {
			obs.Logger.WarnContext(ctx, "Skipping optimization request - no visits or fleet to process",
				"visitsCount", len(request.Visits),
				"fleetCount", len(request.Fleet))
			msg.Ack()
			return
		}

		// Obtener el mapeo de claves usando la primera visita como ejemplo
		keyMapping := visitFieldMapper.Map(request.Visits[0]) //visitFieldNamesNormalizerWorkflow(ctx, request.Visits[0])
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error obteniendo mapeo de claves", "error", err)
			msg.Nak()
			return
		}

		obs.Logger.InfoContext(ctx, "Mapeo de claves obtenido", "keyMapping", keyMapping)

		// Usar el workflow genérico para normalizar las visitas
		normalizedVisits, err := keyNormalizationWorkflow(ctx, keyMapping, request.Visits)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error normalizando visitas", "error", err)
			msg.Nak()
			return
		}

		vehicleKeyMapping := vehicleFieldMapper.Map(request.Fleet[0])
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error obteniendo mapeo de claves", "error", err)
			msg.Nak()
			return
		}

		normalizedVehicles, err := keyNormalizationWorkflow(ctx, vehicleKeyMapping, request.Fleet)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error normalizando vehículos", "error", err)
			msg.Nak()
			return
		}

		// Reemplazar las visitas originales con las normalizadas
		request.Visits = normalizedVisits
		request.Fleet = normalizedVehicles

		optimizeFleetRequest := request.ToOptimizeFleetRequest()
		optimizeFleetRequest.VehicleKeyMapping = vehicleKeyMapping
		optimizeFleetRequest.VisitKeyMapping = keyMapping

		obs.Logger.InfoContext(ctx, "Optimize fleet request", "input", optimizeFleetRequest)

		// Publicar el evento de optimización de flota
		eventPayload, _ := json.Marshal(optimizeFleetRequest)

		eventCtx := sharedcontext.AddEventContextToBaggage(ctx,
			sharedcontext.EventContext{
				EntityType: "optimization",
				EventType:  "optimizationRequested",
			})

		if err := publish(eventCtx, domain.Outbox{
			Payload: eventPayload,
		}); err != nil {
			obs.Logger.ErrorContext(ctx, "Error publicando evento de optimización", "error", err)
			msg.Nak()
			return
		}

		obs.Logger.InfoContext(ctx, "OPTIMIZATION_REQUEST_PUBLISHED",
			"planReferenceID", optimizeFleetRequest.PlanReferenceID,
			"vehiclesCount", len(optimizeFleetRequest.Vehicles),
			"visitsCount", len(optimizeFleetRequest.Visits))

		obs.Logger.InfoContext(ctx, "Agent optimization request received", "input", request)
		msg.Ack()
	})
}
