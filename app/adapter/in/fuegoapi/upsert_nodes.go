package fuegoapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/gcppublisher"
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
		upsertNodes,
		httpserver.New,
		gcppublisher.NewApplicationEvents,
		observability.NewObservability)
}

func upsertNodes(
	s httpserver.Server,
	publish gcppublisher.ApplicationEvents,
	obs observability.Observability) {
	fuego.Post(s.Manager, "/nodes",
		func(c fuego.ContextWithBody[request.UpsertNodeRequest]) (response.UpsertNodeResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "upsertNodes")
			defer span.End()

			requestBody, err := c.Body()
			if err != nil {
				return response.UpsertNodeResponse{}, err
			}

			eventPayload, _ := json.Marshal(requestBody)

			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "node",
					EventType:  "nodeSubmitted",
				})

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return response.UpsertNodeResponse{}, fuego.HTTPError{
					Title:  "error submitting node",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx,
				"NODE_SUBMITTED",
				slog.Any("payload", requestBody))

			return response.UpsertNodeResponse{
				Message: "Node submitted successfully",
			}, nil
		},
		option.Summary("upsert node"),
		option.Header("tenant", "api tenant", param.Required()),
		option.Tags(tagNodes))
}
