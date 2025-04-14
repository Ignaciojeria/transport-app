package request

import (
	"time"
	"transport-app/app/adapter/in/fuegoapi/mapper"
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
	UnassignedOrders []UpsertPlanOrderRequest `json:"unassignedOrders"`
	Routes           []struct {
		ReferenceID string `json:"referenceID"`
		EndLocation struct {
			NodeReferenceID string  `json:"nodeReferenceID"`
			Latitude        float64 `json:"latitude"`
			Longitude       float64 `json:"longitude"`
		} `json:"endLocation"`
		Container struct {
			Lpn string `json:"lpn"`
		} `json:"container"`
		Operator struct {
			Email string `json:"email"`
		} `json:"operator"`
		Vehicle *struct {
			Plate string `json:"plate"`
		} `json:"vehicle,omitempty"`
		Driver *struct {
			Email string `json:"email"`
		} `json:"driver,omitempty"`
		Visits []struct {
			Sequence  *int                     `json:"sequence"`
			Latitude  float64                  `json:"latitude"`
			Longitude float64                  `json:"longitude"`
			Orders    []UpsertPlanOrderRequest `json:"orders"`
		} `json:"visits"`
	} `json:"routes"`
}

type UpsertPlanOrderRequest struct {
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
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`

			Contact struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				Documents  []struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"documents"`
				FullName string `json:"fullName"`
			} `json:"contact"`

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
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`

			Contact struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				Documents  []struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"documents"`
				FullName string `json:"fullName"`
			} `json:"contact"`

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

// Map convierte un SearchOrdersResponse a un objeto domain.Order
func (res UpsertPlanOrderRequest) Map() domain.Order {

	order := domain.Order{
		ReferenceID:          domain.ReferenceID(res.ReferenceID),
		DeliveryInstructions: res.DeliveryInstructions,

		// Mapear headers desde BusinessIdentifiers
		Headers: domain.Headers{
			Commerce: res.BusinessIdentifiers.Commerce,
			Consumer: res.BusinessIdentifiers.Consumer,
		},
		// Mapear utilizando los mappers genéricos existentes
		//Items:                   mapper.MapItemsToDomain(res.Items),
		OrderType:               mapper.MapOrderTypeToDomain(res.OrderType),
		References:              mapper.MapReferencesToDomain(res.References),
		TransportRequirements:   mapper.MapReferencesToDomain(res.TransportRequirements),
		Packages:                mapper.MapPackagesToDomain(res.Packages),
		Origin:                  mapper.MapNodeInfoToDomain(res.Origin.NodeInfo, res.Origin.AddressInfo),
		Destination:             mapper.MapNodeInfoToDomain(res.Destination.NodeInfo, res.Destination.AddressInfo),
		CollectAvailabilityDate: mapper.MapCollectAvailabilityDateToDomain(res.CollectAvailabilityDate),
		PromisedDate:            mapper.MapPromisedDateToDomain(res.PromisedDate),
		OrderStatus:             mapper.MapOrderStatusToDomain(res.OrderStatus),
	}

	return order
}

// MapOrdersToSearchOrdersResponse convierte un slice de domain.Order a un slice de SearchOrdersResponse
func MapOrdersToPlanOrder(orders []domain.Order) []UpsertPlanOrderRequest {
	var responses []UpsertPlanOrderRequest
	for _, order := range orders {
		response := UpsertPlanOrderRequest{}
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
		//	response.Items = mapper.MapItemsFromDomain(order.Items)
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
				order.SequenceNumber = visitData.Sequence
				orders = append(orders, order)
			}
		}

		route := domain.Route{
			//	ReferenceID: routeData.ReferenceID,
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
				//	ReferenceID: routeData.Operator.ReferenceID,
			},
			Orders: orders,
		}

		// Mapear vehículo si existe
		if routeData.Vehicle != nil {
			route.Vehicle = domain.Vehicle{
				Plate: routeData.Vehicle.Plate,
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
