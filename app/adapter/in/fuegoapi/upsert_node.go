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
		upsertNode, httpserver.New,
		tidbrepository.NewEnsureOrganizationForCountry,
		gcppublisher.NewApplicationEvents,
	)
}
func upsertNode(
	s httpserver.Server,
	ensureOrg tidbrepository.EnsureOrganizationForCountry,
	outbox gcppublisher.ApplicationEvents) {
	fuego.Post(s.Manager, "/nodes",
		func(c fuego.ContextWithBody[request.UpsertNodeRequest]) (response.UpsertNodeResponse, error) {

			requestBody, err := c.Body()
			if err != nil {
				return response.UpsertNodeResponse{}, err
			}

			organization := domain.Organization{}
			organization.SetKey(c.Header("organization"))
			if err != nil {
				return response.UpsertNodeResponse{}, fuego.HTTPError{
					Title:  "error creating order",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			requestBodyBytes, _ := json.Marshal(requestBody)
			if err := outbox(c.Context(), domain.Outbox{
				Attributes: map[string]string{
					"entityType":   "node",
					"eventType":    "nodeSubmitted",
					"country":      organization.Country.Alpha2(),
					"organization": organization.GetOrgKey(),
					"consumer":     c.Header("consumer"),
					"commerce":     c.Header("commerce"),
					"referenceID":  requestBody.ReferenceID,
				},
				Status:  "pending",
				Payload: requestBodyBytes,
			}); err != nil {
				return response.UpsertNodeResponse{}, err
			}
			return response.UpsertNodeResponse{
				Message: "upsert node submitted",
			}, nil
		},
		option.Summary("upsertNode"),
		option.Header("consumer", "api consumer key", param.Required()),
		option.Header("commerce", "api commerce key", param.Required()),
		option.Tags(tagNetwork),
		option.Tags(tagEndToEndOperator))
}
