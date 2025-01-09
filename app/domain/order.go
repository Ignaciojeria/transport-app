package domain

import (
	"regexp"

	"github.com/joomcode/errorx"
)

type Order struct {
	ID                      int64
	ReferenceID             ReferenceID             `json:"id"`
	Organization            Organization            `json:"organization"`
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

func (o Order) ValidatePromisedDate() error {
	// Validar formato de fecha yyyy-mm-dd
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

	// Validar el rango de fechas
	if o.PromisedDate.DateRange.StartDate != "" && !dateRegex.MatchString(o.PromisedDate.DateRange.StartDate) {
		return errorx.Decorate(
			ErrInvalidDateFormat.New("invalid startDate"),
			"startDate: %s, expected format yyyy-mm-dd",
			o.PromisedDate.DateRange.StartDate,
		)
	}
	if o.PromisedDate.DateRange.EndDate != "" && !dateRegex.MatchString(o.PromisedDate.DateRange.EndDate) {
		return errorx.Decorate(
			ErrInvalidDateFormat.New("invalid endDate"),
			"endDate: %s, expected format yyyy-mm-dd",
			o.PromisedDate.DateRange.EndDate,
		)
	}

	// Validar los rangos horarios
	if o.PromisedDate.TimeRange.StartTime != "" && !timeRegex.MatchString(o.PromisedDate.TimeRange.StartTime) {
		return errorx.Decorate(
			ErrInvalidTimeFormat.New("invalid startTime"),
			"startTime: %s, expected format hh:mm",
			o.PromisedDate.TimeRange.StartTime,
		)
	}
	if o.PromisedDate.TimeRange.EndTime != "" && !timeRegex.MatchString(o.PromisedDate.TimeRange.EndTime) {
		return errorx.Decorate(
			ErrInvalidTimeFormat.New("invalid endTime"),
			"endTime: %s, expected format hh:mm",
			o.PromisedDate.TimeRange.EndTime,
		)
	}

	return nil
}

func (o Order) ValidateCollectAvailabilityDate() error {
	// Validar formato de fecha yyyy-mm-dd
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

	// Validar la fecha
	if o.CollectAvailabilityDate.Date != "" && !dateRegex.MatchString(o.CollectAvailabilityDate.Date) {
		return errorx.Decorate(
			ErrInvalidDateFormat.New("invalid date"),
			"collect date: %s, expected format yyyy-mm-dd",
			o.CollectAvailabilityDate.Date,
		)
	}

	// Validar el rango horario
	if o.CollectAvailabilityDate.TimeRange.StartTime != "" && !timeRegex.MatchString(o.CollectAvailabilityDate.TimeRange.StartTime) {
		return errorx.Decorate(
			ErrInvalidTimeFormat.New("invalid startTime"),
			"collect startTime: %s, expected format hh:mm",
			o.CollectAvailabilityDate.TimeRange.StartTime,
		)
	}
	if o.CollectAvailabilityDate.TimeRange.EndTime != "" && !timeRegex.MatchString(o.CollectAvailabilityDate.TimeRange.EndTime) {
		return errorx.Decorate(
			ErrInvalidTimeFormat.New("invalid endTime"),
			"collect endTime: %s, expected format hh:mm",
			o.CollectAvailabilityDate.TimeRange.EndTime,
		)
	}

	return nil
}

func (o Order) IsOriginAndDestinationContactEqual() bool {
	originContact := o.Origin.AddressInfo.Contact
	destinationContact := o.Destination.AddressInfo.Contact

	return originContact.FullName == destinationContact.FullName &&
		originContact.Email == destinationContact.Email &&
		originContact.Phone == destinationContact.Phone &&
		originContact.NationalID == destinationContact.NationalID
}

func (o Order) IsOriginAndDestinationAddressEqual() bool {
	originAddress := o.Origin.AddressInfo.RawAddress()
	destinationAddress := o.Destination.AddressInfo.RawAddress()

	return originAddress == destinationAddress
}

func (o Order) IsOriginAndDestinationNodeEqual() bool {
	return o.Origin.NodeInfo.ReferenceID == o.Destination.NodeInfo.ReferenceID
}

type ReferenceID string

type References struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type NodeInfo struct {
	ReferenceID ReferenceID  `json:"referenceId"`
	Name        *string      `json:"name"`
	Type        string       `json:"type"`
	Operator    Operator     `json:"operator"`
	References  []References `json:"references"`
}

type Origin struct {
	NodeInfo    NodeInfo    `json:"nodeInfo"`
	AddressInfo AddressInfo `json:"addressInfo"`
}

type Documents struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Contact struct {
	FullName   string      `json:"fullName"`
	Email      string      `json:"email"`
	Phone      string      `json:"phone"`
	NationalID string      `json:"nationalID"`
	Documents  []Documents `json:"documents"`
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

func (addr AddressInfo) RawAddress() string {
	return concatenateWithCommas(addr.AddressLine1, addr.AddressLine2, addr.AddressLine3)
}

func concatenateWithCommas(values ...string) string {
	result := ""
	for _, value := range values {
		if value != "" {
			if result != "" {
				result += ", "
			}
			result += value
		}
	}
	return result
}

type Operator struct {
	Contact Contact `json:"contact"`
	Type    string  `json:"type"`
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
	ID        int64
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

type BusinessIdentifiers struct {
	Commerce string `json:"commerce"`
	Consumer string `json:"consumer"`
}
