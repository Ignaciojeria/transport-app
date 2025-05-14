package domain

import "time"

type ConfirmDelivery struct {
	RequestBody       []byte
	HandledAt         time.Time
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
