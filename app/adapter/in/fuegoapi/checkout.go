package fuegoapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/gcppublisher"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
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

			organization := domain.Organization{
				Key:     c.Header("organization-key"),
				Country: countries.ByName(c.Header("country")),
			}

			org, err := ensureOrg(c.Context(), organization)
			if err != nil {
				return response.CheckoutResponse{}, fuego.HTTPError{
					Title:  "error creating order",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			requestBodyBytes, _ := json.Marshal(requestBody)
			orgIDString := strconv.FormatInt(org.OrganizationCountryID, 10)
			if err := outbox(c.Context(), domain.Outbox{
				Attributes: map[string]string{
					"entityType":            "checkout",
					"eventType":             "checkoutSubmitted",
					"country":               countries.ByName(c.Header("country")).Alpha2(),
					"organizationCountryID": orgIDString,
					"consumer":              c.Header("consumer"),
					"commerce":              c.Header("commerce"),
				},
				Status:       "pending",
				Organization: org,
				Payload:      requestBodyBytes,
			}); err != nil {
				return response.CheckoutResponse{}, err
			}
			return response.CheckoutResponse{
				Message: "checkout submission succedded",
			}, nil
		},
		option.Summary("ordersCheckout"),
		option.Tags(tagOrders),
		option.Tags(tagEndToEndOperator),
	)
}
