package domain

import "time"

type ConfirmDelivery struct {
	RequestBody       []byte
	HandledAt         time.Time
	Recipient         Recipient
	EvidencePhotos    EvidencePhotos
	Latitude          float64
	Longitude         float64
	NonDeliveryReason NonDeliveryReason
}

type Recipient struct {
	FullName   string
	NationalID string
}
