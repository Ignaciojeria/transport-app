package natspublisher

import (
	"context"
	"encoding/json"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/sharedcontext"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ApplicationEvents func(ctx context.Context, outbox domain.Outbox) error

func init() {
	ioc.Registry(
		NewApplicationEvents,
		natsconn.NewJetStream,
		configuration.NewConf)
}

func NewApplicationEvents(
	js jetstream.JetStream,
	conf configuration.Conf) ApplicationEvents {
	return func(ctx context.Context, outbox domain.Outbox) error {
		// Crear el mismo formato que Google Pub/Sub
		msg := &pubsub.Message{
			Attributes: map[string]string{},
			Data:       outbox.Payload,
		}

		// 📦 Propagar baggage y trace context al pubsub message
		otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(msg.Attributes))
		sharedcontext.CopyBaggageToAttributesCamelCase(ctx, msg.Attributes)

		// Agregar atributos del outbox
		for k, v := range outbox.Attributes {
			msg.Attributes[k] = v
		}

		// Serializar el pubsub.Message a JSON
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		// Crear headers de NATS para filtrado adicional usando el baggage
		headers := nats.Header{}

		// Extraer EventType del baggage
		var eventType string
		if value, exists := sharedcontext.GetEventTypeFromContext(ctx); exists {
			eventType = value
			headers.Set("eventType", eventType)
		}

		var tenant string
		// Extraer TenantID del baggage
		if tenantID := sharedcontext.TenantIDFromContext(ctx); tenantID != uuid.Nil {
			tenant = tenantID.String()
		}

		// Extraer TenantCountry del baggage
		if tenantCountry := sharedcontext.TenantCountryFromContext(ctx); tenantCountry != "" {
			tenant = tenant + "-" + tenantCountry
			headers.Set("tenant", tenant)
		}

		// Extraer Channel del baggage
		if channel := sharedcontext.ChannelFromContext(ctx); channel != "" {
			headers.Set("channel", channel)
		}

		// Extraer EntityType del baggage
		if entityType := sharedcontext.EntityTypeFromContext(ctx); entityType != "" {
			headers.Set("entityType", entityType)
		}

		if tenant == "" {
			tenant = "no-tenant"
		}

		// Publicar el mensaje serializado a JetStream
		_, err = js.PublishMsg(ctx, &nats.Msg{
			Subject: conf.TRANSPORT_APP_TOPIC + "." + tenant + "." + eventType,
			Header:  headers,
			Data:    msgBytes,
		})
		return err
	}
}
