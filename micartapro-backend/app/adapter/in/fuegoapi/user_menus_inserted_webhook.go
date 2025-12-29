package fuegoapi

import (
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(userMenusInsertedWebhook, httpserver.New,
		observability.NewObservability,
		eventprocessing.NewPublisherStrategy,
	)
}
func userMenusInsertedWebhook(
	s httpserver.Server,
	obs observability.Observability,
	publisherManager eventprocessing.PublisherManager) {
	fuego.Post(s.Manager, "/webhooks/user-menus-inserted",
		func(c fuego.ContextWithBody[events.UserMenusInsertedWebhook]) (any, error) {
			requestBody, _ := c.Body()
			requestBody.CreatedAtToISO8601()
			obs.Logger.Info("userMenusInsertedWebhook request received", "requestBody", requestBody)
			return "unimplemented", nil
		}, option.Summary("userMenusInsertedWebhook"), option.Tags("webhooks"))
}
