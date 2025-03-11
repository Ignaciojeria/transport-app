package fuegoapi

import (
	"encoding/json"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/gcppublisher"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(
		upsertVehicle,
		httpserver.New,
		tidbrepository.NewEnsureOrganizationForCountry,
		gcppublisher.NewApplicationEvents)
}
func upsertVehicle(
	s httpserver.Server,
	ensureOrg tidbrepository.EnsureOrganizationForCountry,
	outbox gcppublisher.ApplicationEvents) {
	fuego.Post(s.Manager, "/vehicles",
		func(c fuego.ContextWithBody[request.UpsertVehicleRequest]) (response.UpsertVehicleResponse, error) {

			requestBody, err := c.Body()
			if err != nil {
				return response.UpsertVehicleResponse{}, err
			}

			organization := domain.Organization{}
			organization.SetKey(c.Header("organization"))
			org, err := ensureOrg(c.Context(), organization)
			if err != nil {
				return response.UpsertVehicleResponse{}, fuego.HTTPError{
					Title:  "error creating order",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			requestBodyBytes, _ := json.Marshal(requestBody)

			if err := outbox(c.Context(), domain.Outbox{
				Attributes: map[string]string{
					"entityType":   "vehicle",
					"eventType":    "vehicleSubmitted",
					"country":      organization.Country.Alpha2(),
					"organization": organization.GetOrgKey(),
					"consumer":     c.Header("consumer"),
					"commerce":     c.Header("commerce"),
					"referenceID":  requestBody.ReferenceID,
				},
				Status:       "pending",
				Organization: org,
				Payload:      requestBodyBytes,
			}); err != nil {
				return response.UpsertVehicleResponse{}, err
			}

			return response.UpsertVehicleResponse{
				Message: "upsert vehicle submitted",
			}, nil
		},
		option.Summary("upsertVehicle"),
		option.Header("consumer", "api consumer key", param.Required()),
		option.Header("commerce", "api commerce key", param.Required()),
		option.Tags(tagFleets),
	)
}
