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
)

func init() {
	ioc.Registry(
		upsertPlan,
		httpserver.New,
		tidbrepository.NewEnsureOrganizationForCountry,
		gcppublisher.NewApplicationEvents)
}
func upsertPlan(
	s httpserver.Server,
	ensureOrg tidbrepository.EnsureOrganizationForCountry,
	outbox gcppublisher.ApplicationEvents) {
	fuego.Post(s.Manager, "/plans",
		func(c fuego.ContextWithBody[request.UpsertPlanRequest]) (response.UpsertPlanResponse, error) {
			requestBody, err := c.Body()
			if err != nil {
				return response.UpsertPlanResponse{}, err
			}

			organization := domain.Organization{}
			organization.SetKey(c.Header("organization"))
			if err != nil {
				return response.UpsertPlanResponse{}, fuego.HTTPError{
					Title:  "error creating daily plan",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			requestBodyBytes, _ := json.Marshal(requestBody)
			if err := outbox(c.Context(), domain.Outbox{
				Attributes: map[string]string{
					"entityType":   "plan",
					"eventType":    "upsertPlanSubmitted",
					"country":      organization.Country.Alpha2(),
					"organization": organization.GetOrgKey(),
					"consumer":     c.Header("consumer"),
					"commerce":     c.Header("commerce"),
					"referenceID":  requestBody.ReferenceID,
				},
				Status:  "pending",
				Payload: requestBodyBytes,
			}); err != nil {
				return response.UpsertPlanResponse{}, err
			}
			return response.UpsertPlanResponse{
				Message: "daily plan submitted",
			}, err
		},
		option.Summary("upsertPlan"),
		option.Tags(tagEndToEndOperator, tagPlanning),
	)
}
