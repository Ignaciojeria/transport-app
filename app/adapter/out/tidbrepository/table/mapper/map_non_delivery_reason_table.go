package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapNonDeliveryReasonTable(ctx context.Context, e domain.NonDeliveryReason) table.NonDeliveryReason {
	return table.NonDeliveryReason{
		ReferenceID: e.ReferenceID,
		Reason:      e.Reason,
		Details:     e.Details,
		DocumentID:  string(e.DocID(ctx)),
		TenantID:    sharedcontext.TenantIDFromContext(ctx),
	}
}
