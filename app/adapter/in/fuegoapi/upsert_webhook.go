package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(upsertWebhook, httpserver.New)
}
func upsertWebhook(s httpserver.Server) {
	fuego.Post(s.Manager, "/webhooks",
		func(c fuego.ContextWithBody[request.UpsertWebhookRequest]) error {

			return nil
		},
		option.Summary("upsert webhook"),
		option.Tags("webhooks"))
}
