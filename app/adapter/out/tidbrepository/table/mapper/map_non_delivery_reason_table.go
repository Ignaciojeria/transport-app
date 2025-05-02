package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapNonDeliveryReasonTable(ctx context.Context, e domain.NonDeliveryReason) table.NonDeliveryReason {
	return table.NonDeliveryReason{
		ReferenceID: e.ReferenceID,
		Reason:      e.Reason,
		DocumentID:  string(e.DocID(ctx)),
	}
}
