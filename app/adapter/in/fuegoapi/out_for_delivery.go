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
		outForDelivery,
		httpserver.New,
		gcppublisher.NewApplicationEvents)
}
func outForDelivery(
	s httpserver.Server,
	outbox gcppublisher.ApplicationEvents) {
	fuego.Post(s.Manager, "/out-for-delivery",
		func(c fuego.ContextWithBody[request.OutForDeliveryRequest]) (
			response.OutForDeliveryResponse, error) {
			requestBody, err := c.Body()
			if err != nil {
				return response.OutForDeliveryResponse{}, err
			}
			organization := domain.Organization{}
			organization.SetKey(c.Header("organization"))
			requestBodyBytes, _ := json.Marshal(requestBody)
			if err := outbox(c.Context(), domain.Outbox{
				Attributes: map[string]string{
					"entityType":   "outForDelivery",
					"eventType":    "outForDeliverySubmitted",
					"country":      organization.Country.Alpha2(),
					"organization": organization.GetOrgKey(),
					"consumer":     c.Header("consumer"),
					"commerce":     c.Header("commerce"),
				},
				Status:       "pending",
				Organization: organization,
				Payload:      requestBodyBytes,
			}); err != nil {
				return response.OutForDeliveryResponse{}, err
			}
			return response.OutForDeliveryResponse{
				Message: "out for delivery submission succedded",
			}, nil
		},
		option.Summary("outForDelivery"),
		option.Tags(tagDelivery),
	)
}
