package eventprocessing

import (
	"context"
	"encoding/json"
	"fmt"
	"micartapro/app/shared/constants"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/ioc"
	"go.opentelemetry.io/otel/trace"

	"cloud.google.com/go/pubsub/v2"
)

// =============================================================
// GCP IMPLEMENTATION OF THE EVENT BUS
// =============================================================
type GcpPublisherManager struct {
	client *pubsub.Client
}

// IoC: register this manager + the GCP client
func init() {
	ioc.Register(NewGcpPublisherManager)
}

// Constructor for the GCP Publisher (no topic defined here)
func NewGcpPublisherManager(c *pubsub.Client) PublisherManager {
	return &GcpPublisherManager{client: c}
}

// =============================================================
// Publish → dynamic topic, CloudEvents-encoded payload
// =============================================================
func (p *GcpPublisherManager) Publish(
	ctx context.Context,
	request PublishRequest,
) error {

	// DomainEvent → CloudEvent
	ce := request.Event.ToCloudEvent(request.Source)

	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()
	if sc.IsValid() {
		ce.SetExtension(constants.CloudEventExtensionTraceID, sc.TraceID().String())
		ce.SetExtension(constants.CloudEventExtensionSpanID, sc.SpanID().String())
	}

	// Agregar idempotency key del contexto como extensión del CloudEvent
	if idempotencyKey, ok := sharedcontext.IdempotencyKeyFromContext(ctx); ok {
		ce.SetExtension(constants.CloudEventExtensionIdempotencyKey, idempotencyKey)
	}

	// Agregar user ID del contexto como extensión del CloudEvent
	if userID, ok := sharedcontext.UserIDFromContext(ctx); ok {
		ce.SetExtension(constants.CloudEventExtensionUserID, userID)
	}

	// Agregar version ID del contexto como extensión del CloudEvent
	if versionID, ok := sharedcontext.VersionIDFromContext(ctx); ok {
		ce.SetExtension(constants.CloudEventExtensionVersionID, versionID)
	}

	// Encode CloudEvent JSON
	bytes, err := json.Marshal(ce)
	if err != nil {
		return err
	}

	// Build attributes for Pub/Sub filtering
	// Incluimos campos principales del CloudEvent + extensiones para filtrado eficiente
	attrs := make(map[string]string)

	// Campos principales del CloudEvent (útil para filtrado sin deserializar)
	// Usamos guiones (-) siguiendo la convención CloudEvents estándar
	if ce.Type() != "" {
		attrs["ce-type"] = ce.Type()
	}
	if ce.Source() != "" {
		attrs["ce-source"] = ce.Source()
	}
	if ce.Subject() != "" {
		attrs["ce-subject"] = ce.Subject()
	}
	if ce.ID() != "" {
		attrs["ce-id"] = ce.ID()
	}

	// Extensiones del CloudEvent (para filtrado adicional)
	for k, v := range ce.Context.GetExtensions() {
		attrs[k] = fmt.Sprintf("%v", v)
	}

	// Dynamic topic resolution (very cheap)
	pubTopic := p.client.Publisher(request.Topic)
	pubTopic.EnableMessageOrdering = true

	// Publish message
	_, err = pubTopic.Publish(ctx, &pubsub.Message{
		Data:        bytes,
		Attributes:  attrs,
		OrderingKey: request.OrderingKey,
	}).Get(ctx)

	return err
}
