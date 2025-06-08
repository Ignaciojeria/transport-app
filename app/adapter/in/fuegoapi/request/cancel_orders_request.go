package request

import (
	"context"
	"errors"
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
		DeliveryUnits []struct {
			Items []struct {
				Sku string `json:"sku" example:"SKU123"`
			} `json:"items"`
			Lpn string `json:"lpn" example:"ABC123"`
		} `json:"deliveryUnits"`
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
		deliveryUnits := make(domain.DeliveryUnits, 0, len(order.DeliveryUnits))
		// Mapear unidades de entrega
		if len(order.DeliveryUnits) == 0 {
			order.DeliveryUnits = append(order.DeliveryUnits, struct {
				Items []struct {
					Sku string `json:"sku" example:"SKU123"`
				} `json:"items"`
				Lpn string `json:"lpn" example:"ABC123"`
			}{})
		}
		for _, du := range order.DeliveryUnits {
			items := make([]domain.Item, 0, len(du.Items))
			for _, item := range du.Items {
				items = append(items, domain.Item{
					Sku: item.Sku,
				})
			}
			deliveryUnits = append(deliveryUnits, domain.DeliveryUnit{
				Items: items,
				Lpn:   du.Lpn,
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
			})
		}

		domainOrder := domain.Order{
			Headers: domain.Headers{
				Consumer: order.BusinessIdentifiers.Consumer,
				Commerce: order.BusinessIdentifiers.Commerce,
				Channel:  sharedcontext.ChannelFromContext(ctx),
			},
			ReferenceID:   domain.ReferenceID(order.ReferenceID),
			DeliveryUnits: deliveryUnits,
		}
		orders = append(orders, domainOrder)
	}

	return domain.Route{
		Orders: orders,
	}
}

func (r CancelOrdersRequest) Validate() error {
	for _, order := range r.Orders {
		for _, du := range order.DeliveryUnits {
			// Si no hay LPN ni SKUs, retornar error
			if du.Lpn == "" && len(du.Items) == 0 {
				return errors.New("delivery unit must have either LPN or SKUs")
			}
		}
	}
	return nil
}
