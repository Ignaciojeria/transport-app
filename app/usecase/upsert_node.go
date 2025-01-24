package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertNode func(context.Context, domain.NodeInfo) error

func init() {
	ioc.Registry(
		NewUpsertNode,
		tidbrepository.NewUpsertContact,
		tidbrepository.NewUpsertOperator,
		tidbrepository.NewUpsertAddressInfo,
		tidbrepository.NewUpsertNodeInfo,
	)
}

func NewUpsertNode(
	upsertContact tidbrepository.UpsertContact,
	upsertOperator tidbrepository.UpsertOperator,
	upsertAddressInfo tidbrepository.UpsertAddressInfo,
	upsertNodeInfo tidbrepository.UpsertNodeInfo,
) UpsertNode {
	return func(ctx context.Context, origin domain.NodeInfo) error {
		/*
			o, err := query(ctx, origin)
			if err != nil {
				return err
			}
			o.UpdateIfChanged(origin)
			if err := upsert(ctx, o); err != nil {
				return err
			}
		*/
		return nil
	}
}
