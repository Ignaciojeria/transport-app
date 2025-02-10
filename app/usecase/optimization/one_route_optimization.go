package optimization

import (
	"context"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type oneRouteOptimization func(context.Context, domain.Plan) (domain.Plan, error)

func init() {
	ioc.Registry(newOneRouteOptimization)
}
func newOneRouteOptimization() oneRouteOptimization {
	return func(ctx context.Context, p domain.Plan) (domain.Plan, error) {
		for _, unnasignedOrder := range p.UnassignedOrders {
			p.Routes[0].Orders = append(p.Routes[0].Orders, unnasignedOrder)
		}
		return p, nil
	}
}
