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
		checkout,
		httpserver.New,
		tidbrepository.NewEnsureOrganizationForCountry,
		gcppublisher.NewApplicationEvents)
}
func checkout(
	s httpserver.Server,
	ensureOrg tidbrepository.EnsureOrganizationForCountry,
	outbox gcppublisher.ApplicationEvents) {
	fuego.Post(s.Manager, "/checkouts",
		func(c fuego.ContextWithBody[request.CheckoutRequest]) (response.CheckoutResponse, error) {
			requestBody, err := c.Body()
			if err != nil {
				return response.CheckoutResponse{}, err
			}

			organization := domain.Organization{}
			organization.SetKey(c.Header("organization"))
			requestBodyBytes, _ := json.Marshal(requestBody)
			if err := outbox(c.Context(), domain.Outbox{
				Attributes: map[string]string{
					"entityType":   "checkout",
					"eventType":    "checkoutSubmitted",
					"country":      organization.Country.Alpha2(),
					"organization": organization.GetOrgKey(),
					"consumer":     c.Header("consumer"),
					"commerce":     c.Header("commerce"),
				},
				Status:       "pending",
				Organization: organization,
				Payload:      requestBodyBytes,
			}); err != nil {
				return response.CheckoutResponse{}, err
			}
			return response.CheckoutResponse{
				Message: "checkout submission succedded",
			}, nil
		},
		option.Summary("checkout"),
		option.Tags(tagCheckouts),
		option.Tags(tagEndToEndOperator),
	)
}
