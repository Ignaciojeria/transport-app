package response

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
	UnassignedOrders []struct {
		Longitude   float64 `json:"longitude"`
		Latitude    float64 `json:"latitude"`
		ReferenceID string  `json:"referenceID"`
	} `json:"unassignedOrders"`
	Routes []struct {
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
			ReferenceID string                 `json:"referenceID"`
			Sequence    int                    `json:"sequence"`
			Latitude    float64                `json:"latitude"`
			Longitude   float64                `json:"longitude"`
			Orders      []SearchOrdersResponse `json:"orders"`
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
			Latitude:        p.Origin.AddressInfo.PlanCorrectedLocation[1],
			Longitude:       p.Origin.AddressInfo.PlanCorrectedLocation[0],
		},
	}

	// Map unassigned orders
	for _, order := range p.UnassignedOrders {
		unassignedOrder := struct {
			Longitude   float64 `json:"longitude"`
			Latitude    float64 `json:"latitude"`
			ReferenceID string  `json:"referenceID"`
		}{
			Longitude:   order.Destination.AddressInfo.PlanCorrectedLocation[0],
			Latitude:    order.Destination.AddressInfo.PlanCorrectedLocation[1],
			ReferenceID: string(order.ReferenceID),
		}
		response.UnassignedOrders = append(response.UnassignedOrders, unassignedOrder)
	}

	// Map routes
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
				ReferenceID string                 `json:"referenceID"`
				Sequence    int                    `json:"sequence"`
				Latitude    float64                `json:"latitude"`
				Longitude   float64                `json:"longitude"`
				Orders      []SearchOrdersResponse `json:"orders"`
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
				Latitude:        route.Destination.AddressInfo.PlanCorrectedLocation[1],
				Longitude:       route.Destination.AddressInfo.PlanCorrectedLocation[0],
			},
		}

		// Agrupar órdenes por sequenceNumber
		ordersBySequence := make(map[int][]domain.Order)
		for _, order := range route.Orders {
			if order.SequenceNumber != nil {
				seq := *order.SequenceNumber
				ordersBySequence[seq] = append(ordersBySequence[seq], order)
			}
		}

		// Crear una visita para cada secuencia
		for seq, ordersInSequence := range ordersBySequence {
			if len(ordersInSequence) > 0 {
				// Usar la ubicación de la primera orden del grupo
				firstOrder := ordersInSequence[0]

				visit := struct {
					ReferenceID string                 `json:"referenceID"`
					Sequence    int                    `json:"sequence"`
					Latitude    float64                `json:"latitude"`
					Longitude   float64                `json:"longitude"`
					Orders      []SearchOrdersResponse `json:"orders"`
				}{
					ReferenceID: string(firstOrder.ReferenceID),
					Sequence:    seq,
					Latitude:    firstOrder.Destination.AddressInfo.PlanCorrectedLocation[1],
					Longitude:   firstOrder.Destination.AddressInfo.PlanCorrectedLocation[0],
					Orders:      MapSearchOrdersResponse(ordersInSequence),
				}

				mappedRoute.Visits = append(mappedRoute.Visits, visit)
			}
		}

		response.Routes = append(response.Routes, mappedRoute)
	}

	return response
}
