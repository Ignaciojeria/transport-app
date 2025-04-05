package response

import (
	"transport-app/app/adapter/in/fuegoapi/mapper"
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
	UnassignedOrders []SearchPlanByReferenceOrdersResponse `json:"unassignedOrders"`
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
			Sequence  int                                   `json:"sequence"`
			Latitude  float64                               `json:"latitude"`
			Longitude float64                               `json:"longitude"`
			Orders    []SearchPlanByReferenceOrdersResponse `json:"orders"`
		} `json:"visits"`
	} `json:"routes"`
}

type SearchPlanByReferenceOrdersResponse struct {
	ReferenceID             string `json:"referenceID"`
	DeliveryInstructions    string `json:"deliveryInstructions"`
	CollectAvailabilityDate struct {
		Date      string `json:"date"`
		TimeRange struct {
			EndTime   string `json:"endTime"`
			StartTime string `json:"startTime"`
		} `json:"timeRange"`
	} `json:"collectAvailabilityDate"`
	Destination struct {
		AddressInfo struct {
			ProviderAddress string `json:"providerAddress"`
			AddressLine1    string `json:"addressLine1"`
			AddressLine2    string `json:"addressLine2"`
			AddressLine3    string `json:"addressLine3"`
			Contact         struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				Documents  []struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"documents"`
				FullName string `json:"fullName"`
			} `json:"contact"`
			Locality  string  `json:"locality"`
			District  string  `json:"district"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Province  string  `json:"province"`
			State     string  `json:"state"`
			TimeZone  string  `json:"timeZone"`
			ZipCode   string  `json:"zipCode"`
		} `json:"addressInfo"`
		DeliveryInstructions string `json:"deliveryInstructions"`
		NodeInfo             struct {
			ReferenceID string `json:"referenceID"`
			Name        string `json:"name"`
		} `json:"nodeInfo"`
	} `json:"destination"`
	Items []struct {
		Description string `json:"description"`
		Dimensions  struct {
			Length float64 `json:"length"`
			Height float64 `json:"height"`
			Unit   string  `json:"unit"`
			Width  float64 `json:"width"`
		} `json:"dimensions"`
		Insurance struct {
			Currency  string  `json:"currency"`
			UnitValue float64 `json:"unitValue"`
		} `json:"insurance"`
		LogisticCondition string `json:"logisticCondition"`
		Quantity          struct {
			QuantityNumber int    `json:"quantityNumber"`
			QuantityUnit   string `json:"quantityUnit"`
		} `json:"quantity"`
		Sku    string `json:"sku"`
		Weight struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
		} `json:"weight"`
	} `json:"items"`
	OrderType struct {
		Description string `json:"description"`
		Type        string `json:"type"`
	} `json:"orderType"`
	Origin struct {
		AddressInfo struct {
			ProviderAddress string `json:"providerAddress"`
			AddressLine1    string `json:"addressLine1"`
			AddressLine2    string `json:"addressLine2"`
			AddressLine3    string `json:"addressLine3"`
			Contact         struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				Documents  []struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"documents"`
				FullName string `json:"fullName"`
			} `json:"contact"`
			Locality  string  `json:"locality"`
			District  string  `json:"district"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Province  string  `json:"province"`
			State     string  `json:"state"`
			TimeZone  string  `json:"timeZone"`
			ZipCode   string  `json:"zipCode"`
		} `json:"addressInfo"`
		NodeInfo struct {
			ReferenceID string `json:"referenceID"`
			Name        string `json:"name"`
		} `json:"nodeInfo"`
	} `json:"origin"`
	Packages []struct {
		Dimensions struct {
			Length float64 `json:"length"`
			Height float64 `json:"height"`
			Unit   string  `json:"unit"`
			Width  float64 `json:"width"`
		} `json:"dimensions"`
		Insurance struct {
			Currency  string  `json:"currency"`
			UnitValue float64 `json:"unitValue"`
		} `json:"insurance"`
		ItemReferences []struct {
			Quantity struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			} `json:"quantity"`
			Sku string `json:"sku"`
		} `json:"itemReferences"`
		Lpn    string `json:"lpn"`
		Weight struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
		} `json:"weight"`
	} `json:"packages"`
	PromisedDate struct {
		DateRange struct {
			EndDate   string `json:"endDate"`
			StartDate string `json:"startDate"`
		} `json:"dateRange"`
		ServiceCategory string `json:"serviceCategory"`
		TimeRange       struct {
			EndTime   string `json:"endTime"`
			StartTime string `json:"startTime"`
		} `json:"timeRange"`
	} `json:"promisedDate"`
	References []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"references"`
	TransportRequirements []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"transportRequirements"`
	BusinessIdentifiers struct {
		Commerce string `json:"commerce"`
		Consumer string `json:"consumer"`
	} `json:"businessIdentifiers"`
	OrderStatus struct {
		ID        int64
		Status    string `json:"status"`
		CreatedAt string `json:"createdAt"`
	} `json:"orderStatus"`
}

// MapOrdersToSearchOrdersResponse convierte un slice de domain.Order a un slice de SearchOrdersResponse
func MapSearchPlanByReferenceOrdersResponse(orders []domain.Order) []SearchPlanByReferenceOrdersResponse {
	var responses []SearchPlanByReferenceOrdersResponse
	for _, order := range orders {
		response := SearchPlanByReferenceOrdersResponse{}
		response.ReferenceID = string(order.ReferenceID)
		response.OrderStatus = mapper.MapOrderStatusFromDomain(order.OrderStatus)
		response.DeliveryInstructions = order.DeliveryInstructions

		// Identificadores de negocio
		response.BusinessIdentifiers.Commerce = order.Headers.Commerce
		response.BusinessIdentifiers.Consumer = order.Headers.Consumer

		// OrderType
		response.OrderType = mapper.MapOrderTypeFromDomain(order.OrderType)

		// Origin y Destination
		originNodeInfo, originAddressInfo := mapper.MapNodeInfoToResponseNodeInfo(order.Origin)
		response.Origin.NodeInfo = originNodeInfo
		response.Origin.AddressInfo = originAddressInfo

		destNodeInfo, destAddressInfo := mapper.MapNodeInfoToResponseNodeInfo(order.Destination)
		response.Destination.NodeInfo = destNodeInfo
		response.Destination.AddressInfo = destAddressInfo
		response.Destination.DeliveryInstructions = order.DeliveryInstructions

		// Items, Packages
		response.Items = mapper.MapItemsFromDomain(order.Items)
		response.Packages = mapper.MapPackagesFromDomain(order.Packages)

		// References, Requirements
		response.References = mapper.MapReferencesFromDomain(order.References)
		response.TransportRequirements = mapper.MapReferencesFromDomain(order.TransportRequirements)

		// Fechas
		response.CollectAvailabilityDate = mapper.MapCollectAvailabilityDateFromDomain(order.CollectAvailabilityDate)
		response.PromisedDate = mapper.MapPromisedDateFromDomain(order.PromisedDate)
		responses = append(responses, response)
	}
	return responses
}

func MapSearchPlanByReferenceResponse(p domain.Plan) SearchPlanByReferenceResponse {
	result := SearchPlanByReferenceResponse{
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
	result.UnassignedOrders = MapSearchPlanByReferenceOrdersResponse(p.UnassignedOrders)
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
				Sequence  int                                   `json:"sequence"`
				Latitude  float64                               `json:"latitude"`
				Longitude float64                               `json:"longitude"`
				Orders    []SearchPlanByReferenceOrdersResponse `json:"orders"`
			} `json:"visits"`
		}{
			ReferenceID: string(route.ReferenceID),
			Operator: struct {
				ReferenceID string `json:"referenceID"`
			}{
				//	ReferenceID: route.Operator.ReferenceID,
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
					Sequence  int                                   `json:"sequence"`
					Latitude  float64                               `json:"latitude"`
					Longitude float64                               `json:"longitude"`
					Orders    []SearchPlanByReferenceOrdersResponse `json:"orders"`
				}{
					Sequence:  seq,
					Latitude:  firstOrder.Destination.AddressInfo.Location[1],
					Longitude: firstOrder.Destination.AddressInfo.Location[0],
					Orders:    MapSearchPlanByReferenceOrdersResponse(ordersInSequence),
				}

				mappedRoute.Visits = append(mappedRoute.Visits, visit)
			}
		}

		// Agregar una visita especial para órdenes sin secuencia
		if len(unorderedOrders) > 0 {
			firstOrder := unorderedOrders[0] // Usamos la primera orden para la ubicación

			visit := struct {
				Sequence  int                                   `json:"sequence"`
				Latitude  float64                               `json:"latitude"`
				Longitude float64                               `json:"longitude"`
				Orders    []SearchPlanByReferenceOrdersResponse `json:"orders"`
			}{
				Sequence:  -1, // Indica que no tienen secuencia asignada
				Latitude:  firstOrder.Destination.AddressInfo.Location[1],
				Longitude: firstOrder.Destination.AddressInfo.Location[0],
				Orders:    MapSearchPlanByReferenceOrdersResponse(unorderedOrders),
			}

			mappedRoute.Visits = append(mappedRoute.Visits, visit)
		}
		result.Routes = append(result.Routes, mappedRoute)
	}

	return result
}
