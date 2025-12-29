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
			ID:          wh.Record.MenuID,
			CoverImage:  "https://storage.googleapis.com/micartapro-menus/menus/01KCW67YKSV455GBVDT88S4072/gallery/portadav2.webp",
			FooterImage: "https://storage.googleapis.com/micartapro-menus/menus/01KCW67YKSV455GBVDT88S4072/gallery/logov2.webp",
			BusinessInfo: events.BusinessInfo{
				BusinessName: "cadorago",
				Whatsapp:     "+56957857558",
				BusinessHours: []string{
					"Lunes a Martes: 9h a 16h",
					"Miércoles, Jueves, Sábado y Domingo: hasta las 20h",
					"Viernes: Cerrado",
				},
			},
			Menu: []events.MenuCategory{
				{
					Title: "Menú",
					Items: []events.MenuItem{
						{
							Title: "Pollo a la plancha",
							Sides: []events.Side{
								{Name: "Con puré", Price: 3000},
								{Name: "Con arroz", Price: 3000},
							},
						},
						{
							Title: "Completo italiano",
							Price: 2500,
						},
						{
							Title: "Hamburguesa",
							Sides: []events.Side{
								{Name: "Sola", Price: 4100},
							},
						},
						{
							Title: "chacareros",
							Price: 7000,
						},
					},
				},
				{
					Title: "Postres",
					Items: []events.MenuItem{
						{
							Title: "mote con huesillo",
							Price: 4000,
						},
						{
							Title: "leche asada",
							Price: 3000,
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
