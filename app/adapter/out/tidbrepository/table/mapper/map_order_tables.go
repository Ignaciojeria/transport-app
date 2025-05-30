package mapper

import (
	"context"
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
		CollectAvailabilityTimeRangeStart: order.CollectAvailabilityDate.TimeRange.StartTime,
		CollectAvailabilityTimeRangeEnd:   order.CollectAvailabilityDate.TimeRange.EndTime,
		PromisedDateRangeStart:            safePtrTime(order.PromisedDate.DateRange.StartDate),
		PromisedDateRangeEnd:              safePtrTime(order.PromisedDate.DateRange.EndDate),
		PromisedTimeRangeStart:            order.PromisedDate.TimeRange.StartTime,
		PromisedTimeRangeEnd:              order.PromisedDate.TimeRange.EndTime,
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

func mapItemsToTable(items []domain.Item) table.JSONItems {
	mapped := make(table.JSONItems, len(items))
	for i, item := range items {
		mapped[i] = table.Items{
			Sku:            item.Sku,
			Skills:         item.Skills,
			QuantityNumber: item.Quantity.QuantityNumber,
			QuantityUnit:   item.Quantity.QuantityUnit,
			JSONInsurance: table.JSONInsurance{
				UnitValue: item.Insurance.UnitValue,
				Currency:  item.Insurance.Currency,
			},
			Description: item.Description,
			JSONDimensions: table.JSONDimensions{
				Height: item.Dimensions.Height,
				Width:  item.Dimensions.Width,
				Length: item.Dimensions.Length,
				Unit:   item.Dimensions.Unit,
			},
			JSONWeight: table.JSONWeight{
				WeightValue: item.Weight.Value,
				WeightUnit:  item.Weight.Unit,
			},
		}
	}

	return mapped
}

func MapPackagesToTable(ctx context.Context, packages []domain.DeliveryUnit) []table.DeliveryUnit {
	mapped := make([]table.DeliveryUnit, len(packages))
	for i, pkg := range packages {
		mapped[i] = table.DeliveryUnit{
			TenantID: sharedcontext.TenantIDFromContext(ctx),
			Lpn:      pkg.Lpn,
			JSONDimensions: table.JSONDimensions{
				Height: pkg.Dimensions.Height,
				Width:  pkg.Dimensions.Width,
				Length: pkg.Dimensions.Length,
				Unit:   pkg.Dimensions.Unit,
			},
			JSONWeight: table.JSONWeight{
				WeightValue: pkg.Weight.Value,
				WeightUnit:  pkg.Weight.Unit,
			},
			JSONInsurance: table.JSONInsurance{
				UnitValue: pkg.Insurance.UnitValue,
				Currency:  pkg.Insurance.Currency,
			},
		}
	}
	return mapped
}

func MapPackageToTable(ctx context.Context, pkg domain.DeliveryUnit) table.DeliveryUnit {
	return table.DeliveryUnit{
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
		DocumentID: pkg.DocID(ctx).String(),
		Lpn:        pkg.Lpn,
		JSONDimensions: table.JSONDimensions{
			Height: pkg.Dimensions.Height,
			Width:  pkg.Dimensions.Width,
			Length: pkg.Dimensions.Length,
			Unit:   pkg.Dimensions.Unit,
		},
		JSONWeight: table.JSONWeight{
			WeightValue: pkg.Weight.Value,
			WeightUnit:  pkg.Weight.Unit,
		},
		JSONInsurance: table.JSONInsurance{
			UnitValue: pkg.Insurance.UnitValue,
			Currency:  pkg.Insurance.Currency,
		},
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
		StateDoc:         address.State.DocID(ctx).String(),
		ProvinceDoc:      address.Province.DocID(ctx).String(),
		DistrictDoc:      address.District.DocID(ctx).String(),
		AddressLine1:     address.AddressLine1,
		DocumentID:       string(address.DocID(ctx)),
		Latitude:         address.Coordinates.Point.Lat(),
		Longitude:        address.Coordinates.Point.Lon(),
		ZipCode:          address.ZipCode,
		TimeZone:         address.TimeZone,
		CoordinateSource: address.Coordinates.Source,
	}
}
