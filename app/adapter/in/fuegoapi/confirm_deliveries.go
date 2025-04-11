package fuegoapi

import (
	"encoding/json"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/gcppublisher"
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
		gcppublisher.NewApplicationEvents)
}
func checkout(
	s httpserver.Server,
	outbox gcppublisher.ApplicationEvents) {
	fuego.Post(s.Manager, "/confirm-deliveries",
		func(c fuego.ContextWithBody[request.ConfirmDeliveriesRequest]) (
			response.ConfirmDeliveriesResponse, error) {
			requestBody, err := c.Body()
			if err != nil {
				return response.ConfirmDeliveriesResponse{}, err
			}

			organization := domain.Organization{}
			organization.SetKey(c.Header("organization"))
			requestBodyBytes, _ := json.Marshal(requestBody)
			if err := outbox(c.Context(), domain.Outbox{
				Attributes: map[string]string{
					"entityType":   "confirmDeliveries",
					"eventType":    "confirmDeliveriesSubmitted",
					"country":      organization.Country.Alpha2(),
					"organization": organization.GetOrgKey(),
					"consumer":     c.Header("consumer"),
					"commerce":     c.Header("commerce"),
				},
				Status:  "pending",
				Payload: requestBodyBytes,
			}); err != nil {
				return response.ConfirmDeliveriesResponse{}, err
			}
			return response.ConfirmDeliveriesResponse{
				Message: "checkout submission succedded",
			}, nil
		},
		option.Summary("confirmDeliveries"),
		option.Tags(tagDelivery),
		option.Tags(tagEndToEndOperator),
	)
}
