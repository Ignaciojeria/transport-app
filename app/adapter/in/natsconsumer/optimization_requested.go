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

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newOptimizationRequestedConsumer,
		natsconn.NewJetStream,
		workers.NewFleetOptimizer,
		observability.NewObservability,
		configuration.NewConf,
	)
}

func newOptimizationRequestedConsumer(
	js jetstream.JetStream,
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
		return nil, err
	}

	return consumer.Consume(func(msg jetstream.Msg) {
		//accessToken :=

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

		// Deserializar el payload como OptimizeFleetRequest
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
