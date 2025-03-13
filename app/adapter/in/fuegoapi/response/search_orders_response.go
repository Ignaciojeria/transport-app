package response

import (
	"transport-app/app/adapter/in/fuegoapi/mapper"
	"transport-app/app/domain"
)

type SearchOrdersResponse struct {
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
			ReferenceID string `json:"referenceId"`
			Name        string `json:"name"`
		} `json:"nodeInfo"`
	} `json:"destination"`
	Items []struct {
		Description string `json:"description"`
		Dimensions  struct {
			Depth  float64 `json:"depth"`
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
		ReferenceID string `json:"referenceId"`
		Weight      struct {
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
			ReferenceID string `json:"referenceId"`
			Name        string `json:"name"`
		} `json:"nodeInfo"`
	} `json:"origin"`
	Packages []struct {
		Dimensions struct {
			Depth  float64 `json:"depth"`
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
			ReferenceID string `json:"referenceId"`
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

// Map convierte un SearchOrdersResponse a un objeto domain.Order
func (res SearchOrdersResponse) Map() domain.Order {
	order := domain.Order{
		ReferenceID:          domain.ReferenceID(res.ReferenceID),
		DeliveryInstructions: res.DeliveryInstructions,

		// Mapear headers desde BusinessIdentifiers
		Headers: domain.Headers{
			Commerce: res.BusinessIdentifiers.Commerce,
			Consumer: res.BusinessIdentifiers.Consumer,
		},

		// Mapear utilizando los mappers genéricos existentes
		Items:                   mapper.MapItemsToDomain(res.Items),
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
func MapOrdersToSearchOrdersResponse(orders []domain.Order) []SearchOrdersResponse {
	var responses []SearchOrdersResponse
	for _, order := range orders {
		response := SearchOrdersResponse{}
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
