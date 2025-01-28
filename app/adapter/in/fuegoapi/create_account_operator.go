package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(createAccountOperator, httpserver.New)
}
func createAccountOperator(s httpserver.Server) {
	fuego.Post(s.Manager, "/operator",
		func(c fuego.ContextWithBody[request.CreateAccountOperatorRequest]) (any, error) {

			return "unimplemented", nil
		},
		option.Summary("createAccountOperator"),
		option.Tags(tagAccounts),
		option.Tags(tagEndToEndOperator),
	)
}
