package usecase

import (
	"context"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchPlan func(ctx context.Context, referenceID string) (domain.Plan, error)

func init() {
	ioc.Registry(NewSearchPlan)
}

func NewSearchPlan() SearchPlan {
	return func(ctx context.Context, referenceID string) (domain.Plan, error) {
		return domain.Plan{}, nil
	}
}
