package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapStateTable(ctx context.Context, s domain.State) table.State {
	return table.State{
		Name:       s.String(),
		DocumentID: s.DocID(ctx).String(),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
	}
}
