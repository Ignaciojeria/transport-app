package domain

type Plan struct {
	Headers
	ID             int64
	ReferenceID    string
	Date           string
	Routes         []Route
	PlanningStatus PlanningStatus
	PlanType       PlanType
}

func (p Plan) UpdateIfChanged(newPlan Plan) Plan {
	// Actualizar campos básicos si no están vacíos
	if newPlan.ID != 0 {
		p.ID = newPlan.ID
	}
	if newPlan.ReferenceID != "" {
		p.ReferenceID = newPlan.ReferenceID
	}
	if newPlan.Date != "" {
		p.Date = newPlan.Date
	}

	// Actualizar PlanningStatus
	if newPlan.PlanningStatus.ID != 0 || newPlan.PlanningStatus.Value != "" {
		p.PlanningStatus = newPlan.PlanningStatus
	}

	// Actualizar PlanType
	if newPlan.PlanType.ID != 0 || newPlan.PlanType.Value != "" {
		p.PlanType = newPlan.PlanType
	}

	// Actualizar Routes si hay nuevas rutas
	if len(newPlan.Routes) > 0 {
		p.Routes = newPlan.Routes
	}

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
	return ps
}

type Route struct {
	Organization
	ID       int64
	Vehicle  Vehicle
	Operator Operator
	Orders   []Order
}

func (r Route) UpdateIfChanged(newRoute Route) Route {
	if newRoute.ID != 0 {
		r.ID = newRoute.ID
	}
	r.Operator = r.Operator.UpdateIfChanged(newRoute.Operator)
	r.Orders = newRoute.Orders
	return r
}
