package domain

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/joomcode/errorx"
)

type Order struct {
	ID                      int64
	ReferenceID             ReferenceID             `json:"referenceID"`
	Organization            Organization            `json:"organization"`
	BusinessIdentifiers     BusinessIdentifiers     `json:"businessIdentifiers"`
	OrderStatus             OrderStatus             `json:"orderStatus"`
	OrderType               OrderType               `json:"orderType"`
	References              []Reference             `json:"references"`
	Origin                  Origin                  `json:"origin"`
	Destination             Destination             `json:"destination"`
	Items                   []Item                  `json:"items"`
	Packages                []Package               `json:"packages"`
	CollectAvailabilityDate CollectAvailabilityDate `json:"collectAvailabilityDate"`
	PromisedDate            PromisedDate            `json:"promisedDate"`
	Visits                  []Visit                 `json:"visits"`
	TransportRequirements   []Reference             `json:"transportRequirements"`
}

func (o Order) Validate() error {
	// Validar fechas de disponibilidad de recolección
	if err := o.ValidateCollectAvailabilityDate(); err != nil {
		return errorx.Decorate(err, "validation failed for CollectAvailabilityDate")
	}

	// Validar fechas prometidas
	if err := o.ValidatePromisedDate(); err != nil {
		return errorx.Decorate(err, "validation failed for PromisedDate")
	}

	if err := o.ValidatePackages(); err != nil {
		return errorx.Decorate(err, "validation failed for Packages")
	}

	// Validar otras reglas de dominio (si las hay)
	// Por ejemplo, puedes agregar reglas adicionales aquí

	return nil
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

type Reference struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type NodeInfo struct {
	ReferenceID ReferenceID `json:"referenceId"`
	Name        *string     `json:"name"`
	Type        string      `json:"type"`
	Operator    Operator    `json:"operator"`
	References  []Reference `json:"references"`
}

type Origin struct {
	NodeInfo    NodeInfo    `json:"nodeInfo"`
	AddressInfo AddressInfo `json:"addressInfo"`
}

type Document struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Contact struct {
	FullName   string     `json:"fullName"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	NationalID string     `json:"nationalID"`
	Documents  []Document `json:"documents"`
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
	Latitude     float32 `json:"latitude"`
	Longitude    float32 `json:"longitude"`
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
	UnitValue float64 `json:"unitValue"`
	Currency  string  `json:"currency"`
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

type Item struct {
	ReferenceID       ReferenceID `json:"referenceId"`
	LogisticCondition string      `json:"logisticCondition"`
	Quantity          Quantity    `json:"quantity"`
	Insurance         Insurance   `json:"insurance"`
	Description       string      `json:"description"`
	Dimensions        Dimensions  `json:"dimensions"`
	Weight            Weight      `json:"weight"`
}

type ItemReference struct {
	ReferenceID ReferenceID `json:"referenceId"`
	Quantity    Quantity    `json:"quantity"`
}

type Package struct {
	ID             int64
	Lpn            string          `json:"lpn"`
	Dimensions     Dimensions      `json:"dimensions"`
	Weight         Weight          `json:"weight"`
	Insurance      Insurance       `json:"insurance"`
	ItemReferences []ItemReference `json:"itemReferences"`
}

func (o *Order) ValidatePackages() error {
	// Crear un mapa para verificar rápidamente si un ReferenceID pertenece a los ítems de la orden
	itemMap := make(map[ReferenceID]bool)
	for _, item := range o.Items {
		itemMap[item.ReferenceID] = true
	}

	if len(o.Packages) == 1 {
		// Si solo hay un paquete y no tiene referencias de ítems, asignar todos los ítems de la orden
		if len(o.Packages[0].ItemReferences) == 0 {
			for _, item := range o.Items {
				o.Packages[0].ItemReferences = append(o.Packages[0].ItemReferences, ItemReference{
					ReferenceID: item.ReferenceID,
					Quantity:    item.Quantity,
				})
			}
		} else {
			// Validar que todas las referencias del paquete sean válidas
			for _, ref := range o.Packages[0].ItemReferences {
				if !itemMap[ref.ReferenceID] {
					return fmt.Errorf("validation failed: item reference ID '%s' in package is not part of the order items", ref.ReferenceID)
				}
			}
		}
	} else {
		// Si hay más de un paquete, validar que todos los paquetes tengan referencias de ítems
		for _, p := range o.Packages {
			if len(p.ItemReferences) == 0 {
				return errors.New("validation failed: packages with no item references must be explicitly defined when there are multiple packages")
			}

			// Validar que todas las referencias del paquete sean válidas
			for _, ref := range p.ItemReferences {
				if !itemMap[ref.ReferenceID] {
					return fmt.Errorf("validation failed: item reference ID '%s' in package is not part of the order items", ref.ReferenceID)
				}
			}
		}
	}

	return nil
}

func (p *Package) UpdateIfChanged(newPackage Package) (updatedPackage Package, needsUpdate bool) {
	// Crear una copia del paquete actual
	updatedPackage = *p
	needsUpdate = false

	// Comparar y actualizar dimensiones
	if newPackage.Dimensions != (Dimensions{}) {
		if p.Dimensions != newPackage.Dimensions {
			updatedPackage.Dimensions = newPackage.Dimensions
			needsUpdate = true
		}
	}

	// Comparar y actualizar peso
	if newPackage.Weight != (Weight{}) {
		if p.Weight != newPackage.Weight {
			updatedPackage.Weight = newPackage.Weight
			needsUpdate = true
		}
	}

	// Comparar y actualizar seguro
	if newPackage.Insurance != (Insurance{}) {
		if p.Insurance != newPackage.Insurance {
			updatedPackage.Insurance = newPackage.Insurance
			needsUpdate = true
		}
	}

	// Comparar y actualizar referencias de ítems
	if len(newPackage.ItemReferences) > 0 {
		if len(p.ItemReferences) != len(newPackage.ItemReferences) || !compareItemReferences(p.ItemReferences, newPackage.ItemReferences) {
			updatedPackage.ItemReferences = newPackage.ItemReferences
			needsUpdate = true
		}
	}

	return updatedPackage, needsUpdate
}

// Función auxiliar para comparar arreglos de referencias de ítems
func compareItemReferences(oldRefs, newRefs []ItemReference) bool {
	if len(oldRefs) != len(newRefs) {
		return false
	}
	for i := range oldRefs {
		if oldRefs[i] != newRefs[i] {
			return false
		}
	}
	return true
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
