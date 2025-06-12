package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapAddressInfoTable(ctx context.Context, e domain.AddressInfo) table.AddressInfo {
	return table.AddressInfo{
		PoliticalAreaDoc:        e.PoliticalArea.DocID(ctx).String(),
		AddressLine1:            e.AddressLine1,
		AddressLine2:            e.AddressLine2,
		DocumentID:              string(e.DocID(ctx)),
		Latitude:                e.Coordinates.Point.Lat(),
		Longitude:               e.Coordinates.Point.Lon(),
		ZipCode:                 e.ZipCode,
		TenantID:                sharedcontext.TenantIDFromContext(ctx),
		CoordinateSource:        e.Coordinates.Source,
		CoordinateConfidence:    e.Coordinates.Confidence.Level,
		CoordinateMessage:       e.Coordinates.Confidence.Message,
		CoordinateReason:        e.Coordinates.Confidence.Reason,
		PoliticalAreaConfidence: e.PoliticalArea.Confidence.Level,
		PoliticalAreaMessage:    e.PoliticalArea.Confidence.Message,
		PoliticalAreaReason:     e.PoliticalArea.Confidence.Reason,
	}
}
