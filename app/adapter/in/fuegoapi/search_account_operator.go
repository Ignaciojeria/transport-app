package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(
		searchAccountOperator,
		httpserver.New,
		usecase.NewSearchAccountOperator)
}
func searchAccountOperator(
	s httpserver.Server,
	search usecase.SearchAccountOperator) {
	fuego.Get(s.Manager, "/operators",
		func(c fuego.ContextNoBody) (response.SearchAccountResponse, error) {
			operator := domain.Operator{
				Contact: domain.Contact{
					Email: c.QueryParam("email"),
				},
			}
			operator.Organization.SetKey(c.Header("organization"))
			operator, err := search(c.Context(), operator)
			if err != nil {
				return response.SearchAccountResponse{}, err
			}
			return response.MapSearchAccountOperatorResponse(operator), nil
		},
		option.Summary("searchAccountOperator"),
		option.Header("organization", "api organization key", param.Required()),
		option.Tags(tagAccounts, tagEndToEndOperator),
		option.Query("email", "Filter By Operator Email"))
}
