package gcppublisher

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/gcppubsub"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
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

		message := &pubsub.Message{
			Attributes: map[string]string{
				"eventType":   outbox.EventType,
				"entityType":  outbox.EntityType,
				"referenceId": outbox.ReferenceID,
				"createdAt":   outbox.CreatedAt,
				"updatedAt":   outbox.UpdatedAt,
			},
			Data: outbox.Payload,
		}

		result := topic.Publish(ctx, message)
		// Get the server-generated message ID.
		_, err := result.Get(ctx)

		if err != nil {
			return err
		}

		return nil
	}
}
