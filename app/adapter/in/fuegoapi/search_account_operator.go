package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		searchAccountOperator,
		httpserver.New)
}
func searchAccountOperator(s httpserver.Server) {
	fuego.Get(s.Manager, "/operator",
		func(c fuego.ContextNoBody) (response.SearchAccountsResponse, error) {
			return response.SearchAccountsResponse{}, nil
		},
		option.Summary("searchAccountOperator"),
		option.Tags(tagAccounts, tagEndToEndOperator),
		option.Query("email", "Filter By Operator Email"))
}
