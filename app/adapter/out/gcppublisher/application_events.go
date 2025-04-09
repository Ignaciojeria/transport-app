package gcppublisher

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/gcppubsub"
	"transport-app/app/shared/sharedcontext"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ApplicationEvents func(ctx context.Context, outbox domain.Outbox) error

func init() {
	ioc.Registry(
		NewApplicationEvents,
		gcppubsub.NewClient,
		configuration.NewConf)
}
func NewApplicationEvents(
	c *pubsub.Client,
	conf configuration.Conf) ApplicationEvents {
	topicName := conf.OUTBOX_TOPIC_NAME
	topic := c.Topic(topicName)
	return func(ctx context.Context, outbox domain.Outbox) error {
		msg := &pubsub.Message{
			Attributes: map[string]string{},
			Data:       outbox.Payload,
		}
		// 📦 Propagar baggage y trace context al pubsub message
		otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(msg.Attributes))
		sharedcontext.CopyBaggageToAttributesCamelCase(ctx, msg.Attributes)
		result := topic.Publish(ctx, msg)
		_, err := result.Get(ctx)
		return err
	}
}
