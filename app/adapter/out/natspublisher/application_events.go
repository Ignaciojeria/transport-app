package natspublisher

import (
	"context"
	"encoding/json"
	"transport-app/app/domain"
	"transport-app/app/shared/chunker"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase"

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
		configuration.NewConf,
		natsconn.NewKeyValue,
		usecase.NewStoreDataInRedisWorkflow,
	)
}

func NewApplicationEvents(
	js jetstream.JetStream,
	conf configuration.Conf,
	kv jetstream.KeyValue,
	storeDataInRedisWorkflow usecase.StoreDataInRedisWorkflow,
) ApplicationEvents {
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

		// Crear headers de NATS para filtrado adicional usando el baggage
		headers := nats.Header{}

		// Serializar el pubsub.Message a JSON
		// Si el evento es optimizationRequested, chunkear solo el payload y poner los IDs en Data
		eventType, eventTypeExists := sharedcontext.GetEventTypeFromContext(ctx)
		if eventTypeExists && eventType == "optimizationRequested" || eventType == "agentOptimizationRequested" {
			chunks, err := chunker.SplitBytes(outbox.Payload)
			if err != nil {
				return err
			}
			var chunkIDs []string

			for _, chunk := range chunks {
				chunkKey := chunk.ID.String()
				if err != nil {
					return err
				}
				err = storeDataInRedisWorkflow(ctx, chunkKey, chunk.Data)
				if err != nil {
					return err
				}
				chunkIDs = append(chunkIDs, chunkKey)
			}
			// El campo Data del pubsub.Message será el arreglo de IDs serializado
			msg.Data, err = json.Marshal(chunkIDs)
			if err != nil {
				return err
			}
		}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		// Extraer EventType del baggage
		if eventType, eventTypeExists := sharedcontext.GetEventTypeFromContext(ctx); eventTypeExists {
			headers.Set("eventType", eventType)
		}

		var tenant string
		// Extraer TenantID del baggage
		if tenantID := sharedcontext.TenantIDFromContext(ctx); tenantID != uuid.Nil {
			tenant = tenantID.String()
		}

		// Extraer TenantCountry del baggage
		var country string
		if tenantCountry := sharedcontext.TenantCountryFromContext(ctx); tenantCountry != "" {
			country = tenantCountry
			headers.Set("tenant", tenant+"-"+tenantCountry)
		}

		// Extraer Channel del baggage
		if channel := sharedcontext.ChannelFromContext(ctx); channel != "" {
			headers.Set("channel", channel)
		}

		// Extraer EntityType del baggage
		if entityType := sharedcontext.EntityTypeFromContext(ctx); entityType != "" {
			headers.Set("entityType", entityType)
		}

		if accessToken, exists := sharedcontext.AccessTokenFromContext(ctx); exists {
			headers.Set("X-Access-Token", accessToken)
		}

		if tenant == "" {
			tenant = "no-tenant"
		}

		if country == "" {
			country = "no-country"
		}
		// Publicar el mensaje serializado a JetStream
		_, err = js.PublishMsg(ctx, &nats.Msg{
			Subject: conf.TRANSPORT_APP_TOPIC + "." + conf.ENVIRONMENT + "." + tenant + "." + country + "." + eventType,
			Header:  headers,
			Data:    msgBytes,
		})
		return err
	}
}
