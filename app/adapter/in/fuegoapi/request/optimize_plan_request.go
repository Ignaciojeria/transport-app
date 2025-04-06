package request

import (
	"time"
	"transport-app/app/adapter/in/fuegoapi/mapper"
	"transport-app/app/domain"

	"github.com/paulmach/orb"
)

type OptimizePlanRequest struct {
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
	} `json:"unassignedOrders"`
	Routes []struct {
		ReferenceID string `json:"referenceID"`
		EndLocation struct {
			NodeReferenceID string  `json:"nodeReferenceID"`
			Latitude        float64 `json:"latitude"`
			Longitude       float64 `json:"longitude"`
		} `json:"endLocation"`
		Operator struct {
			Email string `json:"email"`
		} `json:"operator"`
		Vehicle *struct {
			Plate string `json:"plate"`
		} `json:"vehicle,omitempty"`
		Driver *struct {
			Email string `json:"email"`
		} `json:"driver,omitempty"`
	} `json:"routes"`
}

func (r OptimizePlanRequest) Map() domain.Plan {
	planDate, err := time.Parse("2006-01-02", r.PlannedDate)
	if err != nil {
		planDate = time.Time{}
	}

	startLocation := domain.NodeInfo{
		ReferenceID: domain.ReferenceID(r.StartLocation.NodeReferenceID),
		AddressInfo: domain.AddressInfo{
			Location: orb.Point{
				r.StartLocation.Longitude,
				r.StartLocation.Latitude,
			},
		},
	}

	var unassignedOrders []domain.Order

	for _, uo := range r.UnassignedOrders {
		order := domain.Order{
			ReferenceID:          domain.ReferenceID(uo.ReferenceID),
			DeliveryInstructions: uo.DeliveryInstructions,
			Headers: domain.Headers{
				Commerce: uo.BusinessIdentifiers.Commerce,
				Consumer: uo.BusinessIdentifiers.Consumer,
			},
			Items:                   mapper.MapItemsToDomain(uo.Items),
			OrderType:               mapper.MapOrderTypeToDomain(uo.OrderType),
			References:              mapper.MapReferencesToDomain(uo.References),
			TransportRequirements:   mapper.MapReferencesToDomain(uo.TransportRequirements),
			Packages:                mapper.MapPackagesToDomain(uo.Packages),
			Origin:                  mapper.MapNodeInfoToDomain(uo.Origin.NodeInfo, uo.Origin.AddressInfo),
			Destination:             mapper.MapNodeInfoToDomain(uo.Destination.NodeInfo, uo.Destination.AddressInfo),
			CollectAvailabilityDate: mapper.MapCollectAvailabilityDateToDomain(uo.CollectAvailabilityDate),
			PromisedDate:            mapper.MapPromisedDateToDomain(uo.PromisedDate),
			OrderStatus:             mapper.MapOrderStatusToDomain(uo.OrderStatus),
		}
		unassignedOrders = append(unassignedOrders, order)
	}

	var routes []domain.Route
	for _, routeData := range r.Routes {
		route := domain.Route{
			//ReferenceID: routeData.ReferenceID,
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
		}

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
