package request

import (
	"time"
	"transport-app/app/domain"
)

type OrdersCheckoutRequest struct {
	OrderReferenceID string    `json:"orderReferenceID"`
	RouteReferenceID string    `json:"routeReferenceID"`
	Status           string    `json:"status"`
	DeliveredAt      time.Time `json:"deliveredAt"`
	Vehicle          struct {
		ReferenceID string `json:"referenceID"`
		Plate       string `json:"plate"`
		Carrier     struct {
			ReferenceID string `json:"referenceID"`
			NationalID  string `json:"nationalID"`
		} `json:"carrier"`
	} `json:"vehicle"`
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
			ReasonReferenceID string `json:"reasonReferenceID"`
			Reason            string `json:"reason"`
			Detail            string `json:"detail"`
		} `json:"failure"`
		Location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
	} `json:"delivery"`
}

func Map(r OrdersCheckoutRequest) domain.OrderCheckout {
	evidencePhotos := make([]domain.EvicencePhotos, len(r.EvidencePhotos))
	for i, photo := range r.EvidencePhotos {
		evidencePhotos[i] = domain.EvicencePhotos{
			URL:     photo.URL,
			Type:    photo.Type,
			TakenAt: photo.TakenAt,
		}
	}

	return domain.OrderCheckout{
		Order: domain.Order{
			ReferenceID: domain.ReferenceID(r.OrderReferenceID),
		},
		Route: domain.Route{
			ReferenceID: r.RouteReferenceID,
		},
		OrderStatus: domain.OrderStatus{
			Status:    r.Status,
			CreatedAt: time.Now(),
		},
		DeliveredAt: r.DeliveredAt,
		Vehicle: domain.Vehicle{
			ReferenceID: r.Vehicle.ReferenceID,
			Plate:       r.Vehicle.Plate,
			Carrier: domain.Carrier{
				ReferenceID: r.Vehicle.Carrier.ReferenceID,
				NationalID:  r.Vehicle.Carrier.NationalID,
			},
		},
		Recipient: domain.Recipient{
			FullName:   r.Recipient.FullName,
			NationalID: r.Recipient.NationalID,
		},
		EvicencePhotos: evidencePhotos,
		Latitude:       float32(r.Delivery.Location.Latitude),
		Longitude:      float32(r.Delivery.Location.Longitude),
		NotDeliveryReason: domain.NotDeliveryReason{
			ReferenceID: r.Delivery.Failure.ReasonReferenceID,
			Detail:      r.Delivery.Failure.Detail,
		},
	}
}
