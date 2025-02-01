package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(searchDispatcherByEmail, httpserver.New)
}
func searchDispatcherByEmail(s httpserver.Server) {
	/*
		fuego.Get(s.Manager, "/account/dispatcher",
			func(c fuego.ContextNoBody) (any, error) {

				return "unimplemented", nil
			},
			option.Summary("searchDispatcherByEmail"),
			option.Tags(tagAccounts),
		)*/
}
