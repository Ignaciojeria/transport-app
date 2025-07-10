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
		routeStarted,
		httpserver.New,
		natspublisher.NewApplicationEvents,
		observability.NewObservability)
}

func routeStarted(
	s httpserver.Server,
	publish natspublisher.ApplicationEvents,
	obs observability.Observability) {
	fuego.Post(s.Manager, "/routes/start",
		func(c fuego.ContextWithBody[request.RouteStartedRequest]) (response.RouteStartedResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "routeStarted")
			defer span.End()

			requestBody, err := c.Body()
			if err != nil {
				return response.RouteStartedResponse{}, err
			}

			eventPayload, _ := json.Marshal(requestBody)

			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "route",
					EventType:  "routeStarted",
				})

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return response.RouteStartedResponse{}, fuego.HTTPError{
					Title:  "error starting route",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx,
				"ROUTE_START_SUBMITTED",
				slog.Any("payload", requestBody))

			return response.RouteStartedResponse{
				Message: "Route started successfully",
			}, nil
		},
		option.Summary("route start"),
		option.Header("tenant", "api tenant (required only for local development)", param.Required()),
		option.Header("channel", "api channel", param.Required()),
		option.Header("X-Access-Token", "api access token"),
		option.Tags(tagRoutes))
}
