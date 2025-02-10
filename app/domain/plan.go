package domain

import "time"

type Plan struct {
	Headers
	ID               int64
	Origin           NodeInfo
	Destination      NodeInfo
	ReferenceID      string
	PlannedDate      time.Time
	UnassignedOrders []Order
	Routes           []Route
	PlanningStatus   PlanningStatus
	PlanType         PlanType
}

func (p Plan) GetOrderSearchFilters() OrderSearchFilters {
	osf := OrderSearchFilters{}
	osf.Organization = p.Organization
	for _, route := range p.Routes {
		for _, order := range route.Orders {
			osf.ReferenceIDs = append(osf.ReferenceIDs, string(order.ReferenceID))
			osf.RouteReferenceID = route.ReferenceID
		}
	}
	return osf
}

func (p Plan) UpdateIfChanged(newPlan Plan) Plan {
	// Actualizar campos básicos si no están vacíos
	if newPlan.ID != 0 {
		p.ID = newPlan.ID
	}
	if newPlan.ReferenceID != "" {
		p.ReferenceID = newPlan.ReferenceID
	}
	// Comparar con zero value de time.Time
	if !newPlan.PlannedDate.IsZero() {
		p.PlannedDate = newPlan.PlannedDate
	}

	// El resto permanece igual
	if newPlan.PlanningStatus.ID != 0 || newPlan.PlanningStatus.Value != "" {
		p.PlanningStatus = newPlan.PlanningStatus
	}

	if newPlan.PlanType.ID != 0 || newPlan.PlanType.Value != "" {
		p.PlanType = newPlan.PlanType
	}

	if len(newPlan.Routes) > 0 {
		p.Routes = newPlan.Routes
	}
	p.Organization = newPlan.Organization
	return p
}

type PlanType struct {
	Organization
	ID    int64
	Value string
}

func (pt PlanType) UpdateIfChanged(newPlanType PlanType) PlanType {
	if newPlanType.ID != 0 {
		pt.ID = newPlanType.ID
	}
	if newPlanType.Value != "" {
		pt.Value = newPlanType.Value
	}
	pt.Organization = newPlanType.Organization
	return pt
}

type PlanningStatus struct {
	Organization
	ID    int64
	Value string
}

func (ps PlanningStatus) UpdateIfChanged(newPlanningStatus PlanningStatus) PlanningStatus {
	if newPlanningStatus.ID != 0 {
		ps.ID = newPlanningStatus.ID
	}
	if newPlanningStatus.Value != "" {
		ps.Value = newPlanningStatus.Value
	}
	ps.Organization = newPlanningStatus.Organization
	return ps
}

type Route struct {
	Organization
	ReferenceID string
	PlanID      int64
	ID          int64
	Vehicle     Vehicle
	Operator    Operator
	Orders      []Order
}

func (r Route) UpdateIfChanged(newRoute Route) Route {
	if newRoute.ID != 0 {
		r.ID = newRoute.ID
	}
	if newRoute.PlanID != 0 {
		r.PlanID = newRoute.PlanID
	}
	if newRoute.ReferenceID != "" {
		r.ReferenceID = newRoute.ReferenceID
	}
	r.Operator = r.Operator.UpdateIfChanged(newRoute.Operator)
	r.Orders = newRoute.Orders
	r.Organization = newRoute.Organization
	return r
}
