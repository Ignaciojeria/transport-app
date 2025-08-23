package fuegoapi

import (
	"encoding/json"
	"fmt"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		getRoute,
		httpserver.New,
		usecase.NewGetDataFromRedisWorkflow)
}

func getRoute(s httpserver.Server, getDataFromRedisWorkflow usecase.GetDataFromRedisWorkflow) {
	fuego.Get(s.Manager, "/routes/{id}",
		func(c fuego.ContextWithBody[request.UpsertRouteRequest]) (any, error) {
			ctx := c.Context()
			routeID := c.PathParam("id")

			// Descargar datos desde Redis
			data, err := getDataFromRedisWorkflow(ctx, routeID)
			if err != nil {
				return nil, fmt.Errorf("error descargando ruta desde Redis: %w", err)
			}

			// Hacer unmarshal de los datos al UpsertRouteRequest
			var routeRequest request.UpsertRouteRequest
			if err := json.Unmarshal(data, &routeRequest); err != nil {
				return nil, fmt.Errorf("error deserializando datos de ruta: %w", err)
			}

			fmt.Println("get route controller call done! :D")
			if c.Header("X-View-Mode") == "sheet-map" {
				return routeRequest.FlattenForExcel(), nil
			}
			return routeRequest, nil
		},
		option.Summary("get route"),
		option.Tags(tagRoutes))
}
