package fuegoapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(
		createOrder,
		httpserver.New,
		tidbrepository.NewEnsureOrganizationForCountry,
		tidbrepository.NewSaveOrderOutbox,
		observability.NewObservability)
}
func createOrder(
	s httpserver.Server,
	ensureOrg tidbrepository.EnsureOrganizationForCountry,
	saveOutboxTrx tidbrepository.SaveOrderOutbox,
	obs observability.Observability) {
	fuego.Post(s.Manager, "/order",
		func(c fuego.ContextWithBody[request.UpsertOrderRequest]) (response.UpsertOrderResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "createOrder")
			defer span.End()
			requestBody, err := c.Body()
			if err != nil {
				return response.UpsertOrderResponse{}, err
			}
			mappedTO := requestBody.Map()
			mappedTO.Organization.Key = c.Header("organization-key")
			mappedTO.Organization.Country = countries.ByName(c.Header("country"))
			mappedTO.BusinessIdentifiers.Consumer = c.Header("consumer")
			mappedTO.BusinessIdentifiers.Commerce = c.Header("commerce")
			if c.Header("consumer") == "" {
				mappedTO.BusinessIdentifiers.Consumer = "UNSPECIFIED"
			}
			if c.Header("commerce") == "" {
				mappedTO.BusinessIdentifiers.Commerce = "UNSPECIFIED"
			}
			if err := mappedTO.Validate(); err != nil {
				return response.UpsertOrderResponse{}, fuego.HTTPError{
					Title:  "error creating order",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}
			org, err := ensureOrg(spanCtx, mappedTO.Organization)
			if err != nil {
				return response.UpsertOrderResponse{}, fuego.HTTPError{
					Title:  "error creating order",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			// Convierte el OrganizationCountryID a string
			orgIDString := strconv.FormatInt(org.OrganizationCountryID, 10)

			eventPayload, _ := json.Marshal(requestBody)
			if _, err := saveOutboxTrx(spanCtx, domain.Outbox{
				Attributes: map[string]string{
					"entityType":            "order",
					"eventType":             "orderSubmitted",
					"country":               countries.ByName(c.Header("country")).Alpha2(),
					"organizationCountryID": orgIDString,
					"consumer":              c.Header("consumer"),
					"commerce":              c.Header("commerce"),
					"referenceID":           requestBody.ReferenceID,
				},
				Payload:      eventPayload,
				Status:       "pending",
				Organization: org,
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
		option.Header("organization-key", "api organization key", param.Required()),
		option.Header("country", "api country", param.Required()),
		option.Header("consumer", "api consumer key"),
		option.Header("commerce", "api commerce key"),
		option.Tags(tagOrders),
	)
}
