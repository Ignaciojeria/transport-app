package menu

import (
	"context"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
)

type OnUserMenusInsertedWebhook func(ctx context.Context, input events.UserMenusInsertedWebhook) error

func init() {
	ioc.Registry(NewOnUserMenusInsertedWebhook, eventprocessing.NewPublisherStrategy)
}

func NewOnUserMenusInsertedWebhook(publisherManager eventprocessing.PublisherManager) OnUserMenusInsertedWebhook {
	return func(ctx context.Context, wh events.UserMenusInsertedWebhook) error {
		// Generar version_id si no viene en el contexto
		var versionID string
		if existingVersionID, ok := sharedcontext.VersionIDFromContext(ctx); ok && existingVersionID != "" {
			versionID = existingVersionID
		} else {
			versionID = uuid.New().String()
			ctx = sharedcontext.WithVersionID(ctx, versionID)
		}

		menuCreateRequest := events.MenuCreateRequest{
			ID:         wh.Record.MenuID,
			CoverImage: "https://storage.googleapis.com/micartapro-menus/core/micartaprov3.webp",
			BusinessInfo: events.BusinessInfo{
				BusinessName:  "cadorago",
				Whatsapp:      "+56957857558",
				BusinessHours: []string{},
			},
			Menu: []events.MenuCategory{},
			DeliveryOptions: []events.DeliveryOption{
				{
					Type:        events.DeliveryOptionDelivery,
					RequireTime: false,
				},
				{
					Type:            events.DeliveryOptionPickup,
					RequireTime:     true,
					TimeRequestType: events.TimeRequestWindow,
					TimeWindows: []events.TimeWindow{
						{
							Start: "09:00",
							End:   "23:59",
						},
					},
				},
			},
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
