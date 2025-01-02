package domain

import "github.com/biter777/countries"

type Order struct {
	ID                      int64
	ReferenceID             ReferenceID             `json:"id"`
	Tenant                  Tenant                  `json:"tenant"`
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
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Operator struct {
	ReferenceID ReferenceID `json:"referenceId"`
	NationalID  string      `json:"nationalId"`
	Type        string      `json:"type"`
	Name        string      `json:"name"`
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
	Province     string  `json:"province"`
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
	ID        int
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
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type DateRange struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type Visit struct {
	Date      string    `json:"date"`
	TimeRange TimeRange `json:"timeRange"`
}

type Tenant struct {
	Country      countries.CountryCode `json:"country"`
	Organization string                `json:"organization"`
	Commerce     string                `json:"commerce"`
	Consumer     string                `json:"consumer"`
}
