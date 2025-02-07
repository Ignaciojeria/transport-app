package domain

import "time"

type Checkout struct {
	Order             Order
	Route             Route
	OrderStatus       OrderStatus
	DeliveredAt       time.Time
	Vehicle           Vehicle
	Recipient         Recipient
	EvidencePhotos    []EvidencePhotos
	Latitude          float32
	Longitude         float32
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
