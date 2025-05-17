package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapAddressInfoTable(ctx context.Context, e domain.AddressInfo) table.AddressInfo {
	return table.AddressInfo{
		StateDoc:     e.State.DocID(ctx).String(),
		ProvinceDoc:  e.Province.DocID(ctx).String(),
		DistrictDoc:  e.District.DocID(ctx).String(),
		AddressLine1: e.AddressLine1,
		DocumentID:   string(e.DocID(ctx)),
		Latitude:     e.Location[1],
		Longitude:    e.Location[0],
		ZipCode:      e.ZipCode,
		TimeZone:     e.TimeZone,
		TenantID:     sharedcontext.TenantIDFromContext(ctx),
	}
}
