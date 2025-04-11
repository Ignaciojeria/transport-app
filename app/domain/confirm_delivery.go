package domain

import "time"

type ConfirmDelivery struct {
	Order             Order
	Route             Route
	OrderStatus       OrderStatus
	DeliveredAt       time.Time
	Vehicle           Vehicle
	Recipient         Recipient
	EvidencePhotos    []EvidencePhotos
	Latitude          float64
	Longitude         float64
	NotDeliveryReason NotDeliveryReason
}

type Recipient struct {
	FullName   string
	NationalID string
}

type EvidencePhotos struct {
	ID      int64
	URL     string
	Type    string
	TakenAt time.Time
}

type NotDeliveryReason struct {
	ID          int64
	ReferenceID string
	Detail      string
}
