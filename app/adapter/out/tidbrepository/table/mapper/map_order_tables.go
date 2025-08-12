package mapper

import (
	"context"
	"strconv"
	"strings"
	"time"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	"github.com/google/uuid"
)

func MapOrderToTable(ctx context.Context, order domain.Order) table.Order {
	tbl := table.Order{
		TenantID:               sharedcontext.TenantIDFromContext(ctx),
		ReferenceID:            string(order.ReferenceID),
		DocumentID:             order.DocID(ctx).String(),
		OrderHeadersDoc:        order.Headers.DocID(ctx).String(),
		OrderTypeDoc:           order.OrderType.DocID(ctx).String(),
		OriginNodeInfoDoc:      order.Origin.DocID(ctx).String(),
		DestinationNodeInfoDoc: order.Destination.DocID(ctx).String(),
		ServiceCategory:        order.PromisedDate.ServiceCategory,
		// Si están disponibles, también mapear los contactos y direcciones
		OriginContactDoc:                  order.Origin.AddressInfo.Contact.DocID(ctx).String(),
		DestinationContactDoc:             order.Destination.AddressInfo.Contact.DocID(ctx).String(),
		OriginAddressInfoDoc:              order.Origin.AddressInfo.DocID(ctx).String(),
		DestinationAddressInfoDoc:         order.Destination.AddressInfo.DocID(ctx).String(),
		ExtraFields:                       order.ExtraFields,
		DeliveryInstructions:              order.DeliveryInstructions,
		CollectAvailabilityDate:           safePtrTime(order.CollectAvailabilityDate.Date),
		CollectAvailabilityTimeRangeStart: safePtrString(order.CollectAvailabilityDate.TimeRange.StartTime),
		CollectAvailabilityTimeRangeEnd:   safePtrString(order.CollectAvailabilityDate.TimeRange.EndTime),
		PromisedDateRangeStart:            safePtrTime(order.PromisedDate.DateRange.StartDate),
		PromisedDateRangeEnd:              safePtrTime(order.PromisedDate.DateRange.EndDate),
		PromisedTimeRangeStart:            safePtrString(order.PromisedDate.TimeRange.StartTime),
		PromisedTimeRangeEnd:              safePtrString(order.PromisedDate.TimeRange.EndTime),
		GroupByType:                       order.GroupBy.Type,
		GroupByValue:                      order.GroupBy.Value,
	}
	return tbl
}

func safePtrTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil // Retorna nil si la fecha es vacía en el dominio
	}
	return &t
}

func safePtrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func parseTime(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return nil
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return nil
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return nil
	}

	// Crear hora sin fecha ni zona
	t := time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC)
	return &t
}

func mapItemsToTable(items []domain.Item) table.JSONItems {
	mapped := make(table.JSONItems, len(items))
	for i, item := range items {

		mapped[i] = table.Items{
			Sku:         item.Sku,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Description: item.Description,
			JSONDimensions: table.JSONDimensions{
				Height: item.Dimensions.Height,
				Width:  item.Dimensions.Width,
				Length: item.Dimensions.Length,
				Unit:   item.Dimensions.Unit,
			},
			Weight: item.Weight,
		}
	}

	return mapped
}

func MapPackagesToTable(ctx context.Context, packages []domain.DeliveryUnit) []table.DeliveryUnit {
	mapped := make([]table.DeliveryUnit, len(packages))
	for i, pkg := range packages {
		var vol, wgt, ins int64
		if pkg.Volume != nil {
			vol = *pkg.Volume
		}
		if pkg.Weight != nil {
			wgt = *pkg.Weight
		}
		if pkg.Price != nil {
			ins = *pkg.Price
		}
		mapped[i] = table.DeliveryUnit{
			TenantID:  sharedcontext.TenantIDFromContext(ctx),
			Lpn:       pkg.Lpn,
			Volume:    vol,
			Weight:    wgt,
			Price:     ins,
			JSONItems: mapItemsToTable(pkg.Items),
		}
	}
	return mapped
}

func MapPackageToTable(ctx context.Context, pkg domain.DeliveryUnit) table.DeliveryUnit {
	var vol, wgt, ins int64
	if pkg.Volume != nil {
		vol = *pkg.Volume
	}
	if pkg.Weight != nil {
		wgt = *pkg.Weight
	}
	if pkg.Price != nil {
		ins = *pkg.Price
	}
	return table.DeliveryUnit{
		TenantID:        sharedcontext.TenantIDFromContext(ctx),
		DocumentID:      pkg.DocID(ctx).String(),
		Lpn:             pkg.Lpn,
		Volume:          vol,
		Weight:          wgt,
		Price:           ins,
		JSONItems:       mapItemsToTable(pkg.Items),
		SizeCategoryDoc: pkg.SizeCategory.DocumentID(ctx).String(),
	}
}

func mapHeadersToTable(c domain.Headers, orgID uuid.UUID) table.OrderHeaders {
	return table.OrderHeaders{
		Commerce: c.Commerce,
		Consumer: c.Consumer,
		Tenant: table.Tenant{
			ID: orgID,
		},
	}
}

func mapOrderTypeToTable(t domain.OrderType, orgID uuid.UUID) table.OrderType {
	return table.OrderType{
		TenantID:    orgID,
		Type:        t.Type,
		Description: t.Description,
	}
}

func MapAddressInfoToTable(ctx context.Context, address domain.AddressInfo) table.AddressInfo {
	return table.AddressInfo{
		TenantID:         sharedcontext.TenantIDFromContext(ctx),
		PoliticalAreaDoc: address.PoliticalArea.DocID(ctx).String(),
		AddressLine1:     address.AddressLine1,
		DocumentID:       string(address.DocID(ctx)),
		Latitude:         address.Coordinates.Point.Lat(),
		Longitude:        address.Coordinates.Point.Lon(),
		ZipCode:          address.ZipCode,
		CoordinateSource: address.Coordinates.Source,
	}
}
