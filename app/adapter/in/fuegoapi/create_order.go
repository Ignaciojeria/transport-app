package fuegoapi

import (
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
		createOrder,
		httpserver.New,
		usecase.NewCreateOrder)
}
func createOrder(s httpserver.Server, createTo usecase.CreateOrder) {
	fuego.Post(s.Manager, "/order",
		func(c fuego.ContextWithBody[model.CreateOrderRequest]) (model.CreateOrderResponse, error) {
			requestBody, err := c.Body()
			if err != nil {
				return model.CreateOrderResponse{}, err
			}
			mappedTO := requestBody.Map()
			mappedTO.Organization.Email = "TODO" //c.Header("api-key")
			mappedTO.Organization.Name = "TODO"
			mappedTO.Organization.Country = countries.ByName(c.Header("country"))
			mappedTO.BusinessIdentifiers.Consumer = c.Header("consumer")
			mappedTO.BusinessIdentifiers.Commerce = c.Header("commerce")
			createdTo, err := createTo(c.Context(), mappedTO)
			return model.CreateOrderResponse{
				ID:      createdTo.ID,
				Message: "order created",
			}, err
		}, option.Summary("createOrder"),
		option.Header("organization-key", "api organization key", param.Required()),
		option.Header("country", "api country", param.Required()),
		option.Header("consumer", "api consumer key"),
		option.Header("commerce", "api commerce key"),
	)
}
