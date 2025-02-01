package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(createAccountPlanner, httpserver.New)
}
func createAccountPlanner(s httpserver.Server) {
	/*
		fuego.Post(s.Manager, "/account/planner",
			func(c fuego.ContextNoBody) (any, error) {

				return "unimplemented", nil
			},
			option.Summary("createAccountPlanner"),
			option.Tags(tagAccounts),
		)*/
}
