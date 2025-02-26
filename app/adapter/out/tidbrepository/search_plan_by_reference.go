package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"
)

type SearchPlanByReference func(context.Context, string) (domain.Plan, error)

func NewSearchPlanByReference(conn tidb.TIDBConnection) SearchPlanByReference {
	return func(ctx context.Context, s string) (domain.Plan, error) {

		return domain.Plan{}, nil
	}
}
