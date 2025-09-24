package fuegoapi

import (
	"fmt"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(authGoogleCallback, httpserver.New)
}
func authGoogleCallback(s httpserver.Server) {
	fuego.Post(s.Manager, "/auth/google/callback",
		func(c fuego.ContextNoBody) (any, error) {
			fmt.Println("authGoogleCallback works! :D")
			return "unimplemented", nil
		}, option.Summary("authGoogleCallback"))
}
