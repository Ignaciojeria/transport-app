package fuegoapi

import (
	"micartapro/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(userCreated, httpserver.New)
}
func userCreated(s httpserver.Server) {
	fuego.Post(s.Manager, "/webhook/user-created",
		func(c fuego.ContextNoBody) (any, error) {

			return "unimplemented", nil
		}, option.Summary("userCreated"))
}
