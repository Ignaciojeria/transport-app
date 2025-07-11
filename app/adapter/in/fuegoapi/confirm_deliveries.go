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
		confirmDeliveries,
		httpserver.New,
		natspublisher.NewApplicationEvents,
		observability.NewObservability)
}
func confirmDeliveries(
	s httpserver.Server,
	publish natspublisher.ApplicationEvents,
	obs observability.Observability) {
	fuego.Post(s.Manager, "/orders/deliveries",
		func(c fuego.ContextWithBody[request.ConfirmDeliveriesRequest]) (response.ConfirmDeliveriesResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "confirmDeliveries")
			defer span.End()

			requestBody, err := c.Body()
			if err != nil {
				return response.ConfirmDeliveriesResponse{}, err
			}

			if err := requestBody.Validate(); err != nil {
				return response.ConfirmDeliveriesResponse{}, fuego.HTTPError{
					Title:  "error validating deliveries",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			eventPayload, _ := json.Marshal(requestBody)

			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "delivery",
					EventType:  "deliveriesSubmitted",
				})

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return response.ConfirmDeliveriesResponse{}, fuego.HTTPError{
					Title:  "error submitting deliveries",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			obs.Logger.InfoContext(spanCtx,
				"DELIVERY_SUBMITTED",
				slog.Any("payload", requestBody))

			return response.ConfirmDeliveriesResponse{
				Message: "Deliveries submitted successfully",
			}, nil
		},
		option.Summary("deliveries"),
		option.Header("tenant", "api tenant (required only for local development)", param.Required()),
		option.Header("channel", "api channel", param.Required()),
		option.Header("X-Access-Token", "api access token"),
		option.Tags(tagOrders),
	)
}
