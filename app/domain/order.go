package domain

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"time"
	"transport-app/app/shared/apperrors"

	"github.com/cockroachdb/errors"
)

type Order struct {
	Headers
	ExtraFields             map[string]string
	ReferenceID             ReferenceID
	SequenceNumber          *int
	OrderType               OrderType
	References              []Reference
	Origin                  NodeInfo
	Destination             NodeInfo
	DeliveryUnits           DeliveryUnits
	CollectAvailabilityDate CollectAvailabilityDate
	PromisedDate            PromisedDate
	UnassignedReason        string
	DeliveryInstructions    string
	GroupBy                 struct {
		Type  string
		Value string
	}
}

func (o *Order) AssignIndexesIfNoLPN() {
	o.DeliveryUnits.assignIndexesIfNoLPN(o.ReferenceID.String())
}

type DeliveryUnits []DeliveryUnit

func (pkgs *DeliveryUnits) assignIndexesIfNoLPN(referenceID string) {
	for i := range *pkgs {
		pkg := &(*pkgs)[i]
		if pkg.Lpn != "" {
			continue
		}
		pkg.noLPNReference = referenceID
	}
}

func (o Order) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, string(o.ReferenceID))
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

	if len(newOrder.DeliveryUnits) > 0 {
		o.DeliveryUnits = newOrder.DeliveryUnits
		changed = true
	}

	// GroupBy
	if newOrder.GroupBy.Type != "" && newOrder.GroupBy.Type != o.GroupBy.Type {
		o.GroupBy.Type = newOrder.GroupBy.Type
		changed = true
	}
	if newOrder.GroupBy.Value != "" && newOrder.GroupBy.Value != o.GroupBy.Value {
		o.GroupBy.Value = newOrder.GroupBy.Value
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

func (o Order) ValidateDeliveryUnits() error {
	// Mapa para trackear LPNs duplicados
	lpnSet := make(map[string]struct{})
	// Mapa para trackear conjuntos de SKUs duplicados
	skuSetMap := make(map[string]struct{})

	for _, pkg := range o.DeliveryUnits {
		// Validar LPNs duplicados
		if pkg.Lpn != "" {
			if _, exists := lpnSet[pkg.Lpn]; exists {
				return apperrors.MarkAsAlertable(
					errors.Wrap(ErrInvalidPackageFormat, "duplicate LPN found: each LPN must be unique within an order"))
			}
			lpnSet[pkg.Lpn] = struct{}{}
			continue
		}

		// Para paquetes sin LPN, validar conjuntos de SKUs duplicados
		if len(pkg.Items) > 0 {
			skus := make([]string, 0, len(pkg.Items))
			for _, item := range pkg.Items {
				skus = append(skus, item.Sku)
			}
			sort.Strings(skus)
			skuKey := strings.Join(skus, ",")

			if _, exists := skuSetMap[skuKey]; exists {
				return apperrors.MarkAsAlertable(
					errors.Wrap(ErrInvalidPackageFormat, "duplicate SKU set found: identical SKU combinations must be grouped in a single package"))
			}
			skuSetMap[skuKey] = struct{}{}
		}
	}

	return nil
}

func (o Order) Validate() error {
	// Validaciones existentes
	if err := o.ValidateCollectAvailabilityDate(); err != nil {
		return apperrors.MarkAsAlertable(errors.Wrap(err, "validation failed for CollectAvailabilityDate"))
	}
	if err := o.ValidatePromisedDate(); err != nil {
		return apperrors.MarkAsAlertable(errors.Wrap(err, "validation failed for PromisedDate"))
	}

	// Validar duplicados en unidades de entrega
	if err := o.ValidateDeliveryUnits(); err != nil {
		return err
	}

	// Validación sobre los paquetes
	for _, pkg := range o.DeliveryUnits {
		if pkg.Lpn == "" {
			if len(pkg.Items) == 0 {
				return apperrors.MarkAsAlertable(
					errors.Wrap(ErrInvalidPackageFormat, "no items found in package without LPN: package without LPN must have at least one item"))
			}
			for _, item := range pkg.Items {
				if item.Sku == "" {
					return apperrors.MarkAsAlertable(
						errors.Wrap(ErrInvalidPackageFormat, "item missing SKU: all items in package without LPN must have a non-empty SKU"))
				}
			}
		}
	}

	return nil
}

func (o Order) ValidatePromisedDate() error {
	timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

	// Validar que EndDate no sea anterior a StartDate cuando ambas están definidas
	if (o.PromisedDate.DateRange.StartDate != time.Time{}) && (o.PromisedDate.DateRange.EndDate != time.Time{}) {
		if o.PromisedDate.DateRange.EndDate.Before(o.PromisedDate.DateRange.StartDate) {
			return apperrors.MarkAsAlertable(errors.Wrap(ErrInvalidDateFormat, "invalid endDate: promised delivery endDate cannot be before startDate"))
		}
	}

	// Validar formato de horas
	if o.PromisedDate.TimeRange.StartTime != "" && !timeRegex.MatchString(o.PromisedDate.TimeRange.StartTime) {
		return apperrors.MarkAsAlertable(errors.Wrapf(
			ErrInvalidTimeFormat,
			"promised delivery startTime: %s, expected format hh:mm",
			o.PromisedDate.TimeRange.StartTime,
		))
	}
	if o.PromisedDate.TimeRange.EndTime != "" && !timeRegex.MatchString(o.PromisedDate.TimeRange.EndTime) {
		return apperrors.MarkAsAlertable(errors.Wrapf(
			ErrInvalidTimeFormat,
			"promised delivery endTime: %s, expected format hh:mm",
			o.PromisedDate.TimeRange.EndTime,
		))
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
		return apperrors.MarkAsAlertable(errors.Wrapf(
			ErrInvalidTimeFormat,
			"collect startTime: %s, expected format hh:mm",
			o.CollectAvailabilityDate.TimeRange.StartTime,
		))
	}
	if o.CollectAvailabilityDate.TimeRange.EndTime != "" && !timeRegex.MatchString(o.CollectAvailabilityDate.TimeRange.EndTime) {
		return apperrors.MarkAsAlertable(errors.Wrapf(
			ErrInvalidTimeFormat,
			"collect endTime: %s, expected format hh:mm",
			o.CollectAvailabilityDate.TimeRange.EndTime,
		))
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
	StartTime string `json:"startTime" example:"10:30"`
	EndTime   string `json:"endTime" example:"21:30"`
}

type DateRange struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
