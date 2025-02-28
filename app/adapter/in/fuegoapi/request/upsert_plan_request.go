package request

import (
	"time"
	"transport-app/app/domain"

	"github.com/paulmach/orb"
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
	UnassignedOrders []SearchOrdersResponse `json:"unassignedOrders"`
	Routes           []struct {
		ReferenceID string `json:"referenceID"`
		EndLocation struct {
			NodeReferenceID string  `json:"nodeReferenceID"`
			Latitude        float64 `json:"latitude"`
			Longitude       float64 `json:"longitude"`
		} `json:"endLocation"`
		Operator struct {
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
			Sequence  int                    `json:"sequence"`
			Latitude  float64                `json:"latitude"`
			Longitude float64                `json:"longitude"`
			Orders    []SearchOrdersResponse `json:"orders"`
		} `json:"visits"`
	} `json:"routes"`
}

func (r UpsertPlanRequest) Map() domain.Plan {
	planDate, err := time.Parse("2006-01-02", r.PlannedDate)
	if err != nil {
		planDate = time.Time{}
	}

	// Mapear startLocation como NodeInfo
	startLocation := domain.NodeInfo{
		ReferenceID: domain.ReferenceID(r.StartLocation.NodeReferenceID),
		AddressInfo: domain.AddressInfo{
			Location: orb.Point{
				r.StartLocation.Longitude,
				r.StartLocation.Latitude,
			},
		},
	}

	// Mapear órdenes no asignadas
	var unassignedOrders []domain.Order
	for _, unassignedOrder := range r.UnassignedOrders {
		order := unassignedOrder.Map()
		unassignedOrders = append(unassignedOrders, order)
	}

	// Mapear rutas
	var routes []domain.Route
	for _, routeData := range r.Routes {
		var orders []domain.Order

		// Convertir visitas a órdenes
		for _, visitData := range routeData.Visits {
			// Mapear las órdenes de la visita
			for _, orderData := range visitData.Orders {
				order := orderData.Map()
				orders = append(orders, order)
			}
		}

		route := domain.Route{
			ReferenceID: routeData.ReferenceID,
			Destination: domain.NodeInfo{
				ReferenceID: domain.ReferenceID(routeData.EndLocation.NodeReferenceID),
				AddressInfo: domain.AddressInfo{
					Location: orb.Point{
						routeData.EndLocation.Longitude,
						routeData.EndLocation.Latitude,
					},
				},
			},
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
		Origin:           startLocation,
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
