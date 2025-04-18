package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchPlan func(ctx context.Context, osf domain.OrderSearchFilters) (domain.Plan, error)

func init() {
	ioc.Registry(
		NewSearchPlan,
		tidbrepository.NewFindOrdersByFilters,
	)
}

func NewSearchPlan(search tidbrepository.FindOrdersByFilters) SearchPlan {
	return func(ctx context.Context, osf domain.OrderSearchFilters) (domain.Plan, error) {
		/*
			orders, err := search(ctx, osf)
			if err != nil {
				return domain.Plan{}, err
			}

			if len(orders) == 0 {
				return domain.Plan{}, fmt.Errorf("no orders found")
			}

			// Agrupar órdenes por ruta
			routesMap := make(map[string]domain.Route)
			unassignedOrders := []domain.Order{}

			for _, order := range orders {
				routeKey := order.Plan.Routes[0].ReferenceID
				if route, exists := routesMap[routeKey]; exists {
					route.Orders = append(route.Orders, order)
					routesMap[routeKey] = route
				} else {
					routesMap[routeKey] = domain.Route{
						PlanID:      0, // Se asignará después
						ReferenceID: routeKey,
						Orders:      []domain.Order{order},
						Destination: order.Destination,
					}
				}
			}

			// Convertir mapa de rutas a slice
			routes := make([]domain.Route, 0, len(routesMap))
			for _, route := range routesMap {
				routes = append(routes, route)
			}

			// Crear el plan
			plan := domain.Plan{
				Origin:           orders[0].Plan.Origin,
				ReferenceID:      osf.PlanReferenceID,
				PlannedDate:      time.Now(),
				UnassignedOrders: unassignedOrders,
				Routes:           routes,
				PlanningStatus:   domain.PlanningStatus{Value: "IN_PROGRESS"},
				PlanType:         domain.PlanType{Value: "DEFAULT"},
			}

			return plan, nil
		*/
		return domain.Plan{}, nil
	}
}
