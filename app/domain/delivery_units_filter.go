package domain

type ReferenceFilter struct {
	Type  string
	Value string
}

type LabelFilter struct {
	Type  string
	Value string
}

type CoordinatesConfidenceLevelFilter struct {
	Min *float64
	Max *float64
}

type PromisedDateRangeFilter struct {
	StartDate *string
	EndDate   *string
}

type DeliveryUnitsFilter struct {
	Pagination                 Pagination
	ReferenceIds               []string
	OriginNodeReferences       []string
	ReferenceType              *string
	ReferenceValue             *string
	Lpns                       []string
	GroupBy                    []string
	LabelType                  *string
	LabelValue                 *string
	Commerces                  []string
	Consumers                  []string
	RequestedFields            map[string]any
	References                 []ReferenceFilter
	Labels                     []LabelFilter
	CoordinatesConfidenceLevel *CoordinatesConfidenceLevelFilter
	PromisedDateRange          *PromisedDateRangeFilter
	OnlyLatestStatus           bool
}
