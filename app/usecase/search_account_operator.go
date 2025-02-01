package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchAccountOperator func(context.Context, domain.Operator) (domain.Operator, error)

func init() {
	ioc.Registry(
		NewSearchAccountOperator,
		tidbrepository.NewSearchOperatorByEmail)
}

func NewSearchAccountOperator(
	search tidbrepository.SearchOperatorByEmail) SearchAccountOperator {
	return func(ctx context.Context, input domain.Operator) (domain.Operator, error) {
		return search(ctx, input)
	}
}
