package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(searchPlannerByEmail, httpserver.New)
}
func searchPlannerByEmail(s httpserver.Server) {
	/*
		fuego.Get(s.Manager, "/account/planner",
			func(c fuego.ContextNoBody) (any, error) {

				return "unimplemented", nil
			},
			option.Summary("searchPlannerByEmail"),
			option.Tags(tagAccounts))
	*/
}
