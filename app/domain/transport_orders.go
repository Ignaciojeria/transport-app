package domain

type TransportOrder struct {
	ReferenceID             ReferenceID             `json:"id" validate:"required"`
	BusinessIdentifiers     BusinessIdentifiers     `json:"businessIdentifiers"`
	OrderStatus             OrderStatus             `json:"orderStatus"`
	OrderType               OrderType               `json:"orderType"`
	References              []References            `json:"references"`
	Origin                  Origin                  `json:"origin"`
	Destination             Destination             `json:"destination"`
	Items                   []Items                 `json:"items"`
	Packages                []Packages              `json:"packages"`
	CollectAvailabilityDate CollectAvailabilityDate `json:"collectAvailabilityDate"`
	PromisedDate            PromisedDate            `json:"promisedDate"`
	Visit                   Visit                   `json:"visit"`
	TransportRequirements   []References            `json:"transportRequirements"`
}

type ReferenceID string

type References struct {
	Type  string `json:"type" validate:"required"`
	Value string `json:"value"`
}

type Operator struct {
	ReferenceID ReferenceID `json:"referenceId"`
	NationalID  string      `json:"nationalId"`
	Type        string      `json:"type" validate:"required"`
	Name        string      `json:"name" validate:"required"`
}

type NodeInfo struct {
	ReferenceID ReferenceID  `json:"referenceId"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Operator    Operator     `json:"operator"`
	References  []References `json:"references"`
}

type ContactMethods struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Documents struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Contact struct {
	FullName       string           `json:"fullName"`
	ContactMethods []ContactMethods `json:"contactmethods"`
	Documents      []Documents      `json:"documents"`
}

type AddressInfo struct {
	Contact      Contact `json:"contact"`
	State        string  `json:"state"`
	County       string  `json:"county"`
	District     string  `json:"district"`
	AddressLine1 string  `json:"addressLine1"`
	AddressLine2 string  `json:"addressLine2"`
	AddressLine3 string  `json:"addressLine3"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	ZipCode      string  `json:"zipCode"`
	TimeZone     string  `json:"timeZone"`
}

type Origin struct {
	NodeInfo    NodeInfo    `json:"nodeInfo"`
	AddressInfo AddressInfo `json:"addressInfo"`
}

type Destination struct {
	DeliveryInstructions string      `json:"deliveryInstructions"`
	NodeInfo             NodeInfo    `json:"nodeInfo"`
	AddressInfo          AddressInfo `json:"addressInfo"`
}

type Quantity struct {
	QuantityNumber int    `json:"quantityNumber"`
	QuantityUnit   string `json:"quantityUnit"`
}

type Insurance struct {
	UnitValue int    `json:"unitValue"`
	Currency  string `json:"currency"`
}

type Dimensions struct {
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Depth  float64 `json:"depth"`
	Unit   string  `json:"unit"`
}

type Weight struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type Items struct {
	ReferenceID       ReferenceID `json:"referenceId"`
	LogisticCondition string      `json:"logisticCondition"`
	Quantity          Quantity    `json:"quantity"`
	Insurance         Insurance   `json:"insurance"`
	Description       string      `json:"description"`
	Dimensions        Dimensions  `json:"dimensions"`
	Weight            Weight      `json:"weight"`
}

type ItemReferences struct {
	ReferenceID ReferenceID `json:"referenceId"`
	Quantity    Quantity    `json:"quantity"`
}

type Packages struct {
	Lpn            string           `json:"lpn"`
	PackageType    string           `json:"packageType"`
	Dimensions     Dimensions       `json:"dimensions"`
	Weight         Weight           `json:"weight"`
	Insurance      Insurance        `json:"insurance"`
	ItemReferences []ItemReferences `json:"itemReferences"`
}

type OrderType struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type OrderStatus struct {
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
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

type TimeRange struct {
	Start string `json:"start" validate:"required,datetime=15:04" example:"9:00"`
	End   string `json:"end" validate:"required,datetime=15:04" example:"23:00"`
}

type DateRange struct {
	StartDate string `json:"startDate" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"endDate" validate:"required,datetime=2006-01-02"`
}

type Visit struct {
	Date      string    `json:"date"`
	TimeRange TimeRange `json:"timeRange"`
}

type BusinessIdentifiers struct {
	Organization string `json:"organization"`
	Commerce     string `json:"commerce"`
	Consumer     string `json:"consumer"`
}
