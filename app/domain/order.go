package domain

import (
	"context"
	"regexp"
	"time"

	"github.com/joomcode/errorx"
)

type Order struct {
	Headers
	AddressLine2            string
	AddressLine3            string
	ReferenceID             ReferenceID
	SequenceNumber          *int
	OrderStatus             OrderStatus
	OrderType               OrderType
	References              []Reference
	Origin                  NodeInfo
	Destination             NodeInfo
	Packages                []Package
	CollectAvailabilityDate CollectAvailabilityDate
	PromisedDate            PromisedDate
	UnassignedReason        string
	DeliveryInstructions    string
	TransportRequirements   []Reference
}

func (o Order) DocID(ctx context.Context) DocumentID {
	return Hash(ctx, string(o.ReferenceID))
}

func (o Order) UpdateIfChanged(newOrder Order) (Order, bool) {
	changed := false

	if newOrder.DeliveryInstructions != "" && newOrder.DeliveryInstructions != o.DeliveryInstructions {
		o.DeliveryInstructions = newOrder.DeliveryInstructions
		changed = true
	}

	if len(newOrder.References) > 0 {
		o.References = newOrder.References
		changed = true
	}

	if len(newOrder.Packages) > 0 {
		o.Packages = newOrder.Packages
		changed = true
	}

	if len(newOrder.TransportRequirements) > 0 {
		o.TransportRequirements = newOrder.TransportRequirements
		changed = true
	}

	// PromisedDate
	if !newOrder.PromisedDate.DateRange.StartDate.IsZero() && !newOrder.PromisedDate.DateRange.StartDate.Equal(o.PromisedDate.DateRange.StartDate) {
		o.PromisedDate.DateRange.StartDate = newOrder.PromisedDate.DateRange.StartDate
		changed = true
	}
	if !newOrder.PromisedDate.DateRange.EndDate.IsZero() && !newOrder.PromisedDate.DateRange.EndDate.Equal(o.PromisedDate.DateRange.EndDate) {
		o.PromisedDate.DateRange.EndDate = newOrder.PromisedDate.DateRange.EndDate
		changed = true
	}
	if newOrder.PromisedDate.TimeRange.StartTime != "" && newOrder.PromisedDate.TimeRange.StartTime != o.PromisedDate.TimeRange.StartTime {
		o.PromisedDate.TimeRange.StartTime = newOrder.PromisedDate.TimeRange.StartTime
		changed = true
	}
	if newOrder.PromisedDate.TimeRange.EndTime != "" && newOrder.PromisedDate.TimeRange.EndTime != o.PromisedDate.TimeRange.EndTime {
		o.PromisedDate.TimeRange.EndTime = newOrder.PromisedDate.TimeRange.EndTime
		changed = true
	}
	if newOrder.PromisedDate.ServiceCategory != "" && newOrder.PromisedDate.ServiceCategory != o.PromisedDate.ServiceCategory {
		o.PromisedDate.ServiceCategory = newOrder.PromisedDate.ServiceCategory
		changed = true
	}

	// CollectAvailabilityDate
	if !newOrder.CollectAvailabilityDate.Date.IsZero() && !newOrder.CollectAvailabilityDate.Date.Equal(o.CollectAvailabilityDate.Date) {
		o.CollectAvailabilityDate.Date = newOrder.CollectAvailabilityDate.Date
		changed = true
	}
	if newOrder.CollectAvailabilityDate.TimeRange.StartTime != "" && newOrder.CollectAvailabilityDate.TimeRange.StartTime != o.CollectAvailabilityDate.TimeRange.StartTime {
		o.CollectAvailabilityDate.TimeRange.StartTime = newOrder.CollectAvailabilityDate.TimeRange.StartTime
		changed = true
	}
	if newOrder.CollectAvailabilityDate.TimeRange.EndTime != "" && newOrder.CollectAvailabilityDate.TimeRange.EndTime != o.CollectAvailabilityDate.TimeRange.EndTime {
		o.CollectAvailabilityDate.TimeRange.EndTime = newOrder.CollectAvailabilityDate.TimeRange.EndTime
		changed = true
	}

	return o, changed
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
