package domain

type OrderFilterInput struct {
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
	RequestedFields      []string
}
