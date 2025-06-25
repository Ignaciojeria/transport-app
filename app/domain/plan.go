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

// AssignIndexesToAllOrders itera sobre todas las órdenes del plan y ejecuta AssignIndexesIfNoLPN() en cada una
func (p *Plan) AssignIndexesToAllOrders() {
	// Iterar sobre órdenes no asignadas
	for i := range p.UnassignedOrders {
		p.UnassignedOrders[i].AssignIndexesIfNoLPN()
	}

	// Iterar sobre órdenes en las rutas
	for i := range p.Routes {
		for j := range p.Routes[i].Orders {
			p.Routes[i].Orders[j].AssignIndexesIfNoLPN()
		}
	}
}

// AssignSequenceNumbersToAllOrders asigna números de secuencia a todas las órdenes del plan
// basándose en su posición en los arreglos
func (p *Plan) AssignSequenceNumbersToAllOrders() {
	sequenceCounter := 1

	// Asignar secuencia a órdenes no asignadas
	for i := range p.UnassignedOrders {
		p.UnassignedOrders[i].SetSequenceNumber(sequenceCounter)
		sequenceCounter++
	}

	// Asignar secuencia a órdenes en las rutas (cada ruta comienza desde 1)
	for i := range p.Routes {
		routeSequenceCounter := 1
		for j := range p.Routes[i].Orders {
			p.Routes[i].Orders[j].SetSequenceNumber(routeSequenceCounter)
			routeSequenceCounter++
		}
	}
}

// AssignSequenceNumbersToAllRoutes asigna números de secuencia a todas las rutas del plan
// Cada ruta comienza desde 1
func (p *Plan) AssignSequenceNumbersToAllRoutes() {
	for i := range p.Routes {
		routeSequenceCounter := 1
		for j := range p.Routes[i].Orders {
			p.Routes[i].Orders[j].SetSequenceNumber(routeSequenceCounter)
			routeSequenceCounter++
		}
	}
}

// AssignSequenceNumbersToRouteOrders asigna números de secuencia solo a las órdenes de una ruta específica
func (p *Plan) AssignSequenceNumbersToRouteOrders(routeIndex int) {
	if routeIndex < 0 || routeIndex >= len(p.Routes) {
		return
	}

	sequenceCounter := 1
	for j := range p.Routes[routeIndex].Orders {
		p.Routes[routeIndex].Orders[j].SetSequenceNumber(sequenceCounter)
		sequenceCounter++
	}
}

// GetTotalOrderCount retorna el número total de órdenes en el plan
func (p *Plan) GetTotalOrderCount() int {
	total := len(p.UnassignedOrders)
	for _, route := range p.Routes {
		total += len(route.Orders)
	}
	return total
}
