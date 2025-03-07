package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(signUp, httpserver.New)
}
func signUp(s httpserver.Server) {
	fuego.Post(s.Manager, "/sign-up",
		func(c fuego.ContextWithBody[request.SignUpRequest]) (any, error) {
			return "unimplemented", nil
		},
		option.Tags(tagAuthentication),
		option.Summary("signUp"))
}
