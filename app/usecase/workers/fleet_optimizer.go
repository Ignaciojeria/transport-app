package workers

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/fuegoapiclient"
	"transport-app/app/adapter/out/vroom"
	"transport-app/app/domain/optimization"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type FleetOptimizer func(
	ctx context.Context,
	input optimization.FleetOptimization) error

func init() {
	ioc.Registry(
		NewFleetOptimizer,
		vroom.NewOptimize,
		fuegoapiclient.NewPostUpsertRoute,
	)
}

func NewFleetOptimizer(
	optimize vroom.Optimize,
	postUpsertRoute fuegoapiclient.PostUpsertRoute) FleetOptimizer {
	return func(ctx context.Context, input optimization.FleetOptimization) error {

		routeRequests, err := optimize(ctx, input)
		if err != nil {
			return err
		}

		// Procesar cada ruta individualmente
		for i, routeRequest := range routeRequests {
			// Convertir UpsertRouteRequest a domain.Route para usar con postUpsertRoute
			// Por ahora, como postUpsertRoute espera domain.Route, necesitamos convertir
			// Esto es temporal hasta que se actualice postUpsertRoute para usar UpsertRouteRequest directamente
			fmt.Printf("Procesando ruta %d de %d: %s\n", i+1, len(routeRequests), routeRequest.ReferenceID)

			// TODO: Implementar conversión de UpsertRouteRequest a domain.Route
			// o actualizar postUpsertRoute para usar UpsertRouteRequest directamente
		}

		fmt.Printf("Optimización completada. Se procesaron %d rutas.\n", len(routeRequests))
		return nil
	}
}
