package request

import (
	"time"
	"transport-app/app/domain"
)

type ConfirmDeliveriesRequest struct {
	Plan struct {
		Routes []struct {
			ReferenceID string `json:"referenceID"`
			Vehicle     struct {
				ReferenceID string `json:"referenceID"`
				Plate       string `json:"plate"`
			} `json:"vehicle"`
			Carrier struct {
				ReferenceID string `json:"referenceID"`
				NationalID  string `json:"nationalID"`
			} `json:"carrier"`
			Orders []struct {
				ReferenceID         string    `json:"referenceID"`
				Status              string    `json:"status"`
				DeliveredAt         time.Time `json:"deliveredAt"`
				BusinessIdentifiers struct {
					Commerce string `json:"commerce"`
					Consumer string `json:"consumer"`
				} `json:"businessIdentifiers"`
				Recipient struct {
					FullName   string `json:"fullName"`
					NationalID string `json:"nationalID"`
				} `json:"recipient"`
				EvidencePhotos []struct {
					URL     string    `json:"url"`
					Type    string    `json:"type"`
					TakenAt time.Time `json:"takenAt"`
				} `json:"evidencePhotos"`
				Delivery struct {
					Failure struct {
						ReferenceID string `json:"referenceID"`
						Reason      string `json:"reason"`
						Detail      string `json:"detail"`
					} `json:"failure"`
					Location struct {
						Latitude  float64 `json:"latitude"`
						Longitude float64 `json:"longitude"`
					} `json:"location"`
				} `json:"delivery"`
			} `json:"orders"`
		} `json:"routes"`
	} `json:"plan"`
}

func MapCheckout(request ConfirmDeliveriesRequest) []domain.Checkout {
	var checkouts []domain.Checkout

	for _, route := range request.Plan.Routes {
		for _, order := range route.Orders {
			evidencePhotos := make([]domain.EvidencePhotos, len(order.EvidencePhotos))
			for i, photo := range order.EvidencePhotos {
				evidencePhotos[i] = domain.EvidencePhotos{
					URL:     photo.URL,
					Type:    photo.Type,
					TakenAt: photo.TakenAt,
				}
			}

			checkouts = append(checkouts, domain.Checkout{
				Order: domain.Order{
					Headers: domain.Headers{
						Consumer: order.BusinessIdentifiers.Consumer,
						Commerce: order.BusinessIdentifiers.Commerce,
					},
					ReferenceID: domain.ReferenceID(order.ReferenceID),
				},
				Route: domain.Route{
					ReferenceID: route.ReferenceID,
				},
				OrderStatus: domain.OrderStatus{
					Status:    order.Status,
					CreatedAt: time.Now(),
				},
				DeliveredAt: order.DeliveredAt,
				Vehicle: domain.Vehicle{
					ReferenceID: route.Vehicle.ReferenceID,
					Plate:       route.Vehicle.Plate,
					Carrier: domain.Carrier{
						ReferenceID: route.Carrier.ReferenceID,
						NationalID:  route.Carrier.NationalID,
					},
				},
				Recipient: domain.Recipient{
					FullName:   order.Recipient.FullName,
					NationalID: order.Recipient.NationalID,
				},
				EvidencePhotos: evidencePhotos,
				Latitude:       order.Delivery.Location.Latitude,
				Longitude:      order.Delivery.Location.Longitude,
				NotDeliveryReason: domain.NotDeliveryReason{
					ReferenceID: order.Delivery.Failure.ReferenceID,
					Detail:      order.Delivery.Failure.Detail,
				},
			})
		}
	}

	return checkouts
}
