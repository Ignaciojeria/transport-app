package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(createAccountSeller, httpserver.New)
}
func createAccountSeller(s httpserver.Server) {
	/*
		fuego.Post(s.Manager, "/account/seller",
			func(c fuego.ContextNoBody) (any, error) {

				return "unimplemented", nil
			},
			option.Summary("createAccountSeller"),
			option.Tags(tagAccounts),
		)*/
}
