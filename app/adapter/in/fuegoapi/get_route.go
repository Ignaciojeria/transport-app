package fuegoapi

import (
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/out/storjbucket"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"storj.io/uplink"
)

func init() {
	ioc.Registry(getRoute, httpserver.New, storjbucket.NewTransportAppBucket)
}

func getRoute(s httpserver.Server, storjBucket *storjbucket.TransportAppBucket) {
	fuego.Get(s.Manager, "/routes/{id}",
		func(c fuego.ContextWithBody[request.UpsertRouteRequest]) (any, error) {
			ctx := c.Context()
			routeID := c.PathParam("id")

			// Generar token de acceso para Storj
			token, err := storjBucket.GenerateEphemeralToken(ctx, 10*time.Minute, uplink.Permission{
				AllowDownload: true,
			})
			if err != nil {
				return nil, fmt.Errorf("error generando token de acceso: %w", err)
			}

			// Descargar datos desde Storj usando el bucket
			data, err := storjBucket.DownloadWithToken(ctx, token, routeID)
			if err != nil {
				return nil, fmt.Errorf("error descargando ruta desde Storj: %w", err)
			}

			// Hacer unmarshal de los datos al UpsertRouteRequest
			var routeRequest request.UpsertRouteRequest
			if err := json.Unmarshal(data, &routeRequest); err != nil {
				return nil, fmt.Errorf("error deserializando datos de ruta: %w", err)
			}

			fmt.Println("get route controller call done! :D")
			return routeRequest, nil
		},
		option.Summary("get route"),
		option.Tags(tagRoutes))
}
