package request

import (
	"context"
	"time"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

type CancelOrdersRequest struct {
	ManualChange struct {
		PerformedBy string `json:"performedBy" example:"juan@example.com"`
	} `json:"manualChange"`
	Orders []struct {
		BusinessIdentifiers struct {
			Commerce string `json:"commerce"`
			Consumer string `json:"consumer"`
		} `json:"businessIdentifiers"`
		ReferenceID string `json:"referenceID"`
	} `json:"orders"`
	CancellationReason struct {
		Detail      string `json:"detail"`
		Reason      string `json:"reason"`
		ReferenceID string `json:"referenceID"`
	} `json:"cancellationReason"`
}

func (r CancelOrdersRequest) Map(ctx context.Context) domain.Route {
	orders := make([]domain.Order, 0, len(r.Orders))
	for _, order := range r.Orders {
		domainOrder := domain.Order{
			Headers: domain.Headers{
				Consumer: order.BusinessIdentifiers.Consumer,
				Commerce: order.BusinessIdentifiers.Commerce,
				Channel:  sharedcontext.ChannelFromContext(ctx),
			},
			ReferenceID: domain.ReferenceID(order.ReferenceID),
			DeliveryUnits: domain.DeliveryUnits{
				{
					ConfirmDelivery: domain.ConfirmDelivery{
						ManualChange: domain.ManualChange{
							PerformedBy: r.ManualChange.PerformedBy,
						},
						HandledAt: time.Now(),
						NonDeliveryReason: domain.NonDeliveryReason{
							Reason:      r.CancellationReason.Reason,
							Details:     r.CancellationReason.Detail,
							ReferenceID: r.CancellationReason.ReferenceID,
						},
					},
				},
			},
		}
		orders = append(orders, domainOrder)
	}

	return domain.Route{
		Orders: orders,
	}
}
