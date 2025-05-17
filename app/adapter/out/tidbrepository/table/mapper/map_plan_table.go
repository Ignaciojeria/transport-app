package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapPlanToTable(ctx context.Context, p domain.Plan) table.Plan {
	return table.Plan{
		ReferenceID:    p.ReferenceID,
		DocumentID:     string(p.DocID(ctx)),
		TenantID:       sharedcontext.TenantIDFromContext(ctx),
		PlannedDate:    p.PlannedDate,
		PlanHeadersDoc: p.Headers.DocID(ctx).String(),
	}
}
