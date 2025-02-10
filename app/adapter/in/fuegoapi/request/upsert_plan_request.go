package request

import (
	"time"
	"transport-app/app/domain"
)

type UpsertPlanRequest struct {
	ReferenceID   string `json:"referenceID"`
	PlannedDate   string `json:"plannedDate"`
	StartLocation struct {
		NodeReferenceID string  `json:"nodeReferenceID"`
		Latitude        float64 `json:"latitude"`
		Longitude       float64 `json:"longitude"`
	} `json:"startLocation"`
	WorkingHours struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"workingHours"`
	UnassignedOrders []struct {
		ReferenceID string `json:"referenceID"`
		Reason      string `json:"reason"`
		Location    struct {
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"longitude"`
		} `json:"location"`
	} `json:"unassignedOrders"`
	Routes []struct {
		ReferenceID string `json:"referenceID"`
		Operator    struct {
			ReferenceID string `json:"referenceID"`
		} `json:"operator"`
		Vehicle *struct {
			ReferenceID string `json:"referenceID"`
			Plate       string `json:"plate"`
		} `json:"vehicle,omitempty"`
		Driver *struct {
			ReferenceID string `json:"referenceID"`
			Email       string `json:"email"`
		} `json:"driver,omitempty"`
		Visits []struct {
			Sequence    int    `json:"sequence"`
			ReferenceID string `json:"referenceID"`
			Location    struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"location"`
			Orders []struct {
				ReferenceID string `json:"referenceID"`
			} `json:"orders"`
		} `json:"visits"`
	} `json:"routes"`
}

func (r UpsertPlanRequest) Map() domain.Plan {
	planDate, err := time.Parse("2006-01-02", r.PlannedDate)
	if err != nil {
		planDate = time.Time{}
	}

	// Mapear Origin como NodeInfo
	origin := domain.NodeInfo{
		ReferenceID: domain.ReferenceID(r.StartLocation.NodeReferenceID),
		AddressInfo: domain.AddressInfo{
			Latitude:  float32(r.StartLocation.Latitude),
			Longitude: float32(r.StartLocation.Longitude),
		},
	}

	// Por ahora Destination será igual a Origin ya que no viene en el request
	// Podrías modificar el request para incluir un destinationLocation si es necesario
	destination := origin

	// Mapear órdenes no asignadas
	var unassignedOrders []domain.Order
	for _, unassignedOrder := range r.UnassignedOrders {
		unassignedOrders = append(unassignedOrders, domain.Order{
			ReferenceID: domain.ReferenceID(unassignedOrder.ReferenceID),
			Destination: domain.NodeInfo{
				AddressInfo: domain.AddressInfo{
					Latitude:  float32(unassignedOrder.Location.Latitude),
					Longitude: float32(unassignedOrder.Location.Longitude),
				},
			},
			UnassignedReason: unassignedOrder.Reason,
		})
	}

	// Mapear rutas
	var routes []domain.Route
	for _, routeData := range r.Routes {
		var orders []domain.Order

		// Convertir visitas a órdenes
		for _, visitData := range routeData.Visits {
			// Crear NodeInfo para el destino de la orden
			destination := domain.NodeInfo{
				ReferenceID: domain.ReferenceID(visitData.ReferenceID),
				AddressInfo: domain.AddressInfo{
					Latitude:  float32(visitData.Location.Latitude),
					Longitude: float32(visitData.Location.Longitude),
				},
			}

			// Mapear las órdenes de la visita
			for _, orderData := range visitData.Orders {
				orders = append(orders, domain.Order{
					ReferenceID: domain.ReferenceID(orderData.ReferenceID),
					Origin:      origin,      // La orden hereda el origin del plan
					Destination: destination, // El destination viene de la visita
				})
			}
		}

		route := domain.Route{
			ReferenceID: routeData.ReferenceID,
			Operator: domain.Operator{
				ReferenceID: routeData.Operator.ReferenceID,
			},
			Orders: orders,
		}

		// Mapear vehículo si existe
		if routeData.Vehicle != nil {
			route.Vehicle = domain.Vehicle{
				ReferenceID: routeData.Vehicle.ReferenceID,
				Plate:       routeData.Vehicle.Plate,
			}
		}

		routes = append(routes, route)
	}

	return domain.Plan{
		ReferenceID:      r.ReferenceID,
		PlannedDate:      planDate,
		Origin:           origin,
		Destination:      destination,
		UnassignedOrders: unassignedOrders,
		Routes:           routes,
		PlanType: domain.PlanType{
			Value: "dailyPlan",
		},
		PlanningStatus: domain.PlanningStatus{
			Value: "planned",
		},
	}
}
