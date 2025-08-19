package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/natsconsumer/model"
	"transport-app/app/adapter/out/storjbucket"
	"transport-app/app/domain"
	canonicaljson "transport-app/app/shared/caonincaljson"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase"

	client "transport-app/app/adapter/out/restyclient/webhook"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"storj.io/uplink"
)

func init() {
	ioc.Registry(
		newFleetOptimizedWebhookSubmitted,
		natsconn.NewJetStream,
		observability.NewObservability,
		configuration.NewConf,
		natsconn.NewKeyValue,
		storjbucket.NewTransportAppBucket,
		client.NewPostWebhook,
		usecase.NewUpsertElectricRouteWorkflow,
	)
}

func newFleetOptimizedWebhookSubmitted(
	js jetstream.JetStream,
	obs observability.Observability,
	conf configuration.Conf,
	kv jetstream.KeyValue,
	storjBucket *storjbucket.TransportAppBucket,
	publishCustomerWebhook client.PostWebhook,
	upsertElectricRoute usecase.UpsertElectricRouteWorkflow,
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
		//accessToken := msg.Headers().Get("X-Access-Token")
		//webhook.Headers["X-Access-Token"] = accessToken
		//webhook.Headers["tenant"] = msg.Headers().Get("tenant")

		// Generar token de acceso para Storj
		token, err := storjBucket.GenerateEphemeralToken(webhookCtx, 10*time.Minute, uplink.Permission{
			AllowDownload: true,
		})
		if err != nil {
			obs.Logger.Error("Error generando token de acceso", "error", err)
			msg.Ack()
			return
		}

		for _, routeID := range input.Routes {
			data, err := storjBucket.DownloadWithToken(webhookCtx, token, routeID)
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
