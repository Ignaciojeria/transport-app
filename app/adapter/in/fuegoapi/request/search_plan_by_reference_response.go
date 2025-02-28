package request

import (
	"transport-app/app/domain"
)

type SearchPlanByReferenceResponse struct {
	ReferenceID   string `json:"referenceID"`
	PlannedDate   string `json:"plannedDate"`
	StartLocation struct {
		NodeReferenceID string  `json:"nodeReferenceID"`
		Latitude        float64 `json:"latitude"`
		Longitude       float64 `json:"longitude"`
	} `json:"startLocation"`
	UnassignedOrders []SearchOrdersResponse `json:"unassignedOrders"`
	Routes           []struct {
		ReferenceID string `json:"referenceID"`
		Operator    struct {
			ReferenceID string `json:"referenceID"`
		} `json:"operator"`
		EndLocation struct {
			NodeReferenceID string  `json:"nodeReferenceID"`
			Latitude        float64 `json:"latitude"`
			Longitude       float64 `json:"longitude"`
		} `json:"endLocation"`
		Visits []struct {
			Sequence  int                    `json:"sequence"`
			Latitude  float64                `json:"latitude"`
			Longitude float64                `json:"longitude"`
			Orders    []SearchOrdersResponse `json:"orders"`
		} `json:"visits"`
	} `json:"routes"`
}

func MapSearchPlanByReferenceResponse(p domain.Plan) SearchPlanByReferenceResponse {
	response := SearchPlanByReferenceResponse{
		ReferenceID: string(p.ReferenceID),
		PlannedDate: p.PlannedDate.Format("2006-01-02"),
		StartLocation: struct {
			NodeReferenceID string  `json:"nodeReferenceID"`
			Latitude        float64 `json:"latitude"`
			Longitude       float64 `json:"longitude"`
		}{
			NodeReferenceID: string(p.Origin.ReferenceID),
			Latitude:        p.Origin.AddressInfo.Location[1],
			Longitude:       p.Origin.AddressInfo.Location[0],
		},
	}
	response.UnassignedOrders = MapSearchOrdersResponse(p.UnassignedOrders)
	// Mapear rutas
	for _, route := range p.Routes {
		mappedRoute := struct {
			ReferenceID string `json:"referenceID"`
			Operator    struct {
				ReferenceID string `json:"referenceID"`
			} `json:"operator"`
			EndLocation struct {
				NodeReferenceID string  `json:"nodeReferenceID"`
				Latitude        float64 `json:"latitude"`
				Longitude       float64 `json:"longitude"`
			} `json:"endLocation"`
			Visits []struct {
				Sequence  int                    `json:"sequence"`
				Latitude  float64                `json:"latitude"`
				Longitude float64                `json:"longitude"`
				Orders    []SearchOrdersResponse `json:"orders"`
			} `json:"visits"`
		}{
			ReferenceID: string(route.ReferenceID),
			Operator: struct {
				ReferenceID string `json:"referenceID"`
			}{
				ReferenceID: route.Operator.ReferenceID,
			},
			EndLocation: struct {
				NodeReferenceID string  `json:"nodeReferenceID"`
				Latitude        float64 `json:"latitude"`
				Longitude       float64 `json:"longitude"`
			}{
				NodeReferenceID: string(route.Destination.ReferenceID),
				Latitude:        route.Destination.AddressInfo.Location[1],
				Longitude:       route.Destination.AddressInfo.Location[0],
			},
		}

		// Agrupar órdenes por sequenceNumber
		ordersBySequence := make(map[int][]domain.Order)
		unorderedOrders := []domain.Order{} // Lista para órdenes sin secuencia

		for _, order := range route.Orders {
			if order.SequenceNumber != nil {
				seq := *order.SequenceNumber
				ordersBySequence[seq] = append(ordersBySequence[seq], order)
			} else {
				unorderedOrders = append(unorderedOrders, order) // Guardar órdenes sin secuencia
			}
		}

		// Crear una visita para cada secuencia asignada
		for seq, ordersInSequence := range ordersBySequence {
			if len(ordersInSequence) > 0 {
				firstOrder := ordersInSequence[0] // Usamos la primera orden para la ubicación

				visit := struct {
					Sequence  int                    `json:"sequence"`
					Latitude  float64                `json:"latitude"`
					Longitude float64                `json:"longitude"`
					Orders    []SearchOrdersResponse `json:"orders"`
				}{
					Sequence:  seq,
					Latitude:  firstOrder.Destination.AddressInfo.CorrectedLocation[1],
					Longitude: firstOrder.Destination.AddressInfo.CorrectedLocation[0],
					Orders:    MapSearchOrdersResponse(ordersInSequence),
				}

				mappedRoute.Visits = append(mappedRoute.Visits, visit)
			}
		}

		// Agregar una visita especial para órdenes sin secuencia
		if len(unorderedOrders) > 0 {
			firstOrder := unorderedOrders[0] // Usamos la primera orden para la ubicación

			visit := struct {
				Sequence  int                    `json:"sequence"`
				Latitude  float64                `json:"latitude"`
				Longitude float64                `json:"longitude"`
				Orders    []SearchOrdersResponse `json:"orders"`
			}{
				Sequence:  -1, // Indica que no tienen secuencia asignada
				Latitude:  firstOrder.Destination.AddressInfo.CorrectedLocation[1],
				Longitude: firstOrder.Destination.AddressInfo.CorrectedLocation[0],
				Orders:    MapSearchOrdersResponse(unorderedOrders),
			}

			mappedRoute.Visits = append(mappedRoute.Visits, visit)
		}

		response.Routes = append(response.Routes, mappedRoute)
	}

	return response
}
