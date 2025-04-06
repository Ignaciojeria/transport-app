package domain

import (
	"time"
)

type Plan struct {
	Headers
	ReferenceID      string
	Origin           NodeInfo
	PlannedDate      time.Time
	UnassignedOrders []Order
	Routes           []Route
	PlanningStatus   PlanningStatus
	PlanType         PlanType
}

func (p Plan) DocID() DocumentID {
	return Hash(p.Organization, p.ReferenceID)
}

func (p Plan) UpdateIfChanged(newPlan Plan) Plan {
	// Comparar con zero value de time.Time
	if !newPlan.PlannedDate.IsZero() {
		p.PlannedDate = newPlan.PlannedDate
	}
	return p
}
