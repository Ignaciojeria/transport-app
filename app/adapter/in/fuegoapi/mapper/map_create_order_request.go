package mapper

import (
	"transport-app/app/adapter/in/fuegoapi/model"
	"transport-app/app/domain"
)

func MapCreateOrderRequest(request model.CreateOrderRequest) domain.Order {
	return domain.Order{
		ReferenceID:             request.ReferenceID,
		OrderType:               request.OrderType,
		References:              request.References,
		Origin:                  request.Origin,
		Destination:             request.Destination,
		Items:                   request.Items,
		Packages:                request.Packages,
		CollectAvailabilityDate: request.CollectAvailabilityDate,
		PromisedDate:            request.PromisedDate,
		Visit:                   request.Visit,
		TransportRequirements:   request.TransportRequirements,
	}
}
