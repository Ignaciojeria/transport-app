package fuegoapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"transport-app/app/adapter/out/natspublisher"
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
		fleetOptimizationAgent,
		httpserver.New,
		natspublisher.NewApplicationEvents,
		observability.NewObservability)
}
func fleetOptimizationAgent(s httpserver.Server, publish natspublisher.ApplicationEvents, obs observability.Observability) {
	fuego.Post(s.Manager, "/agents/fleet-optimization",
		func(c fuego.ContextNoBody) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "optimization")
			defer span.End()

			spanCtx = sharedcontext.WithAccessToken(spanCtx, c.Header("X-Access-Token"))

			requestBody, err := c.Body()
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error getting request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			eventPayload, _ := json.Marshal(requestBody)

			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "optimization",
					EventType:  "agentOptimizationRequested",
				})

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return nil, fuego.HTTPError{
					Title:  "error requesting optimization",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx,
				"OPTIMIZATION_REQUEST_SUBMITTED",
				slog.Any("payload", requestBody))

			return "unimplemented", nil
		}, option.Summary("fleetOptimizationAgent"))
}
