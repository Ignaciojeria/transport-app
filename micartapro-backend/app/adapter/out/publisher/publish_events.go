package publisher

import (
	"context"

	"micartapro-backend/app/shared/infrastructure/eventprocessing"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

// Alias conveniente
type PublishEvents func(ctx context.Context, evt eventprocessing.DomainEvent) error

func init() {
	ioc.Registry(
		NewPublishEvents,
		eventprocessing.NewPublisherStrategy,
	)
}

func NewPublishEvents(pm eventprocessing.PublisherManager) PublishEvents {
	return func(ctx context.Context, evt eventprocessing.DomainEvent) error {
		return pm.Publish(ctx, eventprocessing.PublishRequest{
			Topic:  "micartapro.events",
			Source: "micartapro.api.menu.created",
			Event:  evt,
		})
	}
}
