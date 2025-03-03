package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(geocoding, httpserver.New)
}
func geocoding(s httpserver.Server) {
	fuego.Get(s.Manager, "/insert-your-custom-pattern-here",
		func(c fuego.ContextNoBody) (any, error) {

			return "unimplemented", nil
		},
		option.Tags(tagLocations),
		option.Summary("geocoding"))
}
