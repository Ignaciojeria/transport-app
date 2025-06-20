package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(picking, httpserver.New)
}
func picking(s httpserver.Server) {
	fuego.Post(s.Manager, "/seller/pickup",
		func(c fuego.ContextWithBody[request.PickingVisitConfirmedRequest]) (any, error) {
			return "unimplemented", nil
		}, option.Summary("confirm seller pickup"), option.Tags("logistic train"))
}
