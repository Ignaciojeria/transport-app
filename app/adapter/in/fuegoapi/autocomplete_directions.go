package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(autocompleteDirections, httpserver.New)
}
func autocompleteDirections(s httpserver.Server) {
	fuego.Get(s.Manager, "/autocomplete",
		func(c fuego.ContextNoBody) (any, error) {

			return "unimplemented", nil
		},
		option.Tags(tagLocations),
		option.Summary("autocompleteDirections"),
		option.Query("q", "address"),
		option.Query("limit", "5"),
		option.Query("dedupe", "1"),
	)
}
