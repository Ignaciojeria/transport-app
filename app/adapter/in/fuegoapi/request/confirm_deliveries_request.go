package request

import (
	"context"
	"time"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

type ConfirmDeliveriesRequest struct {
	ManualChange struct {
		PerformedBy string `json:"performedBy" example:"juan@example.com"`
		Reason      string `json:"reason" example:"Corrección tras reclamo de transporte"`
	} `json:"manualChange"`
	Carrier struct {
		Name       string `json:"name" example:"Transportes ABC"`
		NationalID string `json:"nationalID" example:"1234567890"`
	} `json:"carrier"`
	Driver struct {
		Email      string `json:"email" example:"juan@example.com"`
		NationalID string `json:"nationalID" example:"1234567890"`
	} `json:"driver"`
	Route struct {
		Orders []struct {
			BusinessIdentifiers struct {
				Commerce string `json:"commerce"`
				Consumer string `json:"consumer"`
			} `json:"businessIdentifiers"`
			Delivery struct {
				Failure struct {
					Detail      string `json:"detail" example:"no quiso recibir producto porque la caja estaba dañada"`
					Reason      string `json:"reason" example:"CLIENTE_RECHAZA_ENTREGA"`
					ReferenceID string `json:"referenceID" example:"1021"`
				} `json:"failure"`
				HandledAt string `json:"handledAt" example:"2025-06-06T14:30:00Z"`
				Location  struct {
					Latitude  float64 `json:"latitude" example:"19.432607"`
					Longitude float64 `json:"longitude" example:"-99.133209"`
				} `json:"location"`
			} `json:"delivery"`
			EvidencePhotos []struct {
				TakenAt string `json:"takenAt" example:"2025-06-06T14:30:00Z"`
				Type    string `json:"type" example:"HOUSE_PHOTO"`
				URL     string `json:"url" example:"https://ignaciojeria.github.io/"`
			} `json:"evidencePhotos"`
			DeliveryUnits []struct {
				Items []struct {
					Sku string `json:"sku" example:"SKU123"`
				} `json:"items"`
				Lpn string `json:"lpn" example:"ABC123"`
			} `json:"deliveryUnits"`
			Recipient struct {
				FullName   string `json:"fullName" example:"Juan Perez"`
				NationalID string `json:"nationalID" example:"1234567890"`
			} `json:"recipient"`
			ReferenceID string `json:"referenceID"`
		} `json:"orders"`
		ReferenceID string `json:"referenceID"`
	} `json:"route"`
	Vehicle struct {
		Plate string `json:"plate" example:"ABC123"`
	} `json:"vehicle"`
}

func (r ConfirmDeliveriesRequest) Map(ctx context.Context) domain.Route {

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

	// Mapear órdenes
	orders := make([]domain.Order, 0, len(r.Route.Orders))
	for _, order := range r.Route.Orders {
		domainOrder := domain.Order{
			Headers: domain.Headers{
				Consumer: order.BusinessIdentifiers.Consumer,
				Commerce: order.BusinessIdentifiers.Commerce,
				Channel:  sharedcontext.ChannelFromContext(ctx),
			},
			ReferenceID: domain.ReferenceID(order.ReferenceID),
			Destination: domain.NodeInfo{
				AddressInfo: domain.AddressInfo{
					Contact: domain.Contact{
						FullName:   order.Recipient.FullName,
						NationalID: order.Recipient.NationalID,
					},
				},
			},
		}

		// Mapear unidades de entrega
		deliveryUnits := make(domain.DeliveryUnits, 0, len(order.DeliveryUnits))
		for _, du := range order.DeliveryUnits {
			items := make([]domain.Item, 0, len(du.Items))
			for _, item := range du.Items {
				items = append(items, domain.Item{
					Sku: item.Sku,
				})
			}
			t, _ := time.Parse(time.RFC3339, order.Delivery.HandledAt)
			deliveryUnits = append(deliveryUnits, domain.DeliveryUnit{
				Lpn:   du.Lpn,
				Items: items,
				ConfirmDelivery: domain.ConfirmDelivery{
					ManualChange: domain.ManualChange{
						PerformedBy: r.ManualChange.PerformedBy,
						Reason:      r.ManualChange.Reason,
					},
					HandledAt: t,
					Latitude:  order.Delivery.Location.Latitude,
					Longitude: order.Delivery.Location.Longitude,
					EvidencePhotos: func() []domain.EvidencePhoto {
						photos := make([]domain.EvidencePhoto, 0, len(order.EvidencePhotos))
						for _, photo := range order.EvidencePhotos {
							takenAt, _ := time.Parse(time.RFC3339, photo.TakenAt)
							photos = append(photos, domain.EvidencePhoto{
								TakenAt: takenAt,
								Type:    photo.Type,
								URL:     photo.URL,
							})
						}
						return photos
					}(),
					NonDeliveryReason: domain.NonDeliveryReason{
						Reason:      order.Delivery.Failure.Reason,
						Details:     order.Delivery.Failure.Detail,
						ReferenceID: order.Delivery.Failure.ReferenceID,
					},
				},
			})
		}

		domainOrder.DeliveryUnits = deliveryUnits
		orders = append(orders, domainOrder)
	}

	route.Orders = orders

	return route
}
