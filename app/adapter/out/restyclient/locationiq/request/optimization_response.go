package request

import (
	"transport-app/app/domain"
)

type OptimizationResponse struct {
	Code  string `json:"code"`
	Trips []struct {
		Geometry string `json:"geometry"`
		Legs     []struct {
			Steps    []any   `json:"steps"`
			Summary  string  `json:"summary"`
			Weight   float64 `json:"weight"`
			Duration float64 `json:"duration"`
			Distance float64 `json:"distance"`
		} `json:"legs"`
		WeightName string  `json:"weight_name"`
		Weight     float64 `json:"weight"`
		Duration   float64 `json:"duration"`
		Distance   float64 `json:"distance"`
	} `json:"trips"`
	Waypoints []struct {
		WaypointIndex int       `json:"waypoint_index"`
		TripsIndex    int       `json:"trips_index"`
		Hint          string    `json:"hint"`
		Distance      float64   `json:"distance"`
		Name          string    `json:"name"`
		Location      []float64 `json:"location"`
	} `json:"waypoints"`
}

// Map transforma la respuesta de LocationIQ en una actualización del dominio.
// Se asigna la secuencia optimizada (waypoint_index), el punto corregido y
// la distancia corregida (CorrectedDistance) a cada orden.
// Se asume que el primer waypoint es el origen y cada orden corresponde al waypoint en la posición i+1.
func (res OptimizationResponse) Map(route domain.Route) domain.Route {
	if res.Code != "Ok" || len(res.Waypoints) <= 1 {
		return route
	}
	// Iteramos sobre las órdenes.
	// Cada orden se asocia con el waypoint en la posición (índice de la orden + 1)
	for idx := range route.Orders {
		wpIndex := idx + 1 // debido a que el primer waypoint es el origen
		if wpIndex < len(res.Waypoints) {
			wp := res.Waypoints[wpIndex]
			//correctedPoint := orb.Point{wp.Location[0], wp.Location[1]}
			seq := wp.WaypointIndex
			route.Orders[idx].SequenceNumber = &seq
			//route.Orders[idx].Destination.AddressInfo.CorrectedLocation = correctedPoint
			route.Orders[idx].Destination.AddressInfo.CorrectedDistance = wp.Distance
		}
	}

	//route.Orders = route.Orders
	return route
}
