package fuegoapi

import (
	"encoding/json"
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
		createAccountOperator,
		httpserver.New,
		tidbrepository.NewEnsureOrganizationForCountry,
		gcppublisher.NewApplicationEvents)
}
func createAccountOperator(
	s httpserver.Server,
	ensureOrg tidbrepository.EnsureOrganizationForCountry,
	outbox gcppublisher.ApplicationEvents) {
	fuego.Post(s.Manager, "/operators",
		func(c fuego.ContextWithBody[request.CreateAccountOperatorRequest]) (response.CreateAccountResponse, error) {
			requestBody, err := c.Body()
			if err != nil {
				return response.CreateAccountResponse{}, err
			}

			organization := domain.Organization{}
			organization.SetKey(c.Header("organization"))

			requestBodyBytes, _ := json.Marshal(requestBody)
			if err := outbox(c.Context(), domain.Outbox{
				Attributes: map[string]string{
					"entityType":   "operator",
					"eventType":    "operatorSubmitted",
					"country":      organization.Country.Alpha2(),
					"organization": organization.GetOrgKey(),
					"consumer":     c.Header("consumer"),
					"commerce":     c.Header("commerce"),
					"referenceID":  requestBody.ReferenceID,
				},
				Status:  "pending",
				Payload: requestBodyBytes,
			}); err != nil {
				return response.CreateAccountResponse{}, err
			}
			return response.CreateAccountResponse{
				Message: "operator submitted",
			}, err
		},
		option.Summary("createAccountOperator"),
		option.Tags(tagAccounts),
		option.Tags(tagEndToEndOperator),
	)
}
