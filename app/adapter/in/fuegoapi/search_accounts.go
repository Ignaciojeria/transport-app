package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(searchAccounts, httpserver.New)
}
func searchAccounts(s httpserver.Server) {
	fuego.Post(s.Manager, "/account/search",
		func(c fuego.ContextWithBody[request.SearchAccountsRequest]) ([]response.SearchAccountsResponse, error) {
			return nil, nil
		},
		option.Summary("searchAccounts"),
		option.Tags(tagAccounts))
}
