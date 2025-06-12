package domain

type DeliveryUnitProjection struct {
	DestinationPoliticalArea PoliticalAreaProjection
	OriginPoliticalArea      PoliticalAreaProjection
}

type PoliticalAreaProjection struct {
	Code              string
	ConfidenceLevel   float64
	ConfidenceMessage string
	ConfidenceReason  string
}
