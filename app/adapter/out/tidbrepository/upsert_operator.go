package tidbrepository

import (
	"context"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertOperator func(context.Context, domain.Operator) (domain.Operator, error)

func init() {
	ioc.Registry(NewUpsertOperator)
}
func NewUpsertOperator() UpsertOperator {
	return func(ctx context.Context, o domain.Operator) (domain.Operator, error) {
		return domain.Operator{}, nil
	}
}
