package domain

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/joomcode/errorx"
)

type Order struct {
	Headers
	Plan                    Plan
	AddressLine2            string
	AddressLine3            string
	ReferenceID             ReferenceID
	SequenceNumber          *int
	OrderStatus             OrderStatus
	OrderType               OrderType
	References              []Reference
	Origin                  NodeInfo
	Destination             NodeInfo
	Items                   []Item
	Packages                []Package
	CollectAvailabilityDate CollectAvailabilityDate
	PromisedDate            PromisedDate
	UnassignedReason        string
	DeliveryInstructions    string
	TransportRequirements   []Reference
}

func (o Order) DocID() DocumentID {
	return Hash(o.Organization, string(o.ReferenceID))
}

func (o Order) UpdateIfChanged(newOrder Order) (Order, bool) {
	updated := o
	changed := false

	if newOrder.DeliveryInstructions != "" && newOrder.DeliveryInstructions != o.DeliveryInstructions {
		updated.DeliveryInstructions = newOrder.DeliveryInstructions
		changed = true
	}

	if len(newOrder.References) > 0 {
		updated.References = newOrder.References
		changed = true
	}

	if len(newOrder.Packages) > 0 {
		updated.Packages = newOrder.Packages
		changed = true
	}

	if len(newOrder.Items) > 0 {
		updated.Items = newOrder.Items
		changed = true
	}

	if len(newOrder.TransportRequirements) > 0 {
		updated.TransportRequirements = newOrder.TransportRequirements
		changed = true
	}

	// PromisedDate
	if !newOrder.PromisedDate.DateRange.StartDate.IsZero() && !newOrder.PromisedDate.DateRange.StartDate.Equal(o.PromisedDate.DateRange.StartDate) {
		updated.PromisedDate.DateRange.StartDate = newOrder.PromisedDate.DateRange.StartDate
		changed = true
	}
	if !newOrder.PromisedDate.DateRange.EndDate.IsZero() && !newOrder.PromisedDate.DateRange.EndDate.Equal(o.PromisedDate.DateRange.EndDate) {
		updated.PromisedDate.DateRange.EndDate = newOrder.PromisedDate.DateRange.EndDate
		changed = true
	}
	if newOrder.PromisedDate.TimeRange.StartTime != "" && newOrder.PromisedDate.TimeRange.StartTime != o.PromisedDate.TimeRange.StartTime {
		updated.PromisedDate.TimeRange.StartTime = newOrder.PromisedDate.TimeRange.StartTime
		changed = true
	}
	if newOrder.PromisedDate.TimeRange.EndTime != "" && newOrder.PromisedDate.TimeRange.EndTime != o.PromisedDate.TimeRange.EndTime {
		updated.PromisedDate.TimeRange.EndTime = newOrder.PromisedDate.TimeRange.EndTime
		changed = true
	}
	if newOrder.PromisedDate.ServiceCategory != "" && newOrder.PromisedDate.ServiceCategory != o.PromisedDate.ServiceCategory {
		updated.PromisedDate.ServiceCategory = newOrder.PromisedDate.ServiceCategory
		changed = true
	}

	// CollectAvailabilityDate
	if !newOrder.CollectAvailabilityDate.Date.IsZero() && !newOrder.CollectAvailabilityDate.Date.Equal(o.CollectAvailabilityDate.Date) {
		updated.CollectAvailabilityDate.Date = newOrder.CollectAvailabilityDate.Date
		changed = true
	}
	if newOrder.CollectAvailabilityDate.TimeRange.StartTime != "" && newOrder.CollectAvailabilityDate.TimeRange.StartTime != o.CollectAvailabilityDate.TimeRange.StartTime {
		updated.CollectAvailabilityDate.TimeRange.StartTime = newOrder.CollectAvailabilityDate.TimeRange.StartTime
		changed = true
	}
	if newOrder.CollectAvailabilityDate.TimeRange.EndTime != "" && newOrder.CollectAvailabilityDate.TimeRange.EndTime != o.CollectAvailabilityDate.TimeRange.EndTime {
		updated.CollectAvailabilityDate.TimeRange.EndTime = newOrder.CollectAvailabilityDate.TimeRange.EndTime
		changed = true
	}

	return updated, changed
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
	timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

	// Validar que EndDate no sea anterior a StartDate cuando ambas están definidas
	if (o.PromisedDate.DateRange.StartDate != time.Time{}) && (o.PromisedDate.DateRange.EndDate != time.Time{}) {
		if o.PromisedDate.DateRange.EndDate.Before(o.PromisedDate.DateRange.StartDate) {
			return errorx.Decorate(
				ErrInvalidDateFormat.New("invalid endDate"),
				"promised delivery endDate cannot be before startDate",
			)
		}
	}

	// Validar formato de horas
	if o.PromisedDate.TimeRange.StartTime != "" && !timeRegex.MatchString(o.PromisedDate.TimeRange.StartTime) {
		return errorx.Decorate(
			ErrInvalidTimeFormat.New("invalid startTime"),
			"promised delivery startTime: %s, expected format hh:mm",
			o.PromisedDate.TimeRange.StartTime,
		)
	}
	if o.PromisedDate.TimeRange.EndTime != "" && !timeRegex.MatchString(o.PromisedDate.TimeRange.EndTime) {
		return errorx.Decorate(
			ErrInvalidTimeFormat.New("invalid endTime"),
			"promised delivery endTime: %s, expected format hh:mm",
			o.PromisedDate.TimeRange.EndTime,
		)
	}

	return nil
}

func (o Order) ValidateCollectAvailabilityDate() error {
	// Si la fecha está definida (no es zero value), solo validamos que sea una fecha válida
	// Ya no validamos si está en el pasado
	if (o.CollectAvailabilityDate.Date != time.Time{}) {
		// Aquí podrías agregar otras validaciones específicas si las necesitas
		// pero no la comparación con time.Now()
	}

	// Validar el rango horario
	timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)
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
		originContact.PrimaryEmail == destinationContact.PrimaryEmail &&
		originContact.PrimaryPhone == destinationContact.PrimaryPhone &&
		originContact.NationalID == destinationContact.NationalID
}

func (o Order) IsOriginAndDestinationAddressEqual() bool {
	originAddress := o.Origin.AddressInfo.FullAddress()
	destinationAddress := o.Destination.AddressInfo.FullAddress()

	return originAddress == destinationAddress
}

func (o Order) IsOriginAndDestinationNodeEqual() bool {
	return o.Origin.ReferenceID == o.Destination.ReferenceID
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

type PromisedDate struct {
	DateRange       DateRange `json:"dateRange"`
	TimeRange       TimeRange `json:"timeRange"`
	ServiceCategory string    `json:"serviceCategory"`
}

type CollectAvailabilityDate struct {
	Date      time.Time `json:"date"`
	TimeRange TimeRange `json:"timeRange"`
}

type TimeRange struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type DateRange struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
