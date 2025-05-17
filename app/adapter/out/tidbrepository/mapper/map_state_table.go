package mapper

import (
	"context"

	"github.com/your-project/domain"
	"github.com/your-project/sharedcontext"
	"github.com/your-project/table"
)

func MapStateTable(ctx context.Context, state domain.State) table.State {
	return table.State{
		Name:       string(state),
		DocumentID: state.DocID(ctx).String(),
		TenantID:   sharedcontext.GetTenantID(ctx).String(),
	}
}
