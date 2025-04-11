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
		tidbrepository.NewUpsertNodeType,
		tidbrepository.NewUpsertContact,
		tidbrepository.NewUpsertAddressInfo,
		tidbrepository.NewUpsertNodeInfo,
	)
}

func NewUpsertNode(
	upsertNodeType tidbrepository.UpsertNodeType,
	upsertContact tidbrepository.UpsertContact,
	upsertAddressInfo tidbrepository.UpsertAddressInfo,
	upsertNodeInfo tidbrepository.UpsertNodeInfo,
) UpsertNode {
	return func(ctx context.Context, nodeInfo domain.NodeInfo) error {
		err := upsertNodeType(ctx, nodeInfo.NodeType)
		if err != nil {
			return err
		}
		err = upsertContact(ctx, nodeInfo.AddressInfo.Contact)
		if err != nil {
			return err
		}
		err = upsertAddressInfo(ctx, nodeInfo.AddressInfo)
		if err != nil {
			return err
		}
		return upsertNodeInfo(ctx, nodeInfo)
	}
}
