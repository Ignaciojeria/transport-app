package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
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
		newOrderCancellationSubmittedConsumer,
		natsconn.NewJetStream,
		usecase.NewCancelOrder,
		observability.NewObservability,
		configuration.NewConf,
	)
}

func newOrderCancellationSubmittedConsumer(
	js jetstream.JetStream,
	cancelOrder usecase.CancelOrder,
	obs observability.Observability,
	conf configuration.Conf,
) (jetstream.ConsumeContext, error) {
	// Validación para verificar si el nombre de la suscripción está vacío
	if conf.ORDER_CANCELLATION_SUBMITTED_SUBSCRIPTION == "" {
		obs.Logger.Warn("Order cancellation submitted subscription name is empty, skipping consumer initialization")
		// Retornar nil para indicar que no hay consumidor activo
		return nil, nil
	}

	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, conf.TRANSPORT_APP_TOPIC, jetstream.ConsumerConfig{
		Name:          fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.ORDER_CANCELLATION_SUBMITTED_SUBSCRIPTION),
		Durable:       fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.ORDER_CANCELLATION_SUBMITTED_SUBSCRIPTION),
		FilterSubject: conf.TRANSPORT_APP_TOPIC + "." + conf.ENVIRONMENT + ".*.*.ordersCancellationSubmitted",
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

		// Deserializar el payload como CancelOrdersRequest
		var input request.CancelOrdersRequest
		if err := json.Unmarshal(pubsubMsg.Data, &input); err != nil {
			obs.Logger.Error("Error deserializando payload de cancelación de órdenes", "error", err)
			msg.Ack()
			return
		}

		// Procesar la cancelación de órdenes
		if err := cancelOrder(ctx, input.Map(ctx)); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando cancelación de órdenes", "error", err)
			msg.Ack()
			return
		}

		obs.Logger.InfoContext(ctx, "Cancelación de órdenes procesada exitosamente desde NATS",
			"eventType", "ordersCancellationSubmitted")
		msg.Ack()
	})
}
