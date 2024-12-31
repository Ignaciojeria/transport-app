package mapper

import (
	"transport-app/app/adapter/in/fuegoapi/model"
	"transport-app/app/domain"
)

func MapCreateTransportOrderRequest(request model.CreateTransportOrderRequest) domain.TransportOrder {
	return domain.TransportOrder{
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
