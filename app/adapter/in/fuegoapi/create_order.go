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

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(
		createOrder,
		httpserver.New,
		tidbrepository.NewEnsureOrganizationForCountry,
		gcppublisher.NewApplicationEvents,
		observability.NewObservability)
}
func createOrder(
	s httpserver.Server,
	ensureOrg tidbrepository.EnsureOrganizationForCountry,
	saveOutboxTrx gcppublisher.ApplicationEvents,
	obs observability.Observability) {
	fuego.Post(s.Manager, "/orders",
		func(c fuego.ContextWithBody[request.UpsertOrderRequest]) (response.UpsertOrderResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "createOrder")
			defer span.End()
			requestBody, err := c.Body()
			if err != nil {
				return response.UpsertOrderResponse{}, err
			}
			mappedTO := requestBody.Map()
			mappedTO.Organization.SetKey(c.Header("organization"))
			mappedTO.Headers.Consumer = c.Header("consumer")
			mappedTO.Headers.Commerce = c.Header("commerce")
			if c.Header("consumer") == "" {
				mappedTO.Headers.Consumer = "UNSPECIFIED"
			}
			if c.Header("commerce") == "" {
				mappedTO.Headers.Commerce = "UNSPECIFIED"
			}
			if err := mappedTO.Validate(); err != nil {
				return response.UpsertOrderResponse{}, fuego.HTTPError{
					Title:  "error creating order",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}
			eventPayload, _ := json.Marshal(requestBody)
			if err := saveOutboxTrx(spanCtx, domain.Outbox{
				Attributes: map[string]string{
					"entityType":   "order",
					"eventType":    "orderSubmitted",
					"country":      mappedTO.Organization.Country.Alpha2(),
					"organization": mappedTO.Organization.GetOrgKey(),
					"consumer":     c.Header("consumer"),
					"commerce":     c.Header("commerce"),
					"referenceID":  requestBody.ReferenceID,
				},
				Payload: eventPayload,
				Status:  "pending",
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
		}, option.Summary("createOrder"),
		option.Header("organization", "api organization key", param.Required()),
		option.Header("consumer", "api consumer key", param.Required()),
		option.Header("commerce", "api commerce key", param.Required()),
		option.Tags(tagOrders),
		option.Tags(tagEndToEndOperator),
	)
}
