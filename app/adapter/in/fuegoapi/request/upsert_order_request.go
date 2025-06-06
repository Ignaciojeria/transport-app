package request

import (
	"context"
	"transport-app/app/adapter/in/fuegoapi/mapper"
	"transport-app/app/domain"
)

type UpsertOrderRequest struct {
	ReferenceID string            `json:"referenceID" validate:"required" example:"1400001234567890"`
	ExtraFields map[string]string `json:"extraFields"`
	GroupBy     struct {
		Type  string `json:"type" example:"customerOrder"`
		Value string `json:"value" example:"1234567890"`
	} `json:"groupBy"`
	CollectAvailabilityDate struct {
		Date      string `json:"date" example:"2025-03-30"`
		TimeRange struct {
			EndTime   string `json:"endTime" example:"09:00"`
			StartTime string `json:"startTime" example:"19:00"`
		} `json:"timeRange"`
	} `json:"collectAvailabilityDate"`
	Destination struct {
		AddressInfo struct {
			AddressLine1 string `json:"addressLine1" example:"Inglaterra 59"`
			AddressLine2 string `json:"addressLine2" example:"Piso 2214"`
			Contact      struct {
				AdditionalContactMethods []struct {
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
			District    string `json:"district" example:"la florida"`
			Coordinates struct {
				Latitude   float64 `json:"latitude" example:"-33.5147889"`
				Longitude  float64 `json:"longitude" example:"-70.6130425"`
				Source     string  `json:"source" example:"GOOGLE_MAPS"`
				Confidence struct {
					Level   float64 `json:"level" example:"0.1"`
					Message string  `json:"message" example:"DISTRICT_CENTROID"`
					Reason  string  `json:"reason" example:"PROVIDER_RESULT_OUT_OF_DISTRICT"`
				} `json:"confidence"`
			} `json:"coordinates"`
			Province string `json:"province" example:"santiago"`
			State    string `json:"state" example:"region metropolitana de santiago"`
			TimeZone string `json:"timeZone" example:"America/Santiago"`
			ZipCode  string `json:"zipCode" example:"7500000"`
		} `json:"addressInfo"`
		DeliveryInstructions string `json:"deliveryInstructions"`
		NodeInfo             struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
	} `json:"destination"`
	OrderType struct {
		Description string `json:"description"`
		Type        string `json:"type"`
	} `json:"orderType"`
	Origin struct {
		AddressInfo struct {
			AddressLine1 string `json:"addressLine1" example:"Inglaterra 59"`
			AddressLine2 string `json:"addressLine2" example:"Piso 2214"`
			Contact      struct {
				AdditionalContactMethods []struct {
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
			District    string `json:"district" example:"la florida"`
			Coordinates struct {
				Latitude   float64 `json:"latitude" example:"-33.5147889"`
				Longitude  float64 `json:"longitude" example:"-70.6130425"`
				Source     string  `json:"source" example:"GOOGLE_MAPS"`
				Confidence struct {
					Level   float64 `json:"level" example:"0.1"`
					Message string  `json:"message" example:"DISTRICT_CENTROID"`
					Reason  string  `json:"reason" example:"PROVIDER_RESULT_OUT_OF_DISTRICT"`
				} `json:"confidence"`
			} `json:"coordinates"`
			Province string `json:"province" example:"santiago"`
			State    string `json:"state" example:"region metropolitana de santiago"`
			TimeZone string `json:"timeZone" example:"America/Santiago"`
			ZipCode  string `json:"zipCode" example:"7500000"`
		} `json:"addressInfo"`
		NodeInfo struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
	} `json:"origin"`
	DeliveryUnits []struct {
		SizeCategory string `json:"sizeCategory" example:"XL"`
		Dimensions   struct {
			Length float64 `json:"length" example:"100.0"`
			Height float64 `json:"height" example:"100.0"`
			Unit   string  `json:"unit" example:"cm"`
			Width  float64 `json:"width" example:"100.0"`
		} `json:"dimensions"`
		Insurance struct {
			Currency  string  `json:"currency" example:"CLP"`
			UnitValue float64 `json:"unitValue" example:"10000"`
		} `json:"insurance"`
		Items []struct {
			Description string `json:"description" example:"Cama 1 plaza"`
			Dimensions  struct {
				Length float64 `json:"length" example:"100.0"`
				Height float64 `json:"height" example:"100.0"`
				Unit   string  `json:"unit" example:"cm"`
				Width  float64 `json:"width" example:"100.0"`
			} `json:"dimensions"`
			Insurance struct {
				Currency  string  `json:"currency" example:"CLP"`
				UnitValue float64 `json:"unitValue" example:"10000"`
			} `json:"insurance"`
			Skills []struct {
				Type        string `json:"type" example:"requires_refrigeration"`
				Value       string `json:"value" example:"2°C-8°C"`
				Description string `json:"description" example:"Producto refrigerado"`
			} `json:"skills"`
			Quantity struct {
				QuantityNumber int    `json:"quantityNumber" example:"1"`
				QuantityUnit   string `json:"quantityUnit" example:"unit"`
			} `json:"quantity"`
			Sku    string `json:"sku" example:"1234567890"`
			Weight struct {
				Unit  string  `json:"unit" example:"kg"`
				Value float64 `json:"value" example:"1.8"`
			} `json:"weight"`
		} `json:"items"`
		Lpn    string `json:"lpn" example:"1234567890"`
		Labels []struct {
			Type  string `json:"type" example:"skill"`
			Value string `json:"value" example:"requires_refrigeration"`
		} `json:"labels"`
		Weight struct {
			Unit  string  `json:"unit" example:"kg"`
			Value float64 `json:"value" example:"1.8"`
		} `json:"weight"`
	} `json:"deliveryUnits"`
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
}

// Map convierte el request a un objeto de dominio Order
func (req UpsertOrderRequest) Map(ctx context.Context) domain.Order {
	order := domain.Order{
		ReferenceID:             domain.ReferenceID(req.ReferenceID),
		OrderType:               mapper.MapOrderTypeToDomain(req.OrderType),
		References:              mapper.MapReferencesToDomain(req.References),
		Origin:                  mapper.MapNodeInfoToDomain(req.Origin.NodeInfo, req.Origin.AddressInfo),
		Destination:             mapper.MapNodeInfoToDomain(req.Destination.NodeInfo, req.Destination.AddressInfo),
		DeliveryUnits:           mapper.MapPackagesToDomain(req.DeliveryUnits),
		CollectAvailabilityDate: mapper.MapCollectAvailabilityDateToDomain(req.CollectAvailabilityDate),
		PromisedDate:            mapper.MapPromisedDateToDomain(req.PromisedDate),
		DeliveryInstructions:    req.Destination.DeliveryInstructions,
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
