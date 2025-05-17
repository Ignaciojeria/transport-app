package mapper

import (
	"context"

	"github.com/your-project/domain"
	"github.com/your-project/sharedcontext"
	"github.com/your-project/table"
)

func MapProvinceTable(ctx context.Context, province domain.Province) table.Province {
	return table.Province{
		Name:       string(province),
		DocumentID: province.DocID(ctx).String(),
		TenantID:   sharedcontext.TenantIDFromContext(ctx).String(),
	}
}
