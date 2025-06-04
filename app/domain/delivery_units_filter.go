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

type OrderFilter struct {
	ReferenceIds []string
	References   []ReferenceFilter
	OrderType    *OrderTypeFilter
	GroupBy      *GroupByFilter
}

type OrderTypeFilter struct {
	Type        string
	Description string
}

type GroupByFilter struct {
	Type  string
	Value string
}

type DeliveryUnitFilter struct {
	Lpns           []string
	SizeCategories []string
	Labels         []LabelFilter
}

type LocationFilter struct {
	NodeReferences        []string
	AddressLines          []string
	Districts             []string
	Provinces             []string
	States                []string
	ZipCodes              []string
	CoordinatesConfidence *CoordinatesConfidenceLevelFilter
}

type PromisedDateFilter struct {
	DateRange *DateRangeFilter
	TimeRange *TimeRangeFilter
}

type DateRangeFilter struct {
	StartDate *string
	EndDate   *string
}

type TimeRangeFilter struct {
	StartTime *string
	EndTime   *string
}

type CollectAvailabilityFilter struct {
	Dates     []string
	TimeRange *TimeRangeFilter
}

type DeliveryUnitsFilter struct {
	Pagination          Pagination
	RequestedFields     map[string]any
	OnlyLatestStatus    bool
	Order               *OrderFilter
	DeliveryUnit        *DeliveryUnitFilter
	Origin              *LocationFilter
	Destination         *LocationFilter
	PromisedDate        *PromisedDateFilter
	CollectAvailability *CollectAvailabilityFilter
}
