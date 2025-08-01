package fuegoapi

import (
	"encoding/json"
	"log/slog"
	"transport-app/app/adapter/out/natspublisher"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"
	"webhooks/app/adapter/in/fuegoapi/model"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		fleetOptimizedWebhookReception,
		httpserver.New,
		natspublisher.NewApplicationEvents,
		observability.NewObservability)
}

func fleetOptimizedWebhookReception(
	s httpserver.Server,
	publish natspublisher.ApplicationEvents,
	obs observability.Observability,
) {
	fuego.Post(s.Manager, "/webhooks/fleet-optimized",
		func(c fuego.ContextWithBody[model.FleetOptimizedWebhookBody]) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "webhookReception")
			defer span.End()

			requestBody, err := c.Body()
			if err != nil {
				return nil, err
			}

			eventPayload, _ := json.Marshal(requestBody)

			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "plan",
					EventType:  "fleetOptimizedWebhook",
				})
			accessToken := c.Header("X-Access-Token")
			eventCtx = sharedcontext.WithAccessToken(eventCtx, accessToken)
			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return nil, err
			}

			obs.Logger.InfoContext(spanCtx,
				"FLEET_OPTIMIZED_WEBHOOK_RECEIVED",
				slog.Any("payload", requestBody))

			return "ok", nil
		}, option.Summary("fleet optimized"), option.Tags("webhook reception"))
}
