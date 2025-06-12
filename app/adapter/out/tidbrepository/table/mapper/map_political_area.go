package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapPoliticalArea(ctx context.Context, pa domain.PoliticalArea) table.PoliticalArea {
	return table.PoliticalArea{
		DocumentID: string(pa.DocID(ctx)),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
		Code:       pa.Code,
		State:      pa.State,
		Province:   pa.Province,
		District:   pa.District,
		TimeZone:   pa.TimeZone,
	}
}
