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
		ReferenceID string `json:"referenceID" example:"1400001234567890"`
	} `json:"orders"`
	CancellationReason struct {
		Detail      string `json:"detail" example:"no quiso recibir producto porque la caja estaba dañada"`
		Reason      string `json:"reason" example:"CLIENTE_RECHAZA_ENTREGA"`
		ReferenceID string `json:"referenceID" example:"1021"`
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
