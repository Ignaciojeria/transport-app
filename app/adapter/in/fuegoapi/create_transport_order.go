package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/mapper"
	"transport-app/app/adapter/in/fuegoapi/model"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(
		createTransportOrder,
		httpserver.New,
		usecase.NewCreateTransportOrder)
}
func createTransportOrder(s httpserver.Server, createTo usecase.CreateTransportOrder) {
	fuego.Post(s.Manager, "/transport-order",
		func(c fuego.ContextWithBody[model.CreateTransportOrderRequest]) (model.CreateTransportOrderResponse, error) {
			requestBody, err := c.Body()
			if err != nil {
				return model.CreateTransportOrderResponse{}, err
			}
			mappedTO := mapper.MapCreateTransportOrderRequest(requestBody)
			mappedTO.Tenant.Organization = c.Header("organization")
			mappedTO.Tenant.Consumer = c.Header("consumer")
			mappedTO.Tenant.Commerce = c.Header("commerce")
			mappedTO.Tenant.Country = countries.ByName(c.Header("country"))
			createdTo, err := createTo(c.Context(), mappedTO)
			return model.CreateTransportOrderResponse{
				ID:      createdTo.ID,
				Message: "transport order created",
			}, err
		}, option.Summary("createTransportOrder"),
		option.Header("organization", "api organization key", param.Required()),
		option.Header("consumer", "api consumer key", param.Required()),
		option.Header("commerce", "api commerce key", param.Required()),
		option.Header("country", "api country", param.Required()),
	)
}
