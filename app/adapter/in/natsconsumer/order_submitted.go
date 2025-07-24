package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/usecase"

	"transport-app/app/domain"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newNatsConsumer,
		natsconn.NewJetStream,
		usecase.NewCreateOrder,
		observability.NewObservability,
		configuration.NewConf,
		usecase.NewUpsertAddressInfoWorkflow,
		usecase.NewUpsertOrderHeadersWorkflow,
		usecase.NewUpsertContactWorkflow,
		usecase.NewUpsertOrderTypeWorkflow,
		usecase.NewUpsertOrderWorkflow,
		usecase.NewUpsertOrderReferencesWorkflow,
		usecase.NewUpsertOrderDeliveryUnitsWorkflow,
		usecase.NewUpsertDeliveryUnitsWorkflow,
		usecase.NewUpsertDeliveryUnitsLabelsWorkflow,
		usecase.NewUpsertDeliveryUnitsSkillsWorkflow,
		usecase.NewUpsertSizeCategoryWorkflow,
		usecase.NewUpsertDeliveryUnitsHistoryWorkflow,
		usecase.NewUpsertSkillWorkflow,
	)
}

func newNatsConsumer(
	js jetstream.JetStream,
	createOrder usecase.CreateOrder, // se puede eliminar después
	obs observability.Observability,
	conf configuration.Conf,
	upsertAddressInfoWorkflow usecase.UpsertAddressInfoWorkflow,
	upsertOrderHeadersWorkflow usecase.UpsertOrderHeadersWorkflow,
	upsertContactWorkflow usecase.UpsertContactWorkflow,
	upsertOrderTypeWorkflow usecase.UpsertOrderTypeWorkflow,
	upsertOrderWorkflow usecase.UpsertOrderWorkflow,
	upsertOrderReferencesWorkflow usecase.UpsertOrderReferencesWorkflow,
	upsertOrderDeliveryUnitsWorkflow usecase.UpsertOrderDeliveryUnitsWorkflow,
	upsertDeliveryUnitsWorkflow usecase.UpsertDeliveryUnitsWorkflow,
	upsertDeliveryUnitsLabelsWorkflow usecase.UpsertDeliveryUnitsLabelsWorkflow,
	upsertDeliveryUnitsSkillsWorkflow usecase.UpsertDeliveryUnitsSkillsWorkflow,
	upsertSizeCategoryWorkflow usecase.UpsertSizeCategoryWorkflow,
	upsertDeliveryUnitsHistoryWorkflow usecase.UpsertDeliveryUnitsHistoryWorkflow,
	upsertSkillWorkflow usecase.UpsertSkillWorkflow,
) (jetstream.ConsumeContext, error) {
	// Validación para verificar si el nombre de la suscripción está vacío
	if conf.ORDER_SUBMITTED_SUBSCRIPTION == "" {
		obs.Logger.Warn("Order submitted subscription name is empty, skipping consumer initialization")
		// Retornar nil para indicar que no hay consumidor activo
		return nil, nil
	}

	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, conf.TRANSPORT_APP_TOPIC, jetstream.ConsumerConfig{
		Name:          fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.ORDER_SUBMITTED_SUBSCRIPTION),
		Durable:       fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.ORDER_SUBMITTED_SUBSCRIPTION),
		FilterSubject: conf.TRANSPORT_APP_TOPIC + "." + conf.ENVIRONMENT + ".*.*.orderSubmitted",
		MaxAckPending: 5,
		MaxDeliver:    4, // 1 intento inicial + 3 reintentos = 4 total
		BackOff:       []time.Duration{2 * time.Second, 2 * time.Second, 2 * time.Second},
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

		// Deserializar el payload como UpsertOrderRequest
		var input request.UpsertOrderRequest
		if err := json.Unmarshal(pubsubMsg.Data, &input); err != nil {
			obs.Logger.Error("Error deserializando payload de orden", "error", err)
			msg.Nak()
			return
		}

		order := input.Map(ctx)
		order.Origin.AddressInfo.ToLowerAndRemovePunctuation()
		order.Destination.AddressInfo.ToLowerAndRemovePunctuation()
		order.AssignIndexesIfNoLPN()

		// Orquestación usando workflows
		if err := upsertOrderHeadersWorkflow(ctx, order.Headers); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando headers de orden", "error", err)
			msg.Nak()
			return
		}
		if err := upsertContactWorkflow(ctx, order.Origin.AddressInfo.Contact); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando contacto origen", "error", err)
			msg.Nak()
			return
		}
		if err := upsertContactWorkflow(ctx, order.Destination.AddressInfo.Contact); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando contacto destino", "error", err)
			msg.Nak()
			return
		}
		if err := upsertAddressInfoWorkflow(ctx, order.Origin.AddressInfo); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando dirección origen", "error", err)
			msg.Nak()
			return
		}
		if err := upsertAddressInfoWorkflow(ctx, order.Destination.AddressInfo); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando dirección destino", "error", err)
			msg.Nak()
			return
		}
		if err := upsertOrderTypeWorkflow(ctx, order.OrderType); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando tipo de orden", "error", err)
			msg.Nak()
			return
		}
		if err := upsertDeliveryUnitsWorkflow(ctx, order.DeliveryUnits); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando unidades de entrega", "error", err)
			msg.Nak()
			return
		}
		if err := upsertOrderReferencesWorkflow(ctx, order); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando referencias de orden", "error", err)
			msg.Nak()
			return
		}
		if err := upsertOrderDeliveryUnitsWorkflow(ctx, order); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando relación orden-unidades", "error", err)
			msg.Nak()
			return
		}
		if err := upsertOrderWorkflow(ctx, order); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando upsert de orden", "error", err)
			msg.Nak()
			return
		}
		if err := upsertDeliveryUnitsLabelsWorkflow(ctx, order); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando labels de unidades de entrega", "error", err)
			msg.Nak()
			return
		}
		if err := upsertDeliveryUnitsSkillsWorkflow(ctx, order); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando skills de unidades de entrega", "error", err)
			msg.Nak()
			return
		}
		for _, du := range order.DeliveryUnits {
			if err := upsertSizeCategoryWorkflow(ctx, du.SizeCategory); err != nil {
				obs.Logger.ErrorContext(ctx, "Error procesando size category", "error", err)
				msg.Nak()
				return
			}
			for _, skill := range du.Skills {
				if err := upsertSkillWorkflow(ctx, skill); err != nil {
					obs.Logger.ErrorContext(ctx, "Error procesando skill", "error", err)
					msg.Nak()
					return
				}
			}
		}
		plan := domain.Plan{Routes: []domain.Route{{Orders: []domain.Order{order}}}}
		if err := upsertDeliveryUnitsHistoryWorkflow(ctx, plan); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando historial de unidades de entrega", "error", err)
			msg.Nak()
			return
		}

		obs.Logger.InfoContext(ctx, "Orden procesada exitosamente desde NATS",
			"eventType", "orderSubmitted",
			"referenceID", input.ReferenceID)
		msg.Ack()
	})
}
