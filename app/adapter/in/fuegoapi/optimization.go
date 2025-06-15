package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/google/uuid"
)

func init() {
	ioc.Registry(optimization, httpserver.New)
}
func optimization(s httpserver.Server) {
	fuego.Post(s.Manager, "/optimize",
		func(c fuego.ContextWithBody[request.OptimizationRequest]) (response.OptimizationResponse, error) {

			return response.OptimizationResponse{
				OptimizationID: uuid.New().String(),
			}, nil
		}, option.Summary("optimization"), option.Tags("optimization"))
}
