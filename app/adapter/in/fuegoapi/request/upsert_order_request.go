package request

import (
	"context"
	"transport-app/app/adapter/in/fuegoapi/mapper"
	"transport-app/app/domain"
)

type UpsertOrderRequest struct {
	ReferenceID             string            `json:"referenceID" validate:"required"`
	ExtraFields             map[string]string `json:"extraFields"`
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
			Contact      struct {
				AdditionalContactMethods struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"additionalContactMethods"`
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
	OrderType struct {
		Description string `json:"description"`
		Type        string `json:"type"`
	} `json:"orderType"`
	Origin struct {
		AddressInfo struct {
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`
			Contact      struct {
				AdditionalContactMethods struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"additionalContactMethods"`
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
			Skills []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			}
			Quantity struct {
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
}

// Map convierte el request a un objeto de dominio Order
func (req UpsertOrderRequest) Map(ctx context.Context) domain.Order {
	order := domain.Order{
		ReferenceID:             domain.ReferenceID(req.ReferenceID),
		OrderType:               mapper.MapOrderTypeToDomain(req.OrderType),
		References:              mapper.MapReferencesToDomain(req.References),
		Origin:                  mapper.MapNodeInfoToDomain(req.Origin.NodeInfo, req.Origin.AddressInfo),
		Destination:             mapper.MapNodeInfoToDomain(req.Destination.NodeInfo, req.Destination.AddressInfo),
		Packages:                mapper.MapPackagesToDomain(req.Packages),
		CollectAvailabilityDate: mapper.MapCollectAvailabilityDateToDomain(req.CollectAvailabilityDate),
		PromisedDate:            mapper.MapPromisedDateToDomain(req.PromisedDate),
		DeliveryInstructions:    req.Destination.DeliveryInstructions,
		TransportRequirements:   mapper.MapReferencesToDomain(req.TransportRequirements),
		ExtraFields:             req.ExtraFields,
	}
	order.Headers.SetFromContext(ctx)
	if order.Commerce == "" {
		order.Commerce = "empty"
	}
	if order.Consumer == "" {
		order.Consumer = "empty"
	}
	return order
}
