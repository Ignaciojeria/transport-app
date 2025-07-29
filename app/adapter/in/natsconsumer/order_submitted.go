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
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase"

	"transport-app/app/domain"
	canonicaljson "transport-app/app/shared/caonincaljson"

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
		key, err := canonicaljson.HashKey(ctx, "order_headers", order.Headers)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando headers de orden", "error", err)
			msg.Ack()
			return
		}
		upsertHeadersCtx := sharedcontext.WithIdempotencyKey(ctx, key)
		if err := upsertOrderHeadersWorkflow(upsertHeadersCtx, order.Headers); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando headers de orden", "error", err)
			msg.Nak()
			return
		}

		// Contactos
		originContactKey, err := canonicaljson.HashKey(ctx, "contact", order.Origin.AddressInfo.Contact)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para contacto origen", "error", err)
			msg.Nak()
			return
		}
		originContactCtx := sharedcontext.WithIdempotencyKey(ctx, originContactKey)
		if err := upsertContactWorkflow(originContactCtx, order.Origin.AddressInfo.Contact); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando contacto origen", "error", err)
			msg.Nak()
			return
		}

		destContactKey, err := canonicaljson.HashKey(ctx, "contact", order.Destination.AddressInfo.Contact)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para contacto destino", "error", err)
			msg.Nak()
			return
		}
		destContactCtx := sharedcontext.WithIdempotencyKey(ctx, destContactKey)
		if err := upsertContactWorkflow(destContactCtx, order.Destination.AddressInfo.Contact); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando contacto destino", "error", err)
			msg.Nak()
			return
		}

		// Address Info
		originAddressKey, err := canonicaljson.HashKey(ctx, "address_info", order.Origin.AddressInfo)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para dirección origen", "error", err)
			msg.Nak()
			return
		}
		originAddressCtx := sharedcontext.WithIdempotencyKey(ctx, originAddressKey)
		if err := upsertAddressInfoWorkflow(originAddressCtx, order.Origin.AddressInfo); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando dirección origen", "error", err)
			msg.Nak()
			return
		}

		destAddressKey, err := canonicaljson.HashKey(ctx, "address_info", order.Destination.AddressInfo)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para dirección destino", "error", err)
			msg.Nak()
			return
		}
		destAddressCtx := sharedcontext.WithIdempotencyKey(ctx, destAddressKey)
		if err := upsertAddressInfoWorkflow(destAddressCtx, order.Destination.AddressInfo); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando dirección destino", "error", err)
			msg.Nak()
			return
		}

		// Order Type
		orderTypeKey, err := canonicaljson.HashKey(ctx, "order_type", order.OrderType)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para tipo de orden", "error", err)
			msg.Nak()
			return
		}
		orderTypeCtx := sharedcontext.WithIdempotencyKey(ctx, orderTypeKey)
		if err := upsertOrderTypeWorkflow(orderTypeCtx, order.OrderType); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando tipo de orden", "error", err)
			msg.Nak()
			return
		}

		// Delivery Units
		deliveryUnitsKey, err := canonicaljson.HashKey(ctx, "delivery_units", order.DeliveryUnits)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para unidades de entrega", "error", err)
			msg.Nak()
			return
		}
		deliveryUnitsCtx := sharedcontext.WithIdempotencyKey(ctx, deliveryUnitsKey)
		if err := upsertDeliveryUnitsWorkflow(deliveryUnitsCtx, order.DeliveryUnits); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando unidades de entrega", "error", err)
			msg.Nak()
			return
		}

		// Order References
		orderReferencesKey, err := canonicaljson.HashKey(ctx, "order_references", order)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para referencias de orden", "error", err)
			msg.Nak()
			return
		}
		orderReferencesCtx := sharedcontext.WithIdempotencyKey(ctx, orderReferencesKey)
		if err := upsertOrderReferencesWorkflow(orderReferencesCtx, order); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando referencias de orden", "error", err)
			msg.Nak()
			return
		}

		// Order Delivery Units
		orderDeliveryUnitsKey, err := canonicaljson.HashKey(ctx, "order_delivery_units", order)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para relación orden-unidades", "error", err)
			msg.Nak()
			return
		}
		orderDeliveryUnitsCtx := sharedcontext.WithIdempotencyKey(ctx, orderDeliveryUnitsKey)
		if err := upsertOrderDeliveryUnitsWorkflow(orderDeliveryUnitsCtx, order); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando relación orden-unidades", "error", err)
			msg.Nak()
			return
		}

		// Order
		orderKey, err := canonicaljson.HashKey(ctx, "order", order)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para orden", "error", err)
			msg.Nak()
			return
		}
		orderCtx := sharedcontext.WithIdempotencyKey(ctx, orderKey)
		if err := upsertOrderWorkflow(orderCtx, order); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando upsert de orden", "error", err)
			msg.Nak()
			return
		}

		// Delivery Units Labels
		deliveryUnitsLabelsKey, err := canonicaljson.HashKey(ctx, "delivery_units_labels", order)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para labels de unidades de entrega", "error", err)
			msg.Nak()
			return
		}
		deliveryUnitsLabelsCtx := sharedcontext.WithIdempotencyKey(ctx, deliveryUnitsLabelsKey)
		if err := upsertDeliveryUnitsLabelsWorkflow(deliveryUnitsLabelsCtx, order); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando labels de unidades de entrega", "error", err)
			msg.Nak()
			return
		}

		// Delivery Units Skills
		deliveryUnitsSkillsKey, err := canonicaljson.HashKey(ctx, "delivery_units_skills", order)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para skills de unidades de entrega", "error", err)
			msg.Nak()
			return
		}
		deliveryUnitsSkillsCtx := sharedcontext.WithIdempotencyKey(ctx, deliveryUnitsSkillsKey)
		if err := upsertDeliveryUnitsSkillsWorkflow(deliveryUnitsSkillsCtx, order); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando skills de unidades de entrega", "error", err)
			msg.Nak()
			return
		}

		// Size Categories y Skills
		for _, du := range order.DeliveryUnits {
			sizeCategoryKey, err := canonicaljson.HashKey(ctx, "size_category", du.SizeCategory)
			if err != nil {
				obs.Logger.ErrorContext(ctx, "Error generando key para size category", "error", err)
				msg.Nak()
				return
			}
			sizeCategoryCtx := sharedcontext.WithIdempotencyKey(ctx, sizeCategoryKey)
			if err := upsertSizeCategoryWorkflow(sizeCategoryCtx, du.SizeCategory); err != nil {
				obs.Logger.ErrorContext(ctx, "Error procesando size category", "error", err)
				msg.Nak()
				return
			}
			for _, skill := range du.Skills {
				skillKey, err := canonicaljson.HashKey(ctx, "skill", skill)
				if err != nil {
					obs.Logger.ErrorContext(ctx, "Error generando key para skill", "error", err)
					msg.Nak()
					return
				}
				skillCtx := sharedcontext.WithIdempotencyKey(ctx, skillKey)
				if err := upsertSkillWorkflow(skillCtx, skill); err != nil {
					obs.Logger.ErrorContext(ctx, "Error procesando skill", "error", err)
					msg.Nak()
					return
				}
			}
		}

		// Delivery Units History
		plan := domain.Plan{Routes: []domain.Route{{Orders: []domain.Order{order}}}}
		deliveryUnitsHistoryKey, err := canonicaljson.HashKey(ctx, "delivery_units_history", plan)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para historial de unidades de entrega", "error", err)
			msg.Nak()
			return
		}
		deliveryUnitsHistoryCtx := sharedcontext.WithIdempotencyKey(ctx, deliveryUnitsHistoryKey)
		if err := upsertDeliveryUnitsHistoryWorkflow(deliveryUnitsHistoryCtx, plan); err != nil {
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
