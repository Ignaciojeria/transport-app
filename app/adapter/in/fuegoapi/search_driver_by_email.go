package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(searchDriverByEmail, httpserver.New)
}
func searchDriverByEmail(s httpserver.Server) {
	fuego.Get(s.Manager, "/account/driver",
		func(c fuego.ContextNoBody) (any, error) {

			return "unimplemented", nil
		},
		option.Summary("searchDriverByEmail"),
		option.Tags(tagAccounts))
}
