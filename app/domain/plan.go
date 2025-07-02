package domain

import (
	"context"
	"time"
)

type Plan struct {
	Headers
	ReferenceID        string
	UnassignedOrigins  []NodeInfo
	UnassignedVehicles []Vehicle
	UnassignedOrders   []Order
	PlannedDate        time.Time
	Routes             []Route
}

func (p Plan) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, p.ReferenceID)
}

func (p Plan) UpdateIfChanged(newPlan Plan) (Plan, bool) {
	changed := false

	// Comparar con zero value de time.Time
	if !newPlan.PlannedDate.IsZero() && !p.PlannedDate.Equal(newPlan.PlannedDate) {
		p.PlannedDate = newPlan.PlannedDate
		changed = true
	}

	return p, changed
}

// GetUnassignedRoute devuelve la ruta de órdenes no asignadas si existe
func (p Plan) GetUnassignedRoute() (Route, bool) {
	if len(p.UnassignedOrders) == 0 {
		return Route{}, false
	}
	return Route{
		ReferenceID: p.ReferenceID + "-unassigned",
		Vehicle:     Vehicle{},
		Orders:      p.UnassignedOrders,
	}, true
}

// AddUnassignedOrders agrega órdenes sin asignar al plan con la razón especificada
func (p *Plan) AddUnassignedOrders(orders []Order, reason string) {
	for _, order := range orders {
		order.UnassignedReason = reason
		p.UnassignedOrders = append(p.UnassignedOrders, order)
	}
}
