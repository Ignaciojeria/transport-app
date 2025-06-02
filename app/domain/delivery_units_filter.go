package domain

type ReferenceFilter struct {
	Type  string
	Value string
}

type DeliveryUnitsFilter struct {
	Pagination           Pagination
	ReferenceIds         []string
	OriginNodeReferences []string
	ReferenceType        *string
	ReferenceValue       *string
	Lpns                 []string
	GroupBy              []string
	LabelType            *string
	LabelValue           *string
	Commerces            []string
	Consumers            []string
	RequestedFields      map[string]any
	References           []ReferenceFilter
}
