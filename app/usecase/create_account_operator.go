package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateAccountOperator func(ctx context.Context, input domain.Operator) (domain.Operator, error)

func init() {
	ioc.Registry(NewCreateAccountOperator, tidbrepository.NewUpsertContact)
}

func NewCreateAccountOperator(upsertContact tidbrepository.UpsertContact) CreateAccountOperator {
	return func(ctx context.Context, input domain.Operator) (domain.Operator, error) {

		return input, nil
	}
}
