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
)

func init() {
	ioc.Registry(
		optimization,
		httpserver.New,
		gcppublisher.NewApplicationEvents,
		observability.NewObservability)
}

func optimization(
	s httpserver.Server,
	publish gcppublisher.ApplicationEvents,
	obs observability.Observability) {
	fuego.Post(s.Manager, "/optimize",
		func(c fuego.ContextWithBody[request.OptimizationRequest]) (response.OptimizationResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "optimization")
			defer span.End()

			requestBody, err := c.Body()
			if err != nil {
				return response.OptimizationResponse{}, err
			}

			eventPayload, _ := json.Marshal(requestBody)

			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "optimization",
					EventType:  "optimizationRequested",
				})

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return response.OptimizationResponse{}, fuego.HTTPError{
					Title:  "error requesting optimization",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx,
				"OPTIMIZATION_REQUEST_SUBMITTED",
				slog.Any("payload", requestBody))

			return response.OptimizationResponse{
				TraceID: span.SpanContext().TraceID().String(),
			}, nil
		}, option.Summary("fleets planning optimization"), option.Tags("optimization"))
}
