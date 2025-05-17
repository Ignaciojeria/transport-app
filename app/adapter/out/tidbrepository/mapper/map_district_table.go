package mapper

import (
	"context"

	"github.com/your-project/domain"
	"github.com/your-project/sharedcontext"
	"github.com/your-project/table"
)

func MapDistrictTable(ctx context.Context, district domain.District) table.District {
	return table.District{
		Name:       string(district),
		DocumentID: district.DocID(ctx).String(),
		TenantID:   sharedcontext.TenantIDFromContext(ctx).String(),
	}
}
