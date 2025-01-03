package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/model"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(createAccount, httpserver.New)
}
func createAccount(s httpserver.Server) {
	fuego.Post(s.Manager, "/account",
		func(c fuego.ContextWithBody[model.CreateAccountRequest]) (model.CreateAccountResponse, error) {

			return model.CreateAccountResponse{
				Message: "account created",
			}, nil
		}, option.Summary("createAccount"))
}
