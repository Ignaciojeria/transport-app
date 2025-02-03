package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		createPlan,
		httpserver.New)
}
func createPlan(s httpserver.Server) {
	/*
		fuego.Post(s.Manager, "/plans",
			func(c fuego.ContextWithBody[request.UpsertDailyPlanRequest]) (any, error) {

				return "unimplemented", nil
			},
			option.Summary("createPlan"),
			option.Tags(tagPlanning),
		)*/
}
