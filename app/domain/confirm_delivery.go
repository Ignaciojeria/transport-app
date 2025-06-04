package domain

import "time"

type ConfirmDelivery struct {
	RequestBody       []byte
	ManualChange      ManualChange
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

type ManualChange struct {
	PerformedBy string
	Reason      string
}
