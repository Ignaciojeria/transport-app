package fuegoapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/gcppublisher"
	"transport-app/app/adapter/out/tidbrepository"
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
		upsertOrder,
		httpserver.New,
		tidbrepository.NewEnsureOrganizationForCountry,
		gcppublisher.NewApplicationEvents,
		observability.NewObservability)
}
func upsertOrder(
	s httpserver.Server,
	ensureOrg tidbrepository.EnsureOrganizationForCountry,
	publish gcppublisher.ApplicationEvents,
	obs observability.Observability) {
	fuego.Post(s.Manager, "/orders",
		func(c fuego.ContextWithBody[request.UpsertOrderRequest]) (response.UpsertOrderResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "upsertOrder")
			defer span.End()
			requestBody, err := c.Body()
			if err != nil {
				return response.UpsertOrderResponse{}, err
			}
			mappedTO := requestBody.Map(spanCtx)

			if err := mappedTO.Validate(); err != nil {
				return response.UpsertOrderResponse{}, fuego.HTTPError{
					Title:  "error creating order",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}
			eventPayload, _ := json.Marshal(requestBody)

			eventCtx := sharedcontext.AddEventContextToBaggage(spanCtx,
				sharedcontext.EventContext{
					EntityType: "order",
					EventType:  "orderSubmitted",
				})

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				return response.UpsertOrderResponse{}, fuego.HTTPError{
					Title:  "error creating order",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			obs.Logger.InfoContext(spanCtx,
				"ORDER_SUBMISSION_SUCCEEDED",
				slog.Any("payload", requestBody))

			return response.UpsertOrderResponse{
				Message: "Order submitted successfully",
				Status:  "pending",
			}, err
		}, option.Summary("upsertOrder"),
		option.Header("organization", "api organization key", param.Required()),
		option.Header("consumer", "api consumer key", param.Required()),
		option.Header("commerce", "api commerce key", param.Required()),
		option.Tags(tagOrders),
		option.Tags(tagEndToEndOperator),
	)
}
