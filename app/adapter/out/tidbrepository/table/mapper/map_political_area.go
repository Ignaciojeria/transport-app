package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapPoliticalArea(ctx context.Context, pa domain.PoliticalArea) table.PoliticalArea {
	return table.PoliticalArea{
		DocumentID:      string(pa.DocID(ctx)),
		TenantID:        sharedcontext.TenantIDFromContext(ctx),
		Code:            pa.Code,
		AdminAreaLevel1: pa.AdminAreaLevel1,
		AdminAreaLevel2: pa.AdminAreaLevel2,
		AdminAreaLevel3: pa.AdminAreaLevel3,
		AdminAreaLevel4: pa.AdminAreaLevel4,
		TimeZone:        pa.TimeZone,
	}
}
