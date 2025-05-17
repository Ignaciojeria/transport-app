package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapPlanHeaders(ctx context.Context, h domain.Headers) table.PlanHeaders {
	return table.PlanHeaders{
		DocumentID: string(h.DocID(ctx)),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
		Commerce:   h.Commerce,
		Consumer:   h.Consumer,
		Channel:    h.Channel,
	}
}
