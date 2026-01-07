package usecase

import (
	"context"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OnUserMenusInsertedWebhook func(ctx context.Context, input events.UserMenusInsertedWebhook) error

func init() {
	ioc.Registry(NewOnUserMenusInsertedWebhook, eventprocessing.NewPublisherStrategy)
}

func NewOnUserMenusInsertedWebhook(publisherManager eventprocessing.PublisherManager) OnUserMenusInsertedWebhook {
	return func(ctx context.Context, wh events.UserMenusInsertedWebhook) error {
		menuCreateRequest := events.MenuCreateRequest{
			ID:         wh.Record.MenuID,
			CoverImage: "https://storage.googleapis.com/micartapro-menus/core/micartaprov3.webp",
			BusinessInfo: events.BusinessInfo{
				BusinessName:  "cadorago",
				Whatsapp:      "+56957857558",
				BusinessHours: []string{},
			},
			Menu: []events.MenuCategory{},
		}

		err := publisherManager.Publish(ctx, eventprocessing.PublishRequest{
			Topic:       "micartapro.events",
			Source:      "micartapro.webhook.user.menus.inserted",
			OrderingKey: wh.Record.MenuID,
			Event:       menuCreateRequest,
		})
		if err != nil {
			return err
		}
		return nil
	}
}
