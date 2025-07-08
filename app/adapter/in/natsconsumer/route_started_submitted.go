package natsconsumer

import (
	"context"
	"encoding/json"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/usecase"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newRouteStartedSubmittedConsumer,
		natsconn.NewJetStream,
		usecase.NewRouteStarted,
		observability.NewObservability,
		configuration.NewConf,
	)
}

func newRouteStartedSubmittedConsumer(
	js jetstream.JetStream,
	routeStarted usecase.RouteStarted,
	obs observability.Observability,
	conf configuration.Conf,
) (jetstream.ConsumeContext, error) {
	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, conf.TRANSPORT_APP_TOPIC, jetstream.ConsumerConfig{
		Name:          conf.ROUTE_STARTED_SUBMITTED_SUBSCRIPTION,
		Durable:       conf.ROUTE_STARTED_SUBMITTED_SUBSCRIPTION,
		FilterSubject: conf.TRANSPORT_APP_TOPIC + ".*.routeStarted",
		MaxAckPending: 5,
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

		// Deserializar el payload como RouteStartedRequest
		var input request.RouteStartedRequest
		if err := json.Unmarshal(pubsubMsg.Data, &input); err != nil {
			obs.Logger.Error("Error deserializando payload de ruta iniciada", "error", err)
			msg.Ack()
			return
		}

		// Procesar la ruta iniciada
		if err := routeStarted(ctx, input.Map(ctx)); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando ruta iniciada", "error", err)
			msg.Ack()
			return
		}

		obs.Logger.InfoContext(ctx, "Ruta iniciada procesada exitosamente desde NATS",
			"eventType", "routeStarted")
		msg.Ack()
	})
}
