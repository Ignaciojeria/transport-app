package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
)

type SearchPlanByReference func(context.Context, string) (domain.Plan, error)

func NewSearchPlanByReference(conn database.ConnectionFactory) SearchPlanByReference {
	return func(ctx context.Context, s string) (domain.Plan, error) {

		return domain.Plan{}, nil
	}
}
