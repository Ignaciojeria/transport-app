package domain

import (
	"context"
	"time"
)

type Plan struct {
	ReferenceID      string
	Origin           NodeInfo
	PlannedDate      time.Time
	UnassignedOrders []Order
	Routes           []Route
	PlanningStatus   PlanningStatus
	PlanType         PlanType
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
