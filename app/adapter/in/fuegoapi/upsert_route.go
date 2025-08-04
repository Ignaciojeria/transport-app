package fuegoapi

import (
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
	ioc.Registry(upsertRoute, httpserver.New)
	ioc.Registry(getRoute, httpserver.New, storjbucket.NewTransportAppBucket)
}

func upsertRoute(s httpserver.Server) {
	fuego.Post(s.Manager, "/routes",
		func(c fuego.ContextWithBody[request.UpsertRouteRequest]) (any, error) {
			fmt.Println("upsert route controller call done! :D")
			return "unimplemented", nil
		},
		option.Summary("upsert route"),
		option.Tags(tagRoutes))
}

func getRoute(s httpserver.Server, storjBucket *storjbucket.TransportAppBucket) {
	fuego.Get(s.Manager, "/routes/{id}",
		func(c fuego.ContextNoBody) (any, error) {
			ctx := c.Context()
			
			// Obtener el ID de la ruta desde los parámetros de la URL
			routeID := c.PathParam("id")
			if routeID == "" {
				return nil, fmt.Errorf("route ID is required")
			}

			// Generar token efímero para acceder al bucket
			token, err := storjBucket.GenerateEphemeralToken(ctx, time.Hour*1, uplink.Permission{
				AllowDownload: true,
			})
			if err != nil {
				return nil, fmt.Errorf("error generating ephemeral token: %w", err)
			}

			// Descargar la ruta desde el StorJ bucket usando el ID
			routeData, err := storjBucket.DownloadWithToken(ctx, token, routeID)
			if err != nil {
				return nil, fmt.Errorf("error downloading route from bucket: %w", err)
			}

			// Retornar los datos de la ruta (como JSON crudo o deserializar según necesidad)
			return string(routeData), nil
		},
		option.Summary("get route by ID"),
		option.Tags(tagRoutes))
}
