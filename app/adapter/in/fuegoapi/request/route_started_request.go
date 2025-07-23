package request

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

type RouteStartedRequest struct {
	StartedAt string `json:"startedAt" example:"2025-06-06T14:30:00Z"`
	Carrier   struct {
		Name       string `json:"name" example:"Transportes ABC"`
		NationalID string `json:"nationalID" example:"1234567890"`
	} `json:"carrier"`
	Driver struct {
		Email      string `json:"email" example:"juan@example.com"`
		NationalID string `json:"nationalID" example:"1234567890"`
	} `json:"driver"`
	Vehicle struct {
		Plate string `json:"plate" example:"ABC123"`
	} `json:"vehicle"`
	Route struct {
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
			ReferenceID string `json:"referenceID"`
		} `json:"orders"`
		ReferenceID string `json:"referenceID"`
	} `json:"route"`
}

func (r RouteStartedRequest) Map(ctx context.Context) domain.Route {
	route := domain.Route{
		ReferenceID: r.Route.ReferenceID,
		Vehicle: domain.Vehicle{
			Plate: r.Vehicle.Plate,
			Carrier: domain.Carrier{
				Name:       r.Carrier.Name,
				NationalID: r.Carrier.NationalID,
				Driver: domain.Driver{
					Email:      r.Driver.Email,
					NationalID: r.Driver.NationalID,
				},
			},
		},
	}

	// Map orders
	orders := make([]domain.Order, 0, len(r.Route.Orders))

	for _, order := range r.Route.Orders {
		domainOrder := domain.Order{
			Headers: domain.Headers{
				Consumer: order.BusinessIdentifiers.Consumer,
				Commerce: order.BusinessIdentifiers.Commerce,
				Channel:  sharedcontext.ChannelFromContext(ctx),
			},
			ReferenceID: domain.ReferenceID(order.ReferenceID),
		}

		if len(order.DeliveryUnits) == 0 {
			order.DeliveryUnits = append(order.DeliveryUnits, struct {
				Items []struct {
					Sku string "json:\"sku\" example:\"SKU123\""
				} "json:\"items\""
				Lpn string "json:\"lpn\" example:\"ABC123\""
			}{})
		}

		// Map delivery units
		deliveryUnits := make(domain.DeliveryUnits, 0, len(order.DeliveryUnits))
		for _, du := range order.DeliveryUnits {
			items := make([]domain.Item, 0, len(du.Items))
			for _, item := range du.Items {
				items = append(items, domain.Item{
					Sku: item.Sku,
				})
			}
			deliveryUnits = append(deliveryUnits, domain.DeliveryUnit{
				Lpn:   du.Lpn,
				Items: items,
				// Volume, Weight, Insurance: nil (no presentes en este request)
			})
		}

		domainOrder.DeliveryUnits = deliveryUnits
		orders = append(orders, domainOrder)
	}

	route.Orders = orders

	return route
}
