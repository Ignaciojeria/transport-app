package fuegoapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/natspublisher"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(
		upsertWebhook,
		httpserver.New,
		natspublisher.NewApplicationEvents,
		observability.NewObservability)
}

func upsertWebhook(
	s httpserver.Server,
	publish natspublisher.ApplicationEvents,
	obs observability.Observability) {
	fuego.Post(s.Manager, "/webhooks",
		func(c fuego.ContextWithBody[request.UpsertWebhookRequest]) (response.UpsertWebhookResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "upsertWebhook")
			defer span.End()
			requestBody, err := c.Body()
			if err != nil {
				return response.UpsertWebhookResponse{}, err
			}
			mappedTO := requestBody.Map(spanCtx)

			if err := mappedTO.Validate(); err != nil {
				return response.UpsertWebhookResponse{}, fuego.HTTPError{
					Title:  "error creating webhook",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}
			eventPayload, _ := json.Marshal(requestBody)

			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "webhook",
					EventType:  "webhookSubmitted",
				})

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return response.UpsertWebhookResponse{}, fuego.HTTPError{
					Title:  "error creating webhook",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			obs.Logger.InfoContext(spanCtx,
				"WEBHOOK_SUBMITTED",
				slog.Any("payload", requestBody))

			return response.UpsertWebhookResponse{
				Message: "Webhook submitted successfully",
				Status:  "pending",
			}, err
		},
		option.Summary("upsert webhook"),
		option.Header("tenant", "api tenant (required only for local development)", param.Required()),
		option.Header("consumer", "api consumer key", param.Required()),
		option.Header("commerce", "api commerce key", param.Required()),
		option.Header("channel", "api channel key", param.Required()),
		option.Header("X-Access-Token", "api access token"),
		option.Tags("webhooks"))
}
