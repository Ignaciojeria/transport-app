package tidbrepository

import (
	"context"
	"transport-app/app/domain"
)

type UpsertOperator func(context.Context, domain.Operator) (domain.Operator, error)

func NewUpsertOperator() UpsertOperator {
	return func(ctx context.Context, o domain.Operator) (domain.Operator, error) {
		return domain.Operator{}, nil
	}
}
