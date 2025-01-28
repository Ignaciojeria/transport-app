package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(capacityKMeans, httpserver.New)
}
func capacityKMeans(s httpserver.Server) {
	fuego.Post(s.Manager, "/algorithm/capacity-k-means",
		func(c fuego.ContextNoBody) (any, error) {

			return "unimplemented", nil
		},
		option.Summary("capacityKMeans"),
		option.Tags(tagAlgorithm),
	)
}
