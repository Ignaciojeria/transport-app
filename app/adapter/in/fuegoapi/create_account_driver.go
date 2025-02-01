package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(createAccountDriver, httpserver.New)
}
func createAccountDriver(s httpserver.Server) {
	/*
		fuego.Post(s.Manager, "/account/driver",
			func(c fuego.ContextNoBody) (any, error) {

				return "unimplemented", nil
			},
			option.Summary("createAccountDriver"),
			option.Tags(tagAccounts),
		)*/
}
