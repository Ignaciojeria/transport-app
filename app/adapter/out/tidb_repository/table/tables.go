package table

type TransportOrder struct {
	ID                      int64
	ReferenceID             string
	BusinessIdentifiersID   int
	BusinessIdentifiers     BusinessIdentifiers
	OrderStatusID           int
	OrderStatus             OrderStatus
	OrderType               OrderType
	References              []References
	DeliveryInstructions    string
	OriginNodeInfoID        int64
	OriginNodeInfo          NodeInfo
	OriginAddressInfo       AddressInfo
	DestinationNodeInfoID   int64
	DestinationNodeInfo     NodeInfo
	DestinationAddressInfo  AddressInfo
	Items                   []Items
	Packages                []Packages
	CollectAvailabilityDate CollectAvailabilityDate
	PromisedDate            PromisedDate
	Visit                   Visit
	TransportRequirements   []References
}

type References struct {
	Type  string
	Value string
}

type Operator struct {
	ID          int64
	ReferenceID string
	NationalID  string
	Type        string
	Name        string
}

type NodeInfo struct {
	ID          int64
	ReferenceID string
	Name        string
	Type        string
	Operator    Operator
	References  []References
}

type Contacts struct {
	Value string
	Type  string
}

type Documents struct {
	Value string
	Type  string
}

type Contact struct {
	FullName  string
	Contacts  []Contacts
	Documents []Documents
}

type AddressInfo struct {
	Contact     Contact
	State       string
	County      string
	District    string
	FullAddress string
	Latitude    float64
	Longitude   float64
	ZipCode     string
	TimeZone    string
}

type Quantity struct {
	QuantityNumber int
	QuantityUnit   string
}

type Insurance struct {
	UnitValue int
	Currency  string
}

type Dimensions struct {
	Height float64
	Width  float64
	Depth  float64
	Unit   string
}

type Weight struct {
	Value float64
	Unit  string
}

type Items struct {
	ReferenceID       string
	LogisticCondition string
	Quantity          Quantity
	Insurance         Insurance
	Description       string
	Dimensions        Dimensions
	Weight            Weight
}

type ItemReferences struct {
	ReferenceID string
	Quantity    Quantity
}

type Packages struct {
	ID             int64
	Lpn            string
	PackageType    string
	Dimensions     Dimensions
	Weight         Weight
	Insurance      Insurance
	ItemReferences []ItemReferences
}

type OrderType struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type OrderStatus struct {
	ID     int
	Status string `json:"status"`
}

type PromisedDate struct {
	DateRange       DateRange `json:"dateRange"`
	TimeRange       TimeRange `json:"timeRange"`
	ServiceCategory string    `json:"serviceCategory"`
}

type CollectAvailabilityDate struct {
	Date      string    `json:"date"`
	TimeRange TimeRange `json:"timeRange"`
}

type PlannedDispatchDate struct {
	Date      string    `json:"date"`
	TimeRange TimeRange `json:"timeRange"`
}

type TimeRange struct {
	Start string `json:"start" validate:"required,datetime=15:04" example:"9:00"`
	End   string `json:"end" validate:"required,datetime=15:04" example:"23:00"`
}

type DateRange struct {
	StartDate string `json:"startDate" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"endDate" validate:"required,datetime=2006-01-02"`
}

type Visit struct {
	ID                  int64
	PlannedDispatchDate PlannedDispatchDate `json:"plannedDispatchDate"`
}

type BusinessIdentifiers struct {
	ID       int
	Commerce string `json:"commerce"`
	Consumer string `json:"consumer"`
}
