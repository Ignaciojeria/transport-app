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

		outbox.Attributes["createdAt"] = outbox.CreatedAt
		outbox.Attributes["updatedAt"] = outbox.UpdatedAt
		message := &pubsub.Message{
			Attributes: map[string]string{
				"referenceId":           outbox.Attributes["referenceID"],
				"createdAt":             outbox.CreatedAt,
				"updatedAt":             outbox.UpdatedAt,
				"eventType":             outbox.Attributes["eventType"],
				"entityType":            outbox.Attributes["entityType"],
				"organizationCountryID": outbox.Attributes["organizationCountryID"],
				"commerce":              outbox.Attributes["commerce"],
				"consumer":              outbox.Attributes["consumer"],
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
