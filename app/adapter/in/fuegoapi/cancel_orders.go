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
		cancelOrders,
		httpserver.New,
		natspublisher.NewApplicationEvents,
		observability.NewObservability)
}

func cancelOrders(
	s httpserver.Server,
	publish natspublisher.ApplicationEvents,
	obs observability.Observability) {
	fuego.Post(s.Manager, "/orders/cancel",
		func(c fuego.ContextWithBody[request.CancelOrdersRequest]) (response.CancelOrdersResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "cancelOrders")
			defer span.End()

			requestBody, err := c.Body()
			if err != nil {
				return response.CancelOrdersResponse{}, err
			}

			if err := requestBody.Validate(); err != nil {
				return response.CancelOrdersResponse{}, fuego.HTTPError{
					Title:  "error validating cancellations",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			eventPayload, _ := json.Marshal(requestBody)

			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "order",
					EventType:  "ordersCancellationSubmitted",
				})

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return response.CancelOrdersResponse{}, fuego.HTTPError{
					Title:  "error cancelling orders",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			obs.Logger.InfoContext(spanCtx,
				"ORDER_CANCELLATION_SUBMITTED",
				slog.Any("payload", requestBody))

			return response.CancelOrdersResponse{
				Message: "Orders cancellation submitted ok",
			}, nil
		},
		option.Summary("cancel order"),
		option.Header("tenant", "api tenant (required only for local development)", param.Required()),
		option.Header("channel", "api channel", param.Required()),
		option.Header("X-Access-Token", "api access token"),
		option.Tags(tagOrders),
	)
}
