package domain

import "time"

type OrderHistory struct {
	RequestBody       []byte
	Order             Order
	Route             Route
	OrderStatus       OrderStatus
	HandledAt         time.Time
	Vehicle           Vehicle
	Recipient         Recipient
	EvidencePhotos    []EvidencePhotos
	Latitude          float64
	Longitude         float64
	NonDeliveryReason NonDeliveryReason
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
