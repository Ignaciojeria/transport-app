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
