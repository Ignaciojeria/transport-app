package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/natsconsumer/model"
	"transport-app/app/domain"
	canonicaljson "transport-app/app/shared/caonincaljson"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/infrastructure/storj"
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase"

	client "transport-app/app/adapter/out/restyclient/webhook"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newFleetOptimizedWebhookSubmitted,
		natsconn.NewJetStream,
		observability.NewObservability,
		configuration.NewConf,
		natsconn.NewKeyValue,
		usecase.NewGetDataFromRedisWorkflow,
		client.NewPostWebhook,
		usecase.NewUpsertElectricRouteWorkflow,
		storj.NewUplink,
	)
}

func newFleetOptimizedWebhookSubmitted(
	js jetstream.JetStream,
	obs observability.Observability,
	conf configuration.Conf,
	kv jetstream.KeyValue,
	getDataFromRedisWorkflow usecase.GetDataFromRedisWorkflow,
	publishCustomerWebhook client.PostWebhook,
	upsertElectricRoute usecase.UpsertElectricRouteWorkflow,
	storjManager storj.UplinkManager,
) (jetstream.ConsumeContext, error) {
	// Validación para verificar si el nombre de la suscripción está vacío
	if conf.FLEET_OPTIMIZED_WEBHOOK_SUBSCRIPTION == "" {
		obs.Logger.Warn("Webhook submitted subscription name is empty, skipping consumer initialization")
		// Retornar nil para indicar que no hay consumidor activo
		return nil, nil
	}

	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, conf.TRANSPORT_APP_TOPIC, jetstream.ConsumerConfig{
		Name:          fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.FLEET_OPTIMIZED_WEBHOOK_SUBSCRIPTION),
		Durable:       fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.FLEET_OPTIMIZED_WEBHOOK_SUBSCRIPTION),
		FilterSubject: conf.TRANSPORT_APP_TOPIC + "." + conf.ENVIRONMENT + ".*.*.fleetOptimizedWebhook",
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

		// Deserializar el payload como UpsertWebhookRequest
		var input model.FleetOptimizedWebhookBody
		if err := json.Unmarshal(pubsubMsg.Data, &input); err != nil {
			obs.Logger.Error("Error deserializando payload de webhook", "error", err)
			msg.Ack()
			return
		}

		wh := domain.Webhook{
			Type: "fleet-optimized",
		}

		ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(pubsubMsg.Attributes))

		// Generar idempotency key para el webhook
		webhookKey, err := canonicaljson.HashKey(ctx, "fleet_optimized_webhook", input)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para webhook", "error", err)
			msg.Nak()
			return
		}
		webhookCtx := sharedcontext.WithIdempotencyKey(ctx, webhookKey)

		bytes, err := kv.Get(webhookCtx, wh.DocID(webhookCtx).String())
		if err != nil {
			obs.Logger.Error("Error obteniendo webhook", "error", err)
			msg.Ack()
			return
		}

		var webhook domain.Webhook
		if err := json.Unmarshal(bytes.Value(), &webhook); err != nil {
			obs.Logger.Error("Error deserializando webhook", "error", err)
			msg.Ack()
			return
		}

		webhook.Body = input

		for _, routeID := range input.Routes {
			data, err := getDataFromRedisWorkflow(webhookCtx, routeID)
			if err != nil {
				obs.Logger.Error("Error descargando ruta desde Storj", "error", err)
				msg.Ack()
				return
			}
			var routeRequest request.UpsertRouteRequest
			if err := json.Unmarshal(data, &routeRequest); err != nil {
				obs.Logger.Error("Error deserializando datos de ruta", "error", err)
				msg.Ack()
				return
			}

			// Asignar URLs a delivery units antes del upsert
			if err := assignURLsToDeliveryUnits(webhookCtx, &routeRequest, storjManager, obs); err != nil {
				obs.Logger.Error("Error asignando URLs a delivery units", "error", err)
				msg.Nak()
				return
			}

			route, err := routeRequest.Map()
			if err != nil {
				obs.Logger.Error("Error mappeando datos de ruta", "error", err)
				msg.Ack()
				return
			}

			plan := domain.Plan{
				ReferenceID: input.Plan,
			}

			// Generar idempotency key para la inserción de ruta
			routeKey, err := canonicaljson.HashKey(webhookCtx, "upsert_electric_route", map[string]interface{}{
				"routeID": routeID,
				"plan":    input.Plan,
			})
			if err != nil {
				obs.Logger.ErrorContext(webhookCtx, "Error generando key para inserción de ruta", "error", err)
				msg.Nak()
				return
			}
			routeCtx := sharedcontext.WithIdempotencyKey(webhookCtx, routeKey)

			err = upsertElectricRoute(routeCtx, route, plan.DocID(routeCtx).String(), routeRequest)
			if err != nil {
				obs.Logger.Error("Error insertando ruta", "error", err)
				msg.Nak()
				return
			}

		}

		if err := publishCustomerWebhook(webhookCtx, webhook); err != nil {
			obs.Logger.Error("Error posteando webhook", "error", err)
			msg.Nak()
			return
		}

		obs.Logger.InfoContext(webhookCtx, "Webhook procesado exitosamente desde NATS")
		msg.Ack()
	})
}

func assignURLsToDeliveryUnits(ctx context.Context, routeRequest *request.UpsertRouteRequest, storjManager storj.UplinkManager, obs observability.Observability) error {
	uploadTTL := 30 * 24 * time.Hour  // 1 mes para upload
	// Sin expiración para download (10 años)
	downloadTTL := 10 * 365 * 24 * time.Hour

	for visitIdx := range routeRequest.Visits {
		for orderIdx := range routeRequest.Visits[visitIdx].Orders {
			for duIdx := range routeRequest.Visits[visitIdx].Orders[orderIdx].DeliveryUnits {
				du := &routeRequest.Visits[visitIdx].Orders[orderIdx].DeliveryUnits[duIdx]
				
				// Generar objectKey único para este delivery unit
				fileID := uuid.New()
				objectKey := fmt.Sprintf("deliveries/%s/evidence_%s.jpg", 
					du.DocumentID, 
					fileID.String()[:8])

				// Generar upload URL (1 mes)
				uploadURL, err := storjManager.GeneratePreSignedURL(ctx, objectKey, uploadTTL)
				if err != nil {
					obs.Logger.Error("Error generando upload URL", "error", err, "documentID", du.DocumentID)
					return fmt.Errorf("failed to generate upload URL for delivery unit %s: %w", du.DocumentID, err)
				}

				// Generar download URL (sin expiración)
				downloadURL, err := storjManager.GeneratePublicDownloadURL(ctx, objectKey, downloadTTL)
				if err != nil {
					obs.Logger.Error("Error generando download URL", "error", err, "documentID", du.DocumentID)
					return fmt.Errorf("failed to generate download URL for delivery unit %s: %w", du.DocumentID, err)
				}

				// Asignar las URLs al delivery unit
				evidence := request.UpsertRouteEvidences{
					UploadUrl:   uploadURL,
					DownloadUrl: downloadURL,
				}
				du.Evidences = append(du.Evidences, evidence)

				obs.Logger.InfoContext(ctx, "URLs asignadas a delivery unit", "documentID", du.DocumentID, "lpn", du.Lpn)
			}
		}
	}

	return nil
}
