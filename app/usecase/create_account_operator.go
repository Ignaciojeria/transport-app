package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateAccountOperator func(ctx context.Context, input domain.Operator) (domain.Operator, error)

func init() {
	ioc.Registry(
		NewCreateAccountOperator,
		tidbrepository.NewUpsertContact,
		tidbrepository.NewUpsertNodeInfo,
		tidbrepository.NewUpsertOperator)
}

func NewCreateAccountOperator(
	upsertContact tidbrepository.UpsertContact,
	upsertNodeInfo tidbrepository.UpsertNodeInfo,
	upsertOperator tidbrepository.UpsertOperator,
) CreateAccountOperator {
	return func(ctx context.Context, input domain.Operator) (domain.Operator, error) {
		err := upsertContact(ctx, input.Contact)
		if err != nil {
			return domain.Operator{}, err
		}
		err = upsertNodeInfo(ctx, input.OriginNode)
		if err != nil {
			return domain.Operator{}, err
		}

		//input.OriginNode = nodeInfo
		return upsertOperator(ctx, input)
	}
}
